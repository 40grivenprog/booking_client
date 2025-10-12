# 🤖 Booking Client - Telegram Bot

A production-ready Telegram bot for the booking application, providing interactive UI for clients and professionals to manage appointments. Built with Go, following clean architecture and modular design principles.

---

## ✨ Features

### Client Features
- 📝 **Client Registration** - Simple registration flow with name and phone
- 👨‍⚕️ **Professional Selection** - Browse and select professionals
- 📅 **Appointment Booking** - Interactive date and time selection
- 📋 **View Appointments** - See upcoming and pending appointments
- ❌ **Cancel Appointments** - Cancel bookings with reason
- ⌨️ **Rich Keyboard UI** - Inline keyboards for better UX

### Professional Features
- 🔐 **Sign In** - Username/password authentication
- 📊 **Dashboard** - Professional control panel
- ✅ **Confirm Appointments** - Approve pending bookings
- ❌ **Cancel Appointments** - Cancel with reason
- 📅 **View Schedule** - See upcoming appointments by date
- 📈 **Timetable** - Daily schedule view with all appointments
- 🚫 **Set Unavailable** - Mark time periods as unavailable
- 🗓️ **Calendar Navigation** - Month/date navigation for appointments

### Architecture & Code Quality
- 🏗️ **Clean Architecture** - Separation of concerns (Handlers → Services → Repository pattern)
- 📦 **Modular API Service** - Domain-specific services (ClientService, ProfessionalService, etc.)
- 🔄 **Callback Router** - Modular callback query handling (replaces large switch statements)
- ⌨️ **Keyboard Builders** - Dedicated modules for keyboard creation
- 💬 **Message Builders** - Consistent message formatting with builder pattern
- 🔐 **JWT Token Generation** - Secure API authentication
- 🛡️ **Error Handling** - Structured API error parsing and user-friendly messages
- 📝 **Structured Logging** - zerolog with context

### Production Ready
- 🐳 **Containerized** - Docker ready
- ☸️ **Kubernetes Ready** - Helm charts for deployment
- ⚙️ **Configuration** - Environment-based configuration
- 🔒 **Secure** - JWT authentication, no hardcoded secrets
- 📊 **Observability** - Structured logging

---

## 🛠️ Tech Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.23+ |
| **Bot Framework** | go-telegram-bot-api/v5 |
| **HTTP Client** | net/http (custom modular client) |
| **Authentication** | JWT (golang-jwt/jwt/v5) |
| **Logging** | zerolog |
| **Configuration** | godotenv |

---

## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- Telegram Bot Token (from @BotFather)
- Running booking_api instance

### Local Development

#### 1. Create Telegram Bot
```bash
# 1. Open Telegram and message @BotFather
# 2. Send /newbot command
# 3. Follow instructions to create bot
# 4. Copy the bot token (looks like: 123456789:ABCdefGHIjklMNOpqrsTUVwxyz)
```

#### 2. Setup Environment
```bash
cd booking_client

# Create .env file
cat > .env << EOF
   # Telegram Bot Configuration
   TELEGRAM_BOT_TOKEN=your_bot_token_here
   
   # API Configuration
   API_BASE_URL=http://localhost:8080
JWT_SECRET=your-jwt-secret  # Must match booking_api secret

# Logging
LOG_LEVEL=info
EOF
```

#### 3. Start the Bot
   ```bash
# Install dependencies
go mod tidy

# Run the bot
make run

# Or directly with go
go run cmd/bot/main.go
```

#### 4. Test the Bot
   ```bash
# Open Telegram and search for your bot
# Send /start command
# You should see the welcome message
```

---

## 📱 User Guide

### For Clients

#### Registration Flow
1. Send `/start` to the bot
2. Click "🆕 I'm a Client"
3. Enter your first name
4. Enter your last name
5. Enter your phone number
6. ✅ Registration complete!

#### Booking Appointment
1. From dashboard, click "📅 Book Appointment"
2. Select a professional from the list
3. Choose a date from the calendar
4. Select available time slot
5. ✅ Appointment created (status: pending)

#### View Appointments
- **Upcoming**: See all confirmed appointments
- **Pending**: See appointments waiting for confirmation
- Each shows: professional name, date, time

#### Cancel Appointment
1. Go to "📋 My Appointments"
2. Select appointment to cancel
3. Confirm cancellation
4. ✅ Appointment cancelled

---

### For Professionals

#### Sign In Flow
1. Send `/start` to the bot
2. Click "👨‍⚕️ I'm a Professional"
3. Enter your username
4. Enter your password
5. ✅ Signed in - Dashboard opens

#### Dashboard Options
- **📋 Pending Requests** - View pending appointments (awaiting confirmation)
- **📅 Upcoming Appointments** - View upcoming confirmed appointments
- **📆 Timetable** - View daily schedule
- **🚫 Set Unavailable** - Mark time as unavailable

#### Confirm/Cancel Appointments
1. Click "📋 Pending Requests"
2. Select appointment
3. Click "✅ Confirm" or "❌ Cancel"
4. If cancelling, provide reason
5. ✅ Action complete

#### View Timetable
1. Click "📆 Timetable"
2. Select month (navigation: ◀️ prev, ▶️ next)
3. Select date with appointments (marked with 📅)
4. View all appointments for that day

#### Set Unavailable
1. Click "🚫 Set Unavailable"
2. Select start date
3. Select start time
4. Select end date
5. Select end time
6. (Optional) Add description
7. ✅ Period marked as unavailable

---

## 🏗️ Architecture

### Clean Architecture

```
┌─────────────────────────────────────────┐
│      Telegram Bot API                   │
│      (User Interface)                   │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Handler Layer                   │
│  - Client Handlers                      │
│  - Professional Handlers                │
│  - Callback Router                      │
│  - Keyboard Builders                    │
│  - Message Builders                     │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      API Service Layer                  │
│  - ClientService                        │
│  - ProfessionalService                  │
│  - UserService                          │
│  - AppointmentService                   │
│  - HTTP Helpers                         │
│  - JWT Token Generator                  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Booking API                     │
│         (REST API)                      │
└─────────────────────────────────────────┘
```

### Directory Structure

```
booking_client/
├── cmd/bot/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration loading
│   ├── handlers/
│   │   ├── handler.go       # Main handler + router setup
│   │   ├── client/          # Client-side handlers
│   │   │   ├── client_handler.go
│   │   │   ├── registration_handler.go
│   │   │   ├── book_appointment_handler.go
│   │   │   ├── upcoming_appointments_handler.go
│   │   │   ├── pending_appointments_handler.go
│   │   │   ├── cancel_appointment_handler.go
│   │   │   └── helpers.go
│   │   ├── professional/    # Professional-side handlers
│   │   │   ├── professional_handler.go
│   │   │   ├── sign_in_handler.go
│   │   │   ├── pending_appointments_handler.go
│   │   │   ├── upcoming_appointments_handler.go
│   │   │   ├── timetable_handler.go
│   │   │   ├── unavailable_handler.go
│   │   │   └── helpers.go
│   │   ├── keyboards/       # Keyboard builders
│   │   │   ├── client_keyboards.go
│   │   │   └── professional_keyboards.go
│   │   ├── router/          # Callback router
│   │   │   └── callback_router.go
│   │   └── common/          # Shared utilities
│   │       ├── callbacks.go      # Callback constants
│   │       ├── constants.go      # Message constants
│   │       ├── helpers.go        # Helper functions
│   │       ├── message_builder.go # Message builders
│   │       └── notification_service.go
│   ├── services/
│   │   └── api_service/     # Modular API client
│   │       ├── service.go        # Core service
│   │       ├── http_helpers.go   # HTTP methods
│   │       ├── errors.go         # Error types
│   │       ├── error_helpers.go  # Error utilities
│   │       ├── schema.go         # Request DTOs
│   │       ├── user_service.go
│   │       ├── client_service.go
│   │       ├── professional_service.go
│   │       └── appointment_service.go
│   ├── token/               # JWT token generation
│   │   ├── maker.go
│   │   ├── jwt_maker.go
│   │   └── payload.go
│   ├── models/
│   │   ├── user.go
│   │   └── constants.go
│   ├── repository/
│   │   └── user_repository.go
│   └── util/
│       └── timezone.go
├── pkg/telegram/
│   └── bot.go               # Telegram bot wrapper
├── Dockerfile
├── Makefile
└── README.md
```

---

## 🔧 Configuration

### Environment Variables

```bash
# Required
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
API_BASE_URL=http://booking-api:8080
JWT_SECRET=your-jwt-secret-must-match-api

# Optional
LOG_LEVEL=info              # debug, info, warn, error
API_TIMEOUT=30s             # HTTP client timeout
```

### Docker

```bash
# Build image
docker build -t booking-client:latest .

# Run container
docker run -d \
  -e TELEGRAM_BOT_TOKEN="your-token" \
  -e API_BASE_URL="http://booking-api:8080" \
  -e JWT_SECRET="your-secret" \
  booking-client:latest
```

---

## 🔄 Callback Router System

### Before (Old Approach)
```go
// Large switch statement (hard to maintain)
switch data {
case "client":
    // handle
case "professional":
    // handle
case "select_professional_123":
    // handle
// ... 50+ more cases
}
```

### After (Modular Router)
```go
// Setup routes once
func (h *Handler) setupRoutes() {
    // Exact matches
    h.callbackRouter.RegisterExact("client", h.clientHandler.StartRegistration)
    h.callbackRouter.RegisterExact("professional", h.professionalHandler.StartSignIn)
    
    // Prefix matches
    h.callbackRouter.RegisterPrefix("select_professional_", h.handleProfessionalSelection)
}

// Route callback
h.callbackRouter.Route(chatID, callbackData)
```

**Benefits:**
- ✅ Modular and extensible
- ✅ Easy to test individual handlers
- ✅ Type-safe routing
- ✅ Clear separation of concerns

---

## ⌨️ Keyboard Builders

Centralized keyboard creation for consistency:

```go
// Client keyboards
keyboards.CreateDateKeyboard(dates, month)
keyboards.CreateTimeKeyboard(times, professionalID, date)
keyboards.CreateDashboardKeyboard()

// Professional keyboards
keyboards.CreateProfessionalDashboardKeyboard()
keyboards.CreateUnavailableDateKeyboard(month, dates)
keyboards.CreateTimetableKeyboard(date, appointments)
```

**Benefits:**
- ✅ Consistent UI across the app
- ✅ Easy to update keyboard layouts
- ✅ Reusable components
- ✅ Centralized button text management

---

## 💬 Message Builders

Consistent message formatting with builder pattern:

```go
// Success message
message := common.NewSuccessMessage("appointment_created").
    WithData("professional", "Dr. Smith").
    WithData("time", "10:00-11:00").
    Build()

// Appointment message
message := common.NewClientAppointmentMessage(appointment).
    Build()

// Professional appointment message
message := common.NewProfessionalAppointmentMessage(appointment).
    Build()
```

**Benefits:**
- ✅ Consistent message formatting
- ✅ Easy to maintain and update
- ✅ Type-safe message construction
- ✅ Reduces duplication

---

## 🔐 JWT Authentication

### Token Generation

```go
// Create token maker
tokenMaker, err := token.NewJWTMaker(jwtSecret)

// Generate token
token, _, err := tokenMaker.CreateToken("booking-client", 15*time.Minute)

// Use in API requests
req.Header.Set("Authorization", "Bearer " + token)
```

### Token Structure

```json
{
  "service": "booking-client",
  "issued_at": "2024-01-15T10:00:00Z",
  "expired_at": "2024-01-15T10:15:00Z"
}
```

---

## 🛡️ Error Handling

### API Error Parsing

```go
// Make API request
data, err := apiService.GetProfessionals()
if err != nil {
    // Check if it's an API error
    if apiService.IsAPIError(err) {
        // Format user-friendly message
        userMsg := apiService.FormatErrorForUser(err)
        
        // Get request ID for support
        requestID := apiService.GetRequestID(err)
        
        bot.Send(tgbotapi.NewMessage(chatID, userMsg))
        return
    }
    
    // Handle other errors
    log.Error().Err(err).Msg("unexpected error")
}
```

### User-Friendly Messages

API errors are automatically converted to user-friendly messages:
- `not_found` → "Resource not found"
- `validation_error` → "Invalid data"
- `conflict` → "Data conflict"
- `unauthorized` → "Authentication error"
- `internal_error` → "Internal server error"

---

## 🛠️ Development

### Commands

```bash
# Build binary
make build

# Run bot
make run

# Run in development mode
make dev

# Install dependencies
make deps

# Format code
make fmt

# Lint code
make lint

# Run tests
make test

# Clean build artifacts
make clean

# Show help
make help
```

### Adding New Handler

1. **Create handler file** in `internal/handlers/client/` or `internal/handlers/professional/`
2. **Add handler method** to respective handler struct
3. **Register route** in `handler.go` `setupRoutes()`
4. **Add keyboard** in `internal/handlers/keyboards/` if needed
5. **Add constants** in `internal/handlers/common/callbacks.go`

Example:
```go
// 1. Create method in professional_handler.go
func (h *ProfessionalHandler) HandleNewFeature(chatID int64, messageID int) {
    // Implementation
}

// 2. Register in setupRoutes()
h.callbackRouter.RegisterExact(common.CallbackNewFeature, func(chatID int64, _ string) {
    h.professionalHandler.HandleNewFeature(chatID, 0)
})

// 3. Add constant in callbacks.go
const CallbackNewFeature = "new_feature"
```

---

## 📊 Logging

### Structured Logging with zerolog

```go
// Info log
log.Info().
    Int64("chat_id", chatID).
    Str("action", "appointment_created").
    Msg("user created appointment")

// Error log
log.Error().
    Err(err).
    Int64("chat_id", chatID).
    Str("endpoint", "/api/appointments").
    Msg("failed to create appointment")

// Debug log (only in debug mode)
log.Debug().
    Interface("appointment", appointment).
    Msg("appointment details")
```

---

## 🚀 Deployment

### Kubernetes

See [deployment guide](../deployments/DEPLOYMENT_GUIDE.md) for:
- Helm charts
- ArgoCD setup
- Multi-environment deployment
- Secrets management

### Quick Deploy

```bash
# Build and push image
docker build -t YOUR_ECR_URL/booking-client:v1.0.0 .
docker push YOUR_ECR_URL/booking-client:v1.0.0

# Update image tag
./deployments/scripts/update-client-tag.sh v1.0.0 dev

# Commit and push
git add deployments/helm/booking-client/values-dev.yaml
git commit -m "deploy: booking-client v1.0.0 to dev"
git push

# ArgoCD will auto-deploy
```

---

## 🧪 Testing

### Manual Testing

```bash
# Start the bot locally
make run

# In Telegram:
# 1. Search for your bot
# 2. Send /start
# 3. Test client flow: register → book → view → cancel
# 4. Test professional flow: sign in → confirm → view timetable → set unavailable
```

### Test Scenarios

**Client Flow:**
1. ✅ Registration with valid data
2. ✅ Book appointment
3. ✅ View pending appointments
4. ✅ View upcoming appointments
5. ✅ Cancel appointment

**Professional Flow:**
1. ✅ Sign in with valid credentials
2. ✅ View pending appointments
3. ✅ Confirm appointment
4. ✅ Cancel appointment with reason
5. ✅ View timetable
6. ✅ Set unavailable period

---

## 🔒 Security

### Best Practices

1. **Never commit secrets** - Use environment variables
2. **JWT Secret** - Must match booking_api secret
3. **Telegram Token** - Keep confidential, regenerate if exposed
4. **Validate Input** - Always validate user input before API calls
5. **Error Messages** - Don't expose internal details to users
6. **Rate Limiting** - Consider rate limiting for production

---

## 🤝 Contributing

1. Follow clean architecture principles
2. Use message builders for user messages
3. Use keyboard builders for keyboards
4. Register new routes in callback router
5. Add constants to `callbacks.go`
6. Use structured logging
7. Handle errors gracefully
8. Test both client and professional flows

---

## 📄 License

[Your License Here]

---

**Built with ❤️ and Go**
