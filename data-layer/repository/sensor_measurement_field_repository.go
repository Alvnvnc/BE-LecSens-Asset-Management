package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Helper functions for SQL null types
func NullStringFromPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func NullFloat64FromPtr(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: *f,
		Valid:   true,
	}
}

// SensorMeasurementField represents a sensor measurement field in the database
type SensorMeasurementField struct {
	ID                      uuid.UUID       `json:"id"`
	SensorMeasurementTypeID uuid.UUID       `json:"sensor_measurement_type_id"`
	Name                    string          `json:"name"`
	Label                   string          `json:"label"`
	Description             sql.NullString  `json:"description"`
	DataType                string          `json:"data_type"`
	Required                bool            `json:"required"`
	Unit                    sql.NullString  `json:"unit"`
	Min                     sql.NullFloat64 `json:"min"`
	Max                     sql.NullFloat64 `json:"max"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

// SensorMeasurementFieldRepository handles database operations for sensor measurement fields
type SensorMeasurementFieldRepository struct {
	db *sql.DB
}

// NewSensorMeasurementFieldRepository creates a new instance of SensorMeasurementFieldRepository
func NewSensorMeasurementFieldRepository(db *sql.DB) *SensorMeasurementFieldRepository {
	return &SensorMeasurementFieldRepository{
		db: db,
	}
}

// GetAll retrieves all sensor measurement fields
func (r *SensorMeasurementFieldRepository) GetAll(ctx context.Context) ([]*SensorMeasurementField, error) {
	query := `
		SELECT id, sensor_measurement_type_id, name, label, description, data_type, required, unit, min, max, created_at, updated_at
		FROM sensor_measurement_fields
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []*SensorMeasurementField
	for rows.Next() {
		field := &SensorMeasurementField{}
		err := rows.Scan(
			&field.ID,
			&field.SensorMeasurementTypeID,
			&field.Name,
			&field.Label,
			&field.Description,
			&field.DataType,
			&field.Required,
			&field.Unit,
			&field.Min,
			&field.Max,
			&field.CreatedAt,
			&field.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	return fields, nil
}

// Create creates a new sensor measurement field
func (r *SensorMeasurementFieldRepository) Create(ctx context.Context, field *SensorMeasurementField) (*SensorMeasurementField, error) {
	query := `
		INSERT INTO sensor_measurement_fields (
			id, sensor_measurement_type_id, name, label, description, data_type, required, unit, min, max, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) RETURNING id, sensor_measurement_type_id, name, label, description, data_type, required, unit, min, max, created_at, updated_at
	`

	now := time.Now()
	field.ID = uuid.New()
	field.CreatedAt = now
	field.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		field.ID,
		field.SensorMeasurementTypeID,
		field.Name,
		field.Label,
		field.Description,
		field.DataType,
		field.Required,
		field.Unit,
		field.Min,
		field.Max,
		field.CreatedAt,
		field.UpdatedAt,
	).Scan(
		&field.ID,
		&field.SensorMeasurementTypeID,
		&field.Name,
		&field.Label,
		&field.Description,
		&field.DataType,
		&field.Required,
		&field.Unit,
		&field.Min,
		&field.Max,
		&field.CreatedAt,
		&field.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return field, nil
}

// GetByID retrieves a sensor measurement field by its ID
func (r *SensorMeasurementFieldRepository) GetByID(ctx context.Context, id uuid.UUID) (*SensorMeasurementField, error) {
	query := `
		SELECT id, sensor_measurement_type_id, name, label, description, data_type, required, unit, min, max, created_at, updated_at
		FROM sensor_measurement_fields
		WHERE id = $1
	`

	field := &SensorMeasurementField{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&field.ID,
		&field.SensorMeasurementTypeID,
		&field.Name,
		&field.Label,
		&field.Description,
		&field.DataType,
		&field.Required,
		&field.Unit,
		&field.Min,
		&field.Max,
		&field.CreatedAt,
		&field.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return field, nil
}

// GetByMeasurementTypeID retrieves all fields for a measurement type
func (r *SensorMeasurementFieldRepository) GetByMeasurementTypeID(ctx context.Context, measurementTypeID uuid.UUID) ([]*SensorMeasurementField, error) {
	query := `
		SELECT id, sensor_measurement_type_id, name, label, description, data_type, required, unit, min, max, created_at, updated_at
		FROM sensor_measurement_fields
		WHERE sensor_measurement_type_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, measurementTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []*SensorMeasurementField
	for rows.Next() {
		field := &SensorMeasurementField{}
		err := rows.Scan(
			&field.ID,
			&field.SensorMeasurementTypeID,
			&field.Name,
			&field.Label,
			&field.Description,
			&field.DataType,
			&field.Required,
			&field.Unit,
			&field.Min,
			&field.Max,
			&field.CreatedAt,
			&field.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	return fields, nil
}

// Update updates an existing sensor measurement field
func (r *SensorMeasurementFieldRepository) Update(ctx context.Context, field *SensorMeasurementField) (*SensorMeasurementField, error) {
	query := `
		UPDATE sensor_measurement_fields
		SET name = $1, label = $2, description = $3, data_type = $4, required = $5, unit = $6, min = $7, max = $8, updated_at = $9
		WHERE id = $10
		RETURNING id, sensor_measurement_type_id, name, label, description, data_type, required, unit, min, max, created_at, updated_at
	`

	field.UpdatedAt = time.Now()

	err := r.db.QueryRowContext(ctx, query,
		field.Name,
		field.Label,
		field.Description,
		field.DataType,
		field.Required,
		field.Unit,
		field.Min,
		field.Max,
		field.UpdatedAt,
		field.ID,
	).Scan(
		&field.ID,
		&field.SensorMeasurementTypeID,
		&field.Name,
		&field.Label,
		&field.Description,
		&field.DataType,
		&field.Required,
		&field.Unit,
		&field.Min,
		&field.Max,
		&field.CreatedAt,
		&field.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return field, nil
}

// Delete deletes a sensor measurement field
func (r *SensorMeasurementFieldRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sensor_measurement_fields WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetRequiredFields retrieves all required fields for a measurement type
func (r *SensorMeasurementFieldRepository) GetRequiredFields(ctx context.Context, measurementTypeID uuid.UUID) ([]*SensorMeasurementField, error) {
	query := `
		SELECT id, sensor_measurement_type_id, name, label, description, data_type, required, unit, min, max, created_at, updated_at
		FROM sensor_measurement_fields
		WHERE sensor_measurement_type_id = $1 AND required = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, measurementTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []*SensorMeasurementField
	for rows.Next() {
		field := &SensorMeasurementField{}
		err := rows.Scan(
			&field.ID,
			&field.SensorMeasurementTypeID,
			&field.Name,
			&field.Label,
			&field.Description,
			&field.DataType,
			&field.Required,
			&field.Unit,
			&field.Min,
			&field.Max,
			&field.CreatedAt,
			&field.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	return fields, nil
}

// AddDynamicColumns adds dynamic columns to the iot_sensor_readings table based on measurement fields
func (r *SensorMeasurementFieldRepository) AddDynamicColumns(ctx context.Context, fields []SensorMeasurementField) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if columns exist and add them if they don't
	for _, field := range fields {
		// Check if column exists
		var exists bool
		checkQuery := `
			SELECT EXISTS (
				SELECT 1 
				FROM information_schema.columns 
				WHERE table_name = 'iot_sensor_readings' 
				AND column_name = $1
			)`
		err := tx.QueryRowContext(ctx, checkQuery, field.Name).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if column exists: %w", err)
		}

		if !exists {
			// Determine column type based on data type
			var columnType string
			var defaultValue string
			switch field.DataType {
			case "number":
				// Use DOUBLE PRECISION for all numeric values
				columnType = "DOUBLE PRECISION"
				defaultValue = "NULL"

				// Add constraints if min/max values are specified
				if field.Min.Valid || field.Max.Valid {
					constraints := []string{}
					if field.Min.Valid {
						constraints = append(constraints, fmt.Sprintf("CHECK (%s >= %f)", field.Name, field.Min.Float64))
					}
					if field.Max.Valid {
						constraints = append(constraints, fmt.Sprintf("CHECK (%s <= %f)", field.Name, field.Max.Float64))
					}
					columnType += " " + strings.Join(constraints, " AND ")
				}
			case "string":
				columnType = "TEXT"
				defaultValue = "NULL"
			case "boolean":
				columnType = "BOOLEAN"
				defaultValue = "NULL"
			case "array":
				columnType = "JSONB"
				defaultValue = "'[]'::jsonb"
			case "object":
				columnType = "JSONB"
				defaultValue = "'{}'::jsonb"
			default:
				return fmt.Errorf("unsupported data type: %s", field.DataType)
			}

			// Build metadata comment
			metadata := []string{
				fmt.Sprintf("Measurement field: %s", field.Name),
			}
			if field.Description.Valid {
				metadata = append(metadata, fmt.Sprintf("Description: %s", field.Description.String))
			}
			if field.Unit.Valid {
				metadata = append(metadata, fmt.Sprintf("Unit: %s", field.Unit.String))
			}
			metadata = append(metadata, fmt.Sprintf("Required: %v", field.Required))
			if field.Min.Valid {
				metadata = append(metadata, fmt.Sprintf("Min: %f", field.Min.Float64))
			}
			if field.Max.Valid {
				metadata = append(metadata, fmt.Sprintf("Max: %f", field.Max.Float64))
			}

			// Add column with comment for metadata
			addColumnQuery := fmt.Sprintf(`
				ALTER TABLE iot_sensor_readings 
				ADD COLUMN %s %s DEFAULT %s,
				COMMENT ON COLUMN iot_sensor_readings.%s IS '%s'`,
				field.Name,
				columnType,
				defaultValue,
				field.Name,
				strings.Join(metadata, ", "))

			_, err = tx.ExecContext(ctx, addColumnQuery)
			if err != nil {
				return fmt.Errorf("failed to add column %s: %w", field.Name, err)
			}

			// Add index for numeric columns
			if field.DataType == "number" {
				indexQuery := fmt.Sprintf(`
					CREATE INDEX IF NOT EXISTS idx_iot_readings_%s 
					ON iot_sensor_readings(%s)`,
					field.Name,
					field.Name)
				_, err = tx.ExecContext(ctx, indexQuery)
				if err != nil {
					return fmt.Errorf("failed to create index for column %s: %w", field.Name, err)
				}
			}

			// Add NOT NULL constraint if field is required
			if field.Required {
				notNullQuery := fmt.Sprintf(`
					ALTER TABLE iot_sensor_readings 
					ALTER COLUMN %s SET NOT NULL`,
					field.Name)
				_, err = tx.ExecContext(ctx, notNullQuery)
				if err != nil {
					return fmt.Errorf("failed to add NOT NULL constraint for column %s: %w", field.Name, err)
				}
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
