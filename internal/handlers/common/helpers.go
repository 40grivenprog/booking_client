package common

import (
	"fmt"

	"booking_client/internal/models"
	"booking_client/internal/repository"
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

// formatAppointmentDetails formats appointment details for display
func FormatAppointmentDetails(apt *models.ClientAppointment, index int) string {
	date, startTime, endTime := FormatAppointmentTime(apt.StartTime, apt.EndTime)
	return fmt.Sprintf("✍️ Appointment #%d:\n📅 %s\n🕐 %s - %s\n👨‍💼 %s %s\n📝 %s\n\n",
		index+1, date, startTime, endTime,
		apt.Professional.FirstName, apt.Professional.LastName,
		apt.Description)
}

// getUserOrSendError retrieves user from repository or sends error message
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

// formatProfessionalAppointmentDetails formats appointment details for professional display
func FormatProfessionalAppointmentDetails(apt *models.ProfessionalAppointment, index int) string {
	date, startTime, endTime := FormatAppointmentTime(apt.StartTime, apt.EndTime)
	return fmt.Sprintf("✍️ Appointment #%d:\n📅 %s\n🕐 %s - %s\n👤 Client: %s %s\n📝 %s\n\n",
		index+1, date, startTime, endTime,
		apt.Client.FirstName, apt.Client.LastName,
		apt.Description)
}
