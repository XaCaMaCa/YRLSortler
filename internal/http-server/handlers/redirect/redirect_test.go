package redirect_test

import (
	"context"
	"encoding/json"
	"errors"
	"libmods/internal/http-server/handlers/redirect"
	"libmods/internal/lib/api/response"
	"libmods/internal/lib/logger/handlers/slogdiscard"
	"libmods/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// URLGetterMock is a mock implementation of URLGetter interface
type URLGetterMock struct {
	mock.Mock
}

func (m *URLGetterMock) GetURL(alias string) (string, error) {
	args := m.Called(alias)
	return args.String(0), args.Error(1)
}

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		url        string
		mockError  error
		statusCode int
		respError  string
		redirect   bool
		location   string
	}{
		{
			name:       "Success - redirect to URL",
			alias:      "test_alias",
			url:        "https://google.com",
			statusCode: http.StatusFound,
			redirect:   true,
			location:   "https://google.com",
		},
		{
			name:       "URL not found",
			alias:      "nonexistent",
			mockError:  storage.ErrURLNotFound,
			statusCode: http.StatusOK,
			respError:  "not found",
		},
		{
			name:       "Internal error",
			alias:      "test_alias",
			mockError:  errors.New("database error"),
			statusCode: http.StatusOK,
			respError:  "internal error",
		},
		{
			name:       "Empty alias - invalid request",
			alias:      "",
			statusCode: http.StatusOK,
			respError:  "invalid request",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := new(URLGetterMock)

			// Настраиваем mock только если алиас не пустой
			if tc.alias != "" {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).
					Once()
			}

			handler := redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock)

			var req *http.Request
			var err error

			if tc.alias == "" {
				// Для пустого алиаса создаем запрос и устанавливаем контекст с пустым параметром
				req, err = http.NewRequest(http.MethodGet, "/", nil)
				require.NoError(t, err)

				// Создаем контекст с пустым алиасом для тестирования
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("alias", "")
				ctx := req.Context()
				ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
				req = req.WithContext(ctx)
			} else {
				// Создаем chi router для правильной работы URL параметров
				router := chi.NewRouter()
				router.Get("/{alias}", handler)

				req, err = http.NewRequest(http.MethodGet, "/"+tc.alias, nil)
				require.NoError(t, err)

				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				// Проверяем результаты
				require.Equal(t, tc.statusCode, rr.Code)

				if tc.redirect {
					// Проверяем редирект
					require.Equal(t, tc.location, rr.Header().Get("Location"))
				} else {
					// Проверяем JSON ответ с ошибкой
					body := rr.Body.String()

					var resp response.Response

					require.NoError(t, json.Unmarshal([]byte(body), &resp))
					require.Equal(t, response.StatusError, resp.Status)
					require.Equal(t, tc.respError, resp.Error)
				}

				urlGetterMock.AssertExpectations(t)
				return
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Проверяем результаты для пустого алиаса
			require.Equal(t, tc.statusCode, rr.Code)

			if tc.redirect {
				// Проверяем редирект
				require.Equal(t, tc.location, rr.Header().Get("Location"))
			} else {
				// Проверяем JSON ответ с ошибкой
				body := rr.Body.String()

				var resp response.Response

				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				require.Equal(t, response.StatusError, resp.Status)
				require.Equal(t, tc.respError, resp.Error)
			}

			urlGetterMock.AssertExpectations(t)
		})
	}
}
