package repository

import (
	"database/sql"
	"fmt"

	"github.com/lhilove/pitwall/internal/models"
)

// TelemetryRepository handles database operations for telemetry data
type TelemetryRepository struct {
	DB *sql.DB
}

// NewTelemetryRepository creates a new instance of TelemetryRepository
func NewTelemetryRepository(db *sql.DB) *TelemetryRepository {
	return &TelemetryRepository{
		DB: db,
	}
}

// Create inserts a new telemetry record into the database
func (r *TelemetryRepository) Create(t models.Telemetry) error {
	query := `
		INSERT INTO telemetry
		(driver, lap, speed, throttle, brake, gear)
		VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := r.DB.Exec(
		query,
		t.Driver,
		t.Lap,
		t.Speed,
		t.Throttle,
		t.Brake,
		t.Gear,
	)

	return err
}

// GetAll retrieves all telemetry records from the database
func (r *TelemetryRepository) GetAll() ([]models.Telemetry, error) {

	query := `
		SELECT
			driver,
			lap,
			speed,
			throttle,
			brake,
			gear
		FROM telemetry
	`

	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var telemetryList []models.Telemetry

	for rows.Next() {

		var t models.Telemetry

		err := rows.Scan(
			&t.Driver,
			&t.Lap,
			&t.Speed,
			&t.Throttle,
			&t.Brake,
			&t.Gear,
		)

		if err != nil {
			return nil, err
		}

		telemetryList = append(telemetryList, t)
	}

	return telemetryList, nil
}

// GetByDriver retrieves telemetry records for a specific driver
func (r *TelemetryRepository) GetByDriver(driver string) ([]models.Telemetry, error) {

	query := `
		SELECT
			driver,
			lap,
			speed,
			throttle,
			brake,
			gear
		FROM telemetry
		WHERE driver = $1
	`

	rows, err := r.DB.Query(
		query,
		driver,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var telemetry []models.Telemetry

	for rows.Next() {

		var t models.Telemetry

		err := rows.Scan(
			&t.Driver,
			&t.Lap,
			&t.Speed,
			&t.Throttle,
			&t.Brake,
			&t.Gear,
		)

		if err != nil {
			return nil, err
		}

		telemetry = append(telemetry, t)
	}

	return telemetry, nil
}

// GetStats retrieves aggregated statistics about the telemetry data
func (r *TelemetryRepository) GetStats() (models.TelemetryStats, error) {

	query := `
		SELECT
			AVG(speed),
			MAX(speed),
			COUNT(*)
		FROM telemetry
	`

	var stats models.TelemetryStats

	err := r.DB.QueryRow(query).Scan(
		&stats.AverageSpeed,
		&stats.MaxSpeed,
		&stats.TotalRecords,
	)

	return stats, err
}

// GetLeaderboard retrieves a leaderboard of drivers based on their average speed
func (r *TelemetryRepository) GetLeaderboard() ([]models.DriverLeaderboard, error) {

	query := `
		SELECT
			driver,
			AVG(speed) as average_speed,
			COUNT(*) as records
		FROM telemetry
		GROUP BY driver
		ORDER BY average_speed DESC
	`

	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var leaderboard []models.DriverLeaderboard

	for rows.Next() {

		var driver models.DriverLeaderboard

		err := rows.Scan(
			&driver.Driver,
			&driver.AverageSpeed,
			&driver.Records,
		)

		if err != nil {
			return nil, err
		}

		leaderboard = append(leaderboard, driver)
	}

	return leaderboard, nil
}

// GetPaginated retrieves a paginated list of telemetry records
func (r *TelemetryRepository) GetPaginated(limit int, offset int, sort string, order string) ([]models.Telemetry, error) {

	query := fmt.Sprintf(`
		SELECT
			driver,
			lap,
			speed,
			throttle,
			brake,
			gear
		FROM telemetry
		ORDER BY %s %s
		LIMIT $1
		OFFSET $2
	`, sort, order)

	rows, err := r.DB.Query(
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var telemetry []models.Telemetry

	for rows.Next() {

		var t models.Telemetry

		err := rows.Scan(
			&t.Driver,
			&t.Lap,
			&t.Speed,
			&t.Throttle,
			&t.Brake,
			&t.Gear,
		)

		if err != nil {
			return nil, err
		}

		telemetry = append(telemetry, t)
	}

	return telemetry, nil
}
