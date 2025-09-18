package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type QueryFilters struct {
	VehicleID string
	Start     time.Time
	End       time.Time
}

var aggregatedTables = map[string]string{
	"speed": "trend_speed_1min",
	"temp":  "trend_temp_1min",
	"power": "trend_power_1min",
}

var aggregatedColumns = map[string]string{
	"speed": "avg_speed",
	"temp":  "avg_temp",
	"power": "avg_power",
}

func getDuration(start, end time.Time) time.Duration {
	return end.Sub(start)
}

var allowedMetrics = map[string]string{
	"speed":          "odometry_vehicle_speed",
	"temp":           "temperature_ambient",
	"power":          "electric_power_demand",
	"traction_force": "traction_traction_force",
	"brake_pressure": "traction_brake_pressure",
	// add others you want to expose
}

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

func parseQueryFilters(c *gin.Context) (*QueryFilters, bool) {
	vehicle := strings.TrimSpace(c.Query("vehicle_id"))
	startStr := strings.TrimSpace(c.Query("start"))
	endStr := strings.TrimSpace(c.Query("end"))

	if vehicle == "" {
		return nil, false
	}

	var start, end time.Time
	parsedStart, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return nil, false
	}
	start = parsedStart

	parsedEnd, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return nil, false
	}
	end = parsedEnd

	// Validate time range
	if start.After(end) {
		return nil, false
	}

	return &QueryFilters{
		VehicleID: vehicle,
		Start:     start,
		End:       end,
	}, true
}

func validateMetric(metric string) error {
	if metric == "" {
		return fmt.Errorf("metric cannot be empty")
	}
	if _, ok := allowedMetrics[metric]; !ok {
		return fmt.Errorf("invalid metric: %s", metric)
	}
	return nil
}
