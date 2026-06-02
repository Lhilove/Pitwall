package service

import (
	"github.com/lhilove/pitwall/internal/models"
	"github.com/lhilove/pitwall/internal/repository"
)

// TelemetryService provides business logic for telemetry operations
type TelemetryService struct {
	repo *repository.TelemetryRepository
}

// NewTelemetryService creates a new instance of TelemetryService
func NewTelemetryService(repo *repository.TelemetryRepository) *TelemetryService {

	return &TelemetryService{
		repo: repo,
	}
}

// CreateTelemetry processes incoming telemetry data and saves it to the database
func (s *TelemetryService) CreateTelemetry(t models.Telemetry) error {

	return s.repo.Create(t)
}

// GetTelemetry retrieves all telemetry data from the database
func (s *TelemetryService) GetTelemetry() ([]models.Telemetry, error) {
	return s.repo.GetAll()
}

// GetByDriver retrieves telemetry data for a specific driver
func (s *TelemetryService) GetByDriver(driver string) ([]models.Telemetry, error) {

	return s.repo.GetByDriver(driver)
}

// GetStats calculates and returns telemetry statistics
func (s *TelemetryService) GetStats() (models.TelemetryStats, error) {
	return s.repo.GetStats()
}

// GetLeaderboard retrieves the driver leaderboard based on average speed
func (s *TelemetryService) GetLeaderboard() ([]models.DriverLeaderboard, error) {
	return s.repo.GetLeaderboard()
}

// GetPaginated retrieves telemetry data with pagination support
func (s *TelemetryService) GetPaginated(limit int, offset int, sort string, order string) ([]models.Telemetry, error) {

	return s.repo.GetPaginated(limit, offset, sort, order)
}

func (s *TelemetryService) ServeWs(hub string) {
	// This will be implemented in the handler to avoid circular dependencies
}
