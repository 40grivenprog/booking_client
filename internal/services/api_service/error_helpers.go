package api_service

// FormatErrorForUser formats an API error for display to end users
// Returns a user-friendly error message with contact information if request_id is available
func FormatErrorForUser(err error) string {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.FormatUserMessage()
	}
	return err.Error()
}

// GetRequestID extracts the request ID from an API error if available
func GetRequestID(err error) string {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.RequestID
	}
	return ""
}

// IsAPIError checks if the error is an APIError
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}
