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
			// Get latest reading by asset sensor
			iotSensorReadingGroup.GET("/latest/by-asset-sensor/:asset_sensor_id", iotSensorReadingController.GetLatestByAssetSensor)
			// List readings by time range
			iotSensorReadingGroup.GET("/by-time-range", iotSensorReadingController.ListReadingsByTimeRange)
			// Get aggregated data for analytics
			iotSensorReadingGroup.GET("/aggregated", iotSensorReadingController.GetAggregatedData)
			// Get auto-population options for sensor type
			iotSensorReadingGroup.GET("/auto-populate/options", iotSensorReadingController.GetAutoPopulationOptions)
		}

		// SuperAdmin only routes - use SuperAdmin middleware for role validation
		superAdminGroup := router.Group("/api/v1/superadmin/iot-sensor-readings")
		superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// List all readings (for SuperAdmin - across all tenants)
			superAdminGroup.GET("", iotSensorReadingController.ListAllReadings)
			// Create new reading
			superAdminGroup.POST("", iotSensorReadingController.CreateReading)
			// Create batch readings
			superAdminGroup.POST("/batch", iotSensorReadingController.CreateBatchReading)
			// Create reading with auto-population of asset_sensor_id
			superAdminGroup.POST("/auto-populate", iotSensorReadingController.CreateReadingWithAutoPopulation)
			// Validate and create reading with schema validation
			superAdminGroup.POST("/validate", iotSensorReadingController.ValidateAndCreateReading)
			// Update reading
			superAdminGroup.PUT("/:id", iotSensorReadingController.UpdateReading)
			// Delete reading
			superAdminGroup.DELETE("/:id", iotSensorReadingController.DeleteReading)

			// Manual data input routes
			// Create reading from JSON string
			superAdminGroup.POST("/from-json", iotSensorReadingController.CreateFromJSON)
			// Create batch readings from JSON string
			superAdminGroup.POST("/batch-from-json", iotSensorReadingController.CreateBatchFromJSON)
			// Create dummy reading for testing
			superAdminGroup.POST("/dummy", iotSensorReadingController.CreateDummyReading)
			// Create multiple dummy readings for testing
			superAdminGroup.POST("/dummy/multiple", iotSensorReadingController.CreateMultipleDummyReadings)
			// Get JSON template for specific sensor type
			superAdminGroup.GET("/template", iotSensorReadingController.GetJSONTemplate)
			// Get batch JSON template for specific sensor type
			superAdminGroup.GET("/template/batch", iotSensorReadingController.GetBatchJSONTemplate)
			// Create simple reading with minimal fields
			superAdminGroup.POST("/simple", iotSensorReadingController.CreateSimpleReading)

			// Flexible JSON input routes - handles dynamic measurement data
			// Create flexible reading with dynamic JSON structure
			superAdminGroup.POST("/flexible", iotSensorReadingController.CreateFlexibleReading)
			// Create bulk flexible readings with dynamic JSON structure
			superAdminGroup.POST("/flexible/bulk", iotSensorReadingController.CreateBulkFlexibleReadings)
			// Parse text input to flexible JSON structure
			superAdminGroup.POST("/parse-text", iotSensorReadingController.ParseTextToFlexible)
			// List flexible readings with pagination and filtering
			superAdminGroup.GET("/flexible", iotSensorReadingController.ListFlexibleReadings)
			// Get flexible reading by ID with dynamic measurement data
			superAdminGroup.GET("/flexible/:id", iotSensorReadingController.GetFlexibleReading)
			// Create reading from any raw JSON structure
			superAdminGroup.POST("/from-raw-json", iotSensorReadingController.CreateFromRawJSON)
			// Create readings from JSON array with any structure
			superAdminGroup.POST("/from-array", iotSensorReadingController.CreateFromArrayJSON)
		}
	}
}
