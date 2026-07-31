package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueParseRoundTrip(t *testing.T) {
	token, err := Issue("secret", "user-123", time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := Parse("secret", token)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "user-123", claims.Subject)
	assert.Contains(t, claims.Roles, "user")
}

func TestAuthenticateRejectsMissingAndInvalidToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler := Authenticate("secret")(next)

	tests := []struct {
		name   string
		header string
	}{
		{name: "missing token"},
		{name: "invalid token", header: "Bearer not-a-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/private", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})
	}
}
