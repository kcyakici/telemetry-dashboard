package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetKPIs(c *gin.Context, pool *pgxpool.Pool) {
	vehicleID := c.Query("vehicle_id")
	start := c.Query("start")
	end := c.Query("end")

	query := `
		SELECT 
			AVG(speed) as avg_speed,
			MAX(temperature) as max_temp,
			SUM(energy) as total_energy
		FROM telemetry
		WHERE ($1 = '' OR vehicle_id = $1)
		  AND ($2 = '' OR timestamp >= $2::timestamptz)
		  AND ($3 = '' OR timestamp <= $3::timestamptz)
	`

	row := pool.QueryRow(context.Background(), query, vehicleID, start, end)

	var avgSpeed, maxTemp, totalEnergy *float64
	if err := row.Scan(&avgSpeed, &maxTemp, &totalEnergy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"avg_speed":    avgSpeed,
		"max_temp":     maxTemp,
		"total_energy": totalEnergy,
	})
}
