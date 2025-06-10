package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// CreateSensorLogsTable creates the sensor_logs table for comprehensive logging
func CreateSensorLogsTable(cfg *config.Config) error {
	log.Println("Creating sensor_logs table...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// SQL untuk membuat tabel sensor_logs
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sensor_logs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_sensor_id UUID NOT NULL,
		log_type VARCHAR(50) NOT NULL DEFAULT 'system',
		log_level VARCHAR(50) NOT NULL DEFAULT 'info',
		message TEXT NOT NULL,
		component VARCHAR(100) NULL,
		event_type VARCHAR(100) NULL,
		error_code VARCHAR(50) NULL,
		
		-- Connection History Fields (for connection logs)
		connection_type VARCHAR(100) NULL,
		connection_status VARCHAR(50) NULL,
		ip_address INET NULL,
		mac_address VARCHAR(17) NULL,
		network_name VARCHAR(255) NULL,
		connection_duration BIGINT NULL CHECK (connection_duration >= 0),
		
		-- Flexible metadata
		metadata JSONB DEFAULT '{}'::jsonb,
		source_ip INET NULL,
		user_agent TEXT NULL,
		session_id VARCHAR(255) NULL,
		
		-- Timestamps
		recorded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_sensor_logs_asset_sensor_id 
			FOREIGN KEY (asset_sensor_id) REFERENCES asset_sensors(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT chk_sensor_log_type 
			CHECK (log_type IN ('reading', 'connection', 'battery', 'signal', 'error', 'system', 'maintenance')),
		CONSTRAINT chk_sensor_log_level 
			CHECK (log_level IN ('debug', 'info', 'warning', 'error', 'critical')),
		CONSTRAINT chk_sensor_log_component 
			CHECK (component IS NULL OR component IN ('sensor', 'communication', 'battery', 'hardware', 'software', 'network')),
		CONSTRAINT chk_sensor_log_event_type 
			CHECK (event_type IS NULL OR event_type IN ('startup', 'shutdown', 'connected', 'disconnected', 'reading', 'error', 'maintenance', 'calibration', 'alert')),
		CONSTRAINT chk_connection_status 
			CHECK (connection_status IS NULL OR connection_status IN ('connected', 'disconnected', 'failed')),
		CONSTRAINT chk_connection_type 
			CHECK (connection_type IS NULL OR connection_type IN ('wifi', 'cellular', 'lora', 'zigbee', 'bluetooth', 'ethernet', 'other'))
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_tenant_id ON sensor_logs(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_asset_sensor_id ON sensor_logs(asset_sensor_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_log_type ON sensor_logs(log_type);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_log_level ON sensor_logs(log_level);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_component ON sensor_logs(component);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_event_type ON sensor_logs(event_type);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_recorded_at ON sensor_logs(recorded_at);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_created_at ON sensor_logs(created_at);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_error_code ON sensor_logs(error_code);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_session_id ON sensor_logs(session_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_connection_status ON sensor_logs(connection_status);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_connection_type ON sensor_logs(connection_type);
	
	-- Full-text search index for message content
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_message_fts ON sensor_logs USING gin(to_tsvector('english', message));
	
	-- Composite indexes for common queries
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_sensor_type_level ON sensor_logs(asset_sensor_id, log_type, log_level);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_time_range ON sensor_logs(recorded_at DESC, asset_sensor_id);
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create sensor_logs table: %v", err)
	}

	log.Println("sensor_logs table created successfully")
	return nil
}

// CreateSensorLogsTableIfNotExists creates the sensor_logs table only if it doesn't exist
func CreateSensorLogsTableIfNotExists(db *sql.DB) error {
	log.Println("Checking and creating sensor_logs table if not exists...")

	// Check if table exists
	var exists bool
	checkTableSQL := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'sensor_logs'
		);`

	err := db.QueryRow(checkTableSQL).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if sensor_logs table exists: %v", err)
	}

	if exists {
		log.Println("sensor_logs table already exists, skipping creation")
		return nil
	}

	// Table doesn't exist, create it directly using the provided db connection
	log.Println("Creating sensor_logs table...")

	// SQL for creating the sensor_logs table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sensor_logs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_sensor_id UUID NOT NULL,
		log_type VARCHAR(50) NOT NULL DEFAULT 'system',
		log_level VARCHAR(50) NOT NULL DEFAULT 'info',
		message TEXT NOT NULL,
		component VARCHAR(100) NULL,
		event_type VARCHAR(100) NULL,
		error_code VARCHAR(50) NULL,
		
		-- Connection History Fields (for connection logs)
		connection_type VARCHAR(100) NULL,
		connection_status VARCHAR(50) NULL,
		ip_address INET NULL,
		mac_address VARCHAR(17) NULL,
		network_name VARCHAR(255) NULL,
		connection_duration BIGINT NULL CHECK (connection_duration >= 0),
		
		-- Flexible metadata
		metadata JSONB DEFAULT '{}'::jsonb,
		source_ip INET NULL,
		user_agent TEXT NULL,
		session_id VARCHAR(255) NULL,
		
		-- Timestamps
		recorded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_sensor_logs_asset_sensor_id 
			FOREIGN KEY (asset_sensor_id) REFERENCES asset_sensors(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT chk_sensor_log_type 
			CHECK (log_type IN ('reading', 'connection', 'battery', 'signal', 'error', 'system', 'maintenance')),
		CONSTRAINT chk_sensor_log_level 
			CHECK (log_level IN ('debug', 'info', 'warning', 'error', 'critical')),
		CONSTRAINT chk_sensor_log_component 
			CHECK (component IS NULL OR component IN ('sensor', 'communication', 'battery', 'hardware', 'software', 'network')),
		CONSTRAINT chk_sensor_log_event_type 
			CHECK (event_type IS NULL OR event_type IN ('startup', 'shutdown', 'connected', 'disconnected', 'reading', 'error', 'maintenance', 'calibration', 'alert')),
		CONSTRAINT chk_connection_status 
			CHECK (connection_status IS NULL OR connection_status IN ('connected', 'disconnected', 'failed')),
		CONSTRAINT chk_connection_type 
			CHECK (connection_type IS NULL OR connection_type IN ('wifi', 'cellular', 'lora', 'zigbee', 'bluetooth', 'ethernet', 'other'))
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_tenant_id ON sensor_logs(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_asset_sensor_id ON sensor_logs(asset_sensor_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_log_type ON sensor_logs(log_type);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_log_level ON sensor_logs(log_level);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_component ON sensor_logs(component);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_event_type ON sensor_logs(event_type);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_recorded_at ON sensor_logs(recorded_at);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_created_at ON sensor_logs(created_at);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_error_code ON sensor_logs(error_code);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_session_id ON sensor_logs(session_id);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_connection_status ON sensor_logs(connection_status);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_connection_type ON sensor_logs(connection_type);
	
	-- Full-text search index for message content
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_message_fts ON sensor_logs USING gin(to_tsvector('english', message));
	
	-- Composite indexes for common queries
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_sensor_type_level ON sensor_logs(asset_sensor_id, log_type, log_level);
	CREATE INDEX IF NOT EXISTS idx_sensor_logs_time_range ON sensor_logs(recorded_at DESC, asset_sensor_id);
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create sensor_logs table: %v", err)
	}

	log.Println("sensor_logs table created successfully")
	return nil
}

// DropSensorLogsTable drops the sensor_logs table
func DropSensorLogsTable(cfg *config.Config) error {
	log.Println("Dropping sensor_logs table...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Drop the table
	dropTableSQL := `DROP TABLE IF EXISTS sensor_logs CASCADE;`
	_, err = db.Exec(dropTableSQL)
	if err != nil {
		return fmt.Errorf("failed to drop sensor_logs table: %v", err)
	}

	log.Println("sensor_logs table dropped successfully")
	return nil
}
