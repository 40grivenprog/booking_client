package handlers

import (
	"booking_client/internal/models"
	"booking_client/internal/repository"
	"booking_client/pkg/telegram"

	"github.com/rs/zerolog"
)

func getUserOrSendError(repo *repository.UserRepository, bot *telegram.Bot, logger *zerolog.Logger, chatID int64) (*models.User, bool) {
	user, exists := repo.GetUser(chatID)
	if !exists || user == nil {
		text := "❌ User session not found. Please use /start to begin."
		if err := bot.SendMessage(chatID, text); err != nil {
			logger.Error().Err(err).Msg("Failed to send user not found message")
		}
		return nil, false
	}
	return user, true
}
