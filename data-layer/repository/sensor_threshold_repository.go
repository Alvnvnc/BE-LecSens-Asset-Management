package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// SensorThresholdRepository defines the interface for sensor threshold operations
type SensorThresholdRepository interface {
	Create(ctx context.Context, threshold *entity.SensorThreshold) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SensorThreshold, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.SensorThreshold, error)
	GetByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.SensorThreshold, error)
	GetByMeasurementTypeID(ctx context.Context, measurementTypeID uuid.UUID) ([]*entity.SensorThreshold, error)
	Update(ctx context.Context, threshold *entity.SensorThreshold) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.SensorThreshold, int, error)
	ListAll(ctx context.Context, limit, offset int) ([]*entity.SensorThreshold, int, error)
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
	log.Printf("Creating sensor threshold: %+v", threshold)

	// Validate threshold values
	if threshold.MinValue != nil && threshold.MaxValue != nil && *threshold.MinValue >= *threshold.MaxValue {
		return fmt.Errorf("minimum threshold must be less than maximum threshold")
	}

	// Ensure at least one threshold value is set
	if threshold.MinValue == nil && threshold.MaxValue == nil {
		return fmt.Errorf("at least one threshold value (min or max) must be set")
	}

	// Generate new UUID if not set
	if threshold.ID == uuid.Nil {
		threshold.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now()
	if threshold.CreatedAt.IsZero() {
		threshold.CreatedAt = now
	}

	// Check for existing threshold with same field and severity
	existsQuery := `
		SELECT EXISTS (
			SELECT 1 FROM sensor_thresholds 
			WHERE asset_sensor_id = $1 
			AND measurement_field_name = $2 
			AND severity = $3
		)`
	var exists bool
	err := r.DB.QueryRowContext(ctx, existsQuery,
		threshold.AssetSensorID,
		threshold.MeasurementFieldName,
		threshold.Severity).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check for existing threshold: %w", err)
	}
	if exists {
		return fmt.Errorf("threshold already exists for this field and severity")
	}

	query := `
		INSERT INTO sensor_thresholds (
			id, tenant_id, asset_sensor_id, measurement_type_id,
			measurement_field_name, min_value, max_value, severity,
			is_active, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`

	_, err = r.DB.ExecContext(ctx, query,
		threshold.ID,
		threshold.TenantID,
		threshold.AssetSensorID,
		threshold.MeasurementTypeID,
		threshold.MeasurementFieldName,
		threshold.MinValue,
		threshold.MaxValue,
		threshold.Severity,
		threshold.IsActive,
		threshold.CreatedAt,
	)

	if err != nil {
		log.Printf("Error creating sensor threshold: %v", err)
		return fmt.Errorf("failed to create sensor threshold: %w", err)
	}

	log.Printf("Successfully created sensor threshold with ID: %s", threshold.ID)
	return nil
}

// GetByID retrieves a sensor threshold by its ID
func (r *sensorThresholdRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SensorThreshold, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, measurement_type_id,
			   measurement_field_name, min_value, max_value, severity,
			   is_active, created_at, updated_at
		FROM sensor_thresholds
		WHERE id = $1`

	var threshold entity.SensorThreshold
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&threshold.ID,
		&threshold.TenantID,
		&threshold.AssetSensorID,
		&threshold.MeasurementTypeID,
		&threshold.MeasurementFieldName,
		&threshold.MinValue,
		&threshold.MaxValue,
		&threshold.Severity,
		&threshold.IsActive,
		&threshold.CreatedAt,
		&threshold.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sensor threshold: %w", err)
	}

	return &threshold, nil
}

// GetByTenantID retrieves all sensor thresholds for a tenant
func (r *sensorThresholdRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.SensorThreshold, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, measurement_type_id,
			   measurement_field_name, min_value, max_value, severity,
			   is_active, created_at, updated_at
		FROM sensor_thresholds
		WHERE tenant_id = $1
		ORDER BY created_at DESC`

	return r.queryThresholds(ctx, query, tenantID)
}

// GetByAssetSensorID retrieves thresholds for a specific asset sensor
func (r *sensorThresholdRepository) GetByAssetSensorID(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.SensorThreshold, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, measurement_type_id,
			   measurement_field_name, min_value, max_value, severity,
			   is_active, created_at, updated_at
		FROM sensor_thresholds
		WHERE asset_sensor_id = $1
		ORDER BY created_at DESC`

	return r.queryThresholds(ctx, query, assetSensorID)
}

// GetByMeasurementTypeID retrieves thresholds for a specific measurement type
func (r *sensorThresholdRepository) GetByMeasurementTypeID(ctx context.Context, measurementTypeID uuid.UUID) ([]*entity.SensorThreshold, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, measurement_type_id,
			   measurement_field_name, min_value, max_value, severity,
			   is_active, created_at, updated_at
		FROM sensor_thresholds
		WHERE measurement_type_id = $1
		ORDER BY created_at DESC`

	return r.queryThresholds(ctx, query, measurementTypeID)
}

// Update updates an existing sensor threshold
func (r *sensorThresholdRepository) Update(ctx context.Context, threshold *entity.SensorThreshold) error {
	// Validate threshold values
	if threshold.MinValue != nil && threshold.MaxValue != nil && *threshold.MinValue >= *threshold.MaxValue {
		return fmt.Errorf("minimum threshold must be less than maximum threshold")
	}

	// Ensure at least one threshold value is set
	if threshold.MinValue == nil && threshold.MaxValue == nil {
		return fmt.Errorf("at least one threshold value (min or max) must be set")
	}

	// Check if threshold exists
	existing, err := r.GetByID(ctx, threshold.ID)
	if err != nil {
		return fmt.Errorf("failed to check existing threshold: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("threshold not found")
	}

	// Check for duplicate threshold with same field and severity
	if existing.MeasurementFieldName != threshold.MeasurementFieldName || existing.Severity != threshold.Severity {
		existsQuery := `
			SELECT EXISTS (
				SELECT 1 FROM sensor_thresholds 
				WHERE asset_sensor_id = $1 
				AND measurement_field_name = $2 
				AND severity = $3
				AND id != $4
			)`
		var exists bool
		err := r.DB.QueryRowContext(ctx, existsQuery,
			threshold.AssetSensorID,
			threshold.MeasurementFieldName,
			threshold.Severity,
			threshold.ID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check for existing threshold: %w", err)
		}
		if exists {
			return fmt.Errorf("threshold already exists for this field and severity")
		}
	}

	now := time.Now()
	threshold.UpdatedAt = &now

	query := `
		UPDATE sensor_thresholds SET
			min_value = $2,
			max_value = $3,
			severity = $4,
			is_active = $5,
			updated_at = $6
		WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query,
		threshold.ID,
		threshold.MinValue,
		threshold.MaxValue,
		threshold.Severity,
		threshold.IsActive,
		threshold.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update sensor threshold: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sensor threshold not found")
	}

	return nil
}

// Delete removes a sensor threshold by its ID
func (r *sensorThresholdRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sensor_thresholds WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sensor threshold: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sensor threshold not found")
	}

	return nil
}

// List retrieves paginated sensor thresholds for a tenant
func (r *sensorThresholdRepository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.SensorThreshold, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM sensor_thresholds WHERE tenant_id = $1`
	var totalCount int
	err := r.DB.QueryRowContext(ctx, countQuery, tenantID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, tenant_id, asset_sensor_id, measurement_type_id,
			   measurement_field_name, min_value, max_value, severity,
			   is_active, created_at, updated_at
		FROM sensor_thresholds
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.DB.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query sensor thresholds: %w", err)
	}
	defer rows.Close()

	thresholds, err := r.scanRows(rows)
	if err != nil {
		return nil, 0, err
	}

	return thresholds, totalCount, nil
}

// ListAll retrieves all sensor thresholds across all tenants
func (r *sensorThresholdRepository) ListAll(ctx context.Context, limit, offset int) ([]*entity.SensorThreshold, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM sensor_thresholds`
	var totalCount int
	err := r.DB.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, tenant_id, asset_sensor_id, measurement_type_id,
			   measurement_field_name, min_value, max_value, severity,
			   is_active, created_at, updated_at
		FROM sensor_thresholds
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query sensor thresholds: %w", err)
	}
	defer rows.Close()

	thresholds, err := r.scanRows(rows)
	if err != nil {
		return nil, 0, err
	}

	return thresholds, totalCount, nil
}

// Helper methods

// queryThresholds executes a query and returns thresholds
func (r *sensorThresholdRepository) queryThresholds(ctx context.Context, query string, args ...interface{}) ([]*entity.SensorThreshold, error) {
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sensor thresholds: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// scanRows scans database rows into SensorThreshold entities
func (r *sensorThresholdRepository) scanRows(rows *sql.Rows) ([]*entity.SensorThreshold, error) {
	var thresholds []*entity.SensorThreshold

	for rows.Next() {
		var threshold entity.SensorThreshold
		err := rows.Scan(
			&threshold.ID,
			&threshold.TenantID,
			&threshold.AssetSensorID,
			&threshold.MeasurementTypeID,
			&threshold.MeasurementFieldName,
			&threshold.MinValue,
			&threshold.MaxValue,
			&threshold.Severity,
			&threshold.IsActive,
			&threshold.CreatedAt,
			&threshold.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sensor threshold: %w", err)
		}
		thresholds = append(thresholds, &threshold)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return thresholds, nil
}
