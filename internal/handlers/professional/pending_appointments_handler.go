package professional

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	apiService "booking_client/internal/services/api_service"
	"time"
)

// HandlePendingAppointments shows pending appointments for professionals
func (h *ProfessionalHandler) HandlePendingAppointments(chatID int64) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetProfessionalAppointments(user.ID, "pending")
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToLoadPendingAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		h.sendMessage(chatID, common.UIMsgNoPendingAppointments)
		h.ShowDashboard(chatID, user)
		return
	}

	text := common.UIMsgPendingAppointments
	for index, apt := range appointments.Appointments {
		text += common.FormatProfessionalAppointmentDetails(&apt, index)
	}

	keyboard := h.createProfessionalAppointmentsKeyboard(appointments.Appointments, true)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleConfirmAppointment handles professional appointment confirmation
func (h *ProfessionalHandler) HandleConfirmAppointment(chatID int64, appointmentID string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Confirm the appointment
	req := &apiService.ConfirmAppointmentRequest{}

	response, err := h.apiService.ConfirmProfessionalAppointment(user.ID, appointmentID, req)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToConfirmAppointment, err)
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

	// Wait a moment then show dashboard (disappearing message effect)
	go func() {
		time.Sleep(300 * time.Millisecond)
		if user.LastMessageID != nil {
			h.bot.DeleteMessage(chatID, *user.LastMessageID)
		}
	}()
	go func() {
		time.Sleep(3 * time.Second)
		h.bot.DeleteMessage(chatID, confirmedMessageID)
	}()
}

// HandleCancelAppointment starts the professional appointment cancellation process
func (h *ProfessionalHandler) HandleCancelAppointment(chatID int64, appointmentID string) {
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
func (h *ProfessionalHandler) HandleCancellationReason(chatID int64, reason string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	appointmentID := user.SelectedAppointmentID

	// Cancel the appointment
	req := &apiService.CancelAppointmentRequest{
		CancellationReason: reason,
	}

	response, err := h.apiService.CancelProfessionalAppointment(user.ID, appointmentID, req)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToCancelAppointment, err)
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
	h.ShowDashboard(chatID, user)
}
