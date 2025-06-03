package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupAssetAlertRoutes configures all asset alert-related routes
func SetupAssetAlertRoutes(router *gin.Engine, assetAlertController *controller.AssetAlertController) {
	// Group for asset alert routes
	assetAlertGroup := router.Group("/api/v1/asset-alerts")
	{
		// Public routes (requires tenant validation from JWT)
		assetAlertGroup.Use(middleware.TenantMiddleware())
		{
			// List alerts with pagination and filtering
			assetAlertGroup.GET("", assetAlertController.ListAssetAlerts)
			// Get alert by ID
			assetAlertGroup.GET("/:id", assetAlertController.GetAssetAlert)
			// Get alert statistics dashboard
			assetAlertGroup.GET("/statistics", assetAlertController.GetAlertStatistics)
		}

		// Admin routes - use TenantAdmin middleware for role validation
		adminGroup := assetAlertGroup.Group("")
		adminGroup.Use(middleware.TenantAdminMiddleware())
		{
			// Resolve single alert
			adminGroup.PATCH("/:id/resolve", assetAlertController.ResolveAssetAlert)
			// Resolve multiple alerts
			adminGroup.PATCH("/resolve-multiple", assetAlertController.ResolveMultipleAssetAlerts)
			// Delete alert
			adminGroup.DELETE("/:id", assetAlertController.DeleteAssetAlert)
			// Delete multiple alerts
			adminGroup.DELETE("/delete-multiple", assetAlertController.DeleteMultipleAssetAlerts)
		}

		// SuperAdmin only routes - use SuperAdmin middleware for role validation
		superAdminGroup := router.Group("/api/v1/superadmin/asset-alerts")
		superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// Get global alert statistics
			superAdminGroup.GET("/statistics", assetAlertController.GetGlobalAlertStatistics)
		}
	}
}
