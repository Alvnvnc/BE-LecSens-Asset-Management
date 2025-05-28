package controller

import (
	"be-lecsens/asset_management/data-layer/config"
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AssetDocumentController handles HTTP requests for asset document operations
type AssetDocumentController struct {
	assetDocumentService *service.AssetDocumentService
	config               *config.Config
}

// NewAssetDocumentController creates a new AssetDocumentController
func NewAssetDocumentController(assetDocumentService *service.AssetDocumentService, cfg *config.Config) *AssetDocumentController {
	return &AssetDocumentController{
		assetDocumentService: assetDocumentService,
		config:               cfg,
	}
}

// CreateAssetDocument handles POST /api/v1/asset-documents
func (c *AssetDocumentController) CreateAssetDocument(ctx *gin.Context) {
	// Parse multipart form
	err := ctx.Request.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Failed to parse multipart form",
		})
		return
	}

	// Get file from form
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "File is required",
		})
		return
	}

	// Create DTO from form data
	req := &dto.CreateAssetDocumentRequest{
		DocumentType: ctx.PostForm("document_type"),
	}

	// Parse optional asset_id
	if assetIDStr := ctx.PostForm("asset_id"); assetIDStr != "" {
		assetID, err := uuid.Parse(assetIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid asset_id format",
			})
			return
		}
		req.AssetID = &assetID
	}

	// Validate required fields
	if req.DocumentType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "document_type is required",
		})
		return
	}

	// Call service
	document, err := c.assetDocumentService.CreateAssetDocument(ctx, req, fileHeader)
	if err != nil {
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create asset document",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Asset document created successfully",
		"data":    document,
	})
}

// GetAssetDocument handles GET /api/v1/asset-documents/:id
func (c *AssetDocumentController) GetAssetDocument(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid document ID format",
		})
		return
	}

	document, err := c.assetDocumentService.GetAssetDocument(ctx, id)
	if err != nil {
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve asset document",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset document retrieved successfully",
		"data":    document,
	})
}

// GetAssetDocuments handles GET /api/v1/superadmin/asset-documents/documents/:asset_id
func (c *AssetDocumentController) GetAssetDocuments(ctx *gin.Context) {
	assetIDParam := ctx.Param("asset_id")
	assetID, err := uuid.Parse(assetIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset ID format",
		})
		return
	}

	// Get optional document type filter
	docType := ctx.Query("document_type")

	var documents []*dto.AssetDocumentResponse
	if docType != "" {
		documents, err = c.assetDocumentService.GetAssetDocumentsByType(ctx, assetID, docType)
	} else {
		documents, err = c.assetDocumentService.GetAssetDocuments(ctx, assetID)
	}

	if err != nil {
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve asset documents",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset documents retrieved successfully",
		"data":    documents,
	})
}

// ListAssetDocuments handles GET /api/v1/asset-documents
func (c *AssetDocumentController) ListAssetDocuments(ctx *gin.Context) {
	// Parse pagination parameters
	page := 1
	pageSize := 10

	if pageParam := ctx.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := ctx.Query("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	documents, err := c.assetDocumentService.ListAssetDocuments(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to list asset documents",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset documents listed successfully",
		"data":    documents,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateAssetDocument handles PUT /api/v1/asset-documents/:id
func (c *AssetDocumentController) UpdateAssetDocument(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid document ID format",
		})
		return
	}

	// Parse request body or form data
	var req *dto.UpdateAssetDocumentRequest
	var fileHeader *multipart.FileHeader

	// Check if this is a multipart form (file upload)
	if ctx.GetHeader("Content-Type") != "" && strings.Contains(ctx.GetHeader("Content-Type"), "multipart/form-data") {
		// Handle file upload
		err := ctx.Request.ParseMultipartForm(32 << 20) // 32 MB
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Failed to parse multipart form",
			})
			return
		}

		req = &dto.UpdateAssetDocumentRequest{
			DocumentType: getStringPointer(ctx.PostForm("document_type")),
		}

		// Get optional file
		if fh, err := ctx.FormFile("file"); err == nil {
			fileHeader = fh
		}
	} else {
		// Handle JSON request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid request body",
			})
			return
		}
	}

	document, err := c.assetDocumentService.UpdateAssetDocument(ctx, id, req, fileHeader)
	if err != nil {
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to update asset document",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset document updated successfully",
		"data":    document,
	})
}

// DeleteAssetDocument handles DELETE /api/v1/asset-documents/:id
func (c *AssetDocumentController) DeleteAssetDocument(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid document ID format",
		})
		return
	}

	err = c.assetDocumentService.DeleteAssetDocument(ctx, id)
	if err != nil {
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete asset document",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset document deleted successfully",
	})
}

// DeleteAssetDocuments handles DELETE /api/v1/assets/:id/documents
func (c *AssetDocumentController) DeleteAssetDocuments(ctx *gin.Context) {
	assetIDParam := ctx.Param("id")
	assetID, err := uuid.Parse(assetIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset ID format",
		})
		return
	}

	err = c.assetDocumentService.DeleteAssetDocuments(ctx, assetID)
	if err != nil {
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete asset documents",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset documents deleted successfully",
	})
}

// Helper function to convert string to *string
func getStringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// CleanupDuplicateDocuments handles POST /api/v1/asset-documents/cleanup/{assetId}
func (c *AssetDocumentController) CleanupDuplicateDocuments(ctx *gin.Context) {
	// Get asset ID from URL parameter
	assetIDStr := ctx.Param("assetId")
	if assetIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Asset ID is required",
		})
		return
	}

	// Parse asset ID
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset ID format",
		})
		return
	}

	// Call service
	response, err := c.assetDocumentService.CleanupDuplicateDocuments(ctx, assetID)
	if err != nil {
		if validationErr, ok := err.(*common.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation Error",
				"message": validationErr.Message,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to cleanup duplicate documents",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Duplicate documents cleaned up successfully",
		"data":    response,
	})
}

// CleanupAllDuplicateDocuments handles POST /api/v1/asset-documents/cleanup-all
func (c *AssetDocumentController) CleanupAllDuplicateDocuments(ctx *gin.Context) {
	// Call service
	response, err := c.assetDocumentService.CleanupAllDuplicateDocuments(ctx)
	if err != nil {
		if validationErr, ok := err.(*common.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation Error",
				"message": validationErr.Message,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to cleanup all duplicate documents",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "All duplicate documents cleaned up successfully",
		"data":    response,
	})
}

// GetDuplicateDocuments handles GET /api/v1/asset-documents/duplicates
func (c *AssetDocumentController) GetDuplicateDocuments(ctx *gin.Context) {
	// Call service
	response, err := c.assetDocumentService.GetDuplicateDocuments(ctx)
	if err != nil {
		if validationErr, ok := err.(*common.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation Error",
				"message": validationErr.Message,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get duplicate documents",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Duplicate documents retrieved successfully",
		"data":    response,
	})
}

// ReplaceAssetDocument handles POST /api/v1/asset-documents/replace
func (c *AssetDocumentController) ReplaceAssetDocument(ctx *gin.Context) {
	// Parse multipart form
	err := ctx.Request.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Failed to parse multipart form",
		})
		return
	}

	// Get file from form
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "File is required",
		})
		return
	}

	// Parse request data from form
	req := &dto.CreateAssetDocumentRequest{
		DocumentType: ctx.PostForm("document_type"),
	}

	// Parse asset_id from form if provided
	if assetIDStr := ctx.PostForm("asset_id"); assetIDStr != "" {
		assetID, err := uuid.Parse(assetIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid asset_id format",
			})
			return
		}
		req.AssetID = &assetID
	}

	// Validate required fields
	if req.DocumentType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "document_type is required",
		})
		return
	}

	if req.AssetID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "asset_id is required for replacement",
		})
		return
	}

	// Call service to replace the document
	document, err := c.assetDocumentService.ReplaceAssetDocument(ctx, req, fileHeader)
	if err != nil {
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to replace asset document",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset document replaced successfully",
		"data":    document,
	})
}

// GetStorageInfo handles GET /api/v1/asset-documents/storage/:asset_id
func (c *AssetDocumentController) GetStorageInfo(ctx *gin.Context) {
	// Parse asset ID
	assetIDStr := ctx.Param("asset_id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset_id format",
		})
		return
	}

	// Get storage information
	storageInfo, err := c.assetDocumentService.GetStorageInfo(ctx, assetID)
	if err != nil {
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get storage information",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Storage information retrieved successfully",
		"data":    storageInfo,
	})
}
