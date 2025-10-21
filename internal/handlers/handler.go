package handlers

import (
	"context"
	"time"

	"booking_client/internal/common"
	"booking_client/internal/config"
	"booking_client/internal/handlers/client"
	handlersCommon "booking_client/internal/handlers/common"
	"booking_client/internal/handlers/professional"
	"booking_client/internal/handlers/router"
	"booking_client/internal/middleware"
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
	start := time.Now()

	// Create context with request_id and adjusted logger
	ctx, logger := middleware.RequestIDAndLoggerMiddleware(context.Background(), *h.logger)

	// Handle callback queries (inline keyboard buttons)
	if update.CallbackQuery != nil {
		h.handleCallbackQuery(ctx, update.CallbackQuery)
		latency := time.Since(start)
		logger.Info().
			Dur("latency", latency).
			Msg("Callback query processed")
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

	logger.Info().
		Int64("user_id", userID).
		Str("message", text).
		Msg("Received message from user")

	// Handle different commands and states
	switch text {
	case "/start":
		h.handleStart(ctx, chatID)
	case "/dashboard":
		h.handleDashboard(ctx, chatID)
	default:
		h.handleUserInput(ctx, chatID, text)
	}

	latency := time.Since(start)

	logger.Info().
		Dur("latency", latency).
		Msg("Request completed")
}

// setupRoutes registers all callback handlers with the router
func (h *Handler) setupRoutes() {
	// Initial selection
	h.callbackRouter.RegisterExact(handlersCommon.CallbackClient, func(ctx context.Context, chatID int64, _ string) {
		h.clientHandler.StartRegistration(ctx, chatID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackProfessional, func(ctx context.Context, chatID int64, _ string) {
		h.professionalHandler.StartSignIn(ctx, chatID)
	})

	// Client callbacks
	h.callbackRouter.RegisterExact(handlersCommon.CallbackBookAppointment, func(ctx context.Context, chatID int64, _ string) {
		h.clientHandler.HandleBookAppointment(ctx, chatID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackPendingAppointments, func(ctx context.Context, chatID int64, _ string) {
		h.clientHandler.HandlePendingAppointments(ctx, chatID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackUpcomingAppointments, func(ctx context.Context, chatID int64, _ string) {
		h.clientHandler.HandleUpcomingAppointments(ctx, chatID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackCancelBooking, func(ctx context.Context, chatID int64, _ string) {
		h.clientHandler.HandleCancelBooking(ctx, chatID)
	})

	// Professional callbacks
	h.callbackRouter.RegisterExact(handlersCommon.CallbackProfessionalPendingAppointments, func(ctx context.Context, chatID int64, _ string) {
		h.professionalHandler.HandlePendingAppointments(ctx, chatID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackProfessionalUpcomingAppointments, func(ctx context.Context, chatID int64, _ string) {
		h.professionalHandler.HandleUpcomingAppointments(ctx, chatID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackProfessionalTimetable, func(ctx context.Context, chatID int64, _ string) {
		h.professionalHandler.HandleTimetable(ctx, chatID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackSetUnavailable, func(ctx context.Context, chatID int64, _ string) {
		h.professionalHandler.HandleSetUnavailable(ctx, chatID)
	})

	// Unavailable navigation
	h.callbackRouter.RegisterExact(handlersCommon.CallbackPrevUnavailableMonth, func(ctx context.Context, chatID int64, _ string) {
		h.professionalHandler.HandlePrevUnavailableMonth(ctx, chatID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackNextUnavailableMonth, func(ctx context.Context, chatID int64, _ string) {
		h.professionalHandler.HandleNextUnavailableMonth(ctx, chatID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackCancelUnavailable, func(ctx context.Context, chatID int64, _ string) {
		h.professionalHandler.HandleCancelUnavailable(ctx, chatID)
	})

	// Professional timetable navigation
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixPrevTimetableDay, func(ctx context.Context, chatID int64, date string) {
		h.professionalHandler.HandleTimetableDateNavigation(ctx, chatID, date, handlersCommon.DirectionPrev)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixNextTimetableDay, func(ctx context.Context, chatID int64, date string) {
		h.professionalHandler.HandleTimetableDateNavigation(ctx, chatID, date, handlersCommon.DirectionNext)
	})

	// Professional upcoming appointments navigation
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixPrevUpcomingMonth, func(ctx context.Context, chatID int64, month string) {
		h.professionalHandler.HandleUpcomingAppointmentsMonthNavigation(ctx, chatID, month, handlersCommon.DirectionPrev)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixNextUpcomingMonth, func(ctx context.Context, chatID int64, month string) {
		h.professionalHandler.HandleUpcomingAppointmentsMonthNavigation(ctx, chatID, month, handlersCommon.DirectionNext)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectUpcomingDate, func(ctx context.Context, chatID int64, date string) {
		h.professionalHandler.HandleUpcomingAppointmentsDateSelection(ctx, chatID, date)
	})

	// Client booking flow - month navigation
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixPrevMonth, func(ctx context.Context, chatID int64, month string) {
		h.clientHandler.HandleBookAppointmentsMonthNavigation(ctx, chatID, month, handlersCommon.DirectionPrev)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixNextMonth, func(ctx context.Context, chatID int64, month string) {
		h.clientHandler.HandleBookAppointmentsMonthNavigation(ctx, chatID, month, handlersCommon.DirectionNext)
	})

	// Selection callbacks
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectProfessional, func(ctx context.Context, chatID int64, professionalID string) {
		h.clientHandler.HandleProfessionalSelection(ctx, chatID, professionalID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectDate, func(ctx context.Context, chatID int64, date string) {
		h.clientHandler.HandleDateSelection(ctx, chatID, date)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectTime, func(ctx context.Context, chatID int64, startTime string) {
		h.clientHandler.HandleTimeSelection(ctx, chatID, startTime)
	})

	// Appointment actions
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixCancelAppointment, func(ctx context.Context, chatID int64, appointmentID string) {
		h.clientHandler.HandleCancelAppointment(ctx, chatID, appointmentID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixConfirmAppointment, func(ctx context.Context, chatID int64, appointmentID string) {
		h.professionalHandler.HandleConfirmAppointment(ctx, chatID, appointmentID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixCancelProfAppt, func(ctx context.Context, chatID int64, appointmentID string) {
		h.professionalHandler.HandleCancelAppointment(ctx, chatID, appointmentID)
	})

	// Unavailable flow
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectUnavailableDate, func(ctx context.Context, chatID int64, date string) {
		h.professionalHandler.HandleUnavailableDateSelection(ctx, chatID, date)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectUnavailableStart, func(ctx context.Context, chatID int64, startTime string) {
		h.professionalHandler.HandleUnavailableStartTimeSelection(ctx, chatID, startTime)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectUnavailableEnd, func(ctx context.Context, chatID int64, endTime string) {
		h.professionalHandler.HandleUnavailableEndTimeSelection(ctx, chatID, endTime)
	})

	// Back to dashboard (special case - needs user lookup)
	h.callbackRouter.RegisterExact(handlersCommon.CallbackBackToDashboard, func(ctx context.Context, chatID int64, _ string) {
		user, exists := h.apiService.GetUserRepository().GetUser(chatID)
		if !exists || user == nil {
			text := "❌ User session not found. Please use /start to begin."
			if err := h.bot.SendMessage(chatID, text); err != nil {
				// Use base logger for system errors in callback registration
				h.logger.Error().Err(err).Msg("Failed to send user not found message")
			}
			return
		}
		// Show appropriate dashboard based on user role
		if user.Role == "professional" {
			h.professionalHandler.ShowDashboard(ctx, chatID, user)
		} else {
			h.clientHandler.ShowDashboard(ctx, chatID)
		}
	})
}

// handleCallbackQuery handles inline keyboard button presses
func (h *Handler) handleCallbackQuery(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Use logger from context
	logger := common.GetLogger(ctx)
	logger.Info().
		Int64("user_id", callback.From.ID).
		Str("callback_data", data).
		Msg("Received callback query")

	// Answer the callback query to remove loading state
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := h.bot.GetAPI().Request(callbackConfig); err != nil {
		logger := common.GetLogger(ctx)
		logger.Error().Err(err).Msg("Failed to answer callback query")
	}

	// Route the callback to the appropriate handler
	if !h.callbackRouter.Route(ctx, chatID, data) {
		// No handler found - send unknown command message
		h.sendUnknownCommand(ctx, chatID)
	}
}

// handleStart handles the /start command
func (h *Handler) handleStart(ctx context.Context, chatID int64) {
	// Check if user is already registered
	user, err := h.apiService.GetUserByChatID(ctx, chatID)
	if err == nil && user != nil {
		// User is already registered, show appropriate dashboard
		if user.Role == "professional" {
			h.professionalHandler.ShowDashboard(ctx, chatID, user)
		} else {
			h.clientHandler.ShowDashboard(ctx, chatID)
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
		logger := common.GetLogger(ctx)
		logger.Error().Err(err).Msg("Failed to send start message")
	}
}

// handleStart handles the /dashboard command
func (h *Handler) handleDashboard(ctx context.Context, chatID int64) {
	user, exists := h.apiService.GetUserRepository().GetUser(chatID)
	if !exists || user == nil {
		text := "❌ User session not found. Please use /start to begin."
		if err := h.bot.SendMessage(chatID, text); err != nil {
			// Use base logger for system errors in callback registration
			h.logger.Error().Err(err).Msg("Failed to send user not found message")
		}
		return
	}
	// Show appropriate dashboard based on user role
	if user.Role == "professional" {
		h.professionalHandler.ShowDashboard(ctx, chatID, user)
	} else {
		h.clientHandler.ShowDashboard(ctx, chatID)
	}
}

// handleUserInput handles user input based on their current state
func (h *Handler) handleUserInput(ctx context.Context, chatID int64, text string) {
	// Get user from memory to check state
	user, exists := h.apiService.GetUserRepository().GetUser(chatID)
	if !exists {
		h.sendUnknownCommand(ctx, chatID)
		return
	}

	switch user.State {
	case models.StateWaitingForFirstName:
		h.clientHandler.HandleFirstNameInput(ctx, chatID, text)
	case models.StateWaitingForLastName:
		h.clientHandler.HandleLastNameInput(ctx, chatID, text)
	case models.StateWaitingForPhone:
		h.clientHandler.HandlePhoneInput(ctx, chatID, text)
	case models.StateWaitingForUsername:
		h.professionalHandler.HandleUsernameInput(ctx, chatID, text)
	case models.StateWaitingForPassword:
		h.professionalHandler.HandlePasswordInput(ctx, chatID, text)
	case models.StateWaitingForCancellationReason:
		// Check if user is professional or client to handle cancellation appropriately
		user, exists := h.apiService.GetUserRepository().GetUser(chatID)
		if !exists || user == nil {
			h.sendUnknownCommand(ctx, chatID)
			return
		}

		if user.Role == "professional" {
			h.professionalHandler.HandleCancellationReason(ctx, chatID, text)
		} else {
			h.clientHandler.HandleCancellationReason(ctx, chatID, text)
		}
	case models.StateWaitingForUnavailableDescription:
		h.professionalHandler.HandleUnavailableDescription(ctx, chatID, text)
	default:
		h.sendUnknownCommand(ctx, chatID)
	}
}

// sendUnknownCommand sends unknown command message
func (h *Handler) sendUnknownCommand(ctx context.Context, chatID int64) {
	text := `❓ Unknown command

Please use /start to begin.`

	if err := h.bot.SendMessage(chatID, text); err != nil {
		logger := common.GetLogger(ctx)
		logger.Error().Err(err).Msg("Failed to send unknown command message")
	}
}
