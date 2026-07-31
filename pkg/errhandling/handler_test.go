package errhandling

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRespondMapsTypedErrorToJSON(t *testing.T) {
	recorder := httptest.NewRecorder()

	Respond(recorder, New(http.StatusNotFound, "not_found", "task not found"))

	require.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	var body Error
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	assert.Equal(t, "not_found", body.Code)
	assert.Equal(t, "task not found", body.Message)
}

func TestRespondFallsBackToInternalServerError(t *testing.T) {
	recorder := httptest.NewRecorder()

	Respond(recorder, errors.New("database offline"))

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	var body Error
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	assert.Equal(t, "internal_error", body.Code)
	assert.Equal(t, "Internal Server Error", body.Message)
}

func TestRecoverConvertsPanicToJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	middleware := Recover(zerolog.Nop())
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	var body Error
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	assert.Equal(t, "internal_error", body.Code)
	assert.Equal(t, "internal server error", body.Message)
}
