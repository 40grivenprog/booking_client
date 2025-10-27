package professional

import (
	"context"
	"time"

	"booking_client/internal/handlers/common"
	"booking_client/internal/handlers/keyboards"
	"booking_client/internal/models"
	apiService "booking_client/internal/services/api_service"
	"booking_client/pkg/telegram"
	"fmt"

	"github.com/rs/zerolog"
)

// ProfessionalHandler handles all professional-related operations
type ProfessionalHandler struct {
	bot                 *telegram.Bot
	logger              *zerolog.Logger
	apiService          *apiService.APIService
	notificationService *common.NotificationService
	keyboards           *keyboards.ProfessionalKeyboards
}

// NewProfessionalHandler creates a new professional handler
func NewProfessionalHandler(bot *telegram.Bot, logger *zerolog.Logger, apiService *apiService.APIService) *ProfessionalHandler {
	return &ProfessionalHandler{
		bot:                 bot,
		logger:              logger,
		apiService:          apiService,
		notificationService: common.NewNotificationService(bot, logger, apiService),
		keyboards:           keyboards.NewProfessionalKeyboards(logger),
	}
}

// ShowDashboard shows the professional dashboard with appointment options
func (h *ProfessionalHandler) ShowDashboard(ctx context.Context, chatID int64, user *models.User, messageID int) {
	currentUser, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	currentUser.LastMessageID = &messageID
	currentUser.MessagesToDelete = append(currentUser.MessagesToDelete, &messageID)
	currentUser.State = models.StateNone
	messageIDs := append([]*int{}, user.MessagesToDelete...)
	h.apiService.GetUserRepository().SetUser(chatID, currentUser)

	go func() {
		time.Sleep(3 * time.Second)
		for _, messageID := range messageIDs {
			if messageID != nil {
				h.bot.DeleteMessage(chatID, *messageID)
			}
		}
	}()

	text := fmt.Sprintf(common.UIMsgWelcomeBackProfessional, currentUser.LastName, currentUser.Role)
	keyboard := h.createProfessionalDashboardKeyboard()

	h.sendMessageWithKeyboard(chatID, text, keyboard)
}
