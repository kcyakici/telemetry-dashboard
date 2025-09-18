package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetKPIs(c *gin.Context, pool *pgxpool.Pool) {
	filters, valid := parseQueryFilters(c)
	if !valid {
		slog.Warn("invalid KPI request params", "vehicle", c.Query("vehicle_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	slog.Info("handling KPI request",
		"vehicle", filters.VehicleID, "start", filters.Start, "end", filters.End)

	query := `
        SELECT
            AVG(odometry_vehicle_speed),     -- avg_speed
            MAX(temperature_ambient),        -- max_temp
            SUM(electric_power_demand),      -- total_power
            AVG(traction_brake_pressure),    -- avg_brake_pressure
            AVG(status_door_is_open)::float8 -- door_open_ratio
        FROM telemetry
        WHERE vehicle_id = $1
          AND time_iso >= $2::timestamptz
          AND time_iso <= $3::timestamptz
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var avgSpeed, maxTemp, totalPower, avgBrakePressure, doorOpenRatio *float64
	err := pool.QueryRow(ctx, query, filters.VehicleID, filters.Start, filters.End).
		Scan(&avgSpeed, &maxTemp, &totalPower, &avgBrakePressure, &doorOpenRatio)

	if err != nil {
		slog.Error("KPI query failed", "error", err, "vehicle", filters.VehicleID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}

	slog.Info("successfully retrieved KPIs",
		"vehicle", filters.VehicleID, "avg_speed", avgSpeed, "max_temp", maxTemp, "total_power", totalPower, "avg_brake_pressure", avgBrakePressure, "door_open_ratio", doorOpenRatio)
	c.JSON(http.StatusOK, KpiResponse{
		Avg_speed:          avgSpeed,
		Max_temp:           maxTemp,
		Total_power:        totalPower,
		Avg_brake_pressure: avgBrakePressure,
		Door_open_ratio:    doorOpenRatio,
	})
}

type KpiResponse struct {
	Avg_speed          *float64 `json:"avg_speed"`
	Max_temp           *float64 `json:"max_temp"`
	Total_power        *float64 `json:"total_power"`
	Avg_brake_pressure *float64 `json:"avg_brake_pressure"`
	Door_open_ratio    *float64 `json:"door_open_ratio"`
}
