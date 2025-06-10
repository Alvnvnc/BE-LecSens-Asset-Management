package seeder

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

// SensorThresholdSeeder handles seeding sensor thresholds
type SensorThresholdSeeder struct {
	db *sql.DB
}

// NewSensorThresholdSeeder creates a new sensor threshold seeder
func NewSensorThresholdSeeder(db *sql.DB) *SensorThresholdSeeder {
	return &SensorThresholdSeeder{db: db}
}

// Predefined UUIDs for sensor thresholds to ensure consistency
var (
	// Water Quality pH Thresholds
	ThresholdPhWarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222101")
	ThresholdPhCriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222102")

	// Turbidity Thresholds
	ThresholdTurbidityWarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222201")
	ThresholdTurbidityCriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222202")

	// Air Quality PM2.5 Thresholds
	ThresholdPm25WarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222301")
	ThresholdPm25CriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222302")

	// CO2 Thresholds
	ThresholdCo2WarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222303")
	ThresholdCo2CriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222304")

	// Temperature Thresholds
	ThresholdTempWarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222401")
	ThresholdTempCriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222402")

	// Vibration Thresholds
	ThresholdVibrationWarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222501")
	ThresholdVibrationCriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222502")

	// Flow Rate Thresholds
	ThresholdFlowWarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222601")
	ThresholdFlowCriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222602")

	// Pressure Thresholds
	ThresholdPressureWarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222701")
	ThresholdPressureCriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222702")

	// Level Thresholds
	ThresholdLevelWarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222801")
	ThresholdLevelCriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222802")

	// Power Thresholds
	ThresholdPowerWarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222901")
	ThresholdPowerCriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222902")

	// Gas Detection Thresholds
	ThresholdGasWarningID  = uuid.MustParse("22222222-2222-2222-2222-222222222A01")
	ThresholdGasCriticalID = uuid.MustParse("22222222-2222-2222-2222-222222222A02")
)

// getSensorThresholds returns all sensor thresholds to be seeded
func (s *SensorThresholdSeeder) getSensorThresholds() []entity.SensorThreshold {
	now := time.Now()
	tenantID := *GetDefaultTenantID()

	// Define threshold values for different measurements
	phMinWarning, phMaxWarning := 6.0, 8.5
	phMinCritical, phMaxCritical := 5.0, 9.5

	turbidityMaxWarning := 20.0
	turbidityMaxCritical := 50.0

	pm25MaxWarning := 35.0
	pm25MaxCritical := 75.0

	co2MaxWarning := 1000.0
	co2MaxCritical := 2000.0

	tempMinWarning, tempMaxWarning := 10.0, 40.0
	tempMinCritical, tempMaxCritical := 5.0, 50.0

	vibrationMaxWarning := 10.0
	vibrationMaxCritical := 25.0

	flowMinWarning := 50.0
	flowMinCritical := 20.0

	pressureMinWarning, pressureMaxWarning := 2.0, 8.0
	pressureMinCritical, pressureMaxCritical := 1.0, 10.0

	levelMinWarning, levelMaxWarning := 0.5, 5.0
	levelMinCritical, levelMaxCritical := 0.2, 6.0

	powerMaxWarning := 500.0
	powerMaxCritical := 800.0

	gasMaxWarning := 100.0
	gasMaxCritical := 500.0

	return []entity.SensorThreshold{
		// Water Quality pH Thresholds (pH Probe)
		{
			ID:                   ThresholdPhWarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorPHProbeID,
			MeasurementTypeID:    MeasurementTypeWaterQualityID,
			MeasurementFieldName: "ph_value",
			MinValue:             &phMinWarning,
			MaxValue:             &phMaxWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdPhCriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorPHProbeID,
			MeasurementTypeID:    MeasurementTypeWaterQualityID,
			MeasurementFieldName: "ph_value",
			MinValue:             &phMinCritical,
			MaxValue:             &phMaxCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},

		// Turbidity Thresholds (Turbidity Sensor)
		{
			ID:                   ThresholdTurbidityWarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorTurbidityID,
			MeasurementTypeID:    MeasurementTypeTurbidityID,
			MeasurementFieldName: "turbidity_value",
			MaxValue:             &turbidityMaxWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdTurbidityCriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorTurbidityID,
			MeasurementTypeID:    MeasurementTypeTurbidityID,
			MeasurementFieldName: "turbidity_value",
			MaxValue:             &turbidityMaxCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},

		// Air Quality PM2.5 Thresholds (PM2.5 Sensor)
		{
			ID:                   ThresholdPm25WarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorPM25ID,
			MeasurementTypeID:    MeasurementTypeAirQualityID,
			MeasurementFieldName: "pm25_concentration",
			MaxValue:             &pm25MaxWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdPm25CriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorPM25ID,
			MeasurementTypeID:    MeasurementTypeAirQualityID,
			MeasurementFieldName: "pm25_concentration",
			MaxValue:             &pm25MaxCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},

		// CO2 Thresholds (CO2 Sensor)
		{
			ID:                   ThresholdCo2WarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorCO2ID,
			MeasurementTypeID:    MeasurementTypeAirQualityID,
			MeasurementFieldName: "co2_concentration",
			MaxValue:             &co2MaxWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdCo2CriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorCO2ID,
			MeasurementTypeID:    MeasurementTypeAirQualityID,
			MeasurementFieldName: "co2_concentration",
			MaxValue:             &co2MaxCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},

		// Temperature Thresholds (Temperature Sensor)
		{
			ID:                   ThresholdTempWarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorAmbientTempID,
			MeasurementTypeID:    MeasurementTypeTemperatureID,
			MeasurementFieldName: "temperature_value",
			MinValue:             &tempMinWarning,
			MaxValue:             &tempMaxWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdTempCriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorAmbientTempID,
			MeasurementTypeID:    MeasurementTypeTemperatureID,
			MeasurementFieldName: "temperature_value",
			MinValue:             &tempMinCritical,
			MaxValue:             &tempMaxCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},

		// Vibration Thresholds (Vibration Sensor)
		{
			ID:                   ThresholdVibrationWarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorVibrationXID,
			MeasurementTypeID:    MeasurementTypeVibrationID,
			MeasurementFieldName: "vibration_amplitude",
			MaxValue:             &vibrationMaxWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdVibrationCriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorVibrationXID,
			MeasurementTypeID:    MeasurementTypeVibrationID,
			MeasurementFieldName: "vibration_amplitude",
			MaxValue:             &vibrationMaxCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},

		// Flow Rate Thresholds (Flow Sensor)
		{
			ID:                   ThresholdFlowWarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorFlowRateID,
			MeasurementTypeID:    MeasurementTypeFlowID,
			MeasurementFieldName: "flow_rate",
			MinValue:             &flowMinWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdFlowCriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorFlowRateID,
			MeasurementTypeID:    MeasurementTypeFlowID,
			MeasurementFieldName: "flow_rate",
			MinValue:             &flowMinCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},

		// Pressure Thresholds (Pressure Sensor)
		{
			ID:                   ThresholdPressureWarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorPressure1ID,
			MeasurementTypeID:    MeasurementTypePressureID,
			MeasurementFieldName: "pressure_value",
			MinValue:             &pressureMinWarning,
			MaxValue:             &pressureMaxWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdPressureCriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorPressure1ID,
			MeasurementTypeID:    MeasurementTypePressureID,
			MeasurementFieldName: "pressure_value",
			MinValue:             &pressureMinCritical,
			MaxValue:             &pressureMaxCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},

		// Level Thresholds (Level Sensor)
		{
			ID:                   ThresholdLevelWarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorLevelTank1ID,
			MeasurementTypeID:    MeasurementTypeLevelID,
			MeasurementFieldName: "level_value",
			MinValue:             &levelMinWarning,
			MaxValue:             &levelMaxWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdLevelCriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorLevelTank1ID,
			MeasurementTypeID:    MeasurementTypeLevelID,
			MeasurementFieldName: "level_value",
			MinValue:             &levelMinCritical,
			MaxValue:             &levelMaxCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},

		// Power Thresholds (Power Meter)
		{
			ID:                   ThresholdPowerWarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorPowerMeterID,
			MeasurementTypeID:    MeasurementTypePowerID,
			MeasurementFieldName: "active_power",
			MaxValue:             &powerMaxWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdPowerCriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorPowerMeterID,
			MeasurementTypeID:    MeasurementTypePowerID,
			MeasurementFieldName: "active_power",
			MaxValue:             &powerMaxCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},

		// Gas Detection Thresholds (Gas Sensor)
		{
			ID:                   ThresholdGasWarningID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorGasCO2ID,
			MeasurementTypeID:    MeasurementTypeGasDetectionID,
			MeasurementFieldName: "gas_concentration",
			MaxValue:             &gasMaxWarning,
			Severity:             entity.ThresholdSeverityWarning,
			IsActive:             true,
			CreatedAt:            now,
		},
		{
			ID:                   ThresholdGasCriticalID,
			TenantID:             tenantID,
			AssetSensorID:        AssetSensorGasCO2ID,
			MeasurementTypeID:    MeasurementTypeGasDetectionID,
			MeasurementFieldName: "gas_concentration",
			MaxValue:             &gasMaxCritical,
			Severity:             entity.ThresholdSeverityCritical,
			IsActive:             true,
			CreatedAt:            now,
		},
	}
}

// thresholdExists checks if a sensor threshold exists
func (s *SensorThresholdSeeder) thresholdExists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sensor_thresholds WHERE id = $1)`
	var exists bool
	err := s.db.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// createThreshold creates a single sensor threshold
func (s *SensorThresholdSeeder) createThreshold(ctx context.Context, threshold entity.SensorThreshold) error {
	query := `
		INSERT INTO sensor_thresholds (
			id, tenant_id, asset_sensor_id, measurement_type_id,
			measurement_field_name, min_value, max_value, severity,
			is_active, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := s.db.ExecContext(ctx, query,
		threshold.ID,
		threshold.TenantID,
		threshold.AssetSensorID,
		threshold.MeasurementTypeID,
		threshold.MeasurementFieldName,
		threshold.MinValue,
		threshold.MaxValue,
		threshold.Severity,
		threshold.IsActive,
		threshold.CreatedAt,
	)

	return err
}

// SeedThresholds seeds all sensor thresholds
func (s *SensorThresholdSeeder) SeedThresholds(ctx context.Context) error {
	log.Println("Starting sensor threshold seeding...")

	thresholds := s.getSensorThresholds()
	createdCount := 0
	skippedCount := 0

	for _, threshold := range thresholds {
		exists, err := s.thresholdExists(ctx, threshold.ID)
		if err != nil {
			log.Printf("Error checking threshold existence for %s: %v", threshold.MeasurementFieldName, err)
			continue
		}

		if exists {
			log.Printf("Threshold %s for %s already exists, skipping", threshold.Severity, threshold.MeasurementFieldName)
			skippedCount++
			continue
		}

		if err := s.createThreshold(ctx, threshold); err != nil {
			log.Printf("Error creating threshold %s for %s: %v", threshold.Severity, threshold.MeasurementFieldName, err)
			continue
		}

		log.Printf("Created sensor threshold: %s - %s (%s)", threshold.MeasurementFieldName, threshold.Severity, threshold.ID)
		createdCount++
	}

	log.Printf("Sensor threshold seeding completed. Created: %d, Skipped: %d", createdCount, skippedCount)
	return nil
}

// GetThresholdIDByDetails returns the UUID for a threshold by its details
func GetThresholdIDByDetails(assetSensorID uuid.UUID, fieldName string, severity entity.ThresholdSeverity) uuid.UUID {
	// Create a mapping for quick lookups based on asset sensor and field
	thresholdMap := map[string]map[string]map[entity.ThresholdSeverity]uuid.UUID{
		AssetSensorPHProbeID.String(): {
			"ph_value": {
				entity.ThresholdSeverityWarning:  ThresholdPhWarningID,
				entity.ThresholdSeverityCritical: ThresholdPhCriticalID,
			},
		},
		AssetSensorTurbidityID.String(): {
			"turbidity_value": {
				entity.ThresholdSeverityWarning:  ThresholdTurbidityWarningID,
				entity.ThresholdSeverityCritical: ThresholdTurbidityCriticalID,
			},
		},
		AssetSensorPM25ID.String(): {
			"pm25_concentration": {
				entity.ThresholdSeverityWarning:  ThresholdPm25WarningID,
				entity.ThresholdSeverityCritical: ThresholdPm25CriticalID,
			},
		},
		AssetSensorCO2ID.String(): {
			"co2_concentration": {
				entity.ThresholdSeverityWarning:  ThresholdCo2WarningID,
				entity.ThresholdSeverityCritical: ThresholdCo2CriticalID,
			},
		},
		AssetSensorAmbientTempID.String(): {
			"temperature_value": {
				entity.ThresholdSeverityWarning:  ThresholdTempWarningID,
				entity.ThresholdSeverityCritical: ThresholdTempCriticalID,
			},
		},
		AssetSensorVibrationXID.String(): {
			"vibration_amplitude": {
				entity.ThresholdSeverityWarning:  ThresholdVibrationWarningID,
				entity.ThresholdSeverityCritical: ThresholdVibrationCriticalID,
			},
		},
		AssetSensorFlowRateID.String(): {
			"flow_rate": {
				entity.ThresholdSeverityWarning:  ThresholdFlowWarningID,
				entity.ThresholdSeverityCritical: ThresholdFlowCriticalID,
			},
		},
		AssetSensorPressure1ID.String(): {
			"pressure_value": {
				entity.ThresholdSeverityWarning:  ThresholdPressureWarningID,
				entity.ThresholdSeverityCritical: ThresholdPressureCriticalID,
			},
		},
		AssetSensorLevelTank1ID.String(): {
			"level_value": {
				entity.ThresholdSeverityWarning:  ThresholdLevelWarningID,
				entity.ThresholdSeverityCritical: ThresholdLevelCriticalID,
			},
		},
		AssetSensorPowerMeterID.String(): {
			"active_power": {
				entity.ThresholdSeverityWarning:  ThresholdPowerWarningID,
				entity.ThresholdSeverityCritical: ThresholdPowerCriticalID,
			},
		},
		AssetSensorGasCO2ID.String(): {
			"gas_concentration": {
				entity.ThresholdSeverityWarning:  ThresholdGasWarningID,
				entity.ThresholdSeverityCritical: ThresholdGasCriticalID,
			},
		},
	}

	if sensorFields, exists := thresholdMap[assetSensorID.String()]; exists {
		if fieldThresholds, found := sensorFields[fieldName]; found {
			if thresholdID, ok := fieldThresholds[severity]; ok {
				return thresholdID
			}
		}
	}

	return uuid.Nil
}
