package models

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
