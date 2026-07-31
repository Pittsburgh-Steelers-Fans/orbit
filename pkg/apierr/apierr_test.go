package apierr

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIErrorError(t *testing.T) {
	err := &APIError{Code: "bad_request", Message: "name is required"}

	require.Equal(t, "bad_request: name is required", err.Error())
}

func TestWriteJSON(t *testing.T) {
	err := BadRequest("invalid input", map[string]any{"field": "name"})
	recorder := httptest.NewRecorder()

	err.WriteJSON(recorder, http.StatusBadRequest)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var payload APIError
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &payload))
	require.Equal(t, "bad_request", payload.Code)
	require.Equal(t, "invalid input", payload.Message)
	require.Equal(t, "name", payload.Details["field"])
}

func TestConstructors(t *testing.T) {
	require.Equal(t, "not_found", NotFound("missing").Code)
	require.Equal(t, "unauthorized", Unauthorized("login required").Code)
	require.Equal(t, "internal_error", Internal("try again later").Code)
}
