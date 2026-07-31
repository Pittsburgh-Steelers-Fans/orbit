package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type claimsContextKey struct{}

// Authenticate validates Bearer tokens and stores their claims on the request context.
func Authenticate(secret string) func(http.Handler) http.Handler {
	if secret == "" {
		secret = os.Getenv("JWT_SECRET")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				writeAuthError(w, http.StatusUnauthorized, "auth: missing bearer token")
				return
			}

			claims, err := Parse(secret, strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
			if err != nil {
				writeAuthError(w, http.StatusUnauthorized, "auth: invalid bearer token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ClaimsFromContext returns JWT claims previously added by Authenticate.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(*Claims)
	return claims, ok
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
