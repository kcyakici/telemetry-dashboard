package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetTrend(c *gin.Context, pool *pgxpool.Pool) {
	filters, valid := parseQueryFilters(c)
	if !valid {
		slog.Warn("invalid trend request params", "vehicle", c.Query("vehicle_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	metric := c.DefaultQuery("metric", "speed")
	if err := validateMetric(metric); err != nil {
		slog.Warn("invalid trend params", "metric", metric, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "metric is not valid"})
		return
	}

	slog.Info("handling trend request",
		"metric", metric, "vehicle", filters.VehicleID, "start", filters.Start, "end", filters.End)

	queryStr := buildTrendQuery(metric, filters)
	slog.Debug("constructed trend query", "sql", queryStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := pool.Query(ctx, queryStr, filters.VehicleID, filters.Start, filters.End)
	if err != nil {
		slog.Error("trend query failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
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
			slog.Error("row scan failed inside trend", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed"})
			return
		}
		if v != nil {
			result = append(result, Point{Timestamp: ts.Format(time.RFC3339), Value: *v})
		}
	}

	slog.Info("trend query returned rows", "metric", metric, "count", len(result))
	c.JSON(http.StatusOK, result)
}

func buildTrendQuery(metric string, filters *QueryFilters) string {
	col := allowedMetrics[metric]
	duration := getDuration(filters.Start, filters.End)

	var baseQuery, timeCol string
	if duration > 1*time.Hour {
		// Use aggregated tables for better performance
		if table, exists := aggregatedTables[metric]; exists {
			slog.Debug("long time interval selected, querying aggregated table", "interval", duration)
			aggCol := aggregatedColumns[metric]
			baseQuery = fmt.Sprintf(`SELECT bucket AS time_iso, %s AS value FROM %s`, aggCol, table)
			timeCol = "bucket"
		} else {
			// Fallback to raw telemetry
			baseQuery = fmt.Sprintf(`SELECT time_iso, %s AS value FROM telemetry`, col)
			timeCol = "time_iso"
		}
	} else {
		// Raw telemetry for short time ranges
		baseQuery = fmt.Sprintf(`SELECT time_iso, %s AS value FROM telemetry`, col)
		timeCol = "time_iso"
	}

	query := fmt.Sprintf(`
		%s
		WHERE vehicle_id = $1
		  AND %s >= $2::timestamptz
		  AND %s <= $3::timestamptz
		ORDER BY %s
	`, baseQuery, timeCol, timeCol, timeCol)

	return query
}
