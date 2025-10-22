package professional

import (
	"booking_client/internal/handlers/common"
	apiService "booking_client/internal/services/api_service"
	"context"
	"time"
)

// HandleConfirmAppointment handles professional appointment confirmation
func (h *ProfessionalHandler) HandleConfirmAppointment(ctx context.Context, chatID int64, appointmentID string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

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

	confirmedMessageID, err := h.sendMessageWithID(chatID, text)
	if err == nil {
		h.apiService.GetUserRepository().SetUser(chatID, user)
	}

	// Notify client about confirmation
	h.notificationService.NotifyClientAppointmentConfirmation(response)

	go func() {
		time.Sleep(3 * time.Second)
		h.bot.DeleteMessage(chatID, confirmedMessageID)
	}()
}
