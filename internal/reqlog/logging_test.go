package reqlog

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareCallsHandlerAndCapturesStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf)
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})
	handler := Middleware(logger)(next)

	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, buf.String(), `"method":"POST"`)
	assert.Contains(t, buf.String(), `"path":"/tasks"`)
	assert.Contains(t, buf.String(), `"status":201`)
}
