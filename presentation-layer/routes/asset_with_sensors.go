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
		}
	}

	// SuperAdmin only routes - use SuperAdmin middleware for role validation
	superAdminGroup := router.Group("/api/v1/superadmin/assets-with-sensors")
	superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		// Get asset with sensors by ID (for SuperAdmin)
		superAdminGroup.GET("/:id", assetWithSensorsController.GetAssetWithSensors)
		// Create asset with sensors
		superAdminGroup.POST("", assetWithSensorsController.CreateAssetWithSensors)
	}
}
