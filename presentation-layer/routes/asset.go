package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupAssetRoutes configures all asset-related routes
func SetupAssetRoutes(router *gin.Engine, assetController *controller.AssetController) {
	// Group for asset routes
	assetGroup := router.Group("/api/v1/assets")
	{
		// Public routes (requires tenant validation from JWT)
		assetGroup.Use(middleware.TenantMiddleware())
		{
			// List all assets
			assetGroup.GET("", assetController.ListAssets)
			// Get asset by ID
			assetGroup.GET("/:id", assetController.GetAsset)
		}

		// SuperAdmin only routes - use SuperAdmin middleware for role validation
		superAdminGroup := router.Group("/api/v1/superadmin/assets")
		superAdminGroup.Use(middleware.TenantMiddleware())
		superAdminGroup.Use(middleware.RequireSuperAdminMiddleware())
		{
			// List all assets (for SuperAdmin - across all tenants)
			superAdminGroup.GET("", assetController.ListAllAssets)
			// Create new asset
			superAdminGroup.POST("", assetController.CreateAsset)
			// Update asset
			superAdminGroup.PUT("/:id", assetController.UpdateAsset)
			// Delete asset
			superAdminGroup.DELETE("/:id", assetController.DeleteAsset)
			// Assign asset to tenant
			superAdminGroup.POST("/:id/assign", assetController.AssignAssetToTenant)
			// Unassign asset from tenant
			superAdminGroup.POST("/:id/unassign", assetController.UnassignAssetFromTenant)
		}
	}
}
