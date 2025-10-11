package token

import "time"

// Payload contains the token claims data
type Payload struct {
	Service   string    `json:"service"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific service and duration
func NewPayload(service string, duration time.Duration) *Payload {
	return &Payload{
		Service:   service,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
}
