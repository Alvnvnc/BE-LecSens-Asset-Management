package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LocationRepository handles database operations for locations
type LocationRepository struct {
	*BaseRepository
}

// NewLocationRepository creates a new LocationRepository
func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// GetByID retrieves a location by ID
func (r *LocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Location, error) {
	query := `
		SELECT id, region_code, name, description, address, 
		       longitude, latitude, hierarchy_level, is_active,
		       created_at, updated_at
		FROM locations
		WHERE id = $1
	`

	row := r.DB.QueryRowContext(ctx, query, id)

	var location entity.Location
	var updatedAt sql.NullTime
	var regionCode, description, address sql.NullString

	err := row.Scan(
		&location.ID,
		&regionCode,
		&location.Name,
		&description,
		&address,
		&location.Longitude,
		&location.Latitude,
		&location.HierarchyLevel,
		&location.IsActive,
		&location.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("location not found: %v", id)
		}
		return nil, err
	}

	// Handle nullable fields
	if regionCode.Valid {
		location.RegionCode = regionCode.String
	}
	if description.Valid {
		location.Description = description.String
	}
	if address.Valid {
		location.Address = address.String
	}
	if updatedAt.Valid {
		location.UpdatedAt = updatedAt.Time
	}

	return &location, nil
}

// List retrieves locations with pagination
func (r *LocationRepository) List(ctx context.Context, limit, offset int) ([]*entity.Location, error) {
	query := `
		SELECT id, region_code, name, description, address, 
		       longitude, latitude, hierarchy_level, is_active,
		       created_at, updated_at
		FROM locations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*entity.Location
	for rows.Next() {
		var location entity.Location
		var updatedAt sql.NullTime
		var regionCode, description, address sql.NullString

		err := rows.Scan(
			&location.ID,
			&regionCode,
			&location.Name,
			&description,
			&address,
			&location.Longitude,
			&location.Latitude,
			&location.HierarchyLevel,
			&location.IsActive,
			&location.CreatedAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if regionCode.Valid {
			location.RegionCode = regionCode.String
		}
		if description.Valid {
			location.Description = description.String
		}
		if address.Valid {
			location.Address = address.String
		}
		if updatedAt.Valid {
			location.UpdatedAt = updatedAt.Time
		}

		locations = append(locations, &location)
	}

	return locations, nil
}

// Create creates a new location
func (r *LocationRepository) Create(ctx context.Context, location *entity.Location) error {
	location.ID = uuid.New()
	location.CreatedAt = time.Now()
	location.UpdatedAt = time.Now()

	query := `
		INSERT INTO locations (id, region_code, name, description, address, 
		                      longitude, latitude, hierarchy_level, is_active,
		                      created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		location.ID,
		location.RegionCode,
		location.Name,
		location.Description,
		location.Address,
		location.Longitude,
		location.Latitude,
		location.HierarchyLevel,
		location.IsActive,
		location.CreatedAt,
		location.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}

	return nil
}

// Update updates an existing location
func (r *LocationRepository) Update(ctx context.Context, location *entity.Location) error {
	// First get the existing location
	existingLocation, err := r.GetByID(ctx, location.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing location: %w", err)
	}

	// Update only the fields that are provided (non-empty)
	if location.RegionCode != "" {
		existingLocation.RegionCode = location.RegionCode
	}
	if location.Name != "" {
		existingLocation.Name = location.Name
	}
	if location.Description != "" {
		existingLocation.Description = location.Description
	}
	if location.Address != "" {
		existingLocation.Address = location.Address
	}
	if location.Longitude != 0 {
		existingLocation.Longitude = location.Longitude
	}
	if location.Latitude != 0 {
		existingLocation.Latitude = location.Latitude
	}
	if location.HierarchyLevel != 0 {
		existingLocation.HierarchyLevel = location.HierarchyLevel
	}
	existingLocation.IsActive = location.IsActive
	existingLocation.UpdatedAt = time.Now()

	query := `
		UPDATE locations
		SET region_code = $1,
			name = $2,
			description = $3,
			address = $4,
			longitude = $5,
			latitude = $6,
			hierarchy_level = $7,
			is_active = $8,
			updated_at = $9
		WHERE id = $10
	`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		existingLocation.RegionCode,
		existingLocation.Name,
		existingLocation.Description,
		existingLocation.Address,
		existingLocation.Longitude,
		existingLocation.Latitude,
		existingLocation.HierarchyLevel,
		existingLocation.IsActive,
		existingLocation.UpdatedAt,
		existingLocation.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("location not found: %v", location.ID)
	}

	// Update the input location with the merged values
	*location = *existingLocation

	return nil
}

// Delete deletes a location by ID
func (r *LocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM locations WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("location not found: %v", id)
	}

	return nil
}
