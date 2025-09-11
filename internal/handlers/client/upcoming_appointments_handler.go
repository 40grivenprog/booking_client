package client

import "booking_client/internal/handlers/common"

// HandleUpcomingAppointments shows upcoming appointments
func (h *ClientHandler) HandleUpcomingAppointments(chatID int64) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetClientAppointments(user.ID, "confirmed")
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToLoadUpcomingAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		h.sendMessage(chatID, common.UIMsgNoUpcomingAppointments)
		h.ShowDashboard(chatID)
		return
	}

	text := common.UIMsgUpcomingAppointments
	for index, apt := range appointments.Appointments {
		text += common.FormatAppointmentDetails(&apt, index)
	}

	keyboard := h.createAppointmentsKeyboard(appointments.Appointments, common.BtnCancelAppointment)
	id, err := h.bot.SendMessageWithKeyboardAndID(chatID, text, keyboard)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
}
