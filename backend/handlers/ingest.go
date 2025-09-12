package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Telemetry struct {
	VehicleID   string    `json:"vehicle_id"`
	Timestamp   time.Time `json:"timestamp"`
	Speed       float64   `json:"speed"`
	Temperature float64   `json:"temperature"`
	Energy      float64   `json:"energy"`
}

func Ingest(c *gin.Context, pool *pgxpool.Pool) {
	var data []Telemetry

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	batch := &pgx.Batch{}
	for _, t := range data {
		batch.Queue(
			`INSERT INTO telemetry (vehicle_id, timestamp, speed, temperature, energy)
			 VALUES ($1, $2, $3, $4, $5)`,
			t.VehicleID, t.Timestamp, t.Speed, t.Temperature, t.Energy,
		)
	}

	br := pool.SendBatch(context.Background(), batch)
	if err := br.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "inserted": len(data)})
}
