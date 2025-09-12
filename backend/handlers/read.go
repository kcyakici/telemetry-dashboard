package handlers

import (
	"context"
	"fmt"
	"net/http"
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
	// compute min/max for the column under the filters
	minMaxQuery := fmt.Sprintf(`
        SELECT MIN(%s), MAX(%s) FROM telemetry
        WHERE ($1 = '' OR vehicle_id = $1)
    `, col, col)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var min, max *float64
	if err := pool.QueryRow(ctx, minMaxQuery, vehicle).Scan(&min, &max); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "min/max query failed: " + err.Error()})
		return
	}
	if min == nil || max == nil || *min == *max {
		c.JSON(http.StatusOK, gin.H{"buckets": []interface{}{}})
		return
	}

	// create bucketed counts (10 buckets)
	q := fmt.Sprintf(`
    WITH bounds AS (
        SELECT MIN(%[1]s) AS minval, MAX(%[1]s) AS maxval
        FROM telemetry
        WHERE ($1 = '' OR vehicle_id = $1)
    )
    SELECT bucket, COUNT(*) AS cnt
    FROM (
        SELECT width_bucket(%[1]s, bounds.minval, bounds.maxval, 10) AS bucket
        FROM telemetry, bounds
        WHERE ($1 = '' OR vehicle_id = $1)
          AND %[1]s IS NOT NULL
    ) sub
    GROUP BY bucket
    ORDER BY bucket
`, col)

	rows, err := pool.Query(ctx, q, vehicle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bucket query failed: " + err.Error()})
		return
	}
	defer rows.Close()

	type Bucket struct {
		Bucket int `json:"bucket"`
		Count  int `json:"count"`
	}
	var out []Bucket
	for rows.Next() {
		var b Bucket
		if err := rows.Scan(&b.Bucket, &b.Count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed: " + err.Error()})
			return
		}
		out = append(out, b)
	}

	c.JSON(http.StatusOK, gin.H{"buckets": out})

}
