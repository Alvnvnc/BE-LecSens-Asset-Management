package cloudinary

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

// CloudinaryConfig holds Cloudinary configuration
type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

// CloudinaryService handles file operations with Cloudinary
type CloudinaryService struct {
	client *cloudinary.Cloudinary
	config *CloudinaryConfig
}

// NewCloudinaryService creates a new CloudinaryService
func NewCloudinaryService(config *CloudinaryConfig) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(config.CloudName, config.APIKey, config.APISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &CloudinaryService{
		client: cld,
		config: config,
	}, nil
}

// UploadAssetDocument uploads a document file for an asset
func (s *CloudinaryService) UploadAssetDocument(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, assetID uuid.UUID, documentType string) (*UploadResult, error) {
	// Generate unique filename
	filename := s.generateDocumentFilename(fileHeader.Filename, documentType)

	// Create folder path: document_asset/{assetID}/
	folderPath := fmt.Sprintf("document_asset/%s", assetID.String())

	// Upload parameters
	uniqueFilename := false
	overwrite := false
	uploadParams := uploader.UploadParams{
		PublicID:       filename,
		Folder:         folderPath,
		ResourceType:   "auto",          // Auto-detect file type
		UniqueFilename: &uniqueFilename, // Use our custom filename
		Overwrite:      &overwrite,      // Don't overwrite existing files
		Tags:           []string{"asset_document", documentType, assetID.String()},
		Context: map[string]string{
			"asset_id":      assetID.String(),
			"document_type": documentType,
			"uploaded_at":   time.Now().Format(time.RFC3339),
		},
	}

	// Upload file
	result, err := s.client.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to Cloudinary: %w", err)
	}

	return &UploadResult{
		PublicID:         result.PublicID,
		URL:              result.SecureURL,
		Format:           result.Format,
		Version:          result.Version,
		Width:            result.Width,
		Height:           result.Height,
		Bytes:            result.Bytes,
		AssetID:          assetID,
		DocumentType:     documentType,
		OriginalFilename: fileHeader.Filename,
	}, nil
}

// DeleteAssetDocument deletes a document from Cloudinary
func (s *CloudinaryService) DeleteAssetDocument(ctx context.Context, publicID string) error {
	invalidate := true
	_, err := s.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "auto",
		Invalidate:   &invalidate,
	})

	if err != nil {
		return fmt.Errorf("failed to delete file from Cloudinary: %w", err)
	}

	return nil
}

// DeleteAssetDocumentFolder deletes entire folder for an asset
func (s *CloudinaryService) DeleteAssetDocumentFolder(ctx context.Context, assetID uuid.UUID) error {
	folderPath := fmt.Sprintf("document_asset/%s", assetID.String())

	invalidate := true
	// Delete all files in the folder first
	_, err := s.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     folderPath + "/*",
		ResourceType: "auto",
		Invalidate:   &invalidate,
	})

	if err != nil {
		return fmt.Errorf("failed to delete asset document folder: %w", err)
	}

	return nil
}

// GetAssetDocumentURL generates a URL for accessing the document
func (s *CloudinaryService) GetAssetDocumentURL(publicID string, options ...URLOption) string {
	imageAsset, err := s.client.Image(publicID)
	if err != nil {
		return ""
	}

	url, err := imageAsset.String()
	if err != nil {
		return ""
	}

	// Apply any URL transformations if provided
	for _, option := range options {
		url = option(url)
	}

	return url
}

// UploadResult represents the result of a file upload
type UploadResult struct {
	PublicID         string    `json:"public_id"`
	URL              string    `json:"url"`
	Format           string    `json:"format"`
	Version          int       `json:"version"`
	Width            int       `json:"width,omitempty"`
	Height           int       `json:"height,omitempty"`
	Bytes            int       `json:"bytes"`
	AssetID          uuid.UUID `json:"asset_id"`
	DocumentType     string    `json:"document_type"`
	OriginalFilename string    `json:"original_filename"`
}

// URLOption is a function type for URL modifications
type URLOption func(string) string

// generateDocumentFilename creates a unique filename for the document
func (s *CloudinaryService) generateDocumentFilename(originalFilename, documentType string) string {
	// Get file extension
	ext := filepath.Ext(originalFilename)
	// Remove extension from filename
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)
	// Clean filename (remove special characters)
	cleanName := strings.ReplaceAll(nameWithoutExt, " ", "_")
	cleanName = strings.ReplaceAll(cleanName, "(", "")
	cleanName = strings.ReplaceAll(cleanName, ")", "")

	// Generate timestamp
	timestamp := time.Now().Format("20060102_150405")

	// Create unique filename: documentType_originalName_timestamp
	filename := fmt.Sprintf("%s_%s_%s%s", documentType, cleanName, timestamp, ext)

	return filename
}

// ValidateFileType validates if the uploaded file type is allowed
func ValidateFileType(fileHeader *multipart.FileHeader, allowedTypes []string) error {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	// Check if extension is allowed
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return nil
		}
	}

	return fmt.Errorf("file type %s is not allowed", ext)
}

// ValidateFileSize validates if the file size is within limits
func ValidateFileSize(fileHeader *multipart.FileHeader, maxSizeBytes int64) error {
	if fileHeader.Size > maxSizeBytes {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", fileHeader.Size, maxSizeBytes)
	}
	return nil
}

// GetFileContentType returns the content type of the file
func GetFileContentType(file multipart.File) (string, error) {
	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Detect content type
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
