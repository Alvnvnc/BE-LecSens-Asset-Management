package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ThresholdSeverity represents the criticality of a threshold breach
type ThresholdSeverity string

const (
	ThresholdSeverityWarning  ThresholdSeverity = "warning"
	ThresholdSeverityCritical ThresholdSeverity = "critical"
)

// SensorThreshold defines allowable ranges for sensor readings
type SensorThreshold struct {
	ID                uuid.UUID         `json:"id"`
	TenantID          uuid.UUID         `json:"tenant_id"`
	AssetSensorID     uuid.UUID         `json:"asset_sensor_id"`
	SensorTypeID      uuid.UUID         `json:"sensor_type_id"`
	MeasurementField  string            `json:"measurement_field"` // Field name in the measurement data
	Name              string            `json:"name"`
	Description       string            `json:"description,omitempty"`
	MinValue          float64           `json:"min_value"`
	MaxValue          float64           `json:"max_value"`
	Severity          ThresholdSeverity `json:"severity"`
	AlertMessage      string            `json:"alert_message,omitempty"`
	NotificationRules json.RawMessage   `json:"notification_rules,omitempty"` // For configurable notifications
	IsActive          bool              `json:"is_active"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         *time.Time        `json:"updated_at,omitempty"`
}

// NewSensorThreshold creates a new threshold with default values
func NewSensorThreshold() *SensorThreshold {
	now := time.Now()
	return &SensorThreshold{
		ID:        uuid.New(),
		Severity:  ThresholdSeverityWarning,
		IsActive:  true,
		CreatedAt: now,
	}
}

// CheckValue determines if a value is within the threshold range
func (t *SensorThreshold) CheckValue(value float64) bool {
	return value >= t.MinValue && value <= t.MaxValue
}

// IsBreached determines if a value breaches the threshold
func (t *SensorThreshold) IsBreached(value float64) bool {
	return value < t.MinValue || value > t.MaxValue
}

// SetNotificationRules sets the notification rules for this threshold
func (t *SensorThreshold) SetNotificationRules(rules map[string]interface{}) error {
	rulesJSON, err := json.Marshal(rules)
	if err != nil {
		return err
	}
	t.NotificationRules = rulesJSON
	return nil
}

// GetNotificationRules parses the notification rules into a map
func (t *SensorThreshold) GetNotificationRules() (map[string]interface{}, error) {
	if t.NotificationRules == nil {
		return map[string]interface{}{}, nil
	}

	var rules map[string]interface{}
	if err := json.Unmarshal(t.NotificationRules, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}
