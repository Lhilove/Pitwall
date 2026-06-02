package simulator

import (
	"math/rand"

	"github.com/lhilove/pitwall/internal/models"
	"github.com/lhilove/pitwall/internal/queue"
)

// GenerateTelemetry simulates telemetry data for a specific driver and enqueues it for processing
func GenerateTelemetry(driver string, count int) {

	for i := 0; i < count; i++ {

		queue.TelemetryQueue <- models.Telemetry{
			Driver:   driver,
			Lap:      rand.Intn(50) + 1,
			Speed:    rand.Intn(100) + 250,
			Throttle: rand.Intn(100),
			Brake:    rand.Intn(100),
			Gear:     rand.Intn(8) + 1,
		}
	}
}
