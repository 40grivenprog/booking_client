package api_service

import (
	"context"

	"booking_client/internal/models"
)

// CreateAppointment creates a new appointment
func (s *APIService) CreateAppointment(ctx context.Context, req *CreateAppointmentRequest) (*models.CreateAppointmentResponse, error) {
	url := s.buildURL("api", "appointments")

	var response models.CreateAppointmentResponse
	if err := s.makePostRequest(url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
