package client

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/models"
	"booking_client/internal/schema"
	"fmt"
)

// StartRegistration starts the client registration process
func (h *ClientHandler) StartRegistration(chatID int64) {
	// Create a temporary user with state
	tempUser := &models.User{
		ChatID: &chatID,
		Role:   "client",
		State:  models.StateWaitingForFirstName,
	}

	id, err := h.bot.SendMessageWithID(chatID, common.UIMsgClientRegistration)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}

	tempUser.LastMessageID = &id
	tempUser.MessagesToDelete = append(tempUser.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, tempUser)
}

// HandleFirstNameInput handles first name input for client registration
func (h *ClientHandler) HandleFirstNameInput(chatID int64, firstName string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.FirstName = firstName
	user.State = models.StateWaitingForLastName
	h.apiService.GetUserRepository().SetUser(chatID, user)

	id, err := h.bot.SendMessageWithID(chatID, common.SuccessMsgFirstNameSaved)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
}

// HandleLastNameInput handles last name input for client registration
func (h *ClientHandler) HandleLastNameInput(chatID int64, lastName string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.LastName = lastName
	user.State = models.StateWaitingForPhone
	h.apiService.GetUserRepository().SetUser(chatID, user)

	id, err := h.bot.SendMessageWithID(chatID, common.SuccessMsgLastNameSaved)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
}

// HandlePhoneInput handles phone number input for client registration
func (h *ClientHandler) HandlePhoneInput(chatID int64, phone string) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
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

	response, err := h.apiService.RegisterClient(req)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgRegistrationFailed, err)
		return
	}
	user.ID = response.ID
	user.FirstName = response.FirstName
	user.LastName = response.LastName
	user.Role = response.Role
	user.PhoneNumber = response.PhoneNumber
	user.ChatID = response.ChatID
	user.CreatedAt = response.CreatedAt
	user.UpdatedAt = response.UpdatedAt
	// Clear state
	user.State = models.StateNone
	h.apiService.GetUserRepository().SetUser(chatID, user)

	text := fmt.Sprintf(common.SuccessMsgRegistrationSuccessful, response.FirstName, response.LastName, response.Role, chatID)
	keyboard := h.createRegistrationSuccessKeyboard()

	id, err := h.bot.SendMessageWithKeyboardAndID(chatID, text, keyboard)
	if err != nil {
		h.sendError(chatID, common.ErrorMsgFailedToSendMessage, err)
		return
	}
	user.LastMessageID = &id
	user.MessagesToDelete = append(user.MessagesToDelete, &id)
	h.apiService.GetUserRepository().SetUser(chatID, user)
}
