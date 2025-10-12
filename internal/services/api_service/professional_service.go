package api_service

import (
	"booking_client/internal/models"
	"net/url"
)

// RegisterProfessional registers a new professional
func (s *APIService) RegisterProfessional(req *RegisterRequest) (*models.User, error) {
	url := s.buildURL("api", "professionals", "register")

	var response struct {
		User models.User `json:"user"`
	}

	if err := s.makePostRequest(url, req, &response); err != nil {
		return nil, err
	}

	// Store the newly registered user in local storage
	s.userRepository.SetUser(req.ChatID, &response.User)
	s.logger.Debug().Int64("chat_id", req.ChatID).Msg("Newly registered professional stored in local storage")

	return &response.User, nil
}

// SignInProfessional signs in a professional user
func (s *APIService) SignInProfessional(req *ProfessionalSignInRequest) (*models.User, error) {
	url := s.buildURL("api", "professionals", "sign_in")

	var response struct {
		User models.User `json:"user"`
	}

	if err := s.makePostRequest(url, req, &response); err != nil {
		return nil, err
	}

	// Store the signed-in user in local storage
	s.userRepository.SetUser(req.ChatID, &response.User)
	s.logger.Debug().Int64("chat_id", req.ChatID).Msg("Professional signed in and stored in local storage")

	return &response.User, nil
}

// GetProfessionals retrieves all professionals
func (s *APIService) GetProfessionals() (*models.GetProfessionalsResponse, error) {
	url := s.buildURL("api", "professionals")

	var response models.GetProfessionalsResponse
	if err := s.makeGetRequest(url, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalAvailability retrieves availability for a professional on a specific date
func (s *APIService) GetProfessionalAvailability(professionalID, date string) (*models.ProfessionalAvailabilityResponse, error) {
	query := url.Values{}
	query.Set("date", date)

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "availability"}, query)

	var response models.ProfessionalAvailabilityResponse
	if err := s.makeGetRequest(url, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalAppointments retrieves professional appointments with optional status filter
func (s *APIService) GetProfessionalAppointments(professionalID, status string) (*models.GetProfessionalAppointmentsResponse, error) {
	query := url.Values{}
	if status != "" {
		query.Set("status", status)
	}

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "appointments"}, query)

	var response models.GetProfessionalAppointmentsResponse
	if err := s.makeGetRequest(url, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalAppointmentDates retrieves dates with confirmed appointments for a month
func (s *APIService) GetProfessionalAppointmentDates(professionalID, month string) (*models.GetProfessionalAppointmentDatesResponse, error) {
	query := url.Values{}
	if month != "" {
		query.Set("month", month)
	}

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "appointment_dates"}, query)

	var response models.GetProfessionalAppointmentDatesResponse
	if err := s.makeGetRequest(url, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ConfirmProfessionalAppointment confirms an appointment by professional
func (s *APIService) ConfirmProfessionalAppointment(professionalID, appointmentID string, req *ConfirmAppointmentRequest) (*models.ConfirmProfessionalAppointmentResponse, error) {
	url := s.buildURL("api", "professionals", professionalID, "appointments", appointmentID, "confirm")

	var response models.ConfirmProfessionalAppointmentResponse
	if err := s.makePatchRequest(url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CancelProfessionalAppointment cancels an appointment by professional
func (s *APIService) CancelProfessionalAppointment(professionalID, appointmentID string, req *CancelAppointmentRequest) (*models.CancelProfessionalAppointmentResponse, error) {
	url := s.buildURL("api", "professionals", professionalID, "appointments", appointmentID, "cancel")

	var response models.CancelProfessionalAppointmentResponse
	if err := s.makePatchRequest(url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateUnavailableAppointment creates an unavailable appointment for a professional
func (s *APIService) CreateUnavailableAppointment(req *CreateUnavailableAppointmentRequest) (*models.CreateUnavailableAppointmentResponse, error) {
	url := s.buildURL("api", "professionals", req.ProfessionalID, "unavailable_appointments")

	var response models.CreateUnavailableAppointmentResponse
	if err := s.makePostRequest(url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalAppointmentsByDate retrieves professional appointments with status and date filters
func (s *APIService) GetProfessionalAppointmentsByDate(professionalID, status, date string) (*models.GetProfessionalAppointmentsResponse, error) {
	query := url.Values{}
	query.Set("status", status)
	query.Set("date", date)

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "appointments"}, query)

	var response models.GetProfessionalAppointmentsResponse
	if err := s.makeGetRequest(url, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProfessionalTimetable gets the professional's timetable for a specific date
func (s *APIService) GetProfessionalTimetable(professionalID, date string) (*models.GetProfessionalTimetableResponse, error) {
	query := url.Values{}
	query.Set("date", date)

	url := s.buildURLWithQuery([]string{"api", "professionals", professionalID, "timetable"}, query)

	var response models.GetProfessionalTimetableResponse
	if err := s.makeGetRequest(url, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
