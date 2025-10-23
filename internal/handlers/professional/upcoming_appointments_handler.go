package professional

import (
	"context"
	"fmt"
	"time"

	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
)

// HandleUpcomingAppointments shows upcoming appointments for professionals
func (h *ProfessionalHandler) HandleUpcomingAppointments(ctx context.Context, chatID int64, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	h.showUpcomingAppointmentsDatePicker(ctx, chatID, user)
}

// showUpcomingAppointmentsDatePicker shows date picker for upcoming appointments
func (h *ProfessionalHandler) showUpcomingAppointmentsDatePicker(ctx context.Context, chatID int64, user *models.User, month ...string) {
	var targetMonth string
	if len(month) > 0 {
		targetMonth = month[0]
	} else {
		targetMonth = time.Now().Format("2006-01")
	}

	appointmentDates, err := h.apiService.GetProfessionalAppointmentDates(ctx, user.ID, targetMonth)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToLoadAppointments, err)
		return
	}

	text := fmt.Sprintf(common.UIMsgSelectUpcomingAppointmentsDate, targetMonth)
	keyboard := h.createUpcomingAppointmentsDateKeyboard(appointmentDates.Dates, targetMonth)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleUpcomingAppointmentsMonthNavigation handles month navigation for upcoming appointments
func (h *ProfessionalHandler) HandleUpcomingAppointmentsMonthNavigation(ctx context.Context, chatID int64, monthStr string, direction string, messageID int) {
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

	h.showUpcomingAppointmentsDatePicker(ctx, chatID, user, newMonth.Format("2006-01"))
}

// HandleUpcomingAppointmentsDateSelection handles date selection from upcoming appointments picker
func (h *ProfessionalHandler) HandleUpcomingAppointmentsDateSelection(ctx context.Context, chatID int64, dateStr string, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetProfessionalAppointmentsByDate(ctx, user.ID, "confirmed", dateStr)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToLoadAppointments, err)
		return
	}

	if len(appointments.Appointments) == 0 {
		h.sendMessage(chatID, common.UIMsgNoUpcomingAppointments)
		h.ShowDashboard(ctx, chatID, user, 0)
		return
	}

	text := common.UIMsgUpcomingAppointments
	for index, apt := range appointments.Appointments {
		text += common.FormatProfessionalAppointmentDetails(&apt, index)
	}

	keyboard := h.createProfessionalAppointmentsKeyboard(appointments.Appointments, false)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}
