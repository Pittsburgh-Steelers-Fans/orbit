package errhandling

import (
	"encoding/json"
	"errors"
	"net/http"
)

const defaultErrorCode = "internal_error"

// Error is a typed application error that can be rendered as JSON.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

// Error returns the human-readable error message.
func (e Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Code
}

// Unwrap returns the underlying error.
func (e Error) Unwrap() error {
	return e.Err
}

// New creates a typed application error for HTTP responses.
func New(status int, code, message string) Error {
	return Error{Status: status, Code: code, Message: message}
}

// Wrap creates a typed application error that preserves an underlying cause.
func Wrap(status int, code, message string, err error) Error {
	return Error{Status: status, Code: code, Message: message, Err: err}
}

// Respond writes a JSON error response mapped from a typed error.
func Respond(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	payload := Error{Code: defaultErrorCode, Message: http.StatusText(http.StatusInternalServerError)}

	var appErr Error
	if errors.As(err, &appErr) {
		status = appErr.Status
		if status == 0 {
			status = http.StatusInternalServerError
		}
		payload.Code = appErr.Code
		if payload.Code == "" {
			payload.Code = defaultErrorCode
		}
		payload.Message = appErr.Message
		if payload.Message == "" {
			payload.Message = http.StatusText(status)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
