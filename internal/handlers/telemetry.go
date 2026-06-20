package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
	simulator "github.com/lhilove/pitwall/Simulator"
	"github.com/lhilove/pitwall/internal/models"
	"github.com/lhilove/pitwall/internal/queue"
	"github.com/lhilove/pitwall/internal/service"
	"github.com/lhilove/pitwall/internal/websocket"
	"github.com/lhilove/pitwall/internal/workers"
)

// TelemetryHandler handles HTTP requests related to telemetry data
type TelemetryHandler struct {
	service *service.TelemetryService
}

// NewTelemetryHandler creates a new instance of TelemetryHandler
func NewTelemetryHandler(service *service.TelemetryService) *TelemetryHandler {

	return &TelemetryHandler{
		service: service,
	}
}

// CreateTelemetry handles incoming telemetry data and saves it to the database
func (h *TelemetryHandler) CreateTelemetry(c *gin.Context) {

	var telemetry models.Telemetry

	if err := c.ShouldBindJSON(&telemetry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate basic sanity, this matters in real systems
	if telemetry.Speed < 0 || telemetry.Speed > 400 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid speed value"})
		return
	}

	// Enqueue telemetry data for background processing
	queue.TelemetryQueue <- telemetry

	c.JSON(http.StatusAccepted, gin.H{"message": "queued"})

}

// GetPaginated retrieves telemetry data with pagination support
func (h *TelemetryHandler) GetPaginated(c *gin.Context) {

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Validate and sanitize sort and order parameters
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "desc")

	allowedSorts := map[string]bool{
		"id":     true,
		"speed":  true,
		"lap":    true,
		"driver": true,
		"gear":   true,
	}

	if !allowedSorts[sort] {
		sort = "id"
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Convert query parameters to integers
	page, err := strconv.Atoi(pageStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page value"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit value"})
		return
	}

	offset := (page - 1) * limit

	data, err := h.service.GetPaginated(limit, offset, sort, order)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"page": page, "limit": limit, "offset": offset, "data": data})

}

// GetTelemetry retrieves telemetry data, optionally filtered by driver
func (h *TelemetryHandler) GetTelemetry(c *gin.Context) {

	driver := c.Query("driver")

	if driver != "" {

		data, err := h.service.GetByDriver(driver)

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}

		c.JSON(http.StatusOK, data)

		return
	}

	data, err := h.service.GetTelemetry()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, data)
}

// GetStats retrieves telemetry statistics such as average speed, max speed, and total records
func (h *TelemetryHandler) GetStats(c *gin.Context) {

	stats, err := h.service.GetStats()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, stats)
}

// Leaderboard retrieves a leaderboard of drivers based on their average speed
func (h *TelemetryHandler) GetLeaderboard(c *gin.Context) {

	data, err := h.service.GetLeaderboard()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, data)
}

// QueueStats retrieves the current length and capacity of the telemetry queue
func (h *TelemetryHandler) GetQueueStats(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"queue_length": len(queue.TelemetryQueue), "queue_capacity": cap(queue.TelemetryQueue)})
}

// Simulate starts a background simulation of telemetry data for a specific driver
func (h *TelemetryHandler) Simulate(c *gin.Context) {

	go simulator.GenerateTelemetry(
		"Verstappen",
		500,
	)

	c.JSON(http.StatusAccepted, gin.H{"message": "simulation started"})
}

// WorkerStats retrieves the total number of processed telemetry records by the background workers
func (h *TelemetryHandler) WorkerStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"processed_records": workers.ProcessedRecords.Load()})
}

// ServeWs upgrades the HTTP connection to a WebSocket and registers the client with the WebSocket hub
func (h *TelemetryHandler) ServeWs(c *gin.Context) {
	var upgrader = gorilla.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	hub := websocket.NewHub()
	go hub.Run()
	client := &websocket.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	hub.Register <- client // Register the client with the hub
	go client.WritePump()
}

// this is the http handler for comparing two drivers based on their telemetry statistics
func (h *TelemetryHandler) CompareDrivers(c *gin.Context) {
	driverA := c.Query("b")
	driverB := c.Query("a")

	if driverA == "" || driverB == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both driver parameters are required"})
		return
	}

	statsA, statsB, err := h.service.CompareDrivers(driverA, driverB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"driverA": statsA, "driverB": statsB})
}
