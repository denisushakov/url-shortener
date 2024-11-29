package delete

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		respError    string
		expectedCode int
		mockError    error
	}{
		{
			name:         "Success",
			alias:        "alias-to-delete",
			expectedCode: http.StatusOK,
		},
		{
			name:         "URL Not Found",
			alias:        "non-existent-alias",
			respError:    "url not found",
			expectedCode: http.StatusNotFound,
			mockError:    storage.ErrURLNotFound,
		},
		{
			name:         "Internal Error",
			alias:        "alias-with-error",
			respError:    "internal error",
			expectedCode: http.StatusInternalServerError,
			mockError:    assert.AnError,
		},
		/*{
			name:         "Failure - Incorrect Expected Status",
			alias:        "alias-to-delete",
			expectedCode: http.StatusNotFound, // Неправильный ожидаемый код для успешного удаления
		},*/
	}
	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Delete("/url/{alias}", New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			req, err := http.NewRequest(http.MethodDelete, "/"+path.Join("url", tc.alias), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
			if tc.respError != "" {
				assert.Contains(t, rr.Body.String(), tc.respError)
			}
		})
	}
}
