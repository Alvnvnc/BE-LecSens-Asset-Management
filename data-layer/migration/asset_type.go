package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"
)

// CreateAssetTypeTable creates the asset_types table with all necessary columns
func CreateAssetTypeTable(cfg *config.Config) error {
	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create asset_types table with all necessary columns
	query := `
	CREATE TABLE IF NOT EXISTS asset_types (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL UNIQUE,
		category VARCHAR(255),
		description TEXT,
		properties_schema JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create asset_types table: %v", err)
	}

	log.Println("Successfully created asset_types table")
	return nil
}
