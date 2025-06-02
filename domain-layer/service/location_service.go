package service

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// LocationService provides business logic for location operations
type LocationService struct {
	locationRepo *repository.LocationRepository
}

// NewLocationService creates a new LocationService
func NewLocationService(locationRepo *repository.LocationRepository) *LocationService {
	return &LocationService{
		locationRepo: locationRepo,
	}
}

// GetLocationByID retrieves a location by ID
func (s *LocationService) GetLocationByID(ctx context.Context, id uuid.UUID) (*entity.Location, error) {
	return s.locationRepo.GetByID(ctx, id)
}

// ListLocations retrieves a paginated list of locations
func (s *LocationService) ListLocations(ctx context.Context, page, pageSize int) ([]*entity.Location, error) {
	log.Printf("Location Service: Starting ListLocations - page: %d, pageSize: %d", page, pageSize)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	log.Printf("Location Service: Calculated offset: %d", offset)

	locations, err := s.locationRepo.List(ctx, pageSize, offset)
	if err != nil {
		log.Printf("Location Service: Error from repository: %v", err)
		return nil, fmt.Errorf("failed to retrieve locations from repository: %w", err)
	}

	log.Printf("Location Service: Successfully retrieved %d locations", len(locations))
	return locations, nil
}

func (s *LocationService) CreateLocation(ctx context.Context, req dto.CreateLocationRequest) (*dto.LocationResponse, error) {
	// Create location entity
	location := &entity.Location{
		RegionCode:     req.RegionCode,
		Name:           req.Name,
		HierarchyLevel: req.HierarchyLevel,
		IsActive:       true, // Default to active
	}

	// Handle optional fields
	if req.Description != nil {
		location.Description = *req.Description
	}
	if req.Address != nil {
		location.Address = *req.Address
	}
	if req.Longitude != nil {
		location.Longitude = *req.Longitude
	}
	if req.Latitude != nil {
		location.Latitude = *req.Latitude
	}

	// Create in database
	err := s.locationRepo.Create(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	// Convert to response
	return &dto.LocationResponse{
		ID:             location.ID,
		RegionCode:     location.RegionCode,
		Name:           location.Name,
		Description:    &location.Description,
		Address:        &location.Address,
		Longitude:      &location.Longitude,
		Latitude:       &location.Latitude,
		HierarchyLevel: location.HierarchyLevel,
		IsActive:       location.IsActive,
		CreatedAt:      location.CreatedAt,
		UpdatedAt:      location.UpdatedAt,
	}, nil
}

func (s *LocationService) UpdateLocation(ctx context.Context, id uuid.UUID, req dto.UpdateLocationRequest) (*dto.LocationResponse, error) {
	// Get existing location
	location, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	// Update fields
	if req.RegionCode != nil {
		location.RegionCode = *req.RegionCode
	}
	if req.Name != nil {
		location.Name = *req.Name
	}
	if req.Description != nil {
		location.Description = *req.Description
	}
	if req.Address != nil {
		location.Address = *req.Address
	}
	if req.Longitude != nil {
		location.Longitude = *req.Longitude
	}
	if req.Latitude != nil {
		location.Latitude = *req.Latitude
	}
	if req.HierarchyLevel != nil {
		location.HierarchyLevel = *req.HierarchyLevel
	}
	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}

	location.UpdatedAt = time.Now()

	// Update in database
	err = s.locationRepo.Update(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}

	// Convert to response
	return &dto.LocationResponse{
		ID:             location.ID,
		RegionCode:     location.RegionCode,
		Name:           location.Name,
		Description:    &location.Description,
		Address:        &location.Address,
		Longitude:      &location.Longitude,
		Latitude:       &location.Latitude,
		HierarchyLevel: location.HierarchyLevel,
		IsActive:       location.IsActive,
		CreatedAt:      location.CreatedAt,
		UpdatedAt:      location.UpdatedAt,
	}, nil
}

func (s *LocationService) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	// Check if location exists
	_, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get location: %w", err)
	}

	// Delete location
	err = s.locationRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}

	return nil
}
