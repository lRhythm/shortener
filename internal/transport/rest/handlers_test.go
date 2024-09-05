package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/lRhythm/shortener/internal/config"
	"github.com/lRhythm/shortener/internal/logs"
	"github.com/lRhythm/shortener/internal/service"
	"github.com/lRhythm/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestApiCreateHandler(t *testing.T) {
	cfg, _ := config.New()
	db, _ := storage.NewMemory(cfg.File())
	s, _ := New(
		logs.New(),
		cfg,
		service.New(
			service.WithStorage(db),
		),
	)
	f := fiber.New()
	f.Post("/api/shorten", s.apiCreateHandler)
	tests := []struct {
		name        string
		route       string
		req         createRequest
		status      int
		contentType string
	}{
		{
			name:        "1. 201 - success",
			route:       "/api/shorten",
			req:         createRequest{URL: "https://ya.ru"},
			status:      fiber.StatusCreated,
			contentType: "application/json",
		},
		{
			name:   "2. 400 - error: bad url in body.url",
			route:  "/api/shorten",
			req:    createRequest{URL: "WRONG"},
			status: fiber.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			j, err := json.Marshal(test.req)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, test.route, strings.NewReader(string(j)))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := f.Test(req, -1)
			assert.Equal(t, test.status, resp.StatusCode)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)
			r := string(b)
			if test.status == fiber.StatusCreated {
				assert.Equal(t, len(r) > 0, true)
				assert.Equal(t, resp.Header.Get("Content-Type"), test.contentType)
			}
			if test.status == fiber.StatusBadRequest {
				assert.Equal(t, len(r) == 0, true)
			}
		})
	}
}

func TestCreateHandler(t *testing.T) {
	cfg, _ := config.New()
	db, _ := storage.NewMemory(cfg.File())
	s, _ := New(
		logs.New(),
		cfg,
		service.New(
			service.WithStorage(db),
		),
	)
	f := fiber.New()
	f.Post("/", s.createHandler)
	tests := []struct {
		name   string
		route  string
		body   string
		status int
	}{
		{
			name:   "1. 201 - success",
			route:  "/",
			body:   "https://ya.ru",
			status: fiber.StatusCreated,
		},
		{
			name:   "2. 400 - error: bad url in body",
			route:  "/",
			body:   "WRONG",
			status: fiber.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, test.route, strings.NewReader(test.body))
			resp, _ := f.Test(req, -1)
			assert.Equal(t, test.status, resp.StatusCode)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)
			r := string(b)
			if test.status == fiber.StatusCreated {
				assert.Equal(t, len(r) > 0, true)
			}
			if test.status == fiber.StatusBadRequest {
				assert.Equal(t, len(r) == 0, true)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	cfg, _ := config.New()
	db, _ := storage.NewMemory(cfg.File())
	logic := service.New(
		service.WithStorage(db),
	)
	s, _ := New(
		logs.New(),
		cfg,
		logic,
	)
	f := fiber.New()
	f.Get("/:id", s.getHandler)

	// Создание в хранилище сокращенного URL для дальнейшей проверки.
	ou := "https://ya.ru"
	su, _ := logic.CreateShortURL(ou, "http://localhost:8080")
	u, _ := url.Parse(su)

	tests := []struct {
		name     string
		route    string
		id       string
		status   int
		location string
	}{
		{
			name:     "1. 307 - success",
			route:    u.Path,
			status:   fiber.StatusTemporaryRedirect,
			location: ou,
		},
		{
			name:   "2. 400 - error: short url not found",
			route:  "/WRONG",
			status: fiber.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log.Println(test.route)
			req := httptest.NewRequest(http.MethodGet, test.route, nil)
			resp, _ := f.Test(req, -1)
			assert.Equal(t, test.status, resp.StatusCode)
			err := resp.Body.Close()
			require.NoError(t, err)
			if test.status == fiber.StatusTemporaryRedirect {
				assert.Equal(t, test.location, resp.Header.Get("Location"))
			}
		})
	}
}
