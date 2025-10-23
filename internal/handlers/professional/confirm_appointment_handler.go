package professional

import (
	"booking_client/internal/handlers/common"
	apiService "booking_client/internal/services/api_service"
	"context"
)

// HandleConfirmAppointment handles professional appointment confirmation
func (h *ProfessionalHandler) HandleConfirmAppointment(ctx context.Context, chatID int64, appointmentID string, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.LastMessageID = &messageID
	user.MessagesToDelete = append(user.MessagesToDelete, &messageID)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Confirm the appointment
	req := &apiService.ConfirmAppointmentRequest{}

	response, err := h.apiService.ConfirmProfessionalAppointment(ctx, user.ID, appointmentID, req)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToConfirmAppointment, err)
		return
	}

	// Build success message
	date, startTime, endTime := common.FormatAppointmentTime(response.Appointment.StartTime, response.Appointment.EndTime)
	text := common.NewSuccessMessage("appointment_confirmed").
		WithData("date", date).
		WithData("start_time", startTime).
		WithData("end_time", endTime).
		WithData("client_first_name", response.Client.FirstName).
		WithData("client_last_name", response.Client.LastName).
		Build()

	err = h.bot.SendMessage(chatID, text)
	if err == nil {
		h.apiService.GetUserRepository().SetUser(chatID, user)
	}
	h.ShowDashboard(ctx, chatID, user, 0)

	// Notify client about confirmation
	h.notificationService.NotifyClientAppointmentConfirmation(response)
}
