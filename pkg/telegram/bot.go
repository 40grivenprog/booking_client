package telegram

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

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
	workers       int
	updateChan    chan tgbotapi.Update
}

// NewBot creates a new Telegram bot instance
func NewBot(token string, logger *zerolog.Logger) (*Bot, error) {
	return NewBotWithWorkers(token, logger, 1)
}

// NewBotWithWorkers creates a new Telegram bot instance with specified number of workers
func NewBotWithWorkers(token string, logger *zerolog.Logger, workers int) (*Bot, error) {
	// Create secure HTTP client
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	// Create bot API with custom HTTP client
	api, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	// Set debug mode if logger level is debug
	if logger.GetLevel() == zerolog.DebugLevel {
		api.Debug = true
	}

	ctx, cancel := context.WithCancel(context.Background())

	bot := &Bot{
		api:        api,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		workers:    workers,
		updateChan: make(chan tgbotapi.Update, workers*10), // Buffer for workers
	}

	logger.Info().Str("username", api.Self.UserName).Msg("Authorized on account")
	return bot, nil
}

// Start starts the bot and begins polling for updates
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	// Start update receiver
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-b.ctx.Done():
				close(b.updateChan)
				return
			case update := <-updates:
				select {
				case b.updateChan <- update:
				case <-b.ctx.Done():
					return
				}
			}
		}
	}()

	// Start worker pool
	for i := 0; i < b.workers; i++ {
		b.wg.Add(1)
		go b.worker(i)
	}

	b.logger.Info().
		Int("workers", b.workers).
		Msg("Bot started with worker pool")

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

// worker processes updates in a separate goroutine
func (b *Bot) worker(id int) {
	defer b.wg.Done()

	// Panic recovery for worker
	defer func() {
		if r := recover(); r != nil {
			b.logger.Error().
				Int("worker_id", id).
				Interface("panic", r).
				Msg("Worker panic recovered")
		}
	}()

	b.logger.Debug().Int("worker_id", id).Msg("Worker started")

	for {
		select {
		case <-b.ctx.Done():
			b.logger.Debug().Int("worker_id", id).Msg("Worker stopped")
			return
		case update, ok := <-b.updateChan:
			if !ok {
				b.logger.Debug().Int("worker_id", id).Msg("Update channel closed, worker stopping")
				return
			}
			b.handleUpdate(update)
		}
	}
}

// handleUpdate processes incoming updates
func (b *Bot) handleUpdate(update tgbotapi.Update) {
	// Panic recovery for individual update processing
	defer func() {
		if r := recover(); r != nil {
			b.logger.Error().
				Interface("panic", r).
				Msg("Panic recovered in handleUpdate")
		}
	}()

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

// SendMessageWithID sends a message and returns the message ID
func (b *Bot) SendMessageWithID(chatID int64, text string) (int, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	sentMsg, err := b.api.Send(msg)
	if err != nil {
		return 0, err
	}
	return sentMsg.MessageID, nil
}

// SendMessageWithKeyboard sends a message with a custom keyboard
func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	_, err := b.api.Send(msg)
	return err
}

// SendMessageWithKeyboardAndID sends a message with keyboard and returns the message ID
func (b *Bot) SendMessageWithKeyboardAndID(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) (int, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	sentMsg, err := b.api.Send(msg)
	if err != nil {
		return 0, err
	}
	return sentMsg.MessageID, nil
}

// EditMessage edits an existing message
func (b *Bot) EditMessage(chatID int64, messageID int, text string) error {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	_, err := b.api.Send(edit)
	return err
}

// EditMessageWithKeyboard edits an existing message with a custom keyboard
func (b *Bot) EditMessageWithKeyboard(chatID int64, messageID int, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ReplyMarkup = &keyboard
	_, err := b.api.Send(edit)
	return err
}

// DeleteMessage deletes a message
func (b *Bot) DeleteMessage(chatID int64, messageID int) error {
	delete := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := b.api.Send(delete)
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
