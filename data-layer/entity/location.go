package entity

import (
	"time"

	"github.com/google/uuid"
)

// Location represents a physical location in the hierarchy
type Location struct {
	ID             uuid.UUID `json:"id"`
	RegionCode     string    `json:"region_code,omitempty"` // Code representing the region (province/city/district)
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	Address        string    `json:"address,omitempty"`
	Longitude      float64   `json:"longitude"` // Geographic longitude coordinate
	Latitude       float64   `json:"latitude"`  // Geographic latitude coordinate
	HierarchyLevel int       `json:"hierarchy_level"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
