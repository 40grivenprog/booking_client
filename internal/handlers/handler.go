package handlers

import (
	"booking_client/internal/config"
	"booking_client/internal/handlers/client"
	"booking_client/internal/handlers/common"
	"booking_client/internal/handlers/professional"
	"booking_client/internal/handlers/router"
	"booking_client/internal/models"
	apiService "booking_client/internal/services/api_service"
	"booking_client/pkg/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// Handler manages all bot command handlers
type Handler struct {
	bot                 *telegram.Bot
	config              *config.Config
	logger              *zerolog.Logger
	apiService          *apiService.APIService
	clientHandler       *client.ClientHandler
	professionalHandler *professional.ProfessionalHandler
	callbackRouter      *router.CallbackRouter
}

// NewHandler creates a new handler instance
func NewHandler(bot *telegram.Bot, config *config.Config, logger *zerolog.Logger) (*Handler, error) {
	apiService, err := apiService.NewAPIService(config, logger)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		bot:                 bot,
		config:              config,
		logger:              logger,
		apiService:          apiService,
		clientHandler:       client.NewClientHandler(bot, logger, apiService),
		professionalHandler: professional.NewProfessionalHandler(bot, logger, apiService),
		callbackRouter:      router.NewCallbackRouter(logger),
	}

	// Setup callback routes
	h.setupRoutes()

	exactCount, prefixCount := h.callbackRouter.GetStats()
	logger.Info().
		Int("exact_handlers", exactCount).
		Int("prefix_handlers", prefixCount).
		Msg("Callback router initialized")

	return h, nil
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

// setupRoutes registers all callback handlers with the router
func (h *Handler) setupRoutes() {
	// Initial selection
	h.callbackRouter.RegisterExact(common.CallbackClient, func(chatID int64, _ string) {
		h.clientHandler.StartRegistration(chatID)
	})
	h.callbackRouter.RegisterExact(common.CallbackProfessional, func(chatID int64, _ string) {
		h.professionalHandler.StartSignIn(chatID)
	})

	// Client callbacks
	h.callbackRouter.RegisterExact(common.CallbackBookAppointment, func(chatID int64, _ string) {
		h.clientHandler.HandleBookAppointment(chatID)
	})
	h.callbackRouter.RegisterExact(common.CallbackPendingAppointments, func(chatID int64, _ string) {
		h.clientHandler.HandlePendingAppointments(chatID)
	})
	h.callbackRouter.RegisterExact(common.CallbackUpcomingAppointments, func(chatID int64, _ string) {
		h.clientHandler.HandleUpcomingAppointments(chatID)
	})
	h.callbackRouter.RegisterExact(common.CallbackCancelBooking, func(chatID int64, _ string) {
		h.clientHandler.HandleCancelBooking(chatID)
	})

	// Professional callbacks
	h.callbackRouter.RegisterExact(common.CallbackProfessionalPendingAppointments, func(chatID int64, _ string) {
		h.professionalHandler.HandlePendingAppointments(chatID)
	})
	h.callbackRouter.RegisterExact(common.CallbackProfessionalUpcomingAppointments, func(chatID int64, _ string) {
		h.professionalHandler.HandleUpcomingAppointments(chatID)
	})
	h.callbackRouter.RegisterExact(common.CallbackProfessionalTimetable, func(chatID int64, _ string) {
		h.professionalHandler.HandleTimetable(chatID)
	})
	h.callbackRouter.RegisterExact(common.CallbackSetUnavailable, func(chatID int64, _ string) {
		h.professionalHandler.HandleSetUnavailable(chatID)
	})

	// Unavailable navigation
	h.callbackRouter.RegisterExact(common.CallbackPrevUnavailableMonth, func(chatID int64, _ string) {
		h.professionalHandler.HandlePrevUnavailableMonth(chatID)
	})
	h.callbackRouter.RegisterExact(common.CallbackNextUnavailableMonth, func(chatID int64, _ string) {
		h.professionalHandler.HandleNextUnavailableMonth(chatID)
	})
	h.callbackRouter.RegisterExact(common.CallbackCancelUnavailable, func(chatID int64, _ string) {
		h.professionalHandler.HandleCancelUnavailable(chatID)
	})

	// Professional timetable navigation
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixPrevTimetableDay, func(chatID int64, date string) {
		h.professionalHandler.HandleTimetableDateNavigation(chatID, date, common.DirectionPrev)
	})
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixNextTimetableDay, func(chatID int64, date string) {
		h.professionalHandler.HandleTimetableDateNavigation(chatID, date, common.DirectionNext)
	})

	// Professional upcoming appointments navigation
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixPrevUpcomingMonth, func(chatID int64, month string) {
		h.professionalHandler.HandleUpcomingAppointmentsMonthNavigation(chatID, month, common.DirectionPrev)
	})
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixNextUpcomingMonth, func(chatID int64, month string) {
		h.professionalHandler.HandleUpcomingAppointmentsMonthNavigation(chatID, month, common.DirectionNext)
	})
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixSelectUpcomingDate, func(chatID int64, date string) {
		h.professionalHandler.HandleUpcomingAppointmentsDateSelection(chatID, date)
	})

	// Client booking flow - month navigation
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixPrevMonth, func(chatID int64, month string) {
		h.clientHandler.HandleBookAppointmentsMonthNavigation(chatID, month, common.DirectionPrev)
	})
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixNextMonth, func(chatID int64, month string) {
		h.clientHandler.HandleBookAppointmentsMonthNavigation(chatID, month, common.DirectionNext)
	})

	// Selection callbacks
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixSelectProfessional, func(chatID int64, professionalID string) {
		h.clientHandler.HandleProfessionalSelection(chatID, professionalID)
	})
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixSelectDate, func(chatID int64, date string) {
		h.clientHandler.HandleDateSelection(chatID, date)
	})
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixSelectTime, func(chatID int64, startTime string) {
		h.clientHandler.HandleTimeSelection(chatID, startTime)
	})

	// Appointment actions
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixCancelAppointment, func(chatID int64, appointmentID string) {
		h.clientHandler.HandleCancelAppointment(chatID, appointmentID)
	})
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixConfirmAppointment, func(chatID int64, appointmentID string) {
		h.professionalHandler.HandleConfirmAppointment(chatID, appointmentID)
	})
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixCancelProfAppt, func(chatID int64, appointmentID string) {
		h.professionalHandler.HandleCancelAppointment(chatID, appointmentID)
	})

	// Unavailable flow
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixSelectUnavailableDate, func(chatID int64, date string) {
		h.professionalHandler.HandleUnavailableDateSelection(chatID, date)
	})
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixSelectUnavailableStart, func(chatID int64, startTime string) {
		h.professionalHandler.HandleUnavailableStartTimeSelection(chatID, startTime)
	})
	h.callbackRouter.RegisterPrefix(common.CallbackPrefixSelectUnavailableEnd, func(chatID int64, endTime string) {
		h.professionalHandler.HandleUnavailableEndTimeSelection(chatID, endTime)
	})

	// Back to dashboard (special case - needs user lookup)
	h.callbackRouter.RegisterExact(common.CallbackBackToDashboard, func(chatID int64, _ string) {
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
	})
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

	// Route the callback to the appropriate handler
	if !h.callbackRouter.Route(chatID, data) {
		// No handler found - send unknown command message
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
