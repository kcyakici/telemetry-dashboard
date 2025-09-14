package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var allowedMetrics = map[string]string{
	"speed":    "odometry_vehicle_speed",
	"temp":     "temperature_ambient",
	"power":    "electric_power_demand",
	"traction": "traction_traction_force",
	// add others you want to expose
}

func GetKPIs(c *gin.Context, pool *pgxpool.Pool) {
	vehicle := c.Query("vehicle_id")
	start := c.Query("start")
	end := c.Query("end")

	query := `
        SELECT
            AVG(odometry_vehicle_speed),     -- avg_speed
            MAX(temperature_ambient),        -- max_temp
            SUM(electric_power_demand),      -- total_power
            AVG(traction_brake_pressure),    -- avg_brake_pressure
            AVG(status_door_is_open)::float8 -- door_open_ratio
        FROM telemetry
        WHERE ($1 = '' OR vehicle_id = $1)
          AND ($2 = '' OR time_iso >= $2::timestamptz)
          AND ($3 = '' OR time_iso <= $3::timestamptz)
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var avgSpeed, maxTemp, totalPower, avgBrakePressure, doorOpenRatio *float64
	if err := pool.QueryRow(ctx, query, vehicle, start, end).
		Scan(&avgSpeed, &maxTemp, &totalPower, &avgBrakePressure, &doorOpenRatio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"avg_speed":          avgSpeed,
		"max_temp":           maxTemp,
		"total_power":        totalPower,
		"avg_brake_pressure": avgBrakePressure,
		"door_open_ratio":    doorOpenRatio,
	})
}

// Trend: returns time series for a single validated metric
func GetTrend(c *gin.Context, pool *pgxpool.Pool) {
	metric := c.DefaultQuery("metric", "speed")
	col, ok := allowedMetrics[metric]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid metric"})
		return
	}

	vehicle := c.Query("vehicle_id")
	start := c.Query("start")
	end := c.Query("end")

	query := fmt.Sprintf(`
        SELECT time_iso, %s
        FROM telemetry
        WHERE ($1 = '' OR vehicle_id = $1)
          AND ($2 = '' OR time_iso >= $2::timestamptz)
          AND ($3 = '' OR time_iso <= $3::timestamptz)
        ORDER BY time_iso
    `, col)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := pool.Query(ctx, query, vehicle, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed: " + err.Error()})
		return
	}
	defer rows.Close()

	type Point struct {
		Timestamp string  `json:"timestamp"`
		Value     float64 `json:"value"`
	}
	var result []Point
	for rows.Next() {
		var ts time.Time
		var v *float64
		if err := rows.Scan(&ts, &v); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed: " + err.Error()})
			return
		}
		if v != nil {
			result = append(result, Point{Timestamp: ts.Format(time.RFC3339), Value: *v})
		}
	}

	c.JSON(http.StatusOK, result)
}

// Distribution: compute min/max then bucket using width_bucket
// Distribution: compute min/max then bucket using width_bucket
func GetDistribution(c *gin.Context, pool *pgxpool.Pool) {
	metric := c.DefaultQuery("metric", "temp")
	col, ok := allowedMetrics[metric]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid metric"})
		return
	}

	vehicle := c.Query("vehicle_id")
	from := c.Query("from")
	to := c.Query("to")
	binsStr := c.DefaultQuery("bins", "10")
	bins, err := strconv.Atoi(binsStr)
	if err != nil || bins <= 0 {
		bins = 10
	}

	// Compute min/max with the same pattern as KPIs/Trend
	minMaxQuery := fmt.Sprintf(`
		SELECT MIN(%s), MAX(%s)
		FROM telemetry
		WHERE ($1 = '' OR vehicle_id = $1)
		  AND ($2 = '' OR time_iso >= $2::timestamptz)
		  AND ($3 = '' OR time_iso <= $3::timestamptz)
	`, col, col)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var min, max *float64
	if err := pool.QueryRow(ctx, minMaxQuery, vehicle, from, to).Scan(&min, &max); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "min/max query failed: " + err.Error()})
		return
	}

	if min == nil || max == nil || *min == *max {
		c.JSON(http.StatusOK, DistributionResponse{
			Metric:  metric,
			Vehicle: vehicle,
			Bins:    bins,
			Min:     nil,
			Max:     nil,
			From:    from,
			To:      to,
			Buckets: []Bucket{},
		})
		return
	}

	// Bucket query
	q := fmt.Sprintf(`
		WITH bounds AS (
			SELECT MIN(%[1]s) AS minval, MAX(%[1]s) AS maxval
			FROM telemetry
			WHERE ($1 = '' OR vehicle_id = $1)
			  AND ($2 = '' OR time_iso >= $2::timestamptz)
			  AND ($3 = '' OR time_iso <= $3::timestamptz)
		)
		SELECT bucket, COUNT(*) AS cnt, bounds.minval, bounds.maxval
		FROM (
			SELECT width_bucket(%[1]s, bounds.minval, bounds.maxval, %[2]d) AS bucket
			FROM telemetry, bounds
			WHERE ($1 = '' OR vehicle_id = $1)
			  AND ($2 = '' OR time_iso >= $2::timestamptz)
			  AND ($3 = '' OR time_iso <= $3::timestamptz)
			  AND %[1]s IS NOT NULL
		) sub, bounds
		GROUP BY bucket, bounds.minval, bounds.maxval
		ORDER BY bucket
	`, col, bins)

	rows, err := pool.Query(ctx, q, vehicle, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bucket query failed: " + err.Error()})
		return
	}
	defer rows.Close()

	var out []Bucket
	bucketWidth := (*max - *min) / float64(bins)

	for rows.Next() {
		var b Bucket
		var minVal, maxVal float64
		if err := rows.Scan(&b.Bucket, &b.Count, &minVal, &maxVal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed: " + err.Error()})
			return
		}

		b.RangeMin = minVal + float64(b.Bucket-1)*bucketWidth
		b.RangeMax = minVal + float64(b.Bucket)*bucketWidth
		out = append(out, b)
	}

	c.JSON(http.StatusOK, DistributionResponse{
		Metric:  metric,
		Vehicle: vehicle,
		Bins:    bins,
		Min:     min,
		Max:     max,
		From:    from,
		To:      to,
		Buckets: out,
	})
}

func filterParams(c *gin.Context) (*QueryFilters, bool) {
	vehicle := c.Query("vehicle_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	// enforce both or none
	if (startStr == "" && endStr != "") || (startStr != "" && endStr == "") {
		return nil, false
	}

	var start, end *time.Time
	if startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return nil, false
		}
		start = &t
	}
	if endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return nil, false
		}
		end = &t
	}

	return &QueryFilters{
		VehicleID: vehicle,
		Start:     start,
		End:       end,
	}, true
}

type QueryFilters struct {
	VehicleID string
	Start     *time.Time
	End       *time.Time
}

type DistributionResponse struct {
	Metric  string   `json:"metric"`
	Vehicle string   `json:"vehicle"`
	Bins    int      `json:"bins"`
	Min     *float64 `json:"min"`
	Max     *float64 `json:"max"`
	From    string   `json:"from,omitempty"`
	To      string   `json:"to,omitempty"`
	Buckets []Bucket `json:"buckets"`
}

type Bucket struct {
	Bucket   int     `json:"bucket"`
	Count    int     `json:"count"`
	RangeMin float64 `json:"range_min"`
	RangeMax float64 `json:"range_max"`
}
