package api_service

import (
	"booking_client/internal/models"
)

// CreateAppointment creates a new appointment
func (s *APIService) CreateAppointment(req *CreateAppointmentRequest) (*models.CreateAppointmentResponse, error) {
	url := s.buildURL("api", "appointments")

	var response models.CreateAppointmentResponse
	if err := s.makePostRequest(url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
