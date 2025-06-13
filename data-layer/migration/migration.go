package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// MigrateDatabase creates database if it doesn't exist and runs all migrations
func MigrateDatabase(cfg *config.Config) error {
	// Connect directly to the target database (it's already created by PostgreSQL initialization)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	appDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to application database: %v", err)
	}
	defer appDB.Close()

	// Test the connection
	if err := appDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to database '%s'", cfg.DB.Name)

	// Create locations table if not exists
	err = CreateLocationsTableIfNotExists(appDB)
	if err != nil {
		return fmt.Errorf("failed to create locations table: %v", err)
	}

	// Run migrations
	err = runMigrations(cfg)
	if err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("All migrations completed successfully")
	return nil
}

// runMigrations executes all database migrations in order
func runMigrations(cfg *config.Config) error {
	log.Println("Running database migrations...")

	// Run asset type migration
	if err := CreateAssetTypeTable(cfg); err != nil {
		return fmt.Errorf("asset type migration failed: %v", err)
	}

	// Run asset migration
	if err := CreateAssetTable(cfg); err != nil {
		return fmt.Errorf("asset migration failed: %v", err)
	}

	// Run asset document migration
	if err := CreateAssetDocumentTable(cfg); err != nil {
		return fmt.Errorf("asset document migration failed: %v", err)
	}

	// Connect to database for sensor migrations
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database for sensor migrations: %v", err)
	}
	defer db.Close()

	// Run sensor type migration
	if err := CreateSensorTypeTables(db); err != nil {
		return fmt.Errorf("sensor type migration failed: %v", err)
	}

	// Create sensor measurement type tables
	log.Println("Creating sensor measurement type tables...")
	err = CreateSensorMeasurementTypeTables(db)
	if err != nil {
		return fmt.Errorf("sensor measurement type migration failed: %v", err)
	}
	log.Println("Sensor measurement type tables created successfully")

	// Run sensor measurement field migration
	if err := CreateSensorMeasurementFieldTables(db); err != nil {
		return fmt.Errorf("sensor measurement field migration failed: %v", err)
	}

	// Run asset sensors migration (required before IoT sensor readings)
	log.Println("Creating asset sensors table...")
	if err := CreateAssetSensorTableIfNotExists(db); err != nil {
		return fmt.Errorf("asset sensor migration failed: %v", err)
	}
	log.Println("Asset sensors table created successfully")

	// Ensure required tables exist before creating IoT sensor readings table
	var sensorTypesExists, assetSensorsExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'sensor_types')").Scan(&sensorTypesExists)
	if err != nil {
		return fmt.Errorf("failed to check if sensor_types table exists: %v", err)
	}

	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'asset_sensors')").Scan(&assetSensorsExists)
	if err != nil {
		return fmt.Errorf("failed to check if asset_sensors table exists: %v", err)
	}

	if !sensorTypesExists {
		return fmt.Errorf("sensor_types table does not exist, required for IoT sensor readings table")
	}

	if !assetSensorsExists {
		return fmt.Errorf("asset_sensors table does not exist, required for IoT sensor readings table")
	}

	// Run IoT sensor reading migration with flexible measurement support
	log.Println("Creating IoT sensor readings table with flexible measurement support...")
	if err := CreateIoTSensorReadingTableDirect(db); err != nil {
		return fmt.Errorf("iot sensor reading migration failed: %v", err)
	}
	log.Println("IoT sensor readings table created successfully")

	// Run sensor threshold migration
	log.Println("Creating sensor thresholds table...")
	if err := CreateSensorThresholdTableIfNotExists(db); err != nil {
		return fmt.Errorf("sensor threshold migration failed: %v", err)
	}
	log.Println("Sensor thresholds table created successfully")

	// Run asset alert migration
	log.Println("Creating asset alerts table...")
	if err := CreateAssetAlertTableIfNotExists(db); err != nil {
		return fmt.Errorf("asset alert migration failed: %v", err)
	}
	log.Println("Asset alerts table created successfully")

	// Run asset activity migration
	log.Println("Creating asset activities table...")
	if err := CreateAssetActivityTableIfNotExists(db); err != nil {
		return fmt.Errorf("asset activity migration failed: %v", err)
	}
	log.Println("Asset activities table created successfully")

	// Run sensor status migration
	log.Println("Creating sensor status table...")
	if err := CreateSensorStatusTableIfNotExists(db); err != nil {
		return fmt.Errorf("sensor status migration failed: %v", err)
	}
	log.Println("Sensor status table created successfully")

	// Run sensor logs migration
	log.Println("Creating sensor logs table...")
	if err := CreateSensorLogsTableIfNotExists(db); err != nil {
		return fmt.Errorf("sensor logs migration failed: %v", err)
	}
	log.Println("Sensor logs table created successfully")

	return nil
}
