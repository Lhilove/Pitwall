package models

type Telemetry struct {
	Driver   string `json:"driver" binding:"required"`
	Lap      int    `json:"lap"`
	Speed    int    `json:"speed"`
	Throttle int    `json:"throttle"`
	Brake    int    `json:"brake"`
	Gear     int    `json:"gear"`
}
