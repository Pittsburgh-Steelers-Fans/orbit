package ratelimit

import "net/http"

// KeyFunc derives a rate-limit key from an HTTP request.
type KeyFunc func(*http.Request) string

// Middleware rejects requests with HTTP 429 when the limiter denies the key.
func Middleware(limiter *Limiter, keyFunc KeyFunc) func(http.Handler) http.Handler {
	if keyFunc == nil {
		keyFunc = func(r *http.Request) string { return r.RemoteAddr }
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow(keyFunc(r)) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"rate limit exceeded"}` + "\n"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
