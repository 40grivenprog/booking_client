package util

import (
	"fmt"
	"time"
)

var (
	// AppTimezone is the timezone used throughout the client application
	AppTimezone *time.Location
)

// InitTimezone initializes the client application timezone
func InitTimezone() error {
	// Try to load specific timezone (UTC+2 for Europe/Berlin)
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		// Try alternative timezone names
		alternatives := []string{"Europe/Berlin", "CET", "CEST", "UTC+2"}
		for _, tz := range alternatives {
			if loc, err = time.LoadLocation(tz); err == nil {
				AppTimezone = loc
				return nil
			}
		}

		// Fallback to local timezone if all timezone loading fails
		AppTimezone = time.Local
		return fmt.Errorf("unknown time zone Europe/Berlin: %w", err)
	}

	AppTimezone = loc
	return nil
}

// GetAppTimezone returns the application timezone
func GetAppTimezone() *time.Location {
	if AppTimezone == nil {
		// Initialize if not already done
		InitTimezone()
	}
	return AppTimezone
}

// NowInAppTimezone returns current time in application timezone
func NowInAppTimezone() time.Time {
	return time.Now().In(GetAppTimezone())
}

// ConvertToAppTimezone converts a time to application timezone for storage
func ConvertToAppTimezone(t time.Time) time.Time {
	return t.In(GetAppTimezone())
}
