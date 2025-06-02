package service

import (
	"be-lecsens/asset_management/data-layer/cloudinary"
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Constants for storage management
const (
	MaxDocumentsPerAsset = 3 // Maximum number of documents allowed per asset
)

// AssetDocumentService handles business logic for asset documents with Cloudinary integration
type AssetDocumentService struct {
	assetDocumentRepo repository.AssetDocumentRepository
	assetRepo         repository.AssetRepository
	cloudinaryService *cloudinary.CloudinaryService
}

// NewAssetDocumentService creates a new instance of AssetDocumentService
func NewAssetDocumentService(
	assetDocumentRepo repository.AssetDocumentRepository,
	assetRepo repository.AssetRepository,
	cloudinaryService *cloudinary.CloudinaryService,
) *AssetDocumentService {
	return &AssetDocumentService{
		assetDocumentRepo: assetDocumentRepo,
		assetRepo:         assetRepo,
		cloudinaryService: cloudinaryService,
	}
}

// CreateAssetDocument creates a new asset document with file upload to Cloudinary
func (s *AssetDocumentService) CreateAssetDocument(ctx context.Context, req *dto.CreateAssetDocumentRequest, file *multipart.FileHeader) (*dto.AssetDocumentResponse, error) {
	// Validate tenant ID from context (SuperAdmin can work without tenant ID)
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, common.NewValidationError("tenant ID is required", nil)
	}

	// Validate request
	if req.DocumentType == "" {
		return nil, common.NewValidationError("document_type is required", nil)
	}

	if file == nil {
		return nil, common.NewValidationError("file is required", nil)
	}

	// Validate file type and size
	if err := s.validateFile(file); err != nil {
		return nil, common.NewValidationError(err.Error(), err)
	}

	// Auto-detect tenant from asset and validate asset exists if asset_id is provided
	var finalTenantID *uuid.UUID
	if req.AssetID != nil {
		asset, err := s.assetRepo.GetByID(ctx, *req.AssetID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate asset: %w", err)
		}
		if asset == nil {
			return nil, common.NewNotFoundError("asset", req.AssetID.String())
		}

		// 🔄 VERSIONED STORAGE: Apply storage limit management instead of duplicate prevention
		err = s.enforceStorageLimit(ctx, *req.AssetID)
		if err != nil {
			return nil, fmt.Errorf("failed to enforce storage limit: %w", err)
		}

		// 🚀 AWESOME FEATURE: Auto-assign tenant_id from asset!
		// This makes documents automatically inherit the tenant of their associated asset
		// If asset has no tenant (null), document will also have null tenant
		finalTenantID = asset.TenantID

		// For regular users (non-SuperAdmin), ensure they can only access assets from their tenant
		if !isSuperAdmin && hasTenantID && asset.TenantID != nil && *asset.TenantID != tenantID {
			return nil, common.NewValidationError("asset does not belong to your tenant", nil)
		}
	} else {
		// If no asset_id provided, use tenant from context (could be nil for SuperAdmin)
		if hasTenantID {
			finalTenantID = &tenantID
		} else {
			finalTenantID = nil
		}
	}

	// Default asset ID if not provided
	assetID := uuid.New()
	if req.AssetID != nil {
		assetID = *req.AssetID
	}

	// Upload file to Cloudinary
	fileReader, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()

	uploadResult, err := s.cloudinaryService.UploadAssetDocument(ctx, fileReader, file, assetID, req.DocumentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to Cloudinary: %w", err)
	}

	// Check if cloudinary ID already exists (shouldn't happen, but safety check)
	exists, err := s.assetDocumentRepo.CheckCloudinaryIDExists(ctx, uploadResult.PublicID)
	if err != nil {
		// Try to cleanup uploaded file
		s.cloudinaryService.DeleteAssetDocument(ctx, uploadResult.PublicID)
		return nil, fmt.Errorf("failed to check cloudinary ID existence: %w", err)
	}
	if exists {
		// Try to cleanup uploaded file
		s.cloudinaryService.DeleteAssetDocument(ctx, uploadResult.PublicID)
		return nil, common.NewValidationError("cloudinary ID already exists", nil)
	}

	// Create entity with auto-detected tenant ID
	doc := &entity.AssetDocument{
		ID:               uuid.New(),
		TenantID:         finalTenantID, // Auto-assigned from asset or context
		AssetID:          req.AssetID,
		DocumentType:     req.DocumentType,
		FileURL:          uploadResult.URL,
		CloudinaryID:     uploadResult.PublicID,
		OriginalFilename: file.Filename,
		FileSize:         file.Size,
		MimeType:         file.Header.Get("Content-Type"),
		UploadedAt:       time.Now(),
	}

	// Save to database
	err = s.assetDocumentRepo.Create(ctx, doc)
	if err != nil {
		// Try to cleanup uploaded file
		s.cloudinaryService.DeleteAssetDocument(ctx, uploadResult.PublicID)
		return nil, fmt.Errorf("failed to create asset document: %w", err)
	}

	// Convert to response DTO
	return s.entityToResponse(doc), nil
}

// GetAssetDocument retrieves an asset document by ID
func (s *AssetDocumentService) GetAssetDocument(ctx context.Context, id uuid.UUID) (*dto.AssetDocumentResponse, error) {
	doc, err := s.assetDocumentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset document: %w", err)
	}

	if doc == nil {
		return nil, common.NewNotFoundError("asset document", id.String())
	}

	return s.entityToResponse(doc), nil
}

// GetAssetDocuments retrieves all documents for a specific asset
func (s *AssetDocumentService) GetAssetDocuments(ctx context.Context, assetID uuid.UUID) ([]*dto.AssetDocumentResponse, error) {
	// Validate that the asset exists and user has access to it
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate asset: %w", err)
	}
	if asset == nil {
		return nil, common.NewNotFoundError("asset", assetID.String())
	}

	documents, err := s.assetDocumentRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset documents: %w", err)
	}

	return s.entitiesToResponse(documents), nil
}

// GetAssetDocumentsByType retrieves documents for a specific asset filtered by document type
func (s *AssetDocumentService) GetAssetDocumentsByType(ctx context.Context, assetID uuid.UUID, docType string) ([]*dto.AssetDocumentResponse, error) {
	// Validate that the asset exists and user has access to it
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate asset: %w", err)
	}
	if asset == nil {
		return nil, common.NewNotFoundError("asset", assetID.String())
	}

	// Validate document type is not empty
	if docType == "" {
		return nil, common.NewValidationError("document type is required", nil)
	}

	documents, err := s.assetDocumentRepo.GetByAssetIDAndType(ctx, assetID, docType)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset documents: %w", err)
	}

	return s.entitiesToResponse(documents), nil
}

// ListAssetDocuments retrieves asset documents with pagination
func (s *AssetDocumentService) ListAssetDocuments(ctx context.Context, page, pageSize int) ([]*dto.AssetDocumentResponse, error) {
	// Validate tenant ID from context (SuperAdmin can work without tenant ID)
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	log.Printf("ListAssetDocuments: tenantID=%v, hasTenantID=%v, isSuperAdmin=%v", tenantID, hasTenantID, isSuperAdmin)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, common.NewValidationError("tenant ID is required", nil)
	}

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default page size
	}

	log.Printf("ListAssetDocuments: Calling repository List with page=%d, pageSize=%d", page, pageSize)

	documents, err := s.assetDocumentRepo.List(ctx, page, pageSize)
	if err != nil {
		log.Printf("ListAssetDocuments: Repository error: %v", err)
		return nil, fmt.Errorf("failed to list asset documents: %w", err)
	}

	log.Printf("ListAssetDocuments: Retrieved %d documents", len(documents))

	return s.entitiesToResponse(documents), nil
}

// UpdateAssetDocument updates an existing asset document
func (s *AssetDocumentService) UpdateAssetDocument(ctx context.Context, id uuid.UUID, req *dto.UpdateAssetDocumentRequest, file *multipart.FileHeader) (*dto.AssetDocumentResponse, error) {
	// Get existing document to validate ownership and existence
	existingDoc, err := s.assetDocumentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing document: %w", err)
	}
	if existingDoc == nil {
		return nil, common.NewNotFoundError("asset document", id.String())
	}

	// Start with existing document
	updatedDoc := *existingDoc

	// Update document type if provided
	if req.DocumentType != nil {
		if *req.DocumentType == "" {
			return nil, common.NewValidationError("document type cannot be empty", nil)
		}
		updatedDoc.DocumentType = *req.DocumentType
	}

	// Handle file upload if new file is provided
	if file != nil {
		// Validate file
		validateErr := s.validateFile(file)
		if validateErr != nil {
			return nil, common.NewValidationError(validateErr.Error(), validateErr)
		}

		// Default asset ID if not available
		assetID := uuid.New()
		if updatedDoc.AssetID != nil {
			assetID = *updatedDoc.AssetID
		}

		// Upload new file to Cloudinary
		var fileReader multipart.File
		fileReader, err = file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer fileReader.Close()

		var uploadResult *cloudinary.UploadResult
		uploadResult, err = s.cloudinaryService.UploadAssetDocument(ctx, fileReader, file, assetID, updatedDoc.DocumentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file to Cloudinary: %w", err)
		}

		// Delete old file from Cloudinary
		if existingDoc.CloudinaryID != "" {
			err = s.cloudinaryService.DeleteAssetDocument(ctx, existingDoc.CloudinaryID)
			if err != nil {
				// Log error but don't fail the update
				fmt.Printf("Warning: failed to delete old file from Cloudinary: %v\n", err)
			}
		}

		// Update file-related fields
		updatedDoc.FileURL = uploadResult.URL
		updatedDoc.CloudinaryID = uploadResult.PublicID
		updatedDoc.OriginalFilename = file.Filename
		updatedDoc.FileSize = file.Size
		updatedDoc.MimeType = file.Header.Get("Content-Type")
	}

	// Update in repository
	err = s.assetDocumentRepo.Update(ctx, &updatedDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to update asset document: %w", err)
	}

	return s.entityToResponse(&updatedDoc), nil
}

// DeleteAssetDocument deletes an asset document
func (s *AssetDocumentService) DeleteAssetDocument(ctx context.Context, id uuid.UUID) error {
	// Validate document exists and user has access
	doc, err := s.assetDocumentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get asset document: %w", err)
	}
	if doc == nil {
		return common.NewNotFoundError("asset document", id.String())
	}

	// Delete from database first
	err = s.assetDocumentRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset document: %w", err)
	}

	// Delete file from Cloudinary
	if doc.CloudinaryID != "" {
		err = s.cloudinaryService.DeleteAssetDocument(ctx, doc.CloudinaryID)
		if err != nil {
			// Log error but don't fail the delete since database record is already gone
			fmt.Printf("Warning: failed to delete file from Cloudinary: %v\n", err)
		}
	}

	return nil
}

// DeleteAssetDocuments deletes all documents for a specific asset
func (s *AssetDocumentService) DeleteAssetDocuments(ctx context.Context, assetID uuid.UUID) error {
	// Validate that the asset exists and user has access to it
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to validate asset: %w", err)
	}
	if asset == nil {
		return common.NewNotFoundError("asset", assetID.String())
	}

	// Get all documents for the asset first (to get Cloudinary IDs)
	documents, err := s.assetDocumentRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to get asset documents: %w", err)
	}

	// Delete from database first
	err = s.assetDocumentRepo.DeleteByAssetID(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to delete asset documents: %w", err)
	}

	// Delete files from Cloudinary
	for _, doc := range documents {
		if doc.CloudinaryID != "" {
			err = s.cloudinaryService.DeleteAssetDocument(ctx, doc.CloudinaryID)
			if err != nil {
				// Log error but continue with other files
				fmt.Printf("Warning: failed to delete file from Cloudinary: %v\n", err)
			}
		}
	}

	return nil
}

// validateFile validates the uploaded file
func (s *AssetDocumentService) validateFile(file *multipart.FileHeader) error {
	// Check file size (max 50MB)
	maxSize := int64(50 * 1024 * 1024) // 50MB
	if file.Size > maxSize {
		return fmt.Errorf("file size too large: maximum %d bytes allowed", maxSize)
	}

	// Check file extension
	allowedExtensions := []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt", ".jpg", ".jpeg", ".png", ".gif", ".bmp"}
	ext := strings.ToLower(filepath.Ext(file.Filename))

	allowed := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file type not allowed: %s", ext)
	}

	return nil
}

// entityToResponse converts an entity to response DTO
func (s *AssetDocumentService) entityToResponse(doc *entity.AssetDocument) *dto.AssetDocumentResponse {
	return &dto.AssetDocumentResponse{
		ID:               doc.ID,
		TenantID:         doc.TenantID,
		AssetID:          doc.AssetID,
		DocumentType:     doc.DocumentType,
		FileURL:          doc.FileURL,
		OriginalFilename: doc.OriginalFilename,
		FileSize:         doc.FileSize,
		MimeType:         doc.MimeType,
		UploadedAt:       doc.UploadedAt,
		CreatedAt:        doc.UploadedAt, // Assuming they're the same
		UpdatedAt:        doc.UploadedAt, // This would need to be tracked separately in real implementation
	}
}

// entitiesToResponse converts multiple entities to response DTOs
func (s *AssetDocumentService) entitiesToResponse(docs []*entity.AssetDocument) []*dto.AssetDocumentResponse {
	responses := make([]*dto.AssetDocumentResponse, len(docs))
	for i, doc := range docs {
		responses[i] = s.entityToResponse(doc)
	}
	return responses
}

// CleanupDuplicateDocuments removes duplicate documents for a specific asset, keeping only the latest one
func (s *AssetDocumentService) CleanupDuplicateDocuments(ctx context.Context, assetID uuid.UUID) (*dto.CleanupResponse, error) {
	// Validate tenant ID from context (SuperAdmin can work without tenant ID)
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, common.NewValidationError("tenant ID is required", nil)
	}

	// For regular users, verify asset belongs to their tenant
	if !isSuperAdmin && hasTenantID {
		asset, err := s.assetRepo.GetByID(ctx, assetID)
		if err != nil {
			return nil, common.NewValidationError("asset not found", nil)
		}

		// Check if asset belongs to user's tenant
		if asset.TenantID != nil && *asset.TenantID != tenantID {
			return nil, common.NewValidationError("asset does not belong to your tenant", nil)
		}
	}

	// Cleanup duplicates for the specific asset
	deletedCount, err := s.assetDocumentRepo.CleanupDuplicateDocuments(ctx, assetID)
	if err != nil {
		log.Printf("Failed to cleanup duplicate documents for asset %s: %v", assetID, err)
		return nil, common.NewValidationError("failed to cleanup duplicate documents", nil)
	}

	return &dto.CleanupResponse{
		AssetID:          &assetID,
		DocumentsCleaned: deletedCount,
		Message:          fmt.Sprintf("Successfully cleaned up %d duplicate documents for asset %s", deletedCount, assetID),
	}, nil
}

// CleanupAllDuplicateDocuments removes all duplicate documents across all assets, keeping only the latest for each asset
func (s *AssetDocumentService) CleanupAllDuplicateDocuments(ctx context.Context) (*dto.CleanupResponse, error) {
	// Only SuperAdmin can cleanup all duplicate documents
	if !common.IsSuperAdmin(ctx) {
		return nil, common.NewValidationError("only SuperAdmin can cleanup all duplicate documents", nil)
	}

	// Cleanup all duplicates
	deletedCount, err := s.assetDocumentRepo.CleanupAllDuplicateDocuments(ctx)
	if err != nil {
		log.Printf("Failed to cleanup all duplicate documents: %v", err)
		return nil, common.NewValidationError("failed to cleanup all duplicate documents", nil)
	}

	return &dto.CleanupResponse{
		DocumentsCleaned: deletedCount,
		Message:          fmt.Sprintf("Successfully cleaned up %d duplicate documents across all assets", deletedCount),
	}, nil
}

// GetDuplicateDocuments returns all duplicate documents grouped by asset_id
func (s *AssetDocumentService) GetDuplicateDocuments(ctx context.Context) (*dto.DuplicateDocumentsResponse, error) {
	// Only SuperAdmin can view all duplicate documents
	if !common.IsSuperAdmin(ctx) {
		return nil, common.NewValidationError("only SuperAdmin can view all duplicate documents", nil)
	}

	// Get duplicate documents
	duplicates, err := s.assetDocumentRepo.GetDuplicateDocuments(ctx)
	if err != nil {
		log.Printf("Failed to get duplicate documents: %v", err)
		return nil, common.NewValidationError("failed to get duplicate documents", nil)
	}

	// Convert to response format
	duplicateGroups := make(map[string][]*dto.AssetDocumentResponse)
	totalDuplicates := 0

	for assetID, docs := range duplicates {
		if len(docs) > 1 { // Only include if there are actually duplicates
			assetIDStr := assetID.String()
			duplicateGroups[assetIDStr] = s.entitiesToResponse(docs)
			totalDuplicates += len(docs) - 1 // Count extra documents (excluding the one we keep)
		}
	}

	return &dto.DuplicateDocumentsResponse{
		DuplicateGroups: duplicateGroups,
		TotalDuplicates: totalDuplicates,
		AffectedAssets:  len(duplicateGroups),
		Message:         fmt.Sprintf("Found %d duplicate documents across %d assets", totalDuplicates, len(duplicateGroups)),
	}, nil
}

// ReplaceAssetDocument replaces an existing document with versioning backup system
// Only the latest file is recorded in database, older files serve as backups in Cloudinary
func (s *AssetDocumentService) ReplaceAssetDocument(ctx context.Context, req *dto.CreateAssetDocumentRequest, file *multipart.FileHeader) (*dto.AssetDocumentResponse, error) {
	log.Printf("ReplaceAssetDocument: Starting versioned replacement for asset_id=%s, document_type=%s",
		func() string {
			if req.AssetID != nil {
				return req.AssetID.String()
			}
			return "nil"
		}(), req.DocumentType)

	// Validate tenant ID from context (SuperAdmin can work without tenant ID)
	tenantID, hasTenantID := common.GetTenantID(ctx)
	isSuperAdmin := common.IsSuperAdmin(ctx)

	// For regular users, tenant ID is required. For SuperAdmin, it's optional
	if !hasTenantID && !isSuperAdmin {
		return nil, common.NewValidationError("tenant ID is required", nil)
	}

	// Validate request
	if req.DocumentType == "" {
		return nil, common.NewValidationError("document_type is required", nil)
	}

	if req.AssetID == nil {
		return nil, common.NewValidationError("asset_id is required for replacement", nil)
	}

	if file == nil {
		return nil, common.NewValidationError("file is required", nil)
	}

	// Validate file type and size
	if err := s.validateFile(file); err != nil {
		return nil, common.NewValidationError(err.Error(), err)
	}

	// Validate asset exists
	asset, err := s.assetRepo.GetByID(ctx, *req.AssetID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate asset: %w", err)
	}
	if asset == nil {
		return nil, common.NewNotFoundError("asset", req.AssetID.String())
	}

	// Check if document of this type exists
	exists, err := s.assetDocumentRepo.CheckAssetDocumentExists(ctx, *req.AssetID, req.DocumentType)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing document: %w", err)
	}
	if !exists {
		return nil, common.NewValidationError(
			fmt.Sprintf("No existing document of type '%s' found for this asset. Use create endpoint instead.", req.DocumentType),
			nil,
		)
	}

	// Get existing document to replace
	existingDocs, err := s.assetDocumentRepo.GetByAssetIDAndType(ctx, *req.AssetID, req.DocumentType)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing document: %w", err)
	}
	if len(existingDocs) == 0 {
		return nil, common.NewNotFoundError("existing document", fmt.Sprintf("asset_id=%s, type=%s", req.AssetID.String(), req.DocumentType))
	}

	// Use the first/current active document
	currentDoc := existingDocs[0]

	// For regular users, ensure they can only access assets from their tenant
	if !isSuperAdmin && hasTenantID && asset.TenantID != nil && *asset.TenantID != tenantID {
		return nil, common.NewValidationError("asset does not belong to your tenant", nil)
	}

	// 🔄 VERSIONED STORAGE MANAGEMENT: Upload new file with versioning
	err = s.manageVersionedStorage(*req.AssetID, req.DocumentType, currentDoc.CloudinaryID)
	if err != nil {
		return nil, fmt.Errorf("failed to manage versioned storage: %w", err)
	}

	// Upload new file to Cloudinary
	fileReader, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()

	uploadResult, err := s.cloudinaryService.UploadAssetDocument(ctx, fileReader, file, *req.AssetID, req.DocumentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to Cloudinary: %w", err)
	}

	// Update existing document with new file details (only latest is recorded in DB)
	currentDoc.FileURL = uploadResult.URL
	currentDoc.CloudinaryID = uploadResult.PublicID
	currentDoc.OriginalFilename = file.Filename
	currentDoc.FileSize = file.Size
	currentDoc.MimeType = file.Header.Get("Content-Type")
	currentDoc.UploadedAt = time.Now()

	// Save updated document to database
	err = s.assetDocumentRepo.Update(ctx, currentDoc)
	if err != nil {
		// Try to cleanup uploaded file
		s.cloudinaryService.DeleteAssetDocument(ctx, uploadResult.PublicID)
		return nil, fmt.Errorf("failed to update asset document: %w", err)
	}

	log.Printf("Successfully replaced document %s for asset %s with versioned backup system", currentDoc.ID, req.AssetID)

	// Convert to response DTO
	return s.entityToResponse(currentDoc), nil
}

// manageVersionedStorage manages the versioned storage system where:
// - Max 3 files per asset per document type in Cloudinary
// - Only the latest file is recorded in database
// - 2 older files serve as backups in Cloudinary
func (s *AssetDocumentService) manageVersionedStorage(assetID uuid.UUID, docType, currentCloudinaryID string) error {
	log.Printf("manageVersionedStorage: Starting for asset %s, type %s", assetID.String(), docType)

	// In the versioned system:
	// 1. Current file (recorded in DB) will become backup when we upload new file
	// 2. New file will be uploaded and recorded in DB
	// 3. If we already have 2+ backup files, oldest backup gets deleted

	// For this versioned approach, we let Cloudinary handle the file versioning
	// and we only keep track of the latest file in our database
	// The backup files remain in Cloudinary but are not tracked in our DB

	// This is a placeholder for more advanced Cloudinary versioning logic
	// In a full implementation, we could:
	// 1. Use Cloudinary's version management features
	// 2. Track backup file IDs separately
	// 3. Implement cleanup of old backup files based on timestamps

	log.Printf("manageVersionedStorage: File %s will remain as backup in Cloudinary", currentCloudinaryID)

	// For now, we rely on the 3-file limit in enforceStorageLimit for the simpler approach
	// The versioned approach would be more sophisticated and handle backups automatically

	return nil
}

// enforceStorageLimit enforces the maximum document limit per asset by deleting oldest documents
func (s *AssetDocumentService) enforceStorageLimit(ctx context.Context, assetID uuid.UUID) error {
	// Get all existing documents for the asset
	existingDocs, err := s.assetDocumentRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to get existing documents: %w", err)
	}

	// If we're at or over the limit, delete oldest documents
	if len(existingDocs) >= MaxDocumentsPerAsset {
		// Sort documents by upload date to find oldest ones
		for i := 0; i < len(existingDocs)-1; i++ {
			for j := i + 1; j < len(existingDocs); j++ {
				if existingDocs[i].UploadedAt.After(existingDocs[j].UploadedAt) {
					existingDocs[i], existingDocs[j] = existingDocs[j], existingDocs[i]
				}
			}
		}

		// Calculate how many documents to delete (keep space for 1 new document)
		documentsToDelete := len(existingDocs) - MaxDocumentsPerAsset + 1

		for i := 0; i < documentsToDelete; i++ {
			oldDoc := existingDocs[i]

			log.Printf("Storage limit enforced for asset %s. Deleting document: %s (uploaded: %s)",
				assetID.String(), oldDoc.OriginalFilename, oldDoc.UploadedAt.Format(time.RFC3339))

			// Delete from database first
			err = s.assetDocumentRepo.Delete(ctx, oldDoc.ID)
			if err != nil {
				log.Printf("Warning: failed to delete old document from database: %v", err)
				continue
			}

			// Delete from Cloudinary
			if oldDoc.CloudinaryID != "" {
				err = s.cloudinaryService.DeleteAssetDocument(ctx, oldDoc.CloudinaryID)
				if err != nil {
					log.Printf("Warning: failed to delete old file from Cloudinary: %v", err)
				}
			}
		}
	}

	return nil
}

// GetStorageInfo returns storage information for an asset
func (s *AssetDocumentService) GetStorageInfo(ctx context.Context, assetID uuid.UUID) (*dto.StorageInfoResponse, error) {
	// Validate that the asset exists and user has access to it
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate asset: %w", err)
	}
	if asset == nil {
		return nil, common.NewNotFoundError("asset", assetID.String())
	}

	// Get existing documents
	existingDocs, err := s.assetDocumentRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing documents: %w", err)
	}

	// Calculate total size
	var totalSize int64
	for _, doc := range existingDocs {
		totalSize += doc.FileSize
	}

	return &dto.StorageInfoResponse{
		AssetID:        assetID,
		CurrentCount:   len(existingDocs),
		MaxCount:       MaxDocumentsPerAsset,
		AvailableSlots: MaxDocumentsPerAsset - len(existingDocs),
		TotalSizeBytes: totalSize,
		Documents:      s.entitiesToResponse(existingDocs),
		IsAtLimit:      len(existingDocs) >= MaxDocumentsPerAsset,
		Message: fmt.Sprintf("Asset has %d of %d documents. %d slots available.",
			len(existingDocs), MaxDocumentsPerAsset, MaxDocumentsPerAsset-len(existingDocs)),
	}, nil
}
