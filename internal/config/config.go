package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	TelegramToken string
	APIBaseURL    string
	JWTSecret     string
	Debug         bool
	Port          int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{}

	// Telegram Bot Token
	config.TelegramToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if config.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// API Base URL
	config.APIBaseURL = os.Getenv("API_BASE_URL")
	if config.APIBaseURL == "" {
		config.APIBaseURL = "http://localhost:8080" // Default API URL
	}

	// JWT Secret (must match booking_api secret)
	config.JWTSecret = os.Getenv("JWT_SECRET")
	if config.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	// Debug mode
	debugStr := os.Getenv("DEBUG")
	if debugStr != "" {
		debug, err := strconv.ParseBool(debugStr)
		if err != nil {
			return nil, fmt.Errorf("invalid DEBUG value: %v", err)
		}
		config.Debug = debug
	}

	// Port (for webhook if needed)
	portStr := os.Getenv("PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT value: %v", err)
		}
		config.Port = port
	} else {
		config.Port = 8081 // Default port
	}

	return config, nil
}
