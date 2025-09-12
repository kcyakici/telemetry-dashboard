package main

import (
	"log"

	"telemetry-dashboard/db"
	"telemetry-dashboard/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer conn.Close()

	router := gin.Default()

	router.POST("/ingest", func(c *gin.Context) { handlers.Ingest(c, conn) })
	// router.GET("/kpis", func(c *gin.Context) { handlers.GetKPIs(c, conn) })
	// router.GET("/trend", func(c *gin.Context) { handlers.GetTrend(c, conn) })
	// router.GET("/distribution", func(c *gin.Context) { handlers.GetDistribution(c, conn) })

	log.Println("Server running at :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
