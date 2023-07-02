package redirect_test

import (
	"github.com/go-chi/chi/v5"
	"github.com/idsulik/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/idsulik/url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"github.com/idsulik/url-shortener/internal/logger"
	"github.com/idsulik/url-shortener/internal/storage"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		statusCode int
	}{
		{
			name:       "Success",
			alias:      "test_alias",
			statusCode: http.StatusFound,
		},
		{
			name:       "Empty alias",
			alias:      "",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Alias not found",
			alias:      "not_exists_alias",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewUrlGetter(t)
			handler := redirect.New(logger.New("test"), urlGetterMock)
			r := chi.NewRouter()
			r.Get("/{alias}", handler)

			if tc.alias != "" {
				if tc.alias == "not_exists_alias" {
					urlGetterMock.
						On("GetUrl", tc.alias).
						Return("", storage.ErrUrlNotFound).
						Once()
				} else {
					urlGetterMock.
						On("GetUrl", tc.alias).
						Return("https://google.com", nil).
						Once()
				}
			}

			req, err := http.NewRequest(http.MethodGet, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, tc.statusCode, rr.Code)
		})
	}
}
