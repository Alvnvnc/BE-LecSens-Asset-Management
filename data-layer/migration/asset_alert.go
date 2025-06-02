package migration

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateAssetAlertTable creates the asset_alerts table
func CreateAssetAlertTable(db *sql.DB) error {
	log.Println("Creating asset_alerts table...")

	// Create asset_alerts table
	query := `
	CREATE TABLE IF NOT EXISTS asset_alerts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NOT NULL,
		asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
		asset_sensor_id UUID NOT NULL REFERENCES asset_sensors(id) ON DELETE CASCADE,
		threshold_id UUID NOT NULL REFERENCES sensor_thresholds(id) ON DELETE CASCADE,
		alert_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		resolved_time TIMESTAMP WITH TIME ZONE,
		severity VARCHAR(20) NOT NULL CHECK (severity IN ('warning', 'critical'))
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_tenant_id ON asset_alerts(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_asset_id ON asset_alerts(asset_id);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_asset_sensor_id ON asset_alerts(asset_sensor_id);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_threshold_id ON asset_alerts(threshold_id);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_alert_time ON asset_alerts(alert_time);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_resolved_time ON asset_alerts(resolved_time);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_severity ON asset_alerts(severity);
	CREATE INDEX IF NOT EXISTS idx_asset_alerts_unresolved ON asset_alerts(resolved_time) WHERE resolved_time IS NULL;`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create asset_alerts table: %v", err)
	}

	log.Println("Successfully created asset_alerts table")
	return nil
}