package logger_middleware

import (
	"fmt"
	"github.com/idsulik/url-shortener/internal/logger"
	"net/http"
)

func New(logger *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("New request: %s %s", r.Method, r.URL.Path))
			defer func() {
				logger.Info(fmt.Sprintf("Request completed: %s %s", r.Method, r.URL.Path))
			}()
			next.ServeHTTP(w, r)
		})
	}
}
