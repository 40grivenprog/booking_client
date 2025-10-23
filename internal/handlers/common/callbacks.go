package common

// Callback constants for inline keyboard buttons
// These are used to construct callback_data strings that identify which action to perform

const (
	// ========================================
	// EXACT MATCH CALLBACKS (no parameters)
	// ========================================

	// Initial selection
	CallbackClient       = "client"
	CallbackProfessional = "professional"

	// Client callbacks
	CallbackBookAppointment      = "book_appointment"
	CallbackPendingAppointments  = "pending_appointments"
	CallbackUpcomingAppointments = "upcoming_appointments"
	CallbackCancelBooking        = "cancel_booking"

	// Professional callbacks
	CallbackProfessionalPendingAppointments  = "professional_pending_appointments"
	CallbackProfessionalUpcomingAppointments = "professional_upcoming_appointments"
	CallbackProfessionalTimetable            = "professional_timetable"
	CallbackSetUnavailable                   = "set_unavailable"

	// Unavailable navigation
	CallbackPrefixPrevUnavailableMonth = "prev_unavailable_month_"
	CallbackPrefixNextUnavailableMonth = "next_unavailable_month_"
	CallbackCancelUnavailable          = "cancel_unavailable"

	// Common
	CallbackBackToDashboard = "back_to_dashboard"

	// ========================================
	// PREFIX CALLBACKS (with parameters)
	// ========================================

	// Professional timetable navigation
	CallbackPrefixPrevTimetableDay = "prev_timetable_day_"
	CallbackPrefixNextTimetableDay = "next_timetable_day_"

	// Professional upcoming appointments navigation
	CallbackPrefixPrevUpcomingMonth  = "prev_upcoming_month_"
	CallbackPrefixNextUpcomingMonth  = "next_upcoming_month_"
	CallbackPrefixSelectUpcomingDate = "select_upcoming_date_"

	// Client booking flow - month navigation
	CallbackPrefixPrevMonth = "prev_month_"
	CallbackPrefixNextMonth = "next_month_"

	// Selection callbacks
	CallbackPrefixSelectProfessional = "select_professional_"
	CallbackPrefixSelectDate         = "select_date_"
	CallbackPrefixSelectTime         = "select_time_"

	// Appointment actions
	CallbackPrefixCancelAppointment  = "cancel_appointment_"
	CallbackPrefixConfirmAppointment = "confirm_appointment_"
	CallbackPrefixCancelProfAppt     = "cancel_prof_appt_"

	// Unavailable flow
	CallbackPrefixSelectUnavailableDate  = "select_unavailable_date_"
	CallbackPrefixSelectUnavailableStart = "select_unavailable_start_"
	CallbackPrefixSelectUnavailableEnd   = "select_unavailable_end_"
)

// BuildCallback constructs a callback string from a prefix and parameter
func BuildCallback(prefix, param string) string {
	return prefix + param
}
