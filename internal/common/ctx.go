package common

import (
	"booking_client/internal/models"
	"booking_client/internal/repository"
	"booking_client/pkg/telegram"
	"context"

	"github.com/rs/zerolog"
)

const (
	RequestIDKey string = "request_id"
	LoggerKey    string = "logger"
)

// Error messages
const (
	ErrorMsgInvalidState = "❌ Invalid state. Please use /start to begin."
)

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

func GetLogger(ctx context.Context) zerolog.Logger {
	if logger, ok := ctx.Value(LoggerKey).(zerolog.Logger); ok {
		return logger
	}
	return zerolog.Nop()
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func WithLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// GetUserOrSendError retrieves user from repository or sends error message
func GetUserOrSendError(userRepo *repository.UserRepository, bot *telegram.Bot, logger zerolog.Logger, chatID int64) (*models.User, bool) {
	user, exists := userRepo.GetUser(chatID)
	if !exists || user == nil {
		text := "❌ User session not found. Please use /start to begin."
		if err := bot.SendMessage(chatID, text); err != nil {
			logger.Error().Err(err).Msg("Failed to send user not found message")
		}
		return nil, false
	}
	return user, true
}
