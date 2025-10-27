package api_service

import (
	"context"
	"net/http"

	"booking_client/internal/schemas"
)

// CreateAppointment creates a new appointment
func (s *APIService) CreateAppointment(ctx context.Context, req *CreateAppointmentRequest) (*schemas.CreateAppointmentResponse, error) {
	url := s.buildURL("api", "appointments")

	var response schemas.CreateAppointmentResponse
	if err := s.makePostRequest(ctx, url, req, &response, http.StatusCreated); err != nil {
		return nil, err
	}

	return &response, nil
}
