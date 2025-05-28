package dto

import (
	"time"

	"github.com/google/uuid"
)

// SensorMeasurementTypeDTO represents the data transfer object for sensor measurement types
type SensorMeasurementTypeDTO struct {
	ID               uuid.UUID   `json:"id"`
	SensorTypeID     uuid.UUID   `json:"sensor_type_id"`
	Name             string      `json:"name"`
	Description      *string     `json:"description,omitempty"`
	UnitOfMeasure    *string     `json:"unit_of_measure,omitempty"`
	MinAcceptedValue *float64    `json:"min_accepted_value,omitempty"`
	MaxAcceptedValue *float64    `json:"max_accepted_value,omitempty"`
	PropertiesSchema interface{} `json:"properties_schema,omitempty"`
	UIConfiguration  interface{} `json:"ui_configuration,omitempty"`
	Version          int         `json:"version"`
	IsActive         bool        `json:"is_active"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        *time.Time  `json:"updated_at,omitempty"`
}

// CreateSensorMeasurementTypeRequest represents the request structure for creating a sensor measurement type
type CreateSensorMeasurementTypeRequest struct {
	SensorTypeID     uuid.UUID   `json:"sensor_type_id" binding:"required"`
	Name             string      `json:"name" binding:"required"`
	Description      *string     `json:"description,omitempty"`
	UnitOfMeasure    *string     `json:"unit_of_measure,omitempty"`
	MinAcceptedValue *float64    `json:"min_accepted_value,omitempty"`
	MaxAcceptedValue *float64    `json:"max_accepted_value,omitempty"`
	PropertiesSchema interface{} `json:"properties_schema,omitempty"`
	UIConfiguration  interface{} `json:"ui_configuration,omitempty"`
	Version          int         `json:"version" binding:"required,min=1"`
	IsActive         bool        `json:"is_active"`
}

// UpdateSensorMeasurementTypeRequest represents the request structure for updating a sensor measurement type
type UpdateSensorMeasurementTypeRequest struct {
	Name             *string     `json:"name,omitempty"`
	Description      *string     `json:"description,omitempty"`
	UnitOfMeasure    *string     `json:"unit_of_measure,omitempty"`
	MinAcceptedValue *float64    `json:"min_accepted_value,omitempty"`
	MaxAcceptedValue *float64    `json:"max_accepted_value,omitempty"`
	PropertiesSchema interface{} `json:"properties_schema,omitempty"`
	UIConfiguration  interface{} `json:"ui_configuration,omitempty"`
	Version          *int        `json:"version,omitempty" binding:"omitempty,min=1"`
	IsActive         *bool       `json:"is_active,omitempty"`
}

// SensorMeasurementTypeResponse represents the response structure for sensor measurement type operations
type SensorMeasurementTypeResponse struct {
	Data    SensorMeasurementTypeDTO `json:"data"`
	Message string                   `json:"message"`
}

// SensorMeasurementTypeListResponse represents the response structure for listing sensor measurement types
type SensorMeasurementTypeListResponse struct {
	Data    []SensorMeasurementTypeDTO `json:"data"`
	Total   int                        `json:"total"`
	Page    int                        `json:"page"`
	Limit   int                        `json:"limit"`
	Message string                     `json:"message"`
}
