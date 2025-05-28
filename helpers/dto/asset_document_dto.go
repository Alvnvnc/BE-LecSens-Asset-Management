package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateAssetDocumentRequest represents the request to create a new asset document
type CreateAssetDocumentRequest struct {
	AssetID      *uuid.UUID `form:"asset_id" validate:"omitempty,uuid"`
	DocumentType string     `form:"document_type" binding:"required" validate:"required,max=100"`
	// File will be handled separately in multipart form

	// 🚀 AWESOME FEATURES:
	// 1. AUTO-TENANT DETECTION:
	//    When asset_id is provided, the document will automatically inherit
	//    the tenant_id from the associated asset. This means:
	//    - If asset belongs to tenant A, document will belong to tenant A
	//    - If asset has no tenant (null), document will also have null tenant
	//    - This makes tenant management seamless and automatic!
	//
	// 2. AUTO-REPLACEMENT (ONE DOCUMENT PER ASSET):
	//    Each asset can only have ONE document at a time. When uploading a new document:
	//    - If asset already has a document, the OLD document will be AUTOMATICALLY DELETED
	//    - Both database record and Cloudinary file will be removed
	//    - The NEW document becomes the current document for that asset
	//    - This ensures clean data management and prevents file accumulation
	//
	// 3. SUPERADMIN GLOBAL ACCESS:
	//    SuperAdmin users can manage documents across all tenants without restrictions
	//
	// 2. DUPLICATE PREVENTION:
	//    Each asset can only have ONE document per document_type.
	//    If you try to upload a document with the same type for the same asset,
	//    the system will reject it with an error message.
	//    This ensures data integrity and prevents confusion.
}

// UpdateAssetDocumentRequest represents the request body for partial asset document updates
type UpdateAssetDocumentRequest struct {
	DocumentType *string `json:"document_type,omitempty" validate:"omitempty,max=100"`
	// Note: File URL and Cloudinary ID should not be directly updatable by users
}

// AssetDocumentResponse represents the response structure for asset document operations
type AssetDocumentResponse struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         *uuid.UUID `json:"tenant_id"`
	AssetID          *uuid.UUID `json:"asset_id"`
	DocumentType     string     `json:"document_type"`
	FileURL          string     `json:"file_url"`
	OriginalFilename string     `json:"original_filename"`
	FileSize         int64      `json:"file_size"`
	MimeType         string     `json:"mime_type"`
	UploadedAt       time.Time  `json:"uploaded_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// AssetDocumentListResponse represents the response for listing asset documents with pagination
type AssetDocumentListResponse struct {
	Documents  []AssetDocumentResponse `json:"documents"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	Total      int64                   `json:"total"`
	TotalPages int                     `json:"total_pages"`
}

// AssetDocumentsByTypeRequest represents the request to filter documents by type
type AssetDocumentsByTypeRequest struct {
	AssetID      *uuid.UUID `json:"asset_id" validate:"omitempty,uuid"`
	DocumentType string     `json:"document_type" binding:"required" validate:"required,max=100"`
}

// AssetDocumentUploadRequest represents the request for document upload with metadata
type AssetDocumentUploadRequest struct {
	AssetID      *uuid.UUID `form:"asset_id" validate:"omitempty,uuid"`
	DocumentType string     `form:"document_type" binding:"required" validate:"required,max=100"`
	// File will be handled separately in multipart form
}

// AssetDocumentDeleteRequest represents the request to delete a document
type AssetDocumentDeleteRequest struct {
	ID uuid.UUID `json:"id" binding:"required" validate:"required"`
}

// BulkDeleteAssetDocumentsRequest represents the request to delete multiple documents
type BulkDeleteAssetDocumentsRequest struct {
	DocumentIDs []uuid.UUID `json:"document_ids" binding:"required,min=1" validate:"required,min=1"`
}

// AssetDocumentExistsResponse represents the response for checking if a cloudinary ID exists
type AssetDocumentExistsResponse struct {
	Exists       bool   `json:"exists"`
	CloudinaryID string `json:"cloudinary_id"`
}

// AssetDocumentStatsResponse represents statistics about documents for an asset
type AssetDocumentStatsResponse struct {
	AssetID       *uuid.UUID                       `json:"asset_id"`
	TotalCount    int                              `json:"total_count"`
	TypeCounts    map[string]int                   `json:"type_counts"`
	LastUploaded  *time.Time                       `json:"last_uploaded,omitempty"`
	FirstUploaded *time.Time                       `json:"first_uploaded,omitempty"`
	DocumentTypes []AssetDocumentTypeStatsResponse `json:"document_types"`
	TotalSize     int64                            `json:"total_size_bytes"`
}

// AssetDocumentTypeStatsResponse represents statistics for a specific document type
type AssetDocumentTypeStatsResponse struct {
	DocumentType string    `json:"document_type"`
	Count        int       `json:"count"`
	LastUploaded time.Time `json:"last_uploaded"`
}

// CleanupResponse represents the response for cleanup operations
type CleanupResponse struct {
	AssetID          *uuid.UUID `json:"asset_id,omitempty"`
	DocumentsCleaned int        `json:"documents_cleaned"`
	Message          string     `json:"message"`
}

// DuplicateDocumentsResponse represents the response for getting duplicate documents
type DuplicateDocumentsResponse struct {
	DuplicateGroups map[string][]*AssetDocumentResponse `json:"duplicate_groups"`
	TotalDuplicates int                                 `json:"total_duplicates"`
	AffectedAssets  int                                 `json:"affected_assets"`
	Message         string                              `json:"message"`
}

// StorageInfoResponse represents storage information for an asset
type StorageInfoResponse struct {
	AssetID        uuid.UUID                `json:"asset_id"`
	CurrentCount   int                      `json:"current_count"`
	MaxCount       int                      `json:"max_count"`
	AvailableSlots int                      `json:"available_slots"`
	TotalSizeBytes int64                    `json:"total_size_bytes"`
	Documents      []*AssetDocumentResponse `json:"documents"`
	IsAtLimit      bool                     `json:"is_at_limit"`
	Message        string                   `json:"message"`
}
