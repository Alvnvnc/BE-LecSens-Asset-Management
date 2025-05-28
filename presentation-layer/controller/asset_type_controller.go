package controller

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/domain-layer/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AssetTypeController handles HTTP requests related to asset types
type AssetTypeController struct {
	assetTypeService *service.AssetTypeService
}

// NewAssetTypeController creates a new AssetTypeController
func NewAssetTypeController(assetTypeService *service.AssetTypeService) *AssetTypeController {
	return &AssetTypeController{
		assetTypeService: assetTypeService,
	}
}

// GetByID handles GET /asset-types/:id
func (c *AssetTypeController) GetByID(ctx *gin.Context) {
	// Get asset type ID from URL parameter
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset type ID format"})
		return
	}

	assetType, err := c.assetTypeService.GetAssetTypeByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Asset type not found"})
		return
	}

	ctx.JSON(http.StatusOK, assetType)
}

// List handles GET /asset-types
func (c *AssetTypeController) List(ctx *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	log.Printf("AssetTypeController: Starting List request - page: %d, pageSize: %d", page, pageSize)

	assetTypes, err := c.assetTypeService.ListAssetTypes(ctx.Request.Context(), page, pageSize)
	if err != nil {
		log.Printf("AssetTypeController: Error retrieving asset types: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve asset types", "details": err.Error()})
		return
	}

	log.Printf("AssetTypeController: Successfully retrieved %d asset types", len(assetTypes))

	ctx.JSON(http.StatusOK, gin.H{
		"data":     assetTypes,
		"page":     page,
		"pageSize": pageSize,
	})
}

// Create handles POST /asset-types
func (c *AssetTypeController) Create(ctx *gin.Context) {
	var assetType entity.AssetType

	if err := ctx.ShouldBindJSON(&assetType); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if assetType.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Ensure properties_schema has a default value if not provided or is empty
	if len(assetType.PropertiesSchema) == 0 {
		assetType.PropertiesSchema = []byte("{}")
	}

	err := c.assetTypeService.CreateAssetType(ctx.Request.Context(), &assetType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, assetType)
}

// Update handles PUT /asset-types/:id
func (c *AssetTypeController) Update(ctx *gin.Context) {
	// Get asset type ID from URL parameter
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset type ID format"})
		return
	}

	var assetType entity.AssetType
	if err := ctx.ShouldBindJSON(&assetType); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set ID from path parameter
	assetType.ID = id

	// Validate required fields
	if assetType.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Ensure properties_schema has a default value if not provided or is empty
	if len(assetType.PropertiesSchema) == 0 {
		assetType.PropertiesSchema = []byte("{}")
	}

	err = c.assetTypeService.UpdateAssetType(ctx.Request.Context(), &assetType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, assetType)
}

// Delete handles DELETE /asset-types/:id
func (c *AssetTypeController) Delete(ctx *gin.Context) {
	// Get asset type ID from URL parameter
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset type ID format"})
		return
	}

	err = c.assetTypeService.DeleteAssetType(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Asset type deleted successfully"})
}
