package client

import (
	"context"

	"booking_client/internal/handlers/common"
)

// HandleUpcomingAppointments shows upcoming appointments
func (h *ClientHandler) HandleUpcomingAppointments(ctx context.Context, chatID int64) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetClientAppointments(ctx, user.ID, "confirmed")
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToLoadUpcomingAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		id, err := h.sendMessageWithID(chatID, common.UIMsgNoUpcomingAppointments)
		if err != nil {
			h.sendError(ctx, chatID, common.ErrorMsgFailedToSendMessage, err)
			return
		}
		user.LastMessageID = &id
		user.MessagesToDelete = append(user.MessagesToDelete, &id)
		h.apiService.GetUserRepository().SetUser(chatID, user)
		h.ShowDashboard(ctx, chatID)
		return
	}

	text := common.UIMsgUpcomingAppointments
	for index, apt := range appointments.Appointments {
		text += common.FormatAppointmentDetails(&apt, index)
	}

	keyboard := h.createAppointmentsKeyboard(appointments.Appointments, common.BtnCancelAppointment)
	id, err := h.bot.SendMessageWithKeyboardAndID(chatID, text, keyboard)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
}
