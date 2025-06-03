package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupSensorThresholdRoutes configures all sensor threshold-related routes
func SetupSensorThresholdRoutes(router *gin.Engine, sensorThresholdController *controller.SensorThresholdController) {
	// Group for sensor threshold routes
	sensorThresholdGroup := router.Group("/api/v1/sensor-thresholds")
	{
		// Public routes (requires tenant validation from JWT)
		sensorThresholdGroup.Use(middleware.TenantMiddleware())
		{
			// List thresholds with pagination and filtering
			sensorThresholdGroup.GET("", sensorThresholdController.ListSensorThresholds)
			// Get threshold by ID
			sensorThresholdGroup.GET("/:id", sensorThresholdController.GetSensorThreshold)
			// List thresholds by asset sensor ID
			sensorThresholdGroup.GET("/by-asset-sensor/:asset_sensor_id", sensorThresholdController.ListThresholdsByAssetSensor)
			// List thresholds by measurement type ID
			sensorThresholdGroup.GET("/by-measurement-type/:measurement_type_id", sensorThresholdController.ListThresholdsByMeasurementType)
		}

		// Admin routes - use TenantAdmin middleware for role validation
		adminGroup := router.Group("/api/v1/admin/sensor-thresholds")
		adminGroup.Use(middleware.TenantAdminMiddleware())
		{
			// Create new threshold
			adminGroup.POST("", sensorThresholdController.CreateSensorThreshold)
			// Update threshold
			adminGroup.PUT("/:id", sensorThresholdController.UpdateSensorThreshold)
			// Delete threshold
			adminGroup.DELETE("/:id", sensorThresholdController.DeleteSensorThreshold)
			// Activate threshold
			adminGroup.POST("/:id/activate", sensorThresholdController.ActivateSensorThreshold)
			// Deactivate threshold
			adminGroup.POST("/:id/deactivate", sensorThresholdController.DeactivateSensorThreshold)
		}

		// SuperAdmin only routes - use SuperAdmin middleware for role validation
		superAdminGroup := router.Group("/api/v1/superadmin/sensor-thresholds")
		superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// List all thresholds (for SuperAdmin - across all tenants)
			superAdminGroup.GET("", sensorThresholdController.ListAllSensorThresholds)
		}
	}
}
