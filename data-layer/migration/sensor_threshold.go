package migration

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateSensorThresholdTable creates the sensor_thresholds table
func CreateSensorThresholdTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sensor_thresholds (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NOT NULL,
		asset_sensor_id UUID NOT NULL,
		measurement_type_id UUID NOT NULL,
		measurement_field_name VARCHAR(255) NOT NULL,
		min_value DOUBLE PRECISION NULL,
		max_value DOUBLE PRECISION NULL,
		severity VARCHAR(20) NOT NULL DEFAULT 'warning' CHECK (severity IN ('warning', 'critical')),
		is_active BOOLEAN NOT NULL DEFAULT true,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_sensor_thresholds_asset_sensor_id 
			FOREIGN KEY (asset_sensor_id) REFERENCES asset_sensors(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
			
		CONSTRAINT fk_sensor_thresholds_measurement_type_id 
			FOREIGN KEY (measurement_type_id) REFERENCES sensor_measurement_types(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
			
		-- Ensure at least one threshold value is set
		CONSTRAINT chk_threshold_values CHECK (
			min_value IS NOT NULL OR max_value IS NOT NULL
		),
		
		-- Ensure min_value < max_value when both are set
		CONSTRAINT chk_min_max_values CHECK (
			min_value IS NULL OR max_value IS NULL OR min_value < max_value
		),
		
		-- Unique constraint to prevent duplicate thresholds for same field and severity
		CONSTRAINT uq_sensor_threshold_field_severity UNIQUE (
			asset_sensor_id, measurement_field_name, severity
		)
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_tenant_id ON sensor_thresholds(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_asset_sensor_id ON sensor_thresholds(asset_sensor_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_measurement_type_id ON sensor_thresholds(measurement_type_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_measurement_field ON sensor_thresholds(measurement_field_name);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_active ON sensor_thresholds(is_active);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_severity ON sensor_thresholds(severity);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_composite ON sensor_thresholds(asset_sensor_id, measurement_field_name, is_active);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating sensor_thresholds table: %v", err)
		return fmt.Errorf("failed to create sensor_thresholds table: %w", err)
	}

	log.Println("Successfully created sensor_thresholds table")
	return nil
}

// CreateSensorThresholdTableIfNotExists creates the sensor_thresholds table if it doesn't exist
func CreateSensorThresholdTableIfNotExists(db *sql.DB) error {
	log.Println("Creating sensor_thresholds table if it doesn't exist...")
	return CreateSensorThresholdTable(db)
}
