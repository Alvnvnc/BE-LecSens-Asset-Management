package repository

import (
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// SensorMeasurementTypeRepository defines the interface for sensor measurement type operations
type SensorMeasurementTypeRepository interface {
	Create(ctx context.Context, sensorMeasurementType *dto.SensorMeasurementTypeDTO) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.SensorMeasurementTypeDTO, error)
	List(ctx context.Context, offset, limit int) ([]dto.SensorMeasurementTypeDTO, int, error)
	Update(ctx context.Context, sensorMeasurementType *dto.SensorMeasurementTypeDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetActive(ctx context.Context) ([]dto.SensorMeasurementTypeDTO, error)
	GetBySensorTypeID(ctx context.Context, sensorTypeID uuid.UUID) ([]dto.SensorMeasurementTypeDTO, error)
}

// sensorMeasurementTypeRepository implements the repository interface for sensor measurement types
type sensorMeasurementTypeRepository struct {
	db *sql.DB
}

// NewSensorMeasurementTypeRepository creates a new instance of SensorMeasurementTypeRepository
func NewSensorMeasurementTypeRepository(db *sql.DB) SensorMeasurementTypeRepository {
	return &sensorMeasurementTypeRepository{
		db: db,
	}
}

// Create creates a new sensor measurement type
func (r *sensorMeasurementTypeRepository) Create(ctx context.Context, sensorMeasurementType *dto.SensorMeasurementTypeDTO) error {
	query := `
		INSERT INTO sensor_measurement_types (
			id, sensor_type_id, name, description, properties_schema,
			ui_configuration, version, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	propertiesSchema, err := json.Marshal(sensorMeasurementType.PropertiesSchema)
	if err != nil {
		return fmt.Errorf("error marshaling properties schema: %v", err)
	}

	uiConfig, err := json.Marshal(sensorMeasurementType.UIConfiguration)
	if err != nil {
		return fmt.Errorf("error marshaling UI configuration: %v", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		sensorMeasurementType.ID,
		sensorMeasurementType.SensorTypeID,
		sensorMeasurementType.Name,
		sensorMeasurementType.Description,
		propertiesSchema,
		uiConfig,
		sensorMeasurementType.Version,
		sensorMeasurementType.IsActive,
		sensorMeasurementType.CreatedAt,
		sensorMeasurementType.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating sensor measurement type: %v", err)
	}

	return nil
}

// GetByID retrieves a sensor measurement type by its ID
func (r *sensorMeasurementTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*dto.SensorMeasurementTypeDTO, error) {
	query := `
		SELECT id, sensor_type_id, name, description, properties_schema,
			ui_configuration, version, is_active, created_at, updated_at
		FROM sensor_measurement_types
		WHERE id = $1
	`

	var sensorMeasurementType dto.SensorMeasurementTypeDTO
	var propertiesSchema, uiConfig []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sensorMeasurementType.ID,
		&sensorMeasurementType.SensorTypeID,
		&sensorMeasurementType.Name,
		&sensorMeasurementType.Description,
		&propertiesSchema,
		&uiConfig,
		&sensorMeasurementType.Version,
		&sensorMeasurementType.IsActive,
		&sensorMeasurementType.CreatedAt,
		&sensorMeasurementType.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sensor measurement type not found")
		}
		return nil, fmt.Errorf("error getting sensor measurement type: %v", err)
	}

	// Unmarshal JSON fields
	if len(propertiesSchema) > 0 {
		if err := json.Unmarshal(propertiesSchema, &sensorMeasurementType.PropertiesSchema); err != nil {
			return nil, fmt.Errorf("error unmarshaling properties schema: %v", err)
		}
	}

	if len(uiConfig) > 0 {
		if err := json.Unmarshal(uiConfig, &sensorMeasurementType.UIConfiguration); err != nil {
			return nil, fmt.Errorf("error unmarshaling UI configuration: %v", err)
		}
	}

	return &sensorMeasurementType, nil
}

// List retrieves a paginated list of sensor measurement types
func (r *sensorMeasurementTypeRepository) List(ctx context.Context, offset, limit int) ([]dto.SensorMeasurementTypeDTO, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM sensor_measurement_types`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Get paginated data
	query := `
		SELECT id, sensor_type_id, name, description, properties_schema,
			ui_configuration, version, is_active, created_at, updated_at
		FROM sensor_measurement_types
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying sensor measurement types: %v", err)
	}
	defer rows.Close()

	var sensorMeasurementTypes []dto.SensorMeasurementTypeDTO
	for rows.Next() {
		var sensorMeasurementType dto.SensorMeasurementTypeDTO
		var propertiesSchema, uiConfig []byte

		err := rows.Scan(
			&sensorMeasurementType.ID,
			&sensorMeasurementType.SensorTypeID,
			&sensorMeasurementType.Name,
			&sensorMeasurementType.Description,
			&propertiesSchema,
			&uiConfig,
			&sensorMeasurementType.Version,
			&sensorMeasurementType.IsActive,
			&sensorMeasurementType.CreatedAt,
			&sensorMeasurementType.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("error scanning sensor measurement type: %v", err)
		}

		// Unmarshal JSON fields
		if len(propertiesSchema) > 0 {
			if err := json.Unmarshal(propertiesSchema, &sensorMeasurementType.PropertiesSchema); err != nil {
				return nil, 0, fmt.Errorf("error unmarshaling properties schema: %v", err)
			}
		}

		if len(uiConfig) > 0 {
			if err := json.Unmarshal(uiConfig, &sensorMeasurementType.UIConfiguration); err != nil {
				return nil, 0, fmt.Errorf("error unmarshaling UI configuration: %v", err)
			}
		}

		sensorMeasurementTypes = append(sensorMeasurementTypes, sensorMeasurementType)
	}

	return sensorMeasurementTypes, total, nil
}

// Update updates an existing sensor measurement type
func (r *sensorMeasurementTypeRepository) Update(ctx context.Context, sensorMeasurementType *dto.SensorMeasurementTypeDTO) error {
	query := `
		UPDATE sensor_measurement_types
		SET name = $1, description = $2, properties_schema = $3,
			ui_configuration = $4, version = $5, is_active = $6, updated_at = $7
		WHERE id = $8
	`

	propertiesSchema, err := json.Marshal(sensorMeasurementType.PropertiesSchema)
	if err != nil {
		return fmt.Errorf("error marshaling properties schema: %v", err)
	}

	uiConfig, err := json.Marshal(sensorMeasurementType.UIConfiguration)
	if err != nil {
		return fmt.Errorf("error marshaling UI configuration: %v", err)
	}

	result, err := r.db.ExecContext(ctx, query,
		sensorMeasurementType.Name,
		sensorMeasurementType.Description,
		propertiesSchema,
		uiConfig,
		sensorMeasurementType.Version,
		sensorMeasurementType.IsActive,
		sensorMeasurementType.UpdatedAt,
		sensorMeasurementType.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating sensor measurement type: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sensor measurement type not found")
	}

	return nil
}

// Delete deletes a sensor measurement type by its ID
func (r *sensorMeasurementTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sensor_measurement_types WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting sensor measurement type: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sensor measurement type not found")
	}

	return nil
}

// GetActive retrieves all active sensor measurement types
func (r *sensorMeasurementTypeRepository) GetActive(ctx context.Context) ([]dto.SensorMeasurementTypeDTO, error) {
	query := `
		SELECT id, sensor_type_id, name, description, properties_schema,
			ui_configuration, version, is_active, created_at, updated_at
		FROM sensor_measurement_types
		WHERE is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying active sensor measurement types: %v", err)
	}
	defer rows.Close()

	var sensorMeasurementTypes []dto.SensorMeasurementTypeDTO
	for rows.Next() {
		var sensorMeasurementType dto.SensorMeasurementTypeDTO
		var propertiesSchema, uiConfig []byte

		err := rows.Scan(
			&sensorMeasurementType.ID,
			&sensorMeasurementType.SensorTypeID,
			&sensorMeasurementType.Name,
			&sensorMeasurementType.Description,
			&propertiesSchema,
			&uiConfig,
			&sensorMeasurementType.Version,
			&sensorMeasurementType.IsActive,
			&sensorMeasurementType.CreatedAt,
			&sensorMeasurementType.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning sensor measurement type: %v", err)
		}

		// Unmarshal JSON fields
		if len(propertiesSchema) > 0 {
			if err := json.Unmarshal(propertiesSchema, &sensorMeasurementType.PropertiesSchema); err != nil {
				return nil, fmt.Errorf("error unmarshaling properties schema: %v", err)
			}
		}

		if len(uiConfig) > 0 {
			if err := json.Unmarshal(uiConfig, &sensorMeasurementType.UIConfiguration); err != nil {
				return nil, fmt.Errorf("error unmarshaling UI configuration: %v", err)
			}
		}

		sensorMeasurementTypes = append(sensorMeasurementTypes, sensorMeasurementType)
	}

	return sensorMeasurementTypes, nil
}

// GetBySensorTypeID retrieves all measurement types for a specific sensor type
func (r *sensorMeasurementTypeRepository) GetBySensorTypeID(ctx context.Context, sensorTypeID uuid.UUID) ([]dto.SensorMeasurementTypeDTO, error) {
	query := `
		SELECT id, sensor_type_id, name, description, properties_schema,
			ui_configuration, version, is_active, created_at, updated_at
		FROM sensor_measurement_types
		WHERE sensor_type_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sensorTypeID)
	if err != nil {
		return nil, fmt.Errorf("error querying sensor measurement types by sensor type ID: %v", err)
	}
	defer rows.Close()

	var sensorMeasurementTypes []dto.SensorMeasurementTypeDTO
	for rows.Next() {
		var sensorMeasurementType dto.SensorMeasurementTypeDTO
		var propertiesSchema, uiConfig []byte

		err := rows.Scan(
			&sensorMeasurementType.ID,
			&sensorMeasurementType.SensorTypeID,
			&sensorMeasurementType.Name,
			&sensorMeasurementType.Description,
			&propertiesSchema,
			&uiConfig,
			&sensorMeasurementType.Version,
			&sensorMeasurementType.IsActive,
			&sensorMeasurementType.CreatedAt,
			&sensorMeasurementType.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning sensor measurement type: %v", err)
		}

		// Unmarshal JSON fields
		if len(propertiesSchema) > 0 {
			if err := json.Unmarshal(propertiesSchema, &sensorMeasurementType.PropertiesSchema); err != nil {
				return nil, fmt.Errorf("error unmarshaling properties schema: %v", err)
			}
		}

		if len(uiConfig) > 0 {
			if err := json.Unmarshal(uiConfig, &sensorMeasurementType.UIConfiguration); err != nil {
				return nil, fmt.Errorf("error unmarshaling UI configuration: %v", err)
			}
		}

		sensorMeasurementTypes = append(sensorMeasurementTypes, sensorMeasurementType)
	}

	return sensorMeasurementTypes, nil
}
