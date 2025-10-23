package professional

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	apiService "booking_client/internal/services/api_service"
	"context"
)

// HandleCancelAppointment starts the professional appointment cancellation process
func (h *ProfessionalHandler) HandleCancelAppointment(ctx context.Context, chatID int64, appointmentID string, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Store appointment ID and ask for cancellation reason
	user.State = models.StateWaitingForCancellationReason
	user.SelectedAppointmentID = appointmentID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	h.sendMessage(chatID, common.UIMsgCancellationReason)
}

// HandleCancellationReason handles the professional cancellation reason input
func (h *ProfessionalHandler) HandleCancellationReason(ctx context.Context, chatID int64, reason string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	appointmentID := user.SelectedAppointmentID

	// Cancel the appointment
	req := &apiService.CancelAppointmentRequest{
		CancellationReason: reason,
	}

	response, err := h.apiService.CancelProfessionalAppointment(ctx, user.ID, appointmentID, req)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToCancelAppointment, err)
		return
	}

	// Clear state
	user.State = models.StateNone
	user.SelectedAppointmentID = ""
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Build success message
	date, startTime, endTime := common.FormatAppointmentTime(response.Appointment.StartTime, response.Appointment.EndTime)
	text := common.NewSuccessMessage("appointment_cancelled").
		WithData("date", date).
		WithData("start_time", startTime).
		WithData("end_time", endTime).
		WithData("first_name", response.Client.FirstName).
		WithData("last_name", response.Client.LastName).
		WithData("reason", response.Appointment.CancellationReason).
		Build()

	h.sendMessage(chatID, text)

	// Notify client about cancellation
	h.notificationService.NotifyClientProfessionalCancellation(response)
	h.ShowDashboard(ctx, chatID, user, 0)
}
