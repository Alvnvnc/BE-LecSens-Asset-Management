package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// MigrateDatabase creates the database if it doesn't exist and runs all migrations
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

	return nil
}
