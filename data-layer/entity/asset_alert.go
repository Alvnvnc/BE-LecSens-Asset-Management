package entity

import (
	"time"

	"github.com/google/uuid"
)

// AssetAlert represents a notification generated when a sensor reading exceeds thresholds
type AssetAlert struct {
	ID            uuid.UUID         `json:"id"`
	TenantID      uuid.UUID         `json:"tenant_id"`
	AssetID       uuid.UUID         `json:"asset_id"`
	AssetSensorID uuid.UUID         `json:"asset_sensor_id"`
	ThresholdID   uuid.UUID         `json:"threshold_id"`
	AlertTime     time.Time         `json:"alert_time"`
	ResolvedTime  *time.Time        `json:"resolved_time,omitempty"`
	Severity      ThresholdSeverity `json:"severity"`
}
