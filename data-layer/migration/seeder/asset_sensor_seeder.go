package seeder

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type AssetSensorSeeder struct {
	db *sql.DB
}

func NewAssetSensorSeeder(db *sql.DB) *AssetSensorSeeder {
	return &AssetSensorSeeder{db: db}
}

// Predefined asset sensor UUIDs for consistency across seeders
var (
	// Water Quality Sensors
	AssetSensorPHProbeID      = uuid.MustParse("650e8400-e29b-41d4-a716-446655440001")
	AssetSensorTurbidityID    = uuid.MustParse("650e8400-e29b-41d4-a716-446655440002")
	AssetSensorDOProbeID      = uuid.MustParse("650e8400-e29b-41d4-a716-446655440003")
	AssetSensorWaterConductID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440004")
	AssetSensorChlorineID     = uuid.MustParse("650e8400-e29b-41d4-a716-446655440025")
	AssetSensorTDSID          = uuid.MustParse("650e8400-e29b-41d4-a716-446655440026")
	AssetSensorORPID          = uuid.MustParse("650e8400-e29b-41d4-a716-446655440027")

	// Air Quality Sensors
	AssetSensorPM25ID        = uuid.MustParse("650e8400-e29b-41d4-a716-446655440005")
	AssetSensorCO2ID         = uuid.MustParse("650e8400-e29b-41d4-a716-446655440006")
	AssetSensorAirTempID     = uuid.MustParse("650e8400-e29b-41d4-a716-446655440007")
	AssetSensorAirHumidityID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440008")

	// Temperature Sensors
	AssetSensorColdStorageTempID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440009")
	AssetSensorAmbientTempID     = uuid.MustParse("650e8400-e29b-41d4-a716-446655440010")
	AssetSensorDataCenterTempID  = uuid.MustParse("650e8400-e29b-41d4-a716-446655440011")

	// Vibration Sensors
	AssetSensorVibrationXID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440012")
	AssetSensorVibrationYID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440013")

	// Flow Sensors
	AssetSensorFlowRateID  = uuid.MustParse("650e8400-e29b-41d4-a716-446655440014")
	AssetSensorFlowTotalID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440015")

	// Pressure Sensors
	AssetSensorPressure1ID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440016")
	AssetSensorPressure2ID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440017")

	// Level Sensors
	AssetSensorLevelTank1ID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440018")
	AssetSensorLevelTank2ID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440019")

	// Energy Sensors
	AssetSensorPowerMeterID   = uuid.MustParse("650e8400-e29b-41d4-a716-446655440020")
	AssetSensorCurrentMeterID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440021")
	AssetSensorVoltageMeterID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440022")

	// Gas and Other Sensors
	AssetSensorGasCO2ID = uuid.MustParse("650e8400-e29b-41d4-a716-446655440023")
	AssetSensorNoiseID  = uuid.MustParse("650e8400-e29b-41d4-a716-446655440024")
)

// Sensor Type ID constants - should match the IDs in sensor_type_seeder.go
var (
	SensorTypePHID           = uuid.MustParse("03456789-3333-3333-3333-000000000001")
	SensorTypeTurbidityID    = uuid.MustParse("03456789-3333-3333-3333-000000000002")
	SensorTypeDOID           = uuid.MustParse("03456789-3333-3333-3333-000000000003")
	SensorTypeConductivityID = uuid.MustParse("03456789-3333-3333-3333-000000000004")
	SensorTypePM25ID         = uuid.MustParse("03456789-3333-3333-3333-000000000005")
	SensorTypeCO2ID          = uuid.MustParse("03456789-3333-3333-3333-000000000006")
	SensorTypeTemperatureID  = uuid.MustParse("03456789-3333-3333-3333-000000000007")
	SensorTypeHumidityID     = uuid.MustParse("03456789-3333-3333-3333-000000000008")
	SensorTypeVibrationID    = uuid.MustParse("03456789-3333-3333-3333-000000000009")
	SensorTypeFlowID         = uuid.MustParse("03456789-3333-3333-3333-000000000010")
	SensorTypePressureID     = uuid.MustParse("03456789-3333-3333-3333-000000000011")
	SensorTypeLevelID        = uuid.MustParse("03456789-3333-3333-3333-000000000012")
	SensorTypePowerID        = uuid.MustParse("03456789-3333-3333-3333-000000000013")
	SensorTypeGasID          = uuid.MustParse("03456789-3333-3333-3333-000000000014")
	SensorTypeNoiseID        = uuid.MustParse("03456789-3333-3333-3333-000000000015")
	SensorTypeChlorineID     = uuid.MustParse("03456789-3333-3333-3333-000000000016")
	SensorTypeTDSID          = uuid.MustParse("03456789-3333-3333-3333-000000000017")
	SensorTypeORPID          = uuid.MustParse("03456789-3333-3333-3333-000000000018")
)

func (s *AssetSensorSeeder) Seed(ctx context.Context) error {
	log.Println("Starting Asset Sensor seeder...")

	// Check if asset sensors already exist
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM asset_sensors").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing asset sensors: %w", err)
	}

	if count > 0 {
		log.Printf("Asset sensors already exist (%d records), skipping seed", count)
		return nil
	}

	// Get tenant ID
	tenantID := GetDefaultTenantID()

	assetSensors := s.getAssetSensorData(tenantID)

	// Insert asset sensors
	for _, assetSensor := range assetSensors {
		err := s.createAssetSensor(ctx, assetSensor)
		if err != nil {
			log.Printf("Failed to create asset sensor %s: %v", assetSensor.Name, err)
			return err
		}
		log.Printf("Created asset sensor: %s", assetSensor.Name)
	}

	log.Printf("Successfully seeded %d asset sensors", len(assetSensors))
	return nil
}

func (s *AssetSensorSeeder) getAssetSensorData(tenantID *uuid.UUID) []*entity.AssetSensor {
	now := time.Now()

	// Configuration templates for different sensor types
	phConfig := map[string]interface{}{
		"calibration_points": []map[string]interface{}{
			{"ph": 4.0, "voltage": 177.48},
			{"ph": 7.0, "voltage": 0.0},
			{"ph": 10.0, "voltage": -177.48},
		},
		"temperature_compensation": true,
		"measurement_interval":     30,
		"auto_cleaning":            true,
	}

	turbidityConfig := map[string]interface{}{
		"range": map[string]float64{
			"min": 0.0,
			"max": 1000.0,
		},
		"units":                "NTU",
		"measurement_interval": 60,
		"wiper_enabled":        true,
	}

	doConfig := map[string]interface{}{
		"calibration": map[string]interface{}{
			"zero_point":    0.0,
			"span_point":    8.0,
			"last_cal_date": "2024-01-15",
		},
		"temperature_compensation": true,
		"salinity_compensation":    true,
		"pressure_compensation":    true,
	}

	pm25Config := map[string]interface{}{
		"measurement_method": "laser_scattering",
		"size_range": map[string]float64{
			"min": 0.3,
			"max": 10.0,
		},
		"sampling_interval": 180,
		"automatic_zeroing": true,
		"heater_enabled":    true,
	}

	co2Config := map[string]interface{}{
		"measurement_principle": "NDIR",
		"range": map[string]float64{
			"min": 0.0,
			"max": 5000.0,
		},
		"units":                 "ppm",
		"abc_logic_enabled":     true,
		"pressure_compensation": true,
	}

	tempConfig := map[string]interface{}{
		"sensor_type": "PT1000",
		"measurement_range": map[string]float64{
			"min": -40.0,
			"max": 85.0,
		},
		"accuracy":      0.1,
		"response_time": "< 5 seconds",
		"housing":       "IP67",
	}

	vibrationConfig := map[string]interface{}{
		"sensor_type": "accelerometer",
		"frequency_range": map[string]float64{
			"min": 1.0,
			"max": 10000.0,
		},
		"sensitivity":      "100 mV/g",
		"measurement_axis": "triaxial",
		"fft_enabled":      true,
		"peak_detection":   true,
	}

	flowConfig := map[string]interface{}{
		"sensor_type": "electromagnetic",
		"pipe_size":   "DN100",
		"flow_range": map[string]float64{
			"min": 0.0,
			"max": 500.0,
		},
		"units":                "m³/h",
		"totalizer_enabled":    true,
		"empty_pipe_detection": true,
	}

	pressureConfig := map[string]interface{}{
		"sensor_type": "piezoresistive",
		"pressure_range": map[string]float64{
			"min": 0.0,
			"max": 16.0,
		},
		"units":    "bar",
		"accuracy": "±0.1%",
		"output":   "4-20mA",
		"housing":  "316L SS",
	}

	levelConfig := map[string]interface{}{
		"sensor_type": "guided_wave_radar",
		"measurement_range": map[string]float64{
			"min": 0.0,
			"max": 20.0,
		},
		"units":               "meters",
		"accuracy":            "±2mm",
		"process_temperature": 150.0,
		"process_pressure":    16.0,
	}

	powerConfig := map[string]interface{}{
		"meter_type":         "smart_meter",
		"voltage_range":      "85-265V",
		"current_range":      "0-200A",
		"frequency":          "50/60Hz",
		"accuracy_class":     "0.5S",
		"communication":      []string{"Modbus RTU", "Ethernet"},
		"energy_calculation": true,
	}

	// Additional water quality sensor configurations
	chlorineConfig := map[string]interface{}{
		"sensor_type": "electrochemical",
		"measurement_range": map[string]float64{
			"min": 0.0,
			"max": 10.0,
		},
		"units":                    "mg/L",
		"accuracy":                 "±0.1 mg/L",
		"response_time":            "< 30 seconds",
		"temperature_compensation": true,
		"ph_compensation":          true,
		"auto_cleaning":            true,
	}

	tdsConfig := map[string]interface{}{
		"sensor_type": "conductivity_based",
		"measurement_range": map[string]float64{
			"min": 0.0,
			"max": 2000.0,
		},
		"units":                    "ppm",
		"accuracy":                 "±2%",
		"temperature_compensation": true,
		"conductivity_factor":      0.64,
		"calibration_solution":     "1413 μS/cm",
	}

	orpConfig := map[string]interface{}{
		"sensor_type": "platinum_electrode",
		"measurement_range": map[string]float64{
			"min": -1000.0,
			"max": 1000.0,
		},
		"units":                    "mV",
		"accuracy":                 "±1 mV",
		"reference_electrode":      "Ag/AgCl",
		"temperature_compensation": true,
		"auto_calibration":         false,
	}

	return []*entity.AssetSensor{
		// Water Quality Monitoring Sensors
		{
			ID:            AssetSensorPHProbeID,
			TenantID:      tenantID,
			AssetID:       AssetProdLineAID,
			SensorTypeID:  SensorTypePHID,
			Name:          "Production Line A - pH Probe",
			Status:        "active",
			Configuration: jsonRawMessage(phConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorTurbidityID,
			TenantID:      tenantID,
			AssetID:       AssetProdLineAID,
			SensorTypeID:  SensorTypeTurbidityID,
			Name:          "Production Line A - Turbidity Sensor",
			Status:        "active",
			Configuration: jsonRawMessage(turbidityConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorDOProbeID,
			TenantID:      tenantID,
			AssetID:       AssetWaterTreatmentID,
			SensorTypeID:  SensorTypeDOID,
			Name:          "Water Treatment - Dissolved Oxygen Probe",
			Status:        "active",
			Configuration: jsonRawMessage(doConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorWaterConductID,
			TenantID:      tenantID,
			AssetID:       AssetWaterTreatmentID,
			SensorTypeID:  SensorTypeConductivityID,
			Name:          "Water Treatment - Conductivity Sensor",
			Status:        "active",
			Configuration: jsonRawMessage(doConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorChlorineID,
			TenantID:      tenantID,
			AssetID:       AssetWaterTreatmentID,
			SensorTypeID:  SensorTypeChlorineID,
			Name:          "Water Treatment - Chlorine Sensor",
			Status:        "active",
			Configuration: jsonRawMessage(chlorineConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorTDSID,
			TenantID:      tenantID,
			AssetID:       AssetWaterTreatmentID,
			SensorTypeID:  SensorTypeTDSID,
			Name:          "Water Treatment - TDS Sensor",
			Status:        "active",
			Configuration: jsonRawMessage(tdsConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorORPID,
			TenantID:      tenantID,
			AssetID:       AssetWaterTreatmentID,
			SensorTypeID:  SensorTypeORPID,
			Name:          "Water Treatment - ORP Sensor",
			Status:        "active",
			Configuration: jsonRawMessage(orpConfig),
			CreatedAt:     now,
		},

		// Air Quality Monitoring Sensors
		{
			ID:            AssetSensorPM25ID,
			TenantID:      tenantID,
			AssetID:       AssetAirQualityHubID,
			SensorTypeID:  SensorTypePM25ID,
			Name:          "Factory Air - PM2.5 Sensor",
			Status:        "active",
			Configuration: jsonRawMessage(pm25Config),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorCO2ID,
			TenantID:      tenantID,
			AssetID:       AssetAirQualityHubID,
			SensorTypeID:  SensorTypeCO2ID,
			Name:          "Factory Air - CO2 Sensor",
			Status:        "active",
			Configuration: jsonRawMessage(co2Config),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorAirTempID,
			TenantID:      tenantID,
			AssetID:       AssetAirQualityHubID,
			SensorTypeID:  SensorTypeTemperatureID,
			Name:          "Factory Air - Temperature Sensor",
			Status:        "active",
			Configuration: jsonRawMessage(tempConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorAirHumidityID,
			TenantID:      tenantID,
			AssetID:       AssetWarehouseHubID,
			SensorTypeID:  SensorTypeHumidityID,
			Name:          "Warehouse - Humidity Sensor",
			Status:        "maintenance",
			Configuration: jsonRawMessage(tempConfig),
			CreatedAt:     now,
		},

		// Temperature Monitoring Sensors
		{
			ID:            AssetSensorColdStorageTempID,
			TenantID:      tenantID,
			AssetID:       AssetTempMonitoringID,
			SensorTypeID:  SensorTypeTemperatureID,
			Name:          "Cold Storage - Zone 1 Temperature",
			Status:        "active",
			Configuration: jsonRawMessage(tempConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorAmbientTempID,
			TenantID:      tenantID,
			AssetID:       AssetTempMonitoringID,
			SensorTypeID:  SensorTypeTemperatureID,
			Name:          "Cold Storage - Ambient Temperature",
			Status:        "active",
			Configuration: jsonRawMessage(tempConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorDataCenterTempID,
			TenantID:      tenantID,
			AssetID:       AssetDataCenterHubID,
			SensorTypeID:  SensorTypeTemperatureID,
			Name:          "Data Center - Rack Temperature",
			Status:        "active",
			Configuration: jsonRawMessage(tempConfig),
			CreatedAt:     now,
		},

		// Vibration Monitoring Sensors
		{
			ID:            AssetSensorVibrationXID,
			TenantID:      tenantID,
			AssetID:       AssetVibrationHubID,
			SensorTypeID:  SensorTypeVibrationID,
			Name:          "Turbine 1 - X-Axis Vibration",
			Status:        "active",
			Configuration: jsonRawMessage(vibrationConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorVibrationYID,
			TenantID:      tenantID,
			AssetID:       AssetVibrationHubID,
			SensorTypeID:  SensorTypeVibrationID,
			Name:          "Turbine 1 - Y-Axis Vibration",
			Status:        "active",
			Configuration: jsonRawMessage(vibrationConfig),
			CreatedAt:     now,
		},

		// Flow Monitoring Sensors
		{
			ID:            AssetSensorFlowRateID,
			TenantID:      tenantID,
			AssetID:       AssetFlowStationID,
			SensorTypeID:  SensorTypeFlowID,
			Name:          "Pipeline Main - Flow Rate",
			Status:        "active",
			Configuration: jsonRawMessage(flowConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorFlowTotalID,
			TenantID:      tenantID,
			AssetID:       AssetFlowStationID,
			SensorTypeID:  SensorTypeFlowID,
			Name:          "Pipeline Bypass - Flow Rate",
			Status:        "active",
			Configuration: jsonRawMessage(flowConfig),
			CreatedAt:     now,
		},

		// Pressure Monitoring Sensors
		{
			ID:            AssetSensorPressure1ID,
			TenantID:      tenantID,
			AssetID:       AssetPressureHubID,
			SensorTypeID:  SensorTypePressureID,
			Name:          "Reactor 1 - Inlet Pressure",
			Status:        "active",
			Configuration: jsonRawMessage(pressureConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorPressure2ID,
			TenantID:      tenantID,
			AssetID:       AssetPressureHubID,
			SensorTypeID:  SensorTypePressureID,
			Name:          "Reactor 1 - Outlet Pressure",
			Status:        "active",
			Configuration: jsonRawMessage(pressureConfig),
			CreatedAt:     now,
		},

		// Level Monitoring Sensors
		{
			ID:            AssetSensorLevelTank1ID,
			TenantID:      tenantID,
			AssetID:       AssetLevelMonitorID,
			SensorTypeID:  SensorTypeLevelID,
			Name:          "Fuel Tank 1 - Level Sensor",
			Status:        "active",
			Configuration: jsonRawMessage(levelConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorLevelTank2ID,
			TenantID:      tenantID,
			AssetID:       AssetLevelMonitorID,
			SensorTypeID:  SensorTypeLevelID,
			Name:          "Fuel Tank 2 - Level Sensor",
			Status:        "active",
			Configuration: jsonRawMessage(levelConfig),
			CreatedAt:     now,
		},

		// Energy Monitoring Sensors
		{
			ID:            AssetSensorPowerMeterID,
			TenantID:      tenantID,
			AssetID:       AssetEnergyMeterID,
			SensorTypeID:  SensorTypePowerID,
			Name:          "Resort Main - Power Meter",
			Status:        "active",
			Configuration: jsonRawMessage(powerConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorCurrentMeterID,
			TenantID:      tenantID,
			AssetID:       AssetPowerDistribution,
			SensorTypeID:  SensorTypePowerID,
			Name:          "Distribution Panel - Current Meter",
			Status:        "active",
			Configuration: jsonRawMessage(powerConfig),
			CreatedAt:     now,
		},
		{
			ID:            AssetSensorVoltageMeterID,
			TenantID:      tenantID,
			AssetID:       AssetPowerDistribution,
			SensorTypeID:  SensorTypePowerID,
			Name:          "Distribution Panel - Voltage Meter",
			Status:        "active",
			Configuration: jsonRawMessage(powerConfig),
			CreatedAt:     now,
		},

		// Gas and Other Sensors
		{
			ID:            AssetSensorGasCO2ID,
			TenantID:      tenantID,
			AssetID:       AssetWarehouseHubID,
			SensorTypeID:  SensorTypeGasID,
			Name:          "Warehouse - CO2 Gas Detector",
			Status:        "active",
			Configuration: jsonRawMessage(co2Config),
			CreatedAt:     now,
		},
		{
			ID:           AssetSensorNoiseID,
			TenantID:     tenantID,
			AssetID:      AssetAirQualityHubID,
			SensorTypeID: SensorTypeNoiseID,
			Name:         "Factory - Noise Level Monitor",
			Status:       "active",
			Configuration: jsonRawMessage(map[string]interface{}{
				"frequency_weighting": "A",
				"time_weighting":      "fast",
				"measurement_range": map[string]float64{
					"min": 30.0,
					"max": 130.0,
				},
				"units": "dB(A)",
			}),
			CreatedAt: now,
		},
	}
}

func (s *AssetSensorSeeder) createAssetSensor(ctx context.Context, assetSensor *entity.AssetSensor) error {
	// Convert Configuration to string for JSONB
	var configStr string
	if assetSensor.Configuration != nil {
		configStr = string(assetSensor.Configuration)
	} else {
		configStr = "{}"
	}

	query := `
		INSERT INTO asset_sensors (
			id, tenant_id, asset_id, sensor_type_id, name, status,
			configuration, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7::jsonb, $8
		)`

	_, err := s.db.ExecContext(ctx, query,
		assetSensor.ID,
		assetSensor.TenantID,
		assetSensor.AssetID,
		assetSensor.SensorTypeID,
		assetSensor.Name,
		assetSensor.Status,
		configStr,
		assetSensor.CreatedAt,
	)

	return err
}

// Helper function to get asset sensor ID by name
func GetAssetSensorIDByName(db *sql.DB, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := db.QueryRow("SELECT id FROM asset_sensors WHERE name = $1", name).Scan(&id)
	return id, err
}

// Helper function to get all asset sensor IDs
func GetAllAssetSensorIDs(db *sql.DB) ([]uuid.UUID, error) {
	rows, err := db.Query("SELECT id FROM asset_sensors ORDER BY created_at")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// Helper function to get asset sensors by asset ID
func GetAssetSensorsByAssetID(db *sql.DB, assetID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := db.Query("SELECT id FROM asset_sensors WHERE asset_id = $1 ORDER BY created_at", assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// Helper function to get asset sensors by sensor type
func GetAssetSensorsBySensorTypeID(db *sql.DB, sensorTypeID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := db.Query("SELECT id FROM asset_sensors WHERE sensor_type_id = $1 ORDER BY created_at", sensorTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
