package main

import (
	"log"
	"log/slog"
	"os"

	"telemetry-dashboard/db"
	"telemetry-dashboard/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer conn.Close()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, //TODO change with frontend URL
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Host", "User-Agent", "Authorization", "Origin", "Accept", "Accept-Encoding", "Content-Length", "Content-Type", "Content type", "Connection"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// Replace the default logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("logger initialized", "level", "INFO", "format", "JSON")

	router.POST("/ingest-csv", func(c *gin.Context) { handlers.IngestCSV(c, conn) })
	router.GET("/live-trend", func(c *gin.Context) { handlers.LiveTrend(c, conn) })
	router.GET("/kpis", func(c *gin.Context) { handlers.GetKPIs(c, conn) })
	router.GET("/trend", func(c *gin.Context) { handlers.GetTrend(c, conn) })
	router.GET("/distribution", func(c *gin.Context) { handlers.GetDistribution(c, conn) })
	router.GET("/stats", func(c *gin.Context) { handlers.GetStats(c, conn) })

	log.Println("Server running at :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
