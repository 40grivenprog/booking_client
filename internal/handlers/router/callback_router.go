package router

import (
	"strings"

	"github.com/rs/zerolog"
)

// CallbackHandler is a function that handles a callback with a parameter
type CallbackHandler func(chatID int64, param string)

// PrefixHandler holds a prefix and its associated handler
type PrefixHandler struct {
	Prefix  string
	Handler CallbackHandler
}

// CallbackRouter routes callback queries to registered handlers
type CallbackRouter struct {
	exactHandlers  map[string]CallbackHandler
	prefixHandlers []PrefixHandler
	logger         *zerolog.Logger
}

// NewCallbackRouter creates a new callback router
func NewCallbackRouter(logger *zerolog.Logger) *CallbackRouter {
	return &CallbackRouter{
		exactHandlers:  make(map[string]CallbackHandler),
		prefixHandlers: []PrefixHandler{},
		logger:         logger,
	}
}

// RegisterExact registers a handler for an exact callback match
func (r *CallbackRouter) RegisterExact(callback string, handler CallbackHandler) {
	r.exactHandlers[callback] = handler
	r.logger.Debug().Str("callback", callback).Msg("Registered exact callback handler")
}

// RegisterPrefix registers a handler for a callback prefix
// The handler will receive the part after the prefix as a parameter
func (r *CallbackRouter) RegisterPrefix(prefix string, handler CallbackHandler) {
	r.prefixHandlers = append(r.prefixHandlers, PrefixHandler{
		Prefix:  prefix,
		Handler: handler,
	})
	r.logger.Debug().Str("prefix", prefix).Msg("Registered prefix callback handler")
}

// Route routes the callback to the appropriate handler
// Returns true if a handler was found and executed, false otherwise
func (r *CallbackRouter) Route(chatID int64, callbackData string) bool {
	// Try exact match first (most common case)
	if handler, exists := r.exactHandlers[callbackData]; exists {
		r.logger.Debug().
			Int64("chat_id", chatID).
			Str("callback", callbackData).
			Msg("Routing to exact handler")
		handler(chatID, "")
		return true
	}

	// Try prefix matches
	for _, ph := range r.prefixHandlers {
		if strings.HasPrefix(callbackData, ph.Prefix) {
			param := callbackData[len(ph.Prefix):]
			r.logger.Debug().
				Int64("chat_id", chatID).
				Str("prefix", ph.Prefix).
				Str("param", param).
				Msg("Routing to prefix handler")
			ph.Handler(chatID, param)
			return true
		}
	}

	// No handler found
	r.logger.Warn().
		Int64("chat_id", chatID).
		Str("callback", callbackData).
		Msg("No handler found for callback")
	return false
}

// GetStats returns statistics about registered handlers
func (r *CallbackRouter) GetStats() (exactCount int, prefixCount int) {
	return len(r.exactHandlers), len(r.prefixHandlers)
}
