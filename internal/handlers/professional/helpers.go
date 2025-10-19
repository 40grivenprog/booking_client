package professional

import (
	"booking_client/internal/common"
	"booking_client/internal/models"
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Professional-specific helper functions

// sendError sends an error message to the user (ProfessionalHandler version)
func (h *ProfessionalHandler) sendError(ctx context.Context, chatID int64, message string, err error) {
	text := fmt.Sprintf(message, err)
	if err := h.bot.SendMessage(chatID, text); err != nil {
		logger := common.GetLogger(ctx)
		logger.Error().Err(err).Msg("Failed to send error message")
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
func (h *ProfessionalHandler) validateUserState(ctx context.Context, chatID int64, allowedStates []string) (*models.User, bool) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, common.GetLogger(ctx), chatID)
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
// Keyboard wrapper methods for backward compatibility
func (h *ProfessionalHandler) createProfessionalDashboardKeyboard() tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateProfessionalDashboardKeyboard()
}

func (h *ProfessionalHandler) createProfessionalAppointmentsKeyboard(appointments []models.ProfessionalAppointment, showConfirm bool) tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateProfessionalAppointmentsKeyboard(appointments, showConfirm)
}

func (h *ProfessionalHandler) createUnavailableDateKeyboard(currentDate time.Time) tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateUnavailableDateKeyboard(currentDate)
}

func (h *ProfessionalHandler) createUnavailableStartTimeKeyboard(availability *models.ProfessionalAvailabilityResponse) tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateUnavailableStartTimeKeyboard(availability)
}

func (h *ProfessionalHandler) createUnavailableEndTimeKeyboard(startTime string, availability *models.ProfessionalAvailabilityResponse) tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateUnavailableEndTimeKeyboard(startTime, availability)
}

func (h *ProfessionalHandler) createUpcomingAppointmentsDateKeyboard(dates []string, currentMonth string) tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateUpcomingAppointmentsDateKeyboard(dates, currentMonth)
}

func (h *ProfessionalHandler) createTimetableKeyboard(dateStr string, appointments []models.TimetableAppointment) tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateTimetableKeyboard(dateStr, appointments)
}
