package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(
	router *gin.Engine,
	assetController *controller.AssetController,
	assetTypeController *controller.AssetTypeController,
	locationController *controller.LocationController,
	assetDocumentController *controller.AssetDocumentController,
	assetSensorController *controller.AssetSensorController,
	sensorTypeController *controller.SensorTypeController,
	sensorMeasurementFieldController *controller.SensorMeasurementFieldController,
	sensorMeasurementTypeController *controller.SensorMeasurementTypeController,
	iotSensorReadingController *controller.IoTSensorReadingController,
	jwtConfig middleware.JWTConfig,
) {

	// Public routes (no tenant required)
	public := router.Group("/api/v1")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Setup Asset Type routes
	SetupAssetTypeRoutes(router, assetTypeController)

	// Setup Location routes
	SetupLocationRoutes(router, locationController)

	// Setup Asset routes
	SetupAssetRoutes(router, assetController)

	// Setup Asset Document routes
	SetupAssetDocumentRoutes(router, assetDocumentController)

	// Setup Asset Sensor routes
	SetupAssetSensorRoutes(router, assetSensorController)

	// Setup Sensor Type routes
	SetupSensorTypeRoutes(router, sensorTypeController)

	// Setup Sensor Measurement Field routes
	SetupSensorMeasurementFieldRoutes(router, sensorMeasurementFieldController)

	// Setup Sensor Measurement Type routes
	SetupSensorMeasurementTypeRoutes(router, sensorMeasurementTypeController)

	// Setup IoT Sensor Reading routes
	SetupIoTSensorReadingRoutes(router, iotSensorReadingController)
}
