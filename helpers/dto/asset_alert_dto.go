package dto

import (
	"be-lecsens/asset_management/data-layer/entity"
	"time"

	"github.com/google/uuid"
)

// AssetAlertResponse represents the response structure for asset alerts
type AssetAlertResponse struct {
	ID                   uuid.UUID                `json:"id"`
	TenantID             uuid.UUID                `json:"tenant_id"`
	AssetID              uuid.UUID                `json:"asset_id"`
	AssetSensorID        uuid.UUID                `json:"asset_sensor_id"`
	ThresholdID          uuid.UUID                `json:"threshold_id"`
	MeasurementFieldName string                   `json:"measurement_field_name"`
	AlertTime            time.Time                `json:"alert_time"`
	ResolvedTime         *time.Time               `json:"resolved_time,omitempty"`
	Severity             entity.ThresholdSeverity `json:"severity"`
	TriggerValue         float64                  `json:"trigger_value"`
	ThresholdMinValue    *float64                 `json:"threshold_min_value,omitempty"`
	ThresholdMaxValue    *float64                 `json:"threshold_max_value,omitempty"`
	AlertMessage         string                   `json:"alert_message"`
	AlertType            string                   `json:"alert_type"`
	IsResolved           bool                     `json:"is_resolved"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            *time.Time               `json:"updated_at,omitempty"`
}

// AssetAlertListResponse represents the paginated response for listing asset alerts
type AssetAlertListResponse struct {
	Data       []AssetAlertResponse `json:"data"`
	Pagination PaginationInfo       `json:"pagination"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_prev"`
}

// AssetAlertStatisticsResponse represents alert statistics
type AssetAlertStatisticsResponse struct {
	TotalAlerts    int `json:"total_alerts"`
	ActiveAlerts   int `json:"active_alerts"`
	ResolvedAlerts int `json:"resolved_alerts"`
	CriticalAlerts int `json:"critical_alerts"`
	WarningAlerts  int `json:"warning_alerts"`
	Alerts24h      int `json:"alerts_24h"`
	Alerts7d       int `json:"alerts_7d"`
	TotalTenants   int `json:"total_tenants,omitempty"` // Only for global statistics
}

// ResolveMultipleAlertsRequest represents the request for resolving multiple alerts
type ResolveMultipleAlertsRequest struct {
	AlertIDs []uuid.UUID `json:"alert_ids" binding:"required"`
}

// ResolveMultipleAlertsResponse represents the response for resolving multiple alerts
type ResolveMultipleAlertsResponse struct {
	ResolvedCount  int `json:"resolved_count"`
	FailedCount    int `json:"failed_count"`
	TotalRequested int `json:"total_requested"`
}

// DeleteMultipleAlertsRequest represents the request for deleting multiple alerts
type DeleteMultipleAlertsRequest struct {
	AlertIDs []uuid.UUID `json:"alert_ids" binding:"required"`
}

// DeleteMultipleAlertsResponse represents the response for deleting multiple alerts
type DeleteMultipleAlertsResponse struct {
	DeletedCount   int `json:"deleted_count"`
	FailedCount    int `json:"failed_count"`
	TotalRequested int `json:"total_requested"`
}

// AssetAlertFilter represents filter parameters for listing alerts
type AssetAlertFilter struct {
	AssetID       *uuid.UUID                `json:"asset_id,omitempty"`
	AssetSensorID *uuid.UUID                `json:"asset_sensor_id,omitempty"`
	Severity      *entity.ThresholdSeverity `json:"severity,omitempty"`
	IsResolved    *bool                     `json:"is_resolved,omitempty"`
	FromTime      *time.Time                `json:"from_time,omitempty"`
	ToTime        *time.Time                `json:"to_time,omitempty"`
	Page          int                       `json:"page"`
	Limit         int                       `json:"limit"`
}
