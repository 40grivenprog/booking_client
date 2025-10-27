package schemas

import "booking_client/internal/models"

// Client represents a client in the response
type Client struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	ChatID      int64  `json:"chat_id,omitempty"`
}

// ClientRegisterResponse represents the response for client registration
type ClientRegisterResponse struct {
	ID          string  `json:"id"`
	ChatID      *int64  `json:"chat_id,omitempty"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Role        string  `json:"role"`
}

// ClientAppointment represents an appointment with professional details in client context
type ClientAppointment struct {
	ID           string                         `json:"id"`
	Type         string                         `json:"type"`
	StartTime    string                         `json:"start_time"`
	EndTime      string                         `json:"end_time"`
	Status       string                         `json:"status"`
	Description  string                         `json:"description,omitempty"`
	CreatedAt    string                         `json:"created_at"`
	UpdatedAt    string                         `json:"updated_at"`
	Professional *ClientAppointmentProfessional `json:"professional,omitempty"`
}

// ClientAppointmentProfessional represents professional details in appointment context
type ClientAppointmentProfessional struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	ChatID      *int64  `json:"chat_id,omitempty"`
}

// ClientAppointmentClient represents client details in appointment context
type ClientAppointmentClient struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	ChatID      *int64  `json:"chat_id,omitempty"`
}

// GetClientAppointmentsResponse represents the response for getting client appointments
type GetClientAppointmentsResponse struct {
	Appointments []ClientAppointment `json:"appointments"`
}

// CancelClientAppointmentResponse represents the response after cancelling an appointment by client
type CancelClientAppointmentResponse struct {
	Appointment  models.CancelledAppointment   `json:"appointment"`
	Client       ClientAppointmentClient       `json:"client"`
	Professional ClientAppointmentProfessional `json:"professional"`
}
