package service

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// AssetTypeService provides business logic for asset type operations
type AssetTypeService struct {
	assetTypeRepo *repository.AssetTypeRepository
}

// NewAssetTypeService creates a new AssetTypeService
func NewAssetTypeService(assetTypeRepo *repository.AssetTypeRepository) *AssetTypeService {
	return &AssetTypeService{
		assetTypeRepo: assetTypeRepo,
	}
}

// GetAssetTypeByID retrieves an asset type by ID
func (s *AssetTypeService) GetAssetTypeByID(ctx context.Context, id uuid.UUID) (*entity.AssetType, error) {
	return s.assetTypeRepo.GetByID(ctx, id)
}

// ListAssetTypes retrieves a paginated list of asset types
func (s *AssetTypeService) ListAssetTypes(ctx context.Context, page, pageSize int) ([]*entity.AssetType, error) {
	log.Printf("AssetTypeService: Starting ListAssetTypes - page: %d, pageSize: %d", page, pageSize)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	log.Printf("AssetTypeService: Calculated offset: %d", offset)

	assetTypes, err := s.assetTypeRepo.List(ctx, pageSize, offset)
	if err != nil {
		log.Printf("AssetTypeService: Error from repository: %v", err)
		return nil, fmt.Errorf("failed to retrieve asset types from repository: %w", err)
	}

	log.Printf("AssetTypeService: Successfully retrieved %d asset types", len(assetTypes))
	return assetTypes, nil
}

// CreateAssetType creates a new asset type
func (s *AssetTypeService) CreateAssetType(ctx context.Context, assetType *entity.AssetType) error {
	// Set creation time
	now := time.Now()
	assetType.CreatedAt = now

	// Ensure properties_schema has a default value if not provided
	if assetType.PropertiesSchema == nil {
		assetType.PropertiesSchema = json.RawMessage("{}")
	}

	return s.assetTypeRepo.Create(ctx, assetType)
}

// UpdateAssetType updates an existing asset type
func (s *AssetTypeService) UpdateAssetType(ctx context.Context, assetType *entity.AssetType) error {
	// Set update time
	now := time.Now()
	assetType.UpdatedAt = &now

	// Ensure properties_schema has a default value if not provided
	if len(assetType.PropertiesSchema) == 0 {
		assetType.PropertiesSchema = json.RawMessage("{}")
	}

	return s.assetTypeRepo.Update(ctx, assetType)
}

// DeleteAssetType deletes an asset type
func (s *AssetTypeService) DeleteAssetType(ctx context.Context, id uuid.UUID) error {
	return s.assetTypeRepo.Delete(ctx, id)
}
