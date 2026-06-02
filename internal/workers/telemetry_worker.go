package workers

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/lhilove/pitwall/internal/queue"
	"github.com/lhilove/pitwall/internal/service"
	"github.com/lhilove/pitwall/internal/websocket"
)

var ProcessedRecords atomic.Int64

// StartTelemetryWorkerWithHub starts a background worker that processes telemetry data and broadcasts it to WebSocket clients
func StartTelemetryWorker(id int, service *service.TelemetryService, hub *websocket.Hub) {
	// Start a goroutine to process telemetry data from the queue
	go func() {
		for telemetry := range queue.TelemetryQueue {
			ProcessedRecords.Add(1)
			err := service.CreateTelemetry(telemetry)
			if err != nil {
				fmt.Println("worker error:", err)
				continue // Skip broadcasting if save failed
			}

			// convert telemetry to JSON and send to all connected WebSocket clients
			if data, err := json.Marshal(telemetry); err == nil {
				hub.Broadcast <- data
			}
		}
	}()
}
