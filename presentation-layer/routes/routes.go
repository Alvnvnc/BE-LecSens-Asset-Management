package routes

import (
	"be-lecsens/asset_management/data-layer/config"
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/presentation-layer/controller"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(
	router *gin.Engine,
	db *sql.DB,
	assetController *controller.AssetController,
	assetTypeController *controller.AssetTypeController,
	locationController *controller.LocationController,
	assetDocumentController *controller.AssetDocumentController,
	assetSensorController *controller.AssetSensorController,
	sensorTypeController *controller.SensorTypeController,
	sensorMeasurementFieldController *controller.SensorMeasurementFieldController,
	sensorMeasurementTypeController *controller.SensorMeasurementTypeController,
	iotSensorReadingController *controller.IoTSensorReadingController,
	sensorThresholdController *controller.SensorThresholdController,
	assetAlertController *controller.AssetAlertController,
	sensorLogsController *controller.SensorLogsController,
	sensorStatusController *controller.SensorStatusController,
	jwtConfig middleware.JWTConfig,
) {

	// Public routes (no tenant required)
	public := router.Group("/api/v1")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Initialize user management service for middleware that need it
	cfg := config.Load()
	userManagementService := service.NewUserManagementService(
		cfg.User.APIURL,
		cfg.User.APIURL,
		cfg.User.APIURL, // tenantAPIURL - using same URL for now
		cfg.User.APIKey,
		cfg.User.ValidateTokenEndpoint,
		cfg.User.UserInfoEndpoint,
		cfg.User.ValidatePermissionsEndpoint,
		cfg.User.ValidateSuperAdminEndpoint,
	)

	// Test routes to verify middleware functionality
	testGroup := router.Group("/api/v1/test")
	{
		// Test JWT middleware only
		testGroup.GET("/jwt", middleware.JWTMiddleware(jwtConfig, userManagementService), func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			userRole, _ := c.Get("user_role")
			tenantID, _ := c.Get("tenant_id")

			c.JSON(http.StatusOK, gin.H{
				"message":   "JWT middleware test successful",
				"user_id":   userID,
				"user_role": userRole,
				"tenant_id": tenantID,
			})
		})

		// Test TenantMiddleware
		testGroup.GET("/tenant",
			middleware.JWTMiddleware(jwtConfig, userManagementService),
			middleware.TenantMiddleware(),
			func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				userRole, _ := c.Get("user_role")
				tenantID, _ := c.Get("tenant_id")

				c.JSON(http.StatusOK, gin.H{
					"message":   "Tenant middleware test successful",
					"user_id":   userID,
					"user_role": userRole,
					"tenant_id": tenantID,
				})
			})

		// Test TenantAdminMiddleware
		testGroup.GET("/tenant-admin",
			middleware.JWTMiddleware(jwtConfig, userManagementService),
			middleware.TenantAdminMiddleware(),
			func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				userRole, _ := c.Get("user_role")
				tenantID, _ := c.Get("tenant_id")

				c.JSON(http.StatusOK, gin.H{
					"message":   "Tenant admin middleware test successful",
					"user_id":   userID,
					"user_role": userRole,
					"tenant_id": tenantID,
				})
			})

		// Test RequireUserWithTenantMiddleware
		testGroup.GET("/user-tenant",
			middleware.JWTMiddleware(jwtConfig, userManagementService),
			middleware.RequireUserWithTenantMiddleware(),
			func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				userRole, _ := c.Get("user_role")
				tenantID, _ := c.Get("tenant_id")

				c.JSON(http.StatusOK, gin.H{
					"message":   "User with tenant middleware test successful",
					"user_id":   userID,
					"user_role": userRole,
					"tenant_id": tenantID,
				})
			})

		// Test RequireUserRoleMiddleware with multiple roles
		testGroup.GET("/user-role",
			middleware.JWTMiddleware(jwtConfig, userManagementService),
			middleware.RequireUserRoleMiddleware("ADMIN", "USER", "MANAGER"),
			func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				userRole, _ := c.Get("user_role")
				tenantID, _ := c.Get("tenant_id")

				c.JSON(http.StatusOK, gin.H{
					"message":       "User role middleware test successful",
					"user_id":       userID,
					"user_role":     userRole,
					"tenant_id":     tenantID,
					"allowed_roles": []string{"ADMIN", "USER", "MANAGER"},
				})
			})

		// Test SuperAdmin middleware
		testGroup.GET("/superadmin",
			middleware.RequireSuperAdminMiddleware(),
			func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				userRole, _ := c.Get("user_role")
				tenantID, _ := c.Get("tenant_id")

				c.JSON(http.StatusOK, gin.H{
					"message":   "SuperAdmin middleware test successful",
					"user_id":   userID,
					"user_role": userRole,
					"tenant_id": tenantID,
				})
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

	// Setup Asset with Sensors routes
	SetupAssetWithSensorsRoutes(router, db)

	// Setup Sensor Type routes
	SetupSensorTypeRoutes(router, sensorTypeController)

	// Setup Sensor Measurement Field routes
	SetupSensorMeasurementFieldRoutes(router, sensorMeasurementFieldController)

	// Setup Sensor Measurement Type routes
	SetupSensorMeasurementTypeRoutes(router, sensorMeasurementTypeController)

	// Setup IoT Sensor Reading routes
	SetupIoTSensorReadingRoutes(router, iotSensorReadingController)

	// Setup Sensor Threshold routes
	SetupSensorThresholdRoutes(router, sensorThresholdController)

	// Setup Asset Alert routes
	SetupAssetAlertRoutes(router, assetAlertController)

	// Setup Sensor Logs routes
	SetupSensorLogsRoutes(router, sensorLogsController)

	// Setup Sensor Status routes
	SetupSensorStatusRoutes(router, sensorStatusController)
}
