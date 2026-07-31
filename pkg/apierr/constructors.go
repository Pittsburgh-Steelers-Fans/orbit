package apierr

// NotFound returns an API error for missing resources.
func NotFound(message string) *APIError {
	return &APIError{Code: "not_found", Message: message}
}

// Unauthorized returns an API error for missing or invalid credentials.
func Unauthorized(message string) *APIError {
	return &APIError{Code: "unauthorized", Message: message}
}

// BadRequest returns an API error for invalid client input.
func BadRequest(message string, details map[string]any) *APIError {
	return &APIError{Code: "bad_request", Message: message, Details: details}
}

// Internal returns an API error for unexpected server failures.
func Internal(message string) *APIError {
	return &APIError{Code: "internal_error", Message: message}
}
