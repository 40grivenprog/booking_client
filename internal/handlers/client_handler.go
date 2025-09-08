package handlers

import (
	"fmt"
	"time"

	"booking_client/internal/models"
	"booking_client/internal/schema"
	"booking_client/internal/services"
	"booking_client/pkg/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// ClientHandler handles all client-related operations
type ClientHandler struct {
	bot        *telegram.Bot
	logger     *zerolog.Logger
	apiService *services.APIService
}

// NewClientHandler creates a new client handler
func NewClientHandler(bot *telegram.Bot, logger *zerolog.Logger, apiService *services.APIService) *ClientHandler {
	return &ClientHandler{
		bot:        bot,
		logger:     logger,
		apiService: apiService,
	}
}

// StartRegistration starts the client registration process
func (h *ClientHandler) StartRegistration(chatID int64) {
	// Create a temporary user with state
	tempUser := &models.User{
		ChatID: &chatID,
		Role:   "client",
		State:  models.StateWaitingForFirstName,
	}

	// Store in memory for state tracking
	h.apiService.GetUserRepository().SetUser(chatID, tempUser)

	text := `👤 Client Registration

Please enter your first name:`

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send client registration message")
	}
}

// ShowDashboard shows the client dashboard with appointment options
func (h *ClientHandler) ShowDashboard(chatID int64, user *models.User) {
	text := fmt.Sprintf(`👋 Welcome back, %s!

You are registered as a %s.

What would you like to do?`, user.FirstName, user.Role)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 Book Appointment", "book_appointment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏳ My Pending Appointments", "pending_appointments"),
			tgbotapi.NewInlineKeyboardButtonData("📋 My Upcoming Appointments", "upcoming_appointments"),
		),
	)

	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send dashboard message")
	}
}

// HandleFirstNameInput handles first name input for client registration
func (h *ClientHandler) HandleFirstNameInput(chatID int64, firstName string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.FirstName = firstName
	user.State = models.StateWaitingForLastName
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := `✅ First name saved!

Please enter your last name:`

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send last name request")
	}
}

// HandleLastNameInput handles last name input for client registration
func (h *ClientHandler) HandleLastNameInput(chatID int64, lastName string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.LastName = lastName
	user.State = models.StateWaitingForPhone
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := `✅ Last name saved!

Please enter your phone number (optional, or type "skip" to skip):`

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send phone request")
	}
}

// HandlePhoneInput handles phone number input for client registration
func (h *ClientHandler) HandlePhoneInput(chatID int64, phone string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	var phoneNumber *string
	if phone != "skip" && phone != "" {
		phoneNumber = &phone
	}

	// Register the client
	req := &schema.RegisterRequest{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		ChatID:      chatID,
		PhoneNumber: phoneNumber,
		Role:        "client",
	}

	registeredUser, err := h.apiService.RegisterClient(req)
	if err != nil {
		text := fmt.Sprintf("❌ Registration failed: %v", err)
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send registration error")
		}
		return
	}

	// Clear state
	user.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, registeredUser)

	text := fmt.Sprintf(`✅ Registration successful!

Welcome, %s %s!
Role: %s
Chat ID: %d`, registeredUser.FirstName, registeredUser.LastName, registeredUser.Role, chatID)

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send registration success")
	}
}

// HandleBookAppointment starts the appointment booking process
func (h *ClientHandler) HandleBookAppointment(chatID int64) {
	// Validate user state - only allow if not in any specific state
	user, valid := h.validateBookingState(chatID, []string{
		models.StateNone,
	})
	if !valid {
		return
	}

	// Set booking state
	user.State = models.StateWaitingForProfessionalSelection
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Get professionals
	professionals, err := h.apiService.GetProfessionals()
	if err != nil {
		text := fmt.Sprintf("❌ Failed to load professionals: %v", err)
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send error message")
		}
		return
	}

	if len(professionals.Professionals) == 0 {
		text := "❌ No professionals available at the moment."
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send no professionals message")
		}
		return
	}

	// Create keyboard with professionals
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, prof := range professionals.Professionals {
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("👨‍💼 %s %s", prof.FirstName, prof.LastName),
			fmt.Sprintf("select_professional_%s", prof.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	// Add cancel booking button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel Booking", "cancel_booking")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	text := "👨‍💼 Please select a professional:"
	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send professionals list")
	}
}

// HandleProfessionalSelection handles when user selects a professional
func (h *ClientHandler) HandleProfessionalSelection(chatID int64, professionalID string) {
	// Validate user state - only allow if not in any specific state or in booking state
	user, valid := h.validateBookingState(chatID, []string{
		models.StateNone,
		models.StateBookingAppointment,
		models.StateWaitingForProfessionalSelection,
	})
	if !valid {
		return
	}

	user.State = models.StateWaitingForDateSelection
	user.SelectedProfessionalID = professionalID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Show current month dates
	h.showDateSelection(chatID, time.Now())
}

// showDateSelection shows available dates for the current month
func (h *ClientHandler) showDateSelection(chatID int64, currentDate time.Time) {
	text := fmt.Sprintf("📅 Select a date (%s %d):", currentDate.Month(), currentDate.Year())

	// Create keyboard with dates for current month
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

		if len(currentRow) == 7 { // 7 days per row
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Add navigation buttons
	prevButton := tgbotapi.NewInlineKeyboardButtonData("⬅️ Previous", "prev_month")
	nextButton := tgbotapi.NewInlineKeyboardButtonData("Next ➡️", "next_month")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(prevButton, nextButton))

	// Add cancel booking button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel Booking", "cancel_booking")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send date selection")
	}
}

// HandleDateSelection handles when user selects a date
func (h *ClientHandler) HandleDateSelection(chatID int64, date string) {
	// Validate user state - only allow if waiting for date selection
	user, valid := h.validateBookingState(chatID, []string{
		models.StateWaitingForDateSelection,
	})
	if !valid {
		return
	}

	user.State = models.StateWaitingForTimeSelection
	user.SelectedDate = date
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Get availability for selected date
	professionalID := user.SelectedProfessionalID
	availability, err := h.apiService.GetProfessionalAvailability(professionalID, date)
	if err != nil {
		text := fmt.Sprintf("❌ Failed to load availability: %v", err)
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send availability error")
		}
		return
	}

	h.showTimeSelection(chatID, availability)
}

// showTimeSelection shows available time slots
func (h *ClientHandler) showTimeSelection(chatID int64, availability *models.ProfessionalAvailabilityResponse) {
	text := fmt.Sprintf("🕐 Select a time slot for %s:", availability.Date)

	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	for _, slot := range availability.Slots {
		if !slot.Available {
			continue
		}

		// Parse the RFC3339 time and format it as HH:MM for display
		startTime, err := time.Parse(time.RFC3339, slot.StartTime)
		if err != nil {
			h.logger.Error().Err(err).Str("time", slot.StartTime).Msg("Failed to parse time slot")
			continue
		}

		timeDisplay := startTime.Format("15:04")
		button := tgbotapi.NewInlineKeyboardButtonData(
			timeDisplay,
			fmt.Sprintf("select_time_%s", timeDisplay),
		)
		currentRow = append(currentRow, button)

		if len(currentRow) == 3 { // 3 slots per row
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Add cancel booking button
	cancelButton := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel Booking", "cancel_booking")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send time selection")
	}
}

// HandleTimeSelection handles when user selects a time slot
func (h *ClientHandler) HandleTimeSelection(chatID int64, startTime string) {
	// Validate user state - only allow if waiting for time selection
	user, valid := h.validateBookingState(chatID, []string{
		models.StateWaitingForTimeSelection,
	})
	if !valid {
		return
	}

	// Parse start time and calculate end time (1 hour later)
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		text := "❌ Invalid time format"
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send time error")
		}
		return
	}

	end := start.Add(time.Hour)
	date := user.SelectedDate

	// Create proper RFC3339 format datetime strings
	// Parse the date and combine with time
	selectedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		text := "❌ Invalid date format"
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send date error")
		}
		return
	}

	// Combine date with time and convert to UTC
	startDateTime := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(),
		start.Hour(), start.Minute(), 0, 0, time.UTC)
	endDateTime := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(),
		end.Hour(), end.Minute(), 0, 0, time.UTC)

	// Validate that start_time is in the future
	if startDateTime.Before(time.Now()) {
		text := "❌ Cannot book appointments in the past. Please select a future time."
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send past time error")
		}
		return
	}

	// Create appointment with RFC3339 format
	req := &schema.CreateAppointmentRequest{
		ClientID:       user.ID,
		ProfessionalID: user.SelectedProfessionalID,
		StartTime:      startDateTime.Format(time.RFC3339),
		EndTime:        endDateTime.Format(time.RFC3339),
	}

	appointment, err := h.apiService.CreateAppointment(req)
	if err != nil {
		text := fmt.Sprintf("❌ Failed to create appointment: %v", err)
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send appointment error")
		}
		return
	}

	// Clear state and show success
	user.State = models.StateNone
	user.SelectedProfessionalID = ""
	user.SelectedDate = ""
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(`✅ Appointment booked successfully!

📅 Date: %s
🕐 Time: %s - %s
👨‍💼 Professional: %s %s

Your appointment is pending confirmation.`,
		date, startTime, end.Format("15:04"),
		appointment.Professional.FirstName, appointment.Professional.LastName)

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send appointment success")
	}

	// Send notification to professional
	h.notifyProfessionalNewAppointment(appointment)
	h.ShowDashboard(chatID, user)
}

// HandlePrevMonth handles previous month navigation
func (h *ClientHandler) HandlePrevMonth(chatID int64) {
	// Validate user state - only allow if waiting for date selection
	_, valid := h.validateBookingState(chatID, []string{
		models.StateWaitingForDateSelection,
	})
	if !valid {
		return
	}

	// For simplicity, we'll just show current month again
	// In a real implementation, you'd store the current month in user state
	h.showDateSelection(chatID, time.Now().AddDate(0, -1, 0))
}

// HandleNextMonth handles next month navigation
func (h *ClientHandler) HandleNextMonth(chatID int64) {
	// Validate user state - only allow if waiting for date selection
	_, valid := h.validateBookingState(chatID, []string{
		models.StateWaitingForDateSelection,
	})
	if !valid {
		return
	}

	// For simplicity, we'll just show current month again
	// In a real implementation, you'd store the current month in user state
	h.showDateSelection(chatID, time.Now().AddDate(0, 1, 0))
}

// HandleCancelBooking cancels the current booking process and returns to dashboard
func (h *ClientHandler) HandleCancelBooking(chatID int64) {
	// Validate user state - only allow if in booking process
	user, valid := h.validateBookingState(chatID, []string{
		models.StateWaitingForProfessionalSelection,
		models.StateWaitingForDateSelection,
		models.StateWaitingForTimeSelection,
		models.StateBookingAppointment,
	})
	if !valid {
		return
	}

	// Clear all booking-related state
	user.State = models.StateNone
	user.SelectedProfessionalID = ""
	user.SelectedDate = ""
	user.SelectedAppointmentID = ""
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := "❌ Booking cancelled. Returning to dashboard."
	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send cancellation message")
	}

	// Show dashboard
	h.ShowDashboard(chatID, user)
}

// HandlePendingAppointments shows pending appointments
func (h *ClientHandler) HandlePendingAppointments(chatID int64) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetClientAppointments(user.ID, "pending")
	if err != nil {
		text := fmt.Sprintf("❌ Failed to load pending appointments: %v", err)
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send error message")
		}
		return
	}

	if len(appointments.Appointments) == 0 {
		text := "📋 You have no pending appointments."
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send no appointments message")
		}
		return
	}

	text := "⏳ Your Pending Appointments:\n\n"
	var rows [][]tgbotapi.InlineKeyboardButton

	for index, apt := range appointments.Appointments {
		text += fmt.Sprintf("✍️ Appointment #%d:\n📅 %s\n🕐 %s - %s\n👨‍💼 %s %s\n\n",
			index+1,
			apt.StartTime[:10], apt.StartTime[11:16], apt.EndTime[11:16],
			apt.Professional.FirstName, apt.Professional.LastName)

		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("❌ Cancel Appointment #%d", index+1),
			fmt.Sprintf("cancel_appointment_%s", apt.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData("🏠 Back to Dashboard", "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send pending appointments")
	}
}

// HandleUpcomingAppointments shows upcoming appointments
func (h *ClientHandler) HandleUpcomingAppointments(chatID int64) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetClientAppointments(user.ID, "confirmed")
	if err != nil {
		text := fmt.Sprintf("❌ Failed to load upcoming appointments: %v", err)
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send error message")
		}
		return
	}

	if len(appointments.Appointments) == 0 {
		text := "📋 You have no upcoming appointments."
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send no appointments message")
		}
		return
	}

	text := "📋 Your Upcoming Appointments:\n\n"
	var rows [][]tgbotapi.InlineKeyboardButton

	for index, apt := range appointments.Appointments {
		text += fmt.Sprintf("✍️ Appointment #%d:\n📅 %s\n🕐 %s - %s\n👨‍💼 %s %s\n\n",
			index+1,
			apt.StartTime[:10], apt.StartTime[11:16], apt.EndTime[11:16],
			apt.Professional.FirstName, apt.Professional.LastName)

		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("❌ Cancel Appointment #%d", index+1),
			fmt.Sprintf("cancel_appointment_%s", apt.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData("🏠 Back to Dashboard", "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send upcoming appointments")
	}
}

// HandleCancelAppointment starts the appointment cancellation process
func (h *ClientHandler) HandleCancelAppointment(chatID int64, appointmentID string) {
	// Store appointment ID and ask for cancellation reason
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.State = models.StateWaitingForCancellationReason
	user.SelectedAppointmentID = appointmentID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := "Please provide a reason for cancelling this appointment:"
	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send cancellation reason request")
	}
}

// HandleCancellationReason handles the cancellation reason input
func (h *ClientHandler) HandleCancellationReason(chatID int64, reason string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	appointmentID := user.SelectedAppointmentID

	// Cancel the appointment
	req := &schema.CancelAppointmentRequest{
		CancellationReason: reason,
	}

	response, err := h.apiService.CancelClientAppointment(user.ID, appointmentID, req)
	if err != nil {
		text := fmt.Sprintf("❌ Failed to cancel appointment: %v", err)
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send cancellation error")
		}
		return
	}

	// Clear state
	user.State = models.StateNone
	user.SelectedAppointmentID = ""
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(`✅ Appointment cancelled successfully!

📅 Date: %s
🕐 Time: %s - %s
👨‍💼 Professional: %s %s
📝 Reason: %s`,
		response.Appointment.StartTime[:10],
		response.Appointment.StartTime[11:16],
		response.Appointment.EndTime[11:16],
		response.Professional.FirstName,
		response.Professional.LastName,
		response.Appointment.CancellationReason)

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send cancellation success")
	}

	// Notify professional about cancellation
	h.notifyProfessionalCancellation(response)

	h.ShowDashboard(chatID, user)
}

// validateBookingState checks if user is in a valid state for booking-related actions
func (h *ClientHandler) validateBookingState(chatID int64, allowedStates []string) (*models.User, bool) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return nil, false
	}

	// Check if user is in an allowed state
	for _, allowedState := range allowedStates {
		if user.State == allowedState {
			return user, true
		}
	}

	// User is not in an allowed state
	text := "❌ This action is not available in your current state. Please use /start to begin a new session."
	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send invalid state message")
	}
	return nil, false
}

// notifyProfessionalNewAppointment sends notification to professional about new appointment
func (h *ClientHandler) notifyProfessionalNewAppointment(appointment *models.CreateAppointmentResponse) {
	if appointment.Professional.ChatID == 0 {
		return // No chat ID for professional
	}

	text := fmt.Sprintf(`🔔 New Appointment Request!

👤 Client: %s %s
📅 Date: %s
🕐 Time: %s - %s

Please confirm or cancel this appointment.`,
		appointment.Client.FirstName, appointment.Client.LastName,
		appointment.Appointment.StartTime[:10],   // Extract date
		appointment.Appointment.StartTime[11:16], // Extract time
		appointment.Appointment.EndTime[11:16])

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Confirm", fmt.Sprintf("confirm_appointment_%s", appointment.Appointment.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", fmt.Sprintf("cancel_appointment_%s", appointment.Appointment.ID)),
		),
	)

	if err := h.bot.SendMessageWithKeyboard(appointment.Professional.ChatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send professional notification")
	}
}

// notifyProfessionalCancellation sends notification to professional about appointment cancellation
func (h *ClientHandler) notifyProfessionalCancellation(response *models.CancelClientAppointmentResponse) {
	if response.Professional.ChatID == nil || *response.Professional.ChatID == 0 {
		return // No chat ID for professional
	}

	text := fmt.Sprintf(`🔔 Appointment Cancelled

👤 Client: %s %s
📅 Date: %s
🕐 Time: %s - %s
📝 Reason: %s`,
		response.Client.FirstName, response.Client.LastName,
		response.Appointment.StartTime[:10],
		response.Appointment.StartTime[11:16],
		response.Appointment.EndTime[11:16],
		response.Appointment.CancellationReason)

	if err := h.bot.SendMessage(*response.Professional.ChatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send professional cancellation notification")
	}
}
