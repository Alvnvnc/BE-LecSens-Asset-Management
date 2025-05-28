package main

import (
	"be-lecsens/asset_management/data-layer/config"
	"be-lecsens/asset_management/data-layer/migration"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	var csvPath string
	var dropTable bool
	flag.StringVar(&csvPath, "csv", "", "Path to CSV file")
	flag.BoolVar(&dropTable, "drop", false, "Drop existing locations table before import")
	flag.Parse()

	// If CSV path not provided via flag, check positional arguments
	if csvPath == "" && len(flag.Args()) > 0 {
		csvPath = flag.Args()[0]
	}

	// If still no CSV path, use default
	if csvPath == "" {
		csvPath = filepath.Join("data-layer", "migration", "seeder", "kota_kab.csv")
	}
	log.Printf("Using CSV path: %s", csvPath)
	if dropTable {
		log.Println("Warning: Will drop existing locations table")
	}

	// Check if CSV file exists
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		log.Fatalf("CSV file not found: %s", csvPath)
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

	// Handle table based on drop flag
	if dropTable {
		log.Println("Dropping locations table...")
		err = migration.DropLocationsTable(db)
		if err != nil {
			log.Fatalf("Failed to drop locations table: %v", err)
		}
	}

	// Create table if not exists
	err = migration.CreateLocationsTableIfNotExists(db)
	if err != nil {
		log.Fatalf("Failed to create/verify locations table: %v", err)
	}

	// Import locations from CSV
	err = migration.ImportLocationsFromCSV(db, csvPath)
	if err != nil {
		log.Fatalf("Failed to import locations: %v", err)
	}

	log.Println("Location import completed successfully")
}
