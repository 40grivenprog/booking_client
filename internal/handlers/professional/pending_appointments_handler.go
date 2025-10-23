package professional

import (
	"context"

	"booking_client/internal/handlers/common"
)

// HandlePendingAppointments shows pending appointments for professionals
func (h *ProfessionalHandler) HandlePendingAppointments(ctx context.Context, chatID int64, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetProfessionalAppointments(ctx, user.ID, "pending")
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToLoadPendingAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		h.sendMessage(chatID, common.UIMsgNoPendingAppointments)
		h.ShowDashboard(ctx, chatID, user, 0)
		return
	}

	text := common.UIMsgPendingAppointments
	for index, apt := range appointments.Appointments {
		text += common.FormatProfessionalAppointmentDetails(&apt, index)
	}

	keyboard := h.createProfessionalAppointmentsKeyboard(appointments.Appointments, true)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}
