package migration

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateAssetActivityTable creates the asset_activities table
func CreateAssetActivityTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS asset_activities (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NOT NULL,
		asset_id UUID NOT NULL,
		activity_type VARCHAR(50) NOT NULL CHECK (activity_type IN ('maintenance', 'calibration', 'inspection')),
		status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed')),
		scheduled_date TIMESTAMP NOT NULL,
		completed_date TIMESTAMP NULL,
		description TEXT,
		notes TEXT,
		assigned_to UUID NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_asset_activities_asset_id 
			FOREIGN KEY (asset_id) REFERENCES assets(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		
		-- Constraint to ensure completed_date is not before scheduled_date
		CONSTRAINT chk_asset_activities_completed_date 
			CHECK (completed_date IS NULL OR completed_date >= scheduled_date)
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_asset_activities_tenant_id ON asset_activities(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_asset_activities_asset_id ON asset_activities(asset_id);
	CREATE INDEX IF NOT EXISTS idx_asset_activities_activity_type ON asset_activities(activity_type);
	CREATE INDEX IF NOT EXISTS idx_asset_activities_status ON asset_activities(status);
	CREATE INDEX IF NOT EXISTS idx_asset_activities_scheduled_date ON asset_activities(scheduled_date);
	CREATE INDEX IF NOT EXISTS idx_asset_activities_completed_date ON asset_activities(completed_date);
	CREATE INDEX IF NOT EXISTS idx_asset_activities_assigned_to ON asset_activities(assigned_to);
	CREATE INDEX IF NOT EXISTS idx_asset_activities_created_at ON asset_activities(created_at);
	CREATE INDEX IF NOT EXISTS idx_asset_activities_tenant_asset ON asset_activities(tenant_id, asset_id);
	CREATE INDEX IF NOT EXISTS idx_asset_activities_tenant_status ON asset_activities(tenant_id, status);
	CREATE INDEX IF NOT EXISTS idx_asset_activities_overdue ON asset_activities(status, scheduled_date) WHERE status = 'pending';
	`

	// Execute the SQL
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create asset_activities table: %v", err)
	}

	log.Println("Asset activities table created successfully")
	return nil
}

// CreateAssetActivityTableIfNotExists creates the asset_activities table if it doesn't exist
func CreateAssetActivityTableIfNotExists(db *sql.DB) error {
	log.Println("Creating asset_activities table if it doesn't exist...")
	return CreateAssetActivityTable(db)
}
