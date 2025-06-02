package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateLocationRequest represents the request to create a new location
type CreateLocationRequest struct {
	RegionCode     string   `json:"region_code" binding:"required"`
	Name           string   `json:"name" binding:"required"`
	Description    *string  `json:"description,omitempty"`
	Address        *string  `json:"address,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	Latitude       *float64 `json:"latitude,omitempty"`
	HierarchyLevel int      `json:"hierarchy_level" binding:"required"`
}

// UpdateLocationRequest represents the request to update an existing location
type UpdateLocationRequest struct {
	RegionCode     *string  `json:"region_code,omitempty"`
	Name           *string  `json:"name,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Address        *string  `json:"address,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	Latitude       *float64 `json:"latitude,omitempty"`
	HierarchyLevel *int     `json:"hierarchy_level,omitempty"`
	IsActive       *bool    `json:"is_active,omitempty"`
}

// LocationResponse represents the response for location operations
type LocationResponse struct {
	ID             uuid.UUID `json:"id"`
	RegionCode     string    `json:"region_code"`
	Name           string    `json:"name"`
	Description    *string   `json:"description,omitempty"`
	Address        *string   `json:"address,omitempty"`
	Longitude      *float64  `json:"longitude,omitempty"`
	Latitude       *float64  `json:"latitude,omitempty"`
	HierarchyLevel int       `json:"hierarchy_level"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
