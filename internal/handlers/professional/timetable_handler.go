package professional

import (
	"context"
	"fmt"
	"time"

	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
)

// HandleTimetable shows the professional's timetable for the current date
func (h *ProfessionalHandler) HandleTimetable(ctx context.Context, chatID int64) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	currentDate := time.Now().Format("2006-01-02")
	h.showTimetable(ctx, chatID, user, currentDate)
}

// showTimetable shows the professional's timetable for a specific date
func (h *ProfessionalHandler) showTimetable(ctx context.Context, chatID int64, user *models.User, dateStr string) {
	timetable, err := h.apiService.GetProfessionalTimetable(ctx, user.ID, dateStr)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToLoadAppointments, err)
		return
	}

	date, _ := time.Parse("2006-01-02", dateStr)
	formattedDate := date.Format("Monday, January 2, 2006")

	text := fmt.Sprintf(common.UIMsgTimetableEmpty, formattedDate)
	if len(timetable.Appointments) > 0 {
		text = fmt.Sprintf(common.UIMsgTimetableHeader, formattedDate)
		for i, slot := range timetable.Appointments {
			startTime, _ := time.Parse(time.RFC3339, slot.StartTime)
			endTime, _ := time.Parse(time.RFC3339, slot.EndTime)
			text += fmt.Sprintf(common.UIMsgTimetableSlot, i+1, startTime.Format("15:04"), endTime.Format("15:04"), slot.Description)
		}
	}

	keyboard := h.createTimetableKeyboard(dateStr, timetable.Appointments)
	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleTimetableDateNavigation handles timetable date navigation
func (h *ProfessionalHandler) HandleTimetableDateNavigation(ctx context.Context, chatID int64, dateStr string, direction string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	currentDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgInvalidDateFormat, err)
		return
	}

	var newDate time.Time
	if direction == common.DirectionPrev {
		newDate = currentDate.AddDate(0, 0, -1)
	} else {
		newDate = currentDate.AddDate(0, 0, 1)
	}

	newDateStr := newDate.Format("2006-01-02")
	h.showTimetable(ctx, chatID, user, newDateStr)
}
