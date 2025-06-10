package main

import (
	"be-lecsens/asset_management/data-layer/config"
	"be-lecsens/asset_management/data-layer/migration"
	"be-lecsens/asset_management/data-layer/migration/seeder"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file. Using environment variables.")
	}

	// Set environment variable for development if not set
	if os.Getenv("ENVIRONMENT") == "" {
		os.Setenv("ENVIRONMENT", "development")
		log.Println("ENVIRONMENT variable not set. Setting to 'development'")
	}

	// Parse command line flags
	var (
		action      = flag.String("action", "", "Action to perform: drop-table, truncate-table, drop-all, migrate, seed, seed-all")
		tableName   = flag.String("table", "", "Table name (for drop-table, truncate-table)")
		csvPath     = flag.String("csv", "", "Path to CSV file (for seed)")
		seederType  = flag.String("seeder", "", "Seeder type: location, asset-type, sensor-type, measurement-type, asset, asset-sensor, measurement-field, threshold, reading, alert, sensor-status, sensor-logs, or all")
		force       = flag.Bool("force", false, "Force action without confirmation")
		days        = flag.Int("days", 7, "Number of days of historical data to generate (for reading seeder)")
		forceReseed = flag.Bool("force-reseed", false, "Force re-seed even if data exists")
	)
	flag.Parse()

	if *action == "" {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize validator
	validator := seeder.NewValidator(db)

	switch *action {
	case "drop-table":
		if *tableName == "" {
			log.Fatal("Table name is required for drop-table action")
		}
		dropTable(db, *tableName, *force)
	case "truncate-table":
		if *tableName == "" {
			log.Fatal("Table name is required for truncate-table action")
		}
		truncateTable(db, *tableName, *force)
	case "drop-all":
		dropAllTables(db, *force)
	case "migrate":
		runMigrations(cfg)
	case "seed":
		if *seederType == "" {
			log.Fatal("Seeder type is required for seed action. Use -seeder flag")
		}
		runSpecificSeeder(db, validator, *seederType, *csvPath, *days, *forceReseed)
	case "seed-all":
		runAllSeeders(db, validator, *days, *forceReseed)
	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

func printUsage() {
	fmt.Println("Database Management Tool")
	fmt.Println("Usage:")
	fmt.Println("  go run helpers/cmd/cmd.go -action=<action> [options]")
	fmt.Println("")
	fmt.Println("Actions:")
	fmt.Println("  drop-table     Drop a specific table")
	fmt.Println("  truncate-table Truncate (empty) a specific table")
	fmt.Println("  drop-all       Drop all tables")
	fmt.Println("  migrate        Run all migrations")
	fmt.Println("  seed           Run location seeder")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -table=<name>  Table name (required for drop-table, truncate-table)")
	fmt.Println("  -csv=<path>    CSV file path (for seed action)")
	fmt.Println("  -force         Skip confirmation prompts")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run helpers/cmd/cmd.go -action=drop-table -table=assets")
	fmt.Println("  go run helpers/cmd/cmd.go -action=truncate-table -table=locations")
	fmt.Println("  go run helpers/cmd/cmd.go -action=drop-all -force")
	fmt.Println("  go run helpers/cmd/cmd.go -action=migrate")
	fmt.Println("  go run helpers/cmd/cmd.go -action=seed -csv=data-layer/migration/seeder/kota_kab.csv")
}

func dropTable(db *sql.DB, tableName string, force bool) {
	if !force {
		fmt.Printf("Are you sure you want to drop table '%s'? This action cannot be undone. (y/N): ", tableName)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			log.Println("Operation cancelled")
			return
		}
	}

	query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to drop table %s: %v", tableName, err)
	}

	log.Printf("Table '%s' dropped successfully", tableName)
}

func truncateTable(db *sql.DB, tableName string, force bool) {
	if !force {
		fmt.Printf("Are you sure you want to truncate table '%s'? This will delete all data. (y/N): ", tableName)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			log.Println("Operation cancelled")
			return
		}
	}

	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to truncate table %s: %v", tableName, err)
	}

	log.Printf("Table '%s' truncated successfully", tableName)
}

func dropAllTables(db *sql.DB, force bool) {
	if !force {
		fmt.Print("Are you sure you want to drop ALL tables? This action cannot be undone. (y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			log.Println("Operation cancelled")
			return
		}
	}

	// Get all table names
	rows, err := db.Query(`
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public'
	`)
	if err != nil {
		log.Fatalf("Failed to get table names: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Failed to scan table name: %v", err)
		}
		tables = append(tables, tableName)
	}

	if len(tables) == 0 {
		log.Println("No tables found to drop")
		return
	}

	// Drop all tables
	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Failed to drop table %s: %v", table, err)
		} else {
			log.Printf("Dropped table: %s", table)
		}
	}

	log.Println("All tables dropped successfully")
}

func runMigrations(cfg *config.Config) {
	log.Println("Running migrations...")

	err := migration.MigrateDatabase(cfg)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}

// runSpecificSeeder runs a specific seeder based on the seeder type
func runSpecificSeeder(db *sql.DB, validator *seeder.Validator, seederType, csvPath string, days int, forceReseed bool) {
	log.Printf("Running %s seeder...", seederType)

	// Validate tables and foreign keys
	if err := validator.ValidateAllTables(); err != nil {
		log.Fatalf("Table validation failed: %v", err)
	}

	if err := validator.ValidateAllForeignKeys(); err != nil {
		log.Fatalf("Foreign key validation failed: %v", err)
	}

	// Skip data existence check for "all" seeder type
	if seederType != "all" {
		// Check if data exists
		tableName := getTableName(seederType)
		if tableName == "" {
			log.Fatalf("Unknown seeder type: %s", seederType)
		}

		exists, err := validator.ValidateDataExists(tableName)
		if err != nil {
			log.Fatalf("Failed to check if data exists: %v", err)
		}

		if exists && !forceReseed {
			log.Printf("Data already exists in %s. Use -force-reseed to override", tableName)
			return
		}
	}

	// Cleanup if force reseed
	if forceReseed {
		if err := validator.CleanupBeforeSeeding(); err != nil {
			log.Fatalf("Failed to cleanup before seeding: %v", err)
		}
	}

	var seedingErr error
	switch seederType {
	case "location":
		seedingErr = runLocationSeeder(db, csvPath)
	case "asset-type":
		seedingErr = runAssetTypeSeeder(db)
	case "sensor-type":
		seedingErr = runSensorTypeSeeder(db)
	case "measurement-type":
		seedingErr = runMeasurementTypeSeeder(db)
	case "asset":
		seedingErr = runAssetSeeder(db)
	case "asset-sensor":
		seedingErr = runAssetSensorSeeder(db)
	case "measurement-field":
		seedingErr = runMeasurementFieldSeeder(db)
	case "threshold":
		seedingErr = runThresholdSeeder(db)
	case "reading":
		seedingErr = runReadingSeeder(db, days)
	case "alert":
		seedingErr = runAlertSeeder()
	case "sensor-status":
		seedingErr = runSensorStatusSeeder(db)
	case "sensor-logs":
		seedingErr = runSensorLogsSeeder(db)
	case "all":
		seedingErr = runAllSeeders(db, validator, days, forceReseed)
	default:
		log.Fatalf("Unknown seeder type: %s", seederType)
	}

	if seedingErr != nil {
		log.Fatalf("Failed to run %s seeder: %v", seederType, seedingErr)
	}

	log.Printf("%s seeder completed successfully", seederType)
}

// runAllSeeders runs all seeders in the proper dependency order
func runAllSeeders(db *sql.DB, validator *seeder.Validator, days int, forceReseed bool) error {
	log.Println("Running all seeders in dependency order...")

	// Validate tables and foreign keys
	if err := validator.ValidateAllTables(); err != nil {
		return fmt.Errorf("table validation failed: %v", err)
	}

	if err := validator.ValidateAllForeignKeys(); err != nil {
		return fmt.Errorf("foreign key validation failed: %v", err)
	}

	// Cleanup if force reseed
	if forceReseed {
		if err := validator.CleanupBeforeSeeding(); err != nil {
			return fmt.Errorf("failed to cleanup before seeding: %v", err)
		}
	}

	seeders := []struct {
		name      string
		fn        func() error
		waitAfter time.Duration
	}{
		{"Location", func() error {
			locationSeeder := seeder.NewLocationSeeder(db)
			return locationSeeder.Seed(context.Background())
		}, 2 * time.Second},
		{"Asset Type", func() error { return runAssetTypeSeeder(db) }, 2 * time.Second},
		{"Sensor Type", func() error { return runSensorTypeSeeder(db) }, 2 * time.Second},
		{"Measurement Type", func() error { return runMeasurementTypeSeeder(db) }, 2 * time.Second},
		{"Asset", func() error { return runAssetSeeder(db) }, 3 * time.Second},
		{"Asset Sensor", func() error { return runAssetSensorSeeder(db) }, 5 * time.Second},
		{"Measurement Field", func() error { return runMeasurementFieldSeeder(db) }, 2 * time.Second},
		{"Threshold", func() error { return runThresholdSeeder(db) }, 2 * time.Second},
		{"IoT Reading", func() error { return runReadingSeeder(db, days) }, 2 * time.Second},
		{"Sensor Status", func() error { return runSensorStatusSeeder(db) }, 2 * time.Second},
		{"Sensor Logs", func() error { return runSensorLogsSeeder(db) }, 2 * time.Second},
		{"Asset Alert", func() error { return runAlertSeeder() }, 0},
	}

	for _, s := range seeders {
		log.Printf("Running %s seeder...", s.name)
		if err := s.fn(); err != nil {
			return fmt.Errorf("failed to run %s seeder: %v", s.name, err)
		}
		log.Printf("%s seeder completed successfully", s.name)

		if s.waitAfter > 0 {
			log.Printf("Waiting %v for database operations to commit...", s.waitAfter)
			time.Sleep(s.waitAfter)
		}
	}

	log.Println("All seeders completed successfully!")
	return nil
}

// getTableName returns the table name for a given seeder type
func getTableName(seederType string) string {
	switch seederType {
	case "location":
		return "locations"
	case "asset-type":
		return "asset_types"
	case "sensor-type":
		return "sensor_types"
	case "measurement-type":
		return "sensor_measurement_types"
	case "asset":
		return "assets"
	case "asset-sensor":
		return "asset_sensors"
	case "measurement-field":
		return "sensor_measurement_fields"
	case "threshold":
		return "sensor_thresholds"
	case "reading":
		return "iot_sensor_readings"
	case "alert":
		return "asset_alerts"
	case "sensor-status":
		return "sensor_status"
	case "sensor-logs":
		return "sensor_logs"
	default:
		return ""
	}
}

// Individual seeder functions
func runLocationSeeder(db *sql.DB, csvPath string) error {
	// Always use predefined seeder to ensure consistent UUIDs with asset seeder constants
	// This prevents foreign key constraint violations in Docker environment
	log.Println("Using predefined location seeder for consistent asset references...")

	locationSeeder := seeder.NewLocationSeeder(db)
	return locationSeeder.Seed(context.Background())
}

func runAssetTypeSeeder(db *sql.DB) error {
	assetTypeSeeder := seeder.NewAssetTypeSeeder(db)
	return assetTypeSeeder.Seed(context.Background())
}

func runSensorTypeSeeder(db *sql.DB) error {
	sensorTypeSeeder := seeder.NewSensorTypeSeeder(db)
	return sensorTypeSeeder.Seed(context.Background())
}

func runMeasurementTypeSeeder(db *sql.DB) error {
	measurementTypeSeeder := seeder.NewSensorMeasurementTypeSeeder(db)
	return measurementTypeSeeder.Seed(context.Background())
}

func runAssetSeeder(db *sql.DB) error {
	assetSeeder := seeder.NewAssetSeeder(db)
	return assetSeeder.Seed(context.Background())
}

func runAssetSensorSeeder(db *sql.DB) error {
	assetSensorSeeder := seeder.NewAssetSensorSeeder(db)
	return assetSensorSeeder.Seed(context.Background())
}

func runMeasurementFieldSeeder(db *sql.DB) error {
	measurementFieldSeeder := seeder.NewSensorMeasurementFieldSeeder(db)
	return measurementFieldSeeder.Seed(context.Background())
}

func runThresholdSeeder(db *sql.DB) error {
	thresholdSeeder := seeder.NewSensorThresholdSeeder(db)
	return thresholdSeeder.SeedThresholds(context.Background())
}

func runReadingSeeder(db *sql.DB, days int) error {
	readingSeeder := seeder.NewIoTSensorReadingSeeder(db)
	// Note: days parameter is ignored since seeder now generates from Jan 1, 2022 to present
	return readingSeeder.SeedReadings(context.Background(), true)
}

func runAlertSeeder() error {
	// TODO: Implement AssetAlertSeeder when it's available
	// alertSeeder := seeder.NewAssetAlertSeeder(db)
	// return alertSeeder.Seed(context.Background())
	log.Println("AssetAlertSeeder not implemented yet - skipping")
	return nil
}

func runSensorStatusSeeder(db *sql.DB) error {
	sensorStatusSeeder := seeder.NewSensorStatusSeeder(db)
	return sensorStatusSeeder.Seed(context.Background())
}

func runSensorLogsSeeder(db *sql.DB) error {
	sensorLogsSeeder := seeder.NewSensorLogsSeeder(db)
	return sensorLogsSeeder.Seed()
}
