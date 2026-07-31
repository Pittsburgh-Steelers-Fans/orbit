package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLimiterAllowsBurstThenBlocksThenRefills(t *testing.T) {
	current := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
	limiter := newLimiter(2, time.Second, func() time.Time { return current })

	assert.True(t, limiter.Allow("user-1"))
	assert.True(t, limiter.Allow("user-1"))
	assert.False(t, limiter.Allow("user-1"))

	current = current.Add(time.Second)
	assert.True(t, limiter.Allow("user-1"))
	assert.False(t, limiter.Allow("user-1"))
}

func TestLimiterTracksKeysIndependently(t *testing.T) {
	current := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
	limiter := newLimiter(1, time.Minute, func() time.Time { return current })

	assert.True(t, limiter.Allow("user-1"))
	assert.False(t, limiter.Allow("user-1"))
	assert.True(t, limiter.Allow("user-2"))
}

func TestMiddlewareReturnsTooManyRequestsWhenDenied(t *testing.T) {
	limiter := newLimiter(1, time.Minute, time.Now)
	middleware := Middleware(limiter, func(r *http.Request) string { return "client" })
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/", nil))
	require.Equal(t, http.StatusNoContent, first.Code)

	second := httptest.NewRecorder()
	handler.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/", nil))
	require.Equal(t, http.StatusTooManyRequests, second.Code)
	assert.JSONEq(t, `{"error":"rate limit exceeded"}`, second.Body.String())
}
