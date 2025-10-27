package common

import (
	"booking_client/internal/models"
	"booking_client/internal/schemas"
	"fmt"
	"time"
)

// AppointmentMessage builds formatted messages for appointments
type AppointmentMessage struct {
	appointment interface{}
	index       int
}

// NewClientAppointmentMessage creates a message builder for client appointments
func NewClientAppointmentMessage(apt *schemas.ClientAppointment, index int) *AppointmentMessage {
	return &AppointmentMessage{
		appointment: apt,
		index:       index,
	}
}

// NewProfessionalAppointmentMessage creates a message builder for professional appointments
func NewProfessionalAppointmentMessage(apt *schemas.ProfessionalAppointment, index int) *AppointmentMessage {
	return &AppointmentMessage{
		appointment: apt,
		index:       index,
	}
}

// ForClient formats appointment for client view
func (m *AppointmentMessage) ForClient() string {
	apt, ok := m.appointment.(*schemas.ClientAppointment)
	if !ok {
		return ""
	}

	date, startTime, endTime := FormatAppointmentTime(apt.StartTime, apt.EndTime)
	return fmt.Sprintf("✍️ Appointment #%d:\n📅 %s\n🕐 %s - %s\n👨‍💼 %s %s\n📝 %s\n\n",
		m.index+1, date, startTime, endTime,
		apt.Professional.FirstName, apt.Professional.LastName,
		apt.Description)
}

// ForProfessional formats appointment for professional view
func (m *AppointmentMessage) ForProfessional() string {
	apt, ok := m.appointment.(*schemas.ProfessionalAppointment)
	if !ok {
		return ""
	}

	date, startTime, endTime := FormatAppointmentTime(apt.StartTime, apt.EndTime)
	return fmt.Sprintf("✍️ Appointment #%d:\n📅 %s\n🕐 %s - %s\n👤 Client: %s %s\n📝 %s\n\n",
		m.index+1, date, startTime, endTime,
		apt.Client.FirstName, apt.Client.LastName,
		apt.Description)
}

// TimetableSlotMessage builds formatted messages for timetable slots
type TimetableSlotMessage struct {
	slot  *schemas.TimetableAppointment
	index int
}

// NewTimetableSlotMessage creates a message builder for timetable slots
func NewTimetableSlotMessage(slot *schemas.TimetableAppointment, index int) *TimetableSlotMessage {
	return &TimetableSlotMessage{
		slot:  slot,
		index: index,
	}
}

// Format formats timetable slot for display
func (m *TimetableSlotMessage) Format() string {
	startTime, _ := time.Parse(time.RFC3339, m.slot.StartTime)
	endTime, _ := time.Parse(time.RFC3339, m.slot.EndTime)

	return fmt.Sprintf("%d. 🕐 %s - %s | %s\n",
		m.index+1,
		startTime.Format("15:04"),
		endTime.Format("15:04"),
		m.slot.Description)
}

// SuccessMessage builds success messages
type SuccessMessage struct {
	messageType string
	data        map[string]interface{}
}

// NewSuccessMessage creates a success message builder
func NewSuccessMessage(messageType string) *SuccessMessage {
	return &SuccessMessage{
		messageType: messageType,
		data:        make(map[string]interface{}),
	}
}

// WithData adds data to the message
func (m *SuccessMessage) WithData(key string, value interface{}) *SuccessMessage {
	m.data[key] = value
	return m
}

// Build generates the formatted success message
func (m *SuccessMessage) Build() string {
	switch m.messageType {
	case "appointment_confirmed":
		date := m.data["date"].(string)
		startTime := m.data["start_time"].(string)
		endTime := m.data["end_time"].(string)
		clientFirstName := m.data["client_first_name"].(string)
		clientLastName := m.data["client_last_name"].(string)
		return fmt.Sprintf("✅ Appointment confirmed!\n📅 %s\n🕐 %s - %s\n👤 Client: %s %s",
			date, startTime, endTime, clientFirstName, clientLastName)

	case "appointment_cancelled":
		date := m.data["date"].(string)
		startTime := m.data["start_time"].(string)
		endTime := m.data["end_time"].(string)
		firstName := m.data["first_name"].(string)
		lastName := m.data["last_name"].(string)
		reason := m.data["reason"].(string)
		return fmt.Sprintf("✅ Appointment cancelled\n📅 %s\n🕐 %s - %s\n👤 %s %s\n📝 Reason: %s",
			date, startTime, endTime, firstName, lastName, reason)

	case "appointment_created":
		date := m.data["date"].(string)
		startTime := m.data["start_time"].(string)
		endTime := m.data["end_time"].(string)
		professionalFirstName := m.data["professional_first_name"].(string)
		professionalLastName := m.data["professional_last_name"].(string)
		return fmt.Sprintf("✅ Your appointment has been created!\n📅 %s\n🕐 %s - %s\n👨‍💼 Professional: %s %s\n\n⏳ Waiting for confirmation...",
			date, startTime, endTime, professionalFirstName, professionalLastName)

	case "unavailable_period_set":
		date := m.data["date"].(string)
		startTime := m.data["start_time"].(string)
		endTime := m.data["end_time"].(string)
		description := m.data["description"].(string)
		return fmt.Sprintf("✅ Unavailable period set successfully!\n📅 %s\n🕐 %s - %s\n📝 %s",
			date, startTime, endTime, description)

	case "registration_success":
		firstName := m.data["first_name"].(string)
		lastName := m.data["last_name"].(string)
		role := m.data["role"].(string)
		return fmt.Sprintf("✅ Registration successful!\n\n👤 Name: %s %s\n🎭 Role: %s\n\nWelcome aboard! 🎉",
			firstName, lastName, role)

	case "sign_in_success":
		firstName := m.data["first_name"].(string)
		lastName := m.data["last_name"].(string)
		role := m.data["role"].(string)
		username := m.data["username"].(string)
		chatID := m.data["chat_id"].(int64)
		return fmt.Sprintf("✅ Sign in successful!\n\n👤 Name: %s %s\n🎭 Role: %s\n👔 Username: %s\n💬 Chat ID: %d",
			firstName, lastName, role, username, chatID)

	default:
		return "✅ Success!"
	}
}

// WelcomeMessage builds welcome messages
type WelcomeMessage struct {
	user *models.User
}

// NewWelcomeMessage creates a welcome message builder
func NewWelcomeMessage(user *models.User) *WelcomeMessage {
	return &WelcomeMessage{user: user}
}

// ForClient builds welcome message for client
func (m *WelcomeMessage) ForClient() string {
	return fmt.Sprintf("👋 Welcome back, %s!\n\nWhat would you like to do?", m.user.FirstName)
}

// ForProfessional builds welcome message for professional
func (m *WelcomeMessage) ForProfessional() string {
	return fmt.Sprintf("👋 Welcome back, %s!\n\nYou are logged in as: %s\n\nWhat would you like to do?",
		m.user.LastName, m.user.Role)
}

// TimetableMessage builds timetable header messages
type TimetableMessage struct {
	date  string
	slots []schemas.TimetableAppointment
}

// NewTimetableMessage creates a timetable message builder
func NewTimetableMessage(date string, slots []schemas.TimetableAppointment) *TimetableMessage {
	return &TimetableMessage{
		date:  date,
		slots: slots,
	}
}

// BuildHeader builds timetable header
func (m *TimetableMessage) BuildHeader() string {
	dateObj, _ := time.Parse("2006-01-02", m.date)
	formattedDate := dateObj.Format("Monday, January 2, 2006")

	if len(m.slots) == 0 {
		return fmt.Sprintf("📅 Timetable for %s\n\n🎉 No appointments scheduled for this day!", formattedDate)
	}

	return fmt.Sprintf("📅 Timetable for %s\n\nYour appointments:\n\n", formattedDate)
}

// BuildFull builds complete timetable message with slots
func (m *TimetableMessage) BuildFull() string {
	header := m.BuildHeader()
	if len(m.slots) == 0 {
		return header
	}

	message := header
	for i, slot := range m.slots {
		message += NewTimetableSlotMessage(&slot, i).Format()
	}
	return message
}
