package professional

import (
	"context"

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
func (h *ProfessionalHandler) ShowDashboard(ctx context.Context, chatID int64, user *models.User) {
	currentUser, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	currentUser.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, currentUser)

	text := fmt.Sprintf(common.UIMsgWelcomeBackProfessional, currentUser.LastName, currentUser.Role)
	keyboard := h.createProfessionalDashboardKeyboard()

	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// ShowDashboardWithEdit shows the professional dashboard by editing the last message
func (h *ProfessionalHandler) ShowDashboardWithEdit(chatID int64, user *models.User) {
	currentUser, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	currentUser.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, currentUser)

	text := fmt.Sprintf(common.UIMsgWelcomeBackProfessional, currentUser.LastName, currentUser.Role)
	keyboard := h.createProfessionalDashboardKeyboard()

	// If user has a last message ID, edit it; otherwise send a new message
	if currentUser.LastMessageID != nil {
		h.editMessageWithKeyboard(chatID, *currentUser.LastMessageID, text, keyboard)
	} else {
		messageID, err := h.sendMessageWithKeyboardAndID(chatID, text, keyboard)
		if err == nil {
			currentUser.LastMessageID = &messageID
			h.apiService.GetUserRepository().SetUser(chatID, currentUser)
		}
	}
}
