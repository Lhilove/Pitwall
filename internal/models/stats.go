package models

// TelemetryStats represents aggregated statistics for telemetry data
type TelemetryStats struct {
	AverageSpeed float64 `json:"average_speed"`
	MaxSpeed     int     `json:"max_speed"`
	TotalRecords int     `json:"total_records"`
}
