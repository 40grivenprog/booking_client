package professional

import (
	"context"
	"time"

	"booking_client/internal/handlers/common"
)

// HandlePendingAppointments shows pending appointments for professionals
func (h *ProfessionalHandler) HandlePendingAppointments(ctx context.Context, chatID int64, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Delete message for dashboard
	go func() {
		time.Sleep(1 * time.Second)
		h.bot.DeleteMessage(chatID, messageID)
	}()

	appointments, err := h.apiService.GetProfessionalAppointments(ctx, user.ID, "pending")
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToLoadPendingAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		id, err := h.sendMessageWithID(chatID, common.UIMsgNoPendingAppointments)
		if err != nil {
			h.sendError(ctx, chatID, common.ErrorMsgFailedToSendMessage, err)
			return
		}
		user.LastMessageID = &id
		user.MessagesToDelete = append(user.MessagesToDelete, &id)
		h.apiService.GetUserRepository().SetUser(chatID, user)
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
