package client

import (
	"context"
	"time"

	"booking_client/internal/handlers/common"
)

// HandlePendingAppointments shows pending appointments
func (h *ClientHandler) HandlePendingAppointments(ctx context.Context, chatID int64, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	// Delete message for dashboard
	go func() {
		time.Sleep(500 * time.Millisecond)
		h.bot.DeleteMessage(chatID, messageID)
	}()

	appointments, err := h.apiService.GetClientAppointments(ctx, user.ID, "pending")
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
		h.ShowDashboard(ctx, chatID, 0)
		return
	}

	text := common.UIMsgPendingAppointments
	for index, apt := range appointments.Appointments {
		text += common.FormatAppointmentDetails(&apt, index)
	}

	keyboard := h.createAppointmentsKeyboard(appointments.Appointments, common.BtnCancelAppointment)
	err = h.bot.SendMessageWithKeyboard(chatID, text, keyboard)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
}
