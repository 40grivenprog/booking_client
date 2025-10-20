package keyboards

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	"booking_client/internal/util"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// ClientKeyboards handles keyboard creation for client-related operations
type ClientKeyboards struct {
	logger *zerolog.Logger
}

// NewClientKeyboards creates a new ClientKeyboards instance
func NewClientKeyboards(logger *zerolog.Logger) *ClientKeyboards {
	return &ClientKeyboards{
		logger: logger,
	}
}

// CreateDateKeyboard creates a keyboard for date selection
func (kb *ClientKeyboards) CreateDateKeyboard(currentDate time.Time) tgbotapi.InlineKeyboardMarkup {
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
			fmt.Sprintf("select_date_%s", dateStr),
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
	currentMonth := currentDate.Format("2006-01")
	todayMonth := today.Format("2006-01")

	var navButtons []tgbotapi.InlineKeyboardButton

	if currentMonth != todayMonth {
		prevButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnPreviousMonth, "prev_month_"+currentMonth)
		navButtons = append(navButtons, prevButton)
	}

	nextButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnNextMonth, "next_month_"+currentMonth)
	navButtons = append(navButtons, nextButton)

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(navButtons...))

	// Add cancel booking button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelBooking, "cancel_booking")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateTimeKeyboard creates a keyboard for time slot selection
func (kb *ClientKeyboards) CreateTimeKeyboard(availability *models.ProfessionalAvailabilityResponse) tgbotapi.InlineKeyboardMarkup {
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
			fmt.Sprintf("select_time_%s", timeDisplay),
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

	// Add cancel booking button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelBooking, "cancel_booking")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateProfessionalsKeyboard creates a keyboard for professional selection
func (kb *ClientKeyboards) CreateProfessionalsKeyboard(professionals []models.User) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, prof := range professionals {
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("👨‍💼 %s %s", prof.FirstName, prof.LastName),
			fmt.Sprintf("select_professional_%s", prof.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	// Add cancel booking button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnCancelBooking, "cancel_booking")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateAppointmentsKeyboard creates a keyboard for appointment management
func (kb *ClientKeyboards) CreateAppointmentsKeyboard(appointments []models.ClientAppointment, buttonPrefix string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for index, apt := range appointments {
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf(buttonPrefix, index+1),
			fmt.Sprintf("cancel_appointment_%s", apt.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData(common.BtnBackToDashboard, "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateDashboardKeyboard creates the main dashboard keyboard
func (kb *ClientKeyboards) CreateDashboardKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnBookAppointment, "book_appointment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnMyPendingAppointments, "pending_appointments"),
			tgbotapi.NewInlineKeyboardButtonData(common.BtnMyUpcomingAppointments, "upcoming_appointments"),
		),
	)
}

// CreateRegistrationSuccessKeyboard creates keyboard for successful registration
func (kb *ClientKeyboards) CreateRegistrationSuccessKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(common.BtnGoToDashboard, "back_to_dashboard"),
		),
	)
}
