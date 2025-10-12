package api_service

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	ChatID      int64   `json:"chat_id" binding:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Role        string  `json:"role" binding:"required"` // "client" or "professional"
}

// ProfessionalSignInRequest represents a professional sign-in request
type ProfessionalSignInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	ChatID   int64  `json:"chat_id" binding:"required"`
}

// CreateAppointmentRequest represents a request to create an appointment
type CreateAppointmentRequest struct {
	ClientID       string `json:"client_id" binding:"required"`
	ProfessionalID string `json:"professional_id" binding:"required"`
	StartTime      string `json:"start_time" binding:"required"`
	EndTime        string `json:"end_time" binding:"required"`
}

// CancelAppointmentRequest represents a request to cancel an appointment
type CancelAppointmentRequest struct {
	CancellationReason string `json:"cancellation_reason" binding:"required"`
}

// ConfirmAppointmentRequest represents a request to confirm an appointment
type ConfirmAppointmentRequest struct {
	// No additional fields needed for confirmation
}

// CreateUnavailableAppointmentRequest represents a request to create an unavailable appointment
type CreateUnavailableAppointmentRequest struct {
	ProfessionalID string `json:"professional_id" binding:"required"`
	StartAt        string `json:"start_at" binding:"required"`
	EndAt          string `json:"end_at" binding:"required"`
	Description    string `json:"description" binding:"required"`
}
