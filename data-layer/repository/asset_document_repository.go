package repository

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/helpers/common"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AssetDocumentRepository defines the interface for asset document data operations
type AssetDocumentRepository interface {
	Create(ctx context.Context, doc *entity.AssetDocument) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AssetDocument, error)
	GetByAssetID(ctx context.Context, assetID uuid.UUID) ([]*entity.AssetDocument, error)
	GetByAssetIDAndType(ctx context.Context, assetID uuid.UUID, docType string) ([]*entity.AssetDocument, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.AssetDocument, error)
	Update(ctx context.Context, doc *entity.AssetDocument) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByAssetID(ctx context.Context, assetID uuid.UUID) error
	CheckCloudinaryIDExists(ctx context.Context, cloudinaryID string) (bool, error)
	CheckAssetDocumentExists(ctx context.Context, assetID uuid.UUID, documentType string) (bool, error)
	GetExistingDocumentByAssetID(ctx context.Context, assetID uuid.UUID) (*entity.AssetDocument, error)
	CleanupDuplicateDocuments(ctx context.Context, assetID uuid.UUID) (int, error)
	CleanupAllDuplicateDocuments(ctx context.Context) (int, error)
	GetDuplicateDocuments(ctx context.Context) (map[uuid.UUID][]*entity.AssetDocument, error)
}

// assetDocumentRepository handles database operations for asset documents
type assetDocumentRepository struct {
	*BaseRepository
}

// NewAssetDocumentRepository creates a new AssetDocumentRepository
func NewAssetDocumentRepository(db *sql.DB) AssetDocumentRepository {
	return &assetDocumentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create inserts a new asset document into the database
func (r *assetDocumentRepository) Create(ctx context.Context, doc *entity.AssetDocument) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	query := `
		INSERT INTO asset_documents (
			tenant_id, asset_id, document_type, file_url, cloudinary_id, 
			original_filename, file_size, mime_type, uploaded_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id`

	now := time.Now()

	// Use tenant ID from entity (which could be nil for SuperAdmin)
	var entityTenantID *uuid.UUID
	if hasTenantID {
		entityTenantID = &tenantID
	} else {
		entityTenantID = doc.TenantID // Could be nil for SuperAdmin
	}

	err := r.DB.QueryRowContext(
		ctx,
		query,
		entityTenantID,
		doc.AssetID,
		doc.DocumentType,
		doc.FileURL,
		doc.CloudinaryID,
		doc.OriginalFilename,
		doc.FileSize,
		doc.MimeType,
		now,
	).Scan(&doc.ID)

	if err != nil {
		return fmt.Errorf("failed to create asset document: %w", err)
	}

	doc.TenantID = entityTenantID
	doc.UploadedAt = now
	return nil
}

// GetByID retrieves an asset document by its ID
func (r *assetDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AssetDocument, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access any document
		query = `
			SELECT id, tenant_id, asset_id, document_type, file_url, cloudinary_id,
				   original_filename, file_size, mime_type, uploaded_at
			FROM asset_documents
			WHERE id = $1`
		args = []interface{}{id}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT id, tenant_id, asset_id, document_type, file_url, cloudinary_id,
				   original_filename, file_size, mime_type, uploaded_at
			FROM asset_documents
			WHERE id = $1 AND tenant_id = $2`
		args = []interface{}{id, tenantID}
	}

	var doc entity.AssetDocument
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&doc.ID,
		&doc.TenantID,
		&doc.AssetID,
		&doc.DocumentType,
		&doc.FileURL,
		&doc.CloudinaryID,
		&doc.OriginalFilename,
		&doc.FileSize,
		&doc.MimeType,
		&doc.UploadedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get asset document: %w", err)
	}

	return &doc, nil
}

// GetByAssetID retrieves all documents for a specific asset
func (r *assetDocumentRepository) GetByAssetID(ctx context.Context, assetID uuid.UUID) ([]*entity.AssetDocument, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access all documents for the asset
		query = `
			SELECT id, tenant_id, asset_id, document_type, file_url, cloudinary_id,
				   original_filename, file_size, mime_type, uploaded_at
			FROM asset_documents
			WHERE asset_id = $1
			ORDER BY uploaded_at DESC`
		args = []interface{}{assetID}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT id, tenant_id, asset_id, document_type, file_url, cloudinary_id,
				   original_filename, file_size, mime_type, uploaded_at
			FROM asset_documents
			WHERE asset_id = $1 AND tenant_id = $2
			ORDER BY uploaded_at DESC`
		args = []interface{}{assetID, tenantID}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset documents: %w", err)
	}
	defer rows.Close()

	var documents []*entity.AssetDocument
	for rows.Next() {
		var doc entity.AssetDocument
		err := rows.Scan(
			&doc.ID,
			&doc.TenantID,
			&doc.AssetID,
			&doc.DocumentType,
			&doc.FileURL,
			&doc.CloudinaryID,
			&doc.OriginalFilename,
			&doc.FileSize,
			&doc.MimeType,
			&doc.UploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset document: %w", err)
		}
		documents = append(documents, &doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating asset documents: %w", err)
	}

	return documents, nil
}

// GetByAssetIDAndType retrieves documents for a specific asset filtered by document type
func (r *assetDocumentRepository) GetByAssetIDAndType(ctx context.Context, assetID uuid.UUID, docType string) ([]*entity.AssetDocument, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access all documents for the asset
		query = `
			SELECT id, tenant_id, asset_id, document_type, file_url, cloudinary_id,
				   original_filename, file_size, mime_type, uploaded_at
			FROM asset_documents
			WHERE asset_id = $1 AND document_type = $2
			ORDER BY uploaded_at DESC`
		args = []interface{}{assetID, docType}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT id, tenant_id, asset_id, document_type, file_url, cloudinary_id,
				   original_filename, file_size, mime_type, uploaded_at
			FROM asset_documents
			WHERE asset_id = $1 AND tenant_id = $2 AND document_type = $3
			ORDER BY uploaded_at DESC`
		args = []interface{}{assetID, tenantID, docType}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset documents: %w", err)
	}
	defer rows.Close()

	var documents []*entity.AssetDocument
	for rows.Next() {
		var doc entity.AssetDocument
		err := rows.Scan(
			&doc.ID,
			&doc.TenantID,
			&doc.AssetID,
			&doc.DocumentType,
			&doc.FileURL,
			&doc.CloudinaryID,
			&doc.OriginalFilename,
			&doc.FileSize,
			&doc.MimeType,
			&doc.UploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset document: %w", err)
		}
		documents = append(documents, &doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating asset documents: %w", err)
	}

	return documents, nil
}

// List retrieves asset documents with pagination (tenant-filtered for regular users, all for SuperAdmin)
func (r *assetDocumentRepository) List(ctx context.Context, page, pageSize int) ([]*entity.AssetDocument, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, errors.New("tenant ID is required for this operation")
	}

	offset := (page - 1) * pageSize

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can access all documents across all tenants
		query = `
			SELECT id, tenant_id, asset_id, document_type, file_url, cloudinary_id,
				   original_filename, file_size, mime_type, uploaded_at
			FROM asset_documents
			ORDER BY uploaded_at DESC
			LIMIT $1 OFFSET $2`
		args = []interface{}{pageSize, offset}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `
			SELECT id, tenant_id, asset_id, document_type, file_url, cloudinary_id,
				   original_filename, file_size, mime_type, uploaded_at
			FROM asset_documents
			WHERE tenant_id = $1
			ORDER BY uploaded_at DESC
			LIMIT $2 OFFSET $3`
		args = []interface{}{tenantID, pageSize, offset}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list asset documents: %w", err)
	}
	defer rows.Close()

	var documents []*entity.AssetDocument
	for rows.Next() {
		var doc entity.AssetDocument
		err := rows.Scan(
			&doc.ID,
			&doc.TenantID,
			&doc.AssetID,
			&doc.DocumentType,
			&doc.FileURL,
			&doc.CloudinaryID,
			&doc.OriginalFilename,
			&doc.FileSize,
			&doc.MimeType,
			&doc.UploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset document: %w", err)
		}
		documents = append(documents, &doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating asset documents: %w", err)
	}

	return documents, nil
}

// Update modifies an existing asset document
func (r *assetDocumentRepository) Update(ctx context.Context, doc *entity.AssetDocument) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin {
		// SuperAdmin can update any document
		query = `
			UPDATE asset_documents
			SET document_type = $1, file_url = $2, cloudinary_id = $3, 
				original_filename = $4, file_size = $5, mime_type = $6
			WHERE id = $7`
		args = []interface{}{
			doc.DocumentType, doc.FileURL, doc.CloudinaryID,
			doc.OriginalFilename, doc.FileSize, doc.MimeType,
			doc.ID,
		}
	} else {
		// Regular users can only update documents from their tenant
		query = `
			UPDATE asset_documents
			SET document_type = $1, file_url = $2, cloudinary_id = $3, 
				original_filename = $4, file_size = $5, mime_type = $6
			WHERE id = $7 AND tenant_id = $8`
		args = []interface{}{
			doc.DocumentType, doc.FileURL, doc.CloudinaryID,
			doc.OriginalFilename, doc.FileSize, doc.MimeType,
			doc.ID, tenantID,
		}
	}

	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update asset document: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset document not found or access denied")
	}

	return nil
}

// Delete removes an asset document by ID
func (r *assetDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin {
		// SuperAdmin can delete any document
		query = `DELETE FROM asset_documents WHERE id = $1`
		args = []interface{}{id}
	} else {
		// Regular users can only delete documents from their tenant
		query = `DELETE FROM asset_documents WHERE id = $1 AND tenant_id = $2`
		args = []interface{}{id, tenantID}
	}

	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete asset document: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset document not found or access denied")
	}

	return nil
}

// DeleteByAssetID removes all documents for a specific asset
func (r *assetDocumentRepository) DeleteByAssetID(ctx context.Context, assetID uuid.UUID) error {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin {
		// SuperAdmin can delete documents from any asset
		query = `DELETE FROM asset_documents WHERE asset_id = $1`
		args = []interface{}{assetID}
	} else {
		// Regular users can only delete documents from assets in their tenant
		query = `DELETE FROM asset_documents WHERE asset_id = $1 AND tenant_id = $2`
		args = []interface{}{assetID, tenantID}
	}

	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete asset documents: %w", err)
	}

	return nil
}

// CheckCloudinaryIDExists checks if a cloudinary ID already exists for the tenant
func (r *assetDocumentRepository) CheckCloudinaryIDExists(ctx context.Context, cloudinaryID string) (bool, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return false, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin: check globally across all tenants
		query = `SELECT EXISTS(SELECT 1 FROM asset_documents WHERE cloudinary_id = $1)`
		args = []interface{}{cloudinaryID}
	} else {
		// Regular user or SuperAdmin with tenant context: check within tenant
		query = `SELECT EXISTS(SELECT 1 FROM asset_documents WHERE cloudinary_id = $1 AND tenant_id = $2)`
		args = []interface{}{cloudinaryID, tenantID}
	}

	var exists bool
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check cloudinary ID existence: %w", err)
	}

	return exists, nil
}

// CheckAssetDocumentExists checks if an asset already has a document of specific type
func (r *assetDocumentRepository) CheckAssetDocumentExists(ctx context.Context, assetID uuid.UUID, documentType string) (bool, error) {
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return false, errors.New("tenant ID is required for this operation")
	}

	var query string
	var args []interface{}

	if isSuperAdmin && !hasTenantID {
		// SuperAdmin without tenant ID can check across all tenants
		query = `SELECT EXISTS(SELECT 1 FROM asset_documents WHERE asset_id = $1 AND document_type = $2)`
		args = []interface{}{assetID, documentType}
	} else {
		// Regular users or SuperAdmin with tenant ID - filter by tenant
		query = `SELECT EXISTS(SELECT 1 FROM asset_documents WHERE asset_id = $1 AND document_type = $2 AND tenant_id = $3)`
		args = []interface{}{assetID, documentType, tenantID}
	}

	var exists bool
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check asset document existence: %w", err)
	}

	return exists, nil
}

// GetExistingDocumentByAssetID gets the oldest document for an asset (to be replaced)
func (r *assetDocumentRepository) GetExistingDocumentByAssetID(ctx context.Context, assetID uuid.UUID) (*entity.AssetDocument, error) {
	query := `
		SELECT id, tenant_id, asset_id, document_type, file_url, cloudinary_id, 
		       original_filename, file_size, mime_type, uploaded_at, created_at, updated_at
		FROM asset_documents 
		WHERE asset_id = $1 
		ORDER BY created_at ASC 
		LIMIT 1`

	var doc entity.AssetDocument
	var tenantID sql.NullString

	err := r.DB.QueryRowContext(ctx, query, assetID).Scan(
		&doc.ID,
		&tenantID,
		&doc.AssetID,
		&doc.DocumentType,
		&doc.FileURL,
		&doc.CloudinaryID,
		&doc.OriginalFilename,
		&doc.FileSize,
		&doc.MimeType,
		&doc.UploadedAt,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No existing document found
		}
		return nil, fmt.Errorf("failed to get existing document: %w", err)
	}

	// Handle nullable tenant_id
	if tenantID.Valid {
		if tenantUUID, err := uuid.Parse(tenantID.String); err == nil {
			doc.TenantID = &tenantUUID
		}
	}

	return &doc, nil
}

// CleanupDuplicateDocuments removes duplicate documents for the same asset, keeping only the latest one
func (r *assetDocumentRepository) CleanupDuplicateDocuments(ctx context.Context, assetID uuid.UUID) (int, error) {
	// Query to find and delete duplicates, keeping only the latest document per asset
	query := `
		DELETE FROM asset_documents 
		WHERE asset_id = $1 
		AND id NOT IN (
			SELECT id FROM asset_documents 
			WHERE asset_id = $1 
			ORDER BY created_at DESC 
			LIMIT 1
		)`

	result, err := r.DB.ExecContext(ctx, query, assetID)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup duplicate documents: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// CleanupAllDuplicateDocuments removes all duplicate documents across all assets
func (r *assetDocumentRepository) CleanupAllDuplicateDocuments(ctx context.Context) (int, error) {
	// Query to find and delete all duplicates, keeping only the latest document per asset
	query := `
		DELETE FROM asset_documents 
		WHERE id NOT IN (
			SELECT DISTINCT ON (asset_id) id 
			FROM asset_documents 
			ORDER BY asset_id, created_at DESC
		)`

	result, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup all duplicate documents: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// GetDuplicateDocuments returns all duplicate documents grouped by asset_id
func (r *assetDocumentRepository) GetDuplicateDocuments(ctx context.Context) (map[uuid.UUID][]*entity.AssetDocument, error) {
	query := `
		SELECT id, tenant_id, asset_id, document_type, file_url, cloudinary_id, 
		       original_filename, file_size, mime_type, uploaded_at, created_at, updated_at
		FROM asset_documents 
		WHERE asset_id IN (
			SELECT asset_id 
			FROM asset_documents 
			GROUP BY asset_id 
			HAVING COUNT(*) > 1
		)
		ORDER BY asset_id, created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get duplicate documents: %w", err)
	}
	defer rows.Close()

	duplicates := make(map[uuid.UUID][]*entity.AssetDocument)

	for rows.Next() {
		var doc entity.AssetDocument
		var tenantID sql.NullString

		err := rows.Scan(
			&doc.ID,
			&tenantID,
			&doc.AssetID,
			&doc.DocumentType,
			&doc.FileURL,
			&doc.CloudinaryID,
			&doc.OriginalFilename,
			&doc.FileSize,
			&doc.MimeType,
			&doc.UploadedAt,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan duplicate document: %w", err)
		}

		// Handle nullable tenant_id
		if tenantID.Valid {
			if tenantUUID, err := uuid.Parse(tenantID.String); err == nil {
				doc.TenantID = &tenantUUID
			}
		}

		if doc.AssetID != nil {
			duplicates[*doc.AssetID] = append(duplicates[*doc.AssetID], &doc)
		}
	}

	return duplicates, nil
}
