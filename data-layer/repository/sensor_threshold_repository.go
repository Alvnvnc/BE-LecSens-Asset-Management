package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/helpers/common"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SensorThresholdRepository defines the interface for sensor threshold data operations
type SensorThresholdRepository interface {
	Create(ctx context.Context, threshold *entity.SensorThreshold) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SensorThreshold, error)
	GetByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.SensorThreshold, error)
	GetBySensorTypeID(ctx context.Context, sensorTypeID uuid.UUID) ([]*entity.SensorThreshold, error)
	GetByMeasurementField(ctx context.Context, assetSensorID uuid.UUID, measurementField string) ([]*entity.SensorThreshold, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.SensorThreshold, error)
	Update(ctx context.Context, threshold *entity.SensorThreshold) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) error
	GetActiveThresholds(ctx context.Context) ([]*entity.SensorThreshold, error)
	GetThresholdsBySeverity(ctx context.Context, severity entity.ThresholdSeverity) ([]*entity.SensorThreshold, error)
}

// sensorThresholdRepository handles database operations for sensor thresholds
type sensorThresholdRepository struct {
	*BaseRepository
}

// NewSensorThresholdRepository creates a new SensorThresholdRepository
func NewSensorThresholdRepository(db *sql.DB) SensorThresholdRepository {
	return &sensorThresholdRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create inserts a new sensor threshold into the database
func (r *sensorThresholdRepository) Create(ctx context.Context, threshold *entity.SensorThreshold) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	query := `
		INSERT INTO sensor_thresholds (
			tenant_id, asset_sensor_id, sensor_type_id, measurement_field, name, 
			description, min_value, max_value, severity, alert_message, 
			notification_rules, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING id`

	now := time.Now()

	// Use tenant ID from context or entity
	var entityTenantID uuid.UUID
	if hasTenantID {
		entityTenantID = tenantID
	} else {
		entityTenantID = threshold.TenantID
	}

	// Convert notification rules to JSON
	var notificationRulesJSON interface{}
	if threshold.NotificationRules != nil {
		notificationRulesJSON = threshold.NotificationRules
	} else {
		notificationRulesJSON = json.RawMessage("{}")
	}

	err := r.DB.QueryRowContext(
		ctx,
		query,
		entityTenantID,
		threshold.AssetSensorID,
		threshold.SensorTypeID,
		threshold.MeasurementField,
		threshold.Name,
		threshold.Description,
		threshold.MinValue,
		threshold.MaxValue,
		string(threshold.Severity),
		threshold.AlertMessage,
		notificationRulesJSON,
		threshold.IsActive,
		now,
		now,
	).Scan(&threshold.ID)

	if err != nil {
		return fmt.Errorf("failed to create sensor threshold: %w", err)
	}

	threshold.TenantID = entityTenantID
	threshold.CreatedAt = now
	updatedAt := now
	threshold.UpdatedAt = &updatedAt
	return nil
}

// GetByID retrieves a sensor threshold by its ID
func (r *sensorThresholdRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SensorThreshold, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access any threshold
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE id = $1`
		args = []interface{}{id}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE id = $1 AND tenant_id = $2`
		args = []interface{}{id, tenantID}
	}

	var threshold entity.SensorThreshold
	var updatedAt sql.NullTime
	var description, alertMessage sql.NullString
	var notificationRules sql.NullString

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&threshold.ID,
		&threshold.TenantID,
		&threshold.AssetSensorID,
		&threshold.SensorTypeID,
		&threshold.MeasurementField,
		&threshold.Name,
		&description,
		&threshold.MinValue,
		&threshold.MaxValue,
		&threshold.Severity,
		&alertMessage,
		&notificationRules,
		&threshold.IsActive,
		&threshold.CreatedAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor threshold: %w", err)
	}

	// Handle nullable fields
	if description.Valid {
		threshold.Description = description.String
	}
	if alertMessage.Valid {
		threshold.AlertMessage = alertMessage.String
	}
	if notificationRules.Valid {
		threshold.NotificationRules = json.RawMessage(notificationRules.String)
	}
	if updatedAt.Valid {
		threshold.UpdatedAt = &updatedAt.Time
	}

	return &threshold, nil
}

// GetByAssetSensorID retrieves all thresholds for a specific asset sensor
func (r *sensorThresholdRepository) GetByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.SensorThreshold, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access all thresholds for the asset sensor
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE asset_sensor_id = $1
			ORDER BY created_at DESC`
		args = []interface{}{assetSensorID}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE asset_sensor_id = $1 AND tenant_id = $2
			ORDER BY created_at DESC`
		args = []interface{}{assetSensorID, tenantID}
	}

	return r.scanThresholds(ctx, query, args...)
}

// GetBySensorTypeID retrieves all thresholds for a specific sensor type
func (r *sensorThresholdRepository) GetBySensorTypeID(ctx context.Context, sensorTypeID uuid.UUID) ([]*entity.SensorThreshold, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE sensor_type_id = $1
			ORDER BY created_at DESC`
		args = []interface{}{sensorTypeID}
	} else {
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE sensor_type_id = $1 AND tenant_id = $2
			ORDER BY created_at DESC`
		args = []interface{}{sensorTypeID, tenantID}
	}

	return r.scanThresholds(ctx, query, args...)
}

// GetByMeasurementField retrieves thresholds for specific measurement field
func (r *sensorThresholdRepository) GetByMeasurementField(ctx context.Context, assetSensorID uuid.UUID, measurementField string) ([]*entity.SensorThreshold, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE asset_sensor_id = $1 AND measurement_field = $2 AND is_active = true
			ORDER BY created_at DESC`
		args = []interface{}{assetSensorID, measurementField}
	} else {
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE asset_sensor_id = $1 AND measurement_field = $2 AND tenant_id = $3 AND is_active = true
			ORDER BY created_at DESC`
		args = []interface{}{assetSensorID, measurementField, tenantID}
	}

	return r.scanThresholds(ctx, query, args...)
}

// List retrieves sensor thresholds with pagination
func (r *sensorThresholdRepository) List(ctx context.Context, page, pageSize int) ([]*entity.SensorThreshold, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	offset := (page - 1) * pageSize

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access all thresholds across all tenants
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2`
		args = []interface{}{pageSize, offset}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE tenant_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3`
		args = []interface{}{tenantID, pageSize, offset}
	}

	return r.scanThresholds(ctx, query, args...)
}

// Update modifies an existing sensor threshold
func (r *sensorThresholdRepository) Update(ctx context.Context, threshold *entity.SensorThreshold) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	// Convert notification rules to JSON
	var notificationRulesJSON interface{}
	if threshold.NotificationRules != nil {
		notificationRulesJSON = threshold.NotificationRules
	} else {
		notificationRulesJSON = json.RawMessage("{}")
	}

	var query string
	var args []interface{}

	if isSuperAdmin {
		// SuperAdmin can update any threshold
		query = `
			UPDATE sensor_thresholds
			SET measurement_field = $1, name = $2, description = $3, min_value = $4,
				max_value = $5, severity = $6, alert_message = $7, notification_rules = $8,
				is_active = $9, updated_at = $10
			WHERE id = $11`
		args = []interface{}{
			threshold.MeasurementField, threshold.Name, threshold.Description,
			threshold.MinValue, threshold.MaxValue, string(threshold.Severity),
			threshold.AlertMessage, notificationRulesJSON, threshold.IsActive,
			time.Now(), threshold.ID,
		}
	} else {
		// Regular users can only update thresholds from their tenant
		query = `
			UPDATE sensor_thresholds
			SET measurement_field = $1, name = $2, description = $3, min_value = $4,
				max_value = $5, severity = $6, alert_message = $7, notification_rules = $8,
				is_active = $9, updated_at = $10
			WHERE id = $11 AND tenant_id = $12`
		args = []interface{}{
			threshold.MeasurementField, threshold.Name, threshold.Description,
			threshold.MinValue, threshold.MaxValue, string(threshold.Severity),
			threshold.AlertMessage, notificationRulesJSON, threshold.IsActive,
			time.Now(), threshold.ID, tenantID,
		}
	}

	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update sensor threshold: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sensor threshold not found or access denied")
	}

	return nil
}

// Delete removes a sensor threshold by ID
func (r *sensorThresholdRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin {
		// SuperAdmin can delete any threshold
		query = `DELETE FROM sensor_thresholds WHERE id = $1`
		args = []interface{}{id}
	} else {
		// Regular users can only delete thresholds from their tenant
		query = `DELETE FROM sensor_thresholds WHERE id = $1 AND tenant_id = $2`
		args = []interface{}{id, tenantID}
	}

	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete sensor threshold: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sensor threshold not found or access denied")
	}

	return nil
}

// DeleteByAssetSensorID removes all thresholds for a specific asset sensor
func (r *sensorThresholdRepository) DeleteByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin {
		// SuperAdmin can delete thresholds from any asset sensor
		query = `DELETE FROM sensor_thresholds WHERE asset_sensor_id = $1`
		args = []interface{}{assetSensorID}
	} else {
		// Regular users can only delete thresholds from asset sensors in their tenant
		query = `DELETE FROM sensor_thresholds WHERE asset_sensor_id = $1 AND tenant_id = $2`
		args = []interface{}{assetSensorID, tenantID}
	}

	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete sensor thresholds: %w", err)
	}

	return nil
}

// GetActiveThresholds retrieves all active thresholds
func (r *sensorThresholdRepository) GetActiveThresholds(ctx context.Context) ([]*entity.SensorThreshold, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE is_active = true
			ORDER BY created_at DESC`
		args = []interface{}{}
	} else {
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE is_active = true AND tenant_id = $1
			ORDER BY created_at DESC`
		args = []interface{}{tenantID}
	}

	return r.scanThresholds(ctx, query, args...)
}

// GetThresholdsBySeverity retrieves thresholds by severity level
func (r *sensorThresholdRepository) GetThresholdsBySeverity(ctx context.Context, severity entity.ThresholdSeverity) ([]*entity.SensorThreshold, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE severity = $1
			ORDER BY created_at DESC`
		args = []interface{}{string(severity)}
	} else {
		query = `
			SELECT id, tenant_id, asset_sensor_id, sensor_type_id, measurement_field,
				   name, description, min_value, max_value, severity, alert_message,
				   notification_rules, is_active, created_at, updated_at
			FROM sensor_thresholds
			WHERE severity = $1 AND tenant_id = $2
			ORDER BY created_at DESC`
		args = []interface{}{string(severity), tenantID}
	}

	return r.scanThresholds(ctx, query, args...)
}

// scanThresholds is a helper method to scan multiple thresholds from query results
func (r *sensorThresholdRepository) scanThresholds(ctx context.Context, query string, args ...interface{}) ([]*entity.SensorThreshold, error) {
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sensor thresholds: %w", err)
	}
	defer rows.Close()

	var thresholds []*entity.SensorThreshold
	for rows.Next() {
		var threshold entity.SensorThreshold
		var updatedAt sql.NullTime
		var description, alertMessage sql.NullString
		var notificationRules sql.NullString

		err := rows.Scan(
			&threshold.ID,
			&threshold.TenantID,
			&threshold.AssetSensorID,
			&threshold.SensorTypeID,
			&threshold.MeasurementField,
			&threshold.Name,
			&description,
			&threshold.MinValue,
			&threshold.MaxValue,
			&threshold.Severity,
			&alertMessage,
			&notificationRules,
			&threshold.IsActive,
			&threshold.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sensor threshold: %w", err)
		}

		// Handle nullable fields
		if description.Valid {
			threshold.Description = description.String
		}
		if alertMessage.Valid {
			threshold.AlertMessage = alertMessage.String
		}
		if notificationRules.Valid {
			threshold.NotificationRules = json.RawMessage(notificationRules.String)
		}
		if updatedAt.Valid {
			threshold.UpdatedAt = &updatedAt.Time
		}

		thresholds = append(thresholds, &threshold)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sensor thresholds: %w", err)
	}

	return thresholds, nil
}