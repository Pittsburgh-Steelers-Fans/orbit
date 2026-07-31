package apierr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError is a typed error payload returned by HTTP API handlers.
type APIError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// Error returns the stable code and human-readable message for the API error.
func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message == "" {
		return e.Code
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WriteJSON writes the API error as a JSON response with the provided status.
func (e *APIError) WriteJSON(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(e)
}
