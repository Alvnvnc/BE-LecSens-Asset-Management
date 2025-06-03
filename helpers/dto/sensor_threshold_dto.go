package dto

import (
	"be-lecsens/asset_management/data-layer/entity"
	"time"

	"github.com/google/uuid"
)

// SensorThresholdResponse represents the response structure for sensor thresholds
type SensorThresholdResponse struct {
	ID                   uuid.UUID                `json:"id"`
	TenantID             uuid.UUID                `json:"tenant_id"`
	AssetSensorID        uuid.UUID                `json:"asset_sensor_id"`
	MeasurementTypeID    uuid.UUID                `json:"measurement_type_id"`
	MeasurementFieldName string                   `json:"measurement_field_name"`
	MinValue             *float64                 `json:"min_value,omitempty"`
	MaxValue             *float64                 `json:"max_value,omitempty"`
	Severity             entity.ThresholdSeverity `json:"severity"`
	IsActive             bool                     `json:"is_active"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            *time.Time               `json:"updated_at,omitempty"`
}

// SensorThresholdListResponse represents the paginated response for listing sensor thresholds
type SensorThresholdListResponse struct {
	Data       []SensorThresholdResponse `json:"data"`
	Pagination PaginationInfo            `json:"pagination"`
}

// CreateSensorThresholdRequest represents the request structure for creating a sensor threshold
type CreateSensorThresholdRequest struct {
	AssetSensorID        uuid.UUID                `json:"asset_sensor_id" binding:"required"`
	MeasurementTypeID    uuid.UUID                `json:"measurement_type_id" binding:"required"`
	MeasurementFieldName string                   `json:"measurement_field_name" binding:"required"`
	MinValue             *float64                 `json:"min_value,omitempty"`
	MaxValue             *float64                 `json:"max_value,omitempty"`
	Severity             entity.ThresholdSeverity `json:"severity" binding:"required"`
	IsActive             bool                     `json:"is_active"`
}

// UpdateSensorThresholdRequest represents the request structure for updating a sensor threshold
type UpdateSensorThresholdRequest struct {
	MeasurementFieldName string                   `json:"measurement_field_name,omitempty"`
	MinValue             *float64                 `json:"min_value,omitempty"`
	MaxValue             *float64                 `json:"max_value,omitempty"`
	Severity             entity.ThresholdSeverity `json:"severity,omitempty"`
	IsActive             *bool                    `json:"is_active,omitempty"`
}

// SensorThresholdFilter represents filter parameters for listing thresholds
type SensorThresholdFilter struct {
	AssetSensorID     *uuid.UUID                `json:"asset_sensor_id,omitempty"`
	MeasurementTypeID *uuid.UUID                `json:"measurement_type_id,omitempty"`
	Severity          *entity.ThresholdSeverity `json:"severity,omitempty"`
	IsActive          *bool                     `json:"is_active,omitempty"`
	Page              int                       `json:"page"`
	Limit             int                       `json:"limit"`
}

// ToEntity converts CreateSensorThresholdRequest to entity.SensorThreshold
func (r *CreateSensorThresholdRequest) ToEntity(tenantID uuid.UUID) *entity.SensorThreshold {
	return &entity.SensorThreshold{
		TenantID:             tenantID,
		AssetSensorID:        r.AssetSensorID,
		MeasurementTypeID:    r.MeasurementTypeID,
		MeasurementFieldName: r.MeasurementFieldName,
		MinValue:             r.MinValue,
		MaxValue:             r.MaxValue,
		Severity:             r.Severity,
		IsActive:             r.IsActive,
	}
}

// ToEntity converts UpdateSensorThresholdRequest to entity.SensorThreshold
func (r *UpdateSensorThresholdRequest) ToEntity(id uuid.UUID) *entity.SensorThreshold {
	return &entity.SensorThreshold{
		ID:                   id,
		MeasurementFieldName: r.MeasurementFieldName,
		MinValue:             r.MinValue,
		MaxValue:             r.MaxValue,
		Severity:             r.Severity,
		IsActive:             r.IsActive != nil && *r.IsActive,
	}
}

// FromEntity converts entity.SensorThreshold to SensorThresholdResponse
func FromEntity(e *entity.SensorThreshold) *SensorThresholdResponse {
	if e == nil {
		return nil
	}
	return &SensorThresholdResponse{
		ID:                   e.ID,
		TenantID:             e.TenantID,
		AssetSensorID:        e.AssetSensorID,
		MeasurementTypeID:    e.MeasurementTypeID,
		MeasurementFieldName: e.MeasurementFieldName,
		MinValue:             e.MinValue,
		MaxValue:             e.MaxValue,
		Severity:             e.Severity,
		IsActive:             e.IsActive,
		CreatedAt:            e.CreatedAt,
		UpdatedAt:            e.UpdatedAt,
	}
}
