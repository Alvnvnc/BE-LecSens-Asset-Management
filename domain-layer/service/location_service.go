package service

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
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

// UpdateLocation updates an existing location
func (s *LocationService) UpdateLocation(ctx context.Context, location *entity.Location) error {
	log.Printf("Location Service: Updating location - ID: %s", location.ID)

	// Set update time
	location.UpdatedAt = time.Now()

	err := s.locationRepo.Update(ctx, location)
	if err != nil {
		log.Printf("Location Service: Error updating location: %v", err)
		return fmt.Errorf("failed to update location: %w", err)
	}

	log.Printf("Location Service: Successfully updated location: %s", location.ID)
	return nil
}
