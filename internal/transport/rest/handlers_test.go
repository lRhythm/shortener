package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lRhythm/shortener/internal/config"
	"github.com/lRhythm/shortener/internal/logs"
	"github.com/lRhythm/shortener/internal/models"
	"github.com/lRhythm/shortener/internal/service"
	"github.com/lRhythm/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// [POST] /api/shorten
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
	f.Post("/api/shorten", s.registerMiddleware, s.apiCreateHandler)
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
			req:         createRequest{OriginalURL: "https://ya.ru"},
			status:      fiber.StatusCreated,
			contentType: "application/json",
		},
		{
			name:   "2. 400 - error: bad url in body.url",
			route:  "/api/shorten",
			req:    createRequest{OriginalURL: "WRONG"},
			status: fiber.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			j, err := json.Marshal(test.req)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, test.route, strings.NewReader(string(j)))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{
				Name:  cookieUserID,
				Value: uuid.NewString(),
			})
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

// [POST] /api/shorten/batch
func TestApiCreateBatchHandler(t *testing.T) {
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
	f.Post("/api/shorten/batch", s.registerMiddleware, s.apiCreateBatchHandler)
	tests := []struct {
		name        string
		route       string
		req         createItemsRequest
		respCorrID  []string
		status      int
		contentType string
	}{
		{
			name:  "1. 201 - success",
			route: "/api/shorten/batch",
			req: createItemsRequest{
				createItemRequest{
					OriginalURL:   "https://ya.ru",
					CorrelationID: "id1",
				},
				createItemRequest{
					OriginalURL:   "https://yandex.ru",
					CorrelationID: "",
				},
				createItemRequest{
					OriginalURL:   "",
					CorrelationID: "id3",
				},
				createItemRequest{
					OriginalURL:   "",
					CorrelationID: "",
				},
				createItemRequest{
					OriginalURL:   "https://yadi.sk",
					CorrelationID: "id5",
				},
			},
			respCorrID:  []string{"id1", "id5"},
			status:      fiber.StatusCreated,
			contentType: "application/json",
		},
		{
			name:  "2. 400 - error: bad url in body.*.original_url",
			route: "/api/shorten/batch",
			req: createItemsRequest{
				createItemRequest{
					OriginalURL:   "WRONG",
					CorrelationID: "id1",
				},
			},
			status: fiber.StatusBadRequest,
		},
		{
			name:  "3. 400 - error: empty body.*",
			route: "/api/shorten/batch",
			req: createItemsRequest{
				createItemRequest{
					OriginalURL:   "",
					CorrelationID: "",
				},
			},
			status: fiber.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			j, err := json.Marshal(test.req)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, test.route, strings.NewReader(string(j)))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{
				Name:  cookieUserID,
				Value: uuid.NewString(),
			})
			resp, _ := f.Test(req, -1)
			assert.Equal(t, test.status, resp.StatusCode)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)
			r := string(b)
			if test.status == fiber.StatusCreated {
				assert.True(t, len(r) > 0)
				assert.Equal(t, resp.Header.Get("Content-Type"), test.contentType)
				var rb createItemsResponse
				err = json.Unmarshal(b, &rb)
				require.NoError(t, err)
				for _, id := range rb {
					assert.Contains(t, test.respCorrID, id.CorrelationID)
				}
			}
			if test.status == fiber.StatusBadRequest {
				assert.True(t, len(r) == 0)
			}
		})
	}
}

// [POST] /
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
	f.Post("/", s.registerMiddleware, s.createHandler)
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
			req.AddCookie(&http.Cookie{
				Name:  cookieUserID,
				Value: uuid.NewString(),
			})
			resp, _ := f.Test(req, -1)
			assert.Equal(t, test.status, resp.StatusCode)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)
			r := string(b)
			if test.status == fiber.StatusCreated {
				assert.True(t, len(r) > 0)
			}
			if test.status == fiber.StatusBadRequest {
				assert.True(t, len(r) == 0)
			}
		})
	}
}

// [GET] /:id
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
	su, _ := logic.CreateShortURL(ou, "http://localhost:8080", uuid.NewString())
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

// [GET] /api/user/urls
func TestApiUserUrlsGetHandler(t *testing.T) {
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
	f.Get("/api/user/urls", s.authenticateMiddleware, s.apiUserUrlsGetHandler)

	// Создание в хранилище сокращенного URL для дальнейшей проверки.
	uid := uuid.NewString()
	a := "http://example.com" // См. https://pkg.go.dev/net/http/httptest#NewRequest
	ou := "https://ya.ru"
	su, _ := logic.CreateShortURL(ou, a, uid)

	tests := []struct {
		name   string
		status int
		userId string
		resp   models.Rows
	}{
		{
			name:   "1. 200 - success",
			status: fiber.StatusOK,
			userId: uid,
			resp: models.Rows{models.Row{
				ShortURL:    su,
				OriginalURL: ou,
			}},
		},
		{
			name:   "1. 204 - no content",
			status: fiber.StatusNoContent,
			userId: uuid.NewString(),
		},
		{
			name:   "1. 401 - unauthorized",
			status: fiber.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			req.AddCookie(&http.Cookie{
				Name:  cookieUserID,
				Value: test.userId,
			})
			resp, _ := f.Test(req, -1)
			assert.Equal(t, test.status, resp.StatusCode)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if test.status == fiber.StatusOK {
				var rb models.Rows
				err = json.Unmarshal(b, &rb)
				require.NoError(t, err)
				assert.Equal(t, test.resp, rb)
			}
			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}
