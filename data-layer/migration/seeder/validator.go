package seeder

import (
	"database/sql"
	"fmt"
	"log"
)

// Validator handles all validation for seeder data
type Validator struct {
	db *sql.DB
}

// NewValidator creates a new validator instance
func NewValidator(db *sql.DB) *Validator {
	return &Validator{db: db}
}

// ValidateTableExists checks if a table exists in the database
func (v *Validator) ValidateTableExists(tableName string) error {
	var exists bool
	query := `SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = $1
	)`
	err := v.db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if table exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	return nil
}

// ValidateForeignKey checks if a foreign key constraint exists
func (v *Validator) ValidateForeignKey(tableName, columnName, referencedTable string) error {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
			ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.constraint_column_usage ccu
			ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY'
		AND tc.table_name = $1
		AND kcu.column_name = $2
		AND ccu.table_name = $3
	)`
	err := v.db.QueryRow(query, tableName, columnName, referencedTable).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check foreign key: %v", err)
	}
	if !exists {
		return fmt.Errorf("foreign key constraint not found: %s.%s references %s", tableName, columnName, referencedTable)
	}
	return nil
}

// ValidateDataExists checks if data already exists in a table
func (v *Validator) ValidateDataExists(tableName string) (bool, error) {
	var count int
	err := v.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ValidateAllTables validates all required tables exist
func (v *Validator) ValidateAllTables() error {
	requiredTables := []string{
		"locations",
		"asset_types",
		"sensor_types",
		"sensor_measurement_types",
		"assets",
		"asset_sensors",
		"sensor_measurement_fields",
		"sensor_thresholds",
		"sensor_status",
		"sensor_logs",
		"asset_alerts",
	}

	for _, table := range requiredTables {
		if err := v.ValidateTableExists(table); err != nil {
			return fmt.Errorf("table validation failed for %s: %v", table, err)
		}
		log.Printf("Table %s exists", table)
	}
	return nil
}

// ValidateAllForeignKeys validates all required foreign key constraints
func (v *Validator) ValidateAllForeignKeys() error {
	foreignKeys := []struct {
		table, column, referencedTable string
	}{
		{"assets", "location_id", "locations"},
		{"assets", "asset_type_id", "asset_types"},
		{"asset_sensors", "asset_id", "assets"},
		{"asset_sensors", "sensor_type_id", "sensor_types"},
		{"sensor_measurement_fields", "sensor_measurement_type_id", "sensor_measurement_types"},
		{"sensor_thresholds", "asset_sensor_id", "asset_sensors"},
		{"sensor_status", "asset_sensor_id", "asset_sensors"},
		{"sensor_logs", "asset_sensor_id", "asset_sensors"},
		{"asset_alerts", "asset_id", "assets"},
	}

	for _, fk := range foreignKeys {
		if err := v.ValidateForeignKey(fk.table, fk.column, fk.referencedTable); err != nil {
			return fmt.Errorf("foreign key validation failed for %s.%s: %v", fk.table, fk.column, err)
		}
		log.Printf("Foreign key %s.%s -> %s exists", fk.table, fk.column, fk.referencedTable)
	}
	return nil
}

// CleanupBeforeSeeding removes all existing data
func (v *Validator) CleanupBeforeSeeding() error {
	// Check if cleanup has already been done
	var count int
	err := v.db.QueryRow("SELECT COUNT(*) FROM locations").Scan(&count)
	if err == nil && count == 0 {
		log.Println("Database is already empty, skipping cleanup")
		return nil
	}

	tables := []string{
		"iot_sensor_readings", // Add IoT readings table
		"asset_alerts",
		"sensor_logs",
		"sensor_status",
		"sensor_thresholds",
		"sensor_measurement_fields",
		"asset_sensors",
		"assets",
		"sensor_measurement_types",
		"sensor_types",
		"asset_types",
		"locations",
	}

	for _, table := range tables {
		if _, err := v.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
			return fmt.Errorf("failed to truncate %s: %v", table, err)
		}
		log.Printf("Truncated table %s", table)
	}
	return nil
}
