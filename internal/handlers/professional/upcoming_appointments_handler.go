package professional

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	"fmt"
	"time"
)

// HandleUpcomingAppointments shows upcoming appointments for professionals
func (h *ProfessionalHandler) HandleUpcomingAppointments(chatID int64) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	h.showUpcomingAppointmentsDatePicker(chatID, user)
}

// showUpcomingAppointmentsDatePicker shows date picker for upcoming appointments
func (h *ProfessionalHandler) showUpcomingAppointmentsDatePicker(chatID int64, user *models.User, month ...string) {
	var targetMonth string
	if len(month) > 0 {
		targetMonth = month[0]
	} else {
		targetMonth = time.Now().Format("2006-01")
	}

	appointmentDates, err := h.apiService.GetProfessionalAppointmentDates(user.ID, targetMonth)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToLoadAppointments, err)
		return
	}

	text := fmt.Sprintf(common.UIMsgSelectUpcomingAppointmentsDate, targetMonth)
	keyboard := h.createUpcomingAppointmentsDateKeyboard(appointmentDates.Dates, targetMonth)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUpcomingAppointmentsMonthNavigation handles month navigation for upcoming appointments
func (h *ProfessionalHandler) HandleUpcomingAppointmentsMonthNavigation(chatID int64, monthStr string, direction string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	currentMonth, err := time.Parse("2006-01", monthStr)
	if err != nil {
		h.sendMessage(chatID, common.ErrorMsgInvalidDateFormat)
		return
	}

	var newMonth time.Time
	if direction == common.DirectionPrev {
		newMonth = currentMonth.AddDate(0, -1, 0)
	} else {
		newMonth = currentMonth.AddDate(0, 1, 0)
	}

	h.showUpcomingAppointmentsDatePicker(chatID, user, newMonth.Format("2006-01"))
}

// HandleUpcomingAppointmentsDateSelection handles date selection from upcoming appointments picker
func (h *ProfessionalHandler) HandleUpcomingAppointmentsDateSelection(chatID int64, dateStr string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetProfessionalAppointmentsByDate(user.ID, "confirmed", dateStr)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToLoadAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		h.sendMessage(chatID, common.UIMsgNoUpcomingAppointments)
		h.ShowDashboard(chatID, user)
		return
	}

	text := common.UIMsgUpcomingAppointments
	for index, apt := range appointments.Appointments {
		text += common.FormatProfessionalAppointmentDetails(&apt, index)
	}

	keyboard := h.createProfessionalAppointmentsKeyboard(appointments.Appointments, false)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}
