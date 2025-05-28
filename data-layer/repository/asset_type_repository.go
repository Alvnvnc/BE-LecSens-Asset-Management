package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// AssetTypeRepository handles database operations for asset types
type AssetTypeRepository struct {
	*BaseRepository
}

// NewAssetTypeRepository creates a new AssetTypeRepository
func NewAssetTypeRepository(db *sql.DB) *AssetTypeRepository {
	return &AssetTypeRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create inserts a new asset type
func (r *AssetTypeRepository) Create(ctx context.Context, assetType *entity.AssetType) error {
	query := `
		INSERT INTO asset_types (
			id, name, category, description, properties_schema, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	if assetType.ID == uuid.Nil {
		assetType.ID = uuid.New()
	}

	var schemaJSON interface{}
	if len(assetType.PropertiesSchema) > 0 {
		// Validate JSON format before storing
		var temp interface{}
		if err := json.Unmarshal(assetType.PropertiesSchema, &temp); err != nil {
			return fmt.Errorf("invalid JSON format for properties_schema: %v", err)
		}
		schemaJSON = assetType.PropertiesSchema
	} else {
		schemaJSON = []byte("{}")
	}

	_, err := r.DB.ExecContext(
		ctx,
		query,
		assetType.ID,
		assetType.Name,
		assetType.Category,
		assetType.Description,
		schemaJSON,
		assetType.CreatedAt,
		assetType.UpdatedAt,
	)

	return err
}

// GetByID retrieves an asset type by ID
func (r *AssetTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AssetType, error) {
	query := `
		SELECT id, name, category, description, properties_schema, created_at, updated_at
		FROM asset_types
		WHERE id = $1
	`

	row := r.DB.QueryRowContext(ctx, query, id)

	var assetType entity.AssetType
	var schemaJSON []byte
	var updatedAt sql.NullTime
	var description sql.NullString

	err := row.Scan(
		&assetType.ID,
		&assetType.Name,
		&assetType.Category,
		&description,
		&schemaJSON,
		&assetType.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset type not found: %v", id)
		}
		return nil, err
	}

	if description.Valid {
		assetType.Description = description.String
	}

	if len(schemaJSON) > 0 {
		assetType.PropertiesSchema = json.RawMessage(schemaJSON)
	}

	if updatedAt.Valid {
		assetType.UpdatedAt = &updatedAt.Time
	}

	return &assetType, nil
}

// List retrieves asset types with pagination
func (r *AssetTypeRepository) List(ctx context.Context, limit, offset int) ([]*entity.AssetType, error) {
	log.Printf("AssetTypeRepository: Starting List - limit: %d, offset: %d", limit, offset)

	query := `
		SELECT id, name, category, description, properties_schema, created_at, updated_at
		FROM asset_types
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	log.Printf("AssetTypeRepository: Executing query: %s", query)
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Printf("AssetTypeRepository: Error executing query: %v", err)
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var assetTypes []*entity.AssetType

	for rows.Next() {
		var assetType entity.AssetType
		var schemaJSON []byte
		var updatedAt sql.NullTime
		var description sql.NullString

		err := rows.Scan(
			&assetType.ID,
			&assetType.Name,
			&assetType.Category,
			&description,
			&schemaJSON,
			&assetType.CreatedAt,
			&updatedAt,
		)

		if err != nil {
			log.Printf("AssetTypeRepository: Error scanning row: %v", err)
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if description.Valid {
			assetType.Description = description.String
		}

		if len(schemaJSON) > 0 {
			assetType.PropertiesSchema = json.RawMessage(schemaJSON)
		}

		if updatedAt.Valid {
			assetType.UpdatedAt = &updatedAt.Time
		}

		assetTypes = append(assetTypes, &assetType)
	}

	if err = rows.Err(); err != nil {
		log.Printf("AssetTypeRepository: Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Printf("AssetTypeRepository: Successfully retrieved %d asset types", len(assetTypes))
	return assetTypes, nil
}

// Update modifies an existing asset type
func (r *AssetTypeRepository) Update(ctx context.Context, assetType *entity.AssetType) error {
	query := `
		UPDATE asset_types 
		SET name = $2, category = $3, description = $4, properties_schema = $5, updated_at = $6
		WHERE id = $1
	`

	var schemaJSON interface{}
	if len(assetType.PropertiesSchema) > 0 {
		// Validate JSON format before storing
		var temp interface{}
		if err := json.Unmarshal(assetType.PropertiesSchema, &temp); err != nil {
			return fmt.Errorf("invalid JSON format for properties_schema: %v", err)
		}
		schemaJSON = assetType.PropertiesSchema
	} else {
		schemaJSON = []byte("{}")
	}

	result, err := r.DB.ExecContext(
		ctx,
		query,
		assetType.ID,
		assetType.Name,
		assetType.Category,
		assetType.Description,
		schemaJSON,
		assetType.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset type not found: %v", assetType.ID)
	}

	return nil
}

// Delete removes an asset type
func (r *AssetTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM asset_types WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset type not found: %v", id)
	}

	return nil
}
