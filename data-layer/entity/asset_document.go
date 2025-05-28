package entity

import (
	"time"

	"github.com/google/uuid"
)

// AssetDocument represents files or documents associated with an asset
type AssetDocument struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         *uuid.UUID `json:"tenant_id"` // Denormalized from Asset for performance and security, nullable for flexibility
	AssetID          *uuid.UUID `json:"asset_id"`  // Nullable foreign key for flexibility
	DocumentType     string     `json:"document_type"`
	FileURL          string     `json:"file_url"`          // Cloudinary URL
	CloudinaryID     string     `json:"cloudinary_id"`     // Cloudinary public ID for deletion
	OriginalFilename string     `json:"original_filename"` // Original filename when uploaded
	FileSize         int64      `json:"file_size"`         // File size in bytes
	MimeType         string     `json:"mime_type"`         // MIME type of the file
	UploadedAt       time.Time  `json:"uploaded_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
