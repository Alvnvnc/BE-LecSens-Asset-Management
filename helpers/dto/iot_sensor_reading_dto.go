package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CreateIoTSensorReadingRequest represents the request to create a new IoT sensor reading
type CreateIoTSensorReadingRequest struct {
	TenantID        *uuid.UUID      `json:"tenant_id,omitempty"`
	AssetSensorID   *uuid.UUID      `json:"asset_sensor_id,omitempty"`
	SensorTypeID    uuid.UUID       `json:"sensor_type_id" binding:"required" validate:"required"`
	MacAddress      string          `json:"mac_address" binding:"required" validate:"required"`
	Location        string          `json:"location" binding:"required" validate:"required"`
	MeasurementData json.RawMessage `json:"measurement_data" binding:"required" validate:"required"`
	ReadingTime     *time.Time      `json:"reading_time,omitempty"`

	// Legacy fields (for backward compatibility)
	DataX   json.RawMessage `json:"data_x,omitempty"`
	DataY   json.RawMessage `json:"data_y,omitempty"`
	PeakX   json.RawMessage `json:"peak_x,omitempty"`
	PeakY   json.RawMessage `json:"peak_y,omitempty"`
	PPM     *float64        `json:"ppm,omitempty"`
	Label   *string         `json:"label,omitempty"`
	RawData json.RawMessage `json:"raw_data,omitempty"`
}

// UpdateIoTSensorReadingRequest represents the request body for partial IoT sensor reading updates
type UpdateIoTSensorReadingRequest struct {
	TenantID        *uuid.UUID      `json:"tenant_id,omitempty"`
	AssetSensorID   *uuid.UUID      `json:"asset_sensor_id,omitempty"`
	SensorTypeID    *uuid.UUID      `json:"sensor_type_id,omitempty"`
	MacAddress      *string         `json:"mac_address,omitempty"`
	Location        *string         `json:"location,omitempty"`
	MeasurementData json.RawMessage `json:"measurement_data,omitempty"`
	ReadingTime     *time.Time      `json:"reading_time,omitempty"`

	// Legacy fields (for backward compatibility)
	DataX   json.RawMessage `json:"data_x,omitempty"`
	DataY   json.RawMessage `json:"data_y,omitempty"`
	PeakX   json.RawMessage `json:"peak_x,omitempty"`
	PeakY   json.RawMessage `json:"peak_y,omitempty"`
	PPM     *float64        `json:"ppm,omitempty"`
	Label   *string         `json:"label,omitempty"`
	RawData json.RawMessage `json:"raw_data,omitempty"`
}

// IoTSensorReadingResponse represents the response structure for IoT sensor reading operations
type IoTSensorReadingResponse struct {
	ID              uuid.UUID       `json:"id"`
	TenantID        *uuid.UUID      `json:"tenant_id,omitempty"`
	AssetSensorID   *uuid.UUID      `json:"asset_sensor_id,omitempty"`
	SensorTypeID    uuid.UUID       `json:"sensor_type_id"`
	MacAddress      string          `json:"mac_address"`
	Location        string          `json:"location"`
	MeasurementData json.RawMessage `json:"measurement_data"`
	StandardFields  json.RawMessage `json:"standard_fields,omitempty"`
	ReadingTime     string          `json:"reading_time"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       *string         `json:"updated_at,omitempty"`

	// Legacy fields (for backward compatibility)
	DataX   json.RawMessage `json:"data_x,omitempty"`
	DataY   json.RawMessage `json:"data_y,omitempty"`
	PeakX   json.RawMessage `json:"peak_x,omitempty"`
	PeakY   json.RawMessage `json:"peak_y,omitempty"`
	PPM     *float64        `json:"ppm,omitempty"`
	Label   *string         `json:"label,omitempty"`
	RawData json.RawMessage `json:"raw_data,omitempty"`
}

// IoTSensorReadingListResponse represents the response for listing IoT sensor readings with pagination
type IoTSensorReadingListResponse struct {
	Readings   []IoTSensorReadingResponse `json:"readings"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	Total      int64                      `json:"total"`
	TotalPages int                        `json:"total_pages"`
}

// IoTSensorReadingQueryParams represents query parameters for filtering IoT sensor readings
type IoTSensorReadingQueryParams struct {
	TenantID      *uuid.UUID `form:"tenant_id"`
	AssetSensorID *uuid.UUID `form:"asset_sensor_id"`
	SensorTypeID  *uuid.UUID `form:"sensor_type_id"`
	MacAddress    *string    `form:"mac_address"`
	Location      *string    `form:"location"`
	StartTime     *time.Time `form:"start_time" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime       *time.Time `form:"end_time" time_format:"2006-01-02T15:04:05Z07:00"`
	Page          int        `form:"page" binding:"min=1"`
	Limit         int        `form:"limit" binding:"min=1,max=100"`
}
