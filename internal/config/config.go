package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds all configuration for the application
type Config struct {
	// Telegram Bot config
	TelegramToken string `env:"TELEGRAM_BOT_TOKEN" envDefault:""`

	// API config
	APIBaseURL string `env:"API_BASE_URL" envDefault:"http://localhost:8080"`

	// JWT config
	JWTSecret string `env:"JWT_SECRET" envDefault:""`

	// Debug config
	Debug bool `env:"DEBUG" envDefault:"false"`

	// Port config (for webhook if needed)
	Port int `env:"PORT" envDefault:"8081"`

	// Log config
	LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"json"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	// Validate required fields
	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}
