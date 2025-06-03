package routes

import (
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/presentation-layer/controller"
	"database/sql"

	"github.com/gin-gonic/gin"
)

// SetupAssetWithSensorsRoutes sets up routes for asset with sensors operations
func SetupAssetWithSensorsRoutes(router *gin.Engine, db *sql.DB) {
	// Initialize repositories
	assetRepo := repository.NewAssetRepository(db)
	assetSensorRepo := repository.NewAssetSensorRepository(db)
	assetTypeRepo := repository.NewAssetTypeRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	sensorTypeRepo := repository.NewSensorTypeRepository(db)

	// Initialize service
	assetWithSensorsService := service.NewAssetWithSensorsService(
		assetRepo,
		assetSensorRepo,
		assetTypeRepo,
		locationRepo,
		sensorTypeRepo,
		db,
	)

	// Initialize controller
	assetWithSensorsController := controller.NewAssetWithSensorsController(assetWithSensorsService)

	// Group for asset with sensors routes
	assetWithSensorsGroup := router.Group("/api/v1/assets-with-sensors")
	{
		// Public routes (requires tenant validation from JWT) - READ ONLY for regular users
		assetWithSensorsGroup.Use(middleware.TenantMiddleware())
		{
			// Get asset with sensors by ID
			assetWithSensorsGroup.GET("/:id", assetWithSensorsController.GetAssetWithSensors)
			// List assets with sensors (paginated)
			assetWithSensorsGroup.GET("", assetWithSensorsController.ListAssetsWithSensors)
		}
	}

	// SuperAdmin only routes - use SuperAdmin middleware for role validation
	superAdminGroup := router.Group("/api/v1/superadmin/assets-with-sensors")
	superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		// Get asset with sensors by ID (for SuperAdmin)
		superAdminGroup.GET("/:id", assetWithSensorsController.GetAssetWithSensors)
		// List assets with sensors (for SuperAdmin)
		superAdminGroup.GET("", assetWithSensorsController.ListAssetsWithSensors)
		// Create asset with sensors
		superAdminGroup.POST("", assetWithSensorsController.CreateAssetWithSensors)
		// Update asset with sensors
		superAdminGroup.PUT("/:id", assetWithSensorsController.UpdateAssetWithSensors)
		// Delete asset with sensors
		superAdminGroup.DELETE("/:id", assetWithSensorsController.DeleteAssetWithSensors)
		// Bulk create assets with sensors
		superAdminGroup.POST("/bulk", assetWithSensorsController.BulkCreateAssetsWithSensors)
	}
}
