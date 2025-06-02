package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/helpers/common"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AssetAlertRepository defines the interface for asset alert data operations
type AssetAlertRepository interface {
	Create(ctx context.Context, alert *entity.AssetAlert) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AssetAlert, error)
	GetByAssetID(ctx context.Context, assetID uuid.UUID) ([]*entity.AssetAlert, error)
	GetByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.AssetAlert, error)
	GetByThresholdID(ctx context.Context, thresholdID uuid.UUID) ([]*entity.AssetAlert, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.AssetAlert, error)
	Update(ctx context.Context, alert *entity.AssetAlert) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByAssetID(ctx context.Context, assetID uuid.UUID) error
	GetUnresolvedAlerts(ctx context.Context) ([]*entity.AssetAlert, error)
	GetAlertsBySeverity(ctx context.Context, severity entity.ThresholdSeverity) ([]*entity.AssetAlert, error)
	ResolveAlert(ctx context.Context, id uuid.UUID) error
	GetAlertsInTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*entity.AssetAlert, error)
}

// assetAlertRepository handles database operations for asset alerts
type assetAlertRepository struct {
	*BaseRepository
}

// NewAssetAlertRepository creates a new AssetAlertRepository
func NewAssetAlertRepository(db *sql.DB) AssetAlertRepository {
	return &assetAlertRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create inserts a new asset alert into the database
func (r *assetAlertRepository) Create(ctx context.Context, alert *entity.AssetAlert) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	query := `
		INSERT INTO asset_alerts (
			tenant_id, asset_id, asset_sensor_id, threshold_id, 
			alert_time, resolved_time, severity
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING id`

	// Use tenant ID from context or entity
	var entityTenantID uuid.UUID
	if hasTenantID {
		entityTenantID = tenantID
	} else {
		entityTenantID = alert.TenantID
	}

	err := r.DB.QueryRowContext(
		ctx,
		query,
		entityTenantID,
		alert.AssetID,
		alert.AssetSensorID,
		alert.ThresholdID,
		alert.AlertTime,
		alert.ResolvedTime,
		string(alert.Severity),
	).Scan(&alert.ID)

	if err != nil {
		return fmt.Errorf("failed to create asset alert: %w", err)
	}

	alert.TenantID = entityTenantID
	return nil
}

// GetByID retrieves an asset alert by its ID
func (r *assetAlertRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AssetAlert, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access any alert
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE id = $1`
		args = []interface{}{id}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE id = $1 AND tenant_id = $2`
		args = []interface{}{id, tenantID}
	}

	var alert entity.AssetAlert
	var resolvedTime sql.NullTime

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&alert.ID,
		&alert.TenantID,
		&alert.AssetID,
		&alert.AssetSensorID,
		&alert.ThresholdID,
		&alert.AlertTime,
		&resolvedTime,
		&alert.Severity,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get asset alert: %w", err)
	}

	// Handle nullable resolved_time
	if resolvedTime.Valid {
		alert.ResolvedTime = &resolvedTime.Time
	}

	return &alert, nil
}

// GetByAssetID retrieves all alerts for a specific asset
func (r *assetAlertRepository) GetByAssetID(ctx context.Context, assetID uuid.UUID) ([]*entity.AssetAlert, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access all alerts for the asset
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE asset_id = $1
			ORDER BY alert_time DESC`
		args = []interface{}{assetID}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE asset_id = $1 AND tenant_id = $2
			ORDER BY alert_time DESC`
		args = []interface{}{assetID, tenantID}
	}

	return r.scanAlerts(ctx, query, args...)
}

// GetByAssetSensorID retrieves all alerts for a specific asset sensor
func (r *assetAlertRepository) GetByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.AssetAlert, error) {
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
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE asset_sensor_id = $1
			ORDER BY alert_time DESC`
		args = []interface{}{assetSensorID}
	} else {
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE asset_sensor_id = $1 AND tenant_id = $2
			ORDER BY alert_time DESC`
		args = []interface{}{assetSensorID, tenantID}
	}

	return r.scanAlerts(ctx, query, args...)
}

// GetByThresholdID retrieves all alerts for a specific threshold
func (r *assetAlertRepository) GetByThresholdID(ctx context.Context, thresholdID uuid.UUID) ([]*entity.AssetAlert, error) {
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
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE threshold_id = $1
			ORDER BY alert_time DESC`
		args = []interface{}{thresholdID}
	} else {
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE threshold_id = $1 AND tenant_id = $2
			ORDER BY alert_time DESC`
		args = []interface{}{thresholdID, tenantID}
	}

	return r.scanAlerts(ctx, query, args...)
}

// List retrieves asset alerts with pagination
func (r *assetAlertRepository) List(ctx context.Context, page, pageSize int) ([]*entity.AssetAlert, error) {
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
		// SuperAdmin without tenant ID can access all alerts across all tenants
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			ORDER BY alert_time DESC
			LIMIT $1 OFFSET $2`
		args = []interface{}{pageSize, offset}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE tenant_id = $1
			ORDER BY alert_time DESC
			LIMIT $2 OFFSET $3`
		args = []interface{}{tenantID, pageSize, offset}
	}

	return r.scanAlerts(ctx, query, args...)
}

// Update modifies an existing asset alert
func (r *assetAlertRepository) Update(ctx context.Context, alert *entity.AssetAlert) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin {
		// SuperAdmin can update any alert
		query = `
			UPDATE asset_alerts
			SET resolved_time = $1, severity = $2
			WHERE id = $3`
		args = []interface{}{alert.ResolvedTime, string(alert.Severity), alert.ID}
	} else {
		// Regular users can only update alerts from their tenant
		query = `
			UPDATE asset_alerts
			SET resolved_time = $1, severity = $2
			WHERE id = $3 AND tenant_id = $4`
		args = []interface{}{alert.ResolvedTime, string(alert.Severity), alert.ID, tenantID}
	}

	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update asset alert: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset alert not found or access denied")
	}

	return nil
}

// Delete removes an asset alert by ID
func (r *assetAlertRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin {
		// SuperAdmin can delete any alert
		query = `DELETE FROM asset_alerts WHERE id = $1`
		args = []interface{}{id}
	} else {
		// Regular users can only delete alerts from their tenant
		query = `DELETE FROM asset_alerts WHERE id = $1 AND tenant_id = $2`
		args = []interface{}{id, tenantID}
	}

	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete asset alert: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset alert not found or access denied")
	}

	return nil
}

// DeleteByAssetID removes all alerts for a specific asset
func (r *assetAlertRepository) DeleteByAssetID(ctx context.Context, assetID uuid.UUID) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin {
		// SuperAdmin can delete alerts from any asset
		query = `DELETE FROM asset_alerts WHERE asset_id = $1`
		args = []interface{}{assetID}
	} else {
		// Regular users can only delete alerts from assets in their tenant
		query = `DELETE FROM asset_alerts WHERE asset_id = $1 AND tenant_id = $2`
		args = []interface{}{assetID, tenantID}
	}

	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete asset alerts: %w", err)
	}

	return nil
}

// GetUnresolvedAlerts retrieves all unresolved alerts
func (r *assetAlertRepository) GetUnresolvedAlerts(ctx context.Context) ([]*entity.AssetAlert, error) {
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
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE resolved_time IS NULL
			ORDER BY alert_time DESC`
		args = []interface{}{}
	} else {
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE resolved_time IS NULL AND tenant_id = $1
			ORDER BY alert_time DESC`
		args = []interface{}{tenantID}
	}

	return r.scanAlerts(ctx, query, args...)
}

// GetAlertsBySeverity retrieves alerts by severity level
func (r *assetAlertRepository) GetAlertsBySeverity(ctx context.Context, severity entity.ThresholdSeverity) ([]*entity.AssetAlert, error) {
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
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE severity = $1
			ORDER BY alert_time DESC`
		args = []interface{}{string(severity)}
	} else {
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE severity = $1 AND tenant_id = $2
			ORDER BY alert_time DESC`
		args = []interface{}{string(severity), tenantID}
	}

	return r.scanAlerts(ctx, query, args...)
}

// ResolveAlert marks an alert as resolved
func (r *assetAlertRepository) ResolveAlert(ctx context.Context, id uuid.UUID) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin {
		// SuperAdmin can resolve any alert
		query = `UPDATE asset_alerts SET resolved_time = $1 WHERE id = $2`
		args = []interface{}{time.Now(), id}
	} else {
		// Regular users can only resolve alerts from their tenant
		query = `UPDATE asset_alerts SET resolved_time = $1 WHERE id = $2 AND tenant_id = $3`
		args = []interface{}{time.Now(), id, tenantID}
	}

	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to resolve asset alert: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset alert not found or access denied")
	}

	return nil
}

// GetAlertsInTimeRange retrieves alerts within a specific time range
func (r *assetAlertRepository) GetAlertsInTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*entity.AssetAlert, error) {
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
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE alert_time >= $1 AND alert_time <= $2
			ORDER BY alert_time DESC`
		args = []interface{}{startTime, endTime}
	} else {
		query = `
			SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
				   alert_time, resolved_time, severity
			FROM asset_alerts
			WHERE alert_time >= $1 AND alert_time <= $2 AND tenant_id = $3
			ORDER BY alert_time DESC`
		args = []interface{}{startTime, endTime, tenantID}
	}

	return r.scanAlerts(ctx, query, args...)
}

// scanAlerts is a helper method to scan multiple alerts from query results
func (r *assetAlertRepository) scanAlerts(ctx context.Context, query string, args ...interface{}) ([]*entity.AssetAlert, error) {
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var alerts []*entity.AssetAlert

	for rows.Next() {
		var alert entity.AssetAlert
		var resolvedTime sql.NullTime

		err := rows.Scan(
			&alert.ID,
			&alert.TenantID,
			&alert.AssetID,
			&alert.AssetSensorID,
			&alert.ThresholdID,
			&alert.AlertTime,
			&resolvedTime,
			&alert.Severity,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}

		// Handle nullable resolved_time
		if resolvedTime.Valid {
			alert.ResolvedTime = &resolvedTime.Time
		}

		alerts = append(alerts, &alert)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return alerts, nil
}