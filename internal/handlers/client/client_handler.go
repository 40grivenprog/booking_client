package client

import (
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
func (h *ClientHandler) ShowDashboard(chatID int64) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(common.UIMsgWelcomeBack, user.FirstName, user.Role)
	keyboard := h.createDashboardKeyboard()
	messageIDs := append([]*int{}, user.MessagesToDelete...)
	go func() {
		time.Sleep(1 * time.Second)
		for _, messageID := range messageIDs {
			h.bot.DeleteMessage(chatID, *messageID)
		}
	}()
	id, err := h.sendMessageWithKeyboardAndID(chatID, text, keyboard)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.MessagesToDelete = nil
	user.LastMessageID = &id
	h.apiService.GetUserRepository().SetUser(chatID, user)
}
