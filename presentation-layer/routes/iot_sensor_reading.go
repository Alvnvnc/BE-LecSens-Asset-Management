package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupIoTSensorReadingRoutes configures all IoT sensor reading-related routes
func SetupIoTSensorReadingRoutes(router *gin.Engine, iotSensorReadingController *controller.IoTSensorReadingController) {
	// Group for IoT sensor reading routes
	iotSensorReadingGroup := router.Group("/api/v1/iot-sensor-readings")
	{
		// Public routes (requires tenant validation from JWT)
		iotSensorReadingGroup.Use(middleware.TenantMiddleware())
		{
			// List readings with pagination and filtering
			iotSensorReadingGroup.GET("", iotSensorReadingController.ListReadings)
			// Get reading by ID
			iotSensorReadingGroup.GET("/:id", iotSensorReadingController.GetReading)
			// List readings by asset sensor ID
			iotSensorReadingGroup.GET("/by-asset-sensor/:asset_sensor_id", iotSensorReadingController.ListReadingsByAssetSensor)
			// List readings by sensor type ID
			iotSensorReadingGroup.GET("/by-sensor-type/:sensor_type_id", iotSensorReadingController.ListReadingsBySensorType)
			// List readings by MAC address
			iotSensorReadingGroup.GET("/by-mac-address/:mac_address", iotSensorReadingController.ListReadingsByMacAddress)
			// Get latest reading by MAC address
			iotSensorReadingGroup.GET("/latest/by-mac-address/:mac_address", iotSensorReadingController.GetLatestByMacAddress)
			// List readings by time range
			iotSensorReadingGroup.GET("/by-time-range", iotSensorReadingController.ListReadingsByTimeRange)
		}

		// SuperAdmin only routes - use SuperAdmin middleware for role validation
		superAdminGroup := router.Group("/api/v1/superadmin/iot-sensor-readings")
		superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// List all readings (for SuperAdmin - across all tenants)
			superAdminGroup.GET("", iotSensorReadingController.ListAllReadings)
			// Create new reading
			superAdminGroup.POST("", iotSensorReadingController.CreateReading)
			// Update reading
			superAdminGroup.PUT("/:id", iotSensorReadingController.UpdateReading)
			// Delete reading
			superAdminGroup.DELETE("/:id", iotSensorReadingController.DeleteReading)
		}
	}
}
