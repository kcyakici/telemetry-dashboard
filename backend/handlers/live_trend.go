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
	CheckOrigin: func(r *http.Request) bool { return true }, // TODO restrict to frontend URL
	// Add buffer sizes to prevent blocking
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type TelemetryEvent struct {
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
	_ = conn.WriteJSON(errMsg)
}

func LiveTrend(c *gin.Context, pool *pgxpool.Pool) {
	vehicle := c.Query("vehicle_id")
	metric := c.DefaultQuery("metric", "speed")
	// TODO validate metric before upgrading to WebSocket

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("ws upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	slog.Info("client connected", "vehicle", vehicle, "metric", metric)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- Keepalive setup ---
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(appData string) error {
		slog.Debug("pong received", "data", appData)
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Ping goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("ping goroutine panic", "error", r)
				cancel()
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					slog.Debug("ping failed", "metric", metric, "error", err)
					cancel()
					return
				}
				slog.Debug("ping sent", "metric", metric)
			}
		}
	}()

	// Reader goroutine (detect disconnects, enable pong handler)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("reader goroutine panic", "error", r)
				cancel()
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if _, _, err := conn.ReadMessage(); err != nil {
					slog.Debug("client disconnected", "error", err)
					cancel()
					return
				}
			}
		}
	}()

	// DB LISTEN
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

	// Send connection success message
	successMsg := map[string]interface{}{
		"type":    "connected",
		"vehicle": vehicle,
		"metric":  metric,
	}
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err := conn.WriteJSON(successMsg); err != nil {
		slog.Error("failed to send connection success", "error", err)
		return
	}

	// Main loop: forward DB notifications → client
	for {
		select {
		case <-ctx.Done():
			slog.Info("closing LiveTrend handler (context cancelled)")
			return
		default:
			notify, err := db.Conn().WaitForNotification(ctx)
			if err != nil {
				if ctx.Err() != nil {
					slog.Debug("db listener stopped", "reason", ctx.Err())
					return
				}
				slog.Error("wait notify failed", "error", err)
				sendWSError(conn, "database listen error")
				return
			}

			var ev TelemetryEvent
			if err := json.Unmarshal([]byte(notify.Payload), &ev); err != nil {
				slog.Warn("unmarshal notify failed", "payload", notify.Payload, "error", err)
				continue // Don't send error to client for malformed data
			}

			// Filter by vehicle
			if vehicle != "" && ev.VehicleID != vehicle {
				continue
			}

			// Pick metric value
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
				slog.Debug("client write failed, closing connection", "error", err)
				return
			}
		}
	}
}
