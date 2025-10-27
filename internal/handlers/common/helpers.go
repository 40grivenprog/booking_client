package common

import (
	"booking_client/internal/models"
	"booking_client/internal/repository"
	"booking_client/internal/schemas"
	"booking_client/pkg/telegram"

	"github.com/rs/zerolog"
)

// formatAppointmentTime formats appointment time for display
func FormatAppointmentTime(startTime, endTime string) (string, string, string) {
	start := startTime[:10]           // Extract date
	startTimeOnly := startTime[11:16] // Extract start time
	endTimeOnly := endTime[11:16]     // Extract end time
	return start, startTimeOnly, endTimeOnly
}

// FormatAppointmentDetails formats appointment details for client display
// Deprecated: Use NewClientAppointmentMessage(apt, index).ForClient() instead
func FormatAppointmentDetails(apt *schemas.ClientAppointment, index int) string {
	return NewClientAppointmentMessage(apt, index).ForClient()
}

// GetUserOrSendError retrieves user from repository or sends error message
func GetUserOrSendError(userRepo *repository.UserRepository, bot *telegram.Bot, logger *zerolog.Logger, chatID int64) (*models.User, bool) {
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

// FormatProfessionalAppointmentDetails formats appointment details for professional display
// Deprecated: Use NewProfessionalAppointmentMessage(apt, index).ForProfessional() instead
func FormatProfessionalAppointmentDetails(apt *schemas.ProfessionalAppointment, index int) string {
	return NewProfessionalAppointmentMessage(apt, index).ForProfessional()
}
