package dto

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/helpers/common"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SensorLogsDTO represents the data transfer object for sensor logs
type SensorLogsDTO struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      *uuid.UUID `json:"tenant_id,omitempty"`
	AssetSensorID uuid.UUID  `json:"asset_sensor_id"`
	LogType       string     `json:"log_type"`
	LogLevel      string     `json:"log_level"`
	Message       string     `json:"message"`
	Component     *string    `json:"component,omitempty"`
	EventType     *string    `json:"event_type,omitempty"`
	ErrorCode     *string    `json:"error_code,omitempty"`

	// Connection History Data
	ConnectionType     *string `json:"connection_type,omitempty"`
	ConnectionStatus   *string `json:"connection_status,omitempty"`
	IPAddress          *string `json:"ip_address,omitempty"`
	MACAddress         *string `json:"mac_address,omitempty"`
	NetworkName        *string `json:"network_name,omitempty"`
	ConnectionDuration *int64  `json:"connection_duration,omitempty"`

	// Metadata
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	SourceIP  *string         `json:"source_ip,omitempty"`
	UserAgent *string         `json:"user_agent,omitempty"`
	SessionID *string         `json:"session_id,omitempty"`

	// Timestamps
	RecordedAt time.Time  `json:"recorded_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// CreateSensorLogsRequest represents the request structure for creating a sensor log
type CreateSensorLogsRequest struct {
	AssetSensorID uuid.UUID `json:"asset_sensor_id" binding:"required"`
	LogType       string    `json:"log_type" binding:"required" validate:"required,oneof=reading connection battery signal error system maintenance"`
	LogLevel      string    `json:"log_level" binding:"required" validate:"required,oneof=debug info warning error critical"`
	Message       string    `json:"message" binding:"required"`
	Component     *string   `json:"component,omitempty" validate:"omitempty,oneof=sensor communication battery hardware network software"`
	EventType     *string   `json:"event_type,omitempty" validate:"omitempty,oneof=startup shutdown connected disconnected reading error maintenance calibration reset"`
	ErrorCode     *string   `json:"error_code,omitempty"`

	// Connection History Data (when log_type = "connection")
	ConnectionType     *string `json:"connection_type,omitempty" validate:"omitempty,oneof=wifi cellular lora zigbee bluetooth ethernet"`
	ConnectionStatus   *string `json:"connection_status,omitempty" validate:"omitempty,oneof=connected disconnected failed timeout establishing"`
	IPAddress          *string `json:"ip_address,omitempty" validate:"omitempty,ip"`
	MACAddress         *string `json:"mac_address,omitempty" validate:"omitempty,mac"`
	NetworkName        *string `json:"network_name,omitempty"`
	ConnectionDuration *int64  `json:"connection_duration,omitempty" validate:"omitempty,min=0"`

	// Metadata for flexible additional data
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	SourceIP  *string         `json:"source_ip,omitempty" validate:"omitempty,ip"`
	UserAgent *string         `json:"user_agent,omitempty"`
	SessionID *string         `json:"session_id,omitempty"`

	// Timestamps
	RecordedAt *time.Time `json:"recorded_at,omitempty"`
}

// UpdateSensorLogsRequest represents the request structure for updating a sensor log
type UpdateSensorLogsRequest struct {
	LogType   *string `json:"log_type,omitempty" validate:"omitempty,oneof=reading connection battery signal error system maintenance"`
	LogLevel  *string `json:"log_level,omitempty" validate:"omitempty,oneof=debug info warning error critical"`
	Message   *string `json:"message,omitempty"`
	Component *string `json:"component,omitempty" validate:"omitempty,oneof=sensor communication battery hardware network software"`
	EventType *string `json:"event_type,omitempty" validate:"omitempty,oneof=startup shutdown connected disconnected reading error maintenance calibration reset"`
	ErrorCode *string `json:"error_code,omitempty"`

	// Connection History Data
	ConnectionType     *string `json:"connection_type,omitempty" validate:"omitempty,oneof=wifi cellular lora zigbee bluetooth ethernet"`
	ConnectionStatus   *string `json:"connection_status,omitempty" validate:"omitempty,oneof=connected disconnected failed timeout establishing"`
	IPAddress          *string `json:"ip_address,omitempty" validate:"omitempty,ip"`
	MACAddress         *string `json:"mac_address,omitempty" validate:"omitempty,mac"`
	NetworkName        *string `json:"network_name,omitempty"`
	ConnectionDuration *int64  `json:"connection_duration,omitempty" validate:"omitempty,min=0"`

	// Metadata
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	SourceIP  *string         `json:"source_ip,omitempty" validate:"omitempty,ip"`
	UserAgent *string         `json:"user_agent,omitempty"`
	SessionID *string         `json:"session_id,omitempty"`

	// Timestamps
	RecordedAt *time.Time `json:"recorded_at,omitempty"`
}

// SensorLogsFilter represents filter parameters for listing sensor logs
type SensorLogsFilter struct {
	AssetSensorID    *uuid.UUID `json:"asset_sensor_id,omitempty"`
	LogType          *string    `json:"log_type,omitempty"`
	LogLevel         *string    `json:"log_level,omitempty"`
	Component        *string    `json:"component,omitempty"`
	EventType        *string    `json:"event_type,omitempty"`
	ErrorCode        *string    `json:"error_code,omitempty"`
	ConnectionType   *string    `json:"connection_type,omitempty"`
	ConnectionStatus *string    `json:"connection_status,omitempty"`
	SearchMessage    *string    `json:"search_message,omitempty"`
	SeverityLevel    *int       `json:"severity_level,omitempty"`
	HasMetadata      *bool      `json:"has_metadata,omitempty"`
	HasErrors        *bool      `json:"has_errors,omitempty"`
	RecordedAfter    *time.Time `json:"recorded_after,omitempty"`
	RecordedBefore   *time.Time `json:"recorded_before,omitempty"`
	CreatedAfter     *time.Time `json:"created_after,omitempty"`
	CreatedBefore    *time.Time `json:"created_before,omitempty"`
	common.QueryParams
}

// LogAnalyticsRequest represents request for log analytics
type LogAnalyticsRequest struct {
	AssetSensorID  *uuid.UUID `json:"asset_sensor_id,omitempty"`
	LogType        *string    `json:"log_type,omitempty"`
	LogLevel       *string    `json:"log_level,omitempty"`
	StartDate      time.Time  `json:"start_date" binding:"required"`
	EndDate        time.Time  `json:"end_date" binding:"required"`
	GroupBy        string     `json:"group_by" validate:"oneof=hour day week month log_type log_level component"`
	IncludeDetails bool       `json:"include_details"`
}

// SensorLogsResponse represents the response structure for sensor log operations
type SensorLogsResponse struct {
	Data    SensorLogsDTO `json:"data"`
	Message string        `json:"message"`
}

// SensorLogsListResponse represents the response structure for listing sensor logs with pagination
type SensorLogsListResponse struct {
	Data       []SensorLogsDTO           `json:"data"`
	Pagination common.PaginationResponse `json:"pagination"`
	Message    string                    `json:"message"`
}

// LogStatisticsResponse provides aggregated log statistics
type LogStatisticsResponse struct {
	TotalLogs       int64             `json:"total_logs"`
	LogsByType      map[string]int64  `json:"logs_by_type"`
	LogsByLevel     map[string]int64  `json:"logs_by_level"`
	LogsByComponent map[string]int64  `json:"logs_by_component"`
	ErrorRate       float64           `json:"error_rate"`
	RecentErrors    []SensorLogsDTO   `json:"recent_errors"`
	TimeRange       LogTimeRangeStats `json:"time_range"`
}

// LogTimeRangeStats provides time-based statistics
type LogTimeRangeStats struct {
	StartDate    time.Time        `json:"start_date"`
	EndDate      time.Time        `json:"end_date"`
	DurationDays int              `json:"duration_days"`
	LogsPerDay   []TimeSeriesData `json:"logs_per_day"`
}

// TimeSeriesData represents time-series data point
type TimeSeriesData struct {
	Date  time.Time `json:"date"`
	Count int64     `json:"count"`
	Value float64   `json:"value,omitempty"`
}

// ConnectionHistoryResponse provides connection history analysis
type ConnectionHistoryResponse struct {
	AssetSensorID       uuid.UUID         `json:"asset_sensor_id"`
	TotalConnections    int64             `json:"total_connections"`
	TotalDisconnections int64             `json:"total_disconnections"`
	AverageUptime       float64           `json:"average_uptime_hours"`
	UptimePercentage    float64           `json:"uptime_percentage"`
	ConnectionsByType   map[string]int64  `json:"connections_by_type"`
	RecentActivity      []SensorLogsDTO   `json:"recent_activity"`
	TimeRange           LogTimeRangeStats `json:"time_range"`
}

// LogAnalyticsResponse provides comprehensive log analytics
type LogAnalyticsResponse struct {
	Statistics        LogStatisticsResponse     `json:"statistics"`
	ConnectionHistory ConnectionHistoryResponse `json:"connection_history"`
	TrendAnalysis     []TimeSeriesData          `json:"trend_analysis"`
	Anomalies         []LogAnomalyData          `json:"anomalies"`
}

// LogAnomalyData represents detected anomalies in logs
type LogAnomalyData struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Count       int64     `json:"count"`
}

// ToEntity converts CreateSensorLogsRequest to entity.SensorLogs
func (r *CreateSensorLogsRequest) ToEntity(tenantID *uuid.UUID) *entity.SensorLogs {
	now := time.Now()
	recordedAt := now
	if r.RecordedAt != nil {
		recordedAt = *r.RecordedAt
	}

	return &entity.SensorLogs{
		ID:                 uuid.New(),
		TenantID:           tenantID,
		AssetSensorID:      r.AssetSensorID,
		LogType:            r.LogType,
		LogLevel:           r.LogLevel,
		Message:            r.Message,
		Component:          r.Component,
		EventType:          r.EventType,
		ErrorCode:          r.ErrorCode,
		ConnectionType:     r.ConnectionType,
		ConnectionStatus:   r.ConnectionStatus,
		IPAddress:          r.IPAddress,
		MACAddress:         r.MACAddress,
		NetworkName:        r.NetworkName,
		ConnectionDuration: r.ConnectionDuration,
		Metadata:           r.Metadata,
		SourceIP:           r.SourceIP,
		UserAgent:          r.UserAgent,
		SessionID:          r.SessionID,
		RecordedAt:         recordedAt,
		CreatedAt:          now,
	}
}

// ToEntity converts UpdateSensorLogsRequest to entity.SensorLogs for partial updates
func (r *UpdateSensorLogsRequest) ToEntity(id uuid.UUID) *entity.SensorLogs {
	now := time.Now()
	recordedAt := now
	if r.RecordedAt != nil {
		recordedAt = *r.RecordedAt
	}

	log := &entity.SensorLogs{
		ID:         id,
		RecordedAt: recordedAt,
		UpdatedAt:  &now,
	}

	// Only update fields that are provided
	if r.LogType != nil {
		log.LogType = *r.LogType
	}
	if r.LogLevel != nil {
		log.LogLevel = *r.LogLevel
	}
	if r.Message != nil {
		log.Message = *r.Message
	}
	if r.Component != nil {
		log.Component = r.Component
	}
	if r.EventType != nil {
		log.EventType = r.EventType
	}
	if r.ErrorCode != nil {
		log.ErrorCode = r.ErrorCode
	}
	if r.ConnectionType != nil {
		log.ConnectionType = r.ConnectionType
	}
	if r.ConnectionStatus != nil {
		log.ConnectionStatus = r.ConnectionStatus
	}
	if r.IPAddress != nil {
		log.IPAddress = r.IPAddress
	}
	if r.MACAddress != nil {
		log.MACAddress = r.MACAddress
	}
	if r.NetworkName != nil {
		log.NetworkName = r.NetworkName
	}
	if r.ConnectionDuration != nil {
		log.ConnectionDuration = r.ConnectionDuration
	}
	if r.Metadata != nil {
		log.Metadata = r.Metadata
	}
	if r.SourceIP != nil {
		log.SourceIP = r.SourceIP
	}
	if r.UserAgent != nil {
		log.UserAgent = r.UserAgent
	}
	if r.SessionID != nil {
		log.SessionID = r.SessionID
	}

	return log
}

// FromEntity converts entity.SensorLogs to SensorLogsDTO
func FromSensorLogsEntity(e *entity.SensorLogs) *SensorLogsDTO {
	if e == nil {
		return nil
	}
	return &SensorLogsDTO{
		ID:                 e.ID,
		TenantID:           e.TenantID,
		AssetSensorID:      e.AssetSensorID,
		LogType:            e.LogType,
		LogLevel:           e.LogLevel,
		Message:            e.Message,
		Component:          e.Component,
		EventType:          e.EventType,
		ErrorCode:          e.ErrorCode,
		ConnectionType:     e.ConnectionType,
		ConnectionStatus:   e.ConnectionStatus,
		IPAddress:          e.IPAddress,
		MACAddress:         e.MACAddress,
		NetworkName:        e.NetworkName,
		ConnectionDuration: e.ConnectionDuration,
		Metadata:           e.Metadata,
		SourceIP:           e.SourceIP,
		UserAgent:          e.UserAgent,
		SessionID:          e.SessionID,
		RecordedAt:         e.RecordedAt,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}

// FromEntityList converts a slice of entity.SensorLogs to SensorLogsDTO slice
func FromSensorLogsEntityList(entities []*entity.SensorLogs) []SensorLogsDTO {
	dtos := make([]SensorLogsDTO, len(entities))
	for i, entity := range entities {
		if dto := FromSensorLogsEntity(entity); dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}
