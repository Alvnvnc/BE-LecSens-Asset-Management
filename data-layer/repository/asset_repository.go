package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AssetRepository defines the interface for asset data operations
type AssetRepository interface {
	Create(ctx context.Context, asset *entity.Asset) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Asset, error)
	List(ctx context.Context, tenantID *uuid.UUID, page, pageSize int) ([]*entity.Asset, error)
	Update(ctx context.Context, asset *entity.Asset) error
	Delete(ctx context.Context, id uuid.UUID) error
	AssignToTenant(ctx context.Context, assetID, tenantID uuid.UUID) error
	UnassignFromTenant(ctx context.Context, assetID uuid.UUID) error
}

// assetRepository handles database operations for assets
type assetRepository struct {
	db *sql.DB
}

// NewAssetRepository creates a new AssetRepository
func NewAssetRepository(db *sql.DB) AssetRepository {
	return &assetRepository{db: db}
}

// Create inserts a new asset into the database
func (r *assetRepository) Create(ctx context.Context, asset *entity.Asset) error {
	query := `
		INSERT INTO assets (
			id, tenant_id, name, asset_type_id, location_id, status, properties, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`

	now := time.Now()
	_, err := r.db.ExecContext(
		ctx,
		query,
		asset.ID,
		asset.TenantID,
		asset.Name,
		asset.AssetTypeID,
		asset.LocationID,
		asset.Status,
		asset.Properties,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}

	asset.CreatedAt = now
	asset.UpdatedAt = now
	return nil
}

// GetByID retrieves an asset by its ID
func (r *assetRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Asset, error) {
	query := `
		SELECT id, tenant_id, name, asset_type_id, location_id, status, properties, created_at, updated_at
		FROM assets
		WHERE id = $1`

	var asset entity.Asset
	var tenantID sql.NullString
	var properties []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&asset.ID,
		&tenantID,
		&asset.Name,
		&asset.AssetTypeID,
		&asset.LocationID,
		&asset.Status,
		&properties,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	if tenantID.Valid {
		parsedID, err := uuid.Parse(tenantID.String)
		if err != nil {
			return nil, fmt.Errorf("invalid tenant ID format: %w", err)
		}
		asset.TenantID = &parsedID
	}

	if properties != nil {
		asset.Properties = properties
	}

	return &asset, nil
}

// List retrieves a paginated list of assets
func (r *assetRepository) List(ctx context.Context, tenantID *uuid.UUID, page, pageSize int) ([]*entity.Asset, error) {
	query := `
		SELECT id, tenant_id, name, asset_type_id, location_id, status, properties, created_at, updated_at
		FROM assets
		WHERE ($1::uuid IS NULL OR tenant_id = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, tenantID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	defer rows.Close()

	var assets []*entity.Asset
	for rows.Next() {
		var asset entity.Asset
		var tenantID sql.NullString
		var properties []byte
		err := rows.Scan(
			&asset.ID,
			&tenantID,
			&asset.Name,
			&asset.AssetTypeID,
			&asset.LocationID,
			&asset.Status,
			&properties,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}

		if tenantID.Valid {
			parsedID, err := uuid.Parse(tenantID.String)
			if err != nil {
				return nil, fmt.Errorf("invalid tenant ID format: %w", err)
			}
			asset.TenantID = &parsedID
		}

		if properties != nil {
			asset.Properties = properties
		}

		assets = append(assets, &asset)
	}

	return assets, nil
}

// Update modifies an existing asset
func (r *assetRepository) Update(ctx context.Context, asset *entity.Asset) error {
	query := `
		UPDATE assets
		SET name = $1, asset_type_id = $2, location_id = $3, status = $4, properties = $5, tenant_id = $6, updated_at = $7
		WHERE id = $8`

	now := time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		asset.Name,
		asset.AssetTypeID,
		asset.LocationID,
		asset.Status,
		asset.Properties,
		asset.TenantID,
		now,
		asset.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update asset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("asset not found")
	}

	asset.UpdatedAt = now
	return nil
}

// Delete removes an asset by its ID
func (r *assetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM assets WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("asset not found")
	}

	return nil
}

// AssignToTenant assigns an asset to a tenant
func (r *assetRepository) AssignToTenant(ctx context.Context, assetID, tenantID uuid.UUID) error {
	query := `
		UPDATE assets
		SET tenant_id = $1, updated_at = $2
		WHERE id = $3 AND tenant_id IS NULL`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, tenantID, now, assetID)
	if err != nil {
		return fmt.Errorf("failed to assign asset to tenant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("asset not found or already assigned to a tenant")
	}

	return nil
}

// UnassignFromTenant removes the tenant assignment from an asset
func (r *assetRepository) UnassignFromTenant(ctx context.Context, assetID uuid.UUID) error {
	query := `
		UPDATE assets
		SET tenant_id = NULL, updated_at = $1
		WHERE id = $2 AND tenant_id IS NOT NULL`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, assetID)
	if err != nil {
		return fmt.Errorf("failed to unassign asset from tenant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("asset not found or not assigned to any tenant")
	}

	return nil
}
