package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Distribution: compute min/max then bucket using width_bucket
func GetDistribution(c *gin.Context, pool *pgxpool.Pool) {
	filters, valid := parseQueryFilters(c)
	if !valid {
		slog.Warn("invalid distribution request params", "vehicle", c.Query("vehicle_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	metric := c.DefaultQuery("metric", "speed")
	if err := validateMetric(metric); err != nil {
		slog.Warn("invalid distribution params", "metric", metric, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "metric is not valid"})
		return
	}

	slog.Info("handling distribution request",
		"metric", metric, "vehicle", filters.VehicleID, "start", filters.Start, "end", filters.End)

	binsStr := c.DefaultQuery("bins", "10")
	bins, err := strconv.Atoi(binsStr)
	if err != nil || bins <= 5 || bins > 20 {
		slog.Warn("invalid bins param, falling back to default", "binsStr", binsStr, "error", err)
		bins = 10
	}

	col := allowedMetrics[metric]

	// Compute min/max with the same pattern as KPIs/Trend
	minMaxQuery := fmt.Sprintf(`
		SELECT MIN(%s), MAX(%s)
		FROM telemetry
		WHERE vehicle_id = $1
		  AND time_iso >= $2::timestamptz
		  AND time_iso <= $3::timestamptz
	`, col, col)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var min, max *float64
	err = pool.QueryRow(ctx, minMaxQuery, filters.VehicleID, filters.Start, filters.End).Scan(&min, &max)

	if err != nil {
		slog.Error("min max query failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "min/max query failed"})
		return
	}

	if min == nil || max == nil || *min == *max {
		slog.Info("distribution query returned no range",
			"metric", metric, "vehicle", filters.VehicleID)
		c.JSON(http.StatusOK, DistributionResponse{
			Metric:  metric,
			Vehicle: filters.VehicleID,
			Bins:    bins,
			Min:     nil,
			Max:     nil,
			From:    filters.Start,
			To:      filters.End,
			Buckets: []Bucket{},
		})
		return
	}

	// Bucket query
	bucketQuery := fmt.Sprintf(`
		WITH bounds AS (
			SELECT MIN(%[1]s) AS minval, MAX(%[1]s) AS maxval
			FROM telemetry
			WHERE vehicle_id = $1
			  AND time_iso >= $2::timestamptz
			  AND time_iso <= $3::timestamptz
		)
		SELECT bucket, COUNT(*) AS cnt, bounds.minval, bounds.maxval
		FROM (
			SELECT width_bucket(%[1]s, bounds.minval, bounds.maxval, $4) AS bucket
			FROM telemetry, bounds
			WHERE vehicle_id = $1
			  AND time_iso >= $2::timestamptz
			  AND time_iso <= $3::timestamptz
			  AND %[1]s IS NOT NULL
		) sub, bounds
		GROUP BY bucket, bounds.minval, bounds.maxval
		ORDER BY bucket
	`, col)

	rows, err := pool.Query(ctx, bucketQuery, filters.VehicleID, filters.Start, filters.End, bins)
	if err != nil {
		slog.Error("bucket query failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bucket query failed"})
		return
	}
	defer rows.Close()

	var out []Bucket
	bucketWidth := (*max - *min) / float64(bins)

	for rows.Next() {
		var b Bucket
		var minVal, maxVal float64
		if err := rows.Scan(&b.Bucket, &b.Count, &minVal, &maxVal); err != nil {
			slog.Error("row scan failed inside distribution", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed"})
			return
		}

		b.RangeMin = minVal + float64(b.Bucket-1)*bucketWidth
		b.RangeMax = minVal + float64(b.Bucket)*bucketWidth
		out = append(out, b)
	}

	slog.Info("distribution computed successfully",
		"metric", metric, "vehicle", filters.VehicleID, "bins", bins, "bucket_count", len(out))
	c.JSON(http.StatusOK, DistributionResponse{
		Metric:  metric,
		Vehicle: filters.VehicleID,
		Bins:    bins,
		Min:     min,
		Max:     max,
		From:    filters.Start,
		To:      filters.End,
		Buckets: out,
	})
}

type DistributionResponse struct {
	Metric  string    `json:"metric"`
	Vehicle string    `json:"vehicle"`
	Bins    int       `json:"bins"`
	Min     *float64  `json:"min"`
	Max     *float64  `json:"max"`
	From    time.Time `json:"from,omitempty"`
	To      time.Time `json:"to,omitempty"`
	Buckets []Bucket  `json:"buckets"`
}

type Bucket struct {
	Bucket   int     `json:"bucket"`
	Count    int     `json:"count"`
	RangeMin float64 `json:"range_min"`
	RangeMax float64 `json:"range_max"`
}
