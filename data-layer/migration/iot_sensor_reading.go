package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"
)

// CreateIoTSensorReadingTable creates the iot_sensor_readings table with all necessary columns
func CreateIoTSensorReadingTable(cfg *config.Config) error {
	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create iot_sensor_readings table with all necessary columns (no indexes or triggers)
	query := `
	CREATE TABLE IF NOT EXISTS iot_sensor_readings (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID,
		asset_sensor_id UUID,
		sensor_type_id UUID,
		mac_address VARCHAR(255),
		location VARCHAR(255),
		
		-- Dynamic measurement fields (primary data structure)
		measurement_data JSONB,
		standard_fields JSONB,
		
		-- Timestamp fields
		reading_time TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE,
		
		-- Deprecated fields (for backward compatibility)
		data_x JSONB,
		data_y JSONB,
		peak_x JSONB,
		peak_y JSONB,
		ppm DECIMAL(10,6),
		label VARCHAR(255),
		raw_data JSONB
	);`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create iot_sensor_readings table: %v", err)
	}

	log.Println("Successfully created iot_sensor_readings table")
	return nil
}

// CreateIoTSensorReadingTableDirect creates the table using an existing database connection
func CreateIoTSensorReadingTableDirect(db *sql.DB) error {
	// Create iot_sensor_readings table with all necessary columns (no indexes or triggers)
	query := `
	CREATE TABLE IF NOT EXISTS iot_sensor_readings (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID,
		asset_sensor_id UUID,
		sensor_type_id UUID,
		mac_address VARCHAR(255),
		location VARCHAR(255),
		
		-- Dynamic measurement fields (primary data structure)
		measurement_data JSONB,
		standard_fields JSONB,
		
		-- Timestamp fields
		reading_time TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE,
		
		-- Deprecated fields (for backward compatibility)
		data_x JSONB,
		data_y JSONB,
		peak_x JSONB,
		peak_y JSONB,
		ppm DECIMAL(10,6),
		label VARCHAR(255),
		raw_data JSONB
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create iot_sensor_readings table: %v", err)
	}

	log.Println("Successfully created iot_sensor_readings table")
	return nil
}
