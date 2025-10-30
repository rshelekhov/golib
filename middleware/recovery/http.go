package recovery

import (
	"log/slog"
	"net/http"
)

// Middleware creates middleware for recovering from panics
func Middleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the error
					logger.Error("http server panic recovered",
						"error", err,
						"path", r.URL.Path,
						"method", r.Method,
					)

					// Return 500 error
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
