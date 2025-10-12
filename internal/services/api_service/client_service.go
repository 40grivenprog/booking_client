package api_service

import (
	"booking_client/internal/models"
	"net/url"
)

// RegisterClient registers a new client
func (s *APIService) RegisterClient(req *RegisterRequest) (*models.ClientRegisterResponse, error) {
	url := s.buildURL("api", "clients", "register")

	var response models.ClientRegisterResponse
	if err := s.makePostRequest(url, req, &response); err != nil {
		return nil, err
	}

	s.logger.Debug().Int64("chat_id", req.ChatID).Msg("Newly registered client stored in local storage")

	return &response, nil
}

// GetClientAppointments retrieves client appointments with optional status filter
func (s *APIService) GetClientAppointments(clientID, status string) (*models.GetClientAppointmentsResponse, error) {
	query := url.Values{}
	if status != "" {
		query.Set("status", status)
	}

	url := s.buildURLWithQuery([]string{"api", "clients", clientID, "appointments"}, query)

	var response models.GetClientAppointmentsResponse
	if err := s.makeGetRequest(url, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CancelClientAppointment cancels an appointment by client
func (s *APIService) CancelClientAppointment(clientID, appointmentID string, req *CancelAppointmentRequest) (*models.CancelClientAppointmentResponse, error) {
	url := s.buildURL("api", "clients", clientID, "appointments", appointmentID, "cancel")

	var response models.CancelClientAppointmentResponse
	if err := s.makePatchRequest(url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
