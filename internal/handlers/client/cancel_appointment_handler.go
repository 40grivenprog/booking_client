package client

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	apiService "booking_client/internal/services/api_service"
	"fmt"
)

// HandleCancelAppointment starts the appointment cancellation process
func (h *ClientHandler) HandleCancelAppointment(chatID int64, appointmentID string) {
	// Store appointment ID and ask for cancellation reason
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.State = models.StateWaitingForCancellationReason
	user.SelectedAppointmentID = appointmentID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	id, err := h.bot.SendMessageWithID(chatID, common.UIMsgCancellationReason)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
}

// HandleCancellationReason handles the cancellation reason input
func (h *ClientHandler) HandleCancellationReason(chatID int64, reason string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	appointmentID := user.SelectedAppointmentID

	// Cancel the appointment
	req := &apiService.CancelAppointmentRequest{
		CancellationReason: reason,
	}

	response, err := h.apiService.CancelClientAppointment(user.ID, appointmentID, req)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToCancelAppointment, err)
		return
	}

	// Clear state
	user.State = models.StateNone
	user.SelectedAppointmentID = ""
	h.apiService.GetUserRepository().SetUser(chatID, user)

	date, startTime, endTime := common.FormatAppointmentTime(response.Appointment.StartTime, response.Appointment.EndTime)
	text := fmt.Sprintf(common.SuccessMsgAppointmentCancelled,
		date, startTime, endTime,
		response.Professional.FirstName, response.Professional.LastName,
		response.Appointment.CancellationReason)

	h.sendMessage(chatID, text)

	// Notify professional about cancellation
	h.notificationService.NotifyProfessionalCancellation(response)
	h.ShowDashboard(chatID)
}
