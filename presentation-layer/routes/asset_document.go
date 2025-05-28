package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupAssetDocumentRoutes configures all asset document-related routes
func SetupAssetDocumentRoutes(router *gin.Engine, assetDocumentController *controller.AssetDocumentController) {
	// Group for asset document routes
	assetDocumentGroup := router.Group("/api/v1/asset-documents")
	{
		// Public routes (requires tenant validation from JWT) - READ ONLY for regular users
		assetDocumentGroup.Use(middleware.TenantMiddleware())
		{
			// List all asset documents (paginated)
			assetDocumentGroup.GET("", assetDocumentController.ListAssetDocuments)
			// Get asset document by ID
			assetDocumentGroup.GET("/:id", assetDocumentController.GetAssetDocument)
		}
	}

	// SuperAdmin only routes - use SuperAdmin middleware for role validation
	superAdminGroup := router.Group("/api/v1/superadmin/asset-documents")
	superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		// List all asset documents (for SuperAdmin - across all tenants)
		superAdminGroup.GET("", assetDocumentController.ListAssetDocuments)
		// Create new asset document (file upload)
		superAdminGroup.POST("", assetDocumentController.CreateAssetDocument)
		// Replace existing asset document (file upload)
		superAdminGroup.POST("/replace", assetDocumentController.ReplaceAssetDocument)
		// Update asset document
		superAdminGroup.PUT("/:id", assetDocumentController.UpdateAssetDocument)
		// Delete asset document
		superAdminGroup.DELETE("/:id", assetDocumentController.DeleteAssetDocument)

		// Duplicate cleanup routes (SuperAdmin only)
		// Get all duplicate documents
		superAdminGroup.GET("/duplicates", assetDocumentController.GetDuplicateDocuments)
		// Cleanup all duplicate documents
		superAdminGroup.POST("/cleanup-all", assetDocumentController.CleanupAllDuplicateDocuments)
		// Cleanup duplicate documents for specific asset
		superAdminGroup.POST("/cleanup/:assetId", assetDocumentController.CleanupDuplicateDocuments)

		// Storage management routes
		superAdminGroup.GET("/storage/:asset_id", assetDocumentController.GetStorageInfo)

		// Asset-specific document management (using asset_id parameter)
		superAdminGroup.GET("/documents/:asset_id", assetDocumentController.GetAssetDocuments)
		superAdminGroup.DELETE("/documents/:asset_id", assetDocumentController.DeleteAssetDocuments)
	}
}
