package handlers

import (
	"fmt"
	"time"

	"booking_client/internal/models"
	"booking_client/internal/schema"
	"booking_client/internal/services"
	"booking_client/internal/util"
	"booking_client/pkg/telegram"

	"github.com/rs/zerolog"
)

// ProfessionalHandler handles all professional-related operations
type ProfessionalHandler struct {
	bot                 *telegram.Bot
	logger              *zerolog.Logger
	apiService          *services.APIService
	notificationService *NotificationService
}

// NewProfessionalHandler creates a new professional handler
func NewProfessionalHandler(bot *telegram.Bot, logger *zerolog.Logger, apiService *services.APIService) *ProfessionalHandler {
	return &ProfessionalHandler{
		bot:                 bot,
		logger:              logger,
		apiService:          apiService,
		notificationService: NewNotificationService(bot, logger),
	}
}

// StartSignIn starts the professional sign-in process
func (h *ProfessionalHandler) StartSignIn(chatID int64) {
	// Create a temporary user with state
	tempUser := &models.User{
		ChatID: &chatID,
		Role:   "professional",
		State:  models.StateWaitingForUsername,
	}

	// Store in memory for state tracking
	h.apiService.GetUserRepository().SetUser(chatID, tempUser)

	h.sendMessage(chatID, UIMsgProfessionalSignIn)
}

// ShowDashboard shows the professional dashboard with appointment options
func (h *ProfessionalHandler) ShowDashboard(chatID int64, user *models.User) {
	currentUser, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	currentUser.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, currentUser)

	text := fmt.Sprintf(UIMsgWelcomeBackProfessional, currentUser.LastName, currentUser.Role)
	keyboard := h.createProfessionalDashboardKeyboard()

	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUsernameInput handles username input for professional sign-in
func (h *ProfessionalHandler) HandleUsernameInput(chatID int64, username string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.Username = username
	user.State = models.StateWaitingForPassword
	h.apiService.GetUserRepository().SetUser(chatID, user)

	h.sendMessage(chatID, SuccessMsgUsernameSaved)
}

// HandlePasswordInput handles password input for professional sign-in
func (h *ProfessionalHandler) HandlePasswordInput(chatID int64, password string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Sign in the professional
	req := &schema.ProfessionalSignInRequest{
		Username: user.Username,
		Password: password,
		ChatID:   chatID,
	}

	signedInUser, err := h.apiService.SignInProfessional(req)
	if err != nil {
		h.sendError(chatID, ErrorMsgSignInFailed, err)
		return
	}

	// Clear state
	user.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, signedInUser)

	text := fmt.Sprintf(SuccessMsgSignInSuccessful, signedInUser.FirstName, signedInUser.LastName, signedInUser.Role, signedInUser.Username, chatID)
	h.sendMessage(chatID, text)
	h.ShowDashboard(chatID, signedInUser)
}

// HandlePendingAppointments shows pending appointments for professionals
func (h *ProfessionalHandler) HandlePendingAppointments(chatID int64) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetProfessionalAppointments(user.ID, "pending")
	if err != nil {
		h.sendError(chatID, ErrorMsgFailedToLoadPendingAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		h.sendMessage(chatID, UIMsgNoPendingAppointments)
		h.ShowDashboard(chatID, user)
		return
	}

	text := UIMsgPendingAppointments
	for index, apt := range appointments.Appointments {
		text += formatProfessionalAppointmentDetails(&apt, index)
	}

	keyboard := h.createProfessionalAppointmentsKeyboard(appointments.Appointments, true)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUpcomingAppointments shows upcoming appointments for professionals
func (h *ProfessionalHandler) HandleUpcomingAppointments(chatID int64) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Show date picker for upcoming appointments
	h.showUpcomingAppointmentsDatePicker(chatID, user)
}

// showUpcomingAppointmentsDatePicker shows a date picker for upcoming appointments
func (h *ProfessionalHandler) showUpcomingAppointmentsDatePicker(chatID int64, user *models.User, month ...string) {
	// Get current month or use provided month
	var currentMonth string
	if len(month) > 0 {
		currentMonth = month[0]
	} else {
		currentMonth = time.Now().Format("2006-01")
	}

	// Get dates with appointments for current month
	datesResponse, err := h.apiService.GetProfessionalAppointmentDates(user.ID, currentMonth)
	if err != nil {
		h.sendError(chatID, ErrorMsgFailedToLoadUpcomingAppointments, err)
		return
	}

	if len(datesResponse.Dates) == 0 {
		h.sendMessage(chatID, UIMsgNoUpcomingAppointments)
		h.ShowDashboard(chatID, user)
		return
	}

	text := UIMsgSelectUpcomingAppointmentsDate
	keyboard := h.createUpcomingAppointmentsDateKeyboard(datesResponse.Dates, currentMonth)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUpcomingAppointmentsMonthNavigation handles month navigation for upcoming appointments
func (h *ProfessionalHandler) HandleUpcomingAppointmentsMonthNavigation(chatID int64, monthStr string, direction string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Parse current month
	currentMonth, err := time.Parse("2006-01", monthStr)
	if err != nil {
		h.sendError(chatID, ErrorMsgInvalidDateFormat, err)
		return
	}

	// Calculate new month based on direction
	var newMonth time.Time
	if direction == "prev" {
		newMonth = currentMonth.AddDate(0, -1, 0)
	} else {
		newMonth = currentMonth.AddDate(0, 1, 0)
	}

	newMonthStr := newMonth.Format("2006-01")

	// Get dates with appointments for the new month
	datesResponse, err := h.apiService.GetProfessionalAppointmentDates(user.ID, newMonthStr)
	if err != nil {
		h.sendError(chatID, ErrorMsgFailedToLoadUpcomingAppointments, err)
		return
	}

	if len(datesResponse.Dates) == 0 {
		// If no appointments in this month, go back to previous month
		var previousMonth time.Time
		if direction == "prev" {
			// If we went to previous month and it's empty, go back to current month
			previousMonth = time.Now()
		} else {
			// If we went to next month and it's empty, go back to previous month
			previousMonth = newMonth.AddDate(0, -1, 0)
		}

		previousMonthStr := previousMonth.Format("2006-01")
		h.sendMessage(chatID, UIMsgNoUpcomingAppointments)

		// Show the previous month
		h.showUpcomingAppointmentsDatePicker(chatID, user, previousMonthStr)
		return
	}

	text := UIMsgSelectUpcomingAppointmentsDate
	keyboard := h.createUpcomingAppointmentsDateKeyboard(datesResponse.Dates, newMonthStr)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUpcomingAppointmentsDateSelection handles selection of a date for upcoming appointments
func (h *ProfessionalHandler) HandleUpcomingAppointmentsDateSelection(chatID int64, dateStr string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Get appointments for the selected date
	appointments, err := h.apiService.GetProfessionalAppointmentsByDate(user.ID, "confirmed", dateStr)
	if err != nil {
		h.sendError(chatID, ErrorMsgFailedToLoadUpcomingAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		h.sendMessage(chatID, UIMsgNoUpcomingAppointments)
		h.ShowDashboard(chatID, user)
		return
	}

	// Format the date for display
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.sendError(chatID, "Invalid date format", err)
		return
	}

	text := fmt.Sprintf("📅 Upcoming Appointments for %s:\n\n", date.Format("Monday, January 2, 2006"))
	for index, apt := range appointments.Appointments {
		text += formatProfessionalAppointmentDetails(&apt, index)
	}

	keyboard := h.createProfessionalAppointmentsKeyboard(appointments.Appointments, false)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleConfirmAppointment confirms an appointment by professional
func (h *ProfessionalHandler) HandleConfirmAppointment(chatID int64, appointmentID string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Confirm the appointment
	req := &schema.ConfirmAppointmentRequest{}

	response, err := h.apiService.ConfirmProfessionalAppointment(user.ID, appointmentID, req)
	if err != nil {
		h.sendError(chatID, ErrorMsgFailedToConfirmAppointment, err)
		return
	}

	date, startTime, endTime := formatAppointmentTime(response.Appointment.StartTime, response.Appointment.EndTime)
	text := fmt.Sprintf(SuccessMsgAppointmentConfirmed,
		date, startTime, endTime,
		response.Client.FirstName, response.Client.LastName)

	h.sendMessage(chatID, text)

	// Notify client about confirmation
	h.notificationService.NotifyClientAppointmentConfirmation(response)
	h.ShowDashboard(chatID, user)
}

// HandleCancelAppointment starts the professional appointment cancellation process
func (h *ProfessionalHandler) HandleCancelAppointment(chatID int64, appointmentID string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Store appointment ID and ask for cancellation reason
	user.State = models.StateWaitingForCancellationReason
	user.SelectedAppointmentID = appointmentID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	h.sendMessage(chatID, UIMsgCancellationReason)
}

// HandleCancellationReason handles the professional cancellation reason input
func (h *ProfessionalHandler) HandleCancellationReason(chatID int64, reason string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	appointmentID := user.SelectedAppointmentID

	// Cancel the appointment
	req := &schema.CancelAppointmentRequest{
		CancellationReason: reason,
	}

	response, err := h.apiService.CancelProfessionalAppointment(user.ID, appointmentID, req)
	if err != nil {
		h.sendError(chatID, ErrorMsgFailedToCancelAppointment, err)
		return
	}

	// Clear state
	user.State = models.StateNone
	user.SelectedAppointmentID = ""
	h.apiService.GetUserRepository().SetUser(chatID, user)

	date, startTime, endTime := formatAppointmentTime(response.Appointment.StartTime, response.Appointment.EndTime)
	text := fmt.Sprintf(SuccessMsgAppointmentCancelled,
		date, startTime, endTime,
		response.Client.FirstName, response.Client.LastName,
		response.Appointment.CancellationReason)

	h.sendMessage(chatID, text)

	// Notify client about cancellation
	h.notificationService.NotifyClientProfessionalCancellation(response)
	h.ShowDashboard(chatID, user)
}

// HandleSetUnavailable starts the unavailable appointment setting process
func (h *ProfessionalHandler) HandleSetUnavailable(chatID int64) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Set state for unavailable appointment
	user.State = models.StateWaitingForUnavailableDateSelection
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Show current month dates
	h.showUnavailableDateSelection(chatID, time.Now())
}

// showUnavailableDateSelection shows available dates for the current month
func (h *ProfessionalHandler) showUnavailableDateSelection(chatID int64, currentDate time.Time) {
	text := fmt.Sprintf(UIMsgSelectUnavailableDate, currentDate.Month(), currentDate.Year())
	keyboard := h.createUnavailableDateKeyboard(currentDate)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUnavailableDateSelection handles when user selects a date for unavailable time
func (h *ProfessionalHandler) HandleUnavailableDateSelection(chatID int64, date string) {
	user, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForUnavailableDateSelection,
	})
	if !valid {
		return
	}

	user.State = models.StateWaitingForUnavailableStartTime
	user.SelectedDate = date
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Get availability for selected date to show time slots
	availability, err := h.apiService.GetProfessionalAvailability(user.ID, date)
	if err != nil {
		h.sendError(chatID, ErrorMsgFailedToLoadAvailability, err)
		return
	}

	h.showUnavailableStartTimeSelection(chatID, availability)
}

// showUnavailableStartTimeSelection shows available time slots for start time
func (h *ProfessionalHandler) showUnavailableStartTimeSelection(chatID int64, availability *models.ProfessionalAvailabilityResponse) {
	text := fmt.Sprintf(UIMsgSelectUnavailableStartTime, availability.Date)
	keyboard := h.createUnavailableStartTimeKeyboard(availability)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUnavailableStartTimeSelection handles when user selects start time for unavailable period
func (h *ProfessionalHandler) HandleUnavailableStartTimeSelection(chatID int64, startTime string) {
	user, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForUnavailableStartTime,
	})
	if !valid {
		return
	}

	user.State = models.StateWaitingForUnavailableEndTime
	user.SelectedUnavailableStartTime = startTime // Store start time temporarily
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Get availability for the selected date to determine available end times
	availability, err := h.apiService.GetProfessionalAvailability(user.ID, user.SelectedDate)
	if err != nil {
		h.sendError(chatID, ErrorMsgFailedToLoadAvailability, err)
		return
	}

	h.showUnavailableEndTimeSelection(chatID, startTime, availability)
}

// showUnavailableEndTimeSelection shows available time slots for end time
func (h *ProfessionalHandler) showUnavailableEndTimeSelection(chatID int64, startTime string, availability *models.ProfessionalAvailabilityResponse) {
	text := fmt.Sprintf(UIMsgSelectUnavailableEndTime, startTime)

	// Find the first unavailable slot after the selected start time to show warning
	var firstUnavailableSlot *models.TimeSlot

	for _, slot := range availability.Slots {
		slotStart, err := time.Parse(time.RFC3339, slot.StartTime)
		if err != nil {
			continue
		}
		slotStartLocal := slotStart.In(util.GetAppTimezone())
		slotTimeStr := slotStartLocal.Format("15:04")

		// Only consider slots that are after the selected start time
		if !slot.Available && slotTimeStr > startTime {
			firstUnavailableSlot = &slot
			break
		}
	}

	if firstUnavailableSlot != nil {
		unavailableStart, _ := time.Parse(time.RFC3339, firstUnavailableSlot.StartTime)
		unavailableStartLocal := unavailableStart.In(util.GetAppTimezone())

		// Build slot details with enhanced information
		slotDetails := fmt.Sprintf("Unavailable slot at %s", unavailableStartLocal.Format("15:04"))
		if firstUnavailableSlot.Type != "" {
			slotDetails += fmt.Sprintf(" (%s)", firstUnavailableSlot.Type)
		}
		if firstUnavailableSlot.Description != "" {
			slotDetails += fmt.Sprintf(" - %s", firstUnavailableSlot.Description)
		}

		text += fmt.Sprintf("\n\n"+UIMsgUnavailableSlotWarning, unavailableStartLocal.Format("15:04"), slotDetails)
	}

	keyboard := h.createUnavailableEndTimeKeyboard(startTime, availability)

	// If no slots available, show a message
	if len(keyboard.InlineKeyboard) == 1 && len(keyboard.InlineKeyboard[0]) == 1 && keyboard.InlineKeyboard[0][0].Text == "❌ Cancel" {
		text += "\n\n" + UIMsgNoAvailableTimeSlots
	}

	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUnavailableEndTimeSelection handles when user selects end time for unavailable period
func (h *ProfessionalHandler) HandleUnavailableEndTimeSelection(chatID int64, endTime string) {
	user, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForUnavailableEndTime,
	})
	if !valid {
		return
	}

	// Store end time and ask for description
	user.State = models.StateWaitingForUnavailableDescription
	user.SelectedUnavailableEndTime = endTime
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(UIMsgUnavailableDescription, user.SelectedDate, user.SelectedUnavailableStartTime, endTime)
	h.sendMessage(chatID, text)
}

// HandleUnavailableDescription handles when user provides description for unavailable period
func (h *ProfessionalHandler) HandleUnavailableDescription(chatID int64, description string) {
	user, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForUnavailableDescription,
	})
	if !valid {
		return
	}

	// Create unavailable appointment
	start, _ := time.Parse("15:04", user.SelectedUnavailableStartTime)
	end, _ := time.Parse("15:04", user.SelectedUnavailableEndTime)
	date := user.SelectedDate

	// Parse the date and combine with times
	selectedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		h.sendMessage(chatID, ErrorMsgInvalidDateFormat)
		return
	}

	// Combine date with times in application timezone
	startDateTime := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(),
		start.Hour(), start.Minute(), 0, 0, util.GetAppTimezone())
	endDateTime := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(),
		end.Hour(), end.Minute(), 0, 0, util.GetAppTimezone())

	// Create unavailable appointment request
	req := &schema.CreateUnavailableAppointmentRequest{
		ProfessionalID: user.ID,
		StartAt:        startDateTime.Format(time.RFC3339),
		EndAt:          endDateTime.Format(time.RFC3339),
		Description:    description,
	}

	appointment, err := h.apiService.CreateUnavailableAppointment(req)
	if err != nil {
		h.sendError(chatID, ErrorMsgFailedToCreateUnavailableAppointment, err)
		return
	}
	// Clear state
	h.clearUnavailableState(user)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(SuccessMsgUnavailablePeriodSet,
		date, start.Format("15:04"), end.Format("15:04"), appointment.Appointment.Description)

	h.sendMessage(chatID, text)
	h.ShowDashboard(chatID, user)
}

// HandleCancelUnavailable cancels the unavailable appointment setting process
func (h *ProfessionalHandler) HandleCancelUnavailable(chatID int64) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Clear all unavailable-related state
	h.clearUnavailableState(user)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	h.sendMessage(chatID, ErrorMsgUnavailableCancelled)
	h.ShowDashboard(chatID, user)
}

// HandlePrevUnavailableMonth handles previous month navigation for unavailable appointments
func (h *ProfessionalHandler) HandlePrevUnavailableMonth(chatID int64) {
	// Validate user state - only allow if waiting for unavailable date selection
	_, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForUnavailableDateSelection,
	})
	if !valid {
		return
	}

	// For simplicity, we'll just show current month again
	// In a real implementation, you'd store the current month in user state
	h.showUnavailableDateSelection(chatID, time.Now().AddDate(0, -1, 0))
}

// HandleNextUnavailableMonth handles next month navigation for unavailable appointments
func (h *ProfessionalHandler) HandleNextUnavailableMonth(chatID int64) {
	// Validate user state - only allow if waiting for unavailable date selection
	_, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForUnavailableDateSelection,
	})
	if !valid {
		return
	}

	// For simplicity, we'll just show current month again
	// In a real implementation, you'd store the current month in user state
	h.showUnavailableDateSelection(chatID, time.Now().AddDate(0, 1, 0))
}
