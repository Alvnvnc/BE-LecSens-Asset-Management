package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"
)

// CreateAssetTable creates the assets table with all necessary columns
func CreateAssetTable(cfg *config.Config) error {
	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create assets table with all necessary columns
	// Note: tenant_id is UUID but no foreign key constraint since tenant data comes from external API
	query := `
	CREATE TABLE IF NOT EXISTS assets (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID,
		name VARCHAR(255) NOT NULL,
		asset_type_id UUID NOT NULL REFERENCES asset_types(id) ON DELETE CASCADE,
		location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
		status VARCHAR(50) DEFAULT 'active',
		properties JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_assets_tenant_id ON assets(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_assets_asset_type_id ON assets(asset_type_id);
	CREATE INDEX IF NOT EXISTS idx_assets_location_id ON assets(location_id);
	CREATE INDEX IF NOT EXISTS idx_assets_status ON assets(status);
	CREATE INDEX IF NOT EXISTS idx_assets_created_at ON assets(created_at);`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create assets table: %v", err)
	}

	log.Println("Successfully created assets table")
	return nil
}
