package main

import (
	"be-lecsens/asset_management/data-layer/config"
	"be-lecsens/asset_management/data-layer/migration"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
		action    = flag.String("action", "", "Action to perform: drop-table, truncate-table, drop-all, migrate, seed, cleanup-duplicates")
		tableName = flag.String("table", "", "Table name (for drop-table, truncate-table)")
		csvPath   = flag.String("csv", "", "Path to CSV file (for seed)")
		force     = flag.Bool("force", false, "Force action without confirmation")
		assetID   = flag.String("asset-id", "", "Asset ID (for cleanup-duplicates)")
		dryRun    = flag.Bool("dry-run", false, "Show what would be deleted without actually deleting (for cleanup-duplicates)")
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
		runSeeder(db, *csvPath)
	case "cleanup-duplicates":
		runCleanupDuplicates(*assetID, *dryRun)
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
	fmt.Println("  cleanup-duplicates  Clean up duplicate asset documents")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -table=<name>  Table name (required for drop-table, truncate-table)")
	fmt.Println("  -csv=<path>    CSV file path (for seed action)")
	fmt.Println("  -force         Skip confirmation prompts")
	fmt.Println("  -asset-id=<id> Asset ID (for cleanup-duplicates, optional)")
	fmt.Println("  -dry-run       Show what would be deleted without actually deleting")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run helpers/cmd/cmd.go -action=drop-table -table=assets")
	fmt.Println("  go run helpers/cmd/cmd.go -action=truncate-table -table=locations")
	fmt.Println("  go run helpers/cmd/cmd.go -action=drop-all -force")
	fmt.Println("  go run helpers/cmd/cmd.go -action=migrate")
	fmt.Println("  go run helpers/cmd/cmd.go -action=seed -csv=data-layer/migration/seeder/kota_kab.csv")
	fmt.Println("  go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run")
	fmt.Println("  go run helpers/cmd/cmd.go -action=cleanup-duplicates")
	fmt.Println("  go run helpers/cmd/cmd.go -action=cleanup-duplicates -asset-id=123e4567-e89b-12d3-a456-426614174000")
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

func runSeeder(db *sql.DB, csvPath string) {
	log.Println("Running location seeder...")

	// If CSV path not provided, use default
	if csvPath == "" {
		csvPath = filepath.Join("data-layer", "migration", "seeder", "kota_kab.csv")
	}

	// Check if CSV file exists
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		log.Fatalf("CSV file not found: %s", csvPath)
	}

	// Create locations table if not exists
	err := migration.CreateLocationsTableIfNotExists(db)
	if err != nil {
		log.Fatalf("Failed to create locations table: %v", err)
	}

	// Import locations from CSV
	err = migration.ImportLocationsFromCSV(db, csvPath)
	if err != nil {
		log.Fatalf("Failed to import locations: %v", err)
	}

	log.Println("Location seeder completed successfully")
}

// runCleanupDuplicates runs the cleanup duplicates functionality
func runCleanupDuplicates(assetID string, dryRun bool) {
	// Create cleanup command without passing db (it creates its own connection)
	cleanupCmd, err := NewCleanupDuplicatesCommand()
	if err != nil {
		log.Fatalf("Failed to create cleanup command: %v", err)
	}
	defer cleanupCmd.Close()

	ctx := context.Background()

	if assetID != "" {
		// Cleanup specific asset
		log.Printf("Starting cleanup for asset: %s", assetID)
		err = cleanupCmd.CleanupAssetDuplicates(ctx, assetID, dryRun)
	} else {
		// Cleanup all assets
		log.Println("Starting cleanup for all assets")
		err = cleanupCmd.CleanupAllDuplicates(ctx, dryRun)
	}

	if err != nil {
		log.Fatalf("Cleanup failed: %v", err)
	}
}
