package handlers

// Error messages
const (
	ErrorMsgRegistrationFailed               = "❌ Registration failed: %v"
	ErrorMsgFailedToLoadProfessionals        = "❌ Failed to load professionals: %v"
	ErrorMsgNoProfessionals                  = "❌ No professionals available at the moment."
	ErrorMsgFailedToLoadAvailability         = "❌ Failed to load availability: %v"
	ErrorMsgInvalidTimeFormat                = "❌ Invalid time format"
	ErrorMsgInvalidDateFormat                = "❌ Invalid date format"
	ErrorMsgPastTimeNotAllowed               = "❌ Cannot book appointments in the past. Please select a future time."
	ErrorMsgFailedToCreateAppointment        = "❌ Failed to create appointment: %v"
	ErrorMsgFailedToLoadPendingAppointments  = "❌ Failed to load pending appointments: %v"
	ErrorMsgFailedToLoadUpcomingAppointments = "❌ Failed to load upcoming appointments: %v"
	ErrorMsgFailedToCancelAppointment        = "❌ Failed to cancel appointment: %v"
	ErrorMsgInvalidState                     = "❌ This action is not available in your current state. Please use /start to begin a new session."
	ErrorMsgBookingCancelled                 = "❌ Booking cancelled. Returning to dashboard."
)

// Success messages
const (
	SuccessMsgFirstNameSaved         = "✅ First name saved!\n\nPlease enter your last name:"
	SuccessMsgLastNameSaved          = "✅ Last name saved!\n\nPlease enter your phone number (optional, or type \"skip\" to skip):"
	SuccessMsgRegistrationSuccessful = "✅ Registration successful!\n\nWelcome, %s %s!\nRole: %s\nChat ID: %d"
	SuccessMsgAppointmentBooked      = "✅ Appointment booked successfully!\n\n📅 Date: %s\n🕐 Time: %s - %s\n👨‍💼 Professional: %s %s\n\nYour appointment is pending confirmation."
	SuccessMsgAppointmentCancelled   = "✅ Appointment cancelled successfully!\n\n📅 Date: %s\n🕐 Time: %s - %s\n👨‍💼 Professional: %s %s\n📝 Reason: %s"
)

// UI messages
const (
	UIMsgClientRegistration             = "👤 Client Registration\n\nPlease enter your first name:"
	UIMsgWelcomeBack                    = "👋 Welcome back, %s!\n\nYou are registered as a %s.\n\nWhat would you like to do?"
	UIMsgSelectProfessional             = "👨‍💼 Please select a professional:"
	UIMsgSelectDate                     = "📅 Select a date (%s %d):"
	UIMsgSelectTime                     = "🕐 Select a time slot for %s:"
	UIMsgNoPendingAppointments          = "📋 You have no pending appointments."
	UIMsgNoUpcomingAppointments         = "📋 You have no upcoming appointments."
	UIMsgNoUpcomingAppointmentsForMonth = "📋 You have no upcoming appointments for this month(%s)"
	UIMsgPendingAppointments            = "⏳ Your Pending Appointments:\n\n"
	UIMsgUpcomingAppointments           = "📋 Your Upcoming Appointments:\n\n"
	UIMsgCancellationReason             = "Please provide a reason for cancelling this appointment:"
	UIMsgNewAppointmentRequest          = "🔔 New Appointment Request!\n\n👤 Client: %s %s\n📅 Date: %s\n🕐 Time: %s - %s\n📝 Description: %s\n\nPlease confirm or cancel this appointment."
	UIMsgAppointmentCancelled           = "🔔 Appointment Cancelled\n\n👤 Client: %s %s\n📅 Date: %s\n🕐 Time: %s - %s\n📝 Reason: %s"
)

// Button texts
const (
	BtnBookAppointment          = "📅 Book Appointment"
	BtnMyPendingAppointments    = "⏳ My Pending Appointments"
	BtnMyUpcomingAppointments   = "✅ My Upcoming Appointments"
	BtnMyTimetable              = "📅 My Timetable"
	BtnCancelBooking            = "❌ Cancel Booking"
	BtnCancelAppointment        = "❌ Cancel Appointment #%d"
	BtnBackToDashboard          = "🏠 Back to Dashboard"
	BtnGoToDashboard            = "🏠 Go to Dashboard"
	BtnPreviousMonth            = "⬅️ Previous"
	BtnNextMonth                = "Next ➡️"
	BtnConfirmAppointment       = "✅ Confirm"
	BtnCancelAppointmentConfirm = "❌ Cancel"
)

// Professional-specific error messages
const (
	ErrorMsgSignInFailed                         = "❌ Sign in failed: %v"
	ErrorMsgFailedToConfirmAppointment           = "❌ Failed to confirm appointment: %v"
	ErrorMsgFailedToLoadAppointments             = "❌ Failed to load appointments: %v"
	ErrorMsgFailedToCreateUnavailableAppointment = "❌ Failed to create unavailable appointment: %v"
	ErrorMsgUnavailableCancelled                 = "❌ Unavailable appointment setting cancelled. Returning to dashboard."
)

// Professional-specific success messages
const (
	SuccessMsgUsernameSaved        = "✅ Username saved!\n\nPlease enter your password:"
	SuccessMsgSignInSuccessful     = "✅ Sign in successful!\n\nWelcome back, %s %s!\nRole: %s\nUsername: %s\nChat ID: %d"
	SuccessMsgAppointmentConfirmed = "✅ Appointment confirmed successfully!\n\n📅 Date: %s\n🕐 Time: %s - %s\n👤 Client: %s %s"
	SuccessMsgUnavailablePeriodSet = "✅ Unavailable period set successfully!\n\n📅 Date: %s\n🕐 Time: %s - %s\n📝 Description: %s"
)

// Professional-specific UI messages
const (
	UIMsgProfessionalSignIn                 = "👨‍💼 Professional Sign In\n\nPlease enter your username:"
	UIMsgWelcomeBackProfessional            = "👋 Welcome back, %s!\n\nYou are registered as a %s.\n\nWhat would you like to do?"
	UIMsgSelectUnavailableDate              = "📅 Select a date for unavailable time (%s %d):"
	UIMsgSelectUnavailableStartTime         = "🕐 Select start time for unavailable period on %s:"
	UIMsgSelectUnavailableEndTime           = "🕐 Select end time for unavailable period (starting at %s):"
	UIMsgUnavailableDescription             = "📝 Please provide a description for your unavailable period:\n\n📅 Date: %s\n🕐 Time: %s - %s\n\nExample: \"Personal break\", \"Lunch time\", \"Out of office\", etc."
	UIMsgUnavailableSlotWarning             = "⚠️ You can only select times before %s (%s)"
	UIMsgNoAvailableTimeSlots               = "❌ No available time slots before your next unavailable period."
	UIMsgSelectUpcomingAppointmentsDate     = "📅 Here are the dates with upcoming appointments. Select a date to view upcoming appointments:"
	UIMsgTimetableEmpty                     = "📋 No activities scheduled for this day(%s)."
	UIMsgTimetableHeader                    = "📋 Your Timetable for %s:\n\n"
	UIMsgTimetableSlot                      = "📅 Slot #%d:\n🕐 %s - %s\n📝 %s\n\n"
	UIMsgAppointmentConfirmed               = "✅ Appointment Confirmed!\n\n📅 Date: %s\n🕐 Time: %s - %s\n👨‍💼 Professional: %s %s\n\nYour appointment has been confirmed."
	UIMsgAppointmentCancelledByProfessional = "🔔 Appointment Cancelled by Professional\n\n📅 Date: %s\n🕐 Time: %s - %s\n👨‍💼 Professional: %s %s\n📝 Reason: %s"
)

// Professional-specific button texts
const (
	BtnPendingAppointments      = "⏳ Pending Appointments"
	BtnUpcomingAppointments     = "📋 Upcoming Appointments"
	BtnSetUnavailable           = "🚫 Set Unavailable"
	BtnConfirmAppointmentProf   = "✅ Confirm Appointment #%d"
	BtnCancelAppointmentProf    = "❌ Cancel Appointment #%d"
	BtnCancelAppointmentProfAlt = "❌ Cancel Appintment %d"
	BtnPreviousUnavailableMonth = "⬅️ Previous"
	BtnNextUnavailableMonth     = "Next ➡️"
	BtnCancelUnavailable        = "❌ Cancel"
	BtnPreviousTimetableDay     = "⬅️ Previous Day"
	BtnNextTimetableDay         = "Next Day ➡️"
	BtnCancelTimetableSlot      = "❌ Cancel Slot #%d"
)

// Navigation directions
const (
	DirectionPrev = "prev"
	DirectionNext = "next"
)

// Keyboard layouts
const (
	DaysPerRow      = 7
	TimeSlotsPerRow = 3
)
