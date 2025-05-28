package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/helpers/common"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// AssetSensorWithDetails represents an asset sensor with all its related information
type AssetSensorWithDetails struct {
	*entity.AssetSensor
	SensorType struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		Manufacturer string    `json:"manufacturer"`
		Model        string    `json:"model"`
		Version      int       `json:"version"`
		IsActive     bool      `json:"is_active"`
	} `json:"sensor_type"`
	MeasurementTypes []struct {
		ID               uuid.UUID       `json:"id"`
		Name             string          `json:"name"`
		Description      string          `json:"description"`
		UnitOfMeasure    string          `json:"unit_of_measure"`
		PropertiesSchema json.RawMessage `json:"properties_schema"`
		UIConfiguration  json.RawMessage `json:"ui_configuration"`
		Version          int             `json:"version"`
		IsActive         bool            `json:"is_active"`
		Fields           []struct {
			ID          uuid.UUID `json:"id"`
			Name        string    `json:"name"`
			Label       string    `json:"label"`
			Description *string   `json:"description"`
			DataType    string    `json:"data_type"`
			Required    bool      `json:"required"`
			Unit        *string   `json:"unit"`
			Min         *float64  `json:"min"`
			Max         *float64  `json:"max"`
		} `json:"fields"`
	} `json:"measurement_types"`
}

// AssetSensorRepository defines the interface for asset sensor data operations
type AssetSensorRepository interface {
	Create(ctx context.Context, sensor *entity.AssetSensor) error
	GetByID(ctx context.Context, id uuid.UUID) (*AssetSensorWithDetails, error)
	GetByAssetID(ctx context.Context, assetID uuid.UUID) ([]*AssetSensorWithDetails, error)
	GetBySensorTypeID(ctx context.Context, sensorTypeID uuid.UUID) ([]*AssetSensorWithDetails, error)
	List(ctx context.Context, page, pageSize int) ([]*AssetSensorWithDetails, error)
	Update(ctx context.Context, sensor *entity.AssetSensor) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByAssetID(ctx context.Context, assetID uuid.UUID) error
	UpdateLastReading(ctx context.Context, id uuid.UUID, value float64, readings map[string]interface{}) error
	GetActiveSensors(ctx context.Context) ([]*AssetSensorWithDetails, error)
	GetSensorsByStatus(ctx context.Context, status string) ([]*AssetSensorWithDetails, error)
}

// assetSensorRepository handles database operations for asset sensors
type assetSensorRepository struct {
	*BaseRepository
}

// NewAssetSensorRepository creates a new AssetSensorRepository
func NewAssetSensorRepository(db *sql.DB) AssetSensorRepository {
	return &assetSensorRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create inserts a new asset sensor into the database
func (r *assetSensorRepository) Create(ctx context.Context, sensor *entity.AssetSensor) error {
	log.Printf("Starting to create asset sensor in database: %+v", sensor)

	// Convert json.RawMessage to string for JSONB
	var configStr, readingsStr string
	if sensor.Configuration != nil {
		configStr = string(sensor.Configuration)
	} else {
		configStr = "{}"
	}
	if sensor.LastReadingValues != nil {
		readingsStr = string(sensor.LastReadingValues)
	} else {
		readingsStr = "{}"
	}

	query := `
		INSERT INTO asset_sensors (
			tenant_id, asset_id, sensor_type_id, name, status, 
			configuration, last_reading_value, last_reading_time, last_reading_values,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6::jsonb, $7, $8, $9::jsonb, $10, $11
		) RETURNING id`

	now := time.Now()

	log.Printf("Executing query with parameters: tenant_id=%v, asset_id=%v, sensor_type_id=%v, name=%v, status=%v, config=%v",
		sensor.TenantID, sensor.AssetID, sensor.SensorTypeID, sensor.Name, sensor.Status, configStr)

	err := r.DB.QueryRowContext(
		ctx,
		query,
		sensor.TenantID,
		sensor.AssetID,
		sensor.SensorTypeID,
		sensor.Name,
		sensor.Status,
		configStr,
		sensor.LastReadingValue,
		sensor.LastReadingTime,
		readingsStr,
		now,
		now,
	).Scan(&sensor.ID)

	if err != nil {
		log.Printf("Error creating asset sensor in database: %v", err)
		return fmt.Errorf("failed to create asset sensor: %w", err)
	}

	sensor.CreatedAt = now
	sensor.UpdatedAt = &now

	log.Printf("Successfully created asset sensor with ID: %s", sensor.ID)
	return nil
}

// GetByID retrieves an asset sensor by its ID with all related information
func (r *assetSensorRepository) GetByID(ctx context.Context, id uuid.UUID) (*AssetSensorWithDetails, error) {
	query := `
		WITH sensor_details AS (
			SELECT 
				asn.*,
				st.id as st_id, st.name as st_name, st.description as st_description,
				st.manufacturer as st_manufacturer, st.model as st_model,
				st.version as st_version, st.is_active as st_is_active
			FROM asset_sensors asn
			JOIN sensor_types st ON asn.sensor_type_id = st.id
			WHERE asn.id = $1
		)
		SELECT 
			sd.*,
			json_agg(
				json_build_object(
					'id', smt.id,
					'name', smt.name,
					'description', smt.description,
					'unit_of_measure', smt.unit_of_measure,
					'properties_schema', smt.properties_schema,
					'ui_configuration', smt.ui_configuration,
					'version', smt.version,
					'is_active', smt.is_active,
					'fields', (
						SELECT json_agg(
							json_build_object(
								'id', smf.id,
								'name', smf.name,
								'label', smf.label,
								'description', smf.description,
								'data_type', smf.data_type,
								'required', smf.required,
								'unit', smf.unit,
								'min', smf.min,
								'max', smf.max
							)
						)
						FROM sensor_measurement_fields smf
						WHERE smf.sensor_measurement_type_id = smt.id
					)
				)
			) as measurement_types
		FROM sensor_details sd
		LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = sd.st_id
		GROUP BY sd.id, sd.tenant_id, sd.asset_id, sd.sensor_type_id, sd.name, sd.status,
				sd.configuration, sd.last_reading_value, sd.last_reading_time, sd.last_reading_values,
				sd.created_at, sd.updated_at, sd.st_id, sd.st_name, sd.st_description,
				sd.st_manufacturer, sd.st_model, sd.st_version, sd.st_is_active`

	var result AssetSensorWithDetails
	var sensor entity.AssetSensor
	var measurementTypesJSON []byte

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&sensor.ID,
		&sensor.TenantID,
		&sensor.AssetID,
		&sensor.SensorTypeID,
		&sensor.Name,
		&sensor.Status,
		&sensor.Configuration,
		&sensor.LastReadingValue,
		&sensor.LastReadingTime,
		&sensor.LastReadingValues,
		&sensor.CreatedAt,
		&sensor.UpdatedAt,
		&result.SensorType.ID,
		&result.SensorType.Name,
		&result.SensorType.Description,
		&result.SensorType.Manufacturer,
		&result.SensorType.Model,
		&result.SensorType.Version,
		&result.SensorType.IsActive,
		&measurementTypesJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get asset sensor: %w", err)
	}

	result.AssetSensor = &sensor

	// Parse measurement types JSON
	if measurementTypesJSON != nil {
		if err := json.Unmarshal(measurementTypesJSON, &result.MeasurementTypes); err != nil {
			return nil, fmt.Errorf("failed to parse measurement types: %w", err)
		}
	}

	return &result, nil
}

// GetByAssetID retrieves all sensors for a specific asset
func (r *assetSensorRepository) GetByAssetID(ctx context.Context, assetID uuid.UUID) ([]*AssetSensorWithDetails, error) {
	query := `
		WITH sensor_details AS (
			SELECT 
				asn.*,
				st.id as st_id, st.name as st_name, st.description as st_description,
				st.manufacturer as st_manufacturer, st.model as st_model,
				st.version as st_version, st.is_active as st_is_active
			FROM asset_sensors asn
			JOIN sensor_types st ON asn.sensor_type_id = st.id
			WHERE asn.asset_id = $1
		)
		SELECT 
			sd.*,
			json_agg(
				json_build_object(
					'id', smt.id,
					'name', smt.name,
					'description', smt.description,
					'unit_of_measure', smt.unit_of_measure,
					'properties_schema', smt.properties_schema,
					'ui_configuration', smt.ui_configuration,
					'version', smt.version,
					'is_active', smt.is_active,
					'fields', (
						SELECT json_agg(
							json_build_object(
								'id', smf.id,
								'name', smf.name,
								'label', smf.label,
								'description', smf.description,
								'data_type', smf.data_type,
								'required', smf.required,
								'unit', smf.unit,
								'min', smf.min,
								'max', smf.max
							)
						)
						FROM sensor_measurement_fields smf
						WHERE smf.sensor_measurement_type_id = smt.id
					)
				)
			) as measurement_types
		FROM sensor_details sd
		LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = sd.st_id
		GROUP BY sd.id, sd.tenant_id, sd.asset_id, sd.sensor_type_id, sd.name, sd.status,
				sd.configuration, sd.last_reading_value, sd.last_reading_time, sd.last_reading_values,
				sd.created_at, sd.updated_at, sd.st_id, sd.st_name, sd.st_description,
				sd.st_manufacturer, sd.st_model, sd.st_version, sd.st_is_active
		ORDER BY sd.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset sensors: %w", err)
	}
	defer rows.Close()

	var sensors []*AssetSensorWithDetails
	for rows.Next() {
		var result AssetSensorWithDetails
		var sensor entity.AssetSensor
		var measurementTypesJSON []byte

		err := rows.Scan(
			&sensor.ID,
			&sensor.TenantID,
			&sensor.AssetID,
			&sensor.SensorTypeID,
			&sensor.Name,
			&sensor.Status,
			&sensor.Configuration,
			&sensor.LastReadingValue,
			&sensor.LastReadingTime,
			&sensor.LastReadingValues,
			&sensor.CreatedAt,
			&sensor.UpdatedAt,
			&result.SensorType.ID,
			&result.SensorType.Name,
			&result.SensorType.Description,
			&result.SensorType.Manufacturer,
			&result.SensorType.Model,
			&result.SensorType.Version,
			&result.SensorType.IsActive,
			&measurementTypesJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan asset sensor: %w", err)
		}

		result.AssetSensor = &sensor

		// Parse measurement types JSON
		if measurementTypesJSON != nil {
			if err := json.Unmarshal(measurementTypesJSON, &result.MeasurementTypes); err != nil {
				return nil, fmt.Errorf("failed to parse measurement types: %w", err)
			}
		}

		sensors = append(sensors, &result)
	}

	return sensors, nil
}

// GetBySensorTypeID retrieves all sensors of a specific type
func (r *assetSensorRepository) GetBySensorTypeID(ctx context.Context, sensorTypeID uuid.UUID) ([]*AssetSensorWithDetails, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access all sensors of the type
		query = `
			WITH sensor_details AS (
				SELECT 
					asn.*,
					st.id as st_id, st.name as st_name, st.description as st_description,
					st.manufacturer as st_manufacturer, st.model as st_model,
					st.version as st_version, st.is_active as st_is_active
				FROM asset_sensors asn
				JOIN sensor_types st ON asn.sensor_type_id = st.id
				WHERE asn.sensor_type_id = $1
			)
			SELECT 
				sd.*,
				json_agg(
					json_build_object(
						'id', smt.id,
						'name', smt.name,
						'description', smt.description,
						'unit_of_measure', smt.unit_of_measure,
						'properties_schema', smt.properties_schema,
						'ui_configuration', smt.ui_configuration,
						'version', smt.version,
						'is_active', smt.is_active,
						'fields', (
							SELECT json_agg(
								json_build_object(
									'id', smf.id,
									'name', smf.name,
									'label', smf.label,
									'description', smf.description,
									'data_type', smf.data_type,
									'required', smf.required,
									'unit', smf.unit,
									'min', smf.min,
									'max', smf.max
								)
							)
							FROM sensor_measurement_fields smf
							WHERE smf.sensor_measurement_type_id = smt.id
						)
					)
				) as measurement_types
			FROM sensor_details sd
			LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = sd.st_id
			GROUP BY sd.id, sd.tenant_id, sd.asset_id, sd.sensor_type_id, sd.name, sd.status,
					sd.configuration, sd.last_reading_value, sd.last_reading_time, sd.last_reading_values,
					sd.created_at, sd.updated_at, sd.st_id, sd.st_name, sd.st_description,
					sd.st_manufacturer, sd.st_model, sd.st_version, sd.st_is_active
			ORDER BY sd.created_at DESC`
		args = []interface{}{sensorTypeID}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			WITH sensor_details AS (
				SELECT 
					asn.*,
					st.id as st_id, st.name as st_name, st.description as st_description,
					st.manufacturer as st_manufacturer, st.model as st_model,
					st.version as st_version, st.is_active as st_is_active
				FROM asset_sensors asn
				JOIN sensor_types st ON asn.sensor_type_id = st.id
				WHERE asn.sensor_type_id = $1 AND asn.tenant_id = $2
			)
			SELECT 
				sd.*,
				json_agg(
					json_build_object(
						'id', smt.id,
						'name', smt.name,
						'description', smt.description,
						'unit_of_measure', smt.unit_of_measure,
						'properties_schema', smt.properties_schema,
						'ui_configuration', smt.ui_configuration,
						'version', smt.version,
						'is_active', smt.is_active,
						'fields', (
							SELECT json_agg(
								json_build_object(
									'id', smf.id,
									'name', smf.name,
									'label', smf.label,
									'description', smf.description,
									'data_type', smf.data_type,
									'required', smf.required,
									'unit', smf.unit,
									'min', smf.min,
									'max', smf.max
								)
							)
							FROM sensor_measurement_fields smf
							WHERE smf.sensor_measurement_type_id = smt.id
						)
					)
				) as measurement_types
			FROM sensor_details sd
			LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = sd.st_id
			GROUP BY sd.id, sd.tenant_id, sd.asset_id, sd.sensor_type_id, sd.name, sd.status,
					sd.configuration, sd.last_reading_value, sd.last_reading_time, sd.last_reading_values,
					sd.created_at, sd.updated_at, sd.st_id, sd.st_name, sd.st_description,
					sd.st_manufacturer, sd.st_model, sd.st_version, sd.st_is_active
			ORDER BY sd.created_at DESC`
		args = []interface{}{sensorTypeID, tenantID}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset sensors: %w", err)
	}
	defer rows.Close()

	var sensors []*AssetSensorWithDetails
	for rows.Next() {
		var result AssetSensorWithDetails
		var sensor entity.AssetSensor
		var measurementTypesJSON []byte

		err := rows.Scan(
			&sensor.ID,
			&sensor.TenantID,
			&sensor.AssetID,
			&sensor.SensorTypeID,
			&sensor.Name,
			&sensor.Status,
			&sensor.Configuration,
			&sensor.LastReadingValue,
			&sensor.LastReadingTime,
			&sensor.LastReadingValues,
			&sensor.CreatedAt,
			&sensor.UpdatedAt,
			&result.SensorType.ID,
			&result.SensorType.Name,
			&result.SensorType.Description,
			&result.SensorType.Manufacturer,
			&result.SensorType.Model,
			&result.SensorType.Version,
			&result.SensorType.IsActive,
			&measurementTypesJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan asset sensor: %w", err)
		}

		result.AssetSensor = &sensor

		// Parse measurement types JSON
		if measurementTypesJSON != nil {
			if err := json.Unmarshal(measurementTypesJSON, &result.MeasurementTypes); err != nil {
				return nil, fmt.Errorf("failed to parse measurement types: %w", err)
			}
		}

		sensors = append(sensors, &result)
	}

	return sensors, nil
}

// List retrieves asset sensors with pagination
func (r *assetSensorRepository) List(ctx context.Context, page, pageSize int) ([]*AssetSensorWithDetails, error) {
	log.Printf("Starting to list asset sensors with page=%d, pageSize=%d", page, page)
	offset := (page - 1) * pageSize

	query := `
	WITH sensor_details AS (
		SELECT 
			asn.id,
			asn.tenant_id,
			asn.asset_id,
			asn.name,
			asn.sensor_type_id,
			asn.status,
			asn.configuration,
			asn.last_reading_value,
			asn.last_reading_time,
			asn.last_reading_values,
			asn.created_at,
			asn.updated_at,
			jsonb_build_object(
				'id', st.id,
				'name', st.name,
				'description', st.description,
				'manufacturer', st.manufacturer,
				'model', st.model,
				'version', st.version,
				'is_active', st.is_active
			) as sensor_type,
			COALESCE(
				jsonb_agg(
					DISTINCT jsonb_build_object(
						'id', mt.id,
						'name', mt.name,
						'description', mt.description,
						'unit_of_measure', mt.unit_of_measure,
						'properties_schema', mt.properties_schema,
						'ui_configuration', mt.ui_configuration,
						'version', mt.version,
						'is_active', mt.is_active,
						'fields', (
							SELECT jsonb_agg(
								jsonb_build_object(
									'id', mf.id,
									'name', mf.name,
									'label', mf.label,
									'description', mf.description,
									'data_type', mf.data_type,
									'required', mf.required,
									'unit', mf.unit,
									'min', mf.min,
									'max', mf.max
								)
							)
							FROM sensor_measurement_fields mf
							WHERE mf.sensor_measurement_type_id = mt.id
						)
					)
				) FILTER (WHERE mt.id IS NOT NULL),
				'[]'::jsonb
			) as measurement_types
		FROM asset_sensors asn
		LEFT JOIN sensor_types st ON asn.sensor_type_id = st.id
		LEFT JOIN sensor_measurement_types mt ON mt.sensor_type_id = st.id
		GROUP BY asn.id, st.id
		ORDER BY asn.created_at DESC
		LIMIT $1 OFFSET $2
	)
	SELECT * FROM sensor_details;
	`

	rows, err := r.DB.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("failed to list asset sensors: %w", err)
	}
	defer rows.Close()

	var sensors []*AssetSensorWithDetails
	for rows.Next() {
		var sensor AssetSensorWithDetails
		sensor.AssetSensor = &entity.AssetSensor{} // Initialize the AssetSensor
		var sensorTypeJSON, measurementTypesJSON []byte
		var tenantID, assetID, sensorTypeID sql.NullString
		var lastReadingValue sql.NullFloat64
		var lastReadingTime sql.NullTime
		var updatedAt sql.NullTime

		err := rows.Scan(
			&sensor.AssetSensor.ID,
			&tenantID,
			&assetID,
			&sensor.AssetSensor.Name,
			&sensorTypeID,
			&sensor.AssetSensor.Status,
			&sensor.AssetSensor.Configuration,
			&lastReadingValue,
			&lastReadingTime,
			&sensor.AssetSensor.LastReadingValues,
			&sensor.AssetSensor.CreatedAt,
			&updatedAt,
			&sensorTypeJSON,
			&measurementTypesJSON,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan asset sensor row: %w", err)
		}

		// Handle nullable fields
		if tenantID.Valid {
			id, _ := uuid.Parse(tenantID.String)
			sensor.AssetSensor.TenantID = &id
		}
		if assetID.Valid {
			sensor.AssetSensor.AssetID, _ = uuid.Parse(assetID.String)
		}
		if sensorTypeID.Valid {
			sensor.AssetSensor.SensorTypeID, _ = uuid.Parse(sensorTypeID.String)
		}
		if lastReadingValue.Valid {
			sensor.AssetSensor.LastReadingValue = &lastReadingValue.Float64
		}
		if lastReadingTime.Valid {
			sensor.AssetSensor.LastReadingTime = &lastReadingTime.Time
		}
		if updatedAt.Valid {
			sensor.AssetSensor.UpdatedAt = &updatedAt.Time
		}

		// Parse sensor type JSON
		if err := json.Unmarshal(sensorTypeJSON, &sensor.SensorType); err != nil {
			log.Printf("Error parsing sensor type JSON: %v", err)
			return nil, fmt.Errorf("failed to parse sensor type JSON: %w", err)
		}

		// Parse measurement types JSON
		if err := json.Unmarshal(measurementTypesJSON, &sensor.MeasurementTypes); err != nil {
			log.Printf("Error parsing measurement types JSON: %v", err)
			return nil, fmt.Errorf("failed to parse measurement types JSON: %w", err)
		}

		sensors = append(sensors, &sensor)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating asset sensor rows: %w", err)
	}

	return sensors, nil
}

// Update modifies an existing asset sensor
func (r *assetSensorRepository) Update(ctx context.Context, sensor *entity.AssetSensor) error {
	// Convert json.RawMessage to string for JSONB
	var configStr, readingsStr string
	if sensor.Configuration != nil {
		configStr = string(sensor.Configuration)
	} else {
		configStr = "{}"
	}
	if sensor.LastReadingValues != nil {
		readingsStr = string(sensor.LastReadingValues)
	} else {
		readingsStr = "{}"
	}

	query := `
		UPDATE asset_sensors
		SET name = $1, status = $2, configuration = $3::jsonb,
			last_reading_value = $4, last_reading_time = $5, last_reading_values = $6::jsonb,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $7`

	result, err := r.DB.ExecContext(ctx, query,
		sensor.Name, sensor.Status, configStr,
		sensor.LastReadingValue, sensor.LastReadingTime, readingsStr,
		sensor.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update asset sensor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset sensor not found")
	}

	return nil
}

// Delete removes an asset sensor by ID
func (r *assetSensorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM asset_sensors WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset sensor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset sensor not found")
	}

	return nil
}

// DeleteByAssetID removes all sensors for a specific asset
func (r *assetSensorRepository) DeleteByAssetID(ctx context.Context, assetID uuid.UUID) error {
	query := `DELETE FROM asset_sensors WHERE asset_id = $1`

	_, err := r.DB.ExecContext(ctx, query, assetID)
	if err != nil {
		return fmt.Errorf("failed to delete asset sensors: %w", err)
	}

	return nil
}

// UpdateLastReading updates the last reading value and time for a sensor
func (r *assetSensorRepository) UpdateLastReading(ctx context.Context, id uuid.UUID, value float64, readings map[string]interface{}) error {
	now := time.Now()

	// Convert readings map to JSON string
	readingsJSON, err := json.Marshal(readings)
	if err != nil {
		return fmt.Errorf("failed to marshal readings: %w", err)
	}
	readingsStr := string(readingsJSON)

	query := `
		UPDATE asset_sensors
		SET last_reading_value = $1,
			last_reading_time = $2,
			last_reading_values = $3::jsonb,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4`

	result, err := r.DB.ExecContext(ctx, query, value, now, readingsStr, id)
	if err != nil {
		return fmt.Errorf("failed to update sensor reading: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset sensor not found")
	}

	return nil
}

// GetActiveSensors retrieves all active sensors
func (r *assetSensorRepository) GetActiveSensors(ctx context.Context) ([]*AssetSensorWithDetails, error) {
	query := `
		WITH sensor_details AS (
			SELECT 
				asn.*,
				st.id as st_id, st.name as st_name, st.description as st_description,
				st.manufacturer as st_manufacturer, st.model as st_model,
				st.version as st_version, st.is_active as st_is_active
			FROM asset_sensors asn
			JOIN sensor_types st ON asn.sensor_type_id = st.id
			WHERE asn.status = 'active'
		)
		SELECT 
			sd.*,
			json_agg(
				json_build_object(
					'id', smt.id,
					'name', smt.name,
					'description', smt.description,
					'unit_of_measure', smt.unit_of_measure,
					'properties_schema', smt.properties_schema,
					'ui_configuration', smt.ui_configuration,
					'version', smt.version,
					'is_active', smt.is_active,
					'fields', (
						SELECT json_agg(
							json_build_object(
								'id', smf.id,
								'name', smf.name,
								'label', smf.label,
								'description', smf.description,
								'data_type', smf.data_type,
								'required', smf.required,
								'unit', smf.unit,
								'min', smf.min,
								'max', smf.max
							)
						)
						FROM sensor_measurement_fields smf
						WHERE smf.sensor_measurement_type_id = smt.id
					)
				)
			) as measurement_types
		FROM sensor_details sd
		LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = sd.st_id
		GROUP BY sd.id, sd.tenant_id, sd.asset_id, sd.sensor_type_id, sd.name, sd.status,
				sd.configuration, sd.last_reading_value, sd.last_reading_time, sd.last_reading_values,
				sd.created_at, sd.updated_at, sd.st_id, sd.st_name, sd.st_description,
				sd.st_manufacturer, sd.st_model, sd.st_version, sd.st_is_active
		ORDER BY sd.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sensors: %w", err)
	}
	defer rows.Close()

	var sensors []*AssetSensorWithDetails
	for rows.Next() {
		var result AssetSensorWithDetails
		var sensor entity.AssetSensor
		var measurementTypesJSON []byte

		err := rows.Scan(
			&sensor.ID,
			&sensor.TenantID,
			&sensor.AssetID,
			&sensor.SensorTypeID,
			&sensor.Name,
			&sensor.Status,
			&sensor.Configuration,
			&sensor.LastReadingValue,
			&sensor.LastReadingTime,
			&sensor.LastReadingValues,
			&sensor.CreatedAt,
			&sensor.UpdatedAt,
			&result.SensorType.ID,
			&result.SensorType.Name,
			&result.SensorType.Description,
			&result.SensorType.Manufacturer,
			&result.SensorType.Model,
			&result.SensorType.Version,
			&result.SensorType.IsActive,
			&measurementTypesJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan asset sensor: %w", err)
		}

		result.AssetSensor = &sensor

		// Parse measurement types JSON
		if measurementTypesJSON != nil {
			if err := json.Unmarshal(measurementTypesJSON, &result.MeasurementTypes); err != nil {
				return nil, fmt.Errorf("failed to parse measurement types: %w", err)
			}
		}

		sensors = append(sensors, &result)
	}

	return sensors, nil
}

// GetSensorsByStatus retrieves all sensors with a specific status
func (r *assetSensorRepository) GetSensorsByStatus(ctx context.Context, status string) ([]*AssetSensorWithDetails, error) {
	query := `
		WITH sensor_details AS (
			SELECT 
				as.*,
				st.id as st_id, st.name as st_name, st.description as st_description,
				st.manufacturer as st_manufacturer, st.model as st_model,
				st.version as st_version, st.is_active as st_is_active
			FROM asset_sensors asn
			JOIN sensor_types st ON asn.sensor_type_id = st.id
			WHERE asn.status = $1
		)
		SELECT 
			sd.*,
			json_agg(
				json_build_object(
					'id', smt.id,
					'name', smt.name,
					'description', smt.description,
					'unit_of_measure', smt.unit_of_measure,
					'properties_schema', smt.properties_schema,
					'ui_configuration', smt.ui_configuration,
					'version', smt.version,
					'is_active', smt.is_active,
					'fields', (
						SELECT json_agg(
							json_build_object(
								'id', smf.id,
								'name', smf.name,
								'label', smf.label,
								'description', smf.description,
								'data_type', smf.data_type,
								'required', smf.required,
								'unit', smf.unit,
								'min', smf.min,
								'max', smf.max
							)
						)
						FROM sensor_measurement_fields smf
						WHERE smf.sensor_measurement_type_id = smt.id
					)
				)
			) as measurement_types
		FROM sensor_details sd
		LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = sd.st_id
		GROUP BY sd.id, sd.tenant_id, sd.asset_id, sd.sensor_type_id, sd.name, sd.status,
				sd.configuration, sd.last_reading_value, sd.last_reading_time, sd.last_reading_values,
				sd.created_at, sd.updated_at, sd.st_id, sd.st_name, sd.st_description,
				sd.st_manufacturer, sd.st_model, sd.st_version, sd.st_is_active
		ORDER BY sd.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensors by status: %w", err)
	}
	defer rows.Close()

	var sensors []*AssetSensorWithDetails
	for rows.Next() {
		var result AssetSensorWithDetails
		var sensor entity.AssetSensor
		var measurementTypesJSON []byte

		err := rows.Scan(
			&sensor.ID,
			&sensor.TenantID,
			&sensor.AssetID,
			&sensor.SensorTypeID,
			&sensor.Name,
			&sensor.Status,
			&sensor.Configuration,
			&sensor.LastReadingValue,
			&sensor.LastReadingTime,
			&sensor.LastReadingValues,
			&sensor.CreatedAt,
			&sensor.UpdatedAt,
			&result.SensorType.ID,
			&result.SensorType.Name,
			&result.SensorType.Description,
			&result.SensorType.Manufacturer,
			&result.SensorType.Model,
			&result.SensorType.Version,
			&result.SensorType.IsActive,
			&measurementTypesJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan asset sensor: %w", err)
		}

		result.AssetSensor = &sensor

		// Parse measurement types JSON
		if measurementTypesJSON != nil {
			if err := json.Unmarshal(measurementTypesJSON, &result.MeasurementTypes); err != nil {
				return nil, fmt.Errorf("failed to parse measurement types: %w", err)
			}
		}

		sensors = append(sensors, &result)
	}

	return sensors, nil
}
