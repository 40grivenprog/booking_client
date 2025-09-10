package models

// User represents a user in the system
type User struct {
	ID                             string  `json:"id"`
	ChatID                         *int64  `json:"chat_id,omitempty"`
	Username                       string  `json:"username"`
	FirstName                      string  `json:"first_name"`
	LastName                       string  `json:"last_name"`
	Role                           string  `json:"role"` // "client" or "professional"
	PhoneNumber                    *string `json:"phone_number,omitempty"`
	State                          string  `json:"state,omitempty"`                            // Bot interaction state
	SelectedProfessionalID         string  `json:"selected_professional_id,omitempty"`         // Temporary storage for appointment booking
	SelectedDate                   string  `json:"selected_date,omitempty"`                    // Temporary storage for selected date
	SelectedTime                   string  `json:"selected_time,omitempty"`                    // Temporary storage for selected time
	SelectedUnavailableStartTime   string  `json:"selected_unavailable_start_time,omitempty"`  // Temporary storage for selected unavailable start time
	SelectedUnavailableEndTime     string  `json:"selected_unavailable_end_time,omitempty"`    // Temporary storage for selected unavailable end time
	SelectedUnavailableDescription string  `json:"selected_unavailable_description,omitempty"` // Temporary storage for selected unavailable description
	SelectedAppointmentID          string  `json:"selected_appointment_id,omitempty"`          // Temporary storage for appointment cancellation
	CreatedAt                      string  `json:"created_at"`
	UpdatedAt                      string  `json:"updated_at"`
}

// User states for bot interaction
const (
	StateNone                = ""
	StateWaitingForRole      = "waiting_for_role"
	StateWaitingForFirstName = "waiting_for_first_name"
	StateWaitingForLastName  = "waiting_for_last_name"
	StateWaitingForPhone     = "waiting_for_phone"
	StateWaitingForUsername  = "waiting_for_username"
	StateWaitingForPassword  = "waiting_for_password"

	// Appointment booking states
	StateWaitingForProfessionalSelection = "waiting_for_professional_selection"
	StateWaitingForDateSelection         = "waiting_for_date_selection"
	StateWaitingForTimeSelection         = "waiting_for_time_selection"
	StateWaitingForCancellationReason    = "waiting_for_cancellation_reason"
	StateBookingAppointment              = "booking_appointment" // General state during booking process

	// Unavailable appointment states
	StateWaitingForUnavailableDateSelection = "waiting_for_unavailable_date_selection"
	StateWaitingForUnavailableStartTime     = "waiting_for_unavailable_start_time"
	StateWaitingForUnavailableEndTime       = "waiting_for_unavailable_end_time"
	StateWaitingForUnavailableDescription   = "waiting_for_unavailable_description"
)

// Appointment represents an appointment
type Appointment struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Status       string `json:"status"`
	Description  string `json:"description,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Professional *User  `json:"professional,omitempty"`
	Client       *User  `json:"client,omitempty"`
}

// CreateAppointmentResponse represents the response after creating an appointment
type CreateAppointmentResponse struct {
	Appointment  Appointment  `json:"appointment"`
	Client       Client       `json:"client"`
	Professional Professional `json:"professional"`
}

// Client represents a client in the response
type Client struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	ChatID      int64  `json:"chat_id,omitempty"`
}

// Professional represents a professional in the response
type Professional struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	ChatID      int64  `json:"chat_id,omitempty"`
}

// GetClientAppointmentsResponse represents the response for getting client appointments
type GetClientAppointmentsResponse struct {
	Appointments []ClientAppointment `json:"appointments"`
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

// CancelClientAppointmentResponse represents the response after cancelling an appointment by client
type CancelClientAppointmentResponse struct {
	Appointment  CancelledAppointment          `json:"appointment"`
	Client       ClientAppointmentClient       `json:"client"`
	Professional ClientAppointmentProfessional `json:"professional"`
}

// ClientAppointmentClient represents client details in appointment context
type ClientAppointmentClient struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	ChatID      *int64  `json:"chat_id,omitempty"`
}

// CancelledAppointment represents a cancelled appointment
type CancelledAppointment struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	StartTime          string `json:"start_time"`
	EndTime            string `json:"end_time"`
	Status             string `json:"status"`
	Description        string `json:"description,omitempty"`
	CancellationReason string `json:"cancellation_reason"`
	CancelledBy        string `json:"cancelled_by"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// GetProfessionalsResponse represents the response for getting all professionals
type GetProfessionalsResponse struct {
	Professionals []User `json:"professionals"`
}

// ProfessionalAvailabilityResponse represents the response for professional availability
type ProfessionalAvailabilityResponse struct {
	Date  string     `json:"date"`
	Slots []TimeSlot `json:"slots"`
}

// TimeSlot represents a one-hour time slot
type TimeSlot struct {
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Available   bool   `json:"available"`
	Type        string `json:"type,omitempty"`        // "appointment", "unavailable", or empty if available
	Description string `json:"description,omitempty"` // Description with client info if available
}

// GetProfessionalAppointmentsResponse represents the response for getting professional appointments
type GetProfessionalAppointmentsResponse struct {
	Appointments []ProfessionalAppointment `json:"appointments"`
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

// CreateUnavailableAppointmentResponse represents the response after creating an unavailable appointment
type CreateUnavailableAppointmentResponse struct {
	Appointment UnavailableAppointment `json:"appointment"`
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
