package migration

import (
	"database/sql"
	"fmt"
)

// CreateSensorTypeTables creates the necessary tables for sensor types
func CreateSensorTypeTables(db *sql.DB) error {
	// Create sensor_types table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sensor_types (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			manufacturer VARCHAR(255),
			model VARCHAR(255),
			version INTEGER NOT NULL DEFAULT 1,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating sensor_types table: %v", err)
	}

	return nil
}

// DropSensorTypeTables drops the sensor type tables
func DropSensorTypeTables(db *sql.DB) error {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS sensor_types CASCADE;
	`)
	if err != nil {
		return fmt.Errorf("error dropping sensor_types table: %v", err)
	}

	return nil
}
