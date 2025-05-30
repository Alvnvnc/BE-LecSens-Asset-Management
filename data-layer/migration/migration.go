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
	// First, connect to postgres directly to create the database if needed
	postgresConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password)

	db, err := sql.Open("postgres", postgresConn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	// Check if the database exists
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '%s')", cfg.DB.Name)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %v", err)
	}

	// Create the database if it doesn't exist
	if !exists {
		log.Printf("Creating database '%s'...", cfg.DB.Name)
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DB.Name))
		if err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
		log.Printf("Database '%s' created successfully", cfg.DB.Name)
	} else {
		log.Printf("Database '%s' already exists", cfg.DB.Name)
	}

	// Now connect to the actual database to create tables
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	appDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to application database: %v", err)
	}

	// Create locations table if not exists
	err = CreateLocationsTableIfNotExists(appDB)
	if err != nil {
		return fmt.Errorf("failed to create locations table: %v", err)
	}
	defer appDB.Close()

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

	return nil
}
