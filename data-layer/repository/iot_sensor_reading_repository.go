package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

// IoTSensorReadingWithDetails represents an IoT sensor reading with all its related information
type IoTSensorReadingWithDetails struct {
	*entity.IoTSensorReading
	AssetSensor struct {
		ID            uuid.UUID       `json:"id"`
		AssetID       uuid.UUID       `json:"asset_id"`
		Name          string          `json:"name"`
		Status        string          `json:"status"`
		Configuration json.RawMessage `json:"configuration"`
	} `json:"asset_sensor"`
	SensorType struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		Manufacturer string    `json:"manufacturer"`
		Model        string    `json:"model"`
		Version      string    `json:"version"`
		IsActive     bool      `json:"is_active"`
	} `json:"sensor_type"`
	MeasurementTypes []struct {
		ID               uuid.UUID       `json:"id"`
		Name             string          `json:"name"`
		Description      string          `json:"description"`
		PropertiesSchema json.RawMessage `json:"properties_schema"`
		UIConfiguration  json.RawMessage `json:"ui_configuration"`
		Version          string          `json:"version"`
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

// IoTSensorReadingListRequest represents parameters for listing IoT sensor readings
type IoTSensorReadingListRequest struct {
	AssetSensorID *uuid.UUID `json:"asset_sensor_id,omitempty"`
	SensorTypeID  *uuid.UUID `json:"sensor_type_id,omitempty"`
	MacAddress    *string    `json:"mac_address,omitempty"`
	LocationID    *uuid.UUID `json:"location_id,omitempty"`
	FromTime      *time.Time `json:"from_time,omitempty"`
	ToTime        *time.Time `json:"to_time,omitempty"`
	Page          int        `json:"page"`
	PageSize      int        `json:"page_size"`
}

// IoTSensorReadingRepository defines the interface for IoT sensor reading data operations
type IoTSensorReadingRepository interface {
	Create(ctx context.Context, reading *entity.IoTSensorReading) error
	CreateBatch(ctx context.Context, readings []*entity.IoTSensorReading) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.IoTSensorReading, error)
	GetByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID, limit int) ([]*IoTSensorReadingWithDetails, error)
	GetBySensorTypeID(ctx context.Context, sensorTypeID uuid.UUID, limit int) ([]*IoTSensorReadingWithDetails, error)
	GetByMacAddress(ctx context.Context, macAddress string, limit int) ([]*IoTSensorReadingWithDetails, error)
	GetAssetSensorsBySensorType(ctx context.Context, sensorTypeID uuid.UUID) ([]dto.AssetSensorLocationInfo, error)
	List(ctx context.Context, req IoTSensorReadingListRequest) ([]*IoTSensorReadingWithDetails, int, error)
	Update(ctx context.Context, reading *entity.IoTSensorReading) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) error
	GetLatestReading(ctx context.Context, assetSensorID uuid.UUID) (*IoTSensorReadingWithDetails, error)
	GetReadingsInTimeRange(ctx context.Context, assetSensorID uuid.UUID, fromTime, toTime time.Time) ([]*IoTSensorReadingWithDetails, error)
	GetAggregatedData(ctx context.Context, assetSensorID uuid.UUID, fromTime, toTime time.Time, interval string) ([]map[string]interface{}, error)
	ValidateAndCreate(ctx context.Context, reading *entity.IoTSensorReading) (bool, []string, error)
	CreateFlexible(ctx context.Context, reading *entity.IoTSensorReadingFlexible) error
	CreateFlexibleBatch(ctx context.Context, readings []*entity.IoTSensorReadingFlexible) error
	GetFlexibleByID(ctx context.Context, id uuid.UUID) (*entity.IoTSensorReadingFlexible, error)
	ListFlexible(ctx context.Context, req IoTSensorReadingListRequest) ([]*entity.IoTSensorReadingFlexible, int, error)
	ParseTextToFlexibleReading(ctx context.Context, textData, assetSensorID, sensorTypeID, macAddress string) (*entity.IoTSensorReadingFlexible, error)
	GetDB() *sql.DB
}

// iotSensorReadingRepository handles database operations for IoT sensor readings
type iotSensorReadingRepository struct {
	*BaseRepository
}

// NewIoTSensorReadingRepository creates a new IoTSensorReadingRepository
func NewIoTSensorReadingRepository(db *sql.DB) IoTSensorReadingRepository {
	return &iotSensorReadingRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create inserts a new IoT sensor reading into the database
func (r *iotSensorReadingRepository) Create(ctx context.Context, reading *entity.IoTSensorReading) error {
	log.Printf("Starting to create IoT sensor reading in database: %+v", reading)

	// Generate new UUID if not set
	if reading.ID == uuid.Nil {
		reading.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now()
	if reading.CreatedAt.IsZero() {
		reading.CreatedAt = now
	}
	reading.UpdatedAt = &now

	// Build and execute the insert query
	query := `
		INSERT INTO iot_sensor_readings (
			id, tenant_id, asset_sensor_id, sensor_type_id, mac_address, 
			location_id, location_name, measurement_type, measurement_label, 
			measurement_unit, numeric_value, text_value, boolean_value, 
			data_source, original_field_name, reading_time, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)`

	_, err := r.DB.ExecContext(ctx, query,
		reading.ID,
		reading.TenantID,
		reading.AssetSensorID,
		reading.SensorTypeID,
		reading.MacAddress,
		reading.LocationID,
		reading.LocationName,
		reading.MeasurementType,
		reading.MeasurementLabel,
		reading.MeasurementUnit,
		reading.NumericValue,
		reading.TextValue,
		reading.BooleanValue,
		reading.DataSource,
		reading.OriginalFieldName,
		reading.ReadingTime,
		reading.CreatedAt,
		reading.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating IoT sensor reading in database: %v", err)
		return fmt.Errorf("failed to create IoT sensor reading: %w", err)
	}

	log.Printf("Successfully created IoT sensor reading with ID: %s", reading.ID)
	return nil
}

// CreateBatch inserts multiple IoT sensor readings into the database efficiently
func (r *iotSensorReadingRepository) CreateBatch(ctx context.Context, readings []*entity.IoTSensorReading) error {
	if len(readings) == 0 {
		return nil
	}

	log.Printf("Starting to create batch of %d IoT sensor readings", len(readings))

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO iot_sensor_readings (
			id, tenant_id, asset_sensor_id, sensor_type_id, mac_address,
			location_id, location_name, reading_time, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()

	for _, reading := range readings {
		// Set IDs and timestamps
		reading.ID = uuid.New()
		reading.CreatedAt = now
		reading.UpdatedAt = &now

		_, err = stmt.ExecContext(
			ctx,
			reading.ID,
			reading.TenantID,
			reading.AssetSensorID,
			reading.SensorTypeID,
			reading.MacAddress,
			reading.LocationID,
			reading.LocationName,
			reading.ReadingTime,
			now,
			now,
		)

		if err != nil {
			log.Printf("Error inserting reading %s: %v", reading.ID, err)
			return fmt.Errorf("failed to insert reading: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully created batch of %d IoT sensor readings", len(readings))
	return nil
}

// GetByID retrieves an IoT sensor reading by its ID
func (r *iotSensorReadingRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.IoTSensorReading, error) {
	log.Printf("Retrieving IoT sensor reading with ID: %s", id)

	query := `
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address,
			   location_id, location_name, reading_time, created_at, updated_at
		FROM iot_sensor_readings
		WHERE id = $1`

	var reading entity.IoTSensorReading
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&reading.ID,
		&reading.TenantID,
		&reading.AssetSensorID,
		&reading.SensorTypeID,
		&reading.MacAddress,
		&reading.LocationID,
		&reading.LocationName,
		&reading.ReadingTime,
		&reading.CreatedAt,
		&reading.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("IoT sensor reading not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get IoT sensor reading: %w", err)
	}

	log.Printf("Successfully retrieved IoT sensor reading with ID: %s", id)
	return &reading, nil
}

// Update updates an existing IoT sensor reading
func (r *iotSensorReadingRepository) Update(ctx context.Context, reading *entity.IoTSensorReading) error {
	log.Printf("Starting to update IoT sensor reading: %s", reading.ID)

	query := `
		UPDATE iot_sensor_readings 
		SET 
			asset_sensor_id = $2,
			sensor_type_id = $3,
			mac_address = $4,
			location_id = $5,
			location_name = $6,
			reading_time = $7,
			updated_at = $8
		WHERE id = $1`

	now := time.Now()

	result, err := r.DB.ExecContext(
		ctx,
		query,
		reading.ID,
		reading.AssetSensorID,
		reading.SensorTypeID,
		reading.MacAddress,
		reading.LocationID,
		reading.LocationName,
		reading.ReadingTime,
		now,
	)

	if err != nil {
		log.Printf("Error updating IoT sensor reading: %v", err)
		return fmt.Errorf("failed to update IoT sensor reading: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no IoT sensor reading found with ID: %s", reading.ID)
	}

	reading.UpdatedAt = &now

	log.Printf("Successfully updated IoT sensor reading: %s", reading.ID)
	return nil
}

// Delete removes an IoT sensor reading by its ID
func (r *iotSensorReadingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	log.Printf("Starting to delete IoT sensor reading: %s", id)

	query := `DELETE FROM iot_sensor_readings WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting IoT sensor reading: %v", err)
		return fmt.Errorf("failed to delete IoT sensor reading: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no IoT sensor reading found with ID: %s", id)
	}

	log.Printf("Successfully deleted IoT sensor reading: %s", id)
	return nil
}

// DeleteByAssetSensorID removes all IoT sensor readings for a specific asset sensor
func (r *iotSensorReadingRepository) DeleteByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) error {
	log.Printf("Starting to delete IoT sensor readings for asset sensor: %s", assetSensorID)

	query := `DELETE FROM iot_sensor_readings WHERE asset_sensor_id = $1`

	result, err := r.DB.ExecContext(ctx, query, assetSensorID)
	if err != nil {
		log.Printf("Error deleting IoT sensor readings for asset sensor: %v", err)
		return fmt.Errorf("failed to delete IoT sensor readings: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	log.Printf("Successfully deleted %d IoT sensor readings for asset sensor: %s", rowsAffected, assetSensorID)
	return nil
}

// GetLatestReading retrieves the most recent reading for a specific asset sensor
func (r *iotSensorReadingRepository) GetLatestReading(ctx context.Context, assetSensorID uuid.UUID) (*IoTSensorReadingWithDetails, error) {
	query := `
		WITH reading_details AS (
			SELECT 
				isr.*,
				asn.id as asn_id, asn.asset_id as asn_asset_id, asn.name as asn_name,
				asn.status as asn_status, asn.configuration as asn_configuration,
				st.id as st_id, st.name as st_name, st.description as st_description,
				st.manufacturer as st_manufacturer, st.model as st_model,
				st.version as st_version, st.is_active as st_is_active
			FROM iot_sensor_readings isr
			JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
			JOIN sensor_types st ON isr.sensor_type_id = st.id
			WHERE isr.asset_sensor_id = $1
			ORDER BY isr.reading_time DESC
			LIMIT 1
		)
		SELECT 
			rd.*,
			json_agg(
				json_build_object(
					'id', smt.id,
					'name', smt.name,
					'description', smt.description,
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
		FROM reading_details rd
		LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = rd.st_id
		GROUP BY rd.id, rd.tenant_id, rd.asset_sensor_id, rd.sensor_type_id, rd.mac_address,
				rd.location_id, rd.location_name, rd.reading_time, rd.created_at, rd.updated_at,
				rd.asn_id, rd.asn_asset_id, rd.asn_name, rd.asn_status, rd.asn_configuration,
				rd.st_id, rd.st_name, rd.st_description, rd.st_manufacturer, rd.st_model,
				rd.st_version, rd.st_is_active
		ORDER BY rd.reading_time DESC`

	var reading IoTSensorReadingWithDetails
	var measurementTypesJSON []byte

	err := r.DB.QueryRowContext(ctx, query, assetSensorID).Scan(
		&reading.ID,
		&reading.TenantID,
		&reading.AssetSensorID,
		&reading.SensorTypeID,
		&reading.MacAddress,
		&reading.LocationID,
		&reading.LocationName,
		&reading.ReadingTime,
		&reading.CreatedAt,
		&reading.UpdatedAt,
		&reading.AssetSensor.ID,
		&reading.AssetSensor.AssetID,
		&reading.AssetSensor.Name,
		&reading.AssetSensor.Status,
		&reading.AssetSensor.Configuration,
		&reading.SensorType.ID,
		&reading.SensorType.Name,
		&reading.SensorType.Description,
		&reading.SensorType.Manufacturer,
		&reading.SensorType.Model,
		&reading.SensorType.Version,
		&reading.SensorType.IsActive,
		&measurementTypesJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest reading: %w", err)
	}

	// Parse measurement types JSON
	if measurementTypesJSON != nil {
		err = json.Unmarshal(measurementTypesJSON, &reading.MeasurementTypes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal measurement types: %w", err)
		}
	}

	return &reading, nil
}

// GetReadingsInTimeRange retrieves readings within a specified time range
func (r *iotSensorReadingRepository) GetReadingsInTimeRange(ctx context.Context, assetSensorID uuid.UUID, fromTime, toTime time.Time) ([]*IoTSensorReadingWithDetails, error) {
	query := `
		WITH reading_details AS (
			SELECT 
				isr.*,
				asn.id as asn_id, asn.asset_id as asn_asset_id, asn.name as asn_name,
				asn.status as asn_status, asn.configuration as asn_configuration,
				st.id as st_id, st.name as st_name, st.description as st_description,
				st.manufacturer as st_manufacturer, st.model as st_model,
				st.version as st_version, st.is_active as st_is_active
			FROM iot_sensor_readings isr
			JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
			JOIN sensor_types st ON isr.sensor_type_id = st.id
			WHERE isr.asset_sensor_id = $1 
			  AND isr.reading_time >= $2 
			  AND isr.reading_time <= $3
			ORDER BY isr.reading_time ASC
		)
		SELECT 
			rd.*,
			json_agg(
				json_build_object(
					'id', smt.id,
					'name', smt.name,
					'description', smt.description,
					
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
		FROM reading_details rd
		LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = rd.st_id
		GROUP BY rd.id, rd.tenant_id, rd.asset_sensor_id, rd.sensor_type_id, rd.mac_address,
				rd.location, rd.measurement_data, rd.standard_fields, rd.reading_time,
				rd.created_at, rd.updated_at, rd.data_x, rd.data_y, rd.peak_x, rd.peak_y,
				rd.ppm, rd.label, rd.raw_data, rd.asn_id, rd.asn_asset_id, rd.asn_name,
				rd.asn_status, rd.asn_configuration, rd.st_id, rd.st_name, rd.st_description,
				rd.st_manufacturer, rd.st_model, rd.st_version, rd.st_is_active
		ORDER BY rd.reading_time ASC`

	rows, err := r.DB.QueryContext(ctx, query, assetSensorID, fromTime, toTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get IoT sensor readings in time range: %w", err)
	}
	defer rows.Close()

	return r.scanRowsToResults(rows)
}

// GetAggregatedData retrieves aggregated sensor data for analytics and visualization
func (r *iotSensorReadingRepository) GetAggregatedData(ctx context.Context, assetSensorID uuid.UUID, fromTime, toTime time.Time, interval string) ([]map[string]interface{}, error) {
	// Validate interval
	validIntervals := map[string]bool{
		"5m": true, "15m": true, "30m": true, "1h": true, "6h": true, "12h": true, "1d": true,
	}
	if !validIntervals[interval] {
		return nil, fmt.Errorf("invalid interval: %s. Valid intervals: 5m, 15m, 30m, 1h, 6h, 12h, 1d", interval)
	}

	query := `
		SELECT 
			date_trunc($4, reading_time) as time_bucket,
			COUNT(*) as reading_count,
			AVG(CASE WHEN ppm > 0 THEN ppm END) as avg_ppm,
			MIN(CASE WHEN ppm > 0 THEN ppm END) as min_ppm,
			MAX(CASE WHEN ppm > 0 THEN ppm END) as max_ppm,
			json_agg(
				CASE WHEN measurement_data IS NOT NULL 
				THEN measurement_data 
				ELSE NULL END
			) FILTER (WHERE measurement_data IS NOT NULL) as measurement_samples
		FROM iot_sensor_readings
		WHERE asset_sensor_id = $1 
		  AND reading_time >= $2 
		  AND reading_time <= $3
		GROUP BY time_bucket
		ORDER BY time_bucket ASC`

	rows, err := r.DB.QueryContext(ctx, query, assetSensorID, fromTime, toTime, interval)
	if err != nil {
		return nil, fmt.Errorf("failed to get aggregated data: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var timeBucket time.Time
		var readingCount int
		var avgPPM, minPPM, maxPPM sql.NullFloat64
		var measurementSamples json.RawMessage

		err := rows.Scan(&timeBucket, &readingCount, &avgPPM, &minPPM, &maxPPM, &measurementSamples)
		if err != nil {
			return nil, fmt.Errorf("failed to scan aggregated data: %w", err)
		}

		result := map[string]interface{}{
			"time":          timeBucket,
			"reading_count": readingCount,
		}

		if avgPPM.Valid {
			result["avg_ppm"] = avgPPM.Float64
		}
		if minPPM.Valid {
			result["min_ppm"] = minPPM.Float64
		}
		if maxPPM.Valid {
			result["max_ppm"] = maxPPM.Float64
		}

		if measurementSamples != nil {
			var samples []json.RawMessage
			if err := json.Unmarshal(measurementSamples, &samples); err == nil {
				result["measurement_samples"] = samples
			}
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating aggregated data rows: %w", err)
	}

	return results, nil
}

// GetByAssetSensorID retrieves IoT sensor readings for a specific asset sensor
func (r *iotSensorReadingRepository) GetByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID, limit int) ([]*IoTSensorReadingWithDetails, error) {
	// Validate limit
	if limit <= 0 {
		limit = 100 // Default limit
	}
	if limit > 1000 {
		limit = 1000 // Maximum limit
	}

	// Check if the context has tenant information for multi-tenancy
	tenantID, hasTenantContext := common.GetTenantID(ctx)
	role, hasRoleContext := common.GetUserRole(ctx)

	var query string
	var args []interface{}

	if hasRoleContext && role == "SuperAdmin" && !hasTenantContext {
		// SuperAdmin without tenant context - access all data
		query = `
			WITH reading_details AS (
				SELECT 
					isr.*,
					asn.id as asn_id, asn.asset_id as asn_asset_id, asn.name as asn_name,
					asn.status as asn_status, asn.configuration as asn_configuration,
					st.id as st_id, st.name as st_name, st.description as st_description,
					st.manufacturer as st_manufacturer, st.model as st_model,
					st.version as st_version, st.is_active as st_is_active
				FROM iot_sensor_readings isr
				JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
				JOIN sensor_types st ON isr.sensor_type_id = st.id
				WHERE isr.asset_sensor_id = $1
				ORDER BY isr.reading_time DESC
				LIMIT $2
			)
			SELECT 
				rd.*,
				json_agg(
					json_build_object(
						'id', smt.id,
						'name', smt.name,
						'description', smt.description,
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
			FROM reading_details rd
			LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = rd.st_id
			GROUP BY rd.id, rd.tenant_id, rd.asset_sensor_id, rd.sensor_type_id, rd.mac_address,
					rd.location_id, rd.location_name, rd.reading_time, rd.created_at, rd.updated_at,
					rd.asn_id, rd.asn_asset_id, rd.asn_name, rd.asn_status, rd.asn_configuration,
					rd.st_id, rd.st_name, rd.st_description, rd.st_manufacturer, rd.st_model,
					rd.st_version, rd.st_is_active
			ORDER BY rd.reading_time DESC`
		args = []interface{}{assetSensorID, limit}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			WITH reading_details AS (
				SELECT 
					isr.*,
					asn.id as asn_id, asn.asset_id as asn_asset_id, asn.name as asn_name,
					asn.status as asn_status, asn.configuration as asn_configuration,
					st.id as st_id, st.name as st_name, st.description as st_description,
					st.manufacturer as st_manufacturer, st.model as st_model,
					st.version as st_version, st.is_active as st_is_active
				FROM iot_sensor_readings isr
				JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
				JOIN sensor_types st ON isr.sensor_type_id = st.id
				WHERE isr.asset_sensor_id = $1 AND isr.tenant_id = $2
				ORDER BY isr.reading_time DESC
				LIMIT $3
			)
			SELECT 
				rd.*,
				json_agg(
					json_build_object(
						'id', smt.id,
						'name', smt.name,
						'description', smt.description,
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
			FROM reading_details rd
			LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = rd.st_id
			GROUP BY rd.id, rd.tenant_id, rd.asset_sensor_id, rd.sensor_type_id, rd.mac_address,
					rd.location_id, rd.location_name, rd.reading_time, rd.created_at, rd.updated_at,
					rd.asn_id, rd.asn_asset_id, rd.asn_name, rd.asn_status, rd.asn_configuration,
					rd.st_id, rd.st_name, rd.st_description, rd.st_manufacturer, rd.st_model,
					rd.st_version, rd.st_is_active
			ORDER BY rd.reading_time DESC`
		args = []interface{}{assetSensorID, tenantID, limit}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get IoT sensor readings: %w", err)
	}
	defer rows.Close()

	return r.scanRowsToResults(rows)
}

// GetByMacAddress retrieves IoT sensor readings for a specific MAC address
func (r *iotSensorReadingRepository) GetByMacAddress(ctx context.Context, macAddress string, limit int) ([]*IoTSensorReadingWithDetails, error) {
	// Validate limit
	if limit <= 0 {
		limit = 100 // Default limit
	}
	if limit > 1000 {
		limit = 1000 // Maximum limit
	}

	// Check if the context has tenant information for multi-tenancy
	tenantID, hasTenantContext := common.GetTenantID(ctx)
	role, hasRoleContext := common.GetUserRole(ctx)

	var query string
	var args []interface{}

	if hasRoleContext && role == "SuperAdmin" && !hasTenantContext {
		// SuperAdmin without tenant context - access all data
		query = `
			WITH reading_details AS (
				SELECT 
					isr.*,
					asn.id as asn_id, asn.asset_id as asn_asset_id, asn.name as asn_name,
					asn.status as asn_status, asn.configuration as asn_configuration,
					st.id as st_id, st.name as st_name, st.description as st_description,
					st.manufacturer as st_manufacturer, st.model as st_model,
					st.version as st_version, st.is_active as st_is_active
				FROM iot_sensor_readings isr
				JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
				JOIN sensor_types st ON isr.sensor_type_id = st.id
				WHERE isr.mac_address = $1
				ORDER BY isr.reading_time DESC
				LIMIT $2
			)
			SELECT 
				rd.*,
				json_agg(
					json_build_object(
						'id', smt.id,
						'name', smt.name,
						'description', smt.description,
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
			FROM reading_details rd
			LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = rd.st_id
			GROUP BY rd.id, rd.tenant_id, rd.asset_sensor_id, rd.sensor_type_id, rd.mac_address,
					rd.location_id, rd.location_name, rd.reading_time, rd.created_at, rd.updated_at,
					rd.asn_id, rd.asn_asset_id, rd.asn_name, rd.asn_status, rd.asn_configuration,
					rd.st_id, rd.st_name, rd.st_description, rd.st_manufacturer, rd.st_model,
					rd.st_version, rd.st_is_active
			ORDER BY rd.reading_time DESC`
		args = []interface{}{macAddress, limit}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			WITH reading_details AS (
				SELECT 
					isr.*,
					asn.id as asn_id, asn.asset_id as asn_asset_id, asn.name as asn_name,
					asn.status as asn_status, asn.configuration as asn_configuration,
					st.id as st_id, st.name as st_name, st.description as st_description,
					st.manufacturer as st_manufacturer, st.model as st_model,
					st.version as st_version, st.is_active as st_is_active
				FROM iot_sensor_readings isr
				JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
				JOIN sensor_types st ON isr.sensor_type_id = st.id
				WHERE isr.mac_address = $1 AND isr.tenant_id = $2
				ORDER BY isr.reading_time DESC
				LIMIT $3
			)
			SELECT 
				rd.*,
				json_agg(
					json_build_object(
						'id', smt.id,
						'name', smt.name,
						'description', smt.description,
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
			FROM reading_details rd
			LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = rd.st_id
			GROUP BY rd.id, rd.tenant_id, rd.asset_sensor_id, rd.sensor_type_id, rd.mac_address,
					rd.location_id, rd.location_name, rd.reading_time, rd.created_at, rd.updated_at,
					rd.asn_id, rd.asn_asset_id, rd.asn_name, rd.asn_status, rd.asn_configuration,
					rd.st_id, rd.st_name, rd.st_description, rd.st_manufacturer, rd.st_model,
					rd.st_version, rd.st_is_active
			ORDER BY rd.reading_time DESC`
		args = []interface{}{macAddress, tenantID, limit}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get IoT sensor readings: %w", err)
	}
	defer rows.Close()

	return r.scanRowsToResults(rows)
}

// GetBySensorTypeID retrieves IoT sensor readings for a specific sensor type
func (r *iotSensorReadingRepository) GetBySensorTypeID(ctx context.Context, sensorTypeID uuid.UUID, limit int) ([]*IoTSensorReadingWithDetails, error) {
	// Validate limit
	if limit <= 0 {
		limit = 100 // Default limit
	}
	if limit > 1000 {
		limit = 1000 // Maximum limit
	}

	// Check if the context has tenant information for multi-tenancy
	tenantID, hasTenantContext := common.GetTenantID(ctx)
	role, hasRoleContext := common.GetUserRole(ctx)

	var query string
	var args []interface{}

	if hasRoleContext && role == "SuperAdmin" && !hasTenantContext {
		// SuperAdmin without tenant context - access all data
		query = `
			WITH reading_details AS (
				SELECT 
					isr.*,
					asn.id as asn_id, asn.asset_id as asn_asset_id, asn.name as asn_name,
					asn.status as asn_status, asn.configuration as asn_configuration,
					st.id as st_id, st.name as st_name, st.description as st_description,
					st.manufacturer as st_manufacturer, st.model as st_model,
					st.version as st_version, st.is_active as st_is_active
				FROM iot_sensor_readings isr
				JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
				JOIN sensor_types st ON isr.sensor_type_id = st.id
				WHERE isr.sensor_type_id = $1
				ORDER BY isr.reading_time DESC
				LIMIT $2
			)
			SELECT 
				rd.*,
				json_agg(
					json_build_object(
						'id', smt.id,
						'name', smt.name,
						'description', smt.description,
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
			FROM reading_details rd
			LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = rd.st_id
			GROUP BY rd.id, rd.tenant_id, rd.asset_sensor_id, rd.sensor_type_id, rd.mac_address,
					rd.location_id, rd.location_name, rd.reading_time, rd.created_at, rd.updated_at,
					rd.asn_id, rd.asn_asset_id, rd.asn_name, rd.asn_status, rd.asn_configuration,
					rd.st_id, rd.st_name, rd.st_description, rd.st_manufacturer, rd.st_model,
					rd.st_version, rd.st_is_active
			ORDER BY rd.reading_time DESC`
		args = []interface{}{sensorTypeID, limit}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			WITH reading_details AS (
				SELECT 
					isr.*,
					asn.id as asn_id, asn.asset_id as asn_asset_id, asn.name as asn_name,
					asn.status as asn_status, asn.configuration as asn_configuration,
					st.id as st_id, st.name as st_name, st.description as st_description,
					st.manufacturer as st_manufacturer, st.model as st_model,
					st.version as st_version, st.is_active as st_is_active
				FROM iot_sensor_readings isr
				JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
				JOIN sensor_types st ON isr.sensor_type_id = st.id
				WHERE isr.sensor_type_id = $1 AND isr.tenant_id = $2
				ORDER BY isr.reading_time DESC
				LIMIT $3
			)
			SELECT 
				rd.*,
				json_agg(
					json_build_object(
						'id', smt.id,
						'name', smt.name,
						'description', smt.description,
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
			FROM reading_details rd
			LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = rd.st_id
			GROUP BY rd.id, rd.tenant_id, rd.asset_sensor_id, rd.sensor_type_id, rd.mac_address,
					rd.location_id, rd.location_name, rd.reading_time, rd.created_at, rd.updated_at,
					rd.asn_id, rd.asn_asset_id, rd.asn_name, rd.asn_status, rd.asn_configuration,
					rd.st_id, rd.st_name, rd.st_description, rd.st_manufacturer, rd.st_model,
					rd.st_version, rd.st_is_active
			ORDER BY rd.reading_time DESC`
		args = []interface{}{sensorTypeID, tenantID, limit}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get IoT sensor readings: %w", err)
	}
	defer rows.Close()

	return r.scanRowsToResults(rows)
}

// ValidateAndCreate validates a reading against measurement type schemas and creates it
func (r *iotSensorReadingRepository) ValidateAndCreate(ctx context.Context, reading *entity.IoTSensorReading) (bool, []string, error) {
	// Get measurement types and their fields for validation
	measurementTypesQuery := `
		SELECT 
			smt.id, smt.name, smt.description, smt.properties_schema,
			COALESCE(
				(SELECT json_agg(
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
				WHERE smf.sensor_measurement_type_id = smt.id), '[]'::json
			) as fields
		FROM sensor_measurement_types smt
		WHERE smt.sensor_type_id = $1 AND smt.is_active = true`

	rows, err := r.DB.QueryContext(ctx, measurementTypesQuery, reading.SensorTypeID)
	if err != nil {
		return false, nil, fmt.Errorf("failed to get measurement types and fields: %w", err)
	}
	defer rows.Close()

	var validationErrors []string
	var measurementTypes []entity.SensorMeasurementType

	for rows.Next() {
		var mt entity.SensorMeasurementType
		var fieldsJSON []byte

		err := rows.Scan(
			&mt.ID,
			&mt.Name,
			&mt.Description,
			&mt.PropertiesSchema,
			&fieldsJSON,
		)
		if err != nil {
			return false, nil, fmt.Errorf("failed to scan measurement type: %w", err)
		}

		// Parse fields JSON
		if len(fieldsJSON) > 0 && string(fieldsJSON) != "[]" {
			var fieldDefs []map[string]interface{}
			if err := json.Unmarshal(fieldsJSON, &fieldDefs); err != nil {
				return false, nil, fmt.Errorf("failed to parse measurement fields: %w", err)
			}

			// Convert to SensorMeasurementField entities
			for _, fieldDef := range fieldDefs {
				field := entity.SensorMeasurementField{
					Name:     fieldDef["name"].(string),
					Label:    fieldDef["label"].(string),
					DataType: entity.MeasurementDataType(fieldDef["data_type"].(string)),
					Required: fieldDef["required"].(bool),
				}

				if desc, ok := fieldDef["description"].(string); ok && desc != "" {
					field.Description = desc
				}
				if unit, ok := fieldDef["unit"].(string); ok && unit != "" {
					field.Unit = unit
				}
				if minVal, ok := fieldDef["min"].(float64); ok {
					field.Min = &minVal
				}
				if maxVal, ok := fieldDef["max"].(float64); ok {
					field.Max = &maxVal
				}

				mt.Fields = append(mt.Fields, field)
			}
		}

		measurementTypes = append(measurementTypes, mt)
	}

	// If no measurement types found, skip validation but warn
	if len(measurementTypes) == 0 {
		validationErrors = append(validationErrors, "Warning: No measurement types found for sensor type, skipping validation")
	} else {
		// Validate against each measurement type
		for _, mt := range measurementTypes {
			valid, errors := reading.ValidateAgainstMeasurementType(&mt)
			if !valid {
				validationErrors = append(validationErrors, errors...)
			}
		}
	}

	// If validation passes, create the reading
	if len(validationErrors) == 0 {
		err := r.Create(ctx, reading)
		if err != nil {
			return false, validationErrors, fmt.Errorf("failed to create validated reading: %w", err)
		}
	}

	return len(validationErrors) == 0, validationErrors, nil
}

// CreateFlexible inserts a new flexible IoT sensor reading with measurement data
func (r *iotSensorReadingRepository) CreateFlexible(ctx context.Context, reading *entity.IoTSensorReadingFlexible) error {
	// Generate new UUID if not set
	if reading.ID == uuid.Nil {
		reading.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now()
	if reading.CreatedAt.IsZero() {
		reading.CreatedAt = now
	}
	reading.UpdatedAt = &now

	query := `
		INSERT INTO iot_sensor_readings (
			id, tenant_id, asset_sensor_id, sensor_type_id, mac_address,
			location_id, location_name, measurement_type, measurement_label,
			measurement_unit, numeric_value, text_value, boolean_value,
			data_source, original_field_name, reading_time, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 
			$11, $12, $13, $14, $15, $16, $17, $18
		)`

	_, err := r.DB.ExecContext(ctx, query,
		reading.ID,
		reading.TenantID,
		reading.AssetSensorID,
		reading.SensorTypeID,
		reading.MacAddress,
		reading.LocationID,
		reading.LocationName,
		reading.MeasurementType,
		reading.MeasurementLabel,
		reading.MeasurementUnit,
		reading.NumericValue,
		reading.TextValue,
		reading.BooleanValue,
		reading.DataSource,
		reading.OriginalFieldName,
		reading.ReadingTime,
		reading.CreatedAt,
		reading.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create flexible IoT sensor reading: %w", err)
	}

	return nil
}

// CreateFlexibleBatch inserts multiple flexible IoT sensor readings efficiently
func (r *iotSensorReadingRepository) CreateFlexibleBatch(ctx context.Context, readings []*entity.IoTSensorReadingFlexible) error {
	if len(readings) == 0 {
		return nil
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO iot_sensor_readings (
			id, tenant_id, asset_sensor_id, sensor_type_id, mac_address,
			location_id, location_name, measurement_type, measurement_label,
			measurement_unit, numeric_value, text_value, boolean_value,
			data_source, original_field_name, reading_time, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18
		)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()

	for _, reading := range readings {
		// Set ID if not already set
		if reading.ID == uuid.Nil {
			reading.ID = uuid.New()
		}

		// Set timestamps
		if reading.CreatedAt.IsZero() {
			reading.CreatedAt = now
		}
		reading.UpdatedAt = &now

		_, err = stmt.ExecContext(ctx,
			reading.ID,
			reading.TenantID,
			reading.AssetSensorID,
			reading.SensorTypeID,
			reading.MacAddress,
			reading.LocationID,
			reading.LocationName,
			reading.MeasurementType,
			reading.MeasurementLabel,
			reading.MeasurementUnit,
			reading.NumericValue,
			reading.TextValue,
			reading.BooleanValue,
			reading.DataSource,
			reading.OriginalFieldName,
			reading.ReadingTime,
			reading.CreatedAt,
			reading.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert flexible reading: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetFlexibleByID retrieves a flexible IoT sensor reading with measurement data
func (r *iotSensorReadingRepository) GetFlexibleByID(ctx context.Context, id uuid.UUID) (*entity.IoTSensorReadingFlexible, error) {
	var reading entity.IoTSensorReadingFlexible

	query := `
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address,
			   location_id, location_name, measurement_type, measurement_label,
			   measurement_unit, numeric_value, text_value, boolean_value,
			   data_source, original_field_name, reading_time, created_at, updated_at
		FROM iot_sensor_readings
		WHERE id = $1`

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&reading.ID,
		&reading.TenantID,
		&reading.AssetSensorID,
		&reading.SensorTypeID,
		&reading.MacAddress,
		&reading.LocationID,
		&reading.LocationName,
		&reading.MeasurementType,
		&reading.MeasurementLabel,
		&reading.MeasurementUnit,
		&reading.NumericValue,
		&reading.TextValue,
		&reading.BooleanValue,
		&reading.DataSource,
		&reading.OriginalFieldName,
		&reading.ReadingTime,
		&reading.CreatedAt,
		&reading.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("flexible IoT sensor reading not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get flexible IoT sensor reading: %w", err)
	}

	return &reading, nil
}

// ParseTextToFlexibleReading converts text data to a flexible IoT sensor reading
func (r *iotSensorReadingRepository) ParseTextToFlexibleReading(ctx context.Context, textData, assetSensorID, sensorTypeID, macAddress string) (*entity.IoTSensorReadingFlexible, error) {
	assetSensorUUID, err := uuid.Parse(assetSensorID)
	if err != nil {
		return nil, fmt.Errorf("invalid asset sensor ID: %w", err)
	}

	sensorTypeUUID, err := uuid.Parse(sensorTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid sensor type ID: %w", err)
	}

	// Parse text to flexible readings
	readings, err := entity.ParseTextToFlexibleReadings(textData, assetSensorUUID, sensorTypeUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse text data: %w", err)
	}

	if len(readings) == 0 {
		return nil, fmt.Errorf("no valid readings found in text data")
	}

	// Return the first reading for single-reading scenario
	// You might want to modify this logic based on your needs
	firstReading := readings[0]
	if macAddress != "" {
		firstReading.MacAddress = &macAddress
	}

	return firstReading, nil
}

// GetFlexibleByAssetSensorID retrieves flexible readings by asset sensor ID
func (r *iotSensorReadingRepository) GetFlexibleByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID, limit, offset int) ([]*entity.IoTSensorReadingFlexible, error) {
	var readings []*entity.IoTSensorReadingFlexible

	query := `
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address,
			   location_id, location_name, measurement_type, measurement_label,
			   measurement_unit, numeric_value, text_value, boolean_value,
			   data_source, original_field_name, reading_time, created_at, updated_at
		FROM iot_sensor_readings
		WHERE asset_sensor_id = $1
		ORDER BY reading_time DESC, created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.DB.QueryContext(ctx, query, assetSensorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query flexible readings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reading entity.IoTSensorReadingFlexible
		err := rows.Scan(
			&reading.ID,
			&reading.TenantID,
			&reading.AssetSensorID,
			&reading.SensorTypeID,
			&reading.MacAddress,
			&reading.LocationID,
			&reading.LocationName,
			&reading.MeasurementType,
			&reading.MeasurementLabel,
			&reading.MeasurementUnit,
			&reading.NumericValue,
			&reading.TextValue,
			&reading.BooleanValue,
			&reading.DataSource,
			&reading.OriginalFieldName,
			&reading.ReadingTime,
			&reading.CreatedAt,
			&reading.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reading: %w", err)
		}
		readings = append(readings, &reading)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return readings, nil
}

// ListFlexible retrieves flexible IoT sensor readings with filtering and pagination
func (r *iotSensorReadingRepository) ListFlexible(ctx context.Context, req IoTSensorReadingListRequest) ([]*entity.IoTSensorReadingFlexible, int, error) {
	// Validate pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Check if the context has tenant information for multi-tenancy
	tenantID, hasTenantContext := common.GetTenantID(ctx)
	role, hasRoleContext := common.GetUserRole(ctx)

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argCount := 0

	// Tenant filtering
	if hasRoleContext && (role == "SuperAdmin" || role == "SUPERADMIN") && !hasTenantContext {
		// SuperAdmin without tenant context - no tenant filter (access all data)
	} else {
		// Regular users or SuperAdmin with tenant context - filter by tenant
		argCount++
		conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argCount))
		args = append(args, tenantID)
	}

	// Asset sensor filtering
	if req.AssetSensorID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("asset_sensor_id = $%d", argCount))
		args = append(args, *req.AssetSensorID)
	}

	// Sensor type filtering
	if req.SensorTypeID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("sensor_type_id = $%d", argCount))
		args = append(args, *req.SensorTypeID)
	}

	// MAC address filtering
	if req.MacAddress != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("mac_address = $%d", argCount))
		args = append(args, *req.MacAddress)
	}

	// Location filtering
	if req.LocationID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("location_id = $%d", argCount))
		args = append(args, *req.LocationID)
	}

	// Time range filtering
	if req.FromTime != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("reading_time >= $%d", argCount))
		args = append(args, *req.FromTime)
	}

	if req.ToTime != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("reading_time <= $%d", argCount))
		args = append(args, *req.ToTime)
	}

	// Build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM iot_sensor_readings %s", whereClause)
	var total int
	err := r.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count flexible readings: %w", err)
	}

	// Data query with pagination
	offset := (req.Page - 1) * req.PageSize
	dataQuery := fmt.Sprintf(`
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address,
			   location_id, location_name, measurement_type, measurement_label,
			   measurement_unit, numeric_value, text_value, boolean_value,
			   data_source, original_field_name, reading_time, created_at, updated_at
		FROM iot_sensor_readings
		%s
		ORDER BY reading_time DESC, created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argCount+1, argCount+2)

	// Add pagination args
	args = append(args, req.PageSize, offset)

	rows, err := r.DB.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query flexible readings: %w", err)
	}
	defer rows.Close()

	var readings []*entity.IoTSensorReadingFlexible
	for rows.Next() {
		var reading entity.IoTSensorReadingFlexible
		err := rows.Scan(
			&reading.ID,
			&reading.TenantID,
			&reading.AssetSensorID,
			&reading.SensorTypeID,
			&reading.MacAddress,
			&reading.LocationID,
			&reading.LocationName,
			&reading.MeasurementType,
			&reading.MeasurementLabel,
			&reading.MeasurementUnit,
			&reading.NumericValue,
			&reading.TextValue,
			&reading.BooleanValue,
			&reading.DataSource,
			&reading.OriginalFieldName,
			&reading.ReadingTime,
			&reading.CreatedAt,
			&reading.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan flexible reading: %w", err)
		}
		readings = append(readings, &reading)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return readings, total, nil
}

// scanRowsToResults scans database rows into IoTSensorReadingWithDetails structure
func (r *iotSensorReadingRepository) scanRowsToResults(rows *sql.Rows) ([]*IoTSensorReadingWithDetails, error) {
	var readings []*IoTSensorReadingWithDetails

	for rows.Next() {
		var reading IoTSensorReadingWithDetails
		reading.IoTSensorReading = &entity.IoTSensorReading{}
		var measurementTypesJSON []byte

		err := rows.Scan(
			&reading.ID,
			&reading.TenantID,
			&reading.AssetSensorID,
			&reading.SensorTypeID,
			&reading.MacAddress,
			&reading.LocationID,
			&reading.LocationName,
			&reading.ReadingTime,
			&reading.CreatedAt,
			&reading.UpdatedAt,
			&reading.AssetSensor.ID,
			&reading.AssetSensor.AssetID,
			&reading.AssetSensor.Name,
			&reading.AssetSensor.Status,
			&reading.AssetSensor.Configuration,
			&reading.SensorType.ID,
			&reading.SensorType.Name,
			&reading.SensorType.Description,
			&reading.SensorType.Manufacturer,
			&reading.SensorType.Model,
			&reading.SensorType.Version,
			&reading.SensorType.IsActive,
			&measurementTypesJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Parse measurement types JSON
		if measurementTypesJSON != nil {
			err = json.Unmarshal(measurementTypesJSON, &reading.MeasurementTypes)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal measurement types: %w", err)
			}
		}

		readings = append(readings, &reading)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return readings, nil
}

// List retrieves IoT sensor readings with filtering and pagination
func (r *iotSensorReadingRepository) List(ctx context.Context, req IoTSensorReadingListRequest) ([]*IoTSensorReadingWithDetails, int, error) {
	// Validate pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	offset := (req.Page - 1) * req.PageSize

	// Check if the context has tenant information for multi-tenancy
	tenantID, hasTenantContext := common.GetTenantID(ctx)
	role, hasRoleContext := common.GetUserRole(ctx)

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Tenant filtering
	if hasRoleContext && (role == "SuperAdmin" || role == "SUPERADMIN") && !hasTenantContext {
		// SuperAdmin tanpa tenant context - tidak ada filter tenant
	} else {
		// Regular users atau SuperAdmin dengan tenant ID - filter by tenant
		conditions = append(conditions, fmt.Sprintf("isr.tenant_id = $%d", argIndex))
		args = append(args, tenantID)
		argIndex++
	}

	// Asset sensor filtering
	if req.AssetSensorID != nil {
		conditions = append(conditions, fmt.Sprintf("isr.asset_sensor_id = $%d", argIndex))
		args = append(args, *req.AssetSensorID)
		argIndex++
	}

	// Sensor type filtering
	if req.SensorTypeID != nil {
		conditions = append(conditions, fmt.Sprintf("isr.sensor_type_id = $%d", argIndex))
		args = append(args, *req.SensorTypeID)
		argIndex++
	}

	// MAC address filtering
	if req.MacAddress != nil && *req.MacAddress != "" {
		conditions = append(conditions, fmt.Sprintf("isr.mac_address = $%d", argIndex))
		args = append(args, *req.MacAddress)
		argIndex++
	}

	// Location filtering
	if req.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("isr.location_id = $%d", argIndex))
		args = append(args, *req.LocationID)
		argIndex++
	}

	// Time range filtering
	if req.FromTime != nil {
		conditions = append(conditions, fmt.Sprintf("isr.reading_time >= $%d", argIndex))
		args = append(args, *req.FromTime)
		argIndex++
	}
	if req.ToTime != nil {
		conditions = append(conditions, fmt.Sprintf("isr.reading_time <= $%d", argIndex))
		args = append(args, *req.ToTime)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM iot_sensor_readings isr
		JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
		JOIN sensor_types st ON isr.sensor_type_id = st.id
		%s`, whereClause)

	var totalCount int
	err := r.DB.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count IoT sensor readings: %w", err)
	}

	// Data query with pagination
	dataQuery := fmt.Sprintf(`
		WITH reading_details AS (
			SELECT 
				isr.id, isr.tenant_id, isr.asset_sensor_id, isr.sensor_type_id, isr.mac_address,
				isr.location_id, isr.location_name, isr.reading_time, isr.created_at, isr.updated_at,
				asn.id as asn_id, asn.asset_id as asn_asset_id, asn.name as asn_name, asn.status as asn_status, asn.configuration as asn_configuration,
				st.id as st_id, st.name as st_name, st.description as st_description, st.manufacturer as st_manufacturer, st.model as st_model, st.version as st_version, st.is_active as st_is_active
			FROM iot_sensor_readings isr
			JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
			JOIN sensor_types st ON isr.sensor_type_id = st.id
			%s
			ORDER BY isr.reading_time DESC
			LIMIT $%d OFFSET $%d
		)
		SELECT 
			reading_details.id,
			reading_details.tenant_id,
			reading_details.asset_sensor_id,
			reading_details.sensor_type_id,
			reading_details.mac_address,
			reading_details.location_id,
			reading_details.location_name,
			reading_details.reading_time,
			reading_details.created_at,
			reading_details.updated_at,
			reading_details.asn_id,
			reading_details.asn_asset_id,
			reading_details.asn_name,
			reading_details.asn_status,
			reading_details.asn_configuration,
			reading_details.st_id,
			reading_details.st_name,
			reading_details.st_description,
			reading_details.st_manufacturer,
			reading_details.st_model,
			reading_details.st_version,
			reading_details.st_is_active,
			json_agg(
				json_build_object(
					'id', smt.id,
					'name', smt.name,
					'description', smt.description,
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
		FROM reading_details
		LEFT JOIN sensor_measurement_types smt ON smt.sensor_type_id = reading_details.st_id
		GROUP BY reading_details.id, reading_details.tenant_id, reading_details.asset_sensor_id, reading_details.sensor_type_id, reading_details.mac_address,
			reading_details.location_id, reading_details.location_name, reading_details.reading_time, reading_details.created_at, reading_details.updated_at,
			reading_details.asn_id, reading_details.asn_asset_id, reading_details.asn_name, reading_details.asn_status, reading_details.asn_configuration,
			reading_details.st_id, reading_details.st_name, reading_details.st_description, reading_details.st_manufacturer, reading_details.st_model,
			reading_details.st_version, reading_details.st_is_active
		ORDER BY reading_details.reading_time DESC`, whereClause, argIndex, argIndex+1)

	// Add pagination args
	paginationArgs := append(args, req.PageSize, offset)

	rows, err := r.DB.QueryContext(ctx, dataQuery, paginationArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get IoT sensor readings: %w", err)
	}
	defer rows.Close()

	readings, err := r.scanRowsToResults(rows)
	if err != nil {
		return nil, 0, err
	}

	return readings, totalCount, nil
}

// GetAssetSensorsBySensorType retrieves asset sensors with location information by sensor type
func (r *iotSensorReadingRepository) GetAssetSensorsBySensorType(ctx context.Context, sensorTypeID uuid.UUID) ([]dto.AssetSensorLocationInfo, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, fmt.Errorf("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access all sensors
		query = `
			SELECT 
				asn.id as asset_sensor_id,
				asn.name as asset_sensor_name,
				a.id as asset_id,
				a.name as asset_name,
				l.id as location_id,
				l.name as location_name
			FROM asset_sensors asn
			JOIN assets a ON asn.asset_id = a.id
			JOIN locations l ON a.location_id = l.id
			WHERE asn.sensor_type_id = $1
			AND asn.status = 'active'
			ORDER BY l.name, a.name, asn.name`
		args = []interface{}{sensorTypeID}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT 
				asn.id as asset_sensor_id,
				asn.name as asset_sensor_name,
				a.id as asset_id,
				a.name as asset_name,
				l.id as location_id,
				l.name as location_name
			FROM asset_sensors asn
			JOIN assets a ON asn.asset_id = a.id
			JOIN locations l ON a.location_id = l.id
			WHERE asn.sensor_type_id = $1
			AND asn.tenant_id = $2
			AND asn.status = 'active'
			ORDER BY l.name, a.name, asn.name`
		args = []interface{}{sensorTypeID, tenantID}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset sensors by sensor type: %w", err)
	}
	defer rows.Close()

	var result []dto.AssetSensorLocationInfo
	for rows.Next() {
		var info dto.AssetSensorLocationInfo
		err := rows.Scan(
			&info.AssetSensorID,
			&info.AssetSensorName,
			&info.AssetID,
			&info.AssetName,
			&info.LocationID,
			&info.LocationName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset sensor location info: %w", err)
		}
		result = append(result, info)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}

func (r *iotSensorReadingRepository) GetDB() *sql.DB {
	return r.DB
}
