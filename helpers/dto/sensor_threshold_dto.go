package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateSensorThresholdRequest represents the request to create a new sensor threshold
type CreateSensorThresholdRequest struct {
	AssetSensorID     uuid.UUID              `json:"asset_sensor_id" validate:"required"`
	SensorTypeID      uuid.UUID              `json:"sensor_type_id" validate:"required"`
	MeasurementField  string                 `json:"measurement_field" validate:"required"`
	Name              string                 `json:"name" validate:"required"`
	Description       string                 `json:"description,omitempty"`
	MinValue          float64                `json:"min_value" validate:"required"`
	MaxValue          float64                `json:"max_value" validate:"required"`
	Severity          string                 `json:"severity" validate:"required,oneof=warning critical"`
	AlertMessage      string                 `json:"alert_message,omitempty"`
	NotificationRules map[string]interface{} `json:"notification_rules,omitempty"`
	IsActive          *bool                  `json:"is_active,omitempty"`
}

// UpdateSensorThresholdRequest represents the request to update an existing sensor threshold
type UpdateSensorThresholdRequest struct {
	Name              *string                `json:"name,omitempty"`
	Description       *string                `json:"description,omitempty"`
	MinValue          *float64               `json:"min_value,omitempty"`
	MaxValue          *float64               `json:"max_value,omitempty"`
	Severity          *string                `json:"severity,omitempty" validate:"omitempty,oneof=warning critical"`
	AlertMessage      *string                `json:"alert_message,omitempty"`
	NotificationRules map[string]interface{} `json:"notification_rules,omitempty"`
	IsActive          *bool                  `json:"is_active,omitempty"`
}

// SensorThresholdResponse represents the response structure for sensor threshold
type SensorThresholdResponse struct {
	ID                uuid.UUID              `json:"id"`
	TenantID          uuid.UUID              `json:"tenant_id"`
	AssetSensorID     uuid.UUID              `json:"asset_sensor_id"`
	SensorTypeID      uuid.UUID              `json:"sensor_type_id"`
	MeasurementField  string                 `json:"measurement_field"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description,omitempty"`
	MinValue          float64                `json:"min_value"`
	MaxValue          float64                `json:"max_value"`
	Severity          string                 `json:"severity"`
	AlertMessage      string                 `json:"alert_message,omitempty"`
	NotificationRules map[string]interface{} `json:"notification_rules,omitempty"`
	IsActive          bool                   `json:"is_active"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         *time.Time             `json:"updated_at,omitempty"`

	// Additional fields for enhanced response
	AssetSensorName string `json:"asset_sensor_name,omitempty"`
	SensorTypeName  string `json:"sensor_type_name,omitempty"`
	AssetName       string `json:"asset_name,omitempty"`
}

// SensorThresholdListResponse represents the paginated response for threshold list
type SensorThresholdListResponse struct {
	Thresholds []SensorThresholdResponse `json:"thresholds"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	Total      int64                     `json:"total"`
	TotalPages int                       `json:"total_pages"`
}

// SensorThresholdFilterRequest represents filters for threshold listing
type SensorThresholdFilterRequest struct {
	AssetSensorID *uuid.UUID `json:"asset_sensor_id,omitempty" form:"asset_sensor_id"`
	SensorTypeID  *uuid.UUID `json:"sensor_type_id,omitempty" form:"sensor_type_id"`
	Severity      *string    `json:"severity,omitempty" form:"severity"`
	IsActive      *bool      `json:"is_active,omitempty" form:"is_active"`
	Page          int        `json:"page" form:"page" validate:"min=1"`
	Limit         int        `json:"limit" form:"limit" validate:"min=1,max=100"`
}

// ThresholdCheckResult represents the result of threshold checking
type ThresholdCheckResult struct {
	ThresholdID   uuid.UUID `json:"threshold_id"`
	ThresholdName string    `json:"threshold_name"`
	IsBreached    bool      `json:"is_breached"`
	CurrentValue  float64   `json:"current_value"`
	MinValue      float64   `json:"min_value"`
	MaxValue      float64   `json:"max_value"`
	Severity      string    `json:"severity"`
	Message       string    `json:"message,omitempty"`
}

// BulkThresholdCheckResponse represents the response for bulk threshold checking
type BulkThresholdCheckResponse struct {
	AssetSensorID uuid.UUID              `json:"asset_sensor_id"`
	Results       []ThresholdCheckResult `json:"results"`
	AlertsCreated int                    `json:"alerts_created"`
}

// NotificationRule represents a notification rule structure
type NotificationRule struct {
	Type        string                 `json:"type"` // email, sms, webhook, etc.
	Enabled     bool                   `json:"enabled"`
	Recipients  []string               `json:"recipients,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	CooldownMin int                    `json:"cooldown_minutes,omitempty"` // Minimum minutes between notifications
}

// ThresholdStatistics represents statistics for thresholds
type ThresholdStatistics struct {
	TotalThresholds    int `json:"total_thresholds"`
	ActiveThresholds   int `json:"active_thresholds"`
	InactiveThresholds int `json:"inactive_thresholds"`
	WarningThresholds  int `json:"warning_thresholds"`
	CriticalThresholds int `json:"critical_thresholds"`
}

// GetSensorThresholdsRequest represents the request for getting sensor thresholds with filters
type GetSensorThresholdsRequest struct {
	Page             int        `json:"page" form:"page" validate:"min=1"`
	Limit            int        `json:"limit" form:"limit" validate:"min=1,max=100"`
	SortBy           string     `json:"sort_by,omitempty" form:"sort_by"`
	SortOrder        string     `json:"sort_order,omitempty" form:"sort_order" validate:"omitempty,oneof=asc desc"`
	AssetSensorID    *uuid.UUID `json:"asset_sensor_id,omitempty" form:"asset_sensor_id"`
	SensorTypeID     *uuid.UUID `json:"sensor_type_id,omitempty" form:"sensor_type_id"`
	Severity         *string    `json:"severity,omitempty" form:"severity"`
	IsActive         *bool      `json:"is_active,omitempty" form:"is_active"`
	MeasurementField *string    `json:"measurement_field,omitempty" form:"measurement_field"`
}

// CheckThresholdsRequest represents the request for checking thresholds against a value
type CheckThresholdsRequest struct {
	AssetSensorID    uuid.UUID `json:"asset_sensor_id" validate:"required"`
	MeasurementField string    `json:"measurement_field" validate:"required"`
	Value            float64   `json:"value" validate:"required"`
}

// GetThresholdStatisticsRequest represents the request for threshold statistics
type GetThresholdStatisticsRequest struct {
	AssetSensorID *uuid.UUID `json:"asset_sensor_id,omitempty" form:"asset_sensor_id"`
	SensorTypeID  *uuid.UUID `json:"sensor_type_id,omitempty" form:"sensor_type_id"`
	StartDate     *string    `json:"start_date,omitempty" form:"start_date"`
	EndDate       *string    `json:"end_date,omitempty" form:"end_date"`
}
