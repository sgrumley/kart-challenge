package middleware

import (
	"log/slog"
	"net/http"

	"github.com/sgrumley/kart-challenge/pkg/logger"
)

func AddLogger(log slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestLogger := log.With(
				"method", r.Method,
				"path", r.URL.Path,
			)

			ctx := logger.AddLoggerContext(r.Context(), requestLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
