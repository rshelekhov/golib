package cors

import "net/http"

// Middleware creates middleware for handling CORS
func Middleware(origins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			// Check if the origin is allowed
			origin := r.Header.Get("Origin")
			for _, allowed := range origins {
				if allowed == "*" || allowed == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
