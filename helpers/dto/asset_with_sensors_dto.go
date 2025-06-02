package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CreateAssetWithSensorsRequest represents the request to create an asset with sensors
type CreateAssetWithSensorsRequest struct {
	Name        string                      `json:"name" binding:"required" validate:"required"`
	AssetTypeID uuid.UUID                   `json:"asset_type_id" binding:"required" validate:"required"`
	LocationID  uuid.UUID                   `json:"location_id" binding:"required" validate:"required"`
	Status      string                      `json:"status,omitempty"`
	Properties  json.RawMessage             `json:"properties,omitempty"`
	SensorTypes []CreateAssetSensorFromType `json:"sensor_types" binding:"required" validate:"required,min=1"`
}

// CreateAssetSensorFromType represents sensor configuration for asset creation
type CreateAssetSensorFromType struct {
	SensorTypeID  uuid.UUID       `json:"sensor_type_id" binding:"required" validate:"required"`
	Status        string          `json:"status,omitempty"`
	Configuration json.RawMessage `json:"configuration,omitempty"`
}

// AssetWithSensorsResponse represents the response for asset creation with sensors
type AssetWithSensorsResponse struct {
	Asset   AssetResponse         `json:"asset"`
	Sensors []AssetSensorResponse `json:"sensors"`
}

// CreateAssetWithSensorsResult represents the internal result structure
type CreateAssetWithSensorsResult struct {
	AssetID   uuid.UUID
	SensorIDs []uuid.UUID
	CreatedAt time.Time
	TenantID  *uuid.UUID
	Errors    []string
}
