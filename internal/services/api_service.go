package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"booking_client/internal/config"
	"booking_client/internal/models"
	"booking_client/internal/repository"
	"booking_client/internal/schema"

	"github.com/rs/zerolog"
)

// APIService handles communication with the booking API and local storage
type APIService struct {
	baseURL        string
	client         *http.Client
	logger         *zerolog.Logger
	userRepository *repository.UserRepository
}

// NewAPIService creates a new API service
func NewAPIService(config *config.Config, logger *zerolog.Logger) *APIService {
	return &APIService{
		baseURL: config.APIBaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:         logger,
		userRepository: repository.NewUserRepository(),
	}
}

// GetUserByChatID retrieves a user by their chat ID (checks local storage first, then API)
func (s *APIService) GetUserByChatID(chatID int64) (*models.User, error) {
	// First, check local storage
	if user, exists := s.userRepository.GetUser(chatID); exists {
		s.logger.Debug().Int64("chat_id", chatID).Msg("User found in local storage")
		return user, nil
	}

	// If not found locally, fetch from API
	s.logger.Debug().Int64("chat_id", chatID).Msg("User not found in local storage, fetching from API")
	user, err := s.fetchUserFromAPI(chatID)
	if err != nil {
		return nil, err
	}

	// Store in local storage for future use
	s.userRepository.SetUser(chatID, user)
	s.logger.Debug().Int64("chat_id", chatID).Msg("User stored in local storage")

	return user, nil
}

// fetchUserFromAPI fetches a user from the API
func (s *APIService) fetchUserFromAPI(chatID int64) (*models.User, error) {
	url := fmt.Sprintf("%s/api/users/%d", s.baseURL, chatID)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		User models.User `json:"user"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.User, nil
}

// RegisterClient registers a new client
func (s *APIService) RegisterClient(req *schema.RegisterRequest) (*models.User, error) {
	url := fmt.Sprintf("%s/api/clients/register", s.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		User models.User `json:"user"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Store the newly registered user in local storage
	s.userRepository.SetUser(req.ChatID, &response.User)
	s.logger.Debug().Int64("chat_id", req.ChatID).Msg("Newly registered client stored in local storage")

	return &response.User, nil
}

// RegisterProfessional registers a new professional
func (s *APIService) RegisterProfessional(req *schema.RegisterRequest) (*models.User, error) {
	url := fmt.Sprintf("%s/api/professionals/register", s.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		User models.User `json:"user"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Store the newly registered user in local storage
	s.userRepository.SetUser(req.ChatID, &response.User)
	s.logger.Debug().Int64("chat_id", req.ChatID).Msg("Newly registered professional stored in local storage")

	return &response.User, nil
}

// SignInProfessional signs in a professional user
func (s *APIService) SignInProfessional(req *schema.ProfessionalSignInRequest) (*models.User, error) {
	url := fmt.Sprintf("%s/api/professionals/sign_in", s.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		User models.User `json:"user"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Store the signed-in user in local storage
	s.userRepository.SetUser(req.ChatID, &response.User)
	s.logger.Debug().Int64("chat_id", req.ChatID).Msg("Professional signed in and stored in local storage")

	return &response.User, nil
}

// GetProfessionals retrieves all professionals
func (s *APIService) GetProfessionals() (*models.GetProfessionalsResponse, error) {
	url := fmt.Sprintf("%s/api/professionals", s.baseURL)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.GetProfessionalsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetProfessionalAvailability retrieves availability for a professional on a specific date
func (s *APIService) GetProfessionalAvailability(professionalID, date string) (*models.ProfessionalAvailabilityResponse, error) {
	url := fmt.Sprintf("%s/api/professionals/%s/availability?date=%s", s.baseURL, professionalID, date)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.ProfessionalAvailabilityResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// CreateAppointment creates a new appointment
func (s *APIService) CreateAppointment(req *schema.CreateAppointmentRequest) (*models.CreateAppointmentResponse, error) {
	url := fmt.Sprintf("%s/api/appointments", s.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		// Try to parse error response for better user experience
		var errorResp struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error.Message != "" {
			return nil, fmt.Errorf("%s", errorResp.Error.Message)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.CreateAppointmentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetClientAppointments retrieves client appointments with optional status filter
func (s *APIService) GetClientAppointments(clientID, status string) (*models.GetClientAppointmentsResponse, error) {
	url := fmt.Sprintf("%s/api/clients/%s/appointments", s.baseURL, clientID)
	if status != "" {
		url += "?status=" + status
	}

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.GetClientAppointmentsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// CancelClientAppointment cancels an appointment by client
func (s *APIService) CancelClientAppointment(clientID, appointmentID string, req *schema.CancelAppointmentRequest) (*models.CancelClientAppointmentResponse, error) {
	url := fmt.Sprintf("%s/api/clients/%s/appointments/%s/cancel", s.baseURL, clientID, appointmentID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.CancelClientAppointmentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetProfessionalAppointments retrieves professional appointments with optional status filter
func (s *APIService) GetProfessionalAppointments(professionalID, status string) (*models.GetProfessionalAppointmentsResponse, error) {
	url := fmt.Sprintf("%s/api/professionals/%s/appointments", s.baseURL, professionalID)
	if status != "" {
		url += "?status=" + status
	}

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.GetProfessionalAppointmentsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// ConfirmProfessionalAppointment confirms an appointment by professional
func (s *APIService) ConfirmProfessionalAppointment(professionalID, appointmentID string, req *schema.ConfirmAppointmentRequest) (*models.ConfirmProfessionalAppointmentResponse, error) {
	url := fmt.Sprintf("%s/api/professionals/%s/appointments/%s/confirm", s.baseURL, professionalID, appointmentID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.ConfirmProfessionalAppointmentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// CancelProfessionalAppointment cancels an appointment by professional
func (s *APIService) CancelProfessionalAppointment(professionalID, appointmentID string, req *schema.CancelAppointmentRequest) (*models.CancelProfessionalAppointmentResponse, error) {
	url := fmt.Sprintf("%s/api/professionals/%s/appointments/%s/cancel", s.baseURL, professionalID, appointmentID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.CancelProfessionalAppointmentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetUserRepository returns the user repository for direct access if needed
func (s *APIService) GetUserRepository() *repository.UserRepository {
	return s.userRepository
}
