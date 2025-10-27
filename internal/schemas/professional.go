package schemas

import "booking_client/internal/models"

// Professional represents a professional in the response
type Professional struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	ChatID      int64  `json:"chat_id,omitempty"`
}

// CreateAppointmentResponse represents the response after creating an appointment
type CreateAppointmentResponse struct {
	Appointment  models.Appointment `json:"appointment"`
	Client       models.Client      `json:"client"`
	Professional Professional       `json:"professional"`
}

// ProfessionalClient represents a client for professional
type ProfessionalClient struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// PreviousAppointment represents a previous appointment
type PreviousAppointment struct {
	ID          string `json:"id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description"`
}

// GetPreviousAppointmentsByClientResponse represents the response for getting previous appointments
type GetPreviousAppointmentsByClientResponse struct {
	Appointments []PreviousAppointment `json:"appointments"`
}

// ProfessionalAvailabilityResponse represents the response for professional availability
type ProfessionalAvailabilityResponse struct {
	Date  string     `json:"date"`
	Slots []TimeSlot `json:"slots"`
}

// ProfessionalAppointment represents an appointment with client details in professional context
type ProfessionalAppointment struct {
	ID          string                         `json:"id"`
	Type        string                         `json:"type"`
	StartTime   string                         `json:"start_time"`
	EndTime     string                         `json:"end_time"`
	Status      string                         `json:"status"`
	Description string                         `json:"description,omitempty"`
	CreatedAt   string                         `json:"created_at"`
	UpdatedAt   string                         `json:"updated_at"`
	Client      *ProfessionalAppointmentClient `json:"client,omitempty"`
}

// ProfessionalAppointmentClient represents client details in appointment context
type ProfessionalAppointmentClient struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	ChatID      *int64  `json:"chat_id,omitempty"`
}

// GetProfessionalAppointmentsResponse represents the response for getting professional appointments
type GetProfessionalAppointmentsResponse struct {
	Appointments []ProfessionalAppointment `json:"appointments"`
}

// GetProfessionalAppointmentDatesResponse represents the response for getting appointment dates
type GetProfessionalAppointmentDatesResponse struct {
	Month string   `json:"month"`
	Dates []string `json:"dates"`
}

// ConfirmProfessionalAppointmentResponse represents the response after confirming an appointment by professional
type ConfirmProfessionalAppointmentResponse struct {
	Appointment  ProfessionalAppointment       `json:"appointment"`
	Client       ProfessionalAppointmentClient `json:"client"`
	Professional Professional                  `json:"professional"`
}

// CancelProfessionalAppointmentResponse represents the response after cancelling an appointment by professional
type CancelProfessionalAppointmentResponse struct {
	Appointment  CancelledAppointment          `json:"appointment"`
	Client       ProfessionalAppointmentClient `json:"client"`
	Professional Professional                  `json:"professional"`
}

// UnavailableAppointment represents an unavailable appointment
type UnavailableAppointment struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Status      string `json:"status"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateUnavailableAppointmentResponse represents the response after creating an unavailable appointment
type CreateUnavailableAppointmentResponse struct {
	Appointment UnavailableAppointment `json:"appointment"`
}

// TimetableAppointment represents an appointment in the timetable
type TimetableAppointment struct {
	ID          string `json:"id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description"`
}

// GetProfessionalTimetableResponse represents the response for getting professional timetable
type GetProfessionalTimetableResponse struct {
	Date         string                 `json:"date"`
	Appointments []TimetableAppointment `json:"appointments"`
}

// GetProfessionalsResponse represents the response for getting all professionals
type GetProfessionalsResponse struct {
	Professionals []models.User `json:"professionals"`
}
