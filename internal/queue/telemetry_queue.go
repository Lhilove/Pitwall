package queue

import "github.com/lhilove/pitwall/internal/models"

// TelemetryQueue is a buffered channel that holds incoming telemetry data for processing by background workers
var TelemetryQueue = make(chan models.Telemetry, 100)

type QueueMetrics struct {
	QueueLength   int `json:"queue_length"`
	QueueCapacity int `json:"queue_capacity"`
}
