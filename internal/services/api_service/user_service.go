package api_service

import (
	"booking_client/internal/common"
	"booking_client/internal/models"
	"context"
	"fmt"
	"strconv"
)

// GetUserByChatID retrieves a user by their chat ID (checks local storage first, then API)
func (s *APIService) GetUserByChatID(ctx context.Context, chatID int64) (*models.User, error) {
	// First, check local storage
	if user, exists := s.userRepository.GetUser(chatID); exists {
		logger := common.GetLogger(ctx)
		logger.Debug().Int64("chat_id", chatID).Msg("User found in local storage")
		return user, nil
	}

	// If not found locally, fetch from API
	logger := common.GetLogger(ctx)
	logger.Debug().Int64("chat_id", chatID).Msg("User not found in local storage, fetching from API")
	user, err := s.fetchUserFromAPI(ctx, chatID)
	if err != nil {
		return nil, err
	}

	// Store in local storage for future use
	s.userRepository.SetUser(chatID, user)
	logger.Debug().Int64("chat_id", chatID).Msg("User stored in local storage")

	return user, nil
}

// fetchUserFromAPI fetches a user from the API
func (s *APIService) fetchUserFromAPI(ctx context.Context, chatID int64) (*models.User, error) {
	url := s.buildURL("api", "users", strconv.FormatInt(chatID, 10))

	var response struct {
		User models.User `json:"user"`
	}

	requestID := common.GetRequestID(ctx)
	if err := s.makeGetRequestWithContext(ctx, url, &response, requestID); err != nil {
		// Handle 404 specifically for better error message
		if fmt.Sprint(err) == "API returned status 404" {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &response.User, nil
}
