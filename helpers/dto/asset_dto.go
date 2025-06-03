package dto

import (
	"github.com/google/uuid"
)

// AssignAssetToTenantRequest represents the request to assign an asset to a tenant
type AssignAssetToTenantRequest struct {
	TenantID uuid.UUID `json:"tenant_id" binding:"required" validate:"required"`
}

// UpdateAssetRequest represents the request body for partial asset updates
type UpdateAssetRequest struct {
	Name        *string    `json:"name,omitempty"`
	AssetTypeID *uuid.UUID `json:"asset_type_id,omitempty"`
	LocationID  *uuid.UUID `json:"location_id,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Properties  *string    `json:"properties,omitempty"`
}

// AssetResponse represents the response structure for asset operations
type AssetResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	AssetTypeID uuid.UUID  `json:"asset_type_id"`
	LocationID  uuid.UUID  `json:"location_id"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty"`
	Status      string     `json:"status"`
	Properties  *string    `json:"properties,omitempty"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

// AssetListResponse represents the response for listing assets with pagination
type AssetListResponse struct {
	Assets     []AssetResponse `json:"assets"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	Total      int64           `json:"total"`
	TotalPages int             `json:"total_pages"`
}
