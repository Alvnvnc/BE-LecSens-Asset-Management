package migration

import (
	"database/sql"
	"fmt"
)

// CreateSensorMeasurementFieldTables creates the necessary tables for sensor measurement fields
func CreateSensorMeasurementFieldTables(db *sql.DB) error {
	// Create sensor_measurement_fields table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sensor_measurement_fields (
			id UUID PRIMARY KEY,
			sensor_measurement_type_id UUID NOT NULL,
			name VARCHAR(255) NOT NULL,
			label VARCHAR(255) NOT NULL,
			description TEXT,
			data_type VARCHAR(50) NOT NULL,
			required BOOLEAN NOT NULL DEFAULT false,
			unit VARCHAR(50),
			min FLOAT,
			max FLOAT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP,
			CONSTRAINT fk_sensor_type FOREIGN KEY (sensor_measurement_type_id) 
				REFERENCES sensor_measurement_types(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating sensor_measurement_fields table: %v", err)
	}

	// Add constraints using PL/pgSQL block
	_, err = db.Exec(`
		DO $$ 
		BEGIN
			-- Add min/max constraint if it doesn't exist
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'check_min_max' 
				AND conrelid = 'sensor_measurement_fields'::regclass
			) THEN
				ALTER TABLE sensor_measurement_fields
					ADD CONSTRAINT check_min_max 
					CHECK ((min IS NULL AND max IS NULL) 
						OR (min IS NOT NULL AND max IS NOT NULL 
							AND min < max));
			END IF;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("error adding constraints: %v", err)
	}

	return nil
}

// DropSensorMeasurementFieldTables drops the sensor measurement field tables
func DropSensorMeasurementFieldTables(db *sql.DB) error {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS sensor_measurement_fields CASCADE;
	`)
	if err != nil {
		return fmt.Errorf("error dropping sensor_measurement_fields table: %v", err)
	}

	return nil
}
