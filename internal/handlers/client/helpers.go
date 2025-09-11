package client

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	"booking_client/internal/util"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// sendError sends an error message to the user
func (h *ClientHandler) sendError(chatID int64, message string, err error) {
	text := fmt.Sprintf(message, err)
	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send error message")
	}
}

// sendMessage sends a simple message to the user
func (h *ClientHandler) sendMessage(chatID int64, text string) {
	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send message")
	}
}

// sendMessageWithID sends a message and returns the message ID
func (h *ClientHandler) sendMessageWithID(chatID int64, text string) (int, error) {
	return h.bot.SendMessageWithID(chatID, text)
}

// sendMessageWithKeyboard sends a message with inline keyboard
func (h *ClientHandler) sendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send message with keyboard")
	}
}

// sendMessageWithKeyboardAndID sends a message with keyboard and returns the message ID
func (h *ClientHandler) sendMessageWithKeyboardAndID(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) (int, error) {
	return h.bot.SendMessageWithKeyboardAndID(chatID, text, keyboard)
}

// editMessage edits the last message sent to the user
func (h *ClientHandler) editMessage(chatID int64, messageID int, text string) {
	if err := h.bot.EditMessage(chatID, messageID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to edit message")
	}
}

// editMessageWithKeyboard edits the last message with keyboard
func (h *ClientHandler) editMessageWithKeyboard(chatID int64, messageID int, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	if err := h.bot.EditMessageWithKeyboard(chatID, messageID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to edit message with keyboard")
	}
}

// validateUserState checks if user is in a valid state for the given action
func (h *ClientHandler) validateUserState(chatID int64, allowedStates []string) (*models.User, bool) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return nil, false
	}

	// Check if user is in an allowed state
	for _, allowedState := range allowedStates {
		if user.State == allowedState {
			return user, true
		}
	}

	// User is not in an allowed state
	h.sendMessage(chatID, common.ErrorMsgInvalidState)
	return nil, false
}

// clearBookingState clears all booking-related state from user
func (h *ClientHandler) clearBookingState(user *models.User) {
	user.State = models.StateNone
	user.SelectedProfessionalID = ""
	user.SelectedDate = ""
	user.SelectedTime = ""
	user.SelectedAppointmentID = ""
}

// createDateKeyboard creates a keyboard for date selection
func (h *ClientHandler) createDateKeyboard(currentDate time.Time) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	// Get first day of month and number of days
	firstDay := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	today := time.Now()

	for d := firstDay; d.Before(lastDay.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		// Skip past dates (compare only date part, not time)
		if d.Year() < today.Year() ||
			(d.Year() == today.Year() && d.Month() < today.Month()) ||
			(d.Year() == today.Year() && d.Month() == today.Month() && d.Day() < today.Day()) {
			continue
		}

		dateStr := d.Format("2006-01-02")
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d", d.Day()),
			fmt.Sprintf("select_date_%s", dateStr),
		)
		currentRow = append(currentRow, button)

		if len(currentRow) == common.DaysPerRow {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}
	currentMonth := currentDate.Format("2006-01")
	// Add navigation buttons
	prevButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnPreviousMonth, "prev_month_"+currentMonth)
	nextButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnNextMonth, "next_month_"+currentMonth)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(prevButton, nextButton))

	// Add cancel booking button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelBooking, "cancel_booking")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// createTimeKeyboard creates a keyboard for time slot selection
func (h *ClientHandler) createTimeKeyboard(availability *models.ProfessionalAvailabilityResponse) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	for _, slot := range availability.Slots {
		if !slot.Available {
			continue
		}

		// Parse the RFC3339 time and convert to local timezone for display
		startTime, err := time.Parse(time.RFC3339, slot.StartTime)
		if err != nil {
			h.logger.Error().Err(err).Str("time", slot.StartTime).Msg("Failed to parse time slot")
			continue
		}

		// Convert to local timezone for display
		localTime := startTime.In(util.GetAppTimezone())
		timeDisplay := localTime.Format("15:04")
		button := tgbotapi.NewInlineKeyboardButtonData(
			timeDisplay,
			fmt.Sprintf("select_time_%s", timeDisplay),
		)
		currentRow = append(currentRow, button)

		if len(currentRow) == common.TimeSlotsPerRow {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Add cancel booking button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelBooking, "cancel_booking")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// createProfessionalsKeyboard creates a keyboard for professional selection
func (h *ClientHandler) createProfessionalsKeyboard(professionals []models.User) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, prof := range professionals {
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("👨‍💼 %s %s", prof.FirstName, prof.LastName),
			fmt.Sprintf("select_professional_%s", prof.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	// Add cancel booking button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelBooking, "cancel_booking")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// createAppointmentsKeyboard creates a keyboard for appointment management
func (h *ClientHandler) createAppointmentsKeyboard(appointments []models.ClientAppointment, buttonPrefix string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for index, apt := range appointments {
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf(buttonPrefix, index+1),
			fmt.Sprintf("cancel_appointment_%s", apt.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnBackToDashboard, "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// createDashboardKeyboard creates the main dashboard keyboard
func (h *ClientHandler) createDashboardKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnBookAppointment, "book_appointment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnMyPendingAppointments, "pending_appointments"),
			tgbotapi.NewInlineKeyboardButtonData(common.BtnMyUpcomingAppointments, "upcoming_appointments"),
		),
	)
}

// createRegistrationSuccessKeyboard creates keyboard for successful registration
func (h *ClientHandler) createRegistrationSuccessKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnGoToDashboard, "back_to_dashboard"),
		),
	)
}
