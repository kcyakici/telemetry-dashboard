package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// KPIs: avg speed, max temperature, total power (using dataset columns)
func GetKPIs(c *gin.Context, pool *pgxpool.Pool) {
	// optional filters
	vehicle := c.Query("vehicle_id") // may be empty
	start := c.Query("start")
	end := c.Query("end")

	// we map KPIs to specific columns
	avgCol := "odometry_vehicle_speed"
	maxCol := "temperature_ambient"
	sumCol := "electric_power_demand"

	query := fmt.Sprintf(`
        SELECT AVG(%s), MAX(%s), SUM(%s)
        FROM telemetry
        WHERE ($1 = '' OR vehicle_id = $1)
          AND ($2 = '' OR time_iso >= $2::timestamptz)
          AND ($3 = '' OR time_iso <= $3::timestamptz)
    `, avgCol, maxCol, sumCol)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var avg, mx, sum *float64
	if err := pool.QueryRow(ctx, query, vehicle, start, end).Scan(&avg, &mx, &sum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"avg_speed":   avg,
		"max_temp":    mx,
		"total_power": sum,
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
func GetDistribution(c *gin.Context, pool *pgxpool.Pool) {
	metric := c.DefaultQuery("metric", "temp")
	col, ok := allowedMetrics[metric]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid metric"})
		return
	}

	vehicle := c.Query("vehicle_id")
	binsStr := c.DefaultQuery("bins", "10")
	bins, err := strconv.Atoi(binsStr)
	if err != nil || bins <= 0 {
		bins = 10
	}

	// Parse optional time filters
	fromStr := c.Query("from")
	toStr := c.Query("to")
	var fromTime, toTime *time.Time
	if fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			fromTime = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from timestamp"})
			return
		}
	}
	if toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			toTime = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to timestamp"})
			return
		}
	}

	// Dynamic WHERE conditions
	conditions := []string{"($1 = '' OR vehicle_id = $1)"}
	args := []interface{}{vehicle}
	argIdx := 2
	if fromTime != nil {
		conditions = append(conditions, fmt.Sprintf("time_iso >= $%d", argIdx))
		args = append(args, *fromTime)
		argIdx++
	}
	if toTime != nil {
		conditions = append(conditions, fmt.Sprintf("time_iso <= $%d", argIdx))
		args = append(args, *toTime)
		argIdx++
	}
	whereClause := strings.Join(conditions, " AND ")

	// Compute min/max
	minMaxQuery := fmt.Sprintf(`
        SELECT MIN(%s), MAX(%s) FROM telemetry
        WHERE %s
    `, col, col, whereClause)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var min, max *float64
	if err := pool.QueryRow(ctx, minMaxQuery, args...).Scan(&min, &max); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "min/max query failed: " + err.Error()})
		return
	}
	if min == nil || max == nil || *min == *max {
		c.JSON(http.StatusOK, gin.H{"buckets": []interface{}{}})
		return
	}

	// Bucket query
	q := fmt.Sprintf(`
    WITH bounds AS (
        SELECT MIN(%[1]s) AS minval, MAX(%[1]s) AS maxval
        FROM telemetry
        WHERE %[2]s
    )
    SELECT bucket, COUNT(*) AS cnt, bounds.minval, bounds.maxval
    FROM (
        SELECT width_bucket(%[1]s, bounds.minval, bounds.maxval, %[3]d) AS bucket
        FROM telemetry, bounds
        WHERE %[2]s
          AND %[1]s IS NOT NULL
    ) sub, bounds
    GROUP BY bucket, bounds.minval, bounds.maxval
    ORDER BY bucket
`, col, whereClause, bins)

	rows, err := pool.Query(ctx, q, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bucket query failed: " + err.Error()})
		return
	}
	defer rows.Close()

	type Bucket struct {
		Bucket   int     `json:"bucket"`
		Count    int     `json:"count"`
		RangeMin float64 `json:"range_min"`
		RangeMax float64 `json:"range_max"`
	}
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

	c.JSON(http.StatusOK, gin.H{
		"metric":  metric,
		"vehicle": vehicle,
		"bins":    bins,
		"min":     *min,
		"max":     *max,
		"from":    fromTime,
		"to":      toTime,
		"buckets": out,
	})
}
