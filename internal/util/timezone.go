package util

import (
	"time"
)

var (
	// AppTimezone is the timezone used throughout the client application
	AppTimezone *time.Location
)

// InitTimezone initializes the client application timezone
func InitTimezone() error {
	// Load specific timezone (UTC+2 for Europe/Berlin)
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		// Fallback to local timezone if timezone loading fails
		AppTimezone = time.Local
		return err
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
