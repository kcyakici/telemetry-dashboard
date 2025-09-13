package handlers

import (
	"fmt"
	"time"
)

// ParseAndNormalizeTime parses a query parameter into UTC time
func ParseAndNormalizeTime(raw string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}

	// Try strict RFC3339 first (with timezone, e.g. 2019-06-24T18:00:00Z)
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		utc := t.UTC()
		return &utc, nil
	}

	// Try "datetime-local" format (browser usually sends yyyy-MM-ddTHH:mm, no TZ)
	// Example: "2019-06-24T18:00"
	if t, err := time.Parse("2006-01-02T15:04", raw); err == nil {
		utc := t.UTC()
		return &utc, nil
	}

	// Try with seconds but no TZ
	if t, err := time.Parse("2006-01-02T15:04:05", raw); err == nil {
		utc := t.UTC()
		return &utc, nil
	}

	return nil, fmt.Errorf("invalid timestamp format: %s", raw)
}
