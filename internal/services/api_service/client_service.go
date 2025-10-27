package api_service

import (
	"context"
	"net/http"
	"net/url"

	"booking_client/internal/common"
	"booking_client/internal/schemas"
)

// RegisterClient registers a new client
func (s *APIService) RegisterClient(ctx context.Context, req *RegisterRequest) (*schemas.ClientRegisterResponse, error) {
	url := s.buildURL("api", "clients", "register")

	var response schemas.ClientRegisterResponse
	if err := s.makePostRequest(ctx, url, req, &response, http.StatusCreated); err != nil {
		return nil, err
	}

	s.logger.Debug().Int64("chat_id", req.ChatID).Msg("Newly registered client stored in local storage")

	return &response, nil
}

// GetClientAppointments retrieves client appointments with optional status filter
func (s *APIService) GetClientAppointments(ctx context.Context, clientID, status string) (*schemas.GetClientAppointmentsResponse, error) {
	query := url.Values{}
	if status != "" {
		query.Set("status", status)
	}

	url := s.buildURLWithQuery([]string{"api", "clients", clientID, "appointments"}, query)

	var response schemas.GetClientAppointmentsResponse
	requestID := common.GetRequestID(ctx)
	if err := s.makeGetRequestWithContext(ctx, url, &response, requestID); err != nil {
		return nil, err
	}

	return &response, nil
}

// CancelClientAppointment cancels an appointment by client
func (s *APIService) CancelClientAppointment(ctx context.Context, clientID, appointmentID string, req *CancelAppointmentRequest) (*schemas.CancelClientAppointmentResponse, error) {
	url := s.buildURL("api", "clients", clientID, "appointments", appointmentID, "cancel")

	var response schemas.CancelClientAppointmentResponse
	if err := s.makePatchRequest(ctx, url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
