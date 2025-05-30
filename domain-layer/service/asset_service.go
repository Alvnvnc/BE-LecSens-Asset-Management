package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"

	"github.com/google/uuid"
)

// Asset represents an asset in the system
type Asset struct {
	ID          uuid.UUID      `json:"id"`
	TenantID    *uuid.UUID     `json:"tenant_id,omitempty"`
	Name        string         `json:"name"`
	AssetTypeID uuid.UUID      `json:"asset_type_id"`
	LocationID  uuid.UUID      `json:"location_id"`
	Status      string         `json:"status"`
	Properties  map[string]any `json:"properties,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
}

// AssetStatus represents the possible statuses of an asset
type AssetStatus string

const (
	AssetStatusActive      AssetStatus = "active"
	AssetStatusInactive    AssetStatus = "inactive"
	AssetStatusMaintenance AssetStatus = "maintenance"
)

// AssetService handles business logic for assets
type AssetService struct {
	assetRepo     repository.AssetRepository
	assetTypeRepo *repository.AssetTypeRepository
	locationRepo  *repository.LocationRepository
}

// NewAssetService creates a new instance of AssetService
func NewAssetService(
	assetRepo repository.AssetRepository,
	assetTypeRepo *repository.AssetTypeRepository,
	locationRepo *repository.LocationRepository,
) *AssetService {
	return &AssetService{
		assetRepo:     assetRepo,
		assetTypeRepo: assetTypeRepo,
		locationRepo:  locationRepo,
	}
}

// CreateAsset creates a new asset
func (s *AssetService) CreateAsset(ctx context.Context, asset *entity.Asset) error {
	// Validate asset type
	_, err := s.assetTypeRepo.GetByID(ctx, asset.AssetTypeID)
	if err != nil {
		return errors.New("invalid asset type")
	}

	// Validate location
	_, err = s.locationRepo.GetByID(ctx, asset.LocationID)
	if err != nil {
		return errors.New("invalid location")
	}

	// Set timestamps
	now := time.Now()
	asset.CreatedAt = now
	asset.UpdatedAt = now

	// Set default status if not provided
	if asset.Status == "" {
		asset.Status = "active"
	}

	return s.assetRepo.Create(ctx, asset)
}

// GetAsset retrieves an asset by ID
func (s *AssetService) GetAsset(ctx context.Context, id uuid.UUID) (*entity.Asset, error) {
	return s.assetRepo.GetByID(ctx, id)
}

// ListAssets retrieves a list of assets with pagination
func (s *AssetService) ListAssets(ctx context.Context, tenantID *uuid.UUID, page, pageSize int) ([]*entity.Asset, error) {
	return s.assetRepo.List(ctx, tenantID, page, pageSize)
}

// UpdateAsset updates an existing asset
func (s *AssetService) UpdateAsset(ctx context.Context, asset *entity.Asset) error {
	// Get existing asset
	existingAsset, err := s.assetRepo.GetByID(ctx, asset.ID)
	if err != nil {
		return err
	}

	// Validate asset type if changed
	if asset.AssetTypeID != existingAsset.AssetTypeID {
		_, err := s.assetTypeRepo.GetByID(ctx, asset.AssetTypeID)
		if err != nil {
			return errors.New("invalid asset type")
		}
	}

	// Validate location if changed
	if asset.LocationID != existingAsset.LocationID {
		_, err := s.locationRepo.GetByID(ctx, asset.LocationID)
		if err != nil {
			return errors.New("invalid location")
		}
	}

	// Update timestamp
	asset.UpdatedAt = time.Now()

	// Preserve tenant ID
	asset.TenantID = existingAsset.TenantID

	return s.assetRepo.Update(ctx, asset)
}

// DeleteAsset deletes an asset
func (s *AssetService) DeleteAsset(ctx context.Context, id uuid.UUID) error {
	return s.assetRepo.Delete(ctx, id)
}

// AssignAssetToTenant assigns an asset to a tenant
func (s *AssetService) AssignAssetToTenant(ctx context.Context, id, tenantID uuid.UUID) error {
	// Get existing asset
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if asset is already assigned to a tenant
	if asset.TenantID != nil {
		return errors.New("asset is already assigned to a tenant")
	}

	return s.assetRepo.AssignToTenant(ctx, id, tenantID)
}

// UnassignAssetFromTenant removes tenant assignment from an asset
func (s *AssetService) UnassignAssetFromTenant(ctx context.Context, id uuid.UUID) error {
	// Get existing asset
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if asset is assigned to a tenant
	if asset.TenantID == nil {
		return errors.New("asset is not assigned to any tenant")
	}

	return s.assetRepo.UnassignFromTenant(ctx, id)
}

// GetAssetTypeByID retrieves an asset type by ID
func (s *AssetService) GetAssetTypeByID(ctx context.Context, id uuid.UUID) (*entity.AssetType, error) {
	return s.assetTypeRepo.GetByID(ctx, id)
}

// GetLocationByID retrieves a location by ID
func (s *AssetService) GetLocationByID(ctx context.Context, id uuid.UUID) (*entity.Location, error) {
	return s.locationRepo.GetByID(ctx, id)
}

// UpdateAssetPartial updates specific fields of an existing asset
func (s *AssetService) UpdateAssetPartial(ctx context.Context, id uuid.UUID, updateReq interface{}) (*entity.Asset, error) {
	// Get existing asset
	existingAsset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert updateReq to map for easier handling
	updateBytes, err := json.Marshal(updateReq)
	if err != nil {
		return nil, err
	}

	var updateRequest map[string]interface{}
	err = json.Unmarshal(updateBytes, &updateRequest)
	if err != nil {
		return nil, err
	}

	// Create a copy of existing asset for updates
	updatedAsset := *existingAsset

	// Update only provided fields
	if name, exists := updateRequest["name"]; exists && name != nil {
		updatedAsset.Name = name.(string)
	}

	if assetTypeID, exists := updateRequest["asset_type_id"]; exists && assetTypeID != nil {
		var parsedAssetTypeID uuid.UUID

		// Handle different possible types for UUID
		switch v := assetTypeID.(type) {
		case string:
			parsedAssetTypeID, err = uuid.Parse(v)
			if err != nil {
				return nil, errors.New("invalid asset type ID format")
			}
		case uuid.UUID:
			parsedAssetTypeID = v
		default:
			return nil, errors.New("invalid asset type ID format")
		}

		// Validate asset type exists
		_, err = s.assetTypeRepo.GetByID(ctx, parsedAssetTypeID)
		if err != nil {
			return nil, errors.New("invalid asset type")
		}
		updatedAsset.AssetTypeID = parsedAssetTypeID
	}

	if locationID, exists := updateRequest["location_id"]; exists && locationID != nil {
		var parsedLocationID uuid.UUID

		// Handle different possible types for UUID
		switch v := locationID.(type) {
		case string:
			parsedLocationID, err = uuid.Parse(v)
			if err != nil {
				return nil, errors.New("invalid location ID format")
			}
		case uuid.UUID:
			parsedLocationID = v
		default:
			return nil, errors.New("invalid location ID format")
		}

		// Validate location exists
		_, err = s.locationRepo.GetByID(ctx, parsedLocationID)
		if err != nil {
			return nil, errors.New("invalid location")
		}
		updatedAsset.LocationID = parsedLocationID
	}

	if status, exists := updateRequest["status"]; exists && status != nil {
		updatedAsset.Status = status.(string)
	}

	if properties, exists := updateRequest["properties"]; exists && properties != nil {
		propertiesBytes, _ := json.Marshal(properties)
		updatedAsset.Properties = json.RawMessage(propertiesBytes)
	}

	// Handle tenant_id updates - this is for fixing assets without tenant assignment
	if tenantID, exists := updateRequest["tenant_id"]; exists && tenantID != nil {
		var parsedTenantID uuid.UUID

		// Handle different possible types for UUID
		switch v := tenantID.(type) {
		case string:
			parsedTenantID, err = uuid.Parse(v)
			if err != nil {
				return nil, errors.New("invalid tenant ID format")
			}
		case uuid.UUID:
			parsedTenantID = v
		default:
			return nil, errors.New("invalid tenant ID format")
		}

		updatedAsset.TenantID = &parsedTenantID
	} else {
		// IMPORTANT: Preserve existing tenant ID to ensure it doesn't get lost during update
		// This is crucial for multi-tenant data integrity
		updatedAsset.TenantID = existingAsset.TenantID
	}

	// Update timestamp
	updatedAsset.UpdatedAt = time.Now()

	// Update in repository
	err = s.assetRepo.Update(ctx, &updatedAsset)
	if err != nil {
		return nil, err
	}

	return &updatedAsset, nil
}
