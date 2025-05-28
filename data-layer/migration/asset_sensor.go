package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// CreateAssetSensorTable creates the asset_sensors table with proper foreign key constraints
func CreateAssetSensorTable(cfg *config.Config) error {
	log.Println("Creating asset_sensors table...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// SQL untuk membuat tabel asset_sensors
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS asset_sensors (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_id UUID NOT NULL,
		sensor_type_id UUID NOT NULL,
		name VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'active',
		configuration JSONB DEFAULT '{}'::jsonb,
		last_reading_value DOUBLE PRECISION NULL,
		last_reading_time TIMESTAMP NULL,
		last_reading_values JSONB DEFAULT '{}'::jsonb,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_asset_sensors_asset_id 
			FOREIGN KEY (asset_id) REFERENCES assets(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_asset_sensors_sensor_type_id 
			FOREIGN KEY (sensor_type_id) REFERENCES sensor_types(id) 
			ON DELETE RESTRICT ON UPDATE CASCADE
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_asset_sensors_tenant_id ON asset_sensors(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_asset_sensors_asset_id ON asset_sensors(asset_id);
	CREATE INDEX IF NOT EXISTS idx_asset_sensors_sensor_type_id ON asset_sensors(sensor_type_id);
	CREATE INDEX IF NOT EXISTS idx_asset_sensors_status ON asset_sensors(status);
	CREATE INDEX IF NOT EXISTS idx_asset_sensors_created_at ON asset_sensors(created_at);
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create asset_sensors table: %v", err)
	}

	log.Println("Asset sensors table created successfully")
	return nil
}

// CreateAssetSensorTableIfNotExists creates the asset_sensors table if it doesn't exist
func CreateAssetSensorTableIfNotExists(db *sql.DB) error {
	// Check if table exists
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'asset_sensors'
	)`

	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if asset_sensors table exists: %v", err)
	}

	if exists {
		log.Println("Asset sensors table already exists")
		return nil
	}

	log.Println("Creating asset_sensors table...")

	// SQL untuk membuat tabel asset_sensors
	createTableSQL := `
	CREATE TABLE asset_sensors (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_id UUID NOT NULL,
		sensor_type_id UUID NOT NULL,
		name VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'active',
		configuration JSONB DEFAULT '{}'::jsonb,
		last_reading_value DOUBLE PRECISION NULL,
		last_reading_time TIMESTAMP NULL,
		last_reading_values JSONB DEFAULT '{}'::jsonb,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_asset_sensors_asset_id 
			FOREIGN KEY (asset_id) REFERENCES assets(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_asset_sensors_sensor_type_id 
			FOREIGN KEY (sensor_type_id) REFERENCES sensor_types(id) 
			ON DELETE RESTRICT ON UPDATE CASCADE
	);

	-- Create indexes for better query performance
	CREATE INDEX idx_asset_sensors_tenant_id ON asset_sensors(tenant_id);
	CREATE INDEX idx_asset_sensors_asset_id ON asset_sensors(asset_id);
	CREATE INDEX idx_asset_sensors_sensor_type_id ON asset_sensors(sensor_type_id);
	CREATE INDEX idx_asset_sensors_status ON asset_sensors(status);
	CREATE INDEX idx_asset_sensors_created_at ON asset_sensors(created_at);
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create asset_sensors table: %v", err)
	}

	log.Println("Asset sensors table created successfully")
	return nil
}

// DropAssetSensorTable drops the asset_sensors table if it exists
func DropAssetSensorTable(cfg *config.Config) error {
	log.Println("Dropping asset_sensors table...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// SQL untuk drop tabel asset_sensors
	dropTableSQL := `DROP TABLE IF EXISTS asset_sensors CASCADE;`

	// Execute the SQL
	_, err = db.Exec(dropTableSQL)
	if err != nil {
		return fmt.Errorf("failed to drop asset_sensors table: %v", err)
	}

	log.Println("Asset sensors table dropped successfully")
	return nil
}
