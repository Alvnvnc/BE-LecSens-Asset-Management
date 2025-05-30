package dto

import (
	"time"

	"github.com/google/uuid"
)

// SensorTypeDTO represents the data transfer object for sensor type
type SensorTypeDTO struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     *uuid.UUID `json:"tenant_id,omitempty"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Manufacturer string     `json:"manufacturer"`
	Model        string     `json:"model"`
	Version      string     `json:"version"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

// CreateSensorTypeRequest represents the request for creating a new sensor type
type CreateSensorTypeRequest struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	Manufacturer string `json:"manufacturer" validate:"required"`
	Model        string `json:"model" validate:"required"`
	Version      string `json:"version" validate:"required"`
	IsActive     bool   `json:"is_active"`
}

// UpdateSensorTypeRequest represents the request for updating a sensor type
type UpdateSensorTypeRequest struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
	Version      string `json:"version,omitempty"`
	IsActive     bool   `json:"is_active,omitempty"`
}

// SensorTypeResponse represents the response for sensor type operations
type SensorTypeResponse struct {
	Success bool           `json:"success"`
	Data    *SensorTypeDTO `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

// SensorTypeListResponse represents the response for listing sensor types
type SensorTypeListResponse struct {
	Success  bool             `json:"success"`
	Data     []*SensorTypeDTO `json:"data,omitempty"`
	Error    string           `json:"error,omitempty"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}
