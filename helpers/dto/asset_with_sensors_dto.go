package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CreateAssetSensorTypeRequest represents the sensor type configuration in asset creation request
type CreateAssetSensorTypeRequest struct {
	SensorTypeID  uuid.UUID       `json:"sensor_type_id" binding:"required" validate:"required"`
	Status        string          `json:"status,omitempty"` // Default: "active"
	Configuration json.RawMessage `json:"configuration,omitempty"`
}

// CreateAssetWithSensorsRequest represents the request to create an asset with its sensors
type CreateAssetWithSensorsRequest struct {
	Name        string                         `json:"name" binding:"required" validate:"required"`
	AssetTypeID uuid.UUID                      `json:"asset_type_id" binding:"required" validate:"required"`
	LocationID  uuid.UUID                      `json:"location_id" binding:"required" validate:"required"`
	Status      string                         `json:"status,omitempty"` // Default: "active"
	Properties  json.RawMessage                `json:"properties,omitempty"`
	SensorTypes []CreateAssetSensorTypeRequest `json:"sensor_types" binding:"required" validate:"required,dive"`
}

// UpdateAssetWithSensorsRequest represents the request to update an asset with its sensors
type UpdateAssetWithSensorsRequest struct {
	Name        *string                        `json:"name,omitempty"`
	AssetTypeID *uuid.UUID                     `json:"asset_type_id,omitempty"`
	LocationID  *uuid.UUID                     `json:"location_id,omitempty"`
	Status      *string                        `json:"status,omitempty"`
	Properties  json.RawMessage                `json:"properties,omitempty"`
	SensorTypes []CreateAssetSensorTypeRequest `json:"sensor_types,omitempty"`
}

// AssetWithSensorsResponse represents the response structure for asset with sensors operations
type AssetWithSensorsResponse struct {
	Asset   AssetResponse         `json:"asset"`
	Sensors []AssetSensorResponse `json:"sensors"`
}

// AssetWithSensorsListResponse represents the response for listing assets with sensors with pagination
type AssetWithSensorsListResponse struct {
	Assets     []AssetWithSensorsResponse `json:"assets"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	Total      int64                      `json:"total"`
	TotalPages int                        `json:"total_pages"`
}

// CreateAssetWithSensorsResult represents the result of creating a single asset with sensors (used internally)
type CreateAssetWithSensorsResult struct {
	AssetWithSensors *AssetWithSensorsResponse `json:"asset_with_sensors"`
	AssetID          uuid.UUID                 `json:"asset_id"`
	SensorIDs        []uuid.UUID               `json:"sensor_ids"`
	TenantID         *uuid.UUID                `json:"tenant_id,omitempty"`
	CreatedAt        time.Time                 `json:"created_at"`
	Errors           []string                  `json:"errors"`
	Error            error                     `json:"error,omitempty"`
}

// BulkCreateAssetWithSensorsResponse represents the response for bulk asset with sensors creation
type BulkCreateAssetWithSensorsResponse struct {
	Results         []CreateAssetWithSensorsResult `json:"results"`
	SuccessCount    int                            `json:"success_count"`
	ErrorCount      int                            `json:"error_count"`
	Errors          []string                       `json:"errors"`
	TotalRequested  int                            `json:"total_requested"`
	TotalSuccessful int                            `json:"total_successful"`
	TotalFailed     int                            `json:"total_failed"`
}

// BulkCreateError represents an error that occurred during bulk creation
type BulkCreateError struct {
	Index   int                           `json:"index"`
	Request CreateAssetWithSensorsRequest `json:"request"`
	Error   string                        `json:"error"`
}
