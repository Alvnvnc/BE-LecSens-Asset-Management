package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CreateAssetSensorRequest represents the request to create a new asset sensor
type CreateAssetSensorRequest struct {
	AssetID       uuid.UUID       `json:"asset_id" binding:"required" validate:"required"`
	SensorTypeID  uuid.UUID       `json:"sensor_type_id" binding:"required" validate:"required"`
	Name          string          `json:"name" binding:"required" validate:"required"`
	Status        string          `json:"status" binding:"required" validate:"required"`
	Configuration json.RawMessage `json:"configuration,omitempty"`
}

// UpdateAssetSensorRequest represents the request to update an existing asset sensor
type UpdateAssetSensorRequest struct {
	Name          *string          `json:"name,omitempty"`
	Status        *string          `json:"status,omitempty"`
	Configuration *json.RawMessage `json:"configuration,omitempty"`
}

// AssetSensorResponse represents the response structure for asset sensor operations
type AssetSensorResponse struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	AssetID           uuid.UUID       `json:"asset_id"`
	SensorTypeID      uuid.UUID       `json:"sensor_type_id"`
	Name              string          `json:"name"`
	Status            string          `json:"status"`
	Configuration     json.RawMessage `json:"configuration,omitempty"`
	LastReadingValue  *float64        `json:"last_reading_value,omitempty"`
	LastReadingTime   *time.Time      `json:"last_reading_time,omitempty"`
	LastReadingValues json.RawMessage `json:"last_reading_values,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         *time.Time      `json:"updated_at,omitempty"`
}

// AssetSensorListResponse represents the response for listing asset sensors with pagination
type AssetSensorListResponse struct {
	Sensors    []AssetSensorResponse `json:"sensors"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	Total      int64                 `json:"total"`
	TotalPages int                   `json:"total_pages"`
}

// AssetSensorDetailedResponse represents the detailed response structure for asset sensor operations
// This includes sensor type and measurement types information
type AssetSensorDetailedResponse struct {
	ID                uuid.UUID             `json:"id"`
	TenantID          uuid.UUID             `json:"tenant_id"`
	AssetID           uuid.UUID             `json:"asset_id"`
	SensorTypeID      uuid.UUID             `json:"sensor_type_id"`
	Name              string                `json:"name"`
	Status            string                `json:"status"`
	Configuration     json.RawMessage       `json:"configuration,omitempty"`
	LastReadingValue  *float64              `json:"last_reading_value,omitempty"`
	LastReadingTime   *time.Time            `json:"last_reading_time,omitempty"`
	LastReadingValues json.RawMessage       `json:"last_reading_values,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         *time.Time            `json:"updated_at,omitempty"`
	SensorType        SensorTypeInfo        `json:"sensor_type"`
	MeasurementTypes  []MeasurementTypeInfo `json:"measurement_types"`
}

// SensorTypeInfo represents sensor type information in responses
type SensorTypeInfo struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Manufacturer string    `json:"manufacturer"`
	Model        string    `json:"model"`
	Version      string    `json:"version"`
	IsActive     bool      `json:"is_active"`
}

// MeasurementTypeInfo represents measurement type information in responses
type MeasurementTypeInfo struct {
	ID               uuid.UUID              `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	PropertiesSchema json.RawMessage        `json:"properties_schema"`
	UIConfiguration  json.RawMessage        `json:"ui_configuration"`
	Version          string                 `json:"version"`
	IsActive         bool                   `json:"is_active"`
	Fields           []MeasurementFieldInfo `json:"fields"`
}

// MeasurementFieldInfo represents measurement field information in responses
type MeasurementFieldInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	Description *string   `json:"description"`
	DataType    string    `json:"data_type"`
	Required    bool      `json:"required"`
	Unit        *string   `json:"unit"`
	Min         *float64  `json:"min"`
	Max         *float64  `json:"max"`
}

// UpdateSensorReadingRequest represents the request to update sensor readings
type UpdateSensorReadingRequest struct {
	Value    float64                `json:"value" binding:"required"`
	Readings map[string]interface{} `json:"readings,omitempty"`
}
