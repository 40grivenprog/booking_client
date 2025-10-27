package handlers

import (
	"context"

	handlersCommon "booking_client/internal/handlers/common"
)

// setupRoutes registers all callback handlers with the router
func (h *Handler) setupRoutes() {
	// Initial selection
	h.callbackRouter.RegisterExact(handlersCommon.CallbackClient, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.clientHandler.StartRegistration(ctx, chatID, messageID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackProfessional, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.professionalHandler.StartSignIn(ctx, chatID, messageID)
	})

	// Client callbacks
	h.callbackRouter.RegisterExact(handlersCommon.CallbackBookAppointment, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.clientHandler.HandleBookAppointment(ctx, chatID, messageID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackPendingAppointments, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.clientHandler.HandlePendingAppointments(ctx, chatID, messageID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackUpcomingAppointments, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.clientHandler.HandleUpcomingAppointments(ctx, chatID, messageID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackCancelBooking, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.clientHandler.HandleCancelBooking(ctx, chatID, messageID)
	})

	// Professional callbacks
	h.callbackRouter.RegisterExact(handlersCommon.CallbackProfessionalPendingAppointments, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.professionalHandler.HandlePendingAppointments(ctx, chatID, messageID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackProfessionalUpcomingAppointments, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.professionalHandler.HandleUpcomingAppointments(ctx, chatID, messageID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackProfessionalTimetable, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.professionalHandler.HandleTimetable(ctx, chatID, messageID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackSetUnavailable, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.professionalHandler.HandleSetUnavailable(ctx, chatID, messageID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackProfessionalPreviousAppointments, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.professionalHandler.HandlePreviousAppointments(ctx, chatID, messageID)
	})

	// Unavailable navigation
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixPrevUnavailableMonth, func(ctx context.Context, chatID int64, month string, messageID int) {
		h.professionalHandler.HandleUnavailableMonthNavigation(ctx, chatID, month, handlersCommon.DirectionPrev, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixNextUnavailableMonth, func(ctx context.Context, chatID int64, month string, messageID int) {
		h.professionalHandler.HandleUnavailableMonthNavigation(ctx, chatID, month, handlersCommon.DirectionNext, messageID)
	})
	h.callbackRouter.RegisterExact(handlersCommon.CallbackCancelUnavailable, func(ctx context.Context, chatID int64, _ string, messageID int) {
		h.professionalHandler.HandleCancelUnavailable(ctx, chatID, messageID)
	})

	// Professional timetable navigation
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixPrevTimetableDay, func(ctx context.Context, chatID int64, date string, messageID int) {
		h.professionalHandler.HandleTimetableDateNavigation(ctx, chatID, date, handlersCommon.DirectionPrev, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixNextTimetableDay, func(ctx context.Context, chatID int64, date string, messageID int) {
		h.professionalHandler.HandleTimetableDateNavigation(ctx, chatID, date, handlersCommon.DirectionNext, messageID)
	})

	// Professional upcoming appointments navigation
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixPrevUpcomingMonth, func(ctx context.Context, chatID int64, month string, messageID int) {
		h.professionalHandler.HandleUpcomingAppointmentsMonthNavigation(ctx, chatID, month, handlersCommon.DirectionPrev, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixNextUpcomingMonth, func(ctx context.Context, chatID int64, month string, messageID int) {
		h.professionalHandler.HandleUpcomingAppointmentsMonthNavigation(ctx, chatID, month, handlersCommon.DirectionNext, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectUpcomingDate, func(ctx context.Context, chatID int64, date string, messageID int) {
		h.professionalHandler.HandleUpcomingAppointmentsDateSelection(ctx, chatID, date, messageID)
	})

	// Client booking flow - month navigation
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixPrevMonth, func(ctx context.Context, chatID int64, month string, messageID int) {
		h.clientHandler.HandleBookAppointmentsMonthNavigation(ctx, chatID, month, handlersCommon.DirectionPrev, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixNextMonth, func(ctx context.Context, chatID int64, month string, messageID int) {
		h.clientHandler.HandleBookAppointmentsMonthNavigation(ctx, chatID, month, handlersCommon.DirectionNext, messageID)
	})

	// Selection callbacks
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectProfessional, func(ctx context.Context, chatID int64, professionalID string, messageID int) {
		h.clientHandler.HandleProfessionalSelection(ctx, chatID, professionalID, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectDate, func(ctx context.Context, chatID int64, date string, messageID int) {
		h.clientHandler.HandleDateSelection(ctx, chatID, date, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectTime, func(ctx context.Context, chatID int64, startTime string, messageID int) {
		h.clientHandler.HandleTimeSelection(ctx, chatID, startTime, messageID)
	})

	// Appointment actions
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixCancelAppointment, func(ctx context.Context, chatID int64, appointmentID string, messageID int) {
		h.clientHandler.HandleCancelAppointment(ctx, chatID, appointmentID, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixConfirmAppointment, func(ctx context.Context, chatID int64, appointmentID string, messageID int) {
		h.professionalHandler.HandleConfirmAppointment(ctx, chatID, appointmentID, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixCancelProfAppt, func(ctx context.Context, chatID int64, appointmentID string, messageID int) {
		h.professionalHandler.HandleCancelAppointment(ctx, chatID, appointmentID, messageID)
	})

	// Unavailable flow
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectUnavailableDate, func(ctx context.Context, chatID int64, date string, messageID int) {
		h.professionalHandler.HandleUnavailableDateSelection(ctx, chatID, date, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectUnavailableStart, func(ctx context.Context, chatID int64, startTime string, messageID int) {
		h.professionalHandler.HandleUnavailableStartTimeSelection(ctx, chatID, startTime, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectUnavailableEnd, func(ctx context.Context, chatID int64, endTime string, messageID int) {
		h.professionalHandler.HandleUnavailableEndTimeSelection(ctx, chatID, endTime, messageID)
	})

	// Previous appointments flow
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixSelectClient, func(ctx context.Context, chatID int64, clientID string, messageID int) {
		h.professionalHandler.HandleClientSelection(ctx, chatID, clientID, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixPrevPreviousMonth, func(ctx context.Context, chatID int64, month string, messageID int) {
		h.professionalHandler.HandlePreviousAppointmentsMonthNavigation(ctx, chatID, month, handlersCommon.DirectionPrev, messageID)
	})
	h.callbackRouter.RegisterPrefix(handlersCommon.CallbackPrefixNextPreviousMonth, func(ctx context.Context, chatID int64, month string, messageID int) {
		h.professionalHandler.HandlePreviousAppointmentsMonthNavigation(ctx, chatID, month, handlersCommon.DirectionNext, messageID)
	})

	// Back to dashboard (special case - needs user lookup)
	h.callbackRouter.RegisterExact(handlersCommon.CallbackBackToDashboard, func(ctx context.Context, chatID int64, _ string, messageID int) {
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
			h.professionalHandler.ShowDashboard(ctx, chatID, user, messageID)
		} else {
			h.clientHandler.ShowDashboard(ctx, chatID, messageID)
		}
	})
}
