package client

import (
	"booking_client/internal/common"
	"booking_client/internal/models"
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// sendError sends an error message to the user
func (h *ClientHandler) sendError(ctx context.Context, chatID int64, message string, err error) {
	text := fmt.Sprintf(message, err)
	if err := h.bot.SendMessage(chatID, text); err != nil {
		logger := common.GetLogger(ctx)
		logger.Error().Err(err).Msg("Failed to send error message")
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
func (h *ClientHandler) validateUserState(ctx context.Context, chatID int64, allowedStates []string) (*models.User, bool) {
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

// clearBookingState clears all booking-related state from user
func (h *ClientHandler) clearBookingState(user *models.User) {
	user.State = models.StateNone
	user.SelectedProfessionalID = ""
	user.SelectedDate = ""
	user.SelectedTime = ""
	user.SelectedAppointmentID = ""
}

// Keyboard wrapper methods for backward compatibility
func (h *ClientHandler) createDateKeyboard(currentDate time.Time) tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateDateKeyboard(currentDate)
}

func (h *ClientHandler) createTimeKeyboard(availability *models.ProfessionalAvailabilityResponse) tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateTimeKeyboard(availability)
}

func (h *ClientHandler) createProfessionalsKeyboard(professionals []models.User) tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateProfessionalsKeyboard(professionals)
}

func (h *ClientHandler) createAppointmentsKeyboard(appointments []models.ClientAppointment, buttonPrefix string) tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateAppointmentsKeyboard(appointments, buttonPrefix)
}

func (h *ClientHandler) createDashboardKeyboard() tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateDashboardKeyboard()
}

func (h *ClientHandler) createRegistrationSuccessKeyboard() tgbotapi.InlineKeyboardMarkup {
	return h.keyboards.CreateRegistrationSuccessKeyboard()
}
