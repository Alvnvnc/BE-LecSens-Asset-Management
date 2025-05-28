package dto

import (
	"time"

	"github.com/google/uuid"
)

// SensorMeasurementFieldDTO represents a sensor measurement field data transfer object
type SensorMeasurementFieldDTO struct {
	ID                      uuid.UUID `json:"id"`
	SensorMeasurementTypeID uuid.UUID `json:"sensor_measurement_type_id"`
	Name                    string    `json:"name"`
	Label                   string    `json:"label"`
	Description             *string   `json:"description,omitempty"`
	DataType                string    `json:"data_type"`
	Required                bool      `json:"required"`
	Unit                    *string   `json:"unit,omitempty"`
	Min                     *float64  `json:"min,omitempty"`
	Max                     *float64  `json:"max,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// CreateSensorMeasurementFieldRequest represents the request structure for creating a sensor measurement field
type CreateSensorMeasurementFieldRequest struct {
	SensorMeasurementTypeID uuid.UUID `json:"sensor_measurement_type_id" binding:"required"`
	Name                    string    `json:"name" binding:"required"`
	Label                   string    `json:"label" binding:"required"`
	Description             *string   `json:"description,omitempty"`
	DataType                string    `json:"data_type" binding:"required"`
	Required                bool      `json:"required"`
	Unit                    *string   `json:"unit,omitempty"`
	Min                     *float64  `json:"min,omitempty"`
	Max                     *float64  `json:"max,omitempty"`
}

// UpdateSensorMeasurementFieldRequest represents the request structure for updating a sensor measurement field
type UpdateSensorMeasurementFieldRequest struct {
	Name        *string  `json:"name,omitempty"`
	Label       *string  `json:"label,omitempty"`
	Description *string  `json:"description,omitempty"`
	DataType    *string  `json:"data_type,omitempty"`
	Required    *bool    `json:"required,omitempty"`
	Unit        *string  `json:"unit,omitempty"`
	Min         *float64 `json:"min,omitempty"`
	Max         *float64 `json:"max,omitempty"`
}

// SensorMeasurementFieldResponse represents the response structure for sensor measurement field operations
type SensorMeasurementFieldResponse struct {
	Data    SensorMeasurementFieldDTO `json:"data"`
	Message string                    `json:"message"`
}

// SensorMeasurementFieldListResponse represents the response structure for listing sensor measurement fields
type SensorMeasurementFieldListResponse struct {
	Data    []SensorMeasurementFieldDTO `json:"data"`
	Total   int                         `json:"total"`
	Page    int                         `json:"page"`
	Limit   int                         `json:"limit"`
	Message string                      `json:"message"`
}
