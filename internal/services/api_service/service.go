package api_service

import (
	"booking_client/internal/config"
	"booking_client/internal/repository"
	"booking_client/internal/token"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/rs/zerolog"
)

// APIService handles communication with the booking API and local storage
type APIService struct {
	baseURL        string
	client         *http.Client
	logger         *zerolog.Logger
	userRepository *repository.UserRepository
	tokenMaker     token.Maker
}

// NewAPIService creates a new API service
func NewAPIService(config *config.Config, logger *zerolog.Logger) (*APIService, error) {
	// Create token maker
	tokenMaker, err := token.NewJWTMaker(config.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	return &APIService{
		baseURL: config.APIBaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:         logger,
		userRepository: repository.NewUserRepository(),
		tokenMaker:     tokenMaker,
	}, nil
}

// addAuthHeader adds the JWT authorization header to the request
func (s *APIService) addAuthHeader(req *http.Request) error {
	authToken, err := s.tokenMaker.CreateToken("booking_client", 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to create auth token: %w", err)
	}
	fmt.Println("authToken", authToken)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	return nil
}

// addRequestHeaders adds common headers including request_id to the request
func (s *APIService) addRequestHeaders(req *http.Request, requestID string) error {
	// Add auth header
	if err := s.addAuthHeader(req); err != nil {
		return err
	}

	// Add request_id header
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	return nil
}

// GetUserRepository returns the user repository for direct access if needed
func (s *APIService) GetUserRepository() *repository.UserRepository {
	return s.userRepository
}

// buildURL constructs a URL from the base URL and path segments
func (s *APIService) buildURL(pathSegments ...string) string {
	u, _ := url.Parse(s.baseURL)
	u.Path = path.Join(append([]string{u.Path}, pathSegments...)...)
	return u.String()
}

// buildURLWithQuery constructs a URL with query parameters
func (s *APIService) buildURLWithQuery(pathSegments []string, query url.Values) string {
	u, _ := url.Parse(s.baseURL)
	u.Path = path.Join(append([]string{u.Path}, pathSegments...)...)
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	return u.String()
}
