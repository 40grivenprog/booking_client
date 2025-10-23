package client

import (
	"context"
	"time"

	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	apiService "booking_client/internal/services/api_service"
	"fmt"
)

// HandleCancelAppointment starts the appointment cancellation process
func (h *ClientHandler) HandleCancelAppointment(ctx context.Context, chatID int64, appointmentID string, messageID int) {
	// Store appointment ID and ask for cancellation reason
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.State = models.StateWaitingForCancellationReason
	user.SelectedAppointmentID = appointmentID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Delete message for dashboard
	go func() {
		time.Sleep(1 * time.Second)
		h.bot.DeleteMessage(chatID, messageID)
	}()

	id, err := h.bot.SendMessageWithID(chatID, common.UIMsgCancellationReason)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
}

// HandleCancellationReason handles the cancellation reason input
func (h *ClientHandler) HandleCancellationReason(ctx context.Context, chatID int64, reason string, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	appointmentID := user.SelectedAppointmentID

	// Cancel the appointment
	req := &apiService.CancelAppointmentRequest{
		CancellationReason: reason,
	}

	response, err := h.apiService.CancelClientAppointment(ctx, user.ID, appointmentID, req)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToCancelAppointment, err)
		return
	}

	// Clear state
	user.State = models.StateNone
	user.SelectedAppointmentID = ""
	user.LastMessageID = &messageID
	user.MessagesToDelete = append(user.MessagesToDelete, &messageID)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	date, startTime, endTime := common.FormatAppointmentTime(response.Appointment.StartTime, response.Appointment.EndTime)
	text := fmt.Sprintf(common.SuccessMsgAppointmentCancelled,
		date, startTime, endTime,
		response.Professional.FirstName, response.Professional.LastName,
		response.Appointment.CancellationReason)

	h.sendMessage(chatID, text)

	// Notify professional about cancellation
	h.notificationService.NotifyProfessionalCancellation(response)
	h.ShowDashboard(ctx, chatID, 0)
}
