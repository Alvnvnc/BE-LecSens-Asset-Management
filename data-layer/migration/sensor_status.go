package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// CreateSensorStatusTable creates the sensor_status table for real-time sensor monitoring
func CreateSensorStatusTable(cfg *config.Config) error {
	log.Println("Creating sensor_status table...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// SQL untuk membuat tabel sensor_status
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sensor_status (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_sensor_id UUID NOT NULL UNIQUE, -- One status record per sensor
		
		-- Battery Information
		battery_level DOUBLE PRECISION NULL CHECK (battery_level >= 0 AND battery_level <= 100),
		battery_voltage DOUBLE PRECISION NULL CHECK (battery_voltage >= 0),
		battery_status VARCHAR(50) NULL,
		battery_last_charged TIMESTAMP NULL,
		battery_estimated_life INTEGER NULL CHECK (battery_estimated_life >= 0),
		battery_type VARCHAR(50) NULL,
		
		-- Signal Strength Information
		signal_type VARCHAR(100) NULL,
		signal_rssi INTEGER NULL CHECK (signal_rssi >= -120 AND signal_rssi <= 0),
		signal_snr DOUBLE PRECISION NULL,
		signal_quality INTEGER NULL CHECK (signal_quality >= 0 AND signal_quality <= 100),
		signal_frequency DOUBLE PRECISION NULL CHECK (signal_frequency > 0),
		signal_channel INTEGER NULL CHECK (signal_channel >= 0),
		signal_status VARCHAR(50) NULL,
		
		-- Connection Information
		connection_type VARCHAR(100) NULL,
		connection_status VARCHAR(50) NOT NULL DEFAULT 'offline',
		last_connected_at TIMESTAMP NULL,
		last_disconnected_at TIMESTAMP NULL,
		current_ip INET NULL,
		current_network VARCHAR(255) NULL,
		
		-- Additional Status Information
		temperature DOUBLE PRECISION NULL,
		humidity DOUBLE PRECISION NULL CHECK (humidity >= 0 AND humidity <= 100),
		is_online BOOLEAN NOT NULL DEFAULT FALSE,
		last_heartbeat TIMESTAMP NULL,
		firmware_version VARCHAR(100) NULL,
		error_count INTEGER NULL CHECK (error_count >= 0),
		last_error_at TIMESTAMP NULL,
		
		-- Timestamps
		recorded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_sensor_status_asset_sensor_id 
			FOREIGN KEY (asset_sensor_id) REFERENCES asset_sensors(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT chk_battery_status 
			CHECK (battery_status IS NULL OR battery_status IN ('good', 'low', 'critical', 'unknown', 'charging')),
		CONSTRAINT chk_battery_type 
			CHECK (battery_type IS NULL OR battery_type IN ('lithium', 'alkaline', 'rechargeable', 'solar', 'other')),
		CONSTRAINT chk_signal_status 
			CHECK (signal_status IS NULL OR signal_status IN ('excellent', 'good', 'fair', 'poor', 'no_signal')),
		CONSTRAINT chk_signal_type 
			CHECK (signal_type IS NULL OR signal_type IN ('wifi', 'cellular', 'lora', 'zigbee', 'bluetooth', 'ethernet', 'other')),
		CONSTRAINT chk_connection_status 
			CHECK (connection_status IN ('online', 'offline', 'connecting', 'error')),
		CONSTRAINT chk_connection_type 
			CHECK (connection_type IS NULL OR connection_type IN ('wifi', 'cellular', 'lora', 'zigbee', 'bluetooth', 'ethernet', 'other')),
		CONSTRAINT chk_connection_timestamps 
			CHECK (last_connected_at IS NULL OR last_disconnected_at IS NULL OR last_disconnected_at >= last_connected_at OR connection_status = 'online')
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_sensor_status_tenant_id ON sensor_status(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_asset_sensor_id ON sensor_status(asset_sensor_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_connection_status ON sensor_status(connection_status);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_is_online ON sensor_status(is_online);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_battery_level ON sensor_status(battery_level);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_battery_status ON sensor_status(battery_status);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_signal_status ON sensor_status(signal_status);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_signal_rssi ON sensor_status(signal_rssi);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_recorded_at ON sensor_status(recorded_at);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_last_heartbeat ON sensor_status(last_heartbeat);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_error_count ON sensor_status(error_count);
	
	-- Composite indexes for common real-time queries
	CREATE INDEX IF NOT EXISTS idx_sensor_status_health ON sensor_status(is_online, battery_level, signal_rssi);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_alerts ON sensor_status(battery_status, signal_status, error_count);
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create sensor_status table: %v", err)
	}

	log.Println("sensor_status table created successfully")
	return nil
}

// CreateSensorStatusTableIfNotExists creates the sensor_status table only if it doesn't exist
func CreateSensorStatusTableIfNotExists(db *sql.DB) error {
	log.Println("Checking and creating sensor_status table if not exists...")

	// Check if table exists
	var exists bool
	checkTableSQL := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'sensor_status'
		);`

	err := db.QueryRow(checkTableSQL).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if sensor_status table exists: %v", err)
	}

	if exists {
		log.Println("sensor_status table already exists, skipping creation")
		return nil
	}

	// Table doesn't exist, create it directly using the provided db connection
	log.Println("Creating sensor_status table...")

	// SQL for creating the sensor_status table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sensor_status (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_sensor_id UUID NOT NULL UNIQUE, -- One status record per sensor
		
		-- Battery Information
		battery_level DOUBLE PRECISION NULL CHECK (battery_level >= 0 AND battery_level <= 100),
		battery_voltage DOUBLE PRECISION NULL CHECK (battery_voltage >= 0),
		battery_status VARCHAR(50) NULL,
		battery_last_charged TIMESTAMP NULL,
		battery_estimated_life INTEGER NULL CHECK (battery_estimated_life >= 0),
		battery_type VARCHAR(50) NULL,
		
		-- Signal Strength Information
		signal_type VARCHAR(100) NULL,
		signal_rssi INTEGER NULL CHECK (signal_rssi >= -120 AND signal_rssi <= 0),
		signal_snr DOUBLE PRECISION NULL,
		signal_quality INTEGER NULL CHECK (signal_quality >= 0 AND signal_quality <= 100),
		signal_frequency DOUBLE PRECISION NULL CHECK (signal_frequency > 0),
		signal_channel INTEGER NULL CHECK (signal_channel >= 0),
		signal_status VARCHAR(50) NULL,
		
		-- Connection Information
		connection_type VARCHAR(100) NULL,
		connection_status VARCHAR(50) NOT NULL DEFAULT 'offline',
		last_connected_at TIMESTAMP NULL,
		last_disconnected_at TIMESTAMP NULL,
		current_ip INET NULL,
		current_network VARCHAR(255) NULL,
		
		-- Additional Status Information
		temperature DOUBLE PRECISION NULL,
		humidity DOUBLE PRECISION NULL CHECK (humidity >= 0 AND humidity <= 100),
		is_online BOOLEAN NOT NULL DEFAULT FALSE,
		last_heartbeat TIMESTAMP NULL,
		firmware_version VARCHAR(100) NULL,
		error_count INTEGER NULL CHECK (error_count >= 0),
		last_error_at TIMESTAMP NULL,
		
		-- Timestamps
		recorded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_sensor_status_asset_sensor_id 
			FOREIGN KEY (asset_sensor_id) REFERENCES asset_sensors(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT chk_battery_status 
			CHECK (battery_status IS NULL OR battery_status IN ('good', 'low', 'critical', 'unknown', 'charging')),
		CONSTRAINT chk_battery_type 
			CHECK (battery_type IS NULL OR battery_type IN ('lithium', 'alkaline', 'rechargeable', 'solar', 'other')),
		CONSTRAINT chk_signal_status 
			CHECK (signal_status IS NULL OR signal_status IN ('excellent', 'good', 'fair', 'poor', 'no_signal')),
		CONSTRAINT chk_signal_type 
			CHECK (signal_type IS NULL OR signal_type IN ('wifi', 'cellular', 'lora', 'zigbee', 'bluetooth', 'ethernet', 'other')),
		CONSTRAINT chk_connection_status 
			CHECK (connection_status IN ('online', 'offline', 'connecting', 'error')),
		CONSTRAINT chk_connection_type 
			CHECK (connection_type IS NULL OR connection_type IN ('wifi', 'cellular', 'lora', 'zigbee', 'bluetooth', 'ethernet', 'other')),
		CONSTRAINT chk_connection_timestamps 
			CHECK (last_connected_at IS NULL OR last_disconnected_at IS NULL OR last_disconnected_at >= last_connected_at OR connection_status = 'online')
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_sensor_status_tenant_id ON sensor_status(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_asset_sensor_id ON sensor_status(asset_sensor_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_connection_status ON sensor_status(connection_status);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_is_online ON sensor_status(is_online);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_battery_level ON sensor_status(battery_level);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_battery_status ON sensor_status(battery_status);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_signal_status ON sensor_status(signal_status);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_signal_rssi ON sensor_status(signal_rssi);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_recorded_at ON sensor_status(recorded_at);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_last_heartbeat ON sensor_status(last_heartbeat);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_error_count ON sensor_status(error_count);
	
	-- Composite indexes for common real-time queries
	CREATE INDEX IF NOT EXISTS idx_sensor_status_health ON sensor_status(is_online, battery_level, signal_rssi);
	CREATE INDEX IF NOT EXISTS idx_sensor_status_alerts ON sensor_status(battery_status, signal_status, error_count);
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create sensor_status table: %v", err)
	}

	log.Println("sensor_status table created successfully")
	return nil
}

// DropSensorStatusTable drops the sensor_status table
func DropSensorStatusTable(cfg *config.Config) error {
	log.Println("Dropping sensor_status table...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Drop the table
	dropTableSQL := `DROP TABLE IF EXISTS sensor_status CASCADE;`
	_, err = db.Exec(dropTableSQL)
	if err != nil {
		return fmt.Errorf("failed to drop sensor_status table: %v", err)
	}

	log.Println("sensor_status table dropped successfully")
	return nil
}
