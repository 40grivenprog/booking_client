package api_service

import (
	"context"
	"net/http"

	"booking_client/internal/models"
)

// CreateAppointment creates a new appointment
func (s *APIService) CreateAppointment(ctx context.Context, req *CreateAppointmentRequest) (*models.CreateAppointmentResponse, error) {
	url := s.buildURL("api", "appointments")

	var response models.CreateAppointmentResponse
	if err := s.makePostRequest(ctx, url, req, &response, http.StatusCreated); err != nil {
		return nil, err
	}

	return &response, nil
}
