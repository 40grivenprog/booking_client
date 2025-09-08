package telegram

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// UpdateHandler defines the interface for handling updates
type UpdateHandler interface {
	HandleUpdate(update tgbotapi.Update)
}

// Bot wraps the Telegram bot API
type Bot struct {
	api           *tgbotapi.BotAPI
	logger        *zerolog.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	updateHandler UpdateHandler
}

// NewBot creates a new Telegram bot instance
func NewBot(token string, logger *zerolog.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	// Set debug mode if logger level is debug
	if logger.GetLevel() == zerolog.DebugLevel {
		api.Debug = true
	}

	ctx, cancel := context.WithCancel(context.Background())

	bot := &Bot{
		api:    api,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}

	logger.Info().Str("username", api.Self.UserName).Msg("Authorized on account")
	return bot, nil
}

// Start starts the bot and begins polling for updates
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-b.ctx.Done():
				return
			case update := <-updates:
				b.handleUpdate(update)
			}
		}
	}()

	return nil
}

// Stop gracefully stops the bot
func (b *Bot) Stop() {
	b.cancel()
	b.wg.Wait()
}

// SetUpdateHandler sets the update handler
func (b *Bot) SetUpdateHandler(handler UpdateHandler) {
	b.updateHandler = handler
}

// handleUpdate processes incoming updates
func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if b.updateHandler != nil {
		b.updateHandler.HandleUpdate(update)
	} else if update.Message != nil {
		b.logger.Debug().
			Str("message", update.Message.Text).
			Int64("user_id", update.Message.From.ID).
			Msg("Received message")
	}
}

// SendMessage sends a message to a specific chat
func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	return err
}

// SendMessageWithKeyboard sends a message with a custom keyboard
func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	_, err := b.api.Send(msg)
	return err
}

// GetAPI returns the underlying bot API for advanced operations
func (b *Bot) GetAPI() *tgbotapi.BotAPI {
	return b.api
}

// GetLogger returns the logger instance
func (b *Bot) GetLogger() *zerolog.Logger {
	return b.logger
}
