package migration

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateAssetAlertTable creates the asset_alerts table
func CreateAssetAlertTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS asset_alerts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NOT NULL,
		asset_id UUID NOT NULL,
		asset_sensor_id UUID NOT NULL,
		threshold_id UUID NOT NULL,
		measurement_field_name VARCHAR(255) NOT NULL,
		alert_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		resolved_time TIMESTAMP NULL,
		severity VARCHAR(20) NOT NULL CHECK (severity IN ('warning', 'critical')),
		status VARCHAR(20) NOT NULL DEFAULT 'warning' CHECK (status IN ('normal', 'warning', 'critical')),
		trigger_value DOUBLE PRECISION NOT NULL,
		threshold_min_value DOUBLE PRECISION NULL,
		threshold_max_value DOUBLE PRECISION NULL,
		alert_message TEXT NOT NULL,
		alert_type VARCHAR(20) NOT NULL CHECK (alert_type IN ('min_breach', 'max_breach')),
		is_resolved BOOLEAN NOT NULL DEFAULT false,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_asset_alerts_asset_id 
			FOREIGN KEY (asset_id) REFERENCES assets(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_asset_alerts_asset_sensor_id 
			FOREIGN KEY (asset_sensor_id) REFERENCES asset_sensors(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_asset_alerts_threshold_id 
			FOREIGN KEY (threshold_id) REFERENCES sensor_thresholds(id) 
			ON DELETE RESTRICT ON UPDATE CASCADE
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_tenant_id ON asset_alerts(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_asset_id ON asset_alerts(asset_id);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_asset_sensor_id ON asset_alerts(asset_sensor_id);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_threshold_id ON asset_alerts(threshold_id);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_alert_time ON asset_alerts(alert_time);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_severity ON asset_alerts(severity);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_status ON asset_alerts(status);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_is_resolved ON asset_alerts(is_resolved);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_measurement_field ON asset_alerts(measurement_field_name);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_alert_type ON asset_alerts(alert_type);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_active ON asset_alerts(tenant_id, is_resolved);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_recent ON asset_alerts(alert_time DESC);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_composite ON asset_alerts(asset_sensor_id, is_resolved, alert_time);
	`

	// Execute the SQL
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create asset_alerts table: %v", err)
	}

	log.Println("Asset alerts table created successfully")
	return nil
}

// CreateAssetAlertTableIfNotExists creates the asset_alerts table if it doesn't exist
func CreateAssetAlertTableIfNotExists(db *sql.DB) error {
	log.Println("Creating asset_alerts table if it doesn't exist...")
	return CreateAssetAlertTable(db)
}
