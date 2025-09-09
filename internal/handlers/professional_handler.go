package handlers

import (
	"fmt"

	"booking_client/internal/models"
	"booking_client/internal/schema"
	"booking_client/internal/services"
	"booking_client/pkg/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// ProfessionalHandler handles all professional-related operations
type ProfessionalHandler struct {
	bot        *telegram.Bot
	logger     *zerolog.Logger
	apiService *services.APIService
}

// NewProfessionalHandler creates a new professional handler
func NewProfessionalHandler(bot *telegram.Bot, logger *zerolog.Logger, apiService *services.APIService) *ProfessionalHandler {
	return &ProfessionalHandler{
		bot:        bot,
		logger:     logger,
		apiService: apiService,
	}
}

// StartSignIn starts the professional sign-in process
func (h *ProfessionalHandler) StartSignIn(chatID int64) {
	// Create a temporary user with state
	tempUser := &models.User{
		ChatID: &chatID,
		Role:   "professional",
		State:  models.StateWaitingForUsername,
	}

	// Store in memory for state tracking
	h.apiService.GetUserRepository().SetUser(chatID, tempUser)

	text := `👨‍💼 Professional Sign In

Please enter your username:`

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send professional sign-in message")
	}
}

// ShowDashboard shows the professional dashboard with appointment options
func (h *ProfessionalHandler) ShowDashboard(chatID int64, user *models.User) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, user)
	text := fmt.Sprintf(`👋 Welcome back, %s!

You are registered as a %s.

What would you like to do?`, user.LastName, user.Role)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏳ Pending Appointments", "professional_pending_appointments"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Upcoming Appointments", "professional_upcoming_appointments"),
		),
	)

	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send professional dashboard message")
	}
}

// HandleUsernameInput handles username input for professional sign-in
func (h *ProfessionalHandler) HandleUsernameInput(chatID int64, username string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.Username = username
	user.State = models.StateWaitingForPassword
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := `✅ Username saved!

Please enter your password:`

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send password request")
	}
}

// HandlePasswordInput handles password input for professional sign-in
func (h *ProfessionalHandler) HandlePasswordInput(chatID int64, password string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Sign in the professional
	req := &schema.ProfessionalSignInRequest{
		Username: user.Username,
		Password: password,
		ChatID:   chatID,
	}

	signedInUser, err := h.apiService.SignInProfessional(req)
	if err != nil {
		text := fmt.Sprintf("❌ Sign in failed: %v", err)
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send sign-in error")
		}
		return
	}

	// Clear state
	user.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, signedInUser)

	text := fmt.Sprintf(`✅ Sign in successful!

Welcome back, %s %s!
Role: %s
Username: %s
Chat ID: %d`, signedInUser.FirstName, signedInUser.LastName, signedInUser.Role, signedInUser.Username, chatID)

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send sign-in success")
	}
	h.ShowDashboard(chatID, signedInUser)
}

// HandlePendingAppointments shows pending appointments for professionals
func (h *ProfessionalHandler) HandlePendingAppointments(chatID int64) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetProfessionalAppointments(user.ID, "pending")
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
		h.ShowDashboard(chatID, user)
		return
	}

	text := "⏳ Your Pending Appointments:\n\n"
	var rows [][]tgbotapi.InlineKeyboardButton

	for index, apt := range appointments.Appointments {
		text += fmt.Sprintf("✍️ Appointment #%d:\n📅 %s\n🕐 %s - %s\n👤 Client: %s %s\n📝 %s\n\n",
			index+1,
			apt.StartTime[:10], apt.StartTime[11:16], apt.EndTime[11:16],
			apt.Client.FirstName, apt.Client.LastName,
			apt.Description)

		// Add confirm and cancel buttons
		confirmButton := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("✅ Confirm Appointment #%d", index+1),
			fmt.Sprintf("confirm_appointment_%s", apt.ID),
		)
		cancelButton := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("❌ Cancel Appointment #%d", index+1),
			fmt.Sprintf("cancel_prof_appt_%s", apt.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(confirmButton, cancelButton))
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData("🏠 Back to Dashboard", "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send pending appointments")
	}
}

// HandleUpcomingAppointments shows upcoming appointments for professionals
func (h *ProfessionalHandler) HandleUpcomingAppointments(chatID int64) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	appointments, err := h.apiService.GetProfessionalAppointments(user.ID, "confirmed")
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
		h.ShowDashboard(chatID, user)
		return
	}

	text := "📋 Your Upcoming Appointments:\n\n"
	var rows [][]tgbotapi.InlineKeyboardButton

	for index, apt := range appointments.Appointments {
		text += fmt.Sprintf("✍️ Appointment #%d:\n📅 %s\n🕐 %s - %s\n👤 Client: %s %s\n📝 %s\n\n",
			index+1,
			apt.StartTime[:10], apt.StartTime[11:16], apt.EndTime[11:16],
			apt.Client.FirstName, apt.Client.LastName,
			apt.Description)

		// Add cancel button for upcoming appointments
		cancelButton := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("❌ Cancel Appintment %d", index+1),
			fmt.Sprintf("cancel_prof_appt_%s", apt.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelButton))
	}

	// Add back to dashboard button
	backButton := tgbotapi.NewInlineKeyboardButtonData("🏠 Back to Dashboard", "back_to_dashboard")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	if err := h.bot.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send upcoming appointments")
	}
}

// HandleConfirmAppointment confirms an appointment by professional
func (h *ProfessionalHandler) HandleConfirmAppointment(chatID int64, appointmentID string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Confirm the appointment
	req := &schema.ConfirmAppointmentRequest{}

	response, err := h.apiService.ConfirmProfessionalAppointment(user.ID, appointmentID, req)
	if err != nil {
		text := fmt.Sprintf("❌ Failed to confirm appointment: %v", err)
		if err := h.bot.SendMessage(chatID, text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send confirmation error")
		}
		return
	}
	fmt.Println("response", response.Appointment.StartTime)

	text := fmt.Sprintf(`✅ Appointment confirmed successfully!

📅 Date: %s
🕐 Time: %s - %s
👤 Client: %s %s`,
		response.Appointment.StartTime[:10],
		response.Appointment.StartTime[11:16],
		response.Appointment.EndTime[11:16],
		response.Client.FirstName,
		response.Client.LastName)

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send confirmation success")
	}

	// Notify client about confirmation
	h.notifyClientAppointmentConfirmation(response)

	// Show dashboard
	h.ShowDashboard(chatID, user)
}

// HandleCancelAppointment starts the professional appointment cancellation process
func (h *ProfessionalHandler) HandleCancelAppointment(chatID int64, appointmentID string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Store appointment ID and ask for cancellation reason
	user.State = models.StateWaitingForCancellationReason
	user.SelectedAppointmentID = appointmentID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := "Please provide a reason for cancelling this appointment:"
	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send cancellation reason request")
	}
}

// HandleCancellationReason handles the professional cancellation reason input
func (h *ProfessionalHandler) HandleCancellationReason(chatID int64, reason string) {
	user, ok := getUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	appointmentID := user.SelectedAppointmentID

	// Cancel the appointment
	req := &schema.CancelAppointmentRequest{
		CancellationReason: reason,
	}

	response, err := h.apiService.CancelProfessionalAppointment(user.ID, appointmentID, req)
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
👤 Client: %s %s
📝 Reason: %s`,
		response.Appointment.StartTime[:10],
		response.Appointment.StartTime[11:16],
		response.Appointment.EndTime[11:16],
		response.Client.FirstName,
		response.Client.LastName,
		response.Appointment.CancellationReason)

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send cancellation success")
	}

	// Notify client about cancellation
	h.notifyClientProfessionalCancellation(response)

	// Show dashboard
	h.ShowDashboard(chatID, user)
}

// notifyClientAppointmentConfirmation sends notification to client about appointment confirmation
func (h *ProfessionalHandler) notifyClientAppointmentConfirmation(response *models.ConfirmProfessionalAppointmentResponse) {
	if response.Client.ChatID == nil || *response.Client.ChatID == 0 {
		return // No chat ID for client
	}

	text := fmt.Sprintf(`✅ Appointment Confirmed!

📅 Date: %s
🕐 Time: %s - %s
👨‍💼 Professional: %s %s

Your appointment has been confirmed.`,
		response.Appointment.StartTime[:10],
		response.Appointment.StartTime[11:16],
		response.Appointment.EndTime[11:16],
		response.Professional.FirstName,
		response.Professional.LastName)

	if err := h.bot.SendMessage(*response.Client.ChatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send client confirmation notification")
	}
}

// notifyClientProfessionalCancellation sends notification to client about appointment cancellation by professional
func (h *ProfessionalHandler) notifyClientProfessionalCancellation(response *models.CancelProfessionalAppointmentResponse) {
	fmt.Println("response", response.Client.ChatID)
	if response.Client.ChatID == nil || *response.Client.ChatID == 0 {
		return // No chat ID for client
	}

	text := fmt.Sprintf(`🔔 Appointment Cancelled by Professional

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

	if err := h.bot.SendMessage(*response.Client.ChatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send client cancellation notification")
	}
}
