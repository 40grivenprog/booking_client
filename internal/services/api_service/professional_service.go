package api_service

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"booking_client/internal/common"
	"booking_client/internal/models"
	"booking_client/internal/schemas"
)

// RegisterProfessional registers a new professional
func (s *APIService) RegisterProfessional(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	url := s.buildURL("api", "professionals", "register")

	var response struct {
		User models.User `json:"user"`
	}

	if err := s.makePostRequest(ctx, url, req, &response, http.StatusCreated); err != nil {
		return nil, err
	}

	// Store the newly registered user in local storage
	s.userRepository.SetUser(req.ChatID, &response.User)
	s.logger.Debug().Int64("chat_id", req.ChatID).Msg("Newly registered professional stored in local storage")

	return &response.User, nil
}

// SignInProfessional signs in a professional user
func (s *APIService) SignInProfessional(ctx context.Context, req *ProfessionalSignInRequest) (*models.User, error) {
	url := s.buildURL("api", "professionals", "sign_in")

	var response struct {
		User models.User `json:"user"`
	}

	if err := s.makePostRequest(ctx, url, req, &response, http.StatusOK); err != nil {
		return nil, err
	}

	// Store the signed-in user in local storage
	s.userRepository.SetUser(req.ChatID, &response.User)
	s.logger.Debug().Int64("chat_id", req.ChatID).Msg("Professional signed in and stored in local storage")

	return &response.User, nil
}

// GetProfessionals retrieves all professionals
func (s *APIService) GetProfessionals(ctx context.Context) (*schemas.GetProfessionalsResponse, error) {
	url := s.buildURL("api", "professionals")

	var response schemas.GetProfessionalsResponse
	requestID := common.GetRequestID(ctx)
	if err := s.makeGetRequestWithContext(ctx, url, &response, requestID); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalAvailability retrieves availability for a professional on a specific date
func (s *APIService) GetProfessionalAvailability(ctx context.Context, professionalID, date string) (*schemas.ProfessionalAvailabilityResponse, error) {
	query := url.Values{}
	query.Set("date", date)

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "availability"}, query)

	var response schemas.ProfessionalAvailabilityResponse
	requestID := common.GetRequestID(ctx)
	if err := s.makeGetRequestWithContext(ctx, url, &response, requestID); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalAppointments retrieves professional appointments with optional status filter
func (s *APIService) GetProfessionalAppointments(ctx context.Context, professionalID, status string) (*schemas.GetProfessionalAppointmentsResponse, error) {
	query := url.Values{}
	if status != "" {
		query.Set("status", status)
	}

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "appointments"}, query)

	var response schemas.GetProfessionalAppointmentsResponse
	requestID := common.GetRequestID(ctx)
	if err := s.makeGetRequestWithContext(ctx, url, &response, requestID); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalAppointmentDates retrieves dates with confirmed appointments for a month
func (s *APIService) GetProfessionalAppointmentDates(ctx context.Context, professionalID, month string) (*schemas.GetProfessionalAppointmentDatesResponse, error) {
	query := url.Values{}
	if month != "" {
		query.Set("month", month)
	}

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "appointment_dates"}, query)

	var response schemas.GetProfessionalAppointmentDatesResponse
	requestID := common.GetRequestID(ctx)
	if err := s.makeGetRequestWithContext(ctx, url, &response, requestID); err != nil {
		return nil, err
	}

	return &response, nil
}

// ConfirmProfessionalAppointment confirms an appointment by professional
func (s *APIService) ConfirmProfessionalAppointment(ctx context.Context, professionalID, appointmentID string, req *ConfirmAppointmentRequest) (*schemas.ConfirmProfessionalAppointmentResponse, error) {
	url := s.buildURL("api", "professionals", professionalID, "appointments", appointmentID, "confirm")

	var response schemas.ConfirmProfessionalAppointmentResponse
	if err := s.makePatchRequest(ctx, url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CancelProfessionalAppointment cancels an appointment by professional
func (s *APIService) CancelProfessionalAppointment(ctx context.Context, professionalID, appointmentID string, req *CancelAppointmentRequest) (*schemas.CancelProfessionalAppointmentResponse, error) {
	url := s.buildURL("api", "professionals", professionalID, "appointments", appointmentID, "cancel")

	var response schemas.CancelProfessionalAppointmentResponse
	if err := s.makePatchRequest(ctx, url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateUnavailableAppointment creates an unavailable appointment for a professional
func (s *APIService) CreateUnavailableAppointment(ctx context.Context, req *CreateUnavailableAppointmentRequest) (*schemas.CreateUnavailableAppointmentResponse, error) {
	url := s.buildURL("api", "professionals", req.ProfessionalID, "unavailable_appointments")

	var response schemas.CreateUnavailableAppointmentResponse
	if err := s.makePostRequest(ctx, url, req, &response, http.StatusCreated); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalAppointmentsByDate retrieves professional appointments with status and date filters
func (s *APIService) GetProfessionalAppointmentsByDate(ctx context.Context, professionalID, status, date string) (*schemas.GetProfessionalAppointmentsResponse, error) {
	query := url.Values{}
	query.Set("status", status)
	query.Set("date", date)

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "appointments"}, query)

	var response schemas.GetProfessionalAppointmentsResponse
	requestID := common.GetRequestID(ctx)
	if err := s.makeGetRequestWithContext(ctx, url, &response, requestID); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalTimetable gets the professional's timetable for a specific date
func (s *APIService) GetProfessionalTimetable(ctx context.Context, professionalID, date string) (*schemas.GetProfessionalTimetableResponse, error) {
	query := url.Values{}
	query.Set("date", date)

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "timetable"}, query)

	var response schemas.GetProfessionalTimetableResponse
	requestID := common.GetRequestID(ctx)
	if err := s.makeGetRequestWithContext(ctx, url, &response, requestID); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalClients retrieves all clients for a professional
func (s *APIService) GetProfessionalClients(ctx context.Context, professionalID string) ([]schemas.ProfessionalClient, error) {
	url := s.buildURL("api", "professionals", professionalID, "clients")

	var response struct {
		Clients []schemas.ProfessionalClient `json:"clients"`
	}

	requestID := common.GetRequestID(ctx)
	if err := s.makeGetRequestWithContext(ctx, url, &response, requestID); err != nil {
		return nil, err
	}

	return response.Clients, nil
}

// GetPreviousAppointmentsByClient retrieves previous appointments for a specific client
func (s *APIService) GetPreviousAppointmentsByClient(ctx context.Context, professionalID, clientID string, month *time.Time) ([]schemas.PreviousAppointment, error) {
	query := url.Values{}
	query.Set("client_id", clientID)
	if month != nil {
		query.Set("month", month.Format("2006-01"))
	}

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "previous_appointments"}, query)

	var response schemas.GetPreviousAppointmentsByClientResponse
	requestID := common.GetRequestID(ctx)
	if err := s.makeGetRequestWithContext(ctx, url, &response, requestID); err != nil {
		return nil, err
	}

	return response.Appointments, nil
}
