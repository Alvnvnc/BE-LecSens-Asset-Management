package dto

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/helpers/common"
	"time"

	"github.com/google/uuid"
)

// SensorStatusDTO represents the data transfer object for sensor status
type SensorStatusDTO struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      *uuid.UUID `json:"tenant_id,omitempty"`
	AssetSensorID uuid.UUID  `json:"asset_sensor_id"`

	// Battery Information
	BatteryLevel         *float64   `json:"battery_level,omitempty"`
	BatteryVoltage       *float64   `json:"battery_voltage,omitempty"`
	BatteryStatus        *string    `json:"battery_status,omitempty"`
	BatteryLastCharged   *time.Time `json:"battery_last_charged,omitempty"`
	BatteryEstimatedLife *int       `json:"battery_estimated_life,omitempty"`
	BatteryType          *string    `json:"battery_type,omitempty"`

	// Signal Strength Information
	SignalType      *string  `json:"signal_type,omitempty"`
	SignalRSSI      *int     `json:"signal_rssi,omitempty"`
	SignalSNR       *float64 `json:"signal_snr,omitempty"`
	SignalQuality   *int     `json:"signal_quality,omitempty"`
	SignalFrequency *float64 `json:"signal_frequency,omitempty"`
	SignalChannel   *int     `json:"signal_channel,omitempty"`
	SignalStatus    *string  `json:"signal_status,omitempty"`

	// Connection Information
	ConnectionType     *string    `json:"connection_type,omitempty"`
	ConnectionStatus   string     `json:"connection_status"`
	LastConnectedAt    *time.Time `json:"last_connected_at,omitempty"`
	LastDisconnectedAt *time.Time `json:"last_disconnected_at,omitempty"`
	CurrentIP          *string    `json:"current_ip,omitempty"`
	CurrentNetwork     *string    `json:"current_network,omitempty"`

	// Additional Status
	Temperature     *float64   `json:"temperature,omitempty"`
	Humidity        *float64   `json:"humidity,omitempty"`
	IsOnline        bool       `json:"is_online"`
	LastHeartbeat   *time.Time `json:"last_heartbeat,omitempty"`
	FirmwareVersion *string    `json:"firmware_version,omitempty"`
	ErrorCount      *int       `json:"error_count,omitempty"`
	LastErrorAt     *time.Time `json:"last_error_at,omitempty"`

	// Timestamps
	RecordedAt time.Time  `json:"recorded_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// CreateSensorStatusRequest represents the request structure for creating a sensor status
type CreateSensorStatusRequest struct {
	AssetSensorID uuid.UUID `json:"asset_sensor_id" binding:"required"`

	// Battery Information
	BatteryLevel         *float64   `json:"battery_level,omitempty" validate:"omitempty,min=0,max=100"`
	BatteryVoltage       *float64   `json:"battery_voltage,omitempty" validate:"omitempty,min=0"`
	BatteryStatus        *string    `json:"battery_status,omitempty" validate:"omitempty,oneof=good low critical unknown charging"`
	BatteryLastCharged   *time.Time `json:"battery_last_charged,omitempty"`
	BatteryEstimatedLife *int       `json:"battery_estimated_life,omitempty" validate:"omitempty,min=0"`
	BatteryType          *string    `json:"battery_type,omitempty" validate:"omitempty,oneof=lithium alkaline rechargeable solar"`

	// Signal Strength Information
	SignalType      *string  `json:"signal_type,omitempty" validate:"omitempty,oneof=wifi cellular lora zigbee bluetooth"`
	SignalRSSI      *int     `json:"signal_rssi,omitempty" validate:"omitempty,min=-120,max=0"`
	SignalSNR       *float64 `json:"signal_snr,omitempty" validate:"omitempty,min=-40,max=40"`
	SignalQuality   *int     `json:"signal_quality,omitempty" validate:"omitempty,min=0,max=100"`
	SignalFrequency *float64 `json:"signal_frequency,omitempty" validate:"omitempty,min=0"`
	SignalChannel   *int     `json:"signal_channel,omitempty" validate:"omitempty,min=0"`
	SignalStatus    *string  `json:"signal_status,omitempty" validate:"omitempty,oneof=excellent good fair poor no_signal"`

	// Connection Information
	ConnectionType     *string    `json:"connection_type,omitempty" validate:"omitempty,oneof=wifi cellular ethernet lora zigbee bluetooth"`
	ConnectionStatus   string     `json:"connection_status" binding:"required" validate:"required,oneof=online offline connecting error"`
	LastConnectedAt    *time.Time `json:"last_connected_at,omitempty"`
	LastDisconnectedAt *time.Time `json:"last_disconnected_at,omitempty"`
	CurrentIP          *string    `json:"current_ip,omitempty" validate:"omitempty,ip"`
	CurrentNetwork     *string    `json:"current_network,omitempty"`

	// Additional Status
	Temperature     *float64   `json:"temperature,omitempty" validate:"omitempty,min=-50,max=100"`
	Humidity        *float64   `json:"humidity,omitempty" validate:"omitempty,min=0,max=100"`
	IsOnline        bool       `json:"is_online"`
	LastHeartbeat   *time.Time `json:"last_heartbeat,omitempty"`
	FirmwareVersion *string    `json:"firmware_version,omitempty"`
	ErrorCount      *int       `json:"error_count,omitempty" validate:"omitempty,min=0"`
	LastErrorAt     *time.Time `json:"last_error_at,omitempty"`

	// Timestamps
	RecordedAt *time.Time `json:"recorded_at,omitempty"`
}

// UpdateSensorStatusRequest represents the request structure for updating a sensor status
type UpdateSensorStatusRequest struct {
	// Battery Information
	BatteryLevel         *float64   `json:"battery_level,omitempty" validate:"omitempty,min=0,max=100"`
	BatteryVoltage       *float64   `json:"battery_voltage,omitempty" validate:"omitempty,min=0"`
	BatteryStatus        *string    `json:"battery_status,omitempty" validate:"omitempty,oneof=good low critical unknown charging"`
	BatteryLastCharged   *time.Time `json:"battery_last_charged,omitempty"`
	BatteryEstimatedLife *int       `json:"battery_estimated_life,omitempty" validate:"omitempty,min=0"`
	BatteryType          *string    `json:"battery_type,omitempty" validate:"omitempty,oneof=lithium alkaline rechargeable solar"`

	// Signal Strength Information
	SignalType      *string  `json:"signal_type,omitempty" validate:"omitempty,oneof=wifi cellular lora zigbee bluetooth"`
	SignalRSSI      *int     `json:"signal_rssi,omitempty" validate:"omitempty,min=-120,max=0"`
	SignalSNR       *float64 `json:"signal_snr,omitempty" validate:"omitempty,min=-40,max=40"`
	SignalQuality   *int     `json:"signal_quality,omitempty" validate:"omitempty,min=0,max=100"`
	SignalFrequency *float64 `json:"signal_frequency,omitempty" validate:"omitempty,min=0"`
	SignalChannel   *int     `json:"signal_channel,omitempty" validate:"omitempty,min=0"`
	SignalStatus    *string  `json:"signal_status,omitempty" validate:"omitempty,oneof=excellent good fair poor no_signal"`

	// Connection Information
	ConnectionType     *string    `json:"connection_type,omitempty" validate:"omitempty,oneof=wifi cellular ethernet lora zigbee bluetooth"`
	ConnectionStatus   *string    `json:"connection_status,omitempty" validate:"omitempty,oneof=online offline connecting error"`
	LastConnectedAt    *time.Time `json:"last_connected_at,omitempty"`
	LastDisconnectedAt *time.Time `json:"last_disconnected_at,omitempty"`
	CurrentIP          *string    `json:"current_ip,omitempty" validate:"omitempty,ip"`
	CurrentNetwork     *string    `json:"current_network,omitempty"`

	// Additional Status
	Temperature     *float64   `json:"temperature,omitempty" validate:"omitempty,min=-50,max=100"`
	Humidity        *float64   `json:"humidity,omitempty" validate:"omitempty,min=0,max=100"`
	IsOnline        *bool      `json:"is_online,omitempty"`
	LastHeartbeat   *time.Time `json:"last_heartbeat,omitempty"`
	FirmwareVersion *string    `json:"firmware_version,omitempty"`
	ErrorCount      *int       `json:"error_count,omitempty" validate:"omitempty,min=0"`
	LastErrorAt     *time.Time `json:"last_error_at,omitempty"`

	// Timestamps
	RecordedAt *time.Time `json:"recorded_at,omitempty"`
}

// SensorStatusFilter represents filter parameters for listing sensor status
type SensorStatusFilter struct {
	AssetSensorID       *uuid.UUID `json:"asset_sensor_id,omitempty"`
	ConnectionStatus    *string    `json:"connection_status,omitempty"`
	IsOnline            *bool      `json:"is_online,omitempty"`
	BatteryLevelMin     *float64   `json:"battery_level_min,omitempty"`
	BatteryLevelMax     *float64   `json:"battery_level_max,omitempty"`
	BatteryStatus       *string    `json:"battery_status,omitempty"`
	SignalQualityMin    *int       `json:"signal_quality_min,omitempty"`
	SignalQualityMax    *int       `json:"signal_quality_max,omitempty"`
	SignalStatus        *string    `json:"signal_status,omitempty"`
	ConnectionType      *string    `json:"connection_type,omitempty"`
	FirmwareVersion     *string    `json:"firmware_version,omitempty"`
	HasErrors           *bool      `json:"has_errors,omitempty"`
	LastHeartbeatAfter  *time.Time `json:"last_heartbeat_after,omitempty"`
	LastHeartbeatBefore *time.Time `json:"last_heartbeat_before,omitempty"`
	RecordedAfter       *time.Time `json:"recorded_after,omitempty"`
	RecordedBefore      *time.Time `json:"recorded_before,omitempty"`
	common.QueryParams
}

// SensorStatusResponse represents the response structure for sensor status operations
type SensorStatusResponse struct {
	Data    SensorStatusDTO `json:"data"`
	Message string          `json:"message"`
}

// SensorStatusListResponse represents the response structure for listing sensor status with pagination
type SensorStatusListResponse struct {
	Data       []SensorStatusDTO         `json:"data"`
	Pagination common.PaginationResponse `json:"pagination"`
	Message    string                    `json:"message"`
}

// SensorHealthSummaryResponse provides aggregated health information
type SensorHealthSummaryResponse struct {
	TotalSensors      int     `json:"total_sensors"`
	OnlineSensors     int     `json:"online_sensors"`
	OfflineSensors    int     `json:"offline_sensors"`
	LowBattery        int     `json:"low_battery"`
	CriticalBattery   int     `json:"critical_battery"`
	WeakSignal        int     `json:"weak_signal"`
	ErrorSensors      int     `json:"error_sensors"`
	HealthyPercentage float64 `json:"healthy_percentage"`
}

// UpsertSensorStatusRequest represents the request structure for creating or updating a sensor status
type UpsertSensorStatusRequest struct {
	AssetSensorID uuid.UUID `json:"asset_sensor_id" binding:"required"`

	// Battery Information
	BatteryLevel         *float64   `json:"battery_level,omitempty" validate:"omitempty,min=0,max=100"`
	BatteryVoltage       *float64   `json:"battery_voltage,omitempty" validate:"omitempty,min=0"`
	BatteryStatus        *string    `json:"battery_status,omitempty" validate:"omitempty,oneof=good low critical unknown charging"`
	BatteryLastCharged   *time.Time `json:"battery_last_charged,omitempty"`
	BatteryEstimatedLife *int       `json:"battery_estimated_life,omitempty" validate:"omitempty,min=0"`
	BatteryType          *string    `json:"battery_type,omitempty" validate:"omitempty,oneof=lithium alkaline rechargeable solar"`

	// Signal Strength Information
	SignalType      *string  `json:"signal_type,omitempty" validate:"omitempty,oneof=wifi cellular lora zigbee bluetooth"`
	SignalRSSI      *int     `json:"signal_rssi,omitempty" validate:"omitempty,min=-120,max=0"`
	SignalSNR       *float64 `json:"signal_snr,omitempty" validate:"omitempty,min=-40,max=40"`
	SignalQuality   *int     `json:"signal_quality,omitempty" validate:"omitempty,min=0,max=100"`
	SignalFrequency *float64 `json:"signal_frequency,omitempty" validate:"omitempty,min=0"`
	SignalChannel   *int     `json:"signal_channel,omitempty" validate:"omitempty,min=0"`
	SignalStatus    *string  `json:"signal_status,omitempty" validate:"omitempty,oneof=excellent good fair poor no_signal"`

	// Connection Information
	ConnectionType     *string    `json:"connection_type,omitempty" validate:"omitempty,oneof=wifi cellular ethernet lora zigbee bluetooth"`
	ConnectionStatus   string     `json:"connection_status" binding:"required" validate:"required,oneof=online offline connecting error"`
	LastConnectedAt    *time.Time `json:"last_connected_at,omitempty"`
	LastDisconnectedAt *time.Time `json:"last_disconnected_at,omitempty"`
	CurrentIP          *string    `json:"current_ip,omitempty" validate:"omitempty,ip"`
	CurrentNetwork     *string    `json:"current_network,omitempty"`

	// Additional Status
	Temperature     *float64   `json:"temperature,omitempty" validate:"omitempty,min=-50,max=100"`
	Humidity        *float64   `json:"humidity,omitempty" validate:"omitempty,min=0,max=100"`
	IsOnline        bool       `json:"is_online"`
	LastHeartbeat   *time.Time `json:"last_heartbeat,omitempty"`
	FirmwareVersion *string    `json:"firmware_version,omitempty"`
	ErrorCount      *int       `json:"error_count,omitempty" validate:"omitempty,min=0"`
	LastErrorAt     *time.Time `json:"last_error_at,omitempty"`

	// Timestamps
	RecordedAt *time.Time `json:"recorded_at,omitempty"`
}

// SensorHealthAnalyticsResponse provides detailed health analytics over time
type SensorHealthAnalyticsResponse struct {
	Timeframe    string                  `json:"timeframe"` // "24h", "7d", "30d"
	TotalSensors int                     `json:"total_sensors"`
	Analytics    []SensorHealthDataPoint `json:"analytics"`
	Summary      SensorHealthSummary     `json:"summary"`
	Trends       SensorHealthTrends      `json:"trends"`
}

// SensorHealthDataPoint represents a single point in time for health analytics
type SensorHealthDataPoint struct {
	Timestamp       time.Time `json:"timestamp"`
	OnlineSensors   int       `json:"online_sensors"`
	OfflineSensors  int       `json:"offline_sensors"`
	LowBattery      int       `json:"low_battery"`
	CriticalBattery int       `json:"critical_battery"`
	WeakSignal      int       `json:"weak_signal"`
	ErrorSensors    int       `json:"error_sensors"`
}

// SensorHealthSummary provides summary statistics for the timeframe
type SensorHealthSummary struct {
	AverageOnlinePercentage     float64 `json:"average_online_percentage"`
	AverageLowBatteryPercentage float64 `json:"average_low_battery_percentage"`
	AverageWeakSignalPercentage float64 `json:"average_weak_signal_percentage"`
	TotalDowntimeHours          float64 `json:"total_downtime_hours"`
	MostCommonIssue             string  `json:"most_common_issue"`
}

// SensorHealthTrends provides trend information
type SensorHealthTrends struct {
	OnlinePercentageChange     float64 `json:"online_percentage_change"`      // +/- percentage change
	LowBatteryPercentageChange float64 `json:"low_battery_percentage_change"` // +/- percentage change
	WeakSignalPercentageChange float64 `json:"weak_signal_percentage_change"` // +/- percentage change
	ErrorCountChange           int     `json:"error_count_change"`            // +/- change in error count
}

// ToEntity converts CreateSensorStatusRequest to entity.SensorStatus
func (r *CreateSensorStatusRequest) ToEntity(tenantID *uuid.UUID) *entity.SensorStatus {
	now := time.Now()
	recordedAt := now
	if r.RecordedAt != nil {
		recordedAt = *r.RecordedAt
	}

	return &entity.SensorStatus{
		ID:                   uuid.New(),
		TenantID:             tenantID,
		AssetSensorID:        r.AssetSensorID,
		BatteryLevel:         r.BatteryLevel,
		BatteryVoltage:       r.BatteryVoltage,
		BatteryStatus:        r.BatteryStatus,
		BatteryLastCharged:   r.BatteryLastCharged,
		BatteryEstimatedLife: r.BatteryEstimatedLife,
		BatteryType:          r.BatteryType,
		SignalType:           r.SignalType,
		SignalRSSI:           r.SignalRSSI,
		SignalSNR:            r.SignalSNR,
		SignalQuality:        r.SignalQuality,
		SignalFrequency:      r.SignalFrequency,
		SignalChannel:        r.SignalChannel,
		SignalStatus:         r.SignalStatus,
		ConnectionType:       r.ConnectionType,
		ConnectionStatus:     r.ConnectionStatus,
		LastConnectedAt:      r.LastConnectedAt,
		LastDisconnectedAt:   r.LastDisconnectedAt,
		CurrentIP:            r.CurrentIP,
		CurrentNetwork:       r.CurrentNetwork,
		Temperature:          r.Temperature,
		Humidity:             r.Humidity,
		IsOnline:             r.IsOnline,
		LastHeartbeat:        r.LastHeartbeat,
		FirmwareVersion:      r.FirmwareVersion,
		ErrorCount:           r.ErrorCount,
		LastErrorAt:          r.LastErrorAt,
		RecordedAt:           recordedAt,
		CreatedAt:            now,
	}
}

// ToEntity converts UpdateSensorStatusRequest to entity.SensorStatus for partial updates
func (r *UpdateSensorStatusRequest) ToEntity(id uuid.UUID) *entity.SensorStatus {
	now := time.Now()
	recordedAt := now
	if r.RecordedAt != nil {
		recordedAt = *r.RecordedAt
	}

	status := &entity.SensorStatus{
		ID:         id,
		RecordedAt: recordedAt,
		UpdatedAt:  &now,
	}

	// Only update fields that are provided
	if r.BatteryLevel != nil {
		status.BatteryLevel = r.BatteryLevel
	}
	if r.BatteryVoltage != nil {
		status.BatteryVoltage = r.BatteryVoltage
	}
	if r.BatteryStatus != nil {
		status.BatteryStatus = r.BatteryStatus
	}
	if r.BatteryLastCharged != nil {
		status.BatteryLastCharged = r.BatteryLastCharged
	}
	if r.BatteryEstimatedLife != nil {
		status.BatteryEstimatedLife = r.BatteryEstimatedLife
	}
	if r.BatteryType != nil {
		status.BatteryType = r.BatteryType
	}
	if r.SignalType != nil {
		status.SignalType = r.SignalType
	}
	if r.SignalRSSI != nil {
		status.SignalRSSI = r.SignalRSSI
	}
	if r.SignalSNR != nil {
		status.SignalSNR = r.SignalSNR
	}
	if r.SignalQuality != nil {
		status.SignalQuality = r.SignalQuality
	}
	if r.SignalFrequency != nil {
		status.SignalFrequency = r.SignalFrequency
	}
	if r.SignalChannel != nil {
		status.SignalChannel = r.SignalChannel
	}
	if r.SignalStatus != nil {
		status.SignalStatus = r.SignalStatus
	}
	if r.ConnectionType != nil {
		status.ConnectionType = r.ConnectionType
	}
	if r.ConnectionStatus != nil {
		status.ConnectionStatus = *r.ConnectionStatus
	}
	if r.LastConnectedAt != nil {
		status.LastConnectedAt = r.LastConnectedAt
	}
	if r.LastDisconnectedAt != nil {
		status.LastDisconnectedAt = r.LastDisconnectedAt
	}
	if r.CurrentIP != nil {
		status.CurrentIP = r.CurrentIP
	}
	if r.CurrentNetwork != nil {
		status.CurrentNetwork = r.CurrentNetwork
	}
	if r.Temperature != nil {
		status.Temperature = r.Temperature
	}
	if r.Humidity != nil {
		status.Humidity = r.Humidity
	}
	if r.IsOnline != nil {
		status.IsOnline = *r.IsOnline
	}
	if r.LastHeartbeat != nil {
		status.LastHeartbeat = r.LastHeartbeat
	}
	if r.FirmwareVersion != nil {
		status.FirmwareVersion = r.FirmwareVersion
	}
	if r.ErrorCount != nil {
		status.ErrorCount = r.ErrorCount
	}
	if r.LastErrorAt != nil {
		status.LastErrorAt = r.LastErrorAt
	}

	return status
}

// FromEntity converts entity.SensorStatus to SensorStatusDTO
func FromSensorStatusEntity(e *entity.SensorStatus) *SensorStatusDTO {
	if e == nil {
		return nil
	}
	return &SensorStatusDTO{
		ID:                   e.ID,
		TenantID:             e.TenantID,
		AssetSensorID:        e.AssetSensorID,
		BatteryLevel:         e.BatteryLevel,
		BatteryVoltage:       e.BatteryVoltage,
		BatteryStatus:        e.BatteryStatus,
		BatteryLastCharged:   e.BatteryLastCharged,
		BatteryEstimatedLife: e.BatteryEstimatedLife,
		BatteryType:          e.BatteryType,
		SignalType:           e.SignalType,
		SignalRSSI:           e.SignalRSSI,
		SignalSNR:            e.SignalSNR,
		SignalQuality:        e.SignalQuality,
		SignalFrequency:      e.SignalFrequency,
		SignalChannel:        e.SignalChannel,
		SignalStatus:         e.SignalStatus,
		ConnectionType:       e.ConnectionType,
		ConnectionStatus:     e.ConnectionStatus,
		LastConnectedAt:      e.LastConnectedAt,
		LastDisconnectedAt:   e.LastDisconnectedAt,
		CurrentIP:            e.CurrentIP,
		CurrentNetwork:       e.CurrentNetwork,
		Temperature:          e.Temperature,
		Humidity:             e.Humidity,
		IsOnline:             e.IsOnline,
		LastHeartbeat:        e.LastHeartbeat,
		FirmwareVersion:      e.FirmwareVersion,
		ErrorCount:           e.ErrorCount,
		LastErrorAt:          e.LastErrorAt,
		RecordedAt:           e.RecordedAt,
		CreatedAt:            e.CreatedAt,
		UpdatedAt:            e.UpdatedAt,
	}
}

// FromEntityList converts a slice of entity.SensorStatus to SensorStatusDTO slice
func FromSensorStatusEntityList(entities []*entity.SensorStatus) []SensorStatusDTO {
	dtos := make([]SensorStatusDTO, len(entities))
	for i, entity := range entities {
		if dto := FromSensorStatusEntity(entity); dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}

// ToCreateRequest converts UpsertSensorStatusRequest to CreateSensorStatusRequest
func (r *UpsertSensorStatusRequest) ToCreateRequest() *CreateSensorStatusRequest {
	return &CreateSensorStatusRequest{
		AssetSensorID:        r.AssetSensorID,
		BatteryLevel:         r.BatteryLevel,
		BatteryVoltage:       r.BatteryVoltage,
		BatteryStatus:        r.BatteryStatus,
		BatteryLastCharged:   r.BatteryLastCharged,
		BatteryEstimatedLife: r.BatteryEstimatedLife,
		BatteryType:          r.BatteryType,
		SignalType:           r.SignalType,
		SignalRSSI:           r.SignalRSSI,
		SignalSNR:            r.SignalSNR,
		SignalQuality:        r.SignalQuality,
		SignalFrequency:      r.SignalFrequency,
		SignalChannel:        r.SignalChannel,
		SignalStatus:         r.SignalStatus,
		ConnectionType:       r.ConnectionType,
		ConnectionStatus:     r.ConnectionStatus,
		LastConnectedAt:      r.LastConnectedAt,
		LastDisconnectedAt:   r.LastDisconnectedAt,
		CurrentIP:            r.CurrentIP,
		CurrentNetwork:       r.CurrentNetwork,
		Temperature:          r.Temperature,
		Humidity:             r.Humidity,
		IsOnline:             r.IsOnline,
		LastHeartbeat:        r.LastHeartbeat,
		FirmwareVersion:      r.FirmwareVersion,
		ErrorCount:           r.ErrorCount,
		LastErrorAt:          r.LastErrorAt,
		RecordedAt:           r.RecordedAt,
	}
}
