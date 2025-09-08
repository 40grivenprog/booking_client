# Booking Client - Telegram Bot

A Telegram bot for the booking application that allows users to register, book appointments, and manage their bookings.

## Features

- User registration (Client/Professional)
- Appointment booking
- View appointments
- Interactive commands

## Setup

1. **Create a Telegram Bot:**
   - Message @BotFather on Telegram
   - Create a new bot with `/newbot`
   - Get your bot token

2. **Environment Configuration:**
   Create a `.env` file in the root directory:
   ```env
   # Telegram Bot Configuration
   TELEGRAM_BOT_TOKEN=your_bot_token_here
   
   # API Configuration
   API_BASE_URL=http://localhost:8080
   
   # Application Configuration
   DEBUG=false
   PORT=8081
   ```
   
   You can also use the Makefile to create the .env file:
   ```bash
   make env
   ```

3. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

4. **Run the Bot:**
   ```bash
   go run cmd/bot/main.go
   ```

## Available Commands

- `/start` - Start the bot and see welcome message
- `/help` - Show help message with all commands
- `/register` - Register as a client or professional
- `/book` - Book a new appointment
- `/my_appointments` - View your appointments

## Project Structure

```
booking_client/
├── cmd/bot/           # Main application entry point
├── internal/
│   ├── config/        # Configuration management
│   ├── handlers/      # Command handlers
│   ├── models/        # Data models
│   └── services/      # Business logic services
├── pkg/telegram/      # Telegram bot wrapper
└── README.md
```

## Development

### Available Make Commands

```bash
make build    # Build the bot binary
make run      # Run the bot
make clean    # Clean build artifacts
make test     # Run tests
make deps     # Install dependencies
make fmt      # Format code
make lint     # Lint code
make env      # Create .env file from example
make help     # Show help message
```

### Tech Stack

The bot is built with:
- [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) - Telegram Bot API wrapper
- [zerolog](https://github.com/rs/zerolog) - Fast, structured, leveled logging
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading

### Architecture

The bot follows a clean architecture pattern:
- **cmd/bot** - Application entry point
- **internal/config** - Configuration management
- **internal/handlers** - Command and message handlers
- **internal/models** - Data models
- **internal/services** - Business logic and API communication
- **pkg/telegram** - Telegram bot wrapper and utilities
