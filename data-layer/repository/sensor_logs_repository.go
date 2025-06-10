package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/helpers/common"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SensorLogsRepository defines the interface for sensor logs data operations
type SensorLogsRepository interface {
	Create(ctx context.Context, log *entity.SensorLogs) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SensorLogs, error)
	GetBySensorID(ctx context.Context, assetSensorID uuid.UUID, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error)
	GetByLogType(ctx context.Context, logType string, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error)
	GetByLogLevel(ctx context.Context, logLevel string, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error)
	GetConnectionHistory(ctx context.Context, assetSensorID uuid.UUID, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error)
	SearchLogs(ctx context.Context, searchQuery string, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error)
	GetErrorLogs(ctx context.Context, assetSensorID *uuid.UUID, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error)
	DeleteOldLogs(ctx context.Context, olderThan time.Time) (int64, error)
	Update(ctx context.Context, log *entity.SensorLogs) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAssetSensorContext(ctx context.Context, assetSensorID uuid.UUID) (*uuid.UUID, *uuid.UUID, error)
}

// sensorLogsRepository implements SensorLogsRepository
type sensorLogsRepository struct {
	db *sql.DB
}

// NewSensorLogsRepository creates a new sensor logs repository
func NewSensorLogsRepository(db *sql.DB) SensorLogsRepository {
	return &sensorLogsRepository{db: db}
}

// GetAssetSensorContext fetches tenant_id and asset_id for an asset sensor
func (r *sensorLogsRepository) GetAssetSensorContext(ctx context.Context, assetSensorID uuid.UUID) (*uuid.UUID, *uuid.UUID, error) {
	query := `
		SELECT asn.tenant_id, asn.asset_id 
		FROM asset_sensors asn 
		WHERE asn.id = $1
	`

	var tenantID, assetID *uuid.UUID
	err := r.db.QueryRowContext(ctx, query, assetSensorID).Scan(&tenantID, &assetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("asset sensor not found: %s", assetSensorID)
		}
		return nil, nil, fmt.Errorf("failed to get asset sensor context: %w", err)
	}

	return tenantID, assetID, nil
}

// Create inserts a new sensor log entry with automatic tenant_id inheritance
func (r *sensorLogsRepository) Create(ctx context.Context, log *entity.SensorLogs) error {
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}

	// Automatically inherit tenant_id if not set
	if log.TenantID == nil {
		tenantID, _, err := r.GetAssetSensorContext(ctx, log.AssetSensorID)
		if err != nil {
			return fmt.Errorf("failed to inherit tenant context: %w", err)
		}
		log.TenantID = tenantID
	}

	query := `
		INSERT INTO sensor_logs (
			id, tenant_id, asset_sensor_id, log_type, log_level, message, 
			component, event_type, error_code, connection_type, connection_status,
			ip_address, mac_address, network_name, connection_duration,
			metadata, source_ip, user_agent, session_id, recorded_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID, log.TenantID, log.AssetSensorID, log.LogType, log.LogLevel, log.Message,
		log.Component, log.EventType, log.ErrorCode, log.ConnectionType, log.ConnectionStatus,
		log.IPAddress, log.MACAddress, log.NetworkName, log.ConnectionDuration,
		log.Metadata, log.SourceIP, log.UserAgent, log.SessionID, log.RecordedAt, log.CreatedAt,
	)

	return err
}

// GetByID retrieves a sensor log by ID
func (r *sensorLogsRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SensorLogs, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, log_type, log_level, message,
			   component, event_type, error_code, connection_type, connection_status,
			   ip_address, mac_address, network_name, connection_duration,
			   metadata, source_ip, user_agent, session_id, recorded_at, created_at, updated_at
		FROM sensor_logs WHERE id = $1
	`

	log := &entity.SensorLogs{}
	var metadata []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID, &log.TenantID, &log.AssetSensorID, &log.LogType, &log.LogLevel, &log.Message,
		&log.Component, &log.EventType, &log.ErrorCode, &log.ConnectionType, &log.ConnectionStatus,
		&log.IPAddress, &log.MACAddress, &log.NetworkName, &log.ConnectionDuration,
		&metadata, &log.SourceIP, &log.UserAgent, &log.SessionID, &log.RecordedAt, &log.CreatedAt, &log.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if len(metadata) > 0 {
		log.Metadata = json.RawMessage(metadata)
	}

	return log, nil
}

// GetBySensorID retrieves logs for a specific sensor
func (r *sensorLogsRepository) GetBySensorID(ctx context.Context, assetSensorID uuid.UUID, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, log_type, log_level, message,
			   component, event_type, error_code, connection_type, connection_status,
			   ip_address, mac_address, network_name, connection_duration,
			   metadata, source_ip, user_agent, session_id, recorded_at, created_at, updated_at
		FROM sensor_logs 
		WHERE asset_sensor_id = $1
	`

	countQuery := `SELECT COUNT(*) FROM sensor_logs WHERE asset_sensor_id = $1`

	args := []interface{}{assetSensorID}

	return r.executeQuery(ctx, baseQuery, countQuery, args, params)
}

// GetByLogType retrieves logs by log type
func (r *sensorLogsRepository) GetByLogType(ctx context.Context, logType string, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, log_type, log_level, message,
			   component, event_type, error_code, connection_type, connection_status,
			   ip_address, mac_address, network_name, connection_duration,
			   metadata, source_ip, user_agent, session_id, recorded_at, created_at, updated_at
		FROM sensor_logs 
		WHERE log_type = $1
	`

	countQuery := `SELECT COUNT(*) FROM sensor_logs WHERE log_type = $1`

	args := []interface{}{logType}

	return r.executeQuery(ctx, baseQuery, countQuery, args, params)
}

// GetByLogLevel retrieves logs by log level
func (r *sensorLogsRepository) GetByLogLevel(ctx context.Context, logLevel string, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, log_type, log_level, message,
			   component, event_type, error_code, connection_type, connection_status,
			   ip_address, mac_address, network_name, connection_duration,
			   metadata, source_ip, user_agent, session_id, recorded_at, created_at, updated_at
		FROM sensor_logs 
		WHERE log_level = $1
	`

	countQuery := `SELECT COUNT(*) FROM sensor_logs WHERE log_level = $1`

	args := []interface{}{logLevel}

	return r.executeQuery(ctx, baseQuery, countQuery, args, params)
}

// GetConnectionHistory retrieves connection history logs for a sensor
func (r *sensorLogsRepository) GetConnectionHistory(ctx context.Context, assetSensorID uuid.UUID, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, log_type, log_level, message,
			   component, event_type, error_code, connection_type, connection_status,
			   ip_address, mac_address, network_name, connection_duration,
			   metadata, source_ip, user_agent, session_id, recorded_at, created_at, updated_at
		FROM sensor_logs 
		WHERE asset_sensor_id = $1 AND log_type = 'connection'
	`

	countQuery := `SELECT COUNT(*) FROM sensor_logs WHERE asset_sensor_id = $1 AND log_type = 'connection'`

	args := []interface{}{assetSensorID}

	return r.executeQuery(ctx, baseQuery, countQuery, args, params)
}

// SearchLogs searches logs by message content
func (r *sensorLogsRepository) SearchLogs(ctx context.Context, searchQuery string, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, log_type, log_level, message,
			   component, event_type, error_code, connection_type, connection_status,
			   ip_address, mac_address, network_name, connection_duration,
			   metadata, source_ip, user_agent, session_id, recorded_at, created_at, updated_at
		FROM sensor_logs 
		WHERE to_tsvector('english', message) @@ plainto_tsquery('english', $1)
	`

	countQuery := `SELECT COUNT(*) FROM sensor_logs WHERE to_tsvector('english', message) @@ plainto_tsquery('english', $1)`

	args := []interface{}{searchQuery}

	return r.executeQuery(ctx, baseQuery, countQuery, args, params)
}

// GetErrorLogs retrieves error and critical logs
func (r *sensorLogsRepository) GetErrorLogs(ctx context.Context, assetSensorID *uuid.UUID, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error) {
	var baseQuery, countQuery string
	var args []interface{}

	if assetSensorID != nil {
		baseQuery = `
			SELECT id, tenant_id, asset_sensor_id, log_type, log_level, message,
				   component, event_type, error_code, connection_type, connection_status,
				   ip_address, mac_address, network_name, connection_duration,
				   metadata, source_ip, user_agent, session_id, recorded_at, created_at, updated_at
			FROM sensor_logs 
			WHERE asset_sensor_id = $1 AND log_level IN ('error', 'critical')
		`
		countQuery = `SELECT COUNT(*) FROM sensor_logs WHERE asset_sensor_id = $1 AND log_level IN ('error', 'critical')`
		args = []interface{}{*assetSensorID}
	} else {
		baseQuery = `
			SELECT id, tenant_id, asset_sensor_id, log_type, log_level, message,
				   component, event_type, error_code, connection_type, connection_status,
				   ip_address, mac_address, network_name, connection_duration,
				   metadata, source_ip, user_agent, session_id, recorded_at, created_at, updated_at
			FROM sensor_logs 
			WHERE log_level IN ('error', 'critical')
		`
		countQuery = `SELECT COUNT(*) FROM sensor_logs WHERE log_level IN ('error', 'critical')`
		args = []interface{}{}
	}

	return r.executeQuery(ctx, baseQuery, countQuery, args, params)
}

// DeleteOldLogs deletes logs older than specified time
func (r *sensorLogsRepository) DeleteOldLogs(ctx context.Context, olderThan time.Time) (int64, error) {
	query := `DELETE FROM sensor_logs WHERE created_at < $1`

	result, err := r.db.ExecContext(ctx, query, olderThan)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

// Update updates a sensor log entry
func (r *sensorLogsRepository) Update(ctx context.Context, log *entity.SensorLogs) error {
	now := time.Now()
	log.UpdatedAt = &now

	query := `
		UPDATE sensor_logs SET 
			log_type = $2, log_level = $3, message = $4, component = $5, event_type = $6,
			error_code = $7, connection_type = $8, connection_status = $9, ip_address = $10,
			mac_address = $11, network_name = $12, connection_duration = $13, metadata = $14,
			source_ip = $15, user_agent = $16, session_id = $17, updated_at = $18
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID, log.LogType, log.LogLevel, log.Message, log.Component, log.EventType,
		log.ErrorCode, log.ConnectionType, log.ConnectionStatus, log.IPAddress,
		log.MACAddress, log.NetworkName, log.ConnectionDuration, log.Metadata,
		log.SourceIP, log.UserAgent, log.SessionID, log.UpdatedAt,
	)

	return err
}

// Delete deletes a sensor log entry
func (r *sensorLogsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sensor_logs WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// executeQuery is a helper method to execute paginated queries
func (r *sensorLogsRepository) executeQuery(ctx context.Context, baseQuery, countQuery string, args []interface{}, params common.QueryParams) ([]*entity.SensorLogs, *common.PaginationResponse, error) {
	// Count total records
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, nil, fmt.Errorf("failed to count logs: %w", err)
	}

	// Add ORDER BY and pagination
	query := baseQuery + ` ORDER BY recorded_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, params.PageSize, (params.Page-1)*params.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []*entity.SensorLogs
	for rows.Next() {
		log := &entity.SensorLogs{}
		var metadata []byte

		err := rows.Scan(
			&log.ID, &log.TenantID, &log.AssetSensorID, &log.LogType, &log.LogLevel, &log.Message,
			&log.Component, &log.EventType, &log.ErrorCode, &log.ConnectionType, &log.ConnectionStatus,
			&log.IPAddress, &log.MACAddress, &log.NetworkName, &log.ConnectionDuration,
			&metadata, &log.SourceIP, &log.UserAgent, &log.SessionID, &log.RecordedAt, &log.CreatedAt, &log.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan log: %w", err)
		}

		if len(metadata) > 0 {
			log.Metadata = json.RawMessage(metadata)
		}

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating logs: %w", err)
	}

	// Calculate pagination
	totalPages := (totalCount + int64(params.PageSize) - 1) / int64(params.PageSize)
	pagination := &common.PaginationResponse{
		Page:        params.Page,
		PageSize:    params.PageSize,
		TotalPages:  int(totalPages),
		Total:       totalCount,
		HasNext:     params.Page < int(totalPages),
		HasPrevious: params.Page > 1,
	}

	return logs, pagination, nil
}
