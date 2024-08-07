package server

import (
	"github.com/lRhythm/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHandlerResolver - в каждом тесте проверяются создание и получение ссылки одновременно.
func TestHandlerResolver(t *testing.T) {
	s := New(storage.NewInMemory())
	type want struct {
		statusCodeCreate int
		statusCodeGet    int
		url              string
	}
	tests := []struct {
		name       string
		request    string // Адрес маршрута создания сокращённой ссылки.
		body       string // Тело запроса создания сокращённой ссылки.
		requestGet string // Адрес маршрута получения ссылки для подмены.
		want       want
	}{
		{
			name:    "1. success",
			request: "/",
			body:    "https://ya.ru",
			want: want{
				statusCodeCreate: 201,
				statusCodeGet:    307,
				url:              "https://ya.ru",
			},
		},
		{
			name:    "2. create error: url bad format in body",
			request: "/",
			body:    "WRONG",
			want: want{
				statusCodeCreate: 400,
			},
		},
		{
			name:       "3. get error: url not exists in storage",
			request:    "/",
			body:       "https://ya.ru",
			requestGet: "/WRONG",
			want: want{
				statusCodeCreate: 201,
				statusCodeGet:    400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создание сокращённой ссылки.
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			s.handlerResolver(w, request)
			result := w.Result()

			require.Equal(t, tt.want.statusCodeCreate, result.StatusCode)
			// Если ожидаемый статус состояния не 201, то остальные проверки пропускаются.
			if tt.want.statusCodeCreate != 201 {
				return
			}

			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(result.Body)
			b, e := io.ReadAll(result.Body)
			require.NoError(t, e)
			u := string(b)

			require.Greater(t, len(u), 0)

			// Подмена адреса получения ссылки.
			if len(tt.requestGet) > 0 {
				u = tt.requestGet
			}

			// Получение ссылки.
			request = httptest.NewRequest(http.MethodGet, u, nil)
			w = httptest.NewRecorder()
			s.handlerResolver(w, request)
			result = w.Result()

			assert.Equal(t, tt.want.statusCodeGet, result.StatusCode)
			assert.Equal(t, tt.want.url, result.Header.Get("Location"))
		})
	}
}
