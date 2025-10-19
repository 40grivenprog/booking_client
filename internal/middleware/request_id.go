package middleware

import (
	"context"

	"booking_client/internal/common"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// RequestIDAndLoggerMiddleware creates a context with request_id and adjusted logger
func RequestIDAndLoggerMiddleware(ctx context.Context, baseLogger zerolog.Logger) (context.Context, zerolog.Logger) {
	requestID := uuid.New().String()

	// Create adjusted logger with request context
	adjustedLogger := baseLogger.With().
		Str("request_id", requestID).
		Logger()

	ctx = common.WithRequestID(ctx, requestID)
	ctx = common.WithLogger(ctx, adjustedLogger)

	adjustedLogger.Info().Msg("Request started")

	return ctx, adjustedLogger
}
