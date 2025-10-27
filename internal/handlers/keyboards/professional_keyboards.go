package keyboards

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/schemas"
	"booking_client/internal/util"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// ProfessionalKeyboards handles keyboard creation for professional-related operations
type ProfessionalKeyboards struct {
	logger *zerolog.Logger
}

// NewProfessionalKeyboards creates a new ProfessionalKeyboards instance
func NewProfessionalKeyboards(logger *zerolog.Logger) *ProfessionalKeyboards {
	return &ProfessionalKeyboards{
		logger: logger,
	}
}

// CreateProfessionalDashboardKeyboard creates the professional dashboard keyboard
func (kb *ProfessionalKeyboards) CreateProfessionalDashboardKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnMyTimetable, "professional_timetable"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnPendingAppointments, "professional_pending_appointments"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnUpcomingAppointments, "professional_upcoming_appointments"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnSetUnavailable, "set_unavailable"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnPreviousAppointments, "professional_previous_appointments"),
		),
	)
}

// CreateProfessionalAppointmentsKeyboard creates a keyboard for professional appointment management
func (kb *ProfessionalKeyboards) CreateProfessionalAppointmentsKeyboard(appointments []schemas.ProfessionalAppointment, showConfirm bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for index, apt := range appointments {
		if showConfirm {
			// For pending appointments - show both confirm and cancel
			confirmButton := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(common.BtnConfirmAppointmentProf, index+1),
				fmt.Sprintf("confirm_appointment_%s", apt.ID),
			)
			cancelButton := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(common.BtnCancelAppointmentProf, index+1),
				fmt.Sprintf("cancel_prof_appt_%s", apt.ID),
			)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(confirmButton, cancelButton))
		} else {
			// For upcoming appointments - show only cancel
			cancelButton := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(common.BtnCancelAppointmentProfAlt, index+1),
				fmt.Sprintf("cancel_prof_appt_%s", apt.ID),
			)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))
		}
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnBackToDashboard, "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateUnavailableDateKeyboard creates a keyboard for unavailable date selection
func (kb *ProfessionalKeyboards) CreateUnavailableDateKeyboard(currentDate time.Time) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	// Get first day of month and number of days
	firstDay := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	today := time.Now()

	for d := firstDay; d.Before(lastDay.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		// Skip past dates (compare only date part, not time)
		if d.Year() < today.Year() ||
			(d.Year() == today.Year() && d.Month() < today.Month()) ||
			(d.Year() == today.Year() && d.Month() == today.Month() && d.Day() < today.Day()) {
			continue
		}

		dateStr := d.Format("2006-01-02")
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d", d.Day()),
			fmt.Sprintf("select_unavailable_date_%s", dateStr),
		)
		currentRow = append(currentRow, button)

		if len(currentRow) == common.DaysPerRow {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Add navigation buttons
	currentMonth := currentDate.Format("2006-01")
	todayMonth := today.Format("2006-01")

	var navButtons []tgbotapi.InlineKeyboardButton

	if currentMonth != todayMonth {
		prevButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnPreviousUnavailableMonth, "prev_unavailable_month_"+currentMonth)
		navButtons = append(navButtons, prevButton)
	}

	nextButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnNextUnavailableMonth, "next_unavailable_month_"+currentMonth)
	navButtons = append(navButtons, nextButton)

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(navButtons...))

	// Add cancel button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelUnavailable, "cancel_unavailable")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateUnavailableStartTimeKeyboard creates a keyboard for unavailable start time selection
func (kb *ProfessionalKeyboards) CreateUnavailableStartTimeKeyboard(availability *schemas.ProfessionalAvailabilityResponse) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	for _, slot := range availability.Slots {
		if !slot.Available {
			continue
		}

		// Parse the RFC3339 time and convert to local timezone for display
		startTime, err := time.Parse(time.RFC3339, slot.StartTime)
		if err != nil {
			kb.logger.Error().Err(err).Str("time", slot.StartTime).Msg("Failed to parse time slot")
			continue
		}

		// Convert to local timezone for display
		localTime := startTime.In(util.GetAppTimezone())
		timeDisplay := localTime.Format("15:04")
		button := tgbotapi.NewInlineKeyboardButtonData(
			timeDisplay,
			fmt.Sprintf("select_unavailable_start_%s", timeDisplay),
		)
		currentRow = append(currentRow, button)

		if len(currentRow) == common.TimeSlotsPerRow {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Add cancel button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelUnavailable, "cancel_unavailable")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateUnavailableEndTimeKeyboard creates a keyboard for unavailable end time selection
func (kb *ProfessionalKeyboards) CreateUnavailableEndTimeKeyboard(startTime string, availability *schemas.ProfessionalAvailabilityResponse) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	// Find available slots after the selected start time
	for _, slot := range availability.Slots {
		slotStart, err := time.Parse(time.RFC3339, slot.StartTime)
		if err != nil {
			continue
		}
		slotStartLocal := slotStart.In(util.GetAppTimezone())

		// Only show slots that are:
		// 1. Available (not blocked by existing appointments/unavailable periods)
		// 2. At least 1 hour after the selected start time
		// 3. Before the first unavailable slot
		slotTimeStr := slotStartLocal.Format("15:04")
		if slot.Available && slotTimeStr >= startTime {
			// Use the end time of the slot as the end time option
			slotEnd, err := time.Parse(time.RFC3339, slot.EndTime)
			if err != nil {
				continue
			}
			slotEndLocal := slotEnd.In(util.GetAppTimezone())
			timeDisplay := slotEndLocal.Format("15:04")
			button := tgbotapi.NewInlineKeyboardButtonData(
				timeDisplay,
				fmt.Sprintf("select_unavailable_end_%s", timeDisplay),
			)
			currentRow = append(currentRow, button)

			if len(currentRow) == common.TimeSlotsPerRow {
				rows = append(rows, currentRow)
				currentRow = []tgbotapi.InlineKeyboardButton{}
			}
		} else if !slot.Available && slotTimeStr > startTime {
			// Stop at the first unavailable slot after start time
			break
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Add cancel button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelUnavailable, "cancel_unavailable")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateUpcomingAppointmentsDateKeyboard creates a keyboard for upcoming appointments date selection
func (kb *ProfessionalKeyboards) CreateUpcomingAppointmentsDateKeyboard(dates []string, currentMonth string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	// Add month navigation buttons at the top
	// Only show previous button if not current month
	currentTime := time.Now()
	currentMonthTime, _ := time.Parse("2006-01", currentMonth)
	isCurrentMonth := currentMonthTime.Year() == currentTime.Year() && currentMonthTime.Month() == currentTime.Month()

	var navButtons []tgbotapi.InlineKeyboardButton
	if !isCurrentMonth {
		prevButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnPreviousMonth, "prev_upcoming_month_"+currentMonth)
		navButtons = append(navButtons, prevButton)
	}
	nextButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnNextMonth, "next_upcoming_month_"+currentMonth)
	navButtons = append(navButtons, nextButton)

	if len(navButtons) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(navButtons...))
	}

	for _, dateStr := range dates {
		// Parse date to format it nicely
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		// Format as "Sep 10" or "Today" if it's today
		var displayText string
		if date.Year() == time.Now().Year() && date.Month() == time.Now().Month() && date.Day() == time.Now().Day() {
			displayText = "Today"
		} else {
			displayText = date.Format("Jan 02")
		}

		button := tgbotapi.NewInlineKeyboardButtonData(
			displayText,
			fmt.Sprintf("select_upcoming_date_%s", dateStr),
		)
		currentRow = append(currentRow, button)

		if len(currentRow) == common.DaysPerRow {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnBackToDashboard, "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateTimetableKeyboard creates a keyboard for timetable with day navigation and appointment actions
func (kb *ProfessionalKeyboards) CreateTimetableKeyboard(dateStr string, appointments []schemas.TimetableAppointment) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Add day navigation buttons
	currentDate, _ := time.Parse("2006-01-02", dateStr)
	today := time.Now()
	isToday := currentDate.Year() == today.Year() && currentDate.Month() == today.Month() && currentDate.Day() == today.Day()

	var navButtons []tgbotapi.InlineKeyboardButton
	if !isToday {
		prevButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnPreviousTimetableDay, "prev_timetable_day_"+dateStr)
		navButtons = append(navButtons, prevButton)
	}
	nextButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnNextTimetableDay, "next_timetable_day_"+dateStr)
	navButtons = append(navButtons, nextButton)

	if len(navButtons) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(navButtons...))
	}

	// Add appointment cancel buttons
	for i, apt := range appointments {
		cancelButton := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf(common.BtnCancelTimetableSlot, i+1),
			fmt.Sprintf("cancel_appointment_%s", apt.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnBackToDashboard, "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateClientsKeyboard creates a keyboard for client selection
func CreateClientsKeyboard(clients []schemas.ProfessionalClient) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Add client buttons
	for _, client := range clients {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", client.FirstName, client.LastName),
			fmt.Sprintf("select_client_%s", client.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnBackToDashboard, "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreatePreviousAppointmentsNavigationKeyboard creates navigation keyboard for previous appointments
func CreatePreviousAppointmentsNavigationKeyboard(currentMonth time.Time, hasAppointments bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Navigation buttons
	var navButtons []tgbotapi.InlineKeyboardButton

	// Previous month button (always available)
	prevMonthStr := currentMonth.AddDate(0, -1, 0).Format("2006-01")
	prevButton := tgbotapi.NewInlineKeyboardButtonData(
		common.BtnPreviousMonth,
		fmt.Sprintf("prev_previous_month_%s", prevMonthStr),
	)
	navButtons = append(navButtons, prevButton)

	// Next month button (only if not in future)
	nextMonth := currentMonth.AddDate(0, 1, 0)
	if nextMonth.Before(time.Now()) || nextMonth.Format("2006-01") == time.Now().Format("2006-01") {
		nextMonthStr := nextMonth.Format("2006-01")
		nextButton := tgbotapi.NewInlineKeyboardButtonData(
			common.BtnNextMonth,
			fmt.Sprintf("next_previous_month_%s", nextMonthStr),
		)
		navButtons = append(navButtons, nextButton)
	}

	if len(navButtons) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(navButtons...))
	}

	// Back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnBackToDashboard, "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
