package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

// MigrateDatabase creates database if it doesn't exist and runs all migrations
func MigrateDatabase(cfg *config.Config) error {
	// Try to connect to the target database directly first
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	appDB, err := sql.Open("postgres", dsn)
	if err == nil {
		err = appDB.Ping()
		if err == nil {
			log.Printf("Database '%s' already exists and is accessible", cfg.DB.Name)
			// Database exists, proceed with migrations
			defer appDB.Close()
			return runMigrationsWithDB(appDB, cfg)
		}
		appDB.Close()
	}

	// If target database doesn't exist, try to create it
	// Try connecting to different admin databases in order of preference
	adminDatabases := []string{cfg.DB.Name, "postgres", "template1"}
	var adminDB *sql.DB

	for _, dbName := range adminDatabases {
		adminConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, dbName)

		adminDB, err = sql.Open("postgres", adminConn)
		if err != nil {
			continue
		}

		err = adminDB.Ping()
		if err == nil {
			log.Printf("Connected to admin database '%s' for database management", dbName)
			break
		}
		adminDB.Close()
		adminDB = nil
	}

	if adminDB == nil {
		return fmt.Errorf("failed to connect to any admin database (%v) to manage target database '%s'", adminDatabases, cfg.DB.Name)
	}
	defer adminDB.Close()

	// Check if the target database exists (only if we're not already connected to it)
	if adminDatabases[0] != cfg.DB.Name {
		var exists bool
		query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '%s')", cfg.DB.Name)
		err = adminDB.QueryRow(query).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if database exists: %v", err)
		}

		// Create the database if it doesn't exist
		if !exists {
			log.Printf("Creating database '%s'...", cfg.DB.Name)
			_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DB.Name))
			if err != nil {
				return fmt.Errorf("failed to create database: %v", err)
			}
			log.Printf("Database '%s' created successfully", cfg.DB.Name)
		} else {
			log.Printf("Database '%s' already exists", cfg.DB.Name)
		}

		// Now connect to the actual database to create tables
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

		appDB, err = sql.Open("postgres", dsn)
		if err != nil {
			return fmt.Errorf("failed to connect to application database: %v", err)
		}
		defer appDB.Close()
	} else {
		// We're already connected to the target database
		appDB = adminDB
	}

	return runMigrationsWithDB(appDB, cfg)
}

// runMigrationsWithDB runs migrations on the provided database connection
func runMigrationsWithDB(appDB *sql.DB, cfg *config.Config) error {
	// Create locations table if not exists
	err := CreateLocationsTableIfNotExists(appDB)
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

// dropAllTables drops all tables in the correct order to handle dependencies
func dropAllTables(db *sql.DB) error {
	log.Println("Dropping all tables due to schema conflicts...")

	// Drop tables in reverse dependency order
	tablesToDrop := []string{
		"asset_alerts",
		"sensor_thresholds",
		"iot_sensor_readings",
		"asset_sensors",
		"sensor_measurement_fields",
		"sensor_measurement_types",
		"sensor_types",
		"asset_documents",
		"assets",
		"asset_types",
		"locations",
	}

	for _, table := range tablesToDrop {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			log.Printf("Warning: failed to drop table %s: %v", table, err)
		}
	}

	log.Println("All tables dropped successfully")
	return nil
}

// runMigrations executes all database migrations in order
func runMigrations(cfg *config.Config) error {
	log.Println("Running database migrations...")

	// Connect to database for all migrations
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database for migrations: %v", err)
	}
	defer db.Close()

	// Check if we need to drop tables due to schema conflicts
	// We'll attempt to run migrations first, and if we encounter schema errors, drop and recreate
	migrationAttempt := func() error {
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

		// Create locations table if not exists
		err := CreateLocationsTableIfNotExists(db)
		if err != nil {
			return fmt.Errorf("failed to create locations table: %v", err)
		}

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

		// Run IoT sensor reading migration with flexible measurement support
		log.Println("Creating IoT sensor readings table with flexible measurement support...")
		if err := CreateIoTSensorReadingTableDirect(db); err != nil {
			return fmt.Errorf("iot sensor reading migration failed: %v", err)
		}
		log.Println("IoT sensor readings table created successfully")

		// Run sensor threshold migration
		log.Println("Creating sensor thresholds table...")
		if err := CreateSensorThresholdTable(db); err != nil {
			return fmt.Errorf("sensor threshold migration failed: %v", err)
		}
		log.Println("Sensor thresholds table created successfully")

		// Run asset alert migration
		log.Println("Creating asset alerts table...")
		if err := CreateAssetAlertTable(db); err != nil {
			return fmt.Errorf("asset alert migration failed: %v", err)
		}
		log.Println("Asset alerts table created successfully")

		return nil
	}

	// First attempt at migration
	err = migrationAttempt()
	if err != nil {
		// Check if the error is related to missing columns or schema conflicts
		if containsSchemaError(err.Error()) {
			log.Printf("Schema conflict detected: %v", err)
			log.Println("Dropping all tables and recreating...")

			// Drop all tables
			if dropErr := dropAllTables(db); dropErr != nil {
				return fmt.Errorf("failed to drop tables: %v", dropErr)
			}

			// Retry migration after dropping tables
			log.Println("Retrying migrations after dropping tables...")
			err = migrationAttempt()
			if err != nil {
				return fmt.Errorf("migration failed even after dropping tables: %v", err)
			}
		} else {
			return err
		}
	}

	return nil
}

// containsSchemaError checks if the error message indicates a schema conflict
func containsSchemaError(errMsg string) bool {
	schemaErrors := []string{
		"does not exist",
		"already exists",
		"column",
		"relation",
		"constraint",
		"foreign key",
	}

	for _, schemaErr := range schemaErrors {
		if strings.Contains(strings.ToLower(errMsg), strings.ToLower(schemaErr)) {
			return true
		}
	}
	return false
}
