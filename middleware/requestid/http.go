package requestid

import (
	"net/http"

	"github.com/segmentio/ksuid"
)

// HTTPMiddleware creates an HTTP middleware that extracts or generates a request ID
// and adds it to the request context
func HTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := extractFromHTTP(r)
			if requestID != "" {
				w.Header().Set(Header, requestID)
			}

			ctx := WithContext(r.Context(), requestID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// extractFromHTTP extracts the request ID from HTTP headers or generates a new one
func extractFromHTTP(r *http.Request) string {
	if requestID := r.Header.Get(Header); requestID != "" {
		return requestID
	}

	return ksuid.New().String()
}
