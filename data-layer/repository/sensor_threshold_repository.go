package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// SensorThresholdRepository defines the interface for sensor threshold data operations
type SensorThresholdRepository interface {
	Create(ctx context.Context, threshold *entity.SensorThreshold) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SensorThreshold, error)
	List(ctx context.Context, tenantID *uuid.UUID, page, pageSize int) ([]*entity.SensorThreshold, error)
	ListByAssetSensor(ctx context.Context, assetSensorID uuid.UUID, page, pageSize int) ([]*entity.SensorThreshold, error)
	ListBySensorType(ctx context.Context, sensorTypeID uuid.UUID, page, pageSize int) ([]*entity.SensorThreshold, error)
	ListByMeasurementField(ctx context.Context, measurementField string, page, pageSize int) ([]*entity.SensorThreshold, error)
	ListActive(ctx context.Context, tenantID *uuid.UUID, page, pageSize int) ([]*entity.SensorThreshold, error)
	Update(ctx context.Context, threshold *entity.SensorThreshold) error
	Delete(ctx context.Context, id uuid.UUID) error
	ActivateThreshold(ctx context.Context, id uuid.UUID) error
	DeactivateThreshold(ctx context.Context, id uuid.UUID) error
	CheckThresholdBreaches(ctx context.Context, reading *entity.IoTSensorReading) ([]*entity.SensorThreshold, error)
}

// sensorThresholdRepository handles database operations for sensor thresholds
type sensorThresholdRepository struct {
	db *sql.DB
}

// NewSensorThresholdRepository creates a new SensorThresholdRepository
func NewSensorThresholdRepository(db *sql.DB) SensorThresholdRepository {
	return &sensorThresholdRepository{db: db}
}

// Create inserts a new sensor threshold into the database
func (r *sensorThresholdRepository) Create(ctx context.Context, threshold *entity.SensorThreshold) error {
	query := `
		INSERT INTO sensor_thresholds (
			id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field, 
			name, description, min_value, max_value, severity, 
			alert_message, notification_rules, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`

	now := time.Now()
	threshold.CreatedAt = now

	_, err := r.db.ExecContext(
		ctx,
		query,
		threshold.ID,
		threshold.TenantID,
		threshold.AssetSensorID,
		threshold.SensorTypeID,
		threshold.MeasurementField,
		threshold.Name,
		threshold.Description,
		threshold.MinValue,
		threshold.MaxValue,
		threshold.Severity,
		threshold.AlertMessage,
		threshold.NotificationRules,
		threshold.IsActive,
		threshold.CreatedAt,
		threshold.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violations
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("a threshold with the same parameters already exists: %v", err)
		}
		return fmt.Errorf("failed to create sensor threshold: %v", err)
	}

	return nil
}

// GetByID retrieves a sensor threshold by its ID
func (r *sensorThresholdRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SensorThreshold, error) {
	query := `
		SELECT 
			id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field, 
			name, description, min_value, max_value, severity, 
			alert_message, notification_rules, is_active, created_at, updated_at
		FROM sensor_thresholds
		WHERE id = $1`

	var threshold entity.SensorThreshold
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&threshold.ID,
		&threshold.TenantID,
		&threshold.AssetSensorID,
		&threshold.SensorTypeID,
		&threshold.MeasurementField,
		&threshold.Name,
		&threshold.Description,
		&threshold.MinValue,
		&threshold.MaxValue,
		&threshold.Severity,
		&threshold.AlertMessage,
		&threshold.NotificationRules,
		&threshold.IsActive,
		&threshold.CreatedAt,
		&threshold.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("sensor threshold not found")
		}
		return nil, fmt.Errorf("failed to get sensor threshold: %v", err)
	}

	return &threshold, nil
}

// List retrieves a paginated list of sensor thresholds filtered by tenant
func (r *sensorThresholdRepository) List(ctx context.Context, tenantID *uuid.UUID, page, pageSize int) ([]*entity.SensorThreshold, error) {
	offset := (page - 1) * pageSize
	var query string
	var rows *sql.Rows
	var err error

	if tenantID != nil {
		query = `
			SELECT 
				id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field, 
				name, description, min_value, max_value, severity, 
				alert_message, notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE tenant_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3`
		rows, err = r.db.QueryContext(ctx, query, tenantID, pageSize, offset)
	} else {
		query = `
			SELECT 
				id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field, 
				name, description, min_value, max_value, severity, 
				alert_message, notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2`
		rows, err = r.db.QueryContext(ctx, query, pageSize, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list sensor thresholds: %v", err)
	}
	defer rows.Close()

	var thresholds []*entity.SensorThreshold
	for rows.Next() {
		var threshold entity.SensorThreshold
		if err := rows.Scan(
			&threshold.ID,
			&threshold.TenantID,
			&threshold.AssetSensorID,
			&threshold.SensorTypeID,
			&threshold.MeasurementField,
			&threshold.Name,
			&threshold.Description,
			&threshold.MinValue,
			&threshold.MaxValue,
			&threshold.Severity,
			&threshold.AlertMessage,
			&threshold.NotificationRules,
			&threshold.IsActive,
			&threshold.CreatedAt,
			&threshold.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sensor threshold: %v", err)
		}
		thresholds = append(thresholds, &threshold)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over sensor thresholds: %v", err)
	}

	return thresholds, nil
}

// ListByAssetSensor retrieves thresholds for a specific asset sensor
func (r *sensorThresholdRepository) ListByAssetSensor(ctx context.Context, assetSensorID uuid.UUID, page, pageSize int) ([]*entity.SensorThreshold, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT 
			id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field, 
			name, description, min_value, max_value, severity, 
			alert_message, notification_rules, is_active, created_at, updated_at
		FROM sensor_thresholds
		WHERE asset_sensor_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, assetSensorID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list thresholds by asset sensor: %v", err)
	}
	defer rows.Close()

	var thresholds []*entity.SensorThreshold
	for rows.Next() {
		var threshold entity.SensorThreshold
		if err := rows.Scan(
			&threshold.ID,
			&threshold.TenantID,
			&threshold.AssetSensorID,
			&threshold.SensorTypeID,
			&threshold.MeasurementField,
			&threshold.Name,
			&threshold.Description,
			&threshold.MinValue,
			&threshold.MaxValue,
			&threshold.Severity,
			&threshold.AlertMessage,
			&threshold.NotificationRules,
			&threshold.IsActive,
			&threshold.CreatedAt,
			&threshold.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sensor threshold: %v", err)
		}
		thresholds = append(thresholds, &threshold)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over sensor thresholds: %v", err)
	}

	return thresholds, nil
}

// ListBySensorType retrieves thresholds for a specific sensor type
func (r *sensorThresholdRepository) ListBySensorType(ctx context.Context, sensorTypeID uuid.UUID, page, pageSize int) ([]*entity.SensorThreshold, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT 
			id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field, 
			name, description, min_value, max_value, severity, 
			alert_message, notification_rules, is_active, created_at, updated_at
		FROM sensor_thresholds
		WHERE sensor_type_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, sensorTypeID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list thresholds by sensor type: %v", err)
	}
	defer rows.Close()

	var thresholds []*entity.SensorThreshold
	for rows.Next() {
		var threshold entity.SensorThreshold
		if err := rows.Scan(
			&threshold.ID,
			&threshold.TenantID,
			&threshold.AssetSensorID,
			&threshold.SensorTypeID,
			&threshold.MeasurementField,
			&threshold.Name,
			&threshold.Description,
			&threshold.MinValue,
			&threshold.MaxValue,
			&threshold.Severity,
			&threshold.AlertMessage,
			&threshold.NotificationRules,
			&threshold.IsActive,
			&threshold.CreatedAt,
			&threshold.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sensor threshold: %v", err)
		}
		thresholds = append(thresholds, &threshold)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over sensor thresholds: %v", err)
	}

	return thresholds, nil
}

// ListByMeasurementField retrieves thresholds for a specific measurement field
func (r *sensorThresholdRepository) ListByMeasurementField(ctx context.Context, measurementField string, page, pageSize int) ([]*entity.SensorThreshold, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT 
			id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field, 
			name, description, min_value, max_value, severity, 
			alert_message, notification_rules, is_active, created_at, updated_at
		FROM sensor_thresholds
		WHERE measurement_field = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, measurementField, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list thresholds by measurement field: %v", err)
	}
	defer rows.Close()

	var thresholds []*entity.SensorThreshold
	for rows.Next() {
		var threshold entity.SensorThreshold
		if err := rows.Scan(
			&threshold.ID,
			&threshold.TenantID,
			&threshold.AssetSensorID,
			&threshold.SensorTypeID,
			&threshold.MeasurementField,
			&threshold.Name,
			&threshold.Description,
			&threshold.MinValue,
			&threshold.MaxValue,
			&threshold.Severity,
			&threshold.AlertMessage,
			&threshold.NotificationRules,
			&threshold.IsActive,
			&threshold.CreatedAt,
			&threshold.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sensor threshold: %v", err)
		}
		thresholds = append(thresholds, &threshold)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over sensor thresholds: %v", err)
	}

	return thresholds, nil
}

// ListActive retrieves only active thresholds
func (r *sensorThresholdRepository) ListActive(ctx context.Context, tenantID *uuid.UUID, page, pageSize int) ([]*entity.SensorThreshold, error) {
	offset := (page - 1) * pageSize
	var query string
	var rows *sql.Rows
	var err error

	if tenantID != nil {
		query = `
			SELECT 
				id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field, 
				name, description, min_value, max_value, severity, 
				alert_message, notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE tenant_id = $1 AND is_active = true
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3`
		rows, err = r.db.QueryContext(ctx, query, tenantID, pageSize, offset)
	} else {
		query = `
			SELECT 
				id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field, 
				name, description, min_value, max_value, severity, 
				alert_message, notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE is_active = true
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2`
		rows, err = r.db.QueryContext(ctx, query, pageSize, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list active sensor thresholds: %v", err)
	}
	defer rows.Close()

	var thresholds []*entity.SensorThreshold
	for rows.Next() {
		var threshold entity.SensorThreshold
		if err := rows.Scan(
			&threshold.ID,
			&threshold.TenantID,
			&threshold.AssetSensorID,
			&threshold.SensorTypeID,
			&threshold.MeasurementField,
			&threshold.Name,
			&threshold.Description,
			&threshold.MinValue,
			&threshold.MaxValue,
			&threshold.Severity,
			&threshold.AlertMessage,
			&threshold.NotificationRules,
			&threshold.IsActive,
			&threshold.CreatedAt,
			&threshold.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sensor threshold: %v", err)
		}
		thresholds = append(thresholds, &threshold)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over sensor thresholds: %v", err)
	}

	return thresholds, nil
}

// Update updates an existing sensor threshold
func (r *sensorThresholdRepository) Update(ctx context.Context, threshold *entity.SensorThreshold) error {
	query := `
		UPDATE sensor_thresholds SET
			asset_sensor_id = $1,
			sensor_type_id = $2,
			measurement_field = $3,
			name = $4,
			description = $5,
			min_value = $6,
			max_value = $7,
			severity = $8,
			alert_message = $9,
			notification_rules = $10,
			is_active = $11,
			updated_at = $12
		WHERE id = $13 AND tenant_id = $14`

	now := time.Now()
	threshold.UpdatedAt = &now

	result, err := r.db.ExecContext(
		ctx,
		query,
		threshold.AssetSensorID,
		threshold.SensorTypeID,
		threshold.MeasurementField,
		threshold.Name,
		threshold.Description,
		threshold.MinValue,
		threshold.MaxValue,
		threshold.Severity,
		threshold.AlertMessage,
		threshold.NotificationRules,
		threshold.IsActive,
		threshold.UpdatedAt,
		threshold.ID,
		threshold.TenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update sensor threshold: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("sensor threshold not found or you don't have permission to update it")
	}

	return nil
}

// Delete removes a sensor threshold
func (r *sensorThresholdRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sensor_thresholds WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sensor threshold: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("sensor threshold not found")
	}

	return nil
}

// ActivateThreshold sets a threshold to active
func (r *sensorThresholdRepository) ActivateThreshold(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE sensor_thresholds
		SET is_active = true, updated_at = $1
		WHERE id = $2`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("failed to activate sensor threshold: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("sensor threshold not found")
	}

	return nil
}

// DeactivateThreshold sets a threshold to inactive
func (r *sensorThresholdRepository) DeactivateThreshold(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE sensor_thresholds
		SET is_active = false, updated_at = $1
		WHERE id = $2`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate sensor threshold: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("sensor threshold not found")
	}

	return nil
}

// CheckThresholdBreaches checks if a reading breaches any thresholds
func (r *sensorThresholdRepository) CheckThresholdBreaches(ctx context.Context, reading *entity.IoTSensorReading) ([]*entity.SensorThreshold, error) {
	// This is a placeholder implementation. In a real implementation, you would:
	// 1. Get the measurement data from the reading
	// 2. Extract relevant field values
	// 3. Query for active thresholds that match the reading criteria (sensor type, asset sensor)
	// 4. Check each threshold against the corresponding field value
	// 5. Return thresholds that are breached

	query := `
		SELECT 
			id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field, 
			name, description, min_value, max_value, severity, 
			alert_message, notification_rules, is_active, created_at, updated_at
		FROM sensor_thresholds
		WHERE sensor_type_id = $1 
		AND (asset_sensor_id = $2 OR asset_sensor_id IS NULL)
		AND is_active = true`

	rows, err := r.db.QueryContext(ctx, query, reading.SensorTypeID, reading.AssetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for threshold breaches: %v", err)
	}
	defer rows.Close()

	var breachedThresholds []*entity.SensorThreshold
	for rows.Next() {
		var threshold entity.SensorThreshold
		if err := rows.Scan(
			&threshold.ID,
			&threshold.TenantID,
			&threshold.AssetSensorID,
			&threshold.SensorTypeID,
			&threshold.MeasurementField,
			&threshold.Name,
			&threshold.Description,
			&threshold.MinValue,
			&threshold.MaxValue,
			&threshold.Severity,
			&threshold.AlertMessage,
			&threshold.NotificationRules,
			&threshold.IsActive,
			&threshold.CreatedAt,
			&threshold.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sensor threshold: %v", err)
		}

		// For simplicity, we're assuming the field exists in measurement data
		// In a real implementation, you would use more sophisticated JSON parsing
		// This is a simplified implementation
		var measurementData map[string]interface{}
		if err := json.Unmarshal(reading.MeasurementData, &measurementData); err != nil {
			continue // Skip if measurement data is not valid JSON
		}

		fieldValue, ok := measurementData[threshold.MeasurementField].(float64)
		if !ok {
			continue // Skip if field doesn't exist or is not a number
		}

		// Check if the value breaches the threshold
		if threshold.IsBreached(fieldValue) {
			breachedThresholds = append(breachedThresholds, &threshold)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over thresholds: %v", err)
	}

	return breachedThresholds, nil
}
