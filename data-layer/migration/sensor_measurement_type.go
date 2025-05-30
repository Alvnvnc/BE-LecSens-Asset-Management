package migration

import (
	"database/sql"
	"fmt"
)

// CreateSensorMeasurementTypeTables creates the necessary tables for sensor measurement types
func CreateSensorMeasurementTypeTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sensor_measurement_types (
			id UUID PRIMARY KEY,
			sensor_type_id UUID NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			properties_schema JSONB,
			ui_configuration JSONB,
			version VARCHAR(50) NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP,
			CONSTRAINT fk_sensor_type FOREIGN KEY (sensor_type_id) 
				REFERENCES sensor_types(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating sensor_measurement_types table: %v", err)
	}

	// Add version constraint
	_, err = db.Exec(`
		DO $$ 
		BEGIN
			-- Add version constraint if it doesn't exist
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'check_version_positive' 
				AND conrelid = 'sensor_measurement_types'::regclass
			) THEN
				ALTER TABLE sensor_measurement_types
					ADD CONSTRAINT check_version_positive 
					CHECK (version ~ '^[0-9]+\.[0-9]+\.[0-9]+$');
			END IF;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("error adding constraints: %v", err)
	}

	return nil
}

// DropSensorMeasurementTypeTables drops the sensor measurement type tables
func DropSensorMeasurementTypeTables(db *sql.DB) error {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS sensor_measurement_types CASCADE;
	`)
	if err != nil {
		return fmt.Errorf("error dropping sensor_measurement_types table: %v", err)
	}

	return nil
}
