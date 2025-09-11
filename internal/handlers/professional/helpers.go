package professional

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	"booking_client/internal/util"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Professional-specific helper functions

// sendError sends an error message to the user (ProfessionalHandler version)
func (h *ProfessionalHandler) sendError(chatID int64, message string, err error) {
	text := fmt.Sprintf(message, err)
	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send error message")
	}
}

// sendMessage sends a simple message to the user (ProfessionalHandler version)
func (h *ProfessionalHandler) sendMessage(chatID int64, text string) {
	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send message")
	}
}

// sendMessageWithID sends a message and returns the message ID (ProfessionalHandler version)
func (h *ProfessionalHandler) sendMessageWithID(chatID int64, text string) (int, error) {
	return h.bot.SendMessageWithID(chatID, text)
}

// sendMessageWithKeyboard sends a message with inline keyboard (ProfessionalHandler version)
func (h *ProfessionalHandler) sendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send message with keyboard")
	}
}

// sendMessageWithKeyboardAndID sends a message with keyboard and returns the message ID (ProfessionalHandler version)
func (h *ProfessionalHandler) sendMessageWithKeyboardAndID(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) (int, error) {
	return h.bot.SendMessageWithKeyboardAndID(chatID, text, keyboard)
}

// editMessage edits the last message sent to the user (ProfessionalHandler version)
func (h *ProfessionalHandler) editMessage(chatID int64, messageID int, text string) {
	if err := h.bot.EditMessage(chatID, messageID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to edit message")
	}
}

// editMessageWithKeyboard edits the last message with keyboard (ProfessionalHandler version)
func (h *ProfessionalHandler) editMessageWithKeyboard(chatID int64, messageID int, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	if err := h.bot.EditMessageWithKeyboard(chatID, messageID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to edit message with keyboard")
	}
}

// validateUserState checks if user is in a valid state for the given action (ProfessionalHandler version)
func (h *ProfessionalHandler) validateUserState(chatID int64, allowedStates []string) (*models.User, bool) {
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

// clearUnavailableState clears all unavailable-related state from user
func (h *ProfessionalHandler) clearUnavailableState(user *models.User) {
	user.State = models.StateNone
	user.SelectedDate = ""
	user.SelectedUnavailableStartTime = ""
	user.SelectedUnavailableEndTime = ""
	user.SelectedUnavailableDescription = ""
	user.SelectedAppointmentID = ""
}

// createProfessionalDashboardKeyboard creates the professional dashboard keyboard
func (h *ProfessionalHandler) createProfessionalDashboardKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnMyTimetable, "professional_timetable"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnPendingAppointments, "professional_pending_appointments"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnUpcomingAppointments, "professional_upcoming_appointments"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnSetUnavailable, "set_unavailable"),
		),
	)
}

// createProfessionalAppointmentsKeyboard creates a keyboard for professional appointment management
func (h *ProfessionalHandler) createProfessionalAppointmentsKeyboard(appointments []models.ProfessionalAppointment, showConfirm bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for index, apt := range appointments {
		if showConfirm {
			// For pending appointments - show both confirm and cancel
			confirmButton := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(common.BtnConfirmAppointmentProf, index+1),
				fmt.Sprintf("confirm_appointment_%s", apt.ID),
			)
			cancelButton := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(common.BtnCancelAppointmentProf, index+1),
				fmt.Sprintf("cancel_prof_appt_%s", apt.ID),
			)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(confirmButton, cancelButton))
		} else {
			// For upcoming appointments - show only cancel
			cancelButton := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(common.BtnCancelAppointmentProfAlt, index+1),
				fmt.Sprintf("cancel_prof_appt_%s", apt.ID),
			)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))
		}
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnBackToDashboard, "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// createUnavailableDateKeyboard creates a keyboard for unavailable date selection
func (h *ProfessionalHandler) createUnavailableDateKeyboard(currentDate time.Time) tgbotapi.InlineKeyboardMarkup {
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
			fmt.Sprintf("select_unavailable_date_%s", dateStr),
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

	// Add navigation buttons
	prevButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnPreviousUnavailableMonth, "prev_unavailable_month")
	nextButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnNextUnavailableMonth, "next_unavailable_month")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(prevButton, nextButton))

	// Add cancel button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelUnavailable, "cancel_unavailable")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// createUnavailableStartTimeKeyboard creates a keyboard for unavailable start time selection
func (h *ProfessionalHandler) createUnavailableStartTimeKeyboard(availability *models.ProfessionalAvailabilityResponse) tgbotapi.InlineKeyboardMarkup {
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
			fmt.Sprintf("select_unavailable_start_%s", timeDisplay),
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

	// Add cancel button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelUnavailable, "cancel_unavailable")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// createUnavailableEndTimeKeyboard creates a keyboard for unavailable end time selection
func (h *ProfessionalHandler) createUnavailableEndTimeKeyboard(startTime string, availability *models.ProfessionalAvailabilityResponse) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	// Find available slots after the selected start time
	for _, slot := range availability.Slots {
		slotStart, err := time.Parse(time.RFC3339, slot.StartTime)
		if err != nil {
			continue
		}
		slotStartLocal := slotStart.In(util.GetAppTimezone())

		// Only show slots that are:
		// 1. Available (not blocked by existing appointments/unavailable periods)
		// 2. At least 1 hour after the selected start time
		// 3. Before the first unavailable slot
		slotTimeStr := slotStartLocal.Format("15:04")
		if slot.Available && slotTimeStr >= startTime {
			// Use the end time of the slot as the end time option
			slotEnd, err := time.Parse(time.RFC3339, slot.EndTime)
			if err != nil {
				continue
			}
			slotEndLocal := slotEnd.In(util.GetAppTimezone())
			timeDisplay := slotEndLocal.Format("15:04")
			button := tgbotapi.NewInlineKeyboardButtonData(
				timeDisplay,
				fmt.Sprintf("select_unavailable_end_%s", timeDisplay),
			)
			currentRow = append(currentRow, button)

			if len(currentRow) == common.TimeSlotsPerRow {
				rows = append(rows, currentRow)
				currentRow = []tgbotapi.InlineKeyboardButton{}
			}
		} else if !slot.Available && slotTimeStr > startTime {
			// Stop at the first unavailable slot after start time
			break
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Add cancel button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelUnavailable, "cancel_unavailable")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// createUpcomingAppointmentsDateKeyboard creates a keyboard for upcoming appointments date selection
func (h *ProfessionalHandler) createUpcomingAppointmentsDateKeyboard(dates []string, currentMonth string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	// Add month navigation buttons at the top
	// Only show previous button if not current month
	currentTime := time.Now()
	currentMonthTime, _ := time.Parse("2006-01", currentMonth)
	isCurrentMonth := currentMonthTime.Year() == currentTime.Year() && currentMonthTime.Month() == currentTime.Month()

	var navButtons []tgbotapi.InlineKeyboardButton
	if !isCurrentMonth {
		prevButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnPreviousMonth, "prev_upcoming_month_"+currentMonth)
		navButtons = append(navButtons, prevButton)
	}
	nextButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnNextMonth, "next_upcoming_month_"+currentMonth)
	navButtons = append(navButtons, nextButton)

	if len(navButtons) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(navButtons...))
	}

	for _, dateStr := range dates {
		// Parse date to format it nicely
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		// Format as "Sep 10" or "Today" if it's today
		var displayText string
		if date.Year() == time.Now().Year() && date.Month() == time.Now().Month() && date.Day() == time.Now().Day() {
			displayText = "Today"
		} else {
			displayText = date.Format("Jan 02")
		}

		button := tgbotapi.NewInlineKeyboardButtonData(
			displayText,
			fmt.Sprintf("select_upcoming_date_%s", dateStr),
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

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnBackToDashboard, "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// createTimetableKeyboard creates a keyboard for timetable with day navigation and appointment actions
func (h *ProfessionalHandler) createTimetableKeyboard(dateStr string, appointments []models.TimetableAppointment) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Add day navigation buttons
	currentDate, _ := time.Parse("2006-01-02", dateStr)
	today := time.Now()
	isToday := currentDate.Year() == today.Year() && currentDate.Month() == today.Month() && currentDate.Day() == today.Day()

	var navButtons []tgbotapi.InlineKeyboardButton
	if !isToday {
		prevButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnPreviousTimetableDay, "prev_timetable_day_"+dateStr)
		navButtons = append(navButtons, prevButton)
	}
	nextButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnNextTimetableDay, "next_timetable_day_"+dateStr)
	navButtons = append(navButtons, nextButton)

	if len(navButtons) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(navButtons...))
	}

	// Add appointment cancel buttons
	for i, apt := range appointments {
		cancelButton := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf(common.BtnCancelTimetableSlot, i+1),
			fmt.Sprintf("cancel_appointment_%s", apt.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnBackToDashboard, "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
