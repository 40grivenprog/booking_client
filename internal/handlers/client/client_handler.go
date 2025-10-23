package client

import (
	"context"
	"fmt"
	"time"

	"booking_client/internal/handlers/common"
	"booking_client/internal/handlers/keyboards"
	"booking_client/internal/models"
	apiService "booking_client/internal/services/api_service"
	"booking_client/pkg/telegram"

	"github.com/rs/zerolog"
)

// ClientHandler handles all client-related operations
type ClientHandler struct {
	bot                 *telegram.Bot
	logger              *zerolog.Logger
	apiService          *apiService.APIService
	notificationService *common.NotificationService
	keyboards           *keyboards.ClientKeyboards
}

// NewClientHandler creates a new client handler
func NewClientHandler(bot *telegram.Bot, logger *zerolog.Logger, apiService *apiService.APIService) *ClientHandler {
	return &ClientHandler{
		bot:                 bot,
		logger:              logger,
		apiService:          apiService,
		notificationService: common.NewNotificationService(bot, logger, apiService),
		keyboards:           keyboards.NewClientKeyboards(logger),
	}
}

// ShowDashboard shows the client dashboard with appointment options
func (h *ClientHandler) ShowDashboard(ctx context.Context, chatID int64, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.LastMessageID = &messageID
	user.MessagesToDelete = append(user.MessagesToDelete, &messageID)
	user.State = models.StateNone
	messageIDs := append([]*int{}, user.MessagesToDelete...)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(common.UIMsgWelcomeBack, user.FirstName, user.Role)
	keyboard := h.createDashboardKeyboard()
	go func() {
		time.Sleep(3 * time.Second)
		for _, messageID := range messageIDs {
			if messageID != nil {
				h.bot.DeleteMessage(chatID, *messageID)
			}
		}
	}()
	id, err := h.sendMessageWithKeyboardAndID(chatID, text, keyboard)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	user.LastMessageID = &id
	h.apiService.GetUserRepository().SetUser(chatID, user)
}
