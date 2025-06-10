package seeder

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type AssetSeeder struct {
	repo repository.AssetRepository
}

func NewAssetSeeder(db *sql.DB) *AssetSeeder {
	return &AssetSeeder{
		repo: repository.NewAssetRepository(db),
	}
}

// Predefined asset UUIDs for consistency across seeders
var (
	AssetProdLineAID       = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	AssetWaterTreatmentID  = uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")
	AssetAirQualityHubID   = uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")
	AssetTempMonitoringID  = uuid.MustParse("550e8400-e29b-41d4-a716-446655440004")
	AssetDataCenterHubID   = uuid.MustParse("550e8400-e29b-41d4-a716-446655440005")
	AssetVibrationHubID    = uuid.MustParse("550e8400-e29b-41d4-a716-446655440006")
	AssetFlowStationID     = uuid.MustParse("550e8400-e29b-41d4-a716-446655440007")
	AssetPressureHubID     = uuid.MustParse("550e8400-e29b-41d4-a716-446655440008")
	AssetLevelMonitorID    = uuid.MustParse("550e8400-e29b-41d4-a716-446655440009")
	AssetEnergyMeterID     = uuid.MustParse("550e8400-e29b-41d4-a716-446655440010")
	AssetPowerDistribution = uuid.MustParse("550e8400-e29b-41d4-a716-446655440011")
	AssetWarehouseHubID    = uuid.MustParse("550e8400-e29b-41d4-a716-446655440012")
)

// Location constants - should match the IDs in location_seeder.go
var (
	LocationJakartaPusatID   = uuid.MustParse("01234567-1111-1111-1111-000000000001") // Jakarta Pusat
	LocationJakartaSelatanID = uuid.MustParse("01234567-1111-1111-1111-000000000002") // Jakarta Selatan
	LocationJakartaUtaraID   = uuid.MustParse("01234567-1111-1111-1111-000000000003") // Jakarta Utara
	LocationBandungID        = uuid.MustParse("01234567-1111-1111-1111-000000000004") // Bandung
	LocationSurabayaID       = uuid.MustParse("01234567-1111-1111-1111-000000000005") // Surabaya
	LocationYogyakartaID     = uuid.MustParse("01234567-1111-1111-1111-000000000006") // Yogyakarta
	LocationMalangID         = uuid.MustParse("01234567-1111-1111-1111-000000000007") // Malang
	LocationMedanID          = uuid.MustParse("01234567-1111-1111-1111-000000000008") // Medan
	LocationDenpasarID       = uuid.MustParse("01234567-1111-1111-1111-000000000009") // Denpasar
	LocationMakassarID       = uuid.MustParse("01234567-1111-1111-1111-000000000010") // Makassar
)

// Asset Type constants - should match the IDs in asset_type_seeder.go
var (
	AssetTypeWaterQualityID   = uuid.MustParse("02345678-2222-2222-2222-000000000001")
	AssetTypeWaterFlowID      = uuid.MustParse("02345678-2222-2222-2222-000000000002")
	AssetTypeWaterLevelID     = uuid.MustParse("02345678-2222-2222-2222-000000000003")
	AssetTypeWaterPressureID  = uuid.MustParse("02345678-2222-2222-2222-000000000004")
	AssetTypeWaterDroneID     = uuid.MustParse("02345678-2222-2222-2222-000000000005")
	AssetTypeWaterVehicleID   = uuid.MustParse("02345678-2222-2222-2222-000000000006")
	AssetTypeWaterTreatmentID = uuid.MustParse("02345678-2222-2222-2222-000000000007")
)

// GetDefaultTenantID returns a default tenant ID for seeding
func GetDefaultTenantID() *uuid.UUID {
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	return &tenantID
}

// Helper function to create json.RawMessage from interface{}
func jsonRawMessage(data interface{}) json.RawMessage {
	bytes, _ := json.Marshal(data)
	return json.RawMessage(bytes)
}

func (s *AssetSeeder) Seed(ctx context.Context) error {
	log.Println("Starting Asset seeder...")

	// Check if assets already exist using List method with empty tenant
	assets, err := s.repo.List(ctx, nil, 1, 1)
	if err != nil {
		return fmt.Errorf("failed to check existing assets: %w", err)
	}

	if len(assets) > 0 {
		log.Printf("Assets already exist (%d records), skipping seed", len(assets))
		return nil
	}

	tenantID := GetDefaultTenantID()
	assetData := s.getAssetData(tenantID)

	// Insert assets
	for _, asset := range assetData {
		err := s.repo.Create(ctx, asset)
		if err != nil {
			log.Printf("Failed to create asset %s: %v", asset.Name, err)
			return err
		}
		log.Printf("Created asset: %s", asset.Name)
	}

	log.Printf("Successfully seeded %d assets", len(assetData))
	return nil
}

func (s *AssetSeeder) getAssetData(tenantID *uuid.UUID) []*entity.Asset {
	now := time.Now()

	// Configuration templates for different asset types with realistic properties
	waterQualityProperties := map[string]interface{}{
		"description":          "Primary water quality monitoring system",
		"serial_number":        "WQ-PROD-A-001",
		"monitoring_interval":  300,
		"calibration_schedule": "monthly",
		"maintenance_schedule": "quarterly",
		"alert_thresholds": map[string]interface{}{
			"ph_min":        6.5,
			"ph_max":        8.5,
			"turbidity_max": 5.0,
			"do_min":        4.0,
		},
		"location_coordinates": map[string]float64{
			"latitude":  -6.2088,
			"longitude": 106.8456,
		},
	}

	waterFlowProperties := map[string]interface{}{
		"description":         "Water flow monitoring system",
		"serial_number":       "WF-PIPE-001",
		"monitoring_interval": 60,
		"pipe_diameter":       12,
		"measurement_type":    "ultrasonic",
		"alert_thresholds": map[string]interface{}{
			"flow_min": 0.5,
			"flow_max": 100.0,
		},
	}

	waterLevelProperties := map[string]interface{}{
		"description":         "Water level monitoring system",
		"serial_number":       "WL-TANK-001",
		"monitoring_interval": 300,
		"tank_capacity":       1000.0,
		"technology":          "ultrasonic",
		"alert_thresholds": map[string]interface{}{
			"level_min": 10.0,
			"level_max": 90.0,
		},
	}

	waterPressureProperties := map[string]interface{}{
		"description":         "Water pressure monitoring system",
		"serial_number":       "WP-PIPE-001",
		"monitoring_interval": 60,
		"pressure_type":       "gauge",
		"alert_thresholds": map[string]interface{}{
			"pressure_min": 0.1,
			"pressure_max": 10.0,
		},
	}

	waterDroneProperties := map[string]interface{}{
		"description":          "Water sampling drone",
		"serial_number":        "WD-001",
		"flight_time":          45,
		"sampling_capacity":    500,
		"waterproof_rating":    "IP68",
		"maintenance_schedule": "weekly",
	}

	waterVehicleProperties := map[string]interface{}{
		"description":        "Mobile water quality monitoring vehicle",
		"serial_number":      "WV-001",
		"vehicle_type":       "van",
		"equipment_capacity": 500,
		"operating_range":    200,
		"crew_capacity":      3,
	}

	waterTreatmentProperties := map[string]interface{}{
		"description":          "Water treatment system",
		"serial_number":        "WT-001",
		"treatment_capacity":   5000,
		"treatment_methods":    []string{"filtration", "uv", "chemical"},
		"automation_level":     "full_auto",
		"maintenance_schedule": "monthly",
	}

	return []*entity.Asset{
		// Water Quality Monitoring
		{
			ID:          AssetProdLineAID,
			TenantID:    tenantID,
			Name:        "Main Water Treatment Plant Monitor",
			AssetTypeID: AssetTypeWaterQualityID,
			LocationID:  LocationJakartaPusatID,
			Status:      "active",
			Properties:  jsonRawMessage(waterQualityProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          AssetWaterTreatmentID,
			TenantID:    tenantID,
			Name:        "Secondary Water Quality Monitor",
			AssetTypeID: AssetTypeWaterQualityID,
			LocationID:  LocationSurabayaID,
			Status:      "active",
			Properties:  jsonRawMessage(waterQualityProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},

		// Water Flow Monitoring
		{
			ID:          AssetAirQualityHubID,
			TenantID:    tenantID,
			Name:        "Main Pipeline Flow Monitor",
			AssetTypeID: AssetTypeWaterFlowID,
			LocationID:  LocationBandungID,
			Status:      "active",
			Properties:  jsonRawMessage(waterFlowProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},

		// Water Level Monitoring
		{
			ID:          AssetTempMonitoringID,
			TenantID:    tenantID,
			Name:        "Reservoir Level Monitor",
			AssetTypeID: AssetTypeWaterLevelID,
			LocationID:  LocationMedanID,
			Status:      "active",
			Properties:  jsonRawMessage(waterLevelProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          AssetDataCenterHubID,
			TenantID:    tenantID,
			Name:        "Storage Tank Level Monitor",
			AssetTypeID: AssetTypeWaterLevelID,
			LocationID:  LocationYogyakartaID,
			Status:      "active",
			Properties:  jsonRawMessage(waterLevelProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},

		// Water Pressure Monitoring
		{
			ID:          AssetVibrationHubID,
			TenantID:    tenantID,
			Name:        "Distribution Network Pressure Monitor",
			AssetTypeID: AssetTypeWaterPressureID,
			LocationID:  LocationMalangID,
			Status:      "active",
			Properties:  jsonRawMessage(waterPressureProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},

		// Water Sampling Drone
		{
			ID:          AssetFlowStationID,
			TenantID:    tenantID,
			Name:        "Water Quality Sampling Drone",
			AssetTypeID: AssetTypeWaterDroneID,
			LocationID:  LocationJakartaSelatanID,
			Status:      "active",
			Properties:  jsonRawMessage(waterDroneProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},

		// Water Quality Vehicle
		{
			ID:          AssetPressureHubID,
			TenantID:    tenantID,
			Name:        "Mobile Water Quality Lab",
			AssetTypeID: AssetTypeWaterVehicleID,
			LocationID:  LocationMakassarID,
			Status:      "active",
			Properties:  jsonRawMessage(waterVehicleProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},

		// Water Treatment Equipment
		{
			ID:          AssetLevelMonitorID,
			TenantID:    tenantID,
			Name:        "Water Purification System",
			AssetTypeID: AssetTypeWaterTreatmentID,
			LocationID:  LocationJakartaUtaraID,
			Status:      "active",
			Properties:  jsonRawMessage(waterTreatmentProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},

		// Additional Water Quality Monitor
		{
			ID:          AssetEnergyMeterID,
			TenantID:    tenantID,
			Name:        "Industrial Water Quality Monitor",
			AssetTypeID: AssetTypeWaterQualityID,
			LocationID:  LocationDenpasarID,
			Status:      "active",
			Properties:  jsonRawMessage(waterQualityProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},

		// Additional Water Flow Monitor
		{
			ID:          AssetPowerDistribution,
			TenantID:    tenantID,
			Name:        "Secondary Pipeline Flow Monitor",
			AssetTypeID: AssetTypeWaterFlowID,
			LocationID:  LocationJakartaPusatID,
			Status:      "active",
			Properties:  jsonRawMessage(waterFlowProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},

		// Additional Water Treatment Equipment
		{
			ID:          AssetWarehouseHubID,
			TenantID:    tenantID,
			Name:        "Backup Water Treatment System",
			AssetTypeID: AssetTypeWaterTreatmentID,
			LocationID:  LocationJakartaUtaraID,
			Status:      "maintenance",
			Properties:  jsonRawMessage(waterTreatmentProperties),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
}

// Helper functions for external access to asset data

// GetAssetIDByName returns the asset ID for a given asset name
func GetAssetIDByName(name string) (uuid.UUID, bool) {
	assetMap := map[string]uuid.UUID{
		"Production Line A - Water Quality Monitor": AssetProdLineAID,
		"Water Treatment Plant Monitor":             AssetWaterTreatmentID,
		"Factory Air Quality Hub":                   AssetAirQualityHubID,
		"Cold Storage Temperature Monitor":          AssetTempMonitoringID,
		"Data Center Environmental Monitor":         AssetDataCenterHubID,
		"Turbine Vibration Monitoring Hub":          AssetVibrationHubID,
		"Pipeline Flow Monitoring Station":          AssetFlowStationID,
		"Chemical Plant Pressure Hub":               AssetPressureHubID,
		"Fuel Tank Level Monitor":                   AssetLevelMonitorID,
		"Resort Energy Management System":           AssetEnergyMeterID,
		"Main Power Distribution Monitor":           AssetPowerDistribution,
		"Warehouse Environmental Hub":               AssetWarehouseHubID,
	}

	id, exists := assetMap[name]
	return id, exists
}

// GetAllAssetIDs returns all predefined asset IDs
func GetAllAssetIDs() []uuid.UUID {
	return []uuid.UUID{
		AssetProdLineAID,
		AssetWaterTreatmentID,
		AssetAirQualityHubID,
		AssetTempMonitoringID,
		AssetDataCenterHubID,
		AssetVibrationHubID,
		AssetFlowStationID,
		AssetPressureHubID,
		AssetLevelMonitorID,
		AssetEnergyMeterID,
		AssetPowerDistribution,
		AssetWarehouseHubID,
	}
}

// GetAssetsByLocationID returns asset IDs for a specific location
func GetAssetsByLocationID(locationID uuid.UUID) []uuid.UUID {
	locationAssets := map[uuid.UUID][]uuid.UUID{
		LocationJakartaPusatID:   {AssetProdLineAID, AssetPowerDistribution},
		LocationSurabayaID:       {AssetWaterTreatmentID, AssetWarehouseHubID},
		LocationBandungID:        {AssetAirQualityHubID},
		LocationMedanID:          {AssetTempMonitoringID},
		LocationYogyakartaID:     {AssetDataCenterHubID},
		LocationMalangID:         {AssetVibrationHubID},
		LocationJakartaSelatanID: {AssetFlowStationID},
		LocationMakassarID:       {AssetPressureHubID},
		LocationJakartaUtaraID:   {AssetLevelMonitorID},
		LocationDenpasarID:       {AssetEnergyMeterID},
	}

	return locationAssets[locationID]
}

// GetAssetsByTypeID returns asset IDs for a specific asset type
func GetAssetsByTypeID(assetTypeID uuid.UUID) []uuid.UUID {
	typeAssets := map[uuid.UUID][]uuid.UUID{
		AssetTypeWaterQualityID:   {AssetProdLineAID, AssetWaterTreatmentID},
		AssetTypeWaterFlowID:      {AssetAirQualityHubID},
		AssetTypeWaterLevelID:     {AssetTempMonitoringID, AssetDataCenterHubID},
		AssetTypeWaterPressureID:  {AssetVibrationHubID},
		AssetTypeWaterDroneID:     {AssetFlowStationID},
		AssetTypeWaterVehicleID:   {AssetPressureHubID},
		AssetTypeWaterTreatmentID: {AssetLevelMonitorID},
	}

	return typeAssets[assetTypeID]
}
