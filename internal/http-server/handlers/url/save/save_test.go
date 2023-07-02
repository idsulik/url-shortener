package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/idsulik/url-shortener/internal/http-server/handlers/url/save"
	"github.com/idsulik/url-shortener/internal/http-server/handlers/url/save/mocks"
	"github.com/idsulik/url-shortener/internal/logger"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		url        string
		statusCode int
		respError  string
		mockError  error
	}{
		{
			name:       "Success",
			alias:      "test_alias",
			url:        "https://google.com",
			statusCode: http.StatusCreated,
		},
		{
			name:       "Empty alias",
			alias:      "",
			url:        "https://google.com",
			statusCode: http.StatusCreated,
		},
		{
			name:       "Empty URL",
			alias:      "some_alias",
			url:        "",
			statusCode: http.StatusBadRequest,
			respError:  "Key: 'Request.Url' Error:Field validation for 'Url' failed on the 'required' tag\n",
		},
		{
			name:       "Invalid URL",
			alias:      "some_alias",
			url:        "some invalid URL",
			statusCode: http.StatusBadRequest,
			respError:  "Key: 'Request.Url' Error:Field validation for 'Url' failed on the 'url' tag\n",
		},
		{
			name:       "SaveURL Error",
			alias:      "test_alias",
			url:        "https://google.com",
			statusCode: http.StatusBadRequest,
			respError:  "failed to add url",
			mockError:  errors.New("failed to add url"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewUrlSaver(t)
			aliasGeneratorMock := mocks.NewAliasGenerator(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveUrl", mock.AnythingOfType("string"), tc.url).
					Return(int64(1), tc.mockError).
					Once()
			}

			if tc.alias == "" {
				aliasGeneratorMock.
					On("NewAlias", mock.AnythingOfType("int")).
					Return("myAlias").
					Once()
			}

			handler := save.New(logger.New("test"), aliasGeneratorMock, urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.statusCode, rr.Code)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
