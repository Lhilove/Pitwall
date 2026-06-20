package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lhilove/pitwall/internal/db"
	"github.com/lhilove/pitwall/internal/handlers"
	"github.com/lhilove/pitwall/internal/repository"
	"github.com/lhilove/pitwall/internal/service"
	"github.com/lhilove/pitwall/internal/websocket"
	"github.com/lhilove/pitwall/internal/workers"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	db.Connect()

	// Initialize WebSocket hub and start it
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize repositories and services
	repo := repository.NewTelemetryRepository(db.DB)
	svc := service.NewTelemetryService(repo)

	// Set up Gin router and handlers
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	r := gin.Default()

	fmt.Println("Server running on " + port)

	// Telemetry endpoints
	telemetryHandler := handlers.NewTelemetryHandler(svc) // Initialize handler with service

	// workers for background processing, pass hub to workers
	for i := 1; i <= 4; i++ {
		workers.StartTelemetryWorker(i, svc, hub)
	}

	r.POST("/telemetry", telemetryHandler.CreateTelemetry)
	r.POST("/simulate", telemetryHandler.Simulate)
	r.GET("/telemetry", telemetryHandler.GetTelemetry)
	r.GET("/telemetry/all", telemetryHandler.GetPaginated)
	r.GET("/telemetry/stats", telemetryHandler.GetStats)
	r.GET("/telemetry/leaderboard", telemetryHandler.GetLeaderboard)
	r.GET("/queue/stats", telemetryHandler.GetQueueStats)
	r.GET("/metrics", telemetryHandler.WorkerStats)
	r.GET("/telemetry/compare", telemetryHandler.CompareDrivers)
	r.GET("/ws", func(c *gin.Context) {
		handlers.ServeWs(hub, c)
	})

	r.Run(":" + port)

}
