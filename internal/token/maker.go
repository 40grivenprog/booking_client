package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific service and duration
	CreateToken(service string, duration time.Duration) (string, error)
}
