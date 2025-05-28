package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AssetStatus represents the possible states of an asset
type AssetStatus string

const (
	AssetStatusActive      AssetStatus = "active"
	AssetStatusInactive    AssetStatus = "inactive"
	AssetStatusMaintenance AssetStatus = "maintenance"
)

// Asset represents a physical or digital asset in the system
type Asset struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    *uuid.UUID      `json:"tenant_id,omitempty"`
	Name        string          `json:"name"`
	AssetTypeID uuid.UUID       `json:"asset_type_id"`
	LocationID  uuid.UUID       `json:"location_id"`
	Status      string          `json:"status"`
	Properties  json.RawMessage `json:"properties,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
