package repository

import (
	"context"
	"database/sql"
	"fmt"
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

	// Add debug logging
	fmt.Printf("DEBUG: Querying sensor_measurement_fields with ID: %s\n", id.String())

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
			fmt.Printf("DEBUG: No rows found for ID: %s\n", id.String())
			return nil, nil
		}
		fmt.Printf("DEBUG: Database error for ID %s: %v\n", id.String(), err)
		return nil, err
	}

	fmt.Printf("DEBUG: Found field with ID: %s, Name: %s\n", field.ID.String(), field.Name)
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
