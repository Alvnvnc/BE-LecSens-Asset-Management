package migration

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// CreateLocationsTableIfNotExists creates the locations table if it doesn't exist
func CreateLocationsTableIfNotExists(db *sql.DB) error {
	// Check if table exists
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'locations'
		)
	`).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check if locations table exists: %v", err)
	}

	// If table already exists, no need to create it
	if exists {
		log.Println("Locations table already exists")
		return nil
	}

	// Create new locations table aligned with the Location entity
	_, err = db.Exec(`
		CREATE TABLE locations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			region_code VARCHAR(50),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			address TEXT,
			longitude FLOAT,
			latitude FLOAT,
			hierarchy_level INTEGER NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)

	if err != nil {
		return fmt.Errorf("failed to create locations table: %v", err)
	}

	log.Println("Locations table created successfully")
	return nil
}

// DropLocationsTable drops the existing locations table - use with caution!
func DropLocationsTable(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS locations CASCADE`)
	if err != nil {
		return fmt.Errorf("failed to drop locations table: %v", err)
	}
	log.Println("Locations table dropped successfully")
	return nil
}

// ImportLocationsFromCSV imports location data from a CSV file
func ImportLocationsFromCSV(db *sql.DB, csvPath string) error {
	log.Printf("Importing locations from CSV file: %s", csvPath)

	// Open the CSV file
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %v", err)
	}

	// Prepare statement for inserting locations
	stmt, err := db.Prepare(`
		INSERT INTO locations (
			region_code, name, longitude, latitude, hierarchy_level, is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Insert locations
	successCount := 0
	for i, record := range records {
		if i == 0 { // Skip header row
			continue
		}

		// Check if record has enough fields
		if len(record) < 6 {
			log.Printf("Skipping row %d: insufficient fields", i)
			continue
		}

		// CSV columns: [index, id, foreign, name, lat, long]
		regionCode := record[1] // id column as region_code
		name := record[3]       // name column

		// Parse coordinates, handling "null" values
		var latitude, longitude float64
		if record[4] != "null" {
			latitude, _ = strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
		}
		if record[5] != "null" {
			longitude, _ = strconv.ParseFloat(strings.TrimSpace(record[5]), 64)
		}

		// Set hierarchy level based on the name prefix
		hierarchyLevel := 2 // Default for KABUPATEN
		if strings.HasPrefix(name, "KOTA") {
			hierarchyLevel = 3
		}

		// Set location as active by default
		isActive := true

		// Insert location - make sure longitude and latitude are in the correct order
		var id string
		err := stmt.QueryRow(regionCode, name, longitude, latitude, hierarchyLevel, isActive).Scan(&id)
		if err != nil {
			log.Printf("Error inserting location: %v", err)
			continue
		}

		successCount++
	}

	log.Printf("Successfully imported %d locations", successCount)
	return nil
}

// Other existing functions...
