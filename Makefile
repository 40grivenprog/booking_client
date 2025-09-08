# Makefile for Booking Client Telegram Bot

.PHONY: build run clean test deps help

# Default target
all: build

# Build the bot
build:
	@echo "Building booking client bot..."
	go build -o bin/bot cmd/bot/main.go

# Run the bot
run:
	@echo "Running booking client bot..."
	go run cmd/bot/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Create .env file from example
env:
	@echo "Creating .env file from example..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file. Please edit it with your configuration."; \
	else \
		echo ".env file already exists."; \
	fi

# Help
help:
	@echo "Available targets:"
	@echo "  build    - Build the bot binary"
	@echo "  run      - Run the bot"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  deps     - Install dependencies"
	@echo "  fmt      - Format code"
	@echo "  lint     - Lint code"
	@echo "  env      - Create .env file from example"
	@echo "  help     - Show this help message"
