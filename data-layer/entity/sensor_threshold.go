package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ThresholdSeverity represents the criticality of a threshold breach
type ThresholdSeverity string

const (
	ThresholdSeverityWarning  ThresholdSeverity = "warning"
	ThresholdSeverityCritical ThresholdSeverity = "critical"
)

// ThresholdStatus represents the current status of a measurement against thresholds
type ThresholdStatus string

const (
	ThresholdStatusNormal   ThresholdStatus = "normal"
	ThresholdStatusWarning  ThresholdStatus = "warning"
	ThresholdStatusCritical ThresholdStatus = "critical"
)

// SensorThreshold defines alert/warning ranges for sensor measurements
// This defines when to trigger alerts, separate from sensor capability ranges
type SensorThreshold struct {
	ID                   uuid.UUID         `json:"id"`
	TenantID             uuid.UUID         `json:"tenant_id"`
	AssetSensorID        uuid.UUID         `json:"asset_sensor_id"`        // References the asset sensor
	MeasurementTypeID    uuid.UUID         `json:"measurement_type_id"`    // References the measurement type
	MeasurementFieldName string            `json:"measurement_field_name"` // Field name from measurement fields
	MinValue             *float64          `json:"min_value,omitempty"`    // Alert if value < min_value
	MaxValue             *float64          `json:"max_value,omitempty"`    // Alert if value > max_value
	Severity             ThresholdSeverity `json:"severity"`
	IsActive             bool              `json:"is_active"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            *time.Time        `json:"updated_at,omitempty"`
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

// CheckValue determines the status of a value against the thresholds
func (t *SensorThreshold) CheckValue(value float64) ThresholdStatus {
	// If no limits are set, everything is normal
	if t.MinValue == nil && t.MaxValue == nil {
		return ThresholdStatusNormal
	}

	// Check if value is within normal range
	if t.MinValue != nil && t.MaxValue != nil {
		if value >= *t.MinValue && value <= *t.MaxValue {
			return ThresholdStatusNormal
		}
	} else if t.MinValue != nil && value >= *t.MinValue {
		return ThresholdStatusNormal
	} else if t.MaxValue != nil && value <= *t.MaxValue {
		return ThresholdStatusNormal
	}

	// If we get here, the value is outside the normal range
	if t.Severity == ThresholdSeverityCritical {
		return ThresholdStatusCritical
	}
	return ThresholdStatusWarning
}

// IsBreached determines if a value breaches the threshold (triggers alert)
func (t *SensorThreshold) IsBreached(value float64) bool {
	return t.CheckValue(value) != ThresholdStatusNormal
}

// SetThresholds sets the min/max threshold values
func (t *SensorThreshold) SetThresholds(minValue, maxValue *float64) error {
	if minValue != nil && maxValue != nil && *minValue >= *maxValue {
		return fmt.Errorf("minimum threshold must be less than maximum threshold")
	}

	t.MinValue = minValue
	t.MaxValue = maxValue
	return nil
}

// GenerateAlertMessage creates an alert message based on the values
func (t *SensorThreshold) GenerateAlertMessage(value float64, alertType string) string {
	// Generate default message
	switch alertType {
	case "min_breach":
		if t.MinValue != nil {
			return fmt.Sprintf("%s value %.2f is below minimum threshold %.2f",
				t.MeasurementFieldName, value, *t.MinValue)
		}
	case "max_breach":
		if t.MaxValue != nil {
			return fmt.Sprintf("%s value %.2f exceeds maximum threshold %.2f",
				t.MeasurementFieldName, value, *t.MaxValue)
		}
	}

	return fmt.Sprintf("%s value %.2f breached threshold",
		t.MeasurementFieldName, value)
}

// ValidateConfiguration validates the threshold configuration
func (t *SensorThreshold) ValidateConfiguration() error {
	if t.MeasurementFieldName == "" {
		return fmt.Errorf("measurement field name is required")
	}

	// Validate threshold values if both are set
	if t.MinValue != nil && t.MaxValue != nil {
		if *t.MinValue >= *t.MaxValue {
			return fmt.Errorf("minimum threshold must be less than maximum threshold")
		}
	}

	return nil
}

// GetThresholdInfo returns a summary of the threshold configuration
func (t *SensorThreshold) GetThresholdInfo() map[string]interface{} {
	info := map[string]interface{}{
		"id":                t.ID,
		"measurement_field": t.MeasurementFieldName,
		"severity":          t.Severity,
		"is_active":         t.IsActive,
	}

	if t.MinValue != nil {
		info["min_value"] = *t.MinValue
	}
	if t.MaxValue != nil {
		info["max_value"] = *t.MaxValue
	}

	return info
}

// CreateAlertIfBreached checks if a value breaches the threshold and creates an alert
func (t *SensorThreshold) CreateAlertIfBreached(
	value float64,
	tenantID, assetID, assetSensorID uuid.UUID,
) *AssetAlert {
	status := t.CheckValue(value)
	if status == ThresholdStatusNormal {
		return nil // No breach, no alert needed
	}

	alert := CreateAlertFromThresholdBreach(tenantID, assetID, assetSensorID, t, value)
	alert.Severity = ThresholdSeverity(status) // Set severity based on status
	return alert
}
