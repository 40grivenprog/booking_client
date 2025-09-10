package handlers

import (
	"fmt"

	"booking_client/internal/models"
	"booking_client/pkg/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// NotificationService handles all notification-related operations
type NotificationService struct {
	bot    *telegram.Bot
	logger *zerolog.Logger
}

// NewNotificationService creates a new notification service
func NewNotificationService(bot *telegram.Bot, logger *zerolog.Logger) *NotificationService {
	return &NotificationService{
		bot:    bot,
		logger: logger,
	}
}

// NotifyProfessionalNewAppointment sends notification to professional about new appointment
func (ns *NotificationService) NotifyProfessionalNewAppointment(appointment *models.CreateAppointmentResponse) {
	if appointment.Professional.ChatID == 0 {
		return // No chat ID for professional
	}

	date, startTime, endTime := formatAppointmentTime(appointment.Appointment.StartTime, appointment.Appointment.EndTime)

	text := fmt.Sprintf(UIMsgNewAppointmentRequest,
		appointment.Client.FirstName, appointment.Client.LastName,
		date, startTime, endTime,
		appointment.Appointment.Description)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtnConfirmAppointment, fmt.Sprintf("confirm_appointment_%s", appointment.Appointment.ID)),
			tgbotapi.NewInlineKeyboardButtonData(BtnCancelAppointmentConfirm, fmt.Sprintf("cancel_appointment_%s", appointment.Appointment.ID)),
			tgbotapi.NewInlineKeyboardButtonData(BtnBackToDashboard, "back_to_dashboard"),
		),
	)

	if err := ns.bot.SendMessageWithKeyboard(appointment.Professional.ChatID, text, keyboard); err != nil {
		ns.logger.Error().Err(err).Msg("Failed to send professional notification")
	}
}

// NotifyProfessionalCancellation sends notification to professional about appointment cancellation
func (ns *NotificationService) NotifyProfessionalCancellation(response *models.CancelClientAppointmentResponse) {
	if response.Professional.ChatID == nil || *response.Professional.ChatID == 0 {
		return // No chat ID for professional
	}

	date, startTime, endTime := formatAppointmentTime(response.Appointment.StartTime, response.Appointment.EndTime)

	text := fmt.Sprintf(UIMsgAppointmentCancelled,
		response.Client.FirstName, response.Client.LastName,
		date, startTime, endTime,
		response.Appointment.CancellationReason)

	if err := ns.bot.SendMessage(*response.Professional.ChatID, text); err != nil {
		ns.logger.Error().Err(err).Msg("Failed to send professional cancellation notification")
	}
}

// NotifyClientAppointmentConfirmation sends notification to client about appointment confirmation
func (ns *NotificationService) NotifyClientAppointmentConfirmation(response *models.ConfirmProfessionalAppointmentResponse) {
	if response.Client.ChatID == nil || *response.Client.ChatID == 0 {
		return // No chat ID for client
	}

	date, startTime, endTime := formatAppointmentTime(response.Appointment.StartTime, response.Appointment.EndTime)

	text := fmt.Sprintf(UIMsgAppointmentConfirmed,
		date, startTime, endTime,
		response.Professional.FirstName, response.Professional.LastName)

	if err := ns.bot.SendMessage(*response.Client.ChatID, text); err != nil {
		ns.logger.Error().Err(err).Msg("Failed to send client confirmation notification")
	}
}

// NotifyClientProfessionalCancellation sends notification to client about appointment cancellation by professional
func (ns *NotificationService) NotifyClientProfessionalCancellation(response *models.CancelProfessionalAppointmentResponse) {
	if response.Client.ChatID == nil || *response.Client.ChatID == 0 {
		return // No chat ID for client
	}

	date, startTime, endTime := formatAppointmentTime(response.Appointment.StartTime, response.Appointment.EndTime)

	text := fmt.Sprintf(UIMsgAppointmentCancelledByProfessional,
		date, startTime, endTime,
		response.Professional.FirstName, response.Professional.LastName,
		response.Appointment.CancellationReason)

	if err := ns.bot.SendMessage(*response.Client.ChatID, text); err != nil {
		ns.logger.Error().Err(err).Msg("Failed to send client cancellation notification")
	}
}
