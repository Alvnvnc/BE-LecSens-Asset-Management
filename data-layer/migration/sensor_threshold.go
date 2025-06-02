package migration

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateSensorThresholdTable creates the sensor_thresholds table
func CreateSensorThresholdTable(db *sql.DB) error {
	log.Println("Creating sensor_thresholds table...")

	// Create sensor_thresholds table
	query := `
	CREATE TABLE IF NOT EXISTS sensor_thresholds (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NOT NULL,
		asset_sensor_id UUID NOT NULL REFERENCES asset_sensors(id) ON DELETE CASCADE,
		sensor_type_id UUID NOT NULL REFERENCES sensor_types(id) ON DELETE CASCADE,
		measurement_field VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		min_value DOUBLE PRECISION NOT NULL,
		max_value DOUBLE PRECISION NOT NULL,
		severity VARCHAR(20) NOT NULL CHECK (severity IN ('warning', 'critical')),
		alert_message TEXT,
		notification_rules JSONB,
		is_active BOOLEAN NOT NULL DEFAULT true,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_tenant_id ON sensor_thresholds(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_asset_sensor_id ON sensor_thresholds(asset_sensor_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_sensor_type_id ON sensor_thresholds(sensor_type_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_measurement_field ON sensor_thresholds(measurement_field);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_severity ON sensor_thresholds(severity);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_is_active ON sensor_thresholds(is_active);
	CREATE INDEX IF NOT EXISTS idx_sensor_thresholds_created_at ON sensor_thresholds(created_at);
	-- Create trigger for updated_at
	CREATE OR REPLACE FUNCTION update_sensor_thresholds_updated_at()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	-- Drop trigger if exists and create it
	DROP TRIGGER IF EXISTS trigger_update_sensor_thresholds_updated_at ON sensor_thresholds;
	CREATE TRIGGER trigger_update_sensor_thresholds_updated_at
		BEFORE UPDATE ON sensor_thresholds
		FOR EACH ROW
		EXECUTE FUNCTION update_sensor_thresholds_updated_at();`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create sensor_thresholds table: %v", err)
	}

	log.Println("Successfully created sensor_thresholds table")
	return nil
}