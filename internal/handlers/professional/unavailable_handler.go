package professional

import (
	"context"
	"fmt"
	"time"

	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	"booking_client/internal/schemas"
	apiService "booking_client/internal/services/api_service"
	"booking_client/internal/util"
)

// HandleSetUnavailable starts the unavailable appointment setting process
func (h *ProfessionalHandler) HandleSetUnavailable(ctx context.Context, chatID int64, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	// Delete message for dashboard
	go func() {
		time.Sleep(1 * time.Second)
		h.bot.DeleteMessage(chatID, messageID)
	}()

	// Set state for unavailable appointment
	user.State = models.StateWaitingForUnavailableDateSelection
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Show current month dates
	h.showUnavailableDateSelection(chatID, time.Now())
}

// showUnavailableDateSelection shows available dates for the current month
func (h *ProfessionalHandler) showUnavailableDateSelection(chatID int64, currentDate time.Time) {
	text := fmt.Sprintf(common.UIMsgSelectUnavailableDate, currentDate.Month(), currentDate.Year())
	keyboard := h.createUnavailableDateKeyboard(currentDate)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUnavailableDateSelection handles when user selects a date for unavailable time
func (h *ProfessionalHandler) HandleUnavailableDateSelection(ctx context.Context, chatID int64, date string, messageID int) {
	user, valid := h.validateUserState(ctx, chatID, []string{
		models.StateWaitingForUnavailableDateSelection,
	})
	if !valid {
		return
	}

	user.State = models.StateWaitingForUnavailableStartTime
	user.SelectedDate = date
	user.LastMessageID = &messageID
	user.MessagesToDelete = append(user.MessagesToDelete, &messageID)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Get availability for selected date to show time slots
	availability, err := h.apiService.GetProfessionalAvailability(ctx, user.ID, date)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToLoadAvailability, err)
		return
	}

	h.showUnavailableStartTimeSelection(chatID, availability)
}

// showUnavailableStartTimeSelection shows available time slots for start time
func (h *ProfessionalHandler) showUnavailableStartTimeSelection(chatID int64, availability *schemas.ProfessionalAvailabilityResponse) {
	text := fmt.Sprintf(common.UIMsgSelectUnavailableStartTime, availability.Date)
	keyboard := h.createUnavailableStartTimeKeyboard(availability)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUnavailableStartTimeSelection handles when user selects start time for unavailable period
func (h *ProfessionalHandler) HandleUnavailableStartTimeSelection(ctx context.Context, chatID int64, startTime string, messageID int) {
	user, valid := h.validateUserState(ctx, chatID, []string{
		models.StateWaitingForUnavailableStartTime,
	})
	if !valid {
		return
	}

	user.State = models.StateWaitingForUnavailableEndTime
	user.SelectedUnavailableStartTime = startTime // Store start time temporarily
	user.LastMessageID = &messageID
	user.MessagesToDelete = append(user.MessagesToDelete, &messageID)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Get availability for the selected date to determine available end times
	availability, err := h.apiService.GetProfessionalAvailability(ctx, user.ID, user.SelectedDate)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToLoadAvailability, err)
		return
	}

	h.showUnavailableEndTimeSelection(chatID, startTime, availability)
}

// showUnavailableEndTimeSelection shows available time slots for end time
func (h *ProfessionalHandler) showUnavailableEndTimeSelection(chatID int64, startTime string, availability *schemas.ProfessionalAvailabilityResponse) {
	text := fmt.Sprintf(common.UIMsgSelectUnavailableEndTime, startTime)

	// Find the first unavailable slot after the selected start time to show warning
	var firstUnavailableSlot *schemas.TimeSlot

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

		text += fmt.Sprintf("\n\n"+common.UIMsgUnavailableSlotWarning, unavailableStartLocal.Format("15:04"), slotDetails)
	}

	keyboard := h.createUnavailableEndTimeKeyboard(startTime, availability)

	// If no slots available, show a message
	if len(keyboard.InlineKeyboard) == 1 && len(keyboard.InlineKeyboard[0]) == 1 && keyboard.InlineKeyboard[0][0].Text == "❌ Cancel" {
		text += "\n\n" + common.UIMsgNoAvailableTimeSlots
	}

	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUnavailableEndTimeSelection handles when user selects end time for unavailable period
func (h *ProfessionalHandler) HandleUnavailableEndTimeSelection(ctx context.Context, chatID int64, endTime string, messageID int) {
	user, valid := h.validateUserState(ctx, chatID, []string{
		models.StateWaitingForUnavailableEndTime,
	})
	if !valid {
		return
	}

	// Store end time and ask for description
	user.State = models.StateWaitingForUnavailableDescription
	user.SelectedUnavailableEndTime = endTime
	user.LastMessageID = &messageID
	user.MessagesToDelete = append(user.MessagesToDelete, &messageID)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(common.UIMsgUnavailableDescription, user.SelectedDate, user.SelectedUnavailableStartTime, endTime)
	id, err := h.sendMessageWithID(chatID, text)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
}

// HandleUnavailableDescription handles when user provides description for unavailable period
func (h *ProfessionalHandler) HandleUnavailableDescription(ctx context.Context, chatID int64, description string, messageID int) {
	user, valid := h.validateUserState(ctx, chatID, []string{
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
		h.sendMessage(chatID, common.ErrorMsgInvalidDateFormat)
		return
	}

	// Combine date with times in application timezone
	startDateTime := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(),
		start.Hour(), start.Minute(), 0, 0, util.GetAppTimezone())
	endDateTime := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(),
		end.Hour(), end.Minute(), 0, 0, util.GetAppTimezone())

	// Create unavailable appointment request
	req := &apiService.CreateUnavailableAppointmentRequest{
		ProfessionalID: user.ID,
		StartAt:        startDateTime.Format(time.RFC3339),
		EndAt:          endDateTime.Format(time.RFC3339),
		Description:    description,
	}

	appointment, err := h.apiService.CreateUnavailableAppointment(ctx, req)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToCreateUnavailableAppointment, err)
		return
	}
	// Clear state
	h.clearUnavailableState(user)
	user.LastMessageID = &messageID
	user.MessagesToDelete = append(user.MessagesToDelete, &messageID)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(common.SuccessMsgUnavailablePeriodSet,
		date, start.Format("15:04"), end.Format("15:04"), appointment.Appointment.Description)

	h.sendMessage(chatID, text)
	h.ShowDashboard(ctx, chatID, user, 0)
}

// HandleCancelUnavailable cancels the unavailable appointment setting process
func (h *ProfessionalHandler) HandleCancelUnavailable(ctx context.Context, chatID int64, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Clear all unavailable-related state
	h.clearUnavailableState(user)
	user.LastMessageID = &messageID
	user.MessagesToDelete = append(user.MessagesToDelete, &messageID)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	id, err := h.bot.SendMessageWithID(chatID, common.ErrorMsgUnavailableCancelled)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.ShowDashboard(ctx, chatID, user, 0)
}

// HandleUnavailableMonthNavigation handles month navigation for unavailable appointments
func (h *ProfessionalHandler) HandleUnavailableMonthNavigation(ctx context.Context, chatID int64, month string, direction string, messageID int) {
	// Validate user state - only allow if waiting for unavailable date selection
	h.bot.DeleteMessage(chatID, messageID)
	_, valid := h.validateUserState(ctx, chatID, []string{
		models.StateWaitingForUnavailableDateSelection,
	})
	if !valid {
		return
	}
	h.bot.DeleteMessage(chatID, messageID)

	// Parse current month
	currentMonth, err := time.Parse("2006-01", month)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgInvalidDateFormat, err)
		return
	}

	// Calculate new month based on direction
	var newMonth time.Time
	if direction == common.DirectionPrev {
		newMonth = currentMonth.AddDate(0, -1, 0)
	} else {
		newMonth = currentMonth.AddDate(0, 1, 0)
	}

	// Show new month
	h.showUnavailableDateSelection(chatID, newMonth)
}
