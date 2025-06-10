package entity

import (
	"time"

	"github.com/google/uuid"
)

// SensorStatus represents real-time sensor status including battery and signal strength
type SensorStatus struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      *uuid.UUID `json:"tenant_id,omitempty"`
	AssetSensorID uuid.UUID  `json:"asset_sensor_id"`

	// Battery Information
	BatteryLevel         *float64   `json:"battery_level,omitempty"`   // Percentage (0-100)
	BatteryVoltage       *float64   `json:"battery_voltage,omitempty"` // Voltage reading
	BatteryStatus        *string    `json:"battery_status,omitempty"`  // "good", "low", "critical", "unknown", "charging"
	BatteryLastCharged   *time.Time `json:"battery_last_charged,omitempty"`
	BatteryEstimatedLife *int       `json:"battery_estimated_life,omitempty"` // Days remaining
	BatteryType          *string    `json:"battery_type,omitempty"`           // "lithium", "alkaline", "rechargeable", "solar"

	// Signal Strength Information
	SignalType      *string  `json:"signal_type,omitempty"`      // "wifi", "cellular", "lora", "zigbee", etc.
	SignalRSSI      *int     `json:"signal_rssi,omitempty"`      // Received Signal Strength Indicator (dBm)
	SignalSNR       *float64 `json:"signal_snr,omitempty"`       // Signal-to-Noise Ratio (dB)
	SignalQuality   *int     `json:"signal_quality,omitempty"`   // Signal quality percentage (0-100)
	SignalFrequency *float64 `json:"signal_frequency,omitempty"` // Frequency in MHz
	SignalChannel   *int     `json:"signal_channel,omitempty"`   // Channel number
	SignalStatus    *string  `json:"signal_status,omitempty"`    // "excellent", "good", "fair", "poor", "no_signal"

	// Connection Information
	ConnectionType     *string    `json:"connection_type,omitempty"` // Current connection type
	ConnectionStatus   string     `json:"connection_status"`         // "online", "offline", "connecting", "error"
	LastConnectedAt    *time.Time `json:"last_connected_at,omitempty"`
	LastDisconnectedAt *time.Time `json:"last_disconnected_at,omitempty"`
	CurrentIP          *string    `json:"current_ip,omitempty"`
	CurrentNetwork     *string    `json:"current_network,omitempty"`

	// Additional Status
	Temperature     *float64   `json:"temperature,omitempty"`    // Sensor internal temperature
	Humidity        *float64   `json:"humidity,omitempty"`       // Environmental humidity
	IsOnline        bool       `json:"is_online"`                // Real-time online status
	LastHeartbeat   *time.Time `json:"last_heartbeat,omitempty"` // Last heartbeat signal
	FirmwareVersion *string    `json:"firmware_version,omitempty"`
	ErrorCount      *int       `json:"error_count,omitempty"` // Recent error count
	LastErrorAt     *time.Time `json:"last_error_at,omitempty"`

	// Timestamps
	RecordedAt time.Time  `json:"recorded_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// NewSensorStatus creates a new sensor status record with default values
func NewSensorStatus() *SensorStatus {
	now := time.Now()
	return &SensorStatus{
		ID:               uuid.New(),
		ConnectionStatus: "offline",
		IsOnline:         false,
		RecordedAt:       now,
		CreatedAt:        now,
	}
}

// TableName returns the table name for GORM
func (SensorStatus) TableName() string {
	return "sensor_status"
}

// IsLowBattery checks if battery level is below threshold
func (ss *SensorStatus) IsLowBattery(threshold float64) bool {
	return ss.BatteryLevel != nil && *ss.BatteryLevel < threshold
}

// IsCriticalBattery checks if battery level is critically low
func (ss *SensorStatus) IsCriticalBattery() bool {
	return ss.BatteryLevel != nil && *ss.BatteryLevel < 10.0
}

// GetSignalQuality determines signal quality based on RSSI for WiFi
func (ss *SensorStatus) GetSignalQuality() string {
	if ss.SignalRSSI == nil {
		return "unknown"
	}

	rssi := *ss.SignalRSSI
	switch {
	case rssi >= -30:
		return "excellent"
	case rssi >= -50:
		return "good"
	case rssi >= -70:
		return "fair"
	case rssi >= -90:
		return "poor"
	default:
		return "no_signal"
	}
}

// IsSignalStrong checks if signal strength is above threshold
func (ss *SensorStatus) IsSignalStrong(threshold int) bool {
	if ss.SignalRSSI == nil {
		return false
	}
	return *ss.SignalRSSI >= threshold
}

// GetConnectionDuration calculates current connection duration
func (ss *SensorStatus) GetConnectionDuration() *int64 {
	if ss.LastConnectedAt != nil && ss.ConnectionStatus == "online" {
		duration := int64(time.Since(*ss.LastConnectedAt).Seconds())
		return &duration
	}
	return nil
}

// IsHealthy checks overall sensor health
func (ss *SensorStatus) IsHealthy() bool {
	// Check if sensor is online
	if !ss.IsOnline {
		return false
	}

	// Check battery level
	if ss.BatteryLevel != nil && *ss.BatteryLevel < 20.0 {
		return false
	}

	// Check signal quality
	if ss.SignalRSSI != nil && *ss.SignalRSSI < -90 {
		return false
	}

	// Check recent errors
	if ss.ErrorCount != nil && *ss.ErrorCount > 10 {
		return false
	}

	return true
}
