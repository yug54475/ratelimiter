package middleware

import (
	"net/http"

	"ratelimiter/core"
)

// KeyFunc extracts a client identifier from the request (IP, user ID, API key, etc.)
type KeyFunc func(r *http.Request) string

// ByIP is a simple KeyFunc that identifies clients by their IP:port.
func ByIP(r *http.Request) string {
	return r.RemoteAddr
}

// RateLimit wraps any handler with rate limiting logic.
func RateLimit(limiter core.Limiter, keyFn KeyFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := keyFn(r)

			if !limiter.Allow(clientID) {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
