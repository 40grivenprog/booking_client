package api_service

import "fmt"

// ErrorResponse represents the error response from the API
type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

// APIError represents a structured error from the API
type APIError struct {
	StatusCode int
	ErrorType  string
	Message    string
	RequestID  string
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("%s (request_id: %s)", e.Message, e.RequestID)
	}
	return e.Message
}

// FormatUserMessage formats the error message for end users
func (e *APIError) FormatUserMessage() string {
	msg := e.Message
	if e.RequestID != "" {
		msg += fmt.Sprintf("\n\nPlease contact support (maksimfilipenka122@gmail.com) with request_id: %s", e.RequestID)
	}
	return msg
}
