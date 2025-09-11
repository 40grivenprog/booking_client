package handlers

import (
	"booking_client/internal/config"
	"booking_client/internal/handlers/client"
	"booking_client/internal/handlers/common"
	"booking_client/internal/handlers/professional"
	"booking_client/internal/models"
	"booking_client/internal/services"
	"booking_client/pkg/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// Handler manages all bot command handlers
type Handler struct {
	bot                 *telegram.Bot
	config              *config.Config
	logger              *zerolog.Logger
	apiService          *services.APIService
	clientHandler       *client.ClientHandler
	professionalHandler *professional.ProfessionalHandler
}

// NewHandler creates a new handler instance
func NewHandler(bot *telegram.Bot, config *config.Config, logger *zerolog.Logger) *Handler {
	apiService := services.NewAPIService(config, logger)

	return &Handler{
		bot:                 bot,
		config:              config,
		logger:              logger,
		apiService:          apiService,
		clientHandler:       client.NewClientHandler(bot, logger, apiService),
		professionalHandler: professional.NewProfessionalHandler(bot, logger, apiService),
	}
}

// RegisterHandlers registers all command handlers
func (h *Handler) RegisterHandlers() {
	// Set this handler as the update handler for the bot
	h.bot.SetUpdateHandler(h)
}

// HandleUpdate processes incoming updates (implements UpdateHandler interface)
func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	// Handle callback queries (inline keyboard buttons)
	if update.CallbackQuery != nil {
		h.handleCallbackQuery(update.CallbackQuery)
		return
	}

	// Handle regular messages
	if update.Message == nil {
		return
	}

	message := update.Message
	chatID := message.Chat.ID
	userID := message.From.ID
	text := message.Text

	h.logger.Info().
		Int64("user_id", userID).
		Str("message", text).
		Msg("Received message from user")

	// Handle different commands and states
	switch {
	case text == "/start":
		h.handleStart(chatID, int(userID))
	default:
		h.handleUserInput(chatID, text)
	}
}

// handleCallbackQuery handles inline keyboard button presses
func (h *Handler) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	h.logger.Info().
		Int64("user_id", callback.From.ID).
		Str("callback_data", data).
		Msg("Received callback query")

	// Answer the callback query to remove loading state
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := h.bot.GetAPI().Request(callbackConfig); err != nil {
		h.logger.Error().Err(err).Msg("Failed to answer callback query")
	}

	// Handle the callback data
	switch {
	case data == "client":
		h.clientHandler.StartRegistration(chatID)
	case data == "professional":
		h.professionalHandler.StartSignIn(chatID)
	case data == "book_appointment":
		h.clientHandler.HandleBookAppointment(chatID)
	case data == "pending_appointments":
		h.clientHandler.HandlePendingAppointments(chatID)
	case data == "upcoming_appointments":
		h.clientHandler.HandleUpcomingAppointments(chatID)
	case data == "professional_pending_appointments":
		h.professionalHandler.HandlePendingAppointments(chatID)
	case data == "professional_upcoming_appointments":
		h.professionalHandler.HandleUpcomingAppointments(chatID)
	case data == "professional_timetable":
		h.professionalHandler.HandleTimetable(chatID)
	case len(data) >= 19 && data[:19] == "prev_timetable_day_":
		date := data[19:]
		h.professionalHandler.HandleTimetableDateNavigation(chatID, date, common.DirectionPrev)
	case len(data) >= 19 && data[:19] == "next_timetable_day_":
		date := data[19:]
		h.professionalHandler.HandleTimetableDateNavigation(chatID, date, common.DirectionNext)
	case len(data) >= 20 && data[:20] == "prev_upcoming_month_":
		month := data[20:]
		h.professionalHandler.HandleUpcomingAppointmentsMonthNavigation(chatID, month, common.DirectionPrev)
	case len(data) >= 20 && data[:20] == "next_upcoming_month_":
		month := data[20:]
		h.professionalHandler.HandleUpcomingAppointmentsMonthNavigation(chatID, month, common.DirectionNext)
	case len(data) >= 21 && data[:21] == "select_upcoming_date_":
		date := data[21:]
		h.professionalHandler.HandleUpcomingAppointmentsDateSelection(chatID, date)
	case data == "set_unavailable":
		h.professionalHandler.HandleSetUnavailable(chatID)
	case data == "cancel_booking":
		h.clientHandler.HandleCancelBooking(chatID)
	case len(data) > 20 && data[:20] == "select_professional_":
		professionalID := data[20:]
		h.clientHandler.HandleProfessionalSelection(chatID, professionalID)
	case len(data) > 12 && data[:12] == "select_date_":
		date := data[12:]
		h.clientHandler.HandleDateSelection(chatID, date)
	case len(data) > 12 && data[:12] == "select_time_":
		startTime := data[12:]
		h.clientHandler.HandleTimeSelection(chatID, startTime)
	case data == "prev_month":
		h.clientHandler.HandlePrevMonth(chatID)
	case data == "next_month":
		h.clientHandler.HandleNextMonth(chatID)
	case len(data) > 19 && data[:19] == "cancel_appointment_":
		appointmentID := data[19:]
		h.clientHandler.HandleCancelAppointment(chatID, appointmentID)
	case len(data) > 20 && data[:20] == "confirm_appointment_":
		appointmentID := data[20:]
		h.professionalHandler.HandleConfirmAppointment(chatID, appointmentID)
	case len(data) > 20 && data[:17] == "cancel_prof_appt_":
		appointmentID := data[17:]
		h.professionalHandler.HandleCancelAppointment(chatID, appointmentID)
	case len(data) >= 24 && data[:24] == "select_unavailable_date_":
		date := data[24:]
		h.professionalHandler.HandleUnavailableDateSelection(chatID, date)
	case len(data) >= 25 && data[:25] == "select_unavailable_start_":
		startTime := data[25:]
		h.professionalHandler.HandleUnavailableStartTimeSelection(chatID, startTime)
	case len(data) >= 23 && data[:23] == "select_unavailable_end_":
		endTime := data[23:]
		h.professionalHandler.HandleUnavailableEndTimeSelection(chatID, endTime)
	case data == "prev_unavailable_month":
		h.professionalHandler.HandlePrevUnavailableMonth(chatID)
	case data == "next_unavailable_month":
		h.professionalHandler.HandleNextUnavailableMonth(chatID)
	case data == "cancel_unavailable":
		h.professionalHandler.HandleCancelUnavailable(chatID)
	case data == "back_to_dashboard":
		user, exists := h.apiService.GetUserRepository().GetUser(chatID)
		if !exists || user == nil {
			text := "❌ User session not found. Please use /start to begin."
			if err := h.bot.SendMessage(chatID, text); err != nil {
				h.logger.Error().Err(err).Msg("Failed to send user not found message")
			}
			return
		}
		// Show appropriate dashboard based on user role
		if user.Role == "professional" {
			h.professionalHandler.ShowDashboard(chatID, user)
		} else {
			h.clientHandler.ShowDashboard(chatID)
		}
	default:
		h.sendUnknownCommand(chatID)
	}
}

// handleStart handles the /start command
func (h *Handler) handleStart(chatID int64, userID int) {
	// Check if user is already registered
	user, err := h.apiService.GetUserByChatID(chatID)
	if err == nil && user != nil {
		// User is already registered, show appropriate dashboard
		if user.Role == "professional" {
			h.professionalHandler.ShowDashboard(chatID, user)
		} else {
			h.clientHandler.ShowDashboard(chatID)
		}
		return
	}

	// User is not registered, ask for role selection
	welcomeText := `👋 Welcome to the Booking Bot!

Please choose how you want to continue:`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👤 Client", "client"),
			tgbotapi.NewInlineKeyboardButtonData("👨‍💼 Professional", "professional"),
		),
	)

	if err := h.bot.SendMessageWithKeyboard(chatID, welcomeText, keyboard); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send start message")
	}
}

// handleUserInput handles user input based on their current state
func (h *Handler) handleUserInput(chatID int64, text string) {
	// Get user from memory to check state
	user, exists := h.apiService.GetUserRepository().GetUser(chatID)
	if !exists {
		h.sendUnknownCommand(chatID)
		return
	}

	switch user.State {
	case models.StateWaitingForFirstName:
		h.clientHandler.HandleFirstNameInput(chatID, text)
	case models.StateWaitingForLastName:
		h.clientHandler.HandleLastNameInput(chatID, text)
	case models.StateWaitingForPhone:
		h.clientHandler.HandlePhoneInput(chatID, text)
	case models.StateWaitingForUsername:
		h.professionalHandler.HandleUsernameInput(chatID, text)
	case models.StateWaitingForPassword:
		h.professionalHandler.HandlePasswordInput(chatID, text)
	case models.StateWaitingForCancellationReason:
		// Check if user is professional or client to handle cancellation appropriately
		user, exists := h.apiService.GetUserRepository().GetUser(chatID)
		if !exists || user == nil {
			h.sendUnknownCommand(chatID)
			return
		}

		if user.Role == "professional" {
			h.professionalHandler.HandleCancellationReason(chatID, text)
		} else {
			h.clientHandler.HandleCancellationReason(chatID, text)
		}
	case models.StateWaitingForUnavailableDescription:
		h.professionalHandler.HandleUnavailableDescription(chatID, text)
	default:
		h.sendUnknownCommand(chatID)
	}
}

// sendUnknownCommand sends unknown command message
func (h *Handler) sendUnknownCommand(chatID int64) {
	text := `❓ Unknown command

Please use /start to begin.`

	if err := h.bot.SendMessage(chatID, text); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send unknown command message")
	}
}
