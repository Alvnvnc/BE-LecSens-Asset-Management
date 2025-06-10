package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/helpers/common"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SensorStatusRepository defines the interface for sensor status data operations
type SensorStatusRepository interface {
	Create(ctx context.Context, status *entity.SensorStatus) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SensorStatus, error)
	GetBySensorID(ctx context.Context, assetSensorID uuid.UUID) (*entity.SensorStatus, error)
	GetAssetSensorContext(ctx context.Context, assetSensorID uuid.UUID) (*uuid.UUID, *uuid.UUID, error)
	GetOnlineSensors(ctx context.Context, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error)
	GetOfflineSensors(ctx context.Context, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error)
	GetLowBatterySensors(ctx context.Context, threshold float64, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error)
	GetWeakSignalSensors(ctx context.Context, threshold int, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error)
	GetUnhealthySensors(ctx context.Context, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error)
	GetAll(ctx context.Context, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error)
	UpdateBatteryStatus(ctx context.Context, assetSensorID uuid.UUID, batteryLevel *float64, batteryVoltage *float64, batteryStatus *string) error
	UpdateSignalStatus(ctx context.Context, assetSensorID uuid.UUID, rssi *int, snr *float64, quality *int, signalStatus *string) error
	UpdateConnectionStatus(ctx context.Context, assetSensorID uuid.UUID, connectionStatus string, connectionType *string, currentIP *string, currentNetwork *string) error
	UpdateHeartbeat(ctx context.Context, assetSensorID uuid.UUID) error
	UpsertStatus(ctx context.Context, status *entity.SensorStatus) error
	Update(ctx context.Context, status *entity.SensorStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBySensorID(ctx context.Context, assetSensorID uuid.UUID) error
}

// sensorStatusRepository implements SensorStatusRepository
type sensorStatusRepository struct {
	db *sql.DB
}

// NewSensorStatusRepository creates a new sensor status repository
func NewSensorStatusRepository(db *sql.DB) SensorStatusRepository {
	return &sensorStatusRepository{db: db}
}

// GetAssetSensorContext fetches tenant_id and asset_id for an asset sensor
func (r *sensorStatusRepository) GetAssetSensorContext(ctx context.Context, assetSensorID uuid.UUID) (*uuid.UUID, *uuid.UUID, error) {
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

// Create inserts a new sensor status record with automatic tenant_id inheritance
func (r *sensorStatusRepository) Create(ctx context.Context, status *entity.SensorStatus) error {
	if status.ID == uuid.Nil {
		status.ID = uuid.New()
	}

	// Automatically inherit tenant_id if not set
	if status.TenantID == nil {
		tenantID, _, err := r.GetAssetSensorContext(ctx, status.AssetSensorID)
		if err != nil {
			return fmt.Errorf("failed to inherit tenant context: %w", err)
		}
		status.TenantID = tenantID
	}

	query := `
		INSERT INTO sensor_status (
			id, tenant_id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			connection_type, connection_status, last_connected_at, last_disconnected_at,
			current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			firmware_version, error_count, last_error_at, recorded_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)
	`

	_, err := r.db.ExecContext(ctx, query,
		status.ID, status.TenantID, status.AssetSensorID, status.BatteryLevel, status.BatteryVoltage, status.BatteryStatus,
		status.BatteryLastCharged, status.BatteryEstimatedLife, status.BatteryType, status.SignalType, status.SignalRSSI,
		status.SignalSNR, status.SignalQuality, status.SignalFrequency, status.SignalChannel, status.SignalStatus,
		status.ConnectionType, status.ConnectionStatus, status.LastConnectedAt, status.LastDisconnectedAt,
		status.CurrentIP, status.CurrentNetwork, status.Temperature, status.Humidity, status.IsOnline, status.LastHeartbeat,
		status.FirmwareVersion, status.ErrorCount, status.LastErrorAt, status.RecordedAt, status.CreatedAt,
	)

	return err
}

// GetByID retrieves a sensor status by ID
func (r *sensorStatusRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SensorStatus, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			   battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			   signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			   connection_type, connection_status, last_connected_at, last_disconnected_at,
			   current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			   firmware_version, error_count, last_error_at, recorded_at, created_at, updated_at
		FROM sensor_status WHERE id = $1
	`

	status := &entity.SensorStatus{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&status.ID, &status.TenantID, &status.AssetSensorID, &status.BatteryLevel, &status.BatteryVoltage, &status.BatteryStatus,
		&status.BatteryLastCharged, &status.BatteryEstimatedLife, &status.BatteryType, &status.SignalType, &status.SignalRSSI,
		&status.SignalSNR, &status.SignalQuality, &status.SignalFrequency, &status.SignalChannel, &status.SignalStatus,
		&status.ConnectionType, &status.ConnectionStatus, &status.LastConnectedAt, &status.LastDisconnectedAt,
		&status.CurrentIP, &status.CurrentNetwork, &status.Temperature, &status.Humidity, &status.IsOnline, &status.LastHeartbeat,
		&status.FirmwareVersion, &status.ErrorCount, &status.LastErrorAt, &status.RecordedAt, &status.CreatedAt, &status.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return status, nil
}

// GetBySensorID retrieves sensor status by asset sensor ID
func (r *sensorStatusRepository) GetBySensorID(ctx context.Context, assetSensorID uuid.UUID) (*entity.SensorStatus, error) {
	query := `
		SELECT id, tenant_id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			   battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			   signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			   connection_type, connection_status, last_connected_at, last_disconnected_at,
			   current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			   firmware_version, error_count, last_error_at, recorded_at, created_at, updated_at
		FROM sensor_status WHERE asset_sensor_id = $1
	`

	status := &entity.SensorStatus{}
	err := r.db.QueryRowContext(ctx, query, assetSensorID).Scan(
		&status.ID, &status.TenantID, &status.AssetSensorID, &status.BatteryLevel, &status.BatteryVoltage, &status.BatteryStatus,
		&status.BatteryLastCharged, &status.BatteryEstimatedLife, &status.BatteryType, &status.SignalType, &status.SignalRSSI,
		&status.SignalSNR, &status.SignalQuality, &status.SignalFrequency, &status.SignalChannel, &status.SignalStatus,
		&status.ConnectionType, &status.ConnectionStatus, &status.LastConnectedAt, &status.LastDisconnectedAt,
		&status.CurrentIP, &status.CurrentNetwork, &status.Temperature, &status.Humidity, &status.IsOnline, &status.LastHeartbeat,
		&status.FirmwareVersion, &status.ErrorCount, &status.LastErrorAt, &status.RecordedAt, &status.CreatedAt, &status.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return status, nil
}

// GetOnlineSensors retrieves all online sensors
func (r *sensorStatusRepository) GetOnlineSensors(ctx context.Context, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			   battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			   signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			   connection_type, connection_status, last_connected_at, last_disconnected_at,
			   current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			   firmware_version, error_count, last_error_at, recorded_at, created_at, updated_at
		FROM sensor_status 
		WHERE is_online = true
	`

	countQuery := `SELECT COUNT(*) FROM sensor_status WHERE is_online = true`

	return r.executeQuery(ctx, baseQuery, countQuery, []interface{}{}, params)
}

// GetOfflineSensors retrieves all offline sensors
func (r *sensorStatusRepository) GetOfflineSensors(ctx context.Context, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			   battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			   signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			   connection_type, connection_status, last_connected_at, last_disconnected_at,
			   current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			   firmware_version, error_count, last_error_at, recorded_at, created_at, updated_at
		FROM sensor_status 
		WHERE is_online = false
	`

	countQuery := `SELECT COUNT(*) FROM sensor_status WHERE is_online = false`

	return r.executeQuery(ctx, baseQuery, countQuery, []interface{}{}, params)
}

// GetLowBatterySensors retrieves sensors with low battery
func (r *sensorStatusRepository) GetLowBatterySensors(ctx context.Context, threshold float64, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			   battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			   signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			   connection_type, connection_status, last_connected_at, last_disconnected_at,
			   current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			   firmware_version, error_count, last_error_at, recorded_at, created_at, updated_at
		FROM sensor_status 
		WHERE battery_level IS NOT NULL AND battery_level < $1
	`

	countQuery := `SELECT COUNT(*) FROM sensor_status WHERE battery_level IS NOT NULL AND battery_level < $1`
	args := []interface{}{threshold}

	return r.executeQuery(ctx, baseQuery, countQuery, args, params)
}

// GetWeakSignalSensors retrieves sensors with weak signal
func (r *sensorStatusRepository) GetWeakSignalSensors(ctx context.Context, threshold int, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			   battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			   signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			   connection_type, connection_status, last_connected_at, last_disconnected_at,
			   current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			   firmware_version, error_count, last_error_at, recorded_at, created_at, updated_at
		FROM sensor_status 
		WHERE signal_rssi IS NOT NULL AND signal_rssi < $1
	`

	countQuery := `SELECT COUNT(*) FROM sensor_status WHERE signal_rssi IS NOT NULL AND signal_rssi < $1`
	args := []interface{}{threshold}

	return r.executeQuery(ctx, baseQuery, countQuery, args, params)
}

// GetUnhealthySensors retrieves sensors that are considered unhealthy
func (r *sensorStatusRepository) GetUnhealthySensors(ctx context.Context, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			   battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			   signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			   connection_type, connection_status, last_connected_at, last_disconnected_at,
			   current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			   firmware_version, error_count, last_error_at, recorded_at, created_at, updated_at
		FROM sensor_status 
		WHERE is_online = false 
		   OR (battery_level IS NOT NULL AND battery_level < 20.0)
		   OR (signal_rssi IS NOT NULL AND signal_rssi < -90)
		   OR (error_count IS NOT NULL AND error_count > 10)
	`

	countQuery := `
		SELECT COUNT(*) FROM sensor_status 
		WHERE is_online = false 
		   OR (battery_level IS NOT NULL AND battery_level < 20.0)
		   OR (signal_rssi IS NOT NULL AND signal_rssi < -90)
		   OR (error_count IS NOT NULL AND error_count > 10)
	`

	return r.executeQuery(ctx, baseQuery, countQuery, []interface{}{}, params)
}

// GetAll retrieves all sensor statuses with pagination
func (r *sensorStatusRepository) GetAll(ctx context.Context, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error) {
	baseQuery := `
		SELECT id, tenant_id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			   battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			   signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			   connection_type, connection_status, last_connected_at, last_disconnected_at,
			   current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			   firmware_version, error_count, last_error_at, recorded_at, created_at, updated_at
		FROM sensor_status
	`

	countQuery := `SELECT COUNT(*) FROM sensor_status`

	return r.executeQuery(ctx, baseQuery, countQuery, []interface{}{}, params)
}

// UpdateBatteryStatus updates battery-related fields
func (r *sensorStatusRepository) UpdateBatteryStatus(ctx context.Context, assetSensorID uuid.UUID, batteryLevel *float64, batteryVoltage *float64, batteryStatus *string) error {
	now := time.Now()

	query := `
		UPDATE sensor_status SET 
			battery_level = $2, battery_voltage = $3, battery_status = $4, 
			recorded_at = $5, updated_at = $6
		WHERE asset_sensor_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, assetSensorID, batteryLevel, batteryVoltage, batteryStatus, now, now)
	return err
}

// UpdateSignalStatus updates signal-related fields
func (r *sensorStatusRepository) UpdateSignalStatus(ctx context.Context, assetSensorID uuid.UUID, rssi *int, snr *float64, quality *int, signalStatus *string) error {
	now := time.Now()

	query := `
		UPDATE sensor_status SET 
			signal_rssi = $2, signal_snr = $3, signal_quality = $4, signal_status = $5,
			recorded_at = $6, updated_at = $7
		WHERE asset_sensor_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, assetSensorID, rssi, snr, quality, signalStatus, now, now)
	return err
}

// UpdateConnectionStatus updates connection-related fields
func (r *sensorStatusRepository) UpdateConnectionStatus(ctx context.Context, assetSensorID uuid.UUID, connectionStatus string, connectionType *string, currentIP *string, currentNetwork *string) error {
	now := time.Now()
	isOnline := connectionStatus == "online"

	var lastConnectedAt, lastDisconnectedAt *time.Time
	if isOnline {
		lastConnectedAt = &now
	} else {
		lastDisconnectedAt = &now
	}

	query := `
		UPDATE sensor_status SET 
			connection_status = $2, connection_type = $3, current_ip = $4, current_network = $5,
			is_online = $6, last_connected_at = COALESCE($7, last_connected_at), 
			last_disconnected_at = COALESCE($8, last_disconnected_at),
			recorded_at = $9, updated_at = $10
		WHERE asset_sensor_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, assetSensorID, connectionStatus, connectionType, currentIP, currentNetwork, isOnline, lastConnectedAt, lastDisconnectedAt, now, now)
	return err
}

// UpdateHeartbeat updates the last heartbeat timestamp
func (r *sensorStatusRepository) UpdateHeartbeat(ctx context.Context, assetSensorID uuid.UUID) error {
	now := time.Now()

	query := `
		UPDATE sensor_status SET 
			last_heartbeat = $2, is_online = true, recorded_at = $3, updated_at = $4
		WHERE asset_sensor_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, assetSensorID, now, now, now)
	return err
}

// UpsertStatus creates or updates sensor status with automatic tenant_id inheritance
func (r *sensorStatusRepository) UpsertStatus(ctx context.Context, status *entity.SensorStatus) error {
	if status.ID == uuid.Nil {
		status.ID = uuid.New()
	}

	// Automatically inherit tenant_id if not set
	if status.TenantID == nil {
		tenantID, _, err := r.GetAssetSensorContext(ctx, status.AssetSensorID)
		if err != nil {
			return fmt.Errorf("failed to inherit tenant context: %w", err)
		}
		status.TenantID = tenantID
	}

	now := time.Now()
	status.UpdatedAt = &now

	query := `
		INSERT INTO sensor_status (
			id, tenant_id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			connection_type, connection_status, last_connected_at, last_disconnected_at,
			current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			firmware_version, error_count, last_error_at, recorded_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32)
		ON CONFLICT (asset_sensor_id) DO UPDATE SET
			battery_level = EXCLUDED.battery_level,
			battery_voltage = EXCLUDED.battery_voltage,
			battery_status = EXCLUDED.battery_status,
			battery_last_charged = EXCLUDED.battery_last_charged,
			battery_estimated_life = EXCLUDED.battery_estimated_life,
			battery_type = EXCLUDED.battery_type,
			signal_type = EXCLUDED.signal_type,
			signal_rssi = EXCLUDED.signal_rssi,
			signal_snr = EXCLUDED.signal_snr,
			signal_quality = EXCLUDED.signal_quality,
			signal_frequency = EXCLUDED.signal_frequency,
			signal_channel = EXCLUDED.signal_channel,
			signal_status = EXCLUDED.signal_status,
			connection_type = EXCLUDED.connection_type,
			connection_status = EXCLUDED.connection_status,
			last_connected_at = EXCLUDED.last_connected_at,
			last_disconnected_at = EXCLUDED.last_disconnected_at,
			current_ip = EXCLUDED.current_ip,
			current_network = EXCLUDED.current_network,
			temperature = EXCLUDED.temperature,
			humidity = EXCLUDED.humidity,
			is_online = EXCLUDED.is_online,
			last_heartbeat = EXCLUDED.last_heartbeat,
			firmware_version = EXCLUDED.firmware_version,
			error_count = EXCLUDED.error_count,
			last_error_at = EXCLUDED.last_error_at,
			recorded_at = EXCLUDED.recorded_at,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		status.ID, status.TenantID, status.AssetSensorID, status.BatteryLevel, status.BatteryVoltage, status.BatteryStatus,
		status.BatteryLastCharged, status.BatteryEstimatedLife, status.BatteryType, status.SignalType, status.SignalRSSI,
		status.SignalSNR, status.SignalQuality, status.SignalFrequency, status.SignalChannel, status.SignalStatus,
		status.ConnectionType, status.ConnectionStatus, status.LastConnectedAt, status.LastDisconnectedAt,
		status.CurrentIP, status.CurrentNetwork, status.Temperature, status.Humidity, status.IsOnline, status.LastHeartbeat,
		status.FirmwareVersion, status.ErrorCount, status.LastErrorAt, status.RecordedAt, status.CreatedAt, status.UpdatedAt,
	)

	return err
}

// Update updates a sensor status record
func (r *sensorStatusRepository) Update(ctx context.Context, status *entity.SensorStatus) error {
	now := time.Now()
	status.UpdatedAt = &now

	query := `
		UPDATE sensor_status SET 
			battery_level = $2, battery_voltage = $3, battery_status = $4, battery_last_charged = $5,
			battery_estimated_life = $6, battery_type = $7, signal_type = $8, signal_rssi = $9,
			signal_snr = $10, signal_quality = $11, signal_frequency = $12, signal_channel = $13,
			signal_status = $14, connection_type = $15, connection_status = $16, last_connected_at = $17,
			last_disconnected_at = $18, current_ip = $19, current_network = $20, temperature = $21,
			humidity = $22, is_online = $23, last_heartbeat = $24, firmware_version = $25,
			error_count = $26, last_error_at = $27, recorded_at = $28, updated_at = $29
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		status.ID, status.BatteryLevel, status.BatteryVoltage, status.BatteryStatus, status.BatteryLastCharged,
		status.BatteryEstimatedLife, status.BatteryType, status.SignalType, status.SignalRSSI,
		status.SignalSNR, status.SignalQuality, status.SignalFrequency, status.SignalChannel,
		status.SignalStatus, status.ConnectionType, status.ConnectionStatus, status.LastConnectedAt,
		status.LastDisconnectedAt, status.CurrentIP, status.CurrentNetwork, status.Temperature,
		status.Humidity, status.IsOnline, status.LastHeartbeat, status.FirmwareVersion,
		status.ErrorCount, status.LastErrorAt, status.RecordedAt, status.UpdatedAt,
	)

	return err
}

// Delete deletes a sensor status record
func (r *sensorStatusRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sensor_status WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteBySensorID deletes a sensor status record by sensor ID
func (r *sensorStatusRepository) DeleteBySensorID(ctx context.Context, assetSensorID uuid.UUID) error {
	query := `DELETE FROM sensor_status WHERE asset_sensor_id = $1`
	_, err := r.db.ExecContext(ctx, query, assetSensorID)
	return err
}

// executeQuery is a helper method to execute paginated queries
func (r *sensorStatusRepository) executeQuery(ctx context.Context, baseQuery, countQuery string, args []interface{}, params common.QueryParams) ([]*entity.SensorStatus, *common.PaginationResponse, error) {
	// Count total records
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, nil, fmt.Errorf("failed to count sensor status: %w", err)
	}

	// Add ORDER BY and pagination
	query := baseQuery + ` ORDER BY recorded_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, params.PageSize, (params.Page-1)*params.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query sensor status: %w", err)
	}
	defer rows.Close()

	var statuses []*entity.SensorStatus
	for rows.Next() {
		status := &entity.SensorStatus{}

		err := rows.Scan(
			&status.ID, &status.TenantID, &status.AssetSensorID, &status.BatteryLevel, &status.BatteryVoltage, &status.BatteryStatus,
			&status.BatteryLastCharged, &status.BatteryEstimatedLife, &status.BatteryType, &status.SignalType, &status.SignalRSSI,
			&status.SignalSNR, &status.SignalQuality, &status.SignalFrequency, &status.SignalChannel, &status.SignalStatus,
			&status.ConnectionType, &status.ConnectionStatus, &status.LastConnectedAt, &status.LastDisconnectedAt,
			&status.CurrentIP, &status.CurrentNetwork, &status.Temperature, &status.Humidity, &status.IsOnline, &status.LastHeartbeat,
			&status.FirmwareVersion, &status.ErrorCount, &status.LastErrorAt, &status.RecordedAt, &status.CreatedAt, &status.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan sensor status: %w", err)
		}

		statuses = append(statuses, status)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating sensor status: %w", err)
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

	return statuses, pagination, nil
}
