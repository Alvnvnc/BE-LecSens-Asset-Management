package controller

import (
	"be-lecsens/asset_management/data-layer/config"
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TenantResponse represents the response from the tenant microservice
type TenantResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID           uuid.UUID `json:"user_id"`
		TenantID         uuid.UUID `json:"tenant_id"`
		TenantName       string    `json:"tenant_name"`
		UserRoleInTenant string    `json:"user_role_in_tenant"`
		IsActive         bool      `json:"is_active"`
		JoinedAt         string    `json:"joined_at"`
	} `json:"data"`
}

// UpdateAssetRequest represents the request body for partial asset updates
type UpdateAssetRequest struct {
	Name        *string          `json:"name,omitempty"`
	AssetTypeID *uuid.UUID       `json:"asset_type_id,omitempty"`
	LocationID  *uuid.UUID       `json:"location_id,omitempty"`
	Status      *string          `json:"status,omitempty"`
	Properties  *json.RawMessage `json:"properties,omitempty"`
	TenantID    *uuid.UUID       `json:"tenant_id,omitempty"` // Added for fixing assets without tenant
}

// AssetController handles HTTP requests for asset operations
type AssetController struct {
	assetService *service.AssetService
	config       *config.Config
}

// NewAssetController creates a new AssetController
func NewAssetController(assetService *service.AssetService, cfg *config.Config) *AssetController {
	return &AssetController{
		assetService: assetService,
		config:       cfg,
	}
}

// CreateAsset handles the creation of a new asset
func (c *AssetController) CreateAsset(ctx *gin.Context) {
	var asset entity.Asset
	if err := ctx.ShouldBindJSON(&asset); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate new UUID for the asset
	asset.ID = uuid.New()

	// Get tenant ID from context (for SuperAdmin, we need to set it explicitly)
	if tenantIDStr := ctx.GetString("tenant_id"); tenantIDStr != "" {
		if tenantID, err := uuid.Parse(tenantIDStr); err == nil {
			asset.TenantID = &tenantID
		}
	}

	// Create the asset
	if err := c.assetService.CreateAsset(ctx.Request.Context(), &asset); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, asset)
}

// GetAsset handles retrieving an asset by ID
func (c *AssetController) GetAsset(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	asset, err := c.assetService.GetAsset(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	ctx.JSON(http.StatusOK, asset)
}

// ListAssets handles retrieving a list of assets with pagination
func (c *AssetController) ListAssets(ctx *gin.Context) {
	// Get pagination parameters - convert to page-based pagination
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Ensure page is at least 1
	if page < 1 {
		page = 1
	}

	// Ensure limit is reasonable
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get tenant ID from context
	tenantID, _ := common.GetTenantID(ctx.Request.Context())

	assets, err := c.assetService.ListAssets(ctx.Request.Context(), &tenantID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list assets: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, assets)
}

// UpdateAsset handles updating an existing asset (partial update)
func (c *AssetController) UpdateAsset(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	var updateReq UpdateAssetRequest
	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get tenant ID from context for fixing assets without tenant assignment
	if contextTenantID := ctx.GetString("tenant_id"); contextTenantID != "" {
		// Get existing asset first to check if it needs tenant assignment
		existingAsset, err := c.assetService.GetAsset(ctx.Request.Context(), id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
			return
		}

		// If asset doesn't have tenant_id but we have one in context, assign it
		if existingAsset.TenantID == nil {
			if tenantUUID, parseErr := uuid.Parse(contextTenantID); parseErr == nil {
				// Add tenant_id to the update request
				updateReq.TenantID = &tenantUUID
			}
		}
	}

	// Update the asset with partial data
	updatedAsset, err := c.assetService.UpdateAssetPartial(ctx.Request.Context(), id, updateReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedAsset)
}

// DeleteAsset handles deleting an asset
func (c *AssetController) DeleteAsset(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	if err := c.assetService.DeleteAsset(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// AssignAssetToTenant handles assigning an asset to a tenant
func (c *AssetController) AssignAssetToTenant(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	var request dto.AssignAssetToTenantRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.assetService.AssignAssetToTenant(ctx.Request.Context(), id, request.TenantID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UnassignAssetFromTenant handles unassigning an asset from a tenant
func (c *AssetController) UnassignAssetFromTenant(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	if err := c.assetService.UnassignAssetFromTenant(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ListAllAssets handles retrieving all assets for SuperAdmin (across all tenants)
func (c *AssetController) ListAllAssets(ctx *gin.Context) {
	// Check if user has SuperAdmin role
	userRole, exists := ctx.Get("user_role")
	if !exists || userRole != "SUPERADMIN" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "SuperAdmin access required"})
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Ensure page is at least 1
	if page < 1 {
		page = 1
	}

	// Ensure limit is reasonable
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// For SuperAdmin, pass nil as tenantID to get assets from all tenants
	assets, err := c.assetService.ListAssets(ctx.Request.Context(), nil, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list assets: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, assets)
}
