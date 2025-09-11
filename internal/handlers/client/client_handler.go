package client

import (
	"fmt"
	"time"

	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	"booking_client/internal/schema"
	"booking_client/internal/services"
	"booking_client/internal/util"
	"booking_client/pkg/telegram"

	"github.com/rs/zerolog"
)

// ClientHandler handles all client-related operations
type ClientHandler struct {
	bot                 *telegram.Bot
	logger              *zerolog.Logger
	apiService          *services.APIService
	notificationService *common.NotificationService
}

// NewClientHandler creates a new client handler
func NewClientHandler(bot *telegram.Bot, logger *zerolog.Logger, apiService *services.APIService) *ClientHandler {
	return &ClientHandler{
		bot:                 bot,
		logger:              logger,
		apiService:          apiService,
		notificationService: common.NewNotificationService(bot, logger, apiService),
	}
}

// StartRegistration starts the client registration process
func (h *ClientHandler) StartRegistration(chatID int64) {
	// Create a temporary user with state
	tempUser := &models.User{
		ChatID: &chatID,
		Role:   "client",
		State:  models.StateWaitingForFirstName,
	}

	// Store in memory for state tracking
	h.apiService.GetUserRepository().SetUser(chatID, tempUser)

	h.sendMessage(chatID, common.UIMsgClientRegistration)
}

// ShowDashboard shows the client dashboard with appointment options
func (h *ClientHandler) ShowDashboard(chatID int64) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(common.UIMsgWelcomeBack, user.FirstName, user.Role)
	keyboard := h.createDashboardKeyboard()

	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleFirstNameInput handles first name input for client registration
func (h *ClientHandler) HandleFirstNameInput(chatID int64, firstName string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.FirstName = firstName
	user.State = models.StateWaitingForLastName
	h.apiService.GetUserRepository().SetUser(chatID, user)

	h.sendMessage(chatID, common.SuccessMsgFirstNameSaved)
}

// HandleLastNameInput handles last name input for client registration
func (h *ClientHandler) HandleLastNameInput(chatID int64, lastName string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.LastName = lastName
	user.State = models.StateWaitingForPhone
	h.apiService.GetUserRepository().SetUser(chatID, user)

	h.sendMessage(chatID, common.SuccessMsgLastNameSaved)
}

// HandlePhoneInput handles phone number input for client registration
func (h *ClientHandler) HandlePhoneInput(chatID int64, phone string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	var phoneNumber *string
	if phone != "skip" && phone != "" {
		phoneNumber = &phone
	}

	// Register the client
	req := &schema.RegisterRequest{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		ChatID:      chatID,
		PhoneNumber: phoneNumber,
		Role:        "client",
	}

	registeredUser, err := h.apiService.RegisterClient(req)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgRegistrationFailed, err)
		return
	}

	// Clear state
	user.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, registeredUser)

	text := fmt.Sprintf(common.SuccessMsgRegistrationSuccessful, registeredUser.FirstName, registeredUser.LastName, registeredUser.Role, chatID)
	keyboard := h.createRegistrationSuccessKeyboard()

	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleBookAppointment starts the appointment booking process
func (h *ClientHandler) HandleBookAppointment(chatID int64) {
	// Validate user state - only allow if not in any specific state
	user, valid := h.validateUserState(chatID, []string{
		models.StateNone,
	})
	if !valid {
		return
	}

	// Set booking state
	user.State = models.StateWaitingForProfessionalSelection
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Get professionals
	professionals, err := h.apiService.GetProfessionals()
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToLoadProfessionals, err)
		return
	}

	if len(professionals.Professionals) == 0 {
		h.sendMessage(chatID, common.ErrorMsgNoProfessionals)
		h.ShowDashboard(chatID)
		return
	}

	keyboard := h.createProfessionalsKeyboard(professionals.Professionals)
	id, err := h.sendMessageWithKeyboardAndID(chatID, common.UIMsgSelectProfessional, keyboard)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	h.apiService.GetUserRepository().SetUser(chatID, user)
}

// HandleProfessionalSelection handles when user selects a professional
func (h *ClientHandler) HandleProfessionalSelection(chatID int64, professionalID string) {
	// Validate user state - only allow if not in any specific state or in booking state
	user, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForProfessionalSelection,
	})
	if !valid {
		return
	}

	user.State = models.StateWaitingForDateSelection
	user.SelectedProfessionalID = professionalID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Show current month dates
	h.showDateSelection(user, time.Now())

	h.bot.DeleteMessage(chatID, *user.LastMessageID)
}

// showDateSelection shows available dates for the current month
func (h *ClientHandler) showDateSelection(user *models.User, currentDate time.Time) {
	text := fmt.Sprintf(common.UIMsgSelectDate, currentDate.Month(), currentDate.Year())
	keyboard := h.createDateKeyboard(currentDate)
	id, err := h.bot.SendMessageWithKeyboardAndID(*user.ChatID, text, keyboard)
	if err != nil {
		h.sendError(*user.ChatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	h.apiService.GetUserRepository().SetUser(*user.ChatID, user)
}

// HandleDateSelection handles when user selects a date
func (h *ClientHandler) HandleDateSelection(chatID int64, date string) {
	// Validate user state - only allow if waiting for date selection
	user, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForDateSelection,
	})
	if !valid {
		return
	}

	user.State = models.StateWaitingForTimeSelection
	user.SelectedDate = date
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Get availability for selected date
	professionalID := user.SelectedProfessionalID
	availability, err := h.apiService.GetProfessionalAvailability(professionalID, date)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToLoadAvailability, err)
		return
	}

	h.showTimeSelection(chatID, availability)
	h.bot.DeleteMessage(chatID, *user.LastMessageID)
}

// showTimeSelection shows available time slots
func (h *ClientHandler) showTimeSelection(chatID int64, availability *models.ProfessionalAvailabilityResponse) {
	text := fmt.Sprintf(common.UIMsgSelectTime, availability.Date)
	keyboard := h.createTimeKeyboard(availability)
	id, err := h.sendMessageWithKeyboardAndID(chatID, text, keyboard)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}

	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.State = models.StateWaitingForTimeSelection
	user.LastMessageID = &id
	h.apiService.GetUserRepository().SetUser(chatID, user)
}

// HandleTimeSelection handles when user selects a time slot
func (h *ClientHandler) HandleTimeSelection(chatID int64, startTime string) {
	// Validate user state - only allow if waiting for time selection
	user, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForTimeSelection,
	})
	if !valid {
		return
	}

	// Parse start time and calculate end time (1 hour later)
	h.logger.Debug().Str("startTime", startTime).Msg("Parsing start time")
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		h.logger.Error().Err(err).Str("startTime", startTime).Msg("Failed to parse start time")
		h.sendMessage(chatID, common.ErrorMsgInvalidTimeFormat)
		return
	}

	end := start.Add(time.Hour)
	date := user.SelectedDate

	// Create proper RFC3339 format datetime strings
	// Parse the date and combine with time
	selectedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		h.sendMessage(chatID, common.ErrorMsgInvalidDateFormat)
		return
	}

	// Combine date with time in application timezone for storage
	startDateTime := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(),
		start.Hour(), start.Minute(), 0, 0, util.GetAppTimezone())
	endDateTime := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(),
		end.Hour(), end.Minute(), 0, 0, util.GetAppTimezone())

	// Validate that start_time is in the future
	if startDateTime.Before(util.NowInAppTimezone()) {
		h.sendMessage(chatID, common.ErrorMsgPastTimeNotAllowed)
		return
	}

	// Create appointment with RFC3339 format
	req := &schema.CreateAppointmentRequest{
		ClientID:       user.ID,
		ProfessionalID: user.SelectedProfessionalID,
		StartTime:      startDateTime.Format(time.RFC3339),
		EndTime:        endDateTime.Format(time.RFC3339),
	}

	appointment, err := h.apiService.CreateAppointment(req)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToCreateAppointment, err)
		return
	}

	// Clear state and show success
	h.clearBookingState(user)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(common.SuccessMsgAppointmentBooked,
		date, startTime, end.Format("15:04"),
		appointment.Professional.FirstName, appointment.Professional.LastName)

	h.sendMessage(chatID, text)

	// Send notification to professional
	h.notificationService.NotifyProfessionalNewAppointment(appointment)
	h.ShowDashboard(chatID)
	h.bot.DeleteMessage(chatID, *user.LastMessageID)
}

// HandlePrevMonth handles previous month navigation
func (h *ClientHandler) HandlePrevMonth(chatID int64) {
	// Validate user state - only allow if waiting for date selection
	_, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForDateSelection,
	})
	if !valid {
		return
	}

	// For simplicity, we'll just show current month again
	// In a real implementation, you'd store the current month in user state
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	h.showDateSelection(user, time.Now().AddDate(0, -1, 0))
}

// HandleNextMonth handles next month navigation
func (h *ClientHandler) HandleNextMonth(chatID int64) {
	// Validate user state - only allow if waiting for date selection
	_, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForDateSelection,
	})
	if !valid {
		return
	}

	// For simplicity, we'll just show current month again
	// In a real implementation, you'd store the current month in user state
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	h.showDateSelection(user, time.Now().AddDate(0, 1, 0))
}

// HandleCancelBooking cancels the current booking process and returns to dashboard
func (h *ClientHandler) HandleCancelBooking(chatID int64) {
	// Validate user state - only allow if in booking process
	user, valid := h.validateUserState(chatID, []string{
		models.StateWaitingForProfessionalSelection,
		models.StateWaitingForDateSelection,
		models.StateWaitingForTimeSelection,
		models.StateBookingAppointment,
	})
	if !valid {
		return
	}

	// Clear all booking-related state
	h.clearBookingState(user)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	h.sendMessage(chatID, common.ErrorMsgBookingCancelled)
	h.ShowDashboard(chatID)
}

// HandlePendingAppointments shows pending appointments
func (h *ClientHandler) HandlePendingAppointments(chatID int64) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetClientAppointments(user.ID, "pending")
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToLoadPendingAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		h.sendMessage(chatID, common.UIMsgNoPendingAppointments)
		h.ShowDashboard(chatID)
		return
	}

	text := common.UIMsgPendingAppointments
	for index, apt := range appointments.Appointments {
		text += common.FormatAppointmentDetails(&apt, index)
	}

	keyboard := h.createAppointmentsKeyboard(appointments.Appointments, common.BtnCancelAppointment)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUpcomingAppointments shows upcoming appointments
func (h *ClientHandler) HandleUpcomingAppointments(chatID int64) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetClientAppointments(user.ID, "confirmed")
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToLoadUpcomingAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		h.sendMessage(chatID, common.UIMsgNoUpcomingAppointments)
		h.ShowDashboard(chatID)
		return
	}

	text := common.UIMsgUpcomingAppointments
	for index, apt := range appointments.Appointments {
		text += common.FormatAppointmentDetails(&apt, index)
	}

	keyboard := h.createAppointmentsKeyboard(appointments.Appointments, common.BtnCancelAppointment)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleCancelAppointment starts the appointment cancellation process
func (h *ClientHandler) HandleCancelAppointment(chatID int64, appointmentID string) {
	// Store appointment ID and ask for cancellation reason
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.State = models.StateWaitingForCancellationReason
	user.SelectedAppointmentID = appointmentID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	h.sendMessage(chatID, common.UIMsgCancellationReason)
}

// HandleCancellationReason handles the cancellation reason input
func (h *ClientHandler) HandleCancellationReason(chatID int64, reason string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	appointmentID := user.SelectedAppointmentID

	// Cancel the appointment
	req := &schema.CancelAppointmentRequest{
		CancellationReason: reason,
	}

	response, err := h.apiService.CancelClientAppointment(user.ID, appointmentID, req)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToCancelAppointment, err)
		return
	}

	// Clear state
	user.State = models.StateNone
	user.SelectedAppointmentID = ""
	h.apiService.GetUserRepository().SetUser(chatID, user)

	date, startTime, endTime := common.FormatAppointmentTime(response.Appointment.StartTime, response.Appointment.EndTime)
	text := fmt.Sprintf(common.SuccessMsgAppointmentCancelled,
		date, startTime, endTime,
		response.Professional.FirstName, response.Professional.LastName,
		response.Appointment.CancellationReason)

	h.sendMessage(chatID, text)

	// Notify professional about cancellation
	h.notificationService.NotifyProfessionalCancellation(response)
	h.ShowDashboard(chatID)
}
