package models

// DriverLeaderboard represents a driver's performance for the leaderboard
type DriverLeaderboard struct {
	Driver       string  `json:"driver"`
	AverageSpeed float64 `json:"average_speed"`
	Records      int     `json:"records"`
}
