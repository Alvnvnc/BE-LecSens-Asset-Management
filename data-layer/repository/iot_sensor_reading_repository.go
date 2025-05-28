package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// IoTSensorReadingRepository defines the interface for IoT sensor reading data operations
type IoTSensorReadingRepository interface {
	Create(ctx context.Context, reading *entity.IoTSensorReading) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.IoTSensorReading, error)
	List(ctx context.Context, tenantID *uuid.UUID, page, pageSize int) ([]*entity.IoTSensorReading, error)
	ListByAssetSensor(ctx context.Context, assetSensorID uuid.UUID, page, pageSize int) ([]*entity.IoTSensorReading, error)
	ListBySensorType(ctx context.Context, sensorTypeID uuid.UUID, page, pageSize int) ([]*entity.IoTSensorReading, error)
	ListByMacAddress(ctx context.Context, macAddress string, page, pageSize int) ([]*entity.IoTSensorReading, error)
	ListByTimeRange(ctx context.Context, startTime, endTime time.Time, page, pageSize int) ([]*entity.IoTSensorReading, error)
	Update(ctx context.Context, reading *entity.IoTSensorReading) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetLatestByMacAddress(ctx context.Context, macAddress string) (*entity.IoTSensorReading, error)
	GetCountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)
	GetAverageValueByTimeRange(ctx context.Context, sensorTypeID uuid.UUID, startTime, endTime time.Time) (float64, error)
}

// iotSensorReadingRepository handles database operations for IoT sensor readings
type iotSensorReadingRepository struct {
	db *sql.DB
}

// NewIoTSensorReadingRepository creates a new IoTSensorReadingRepository
func NewIoTSensorReadingRepository(db *sql.DB) IoTSensorReadingRepository {
	return &iotSensorReadingRepository{db: db}
}

// Create inserts a new IoT sensor reading into the database
func (r *iotSensorReadingRepository) Create(ctx context.Context, reading *entity.IoTSensorReading) error {
	query := `
		INSERT INTO iot_sensor_readings (
			tenant_id, asset_sensor_id, sensor_type_id, mac_address, location,
			measurement_data, standard_fields, reading_time, created_at, updated_at,
			data_x, data_y, peak_x, peak_y, ppm, label, raw_data
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		) RETURNING id`

	var id uuid.UUID
	now := time.Now()

	err := r.db.QueryRowContext(
		ctx,
		query,
		reading.TenantID,
		reading.AssetSensorID,
		reading.SensorTypeID,
		reading.MacAddress,
		reading.Location,
		reading.MeasurementData,
		reading.StandardFields,
		reading.ReadingTime,
		now,
		reading.UpdatedAt,
		reading.DataX,
		reading.DataY,
		reading.PeakX,
		reading.PeakY,
		reading.PPM,
		reading.Label,
		reading.RawData,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to create IoT sensor reading: %w", err)
	}

	reading.ID = id
	reading.CreatedAt = now
	return nil
}

// GetByID retrieves an IoT sensor reading by its ID
func (r *iotSensorReadingRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.IoTSensorReading, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address, location,
			   measurement_data, standard_fields, reading_time, created_at, updated_at,
			   data_x, data_y, peak_x, peak_y, ppm, label, raw_data
		FROM iot_sensor_readings
		WHERE id = $1`

	var reading entity.IoTSensorReading
	var tenantID, assetSensorID sql.NullString
	var measurementData, standardFields, dataX, dataY, peakX, peakY, rawData []byte
	var ppm sql.NullFloat64
	var label sql.NullString
	var updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&reading.ID,
		&tenantID,
		&assetSensorID,
		&reading.SensorTypeID,
		&reading.MacAddress,
		&reading.Location,
		&measurementData,
		&standardFields,
		&reading.ReadingTime,
		&reading.CreatedAt,
		&updatedAt,
		&dataX,
		&dataY,
		&peakX,
		&peakY,
		&ppm,
		&label,
		&rawData,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get IoT sensor reading: %w", err)
	}

	// Handle nullable fields
	if tenantID.Valid {
		parsedID, err := uuid.Parse(tenantID.String)
		if err != nil {
			return nil, fmt.Errorf("invalid tenant ID format: %w", err)
		}
		reading.TenantID = parsedID
	}

	if assetSensorID.Valid {
		parsedID, err := uuid.Parse(assetSensorID.String)
		if err != nil {
			return nil, fmt.Errorf("invalid asset sensor ID format: %w", err)
		}
		reading.AssetSensorID = parsedID
	}

	if updatedAt.Valid {
		reading.UpdatedAt = &updatedAt.Time
	}

	if ppm.Valid {
		reading.PPM = ppm.Float64
	}

	if label.Valid {
		reading.Label = label.String
	}

	// Handle JSON fields
	if measurementData != nil {
		reading.MeasurementData = json.RawMessage(measurementData)
	}
	if standardFields != nil {
		reading.StandardFields = json.RawMessage(standardFields)
	}
	if dataX != nil {
		reading.DataX = json.RawMessage(dataX)
	}
	if dataY != nil {
		reading.DataY = json.RawMessage(dataY)
	}
	if peakX != nil {
		reading.PeakX = json.RawMessage(peakX)
	}
	if peakY != nil {
		reading.PeakY = json.RawMessage(peakY)
	}
	if rawData != nil {
		reading.RawData = json.RawMessage(rawData)
	}

	return &reading, nil
}

// List retrieves a paginated list of IoT sensor readings
func (r *iotSensorReadingRepository) List(ctx context.Context, tenantID *uuid.UUID, page, pageSize int) ([]*entity.IoTSensorReading, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address, location,
			   measurement_data, standard_fields, reading_time, created_at, updated_at,
			   data_x, data_y, peak_x, peak_y, ppm, label, raw_data
		FROM iot_sensor_readings
		WHERE ($1::uuid IS NULL OR tenant_id = $1)
		ORDER BY reading_time DESC
		LIMIT $2 OFFSET $3`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, tenantID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list IoT sensor readings: %w", err)
	}
	defer rows.Close()

	return r.scanReadings(rows)
}

// ListByAssetSensor retrieves IoT sensor readings for a specific asset sensor
func (r *iotSensorReadingRepository) ListByAssetSensor(ctx context.Context, assetSensorID uuid.UUID, page, pageSize int) ([]*entity.IoTSensorReading, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address, location,
			   measurement_data, standard_fields, reading_time, created_at, updated_at,
			   data_x, data_y, peak_x, peak_y, ppm, label, raw_data
		FROM iot_sensor_readings
		WHERE asset_sensor_id = $1
		ORDER BY reading_time DESC
		LIMIT $2 OFFSET $3`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, assetSensorID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list IoT sensor readings by asset sensor: %w", err)
	}
	defer rows.Close()

	return r.scanReadings(rows)
}

// ListBySensorType retrieves IoT sensor readings for a specific sensor type
func (r *iotSensorReadingRepository) ListBySensorType(ctx context.Context, sensorTypeID uuid.UUID, page, pageSize int) ([]*entity.IoTSensorReading, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address, location,
			   measurement_data, standard_fields, reading_time, created_at, updated_at,
			   data_x, data_y, peak_x, peak_y, ppm, label, raw_data
		FROM iot_sensor_readings
		WHERE sensor_type_id = $1
		ORDER BY reading_time DESC
		LIMIT $2 OFFSET $3`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, sensorTypeID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list IoT sensor readings by sensor type: %w", err)
	}
	defer rows.Close()

	return r.scanReadings(rows)
}

// ListByMacAddress retrieves IoT sensor readings for a specific MAC address
func (r *iotSensorReadingRepository) ListByMacAddress(ctx context.Context, macAddress string, page, pageSize int) ([]*entity.IoTSensorReading, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address, location,
			   measurement_data, standard_fields, reading_time, created_at, updated_at,
			   data_x, data_y, peak_x, peak_y, ppm, label, raw_data
		FROM iot_sensor_readings
		WHERE mac_address = $1
		ORDER BY reading_time DESC
		LIMIT $2 OFFSET $3`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, macAddress, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list IoT sensor readings by MAC address: %w", err)
	}
	defer rows.Close()

	return r.scanReadings(rows)
}

// ListByTimeRange retrieves IoT sensor readings within a specific time range
func (r *iotSensorReadingRepository) ListByTimeRange(ctx context.Context, startTime, endTime time.Time, page, pageSize int) ([]*entity.IoTSensorReading, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address, location,
			   measurement_data, standard_fields, reading_time, created_at, updated_at,
			   data_x, data_y, peak_x, peak_y, ppm, label, raw_data
		FROM iot_sensor_readings
		WHERE reading_time BETWEEN $1 AND $2
		ORDER BY reading_time DESC
		LIMIT $3 OFFSET $4`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, startTime, endTime, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list IoT sensor readings by time range: %w", err)
	}
	defer rows.Close()

	return r.scanReadings(rows)
}

// Update modifies an existing IoT sensor reading
func (r *iotSensorReadingRepository) Update(ctx context.Context, reading *entity.IoTSensorReading) error {
	query := `
		UPDATE iot_sensor_readings
		SET tenant_id = $1, asset_sensor_id = $2, sensor_type_id = $3, mac_address = $4, 
			location = $5, measurement_data = $6, standard_fields = $7, reading_time = $8,
			updated_at = $9, data_x = $10, data_y = $11, peak_x = $12, peak_y = $13,
			ppm = $14, label = $15, raw_data = $16
		WHERE id = $17`

	now := time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		reading.TenantID,
		reading.AssetSensorID,
		reading.SensorTypeID,
		reading.MacAddress,
		reading.Location,
		reading.MeasurementData,
		reading.StandardFields,
		reading.ReadingTime,
		now,
		reading.DataX,
		reading.DataY,
		reading.PeakX,
		reading.PeakY,
		reading.PPM,
		reading.Label,
		reading.RawData,
		reading.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update IoT sensor reading: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("IoT sensor reading not found")
	}

	reading.UpdatedAt = &now
	return nil
}

// Delete removes an IoT sensor reading by its ID
func (r *iotSensorReadingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM iot_sensor_readings WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete IoT sensor reading: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("IoT sensor reading not found")
	}

	return nil
}

// GetLatestByMacAddress retrieves the latest IoT sensor reading for a specific MAC address
func (r *iotSensorReadingRepository) GetLatestByMacAddress(ctx context.Context, macAddress string) (*entity.IoTSensorReading, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, sensor_type_id, mac_address, location,
			   measurement_data, standard_fields, reading_time, created_at, updated_at,
			   data_x, data_y, peak_x, peak_y, ppm, label, raw_data
		FROM iot_sensor_readings
		WHERE mac_address = $1
		ORDER BY reading_time DESC
		LIMIT 1`

	var reading entity.IoTSensorReading
	var tenantID, assetSensorID sql.NullString
	var measurementData, standardFields, dataX, dataY, peakX, peakY, rawData []byte
	var ppm sql.NullFloat64
	var label sql.NullString
	var updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, macAddress).Scan(
		&reading.ID,
		&tenantID,
		&assetSensorID,
		&reading.SensorTypeID,
		&reading.MacAddress,
		&reading.Location,
		&measurementData,
		&standardFields,
		&reading.ReadingTime,
		&reading.CreatedAt,
		&updatedAt,
		&dataX,
		&dataY,
		&peakX,
		&peakY,
		&ppm,
		&label,
		&rawData,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest IoT sensor reading: %w", err)
	}

	// Handle nullable fields (same as GetByID)
	if tenantID.Valid {
		parsedID, err := uuid.Parse(tenantID.String)
		if err != nil {
			return nil, fmt.Errorf("invalid tenant ID format: %w", err)
		}
		reading.TenantID = parsedID
	}

	if assetSensorID.Valid {
		parsedID, err := uuid.Parse(assetSensorID.String)
		if err != nil {
			return nil, fmt.Errorf("invalid asset sensor ID format: %w", err)
		}
		reading.AssetSensorID = parsedID
	}

	if updatedAt.Valid {
		reading.UpdatedAt = &updatedAt.Time
	}

	if ppm.Valid {
		reading.PPM = ppm.Float64
	}

	if label.Valid {
		reading.Label = label.String
	}

	// Handle JSON fields
	if measurementData != nil {
		reading.MeasurementData = json.RawMessage(measurementData)
	}
	if standardFields != nil {
		reading.StandardFields = json.RawMessage(standardFields)
	}
	if dataX != nil {
		reading.DataX = json.RawMessage(dataX)
	}
	if dataY != nil {
		reading.DataY = json.RawMessage(dataY)
	}
	if peakX != nil {
		reading.PeakX = json.RawMessage(peakX)
	}
	if peakY != nil {
		reading.PeakY = json.RawMessage(peakY)
	}
	if rawData != nil {
		reading.RawData = json.RawMessage(rawData)
	}

	return &reading, nil
}

// GetCountByTenant returns the total count of IoT sensor readings for a tenant
func (r *iotSensorReadingRepository) GetCountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM iot_sensor_readings WHERE tenant_id = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count IoT sensor readings: %w", err)
	}

	return count, nil
}

// GetAverageValueByTimeRange calculates average value from measurement_data for a sensor type within time range
func (r *iotSensorReadingRepository) GetAverageValueByTimeRange(ctx context.Context, sensorTypeID uuid.UUID, startTime, endTime time.Time) (float64, error) {
	query := `
		SELECT AVG(CAST(measurement_data->>'value' AS FLOAT))
		FROM iot_sensor_readings
		WHERE sensor_type_id = $1 
		AND reading_time BETWEEN $2 AND $3
		AND measurement_data->>'value' IS NOT NULL`

	var avg sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, sensorTypeID, startTime, endTime).Scan(&avg)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate average value: %w", err)
	}

	if !avg.Valid {
		return 0, nil
	}

	return avg.Float64, nil
}

// scanReadings is a helper function to scan multiple rows into IoT sensor reading entities
func (r *iotSensorReadingRepository) scanReadings(rows *sql.Rows) ([]*entity.IoTSensorReading, error) {
	var readings []*entity.IoTSensorReading

	for rows.Next() {
		var reading entity.IoTSensorReading
		var tenantID, assetSensorID sql.NullString
		var measurementData, standardFields, dataX, dataY, peakX, peakY, rawData []byte
		var ppm sql.NullFloat64
		var label sql.NullString
		var updatedAt sql.NullTime

		err := rows.Scan(
			&reading.ID,
			&tenantID,
			&assetSensorID,
			&reading.SensorTypeID,
			&reading.MacAddress,
			&reading.Location,
			&measurementData,
			&standardFields,
			&reading.ReadingTime,
			&reading.CreatedAt,
			&updatedAt,
			&dataX,
			&dataY,
			&peakX,
			&peakY,
			&ppm,
			&label,
			&rawData,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan IoT sensor reading: %w", err)
		}

		// Handle nullable fields
		if tenantID.Valid {
			parsedID, err := uuid.Parse(tenantID.String)
			if err != nil {
				return nil, fmt.Errorf("invalid tenant ID format: %w", err)
			}
			reading.TenantID = parsedID
		}

		if assetSensorID.Valid {
			parsedID, err := uuid.Parse(assetSensorID.String)
			if err != nil {
				return nil, fmt.Errorf("invalid asset sensor ID format: %w", err)
			}
			reading.AssetSensorID = parsedID
		}

		if updatedAt.Valid {
			reading.UpdatedAt = &updatedAt.Time
		}

		if ppm.Valid {
			reading.PPM = ppm.Float64
		}

		if label.Valid {
			reading.Label = label.String
		}

		// Handle JSON fields
		if measurementData != nil {
			reading.MeasurementData = json.RawMessage(measurementData)
		}
		if standardFields != nil {
			reading.StandardFields = json.RawMessage(standardFields)
		}
		if dataX != nil {
			reading.DataX = json.RawMessage(dataX)
		}
		if dataY != nil {
			reading.DataY = json.RawMessage(dataY)
		}
		if peakX != nil {
			reading.PeakX = json.RawMessage(peakX)
		}
		if peakY != nil {
			reading.PeakY = json.RawMessage(peakY)
		}
		if rawData != nil {
			reading.RawData = json.RawMessage(rawData)
		}

		readings = append(readings, &reading)
	}

	return readings, nil
}
