package professional

import (
	"context"

	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	apiService "booking_client/internal/services/api_service"
)

// StartSignIn starts the professional sign-in process
func (h *ProfessionalHandler) StartSignIn(ctx context.Context, chatID int64) {
	// Create a temporary user with state
	tempUser := &models.User{
		ChatID: &chatID,
		Role:   "professional",
		State:  models.StateWaitingForUsername,
	}

	// Store in memory for state tracking
	h.apiService.GetUserRepository().SetUser(chatID, tempUser)

	h.sendMessage(chatID, common.UIMsgProfessionalSignIn)
}

// HandleUsernameInput handles username input for professional sign-in
func (h *ProfessionalHandler) HandleUsernameInput(ctx context.Context, chatID int64, username string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.Username = username
	user.State = models.StateWaitingForPassword
	h.apiService.GetUserRepository().SetUser(chatID, user)

	h.sendMessage(chatID, common.SuccessMsgUsernameSaved)
}

// HandlePasswordInput handles password input for professional sign-in
func (h *ProfessionalHandler) HandlePasswordInput(ctx context.Context, chatID int64, password string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Sign in the professional
	req := &apiService.ProfessionalSignInRequest{
		Username: user.Username,
		Password: password,
		ChatID:   chatID,
	}

	signedInUser, err := h.apiService.SignInProfessional(ctx, req)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgSignInFailed, err)
		return
	}

	// Clear state
	user.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, signedInUser)

	// Build success message
	text := common.NewSuccessMessage("sign_in_success").
		WithData("first_name", signedInUser.FirstName).
		WithData("last_name", signedInUser.LastName).
		WithData("role", signedInUser.Role).
		WithData("username", signedInUser.Username).
		WithData("chat_id", chatID).
		Build()

	h.sendMessage(chatID, text)
	h.ShowDashboard(ctx, chatID, signedInUser)
}
