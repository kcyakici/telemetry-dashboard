package handlers

import (
	"fmt"
	"time"
)

func parseDuration(start string, end string) (time.Duration, error) {
	var duration time.Duration
	if start != "" && end != "" {
		t1, err1 := time.Parse(time.RFC3339, start)
		t2, err2 := time.Parse(time.RFC3339, end)
		if err1 != nil {
			return 0, fmt.Errorf("invalid time format, must be RFC3339: %s", start)
		}
		if err2 != nil {
			return 0, fmt.Errorf("invalid time format, must be RFC3339: %s", end)
		}
		if t1.After(t2) {
			return 0, fmt.Errorf("end date cannot come before start date")
		}

		duration = t2.Sub(t1)

	}
	return duration, nil
}
