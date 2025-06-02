package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateAssetAlertRequest represents the request to create a new asset alert
type CreateAssetAlertRequest struct {
	AssetID       uuid.UUID `json:"asset_id" validate:"required"`
	AssetSensorID uuid.UUID `json:"asset_sensor_id" validate:"required"`
	ThresholdID   uuid.UUID `json:"threshold_id" validate:"required"`
	Severity      string    `json:"severity" validate:"required,oneof=warning critical"`
	Message       string    `json:"message,omitempty"`
}

// UpdateAssetAlertRequest represents the request to update an existing asset alert
type UpdateAssetAlertRequest struct {
	Severity *string `json:"severity,omitempty" validate:"omitempty,oneof=warning critical"`
	// Note: Alert time should not be updatable
	// Resolved time is handled separately via ResolveAlert method
}

// AssetAlertResponse represents the response structure for asset alert operations
type AssetAlertResponse struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	AssetID       uuid.UUID  `json:"asset_id"`
	AssetSensorID uuid.UUID  `json:"asset_sensor_id"`
	ThresholdID   uuid.UUID  `json:"threshold_id"`
	AlertTime     time.Time  `json:"alert_time"`
	ResolvedTime  *time.Time `json:"resolved_time,omitempty"`
	Severity      string     `json:"severity"`
	IsResolved    bool       `json:"is_resolved"`

	// Additional fields for enhanced response
	AssetName       string `json:"asset_name,omitempty"`
	AssetSensorName string `json:"asset_sensor_name,omitempty"`
	ThresholdName   string `json:"threshold_name,omitempty"`
	SensorTypeName  string `json:"sensor_type_name,omitempty"`
	LocationName    string `json:"location_name,omitempty"`
}

// AssetAlertListResponse represents the paginated response for alert list
type AssetAlertListResponse struct {
	Alerts     []AssetAlertResponse `json:"alerts"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	Total      int64                `json:"total"`
	TotalPages int                  `json:"total_pages"`
}

// AssetAlertFilterRequest represents filters for alert listing
type AssetAlertFilterRequest struct {
	AssetID       *uuid.UUID `json:"asset_id,omitempty" form:"asset_id"`
	AssetSensorID *uuid.UUID `json:"asset_sensor_id,omitempty" form:"asset_sensor_id"`
	ThresholdID   *uuid.UUID `json:"threshold_id,omitempty" form:"threshold_id"`
	Severity      *string    `json:"severity,omitempty" form:"severity" validate:"omitempty,oneof=warning critical"`
	IsResolved    *bool      `json:"is_resolved,omitempty" form:"is_resolved"`
	StartTime     *time.Time `json:"start_time,omitempty" form:"start_time"`
	EndTime       *time.Time `json:"end_time,omitempty" form:"end_time"`
	Page          int        `json:"page" form:"page" validate:"min=1"`
	Limit         int        `json:"limit" form:"limit" validate:"min=1,max=100"`
}

// ResolveAlertRequest represents the request to resolve an alert
type ResolveAlertRequest struct {
	AlertID uuid.UUID `json:"alert_id" validate:"required"`
}

// BulkResolveAlertsRequest represents the request to resolve multiple alerts
type BulkResolveAlertsRequest struct {
	AlertIDs []uuid.UUID `json:"alert_ids" validate:"required,min=1"`
}

// BulkResolveAlertsResponse represents the response for bulk resolve operation
type BulkResolveAlertsResponse struct {
	ResolvedCount int         `json:"resolved_count"`
	FailedCount   int         `json:"failed_count"`
	FailedAlerts  []uuid.UUID `json:"failed_alerts,omitempty"`
	Message       string      `json:"message"`
}

// AlertStatisticsRequest represents the request for alert statistics
type AlertStatisticsRequest struct {
	AssetID   *uuid.UUID `json:"asset_id,omitempty" form:"asset_id"`
	StartTime *time.Time `json:"start_time,omitempty" form:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty" form:"end_time"`
}

// AlertStatistics represents statistics for alerts
type AlertStatistics struct {
	TotalAlerts      int                        `json:"total_alerts"`
	UnresolvedAlerts int                        `json:"unresolved_alerts"`
	ResolvedAlerts   int                        `json:"resolved_alerts"`
	CriticalAlerts   int                        `json:"critical_alerts"`
	WarningAlerts    int                        `json:"warning_alerts"`
	AlertsBySeverity map[string]int             `json:"alerts_by_severity"`
	AlertsByAsset    map[string]int             `json:"alerts_by_asset,omitempty"`
	RecentAlerts     []AssetAlertResponse       `json:"recent_alerts,omitempty"`
	TrendData        *AlertTrendData            `json:"trend_data,omitempty"`
}

// AlertTrendData represents trend data for alerts over time
type AlertTrendData struct {
	Period    string             `json:"period"` // daily, weekly, monthly
	DataPoints []AlertTrendPoint  `json:"data_points"`
}

// AlertTrendPoint represents a single point in alert trend data
type AlertTrendPoint struct {
	Date         string `json:"date"`
	AlertCount   int    `json:"alert_count"`
	CriticalCount int   `json:"critical_count"`
	WarningCount  int   `json:"warning_count"`
}

// AlertSummaryResponse represents a summary of alerts for dashboard
type AlertSummaryResponse struct {
	TotalAlerts       int                    `json:"total_alerts"`
	UnresolvedAlerts  int                    `json:"unresolved_alerts"`
	CriticalAlerts    int                    `json:"critical_alerts"`
	WarningAlerts     int                    `json:"warning_alerts"`
	RecentAlerts      []AssetAlertResponse   `json:"recent_alerts"`
	TopAffectedAssets []AlertAssetSummary    `json:"top_affected_assets"`
}

// AlertAssetSummary represents alert summary for a specific asset
type AlertAssetSummary struct {
	AssetID     uuid.UUID `json:"asset_id"`
	AssetName   string    `json:"asset_name"`
	AlertCount  int       `json:"alert_count"`
	LastAlert   time.Time `json:"last_alert"`
}

// GetAssetAlertsRequest represents the request for getting asset alerts with filters
type GetAssetAlertsRequest struct {
	Page       int        `json:"page" form:"page" validate:"min=1"`
	Limit      int        `json:"limit" form:"limit" validate:"min=1,max=100"`
	SortBy     string     `json:"sort_by,omitempty" form:"sort_by"`
	SortOrder  string     `json:"sort_order,omitempty" form:"sort_order" validate:"omitempty,oneof=asc desc"`
	AssetID    *uuid.UUID `json:"asset_id,omitempty" form:"asset_id"`
	SensorID   *uuid.UUID `json:"sensor_id,omitempty" form:"sensor_id"`
	AlertLevel *string    `json:"alert_level,omitempty" form:"alert_level"`
	Status     *string    `json:"status,omitempty" form:"status"`
	IsResolved *bool      `json:"is_resolved,omitempty" form:"is_resolved"`
	StartDate  *string    `json:"start_date,omitempty" form:"start_date"`
	EndDate    *string    `json:"end_date,omitempty" form:"end_date"`
}

// ResolveAssetAlertRequest represents the request to resolve an asset alert
type ResolveAssetAlertRequest struct {
	ResolveReason *string `json:"resolve_reason,omitempty"`
	ResolvedBy    *string `json:"resolved_by,omitempty"`
}

// BulkResolveAssetAlertsRequest represents the request to resolve multiple asset alerts
type BulkResolveAssetAlertsRequest struct {
	AlertIDs      []uuid.UUID `json:"alert_ids" validate:"required,min=1"`
	ResolveReason *string     `json:"resolve_reason,omitempty"`
	ResolvedBy    *string     `json:"resolved_by,omitempty"`
}

// GetAlertStatisticsRequest represents the request for alert statistics
type GetAlertStatisticsRequest struct {
	AssetID   *uuid.UUID `json:"asset_id,omitempty" form:"asset_id"`
	StartTime *time.Time `json:"start_time,omitempty" form:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty" form:"end_time"`
	Period    *string    `json:"period,omitempty" form:"period" validate:"omitempty,oneof=daily weekly monthly"`
}
