package api_service

import (
	"booking_client/internal/common"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// makeRequest performs an HTTP request with auth header and returns the response body
func (s *APIService) makeRequest(ctx context.Context, method, url string, body interface{}, expectedStatus int) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if err := s.addAuthHeader(req); err != nil {
		return nil, err
	}

	// Add request ID from context
	if requestID := common.GetRequestID(ctx); requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != expectedStatus {
		return nil, s.parseAPIError(resp.StatusCode, respBody)
	}

	return respBody, nil
}

// parseAPIError attempts to parse an API error response
func (s *APIService) parseAPIError(statusCode int, body []byte) error {
	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		// If we can't parse the error response, return the raw body
		return fmt.Errorf("API returned status %d: %s", statusCode, strings.TrimSpace(string(body)))
	}

	// Return structured error
	return &APIError{
		StatusCode: statusCode,
		ErrorType:  errorResp.Error,
		Message:    errorResp.Message,
		RequestID:  errorResp.RequestID,
	}
}

// makePostRequest performs a POST request and unmarshals the response
func (s *APIService) makePostRequest(ctx context.Context, url string, reqBody interface{}, result interface{}, expectedStatus int) error {
	body, err := s.makeRequest(ctx, http.MethodPost, url, reqBody, expectedStatus)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// makePatchRequest performs a PATCH request and unmarshals the response
func (s *APIService) makePatchRequest(ctx context.Context, url string, reqBody interface{}, result interface{}) error {
	body, err := s.makeRequest(ctx, http.MethodPatch, url, reqBody, http.StatusOK)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// makeGetRequestWithContext performs a GET request
func (s *APIService) makeGetRequestWithContext(ctx context.Context, url string, result interface{}, requestID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if err := s.addRequestHeaders(req, requestID); err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return s.parseAPIError(resp.StatusCode, body)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
