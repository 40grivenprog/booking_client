package client

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	"booking_client/internal/schema"
	"booking_client/internal/util"
	"fmt"
	"time"
)

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
	if user.LastMessageID != nil {
		h.bot.DeleteMessage(chatID, *user.LastMessageID)
		user.LastMessageID = nil
	}
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
}

// HandleProfessionalSelection handles when user selects a professional
func (h *ClientHandler) HandleProfessionalSelection(chatID int64, professionalID string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	user.State = models.StateWaitingForDateSelection
	user.SelectedProfessionalID = professionalID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Show current month dates
	h.showDateSelection(user, time.Now())
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
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(*user.ChatID, user)
}

// HandleDateSelection handles when user selects a date
func (h *ClientHandler) HandleDateSelection(chatID int64, date string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
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
}

// HandleUpcomingAppointmentsMonthNavigation handles month navigation for upcoming appointments
func (h *ClientHandler) HandleBookAppointmentsMonthNavigation(chatID int64, monthStr string, direction string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Parse current month
	currentMonth, err := time.Parse("2006-01", monthStr)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgInvalidDateFormat, err)
		return
	}

	// Calculate new month based on direction
	var newMonth time.Time
	if direction == common.DirectionPrev {
		newMonth = currentMonth.AddDate(0, -1, 0)
	} else {
		newMonth = currentMonth.AddDate(0, 1, 0)
	}

	newMonthStr := newMonth.Format("2006-01")
	newMonthTime, err := time.Parse("2006-01", newMonthStr)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgInvalidDateFormat, err)
		return
	}
	h.bot.DeleteMessage(chatID, *user.LastMessageID)
	user.LastMessageID = nil
	h.apiService.GetUserRepository().SetUser(chatID, user)
	h.showDateSelection(user, newMonthTime)
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
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
}

// HandleTimeSelection handles when user selects a time slot
func (h *ClientHandler) HandleTimeSelection(chatID int64, startTime string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
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
	id, err := h.bot.SendMessageWithID(chatID, common.ErrorMsgBookingCancelled)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	for _, messageID := range user.MessagesToDelete {
		h.bot.DeleteMessage(chatID, *messageID)
	}
	user.MessagesToDelete = nil
	user.LastMessageID = &id
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
	h.ShowDashboard(chatID)
}
