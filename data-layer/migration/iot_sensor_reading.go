package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// CreateIoTSensorReadingTable creates the iot_sensor_readings table with proper foreign key constraints
func CreateIoTSensorReadingTable(cfg *config.Config) error {
	log.Println("Creating iot_sensor_readings table...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// SQL untuk membuat tabel iot_sensor_readings dengan dukungan flexible measurement data
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS iot_sensor_readings (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_sensor_id UUID NOT NULL,
		sensor_type_id UUID NOT NULL,
		mac_address VARCHAR(255) NULL,
		location_id UUID NULL,
		location_name VARCHAR(500) NULL,
		
		-- Measurement identification
		measurement_type VARCHAR(100) NOT NULL DEFAULT 'value', -- 'raw_value', 'temperature', 'humidity', etc.
		measurement_label VARCHAR(255) NULL,    -- Human readable label
		measurement_unit VARCHAR(50) NULL,      -- Unit of measurement (°C, μg/m³, %, etc.)
		
		-- Flexible value storage
		numeric_value DOUBLE PRECISION NULL,    -- For numeric values
		text_value TEXT NULL,                   -- For text values
		boolean_value BOOLEAN NULL,             -- For boolean values
		
		-- Additional metadata
		data_source VARCHAR(100) NULL DEFAULT 'json',          -- 'json', 'text', 'csv', etc.
		original_field_name VARCHAR(255) NULL,  -- Original field name from JSON/text
		
		reading_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_iot_readings_asset_sensor_id 
			FOREIGN KEY (asset_sensor_id) REFERENCES asset_sensors(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_iot_readings_sensor_type_id 
			FOREIGN KEY (sensor_type_id) REFERENCES sensor_types(id) 
			ON DELETE RESTRICT ON UPDATE CASCADE,
		CONSTRAINT fk_iot_readings_location_id
			FOREIGN KEY (location_id) REFERENCES locations(id)
			ON DELETE SET NULL ON UPDATE CASCADE,
			
		-- Ensure only one value type is used per record
		CONSTRAINT chk_single_value_type CHECK (
			(numeric_value IS NOT NULL AND text_value IS NULL AND boolean_value IS NULL) OR
			(numeric_value IS NULL AND text_value IS NOT NULL AND boolean_value IS NULL) OR
			(numeric_value IS NULL AND text_value IS NULL AND boolean_value IS NOT NULL) OR
			(numeric_value IS NULL AND text_value IS NULL AND boolean_value IS NULL) -- Allow all NULL
		)
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_iot_readings_tenant_id ON iot_sensor_readings(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_iot_readings_asset_sensor_id ON iot_sensor_readings(asset_sensor_id);
	CREATE INDEX IF NOT EXISTS idx_iot_readings_sensor_type_id ON iot_sensor_readings(sensor_type_id);
	CREATE INDEX IF NOT EXISTS idx_iot_readings_reading_time ON iot_sensor_readings(reading_time);
	CREATE INDEX IF NOT EXISTS idx_iot_readings_created_at ON iot_sensor_readings(created_at);
	CREATE INDEX IF NOT EXISTS idx_iot_readings_mac_address ON iot_sensor_readings(mac_address);
	CREATE INDEX IF NOT EXISTS idx_iot_readings_location_id ON iot_sensor_readings(location_id);
	CREATE INDEX IF NOT EXISTS idx_iot_readings_measurement_type ON iot_sensor_readings(measurement_type);
	CREATE INDEX IF NOT EXISTS idx_iot_readings_numeric_value ON iot_sensor_readings(numeric_value);
	CREATE INDEX IF NOT EXISTS idx_iot_readings_data_source ON iot_sensor_readings(data_source);
	CREATE INDEX IF NOT EXISTS idx_iot_readings_composite ON iot_sensor_readings(asset_sensor_id, measurement_type, reading_time);
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create iot_sensor_readings table: %v", err)
	}

	log.Println("IoT sensor readings table created successfully")
	return nil
}

// CreateIoTSensorReadingTableIfNotExists creates the iot_sensor_readings table if it doesn't exist
func CreateIoTSensorReadingTableIfNotExists(db *sql.DB) error {
	// Check if table exists
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'iot_sensor_readings'
	)`

	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if iot_sensor_readings table exists: %v", err)
	}

	if exists {
		log.Println("IoT sensor readings table already exists")
		return nil
	}

	// Use the same table definition as CreateIoTSensorReadingTable
	createTableSQL := `
	CREATE TABLE iot_sensor_readings (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_sensor_id UUID NOT NULL,
		sensor_type_id UUID NOT NULL,
		mac_address VARCHAR(255) NULL,
		location_id UUID NULL,
		location_name VARCHAR(500) NULL,
		
		-- Measurement identification
		measurement_type VARCHAR(100) NOT NULL DEFAULT 'value', -- 'raw_value', 'temperature', 'humidity', etc.
		measurement_label VARCHAR(255) NULL,    -- Human readable label
		measurement_unit VARCHAR(50) NULL,      -- Unit of measurement (°C, μg/m³, %, etc.)
		
		-- Flexible value storage
		numeric_value DOUBLE PRECISION NULL,    -- For numeric values
		text_value TEXT NULL,                   -- For text values
		boolean_value BOOLEAN NULL,             -- For boolean values
		
		-- Additional metadata
		data_source VARCHAR(100) NULL DEFAULT 'json',          -- 'json', 'text', 'csv', etc.
		original_field_name VARCHAR(255) NULL,  -- Original field name from JSON/text
		
		reading_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_iot_readings_asset_sensor_id 
			FOREIGN KEY (asset_sensor_id) REFERENCES asset_sensors(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_iot_readings_sensor_type_id 
			FOREIGN KEY (sensor_type_id) REFERENCES sensor_types(id) 
			ON DELETE RESTRICT ON UPDATE CASCADE,
		CONSTRAINT fk_iot_readings_location_id
			FOREIGN KEY (location_id) REFERENCES locations(id)
			ON DELETE SET NULL ON UPDATE CASCADE,
			
		-- Ensure only one value type is used per record
		CONSTRAINT chk_single_value_type CHECK (
			(numeric_value IS NOT NULL AND text_value IS NULL AND boolean_value IS NULL) OR
			(numeric_value IS NULL AND text_value IS NOT NULL AND boolean_value IS NULL) OR
			(numeric_value IS NULL AND text_value IS NULL AND boolean_value IS NOT NULL) OR
			(numeric_value IS NULL AND text_value IS NULL AND boolean_value IS NULL) -- Allow all NULL
		)
	);

	-- Create indexes for better query performance
	CREATE INDEX idx_iot_readings_tenant_id ON iot_sensor_readings(tenant_id);
	CREATE INDEX idx_iot_readings_asset_sensor_id ON iot_sensor_readings(asset_sensor_id);
	CREATE INDEX idx_iot_readings_sensor_type_id ON iot_sensor_readings(sensor_type_id);
	CREATE INDEX idx_iot_readings_reading_time ON iot_sensor_readings(reading_time);
	CREATE INDEX idx_iot_readings_created_at ON iot_sensor_readings(created_at);
	CREATE INDEX idx_iot_readings_mac_address ON iot_sensor_readings(mac_address);
	CREATE INDEX idx_iot_readings_location_id ON iot_sensor_readings(location_id);
	CREATE INDEX idx_iot_readings_measurement_type ON iot_sensor_readings(measurement_type);
	CREATE INDEX idx_iot_readings_numeric_value ON iot_sensor_readings(numeric_value);
	CREATE INDEX idx_iot_readings_data_source ON iot_sensor_readings(data_source);
	CREATE INDEX idx_iot_readings_composite ON iot_sensor_readings(asset_sensor_id, measurement_type, reading_time);
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create iot_sensor_readings table: %v", err)
	}

	log.Println("IoT sensor readings table created successfully")
	return nil
}

// CreateIoTSensorReadingTableDirect creates the iot_sensor_readings table directly with database connection
func CreateIoTSensorReadingTableDirect(db *sql.DB) error {
	// Check if table exists
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'iot_sensor_readings'
	)`

	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if iot_sensor_readings table exists: %v", err)
	}

	if exists {
		log.Println("IoT sensor readings table already exists")
		return nil
	}

	// Use the same table definition as CreateIoTSensorReadingTable
	createTableSQL := `
	CREATE TABLE iot_sensor_readings (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_sensor_id UUID NOT NULL,
		sensor_type_id UUID NOT NULL,
		mac_address VARCHAR(255) NULL,
		location_id UUID NULL,
		location_name VARCHAR(500) NULL,
		
		-- Measurement identification
		measurement_type VARCHAR(100) NOT NULL DEFAULT 'value', -- 'raw_value', 'temperature', 'humidity', etc.
		measurement_label VARCHAR(255) NULL,    -- Human readable label
		measurement_unit VARCHAR(50) NULL,      -- Unit of measurement (°C, μg/m³, %, etc.)
		
		-- Flexible value storage
		numeric_value DOUBLE PRECISION NULL,    -- For numeric values
		text_value TEXT NULL,                   -- For text values
		boolean_value BOOLEAN NULL,             -- For boolean values
		
		-- Additional metadata
		data_source VARCHAR(100) NULL DEFAULT 'json',          -- 'json', 'text', 'csv', etc.
		original_field_name VARCHAR(255) NULL,  -- Original field name from JSON/text
		
		reading_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		CONSTRAINT fk_iot_readings_asset_sensor_id 
			FOREIGN KEY (asset_sensor_id) REFERENCES asset_sensors(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_iot_readings_sensor_type_id 
			FOREIGN KEY (sensor_type_id) REFERENCES sensor_types(id) 
			ON DELETE RESTRICT ON UPDATE CASCADE,
		CONSTRAINT fk_iot_readings_location_id
			FOREIGN KEY (location_id) REFERENCES locations(id)
			ON DELETE SET NULL ON UPDATE CASCADE,
			
		-- Ensure only one value type is used per record
		CONSTRAINT chk_single_value_type CHECK (
			(numeric_value IS NOT NULL AND text_value IS NULL AND boolean_value IS NULL) OR
			(numeric_value IS NULL AND text_value IS NOT NULL AND boolean_value IS NULL) OR
			(numeric_value IS NULL AND text_value IS NULL AND boolean_value IS NOT NULL) OR
			(numeric_value IS NULL AND text_value IS NULL AND boolean_value IS NULL) -- Allow all NULL
		)
	);

	-- Create indexes for better query performance
	CREATE INDEX idx_iot_readings_tenant_id ON iot_sensor_readings(tenant_id);
	CREATE INDEX idx_iot_readings_asset_sensor_id ON iot_sensor_readings(asset_sensor_id);
	CREATE INDEX idx_iot_readings_sensor_type_id ON iot_sensor_readings(sensor_type_id);
	CREATE INDEX idx_iot_readings_reading_time ON iot_sensor_readings(reading_time);
	CREATE INDEX idx_iot_readings_created_at ON iot_sensor_readings(created_at);
	CREATE INDEX idx_iot_readings_mac_address ON iot_sensor_readings(mac_address);
	CREATE INDEX idx_iot_readings_location_id ON iot_sensor_readings(location_id);
	CREATE INDEX idx_iot_readings_measurement_type ON iot_sensor_readings(measurement_type);
	CREATE INDEX idx_iot_readings_numeric_value ON iot_sensor_readings(numeric_value);
	CREATE INDEX idx_iot_readings_data_source ON iot_sensor_readings(data_source);
	CREATE INDEX idx_iot_readings_composite ON iot_sensor_readings(asset_sensor_id, measurement_type, reading_time);
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create iot_sensor_readings table: %v", err)
	}

	log.Println("IoT sensor readings table created successfully")
	return nil
}

// DropIoTSensorReadingTable drops the iot_sensor_readings table if it exists
func DropIoTSensorReadingTable(cfg *config.Config) error {
	log.Println("Dropping iot_sensor_readings table...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// SQL untuk drop tabel iot_sensor_readings
	dropTableSQL := `DROP TABLE IF EXISTS iot_sensor_readings CASCADE;`

	// Execute the SQL
	_, err = db.Exec(dropTableSQL)
	if err != nil {
		return fmt.Errorf("failed to drop iot_sensor_readings table: %v", err)
	}

	log.Println("IoT sensor readings table dropped successfully")
	return nil
}
