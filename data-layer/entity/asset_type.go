package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AssetType represents a category of assets in the system
type AssetType struct {
	ID               uuid.UUID       `json:"id"`
	Name             string          `json:"name"`
	Category         string          `json:"category"`
	Description      string          `json:"description"`
	PropertiesSchema json.RawMessage `json:"properties_schema"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        *time.Time      `json:"updated_at,omitempty"`
}
