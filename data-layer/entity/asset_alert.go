package entity

import (
	"time"

	"github.com/google/uuid"
)

// AssetAlert represents a notification generated when a sensor reading exceeds thresholds
type AssetAlert struct {
	ID                   uuid.UUID         `json:"id"`
	TenantID             uuid.UUID         `json:"tenant_id"`
	AssetID              uuid.UUID         `json:"asset_id"`
	AssetSensorID        uuid.UUID         `json:"asset_sensor_id"`
	ThresholdID          uuid.UUID         `json:"threshold_id"`
	MeasurementFieldName string            `json:"measurement_field_name"` // Field yang trigger alert
	AlertTime            time.Time         `json:"alert_time"`
	ResolvedTime         *time.Time        `json:"resolved_time,omitempty"`
	Severity             ThresholdSeverity `json:"severity"`
	Status               ThresholdStatus   `json:"status"`              // Current status of the alert
	TriggerValue         float64           `json:"trigger_value"`       // Nilai yang memicu alert
	ThresholdMinValue    *float64          `json:"threshold_min_value"` // Min threshold saat alert
	ThresholdMaxValue    *float64          `json:"threshold_max_value"` // Max threshold saat alert
	AlertMessage         string            `json:"alert_message"`       // Pesan alert
	AlertType            string            `json:"alert_type"`          // "min_breach", "max_breach"
	IsResolved           bool              `json:"is_resolved"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            *time.Time        `json:"updated_at,omitempty"`
}

// NewAssetAlert creates a new asset alert
func NewAssetAlert() *AssetAlert {
	now := time.Now()
	return &AssetAlert{
		ID:         uuid.New(),
		AlertTime:  now,
		IsResolved: false,
		Status:     ThresholdStatusWarning, // Default to warning
		CreatedAt:  now,
	}
}

// CreateAlertFromThresholdBreach creates an alert when a threshold is breached
func CreateAlertFromThresholdBreach(
	tenantID, assetID, assetSensorID uuid.UUID,
	threshold *SensorThreshold,
	triggerValue float64,
) *AssetAlert {
	alert := NewAssetAlert()
	alert.TenantID = tenantID
	alert.AssetID = assetID
	alert.AssetSensorID = assetSensorID
	alert.ThresholdID = threshold.ID
	alert.MeasurementFieldName = threshold.MeasurementFieldName
	alert.Severity = threshold.Severity
	alert.TriggerValue = triggerValue
	alert.ThresholdMinValue = threshold.MinValue
	alert.ThresholdMaxValue = threshold.MaxValue

	// Determine alert type and status
	if threshold.MinValue != nil && triggerValue < *threshold.MinValue {
		alert.AlertType = "min_breach"
		alert.Status = ThresholdStatusWarning
	} else if threshold.MaxValue != nil && triggerValue > *threshold.MaxValue {
		alert.AlertType = "max_breach"
		alert.Status = ThresholdStatusWarning
	}

	// If severity is critical, update status
	if threshold.Severity == ThresholdSeverityCritical {
		alert.Status = ThresholdStatusCritical
	}

	// Generate alert message using threshold template
	alert.AlertMessage = threshold.GenerateAlertMessage(triggerValue, alert.AlertType)

	return alert
}

// Resolve marks the alert as resolved
func (a *AssetAlert) Resolve() {
	now := time.Now()
	a.ResolvedTime = &now
	a.IsResolved = true
	a.Status = ThresholdStatusNormal
	a.UpdatedAt = &now
}

// UpdateStatus updates the alert status based on new sensor reading
func (a *AssetAlert) UpdateStatus(newValue float64, threshold *SensorThreshold) {
	if a.IsResolved {
		return
	}

	status := threshold.CheckValue(newValue)
	if status == ThresholdStatusNormal {
		a.Resolve()
	} else {
		a.Status = ThresholdStatus(status)
		a.TriggerValue = newValue
		now := time.Now()
		a.UpdatedAt = &now
	}
}

// IsActive returns true if the alert is still active (not resolved)
func (a *AssetAlert) IsActive() bool {
	return !a.IsResolved
}

// GetDuration returns how long the alert has been active
func (a *AssetAlert) GetDuration() time.Duration {
	if a.IsResolved && a.ResolvedTime != nil {
		return a.ResolvedTime.Sub(a.AlertTime)
	}
	return time.Since(a.AlertTime)
}

// GetAlertInfo returns comprehensive alert information
func (a *AssetAlert) GetAlertInfo() map[string]interface{} {
	info := map[string]interface{}{
		"id":                a.ID,
		"asset_id":          a.AssetID,
		"asset_sensor_id":   a.AssetSensorID,
		"threshold_id":      a.ThresholdID,
		"measurement_field": a.MeasurementFieldName,
		"alert_time":        a.AlertTime,
		"severity":          a.Severity,
		"trigger_value":     a.TriggerValue,
		"alert_type":        a.AlertType,
		"alert_message":     a.AlertMessage,
		"is_resolved":       a.IsResolved,
		"duration_seconds":  a.GetDuration().Seconds(),
	}

	if a.ThresholdMinValue != nil {
		info["threshold_min_value"] = *a.ThresholdMinValue
	}
	if a.ThresholdMaxValue != nil {
		info["threshold_max_value"] = *a.ThresholdMaxValue
	}
	if a.ResolvedTime != nil {
		info["resolved_time"] = *a.ResolvedTime
	}

	return info
}
