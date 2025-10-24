package main

import (
	"os"
	"os/signal"
	"syscall"

	"booking_client/internal/config"
	"booking_client/internal/handlers"
	"booking_client/internal/util"
	"booking_client/pkg/telegram"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize timezone
	if err := util.InitTimezone(); err != nil {
		log.Warn().Err(err).Msg("Failed to load timezone, falling back to local timezone")
	} else {
		log.Info().Str("timezone", util.GetAppTimezone().String()).Msg("Timezone initialized")
	}

	// Configure logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Use JSON format for production logs
	log.Logger = log.Output(os.Stdout)

	// Initialize Telegram bot with worker pool (5 workers)
	bot, err := telegram.NewBotWithWorkers(cfg.TelegramToken, &log.Logger, 5)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Telegram bot")
	}

	// Initialize handlers
	handler, err := handlers.NewHandler(bot, cfg, &log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize handlers")
	}

	// Register command handlers
	handler.RegisterHandlers()

	log.Info().Msg("Starting Telegram bot...")

	// Start the bot
	if err := bot.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start bot")
	}

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down bot...")
	bot.Stop()
}
