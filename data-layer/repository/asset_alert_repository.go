package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// AssetAlertRepository defines the interface for asset alert operations
type AssetAlertRepository interface {
	Create(ctx context.Context, alert *entity.AssetAlert) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AssetAlert, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.AssetAlert, error)
	GetByAssetID(ctx context.Context, assetID uuid.UUID) ([]*entity.AssetAlert, error)
	GetByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.AssetAlert, error)
	GetByMeasurementTypeID(ctx context.Context, measurementTypeID uuid.UUID) ([]*entity.AssetAlert, error)
	GetActiveAlerts(ctx context.Context, tenantID uuid.UUID) ([]*entity.AssetAlert, error)
	GetActiveAlertsByAssetSensor(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.AssetAlert, error)
	Update(ctx context.Context, alert *entity.AssetAlert) error
	ResolveAlert(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.AssetAlert, int, error)
	ListAll(ctx context.Context, limit, offset int) ([]*entity.AssetAlert, int, error)
	ListWithFilters(
		ctx context.Context,
		tenantID uuid.UUID,
		limit, offset int,
		assetID *uuid.UUID,
		assetSensorID *uuid.UUID,
		severity *entity.ThresholdSeverity,
		isResolved *bool,
		fromTime *time.Time,
		toTime *time.Time,
	) ([]*entity.AssetAlert, int, error)
	GetAlertStatistics(
		ctx context.Context,
		tenantID uuid.UUID,
		assetID *uuid.UUID,
		fromTime *time.Time,
		toTime *time.Time,
	) (map[string]interface{}, error)
	GetGlobalAlertStatistics(ctx context.Context) (map[string]interface{}, error)
	CreateAlertFromReading(ctx context.Context, reading *entity.IoTSensorReading, threshold *entity.SensorThreshold, value float64) error
	ResolveAlertsForReading(ctx context.Context, reading *entity.IoTSensorReading, threshold *entity.SensorThreshold, value float64) error
	ResolveMultipleAlerts(ctx context.Context, alertIDs []uuid.UUID) (int, int, error)
	DeleteMultipleAlerts(ctx context.Context, alertIDs []uuid.UUID) (int, int, error)
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
	log.Printf("Creating asset alert: %+v", alert)

	// Generate new UUID if not set
	if alert.ID == uuid.Nil {
		alert.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now()
	if alert.CreatedAt.IsZero() {
		alert.CreatedAt = now
	}
	if alert.AlertTime.IsZero() {
		alert.AlertTime = now
	}

	query := `
		INSERT INTO asset_alerts (
			id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			measurement_field_name, alert_time, severity, trigger_value,
			threshold_min_value, threshold_max_value, alert_message,
			alert_type, is_resolved, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`

	_, err := r.DB.ExecContext(ctx, query,
		alert.ID,
		alert.TenantID,
		alert.AssetID,
		alert.AssetSensorID,
		alert.ThresholdID,
		alert.MeasurementFieldName,
		alert.AlertTime,
		alert.Severity,
		alert.TriggerValue,
		alert.ThresholdMinValue,
		alert.ThresholdMaxValue,
		alert.AlertMessage,
		alert.AlertType,
		alert.IsResolved,
		alert.CreatedAt,
	)

	if err != nil {
		log.Printf("Error creating asset alert: %v", err)
		return fmt.Errorf("failed to create asset alert: %w", err)
	}

	log.Printf("Successfully created asset alert with ID: %s", alert.ID)
	return nil
}

// GetByID retrieves an asset alert by its ID
func (r *assetAlertRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AssetAlert, error) {
	query := `
		SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			   measurement_field_name, alert_time, resolved_time, severity,
			   trigger_value, threshold_min_value, threshold_max_value,
			   alert_message, alert_type, is_resolved, created_at, updated_at
		FROM asset_alerts
		WHERE id = $1`

	var alert entity.AssetAlert
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&alert.ID,
		&alert.TenantID,
		&alert.AssetID,
		&alert.AssetSensorID,
		&alert.ThresholdID,
		&alert.MeasurementFieldName,
		&alert.AlertTime,
		&alert.ResolvedTime,
		&alert.Severity,
		&alert.TriggerValue,
		&alert.ThresholdMinValue,
		&alert.ThresholdMaxValue,
		&alert.AlertMessage,
		&alert.AlertType,
		&alert.IsResolved,
		&alert.CreatedAt,
		&alert.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get asset alert: %w", err)
	}

	return &alert, nil
}

// GetByTenantID retrieves all asset alerts for a tenant
func (r *assetAlertRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.AssetAlert, error) {
	query := `
		SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			   measurement_field_name, alert_time, resolved_time, severity,
			   trigger_value, threshold_min_value, threshold_max_value,
			   alert_message, alert_type, is_resolved, created_at, updated_at
		FROM asset_alerts
		WHERE tenant_id = $1
		ORDER BY alert_time DESC`

	return r.queryAlerts(ctx, query, tenantID)
}

// GetByAssetID retrieves alerts for a specific asset
func (r *assetAlertRepository) GetByAssetID(ctx context.Context, assetID uuid.UUID) ([]*entity.AssetAlert, error) {
	query := `
		SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			   measurement_field_name, alert_time, resolved_time, severity,
			   trigger_value, threshold_min_value, threshold_max_value,
			   alert_message, alert_type, is_resolved, created_at, updated_at
		FROM asset_alerts
		WHERE asset_id = $1
		ORDER BY alert_time DESC`

	return r.queryAlerts(ctx, query, assetID)
}

// GetByAssetSensorID retrieves alerts for a specific asset sensor
func (r *assetAlertRepository) GetByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.AssetAlert, error) {
	query := `
		SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			   measurement_field_name, alert_time, resolved_time, severity,
			   trigger_value, threshold_min_value, threshold_max_value,
			   alert_message, alert_type, is_resolved, created_at, updated_at
		FROM asset_alerts
		WHERE asset_sensor_id = $1
		ORDER BY alert_time DESC`

	return r.queryAlerts(ctx, query, assetSensorID)
}

// GetByMeasurementTypeID retrieves alerts for a specific measurement type
func (r *assetAlertRepository) GetByMeasurementTypeID(ctx context.Context, measurementTypeID uuid.UUID) ([]*entity.AssetAlert, error) {
	query := `
		SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			   measurement_field_name, alert_time, resolved_time, severity,
			   trigger_value, threshold_min_value, threshold_max_value,
			   alert_message, alert_type, is_resolved, created_at, updated_at
		FROM asset_alerts
		WHERE threshold_id IN (
			SELECT id FROM sensor_thresholds WHERE measurement_type_id = $1
		)
		ORDER BY alert_time DESC`

	return r.queryAlerts(ctx, query, measurementTypeID)
}

// GetActiveAlerts retrieves all active (unresolved) alerts for a tenant
func (r *assetAlertRepository) GetActiveAlerts(ctx context.Context, tenantID uuid.UUID) ([]*entity.AssetAlert, error) {
	query := `
		SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			   measurement_field_name, alert_time, resolved_time, severity,
			   trigger_value, threshold_min_value, threshold_max_value,
			   alert_message, alert_type, is_resolved, created_at, updated_at
		FROM asset_alerts
		WHERE tenant_id = $1 AND is_resolved = false
		ORDER BY alert_time DESC`

	return r.queryAlerts(ctx, query, tenantID)
}

// GetActiveAlertsByAssetSensor retrieves active alerts for a specific asset sensor
func (r *assetAlertRepository) GetActiveAlertsByAssetSensor(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.AssetAlert, error) {
	query := `
		SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			   measurement_field_name, alert_time, resolved_time, severity,
			   trigger_value, threshold_min_value, threshold_max_value,
			   alert_message, alert_type, is_resolved, created_at, updated_at
		FROM asset_alerts
		WHERE asset_sensor_id = $1 AND is_resolved = false
		ORDER BY alert_time DESC`

	return r.queryAlerts(ctx, query, assetSensorID)
}

// Update updates an existing asset alert
func (r *assetAlertRepository) Update(ctx context.Context, alert *entity.AssetAlert) error {
	now := time.Now()
	alert.UpdatedAt = &now

	query := `
		UPDATE asset_alerts SET
			alert_message = $2,
			is_resolved = $3,
			resolved_time = $4,
			updated_at = $5
		WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query,
		alert.ID,
		alert.AlertMessage,
		alert.IsResolved,
		alert.ResolvedTime,
		alert.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update asset alert: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset alert not found")
	}

	return nil
}

// ResolveAlert marks an alert as resolved
func (r *assetAlertRepository) ResolveAlert(ctx context.Context, id uuid.UUID) error {
	now := time.Now()

	query := `
		UPDATE asset_alerts SET
			is_resolved = true,
			resolved_time = $2,
			updated_at = $2
		WHERE id = $1 AND is_resolved = false`

	result, err := r.DB.ExecContext(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("failed to resolve asset alert: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset alert not found or already resolved")
	}

	return nil
}

// Delete removes an asset alert by its ID
func (r *assetAlertRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM asset_alerts WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset alert: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset alert not found")
	}

	return nil
}

// List retrieves paginated asset alerts for a tenant
func (r *assetAlertRepository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.AssetAlert, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM asset_alerts WHERE tenant_id = $1`
	var totalCount int
	err := r.DB.QueryRowContext(ctx, countQuery, tenantID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			   measurement_field_name, alert_time, resolved_time, severity,
			   trigger_value, threshold_min_value, threshold_max_value,
			   alert_message, alert_type, is_resolved, created_at, updated_at
		FROM asset_alerts
		WHERE tenant_id = $1
		ORDER BY alert_time DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.DB.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query asset alerts: %w", err)
	}
	defer rows.Close()

	alerts, err := r.scanRows(rows)
	if err != nil {
		return nil, 0, err
	}

	return alerts, totalCount, nil
}

// ListWithFilters retrieves paginated asset alerts with filters
func (r *assetAlertRepository) ListWithFilters(
	ctx context.Context,
	tenantID uuid.UUID,
	limit, offset int,
	assetID *uuid.UUID,
	assetSensorID *uuid.UUID,
	severity *entity.ThresholdSeverity,
	isResolved *bool,
	fromTime *time.Time,
	toTime *time.Time,
) ([]*entity.AssetAlert, int, error) {
	// Build the base query
	baseQuery := `FROM asset_alerts WHERE tenant_id = $1`
	countQuery := `SELECT COUNT(*) ` + baseQuery
	selectQuery := `
		SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			   measurement_field_name, alert_time, resolved_time, severity,
			   trigger_value, threshold_min_value, threshold_max_value,
			   alert_message, alert_type, is_resolved, created_at, updated_at
		` + baseQuery

	// Build filter conditions
	args := []interface{}{tenantID}
	argCount := 1

	if assetID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND asset_id = $%d", argCount)
		args = append(args, *assetID)
	}

	if assetSensorID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND asset_sensor_id = $%d", argCount)
		args = append(args, *assetSensorID)
	}

	if severity != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND severity = $%d", argCount)
		args = append(args, *severity)
	}

	if isResolved != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND is_resolved = $%d", argCount)
		args = append(args, *isResolved)
	}

	if fromTime != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND alert_time >= $%d", argCount)
		args = append(args, *fromTime)
	}

	if toTime != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND alert_time <= $%d", argCount)
		args = append(args, *toTime)
	}

	// Get total count
	var totalCount int
	err := r.DB.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Add pagination
	argCount++
	baseQuery += fmt.Sprintf(" ORDER BY alert_time DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Get paginated results
	rows, err := r.DB.QueryContext(ctx, selectQuery+baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query asset alerts: %w", err)
	}
	defer rows.Close()

	alerts, err := r.scanRows(rows)
	if err != nil {
		return nil, 0, err
	}

	return alerts, totalCount, nil
}

// GetAlertStatistics retrieves alert statistics for a tenant
func (r *assetAlertRepository) GetAlertStatistics(
	ctx context.Context,
	tenantID uuid.UUID,
	assetID *uuid.UUID,
	fromTime *time.Time,
	toTime *time.Time,
) (map[string]interface{}, error) {
	// Build the base query
	baseQuery := `FROM asset_alerts WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	argCount := 1

	if assetID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND asset_id = $%d", argCount)
		args = append(args, *assetID)
	}

	if fromTime != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND alert_time >= $%d", argCount)
		args = append(args, *fromTime)
	}

	if toTime != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND alert_time <= $%d", argCount)
		args = append(args, *toTime)
	}

	query := `
		SELECT 
			COUNT(*) as total_alerts,
			COUNT(CASE WHEN is_resolved = false THEN 1 END) as active_alerts,
			COUNT(CASE WHEN is_resolved = true THEN 1 END) as resolved_alerts,
			COUNT(CASE WHEN severity = 'critical' AND is_resolved = false THEN 1 END) as critical_alerts,
			COUNT(CASE WHEN severity = 'warning' AND is_resolved = false THEN 1 END) as warning_alerts,
			COUNT(CASE WHEN alert_time >= NOW() - INTERVAL '24 hours' THEN 1 END) as alerts_24h,
			COUNT(CASE WHEN alert_time >= NOW() - INTERVAL '7 days' THEN 1 END) as alerts_7d
		` + baseQuery

	var stats struct {
		TotalAlerts    int `json:"total_alerts"`
		ActiveAlerts   int `json:"active_alerts"`
		ResolvedAlerts int `json:"resolved_alerts"`
		CriticalAlerts int `json:"critical_alerts"`
		WarningAlerts  int `json:"warning_alerts"`
		Alerts24h      int `json:"alerts_24h"`
		Alerts7d       int `json:"alerts_7d"`
	}

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&stats.TotalAlerts,
		&stats.ActiveAlerts,
		&stats.ResolvedAlerts,
		&stats.CriticalAlerts,
		&stats.WarningAlerts,
		&stats.Alerts24h,
		&stats.Alerts7d,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get alert statistics: %w", err)
	}

	result := map[string]interface{}{
		"total_alerts":    stats.TotalAlerts,
		"active_alerts":   stats.ActiveAlerts,
		"resolved_alerts": stats.ResolvedAlerts,
		"critical_alerts": stats.CriticalAlerts,
		"warning_alerts":  stats.WarningAlerts,
		"alerts_24h":      stats.Alerts24h,
		"alerts_7d":       stats.Alerts7d,
	}

	return result, nil
}

// ListAll retrieves all asset alerts across all tenants
func (r *assetAlertRepository) ListAll(ctx context.Context, limit, offset int) ([]*entity.AssetAlert, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM asset_alerts`
	var totalCount int
	err := r.DB.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, tenant_id, asset_id, asset_sensor_id, threshold_id,
			   measurement_field_name, alert_time, resolved_time, severity,
			   trigger_value, threshold_min_value, threshold_max_value,
			   alert_message, alert_type, is_resolved, created_at, updated_at
		FROM asset_alerts
		ORDER BY alert_time DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query asset alerts: %w", err)
	}
	defer rows.Close()

	alerts, err := r.scanRows(rows)
	if err != nil {
		return nil, 0, err
	}

	return alerts, totalCount, nil
}

// GetGlobalAlertStatistics retrieves alert statistics across all tenants
func (r *assetAlertRepository) GetGlobalAlertStatistics(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_alerts,
			COUNT(CASE WHEN is_resolved = false THEN 1 END) as active_alerts,
			COUNT(CASE WHEN is_resolved = true THEN 1 END) as resolved_alerts,
			COUNT(CASE WHEN severity = 'critical' AND is_resolved = false THEN 1 END) as critical_alerts,
			COUNT(CASE WHEN severity = 'warning' AND is_resolved = false THEN 1 END) as warning_alerts,
			COUNT(CASE WHEN alert_time >= NOW() - INTERVAL '24 hours' THEN 1 END) as alerts_24h,
			COUNT(CASE WHEN alert_time >= NOW() - INTERVAL '7 days' THEN 1 END) as alerts_7d,
			COUNT(DISTINCT tenant_id) as total_tenants
		FROM asset_alerts`

	var stats struct {
		TotalAlerts    int `json:"total_alerts"`
		ActiveAlerts   int `json:"active_alerts"`
		ResolvedAlerts int `json:"resolved_alerts"`
		CriticalAlerts int `json:"critical_alerts"`
		WarningAlerts  int `json:"warning_alerts"`
		Alerts24h      int `json:"alerts_24h"`
		Alerts7d       int `json:"alerts_7d"`
		TotalTenants   int `json:"total_tenants"`
	}

	err := r.DB.QueryRowContext(ctx, query).Scan(
		&stats.TotalAlerts,
		&stats.ActiveAlerts,
		&stats.ResolvedAlerts,
		&stats.CriticalAlerts,
		&stats.WarningAlerts,
		&stats.Alerts24h,
		&stats.Alerts7d,
		&stats.TotalTenants,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get global alert statistics: %w", err)
	}

	result := map[string]interface{}{
		"total_alerts":    stats.TotalAlerts,
		"active_alerts":   stats.ActiveAlerts,
		"resolved_alerts": stats.ResolvedAlerts,
		"critical_alerts": stats.CriticalAlerts,
		"warning_alerts":  stats.WarningAlerts,
		"alerts_24h":      stats.Alerts24h,
		"alerts_7d":       stats.Alerts7d,
		"total_tenants":   stats.TotalTenants,
	}

	return result, nil
}

// Helper methods

// queryAlerts executes a query and returns alerts
func (r *assetAlertRepository) queryAlerts(ctx context.Context, query string, args ...interface{}) ([]*entity.AssetAlert, error) {
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query asset alerts: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// scanRows scans database rows into AssetAlert entities
func (r *assetAlertRepository) scanRows(rows *sql.Rows) ([]*entity.AssetAlert, error) {
	var alerts []*entity.AssetAlert

	for rows.Next() {
		var alert entity.AssetAlert
		err := rows.Scan(
			&alert.ID,
			&alert.TenantID,
			&alert.AssetID,
			&alert.AssetSensorID,
			&alert.ThresholdID,
			&alert.MeasurementFieldName,
			&alert.AlertTime,
			&alert.ResolvedTime,
			&alert.Severity,
			&alert.TriggerValue,
			&alert.ThresholdMinValue,
			&alert.ThresholdMaxValue,
			&alert.AlertMessage,
			&alert.AlertType,
			&alert.IsResolved,
			&alert.CreatedAt,
			&alert.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset alert: %w", err)
		}
		alerts = append(alerts, &alert)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return alerts, nil
}

// CreateAlertFromReading creates an alert from an IoT sensor reading
func (r *assetAlertRepository) CreateAlertFromReading(ctx context.Context, reading *entity.IoTSensorReading, threshold *entity.SensorThreshold, value float64) error {
	// Check if there's already an active alert for this threshold
	activeAlerts, err := r.GetActiveAlertsByAssetSensor(ctx, reading.AssetSensorID)
	if err != nil {
		return fmt.Errorf("failed to check active alerts: %w", err)
	}

	// Check if there's already an active alert for this specific threshold
	for _, alert := range activeAlerts {
		if alert.ThresholdID == threshold.ID && !alert.IsResolved {
			// Update existing alert with new value
			alert.TriggerValue = value
			alert.AlertTime = time.Now()
			alert.AlertMessage = fmt.Sprintf("Threshold breach continues: %s = %.2f (threshold: %.2f - %.2f)",
				threshold.MeasurementFieldName, value,
				*threshold.MinValue, *threshold.MaxValue)
			return r.Update(ctx, alert)
		}
	}

	// Create new alert
	alert := &entity.AssetAlert{
		ID:                   uuid.New(),
		TenantID:             *reading.TenantID,
		AssetSensorID:        reading.AssetSensorID,
		ThresholdID:          threshold.ID,
		MeasurementFieldName: threshold.MeasurementFieldName,
		AlertTime:            time.Now(),
		Severity:             threshold.Severity,
		TriggerValue:         value,
		ThresholdMinValue:    threshold.MinValue,
		ThresholdMaxValue:    threshold.MaxValue,
		IsResolved:           false,
		CreatedAt:            time.Now(),
	}

	// Set alert type and message based on threshold breach
	if threshold.MinValue != nil && value < *threshold.MinValue {
		alert.AlertType = "min_breach"
		alert.AlertMessage = fmt.Sprintf("Value below minimum threshold: %s = %.2f (threshold: %.2f)",
			threshold.MeasurementFieldName, value, *threshold.MinValue)
	} else if threshold.MaxValue != nil && value > *threshold.MaxValue {
		alert.AlertType = "max_breach"
		alert.AlertMessage = fmt.Sprintf("Value above maximum threshold: %s = %.2f (threshold: %.2f)",
			threshold.MeasurementFieldName, value, *threshold.MaxValue)
	}

	return r.Create(ctx, alert)
}

// ResolveAlertsForReading resolves alerts when reading returns to normal
func (r *assetAlertRepository) ResolveAlertsForReading(ctx context.Context, reading *entity.IoTSensorReading, threshold *entity.SensorThreshold, value float64) error {
	// Get active alerts for this sensor
	activeAlerts, err := r.GetActiveAlertsByAssetSensor(ctx, reading.AssetSensorID)
	if err != nil {
		return fmt.Errorf("failed to get active alerts: %w", err)
	}

	// Check if value is within normal range
	isNormal := true
	if threshold.MinValue != nil && value < *threshold.MinValue {
		isNormal = false
	}
	if threshold.MaxValue != nil && value > *threshold.MaxValue {
		isNormal = false
	}

	// If value is normal, resolve all active alerts for this threshold
	if isNormal {
		for _, alert := range activeAlerts {
			if alert.ThresholdID == threshold.ID && !alert.IsResolved {
				alert.IsResolved = true
				now := time.Now()
				alert.ResolvedTime = &now
				alert.AlertMessage = fmt.Sprintf("Alert resolved: %s returned to normal range (%.2f)",
					threshold.MeasurementFieldName, value)
				if err := r.Update(ctx, alert); err != nil {
					return fmt.Errorf("failed to resolve alert: %w", err)
				}
			}
		}
	}

	return nil
}

// ResolveMultipleAlerts resolves multiple asset alerts
func (r *assetAlertRepository) ResolveMultipleAlerts(ctx context.Context, alertIDs []uuid.UUID) (int, int, error) {
	query := `UPDATE asset_alerts SET is_resolved = true, resolved_time = $2, updated_at = $2 WHERE id = ANY($1) AND is_resolved = false`

	result, err := r.DB.ExecContext(ctx, query, pq.Array(alertIDs), time.Now())
	if err != nil {
		return 0, 0, fmt.Errorf("failed to resolve multiple alerts: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), len(alertIDs) - int(rowsAffected), nil
}

// DeleteMultipleAlerts deletes multiple asset alerts
func (r *assetAlertRepository) DeleteMultipleAlerts(ctx context.Context, alertIDs []uuid.UUID) (int, int, error) {
	query := `DELETE FROM asset_alerts WHERE id = ANY($1)`

	result, err := r.DB.ExecContext(ctx, query, pq.Array(alertIDs))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to delete multiple alerts: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), len(alertIDs) - int(rowsAffected), nil
}
