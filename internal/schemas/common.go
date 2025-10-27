package schemas

// Appointment represents an appointment
type Appointment struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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

// TimeSlot represents a one-hour time slot
type TimeSlot struct {
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Available   bool   `json:"available"`
	Type        string `json:"type,omitempty"`        // "appointment", "unavailable", or empty if available
	Description string `json:"description,omitempty"` // Description with client info if available
}
