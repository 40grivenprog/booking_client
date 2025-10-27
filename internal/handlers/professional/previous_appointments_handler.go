package professional

import (
	"booking_client/internal/handlers/common"
	"booking_client/internal/handlers/keyboards"
	"context"
	"fmt"
	"time"
)

// HandlePreviousAppointments shows the list of clients for the professional
func (h *ProfessionalHandler) HandlePreviousAppointments(ctx context.Context, chatID int64, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}

	// Delete message for dashboard
	go func() {
		time.Sleep(1 * time.Second)
		h.bot.DeleteMessage(chatID, messageID)
	}()

	// Get clients for this professional
	clients, err := h.apiService.GetProfessionalClients(ctx, user.ID)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToRetrieveClients, err)
		return
	}

	if len(clients) == 0 {
		h.sendError(ctx, chatID, "No clients found.", nil)
		return
	}

	// Create clients keyboard
	keyboard := keyboards.CreateClientsKeyboard(clients)

	text := "👥 Select a client to view their previous appointments:"

	h.sendMessageWithKeyboard(chatID, text, keyboard)
}

// HandleClientSelection handles when a client is selected
func (h *ProfessionalHandler) HandleClientSelection(ctx context.Context, chatID int64, clientID string, messageID int) {
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok {
		return
	}
	user.LastMessageID = &messageID
	user.MessagesToDelete = append(user.MessagesToDelete, &messageID)

	// Store selected client ID in user state
	user.SelectedClientID = &clientID
	h.apiService.GetUserRepository().SetUser(chatID, user)

	// Show appointments for current month
	currentMonth := time.Now()
	h.showAppointmentsForMonth(ctx, chatID, user.ID, clientID, currentMonth, messageID)
}

// HandlePreviousMonthNavigation handles month navigation for previous appointments
func (h *ProfessionalHandler) HandlePreviousAppointmentsMonthNavigation(ctx context.Context, chatID int64, monthStr string, direction string, messageID int) {
	// Parse month - this is already the target month from the callback
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		h.sendError(ctx, chatID, "Invalid month format.", err)
		return
	}
	h.bot.DeleteMessage(chatID, messageID)

	// Get professional ID and client ID from user
	user, ok := common.GetUserOrSendError(h.apiService.GetUserRepository(), h.bot, h.logger, chatID)
	if !ok || user.SelectedClientID == nil {
		h.sendError(ctx, chatID, "User session or selected client not found.", nil)
		return
	}
	user.MessagesToDelete = append(user.MessagesToDelete, &messageID)
	h.apiService.GetUserRepository().SetUser(chatID, user)

	professionalID := user.ID
	clientID := *user.SelectedClientID

	// Use the month directly (it's already the correct target month from callback)
	h.showAppointmentsForMonth(ctx, chatID, professionalID, clientID, month, messageID)
}

// showAppointmentsForMonth shows appointments for a specific month
func (h *ProfessionalHandler) showAppointmentsForMonth(ctx context.Context, chatID int64, professionalID, clientID string, month time.Time, messageID int) {
	// Get appointments for this month
	appointments, err := h.apiService.GetPreviousAppointmentsByClient(ctx, professionalID, clientID, &month)
	if err != nil {
		h.sendError(ctx, chatID, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	// Create navigation keyboard
	keyboard := keyboards.CreatePreviousAppointmentsNavigationKeyboard(month, len(appointments) > 0)

	// Format appointments text
	text := fmt.Sprintf("📅 Previous appointments for %s %s:\n\n", month.Format("January"), month.Format("2006"))

	if len(appointments) == 0 {
		text += "No appointments found for this month."
	} else {
		for _, apt := range appointments {
			startTime, _ := time.Parse(time.RFC3339, apt.StartTime)
			endTime, _ := time.Parse(time.RFC3339, apt.EndTime)

			text += fmt.Sprintf("📅 %s\n🕐 %s - %s\n",
				startTime.Format("January 2, 2006"),
				startTime.Format("15:04"),
				endTime.Format("15:04"))

			if apt.Description != "" {
				text += fmt.Sprintf("📝 %s\n", apt.Description)
			}
			text += "\n"
		}
	}

	h.sendMessageWithKeyboard(chatID, text, keyboard)
}
