package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type TelemetryEvent struct { // TODO can be merged with telemetry struct
	VehicleID string   `json:"vehicle_id"`
	TimeISO   string   `json:"time_iso"`
	Speed     *float64 `json:"speed,omitempty"`
	Temp      *float64 `json:"temp,omitempty"`
	Power     *float64 `json:"power,omitempty"`
	Traction  *float64 `json:"traction,omitempty"`
	Brake     *float64 `json:"brake,omitempty"`
}

func sendWSError(conn *websocket.Conn, msg string) {
	errMsg := map[string]interface{}{
		"type":  "error",
		"error": msg,
	}
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_ = conn.WriteJSON(errMsg) // ignore error here, we're already in error handling
}

func LiveTrend(c *gin.Context, pool *pgxpool.Pool) {
	vehicle := c.Query("vehicle_id")
	metric := c.DefaultQuery("metric", "speed")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("ws upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	slog.Info("client connected", "vehicle", vehicle, "metric", metric)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := pool.Acquire(ctx)
	if err != nil {
		slog.Error("acquire conn failed", "error", err)
		sendWSError(conn, "database connection failed")
		return
	}
	defer db.Release()

	_, err = db.Exec(ctx, "LISTEN telemetry_channel;")
	if err != nil {
		slog.Error("LISTEN failed", "error", err)
		sendWSError(conn, "failed to subscribe to updates")
		return
	}

	for {
		notify, err := db.Conn().WaitForNotification(ctx)
		if err != nil {
			slog.Error("wait notify failed", "error", err)
			sendWSError(conn, "database listen error")
			return
		}

		var ev TelemetryEvent
		if err := json.Unmarshal([]byte(notify.Payload), &ev); err != nil {
			slog.Warn("unmarshal notify failed", "payload", notify.Payload, "error", err)
			sendWSError(conn, "invalid telemetry data received")
			continue
		}

		if vehicle != "" && ev.VehicleID != vehicle {
			continue
		}

		var value *float64
		switch metric {
		case "speed":
			value = ev.Speed
		case "temp":
			value = ev.Temp
		case "power":
			value = ev.Power
		case "traction":
			value = ev.Traction
		case "brake":
			value = ev.Brake
		default:
			slog.Warn("unsupported metric", "metric", metric)
			sendWSError(conn, "unsupported metric: "+metric)
			continue
		}

		if value == nil {
			continue
		}

		point := map[string]interface{}{
			"type":      "point",
			"timestamp": ev.TimeISO,
			"value":     *value,
		}

		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err := conn.WriteJSON(point); err != nil {
			slog.Info("client disconnected", "error", err)
			break
		}
	}
}
