package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SensorLogs represents comprehensive logging for all sensor activities
type SensorLogs struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      *uuid.UUID `json:"tenant_id,omitempty"`
	AssetSensorID uuid.UUID  `json:"asset_sensor_id"`
	LogType       string     `json:"log_type"`  // "reading", "connection", "battery", "signal", "error", "system"
	LogLevel      string     `json:"log_level"` // "debug", "info", "warning", "error", "critical"
	Message       string     `json:"message"`
	Component     *string    `json:"component,omitempty"`  // "sensor", "communication", "battery", "hardware", "network"
	EventType     *string    `json:"event_type,omitempty"` // "startup", "shutdown", "connected", "disconnected", "reading", "error", "maintenance"
	ErrorCode     *string    `json:"error_code,omitempty"`

	// Connection History Data (when log_type = "connection")
	ConnectionType     *string `json:"connection_type,omitempty"`   // "wifi", "cellular", "lora", "zigbee", etc.
	ConnectionStatus   *string `json:"connection_status,omitempty"` // "connected", "disconnected", "failed"
	IPAddress          *string `json:"ip_address,omitempty"`
	MACAddress         *string `json:"mac_address,omitempty"`
	NetworkName        *string `json:"network_name,omitempty"`
	ConnectionDuration *int64  `json:"connection_duration,omitempty"` // seconds

	// Metadata for flexible additional data
	Metadata  json.RawMessage `json:"metadata,omitempty"` // Additional contextual data
	SourceIP  *string         `json:"source_ip,omitempty"`
	UserAgent *string         `json:"user_agent,omitempty"`
	SessionID *string         `json:"session_id,omitempty"`

	// Timestamps
	RecordedAt time.Time  `json:"recorded_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// NewSensorLogs creates a new sensor log entry with default values
func NewSensorLogs() *SensorLogs {
	now := time.Now()
	return &SensorLogs{
		ID:         uuid.New(),
		LogLevel:   "info",
		LogType:    "system",
		RecordedAt: now,
		CreatedAt:  now,
	}
}

// TableName returns the table name for GORM
func (SensorLogs) TableName() string {
	return "sensor_logs"
}

// IsError checks if the log entry is an error or critical level
func (sl *SensorLogs) IsError() bool {
	return sl.LogLevel == "error" || sl.LogLevel == "critical"
}

// IsConnectionLog checks if this is a connection-related log
func (sl *SensorLogs) IsConnectionLog() bool {
	return sl.LogType == "connection"
}

// IsBatteryLog checks if this is a battery-related log
func (sl *SensorLogs) IsBatteryLog() bool {
	return sl.LogType == "battery"
}

// IsSignalLog checks if this is a signal-related log
func (sl *SensorLogs) IsSignalLog() bool {
	return sl.LogType == "signal"
}

// GetSeverityLevel returns numeric severity level for sorting/filtering
func (sl *SensorLogs) GetSeverityLevel() int {
	switch sl.LogLevel {
	case "debug":
		return 1
	case "info":
		return 2
	case "warning":
		return 3
	case "error":
		return 4
	case "critical":
		return 5
	default:
		return 0
	}
}
