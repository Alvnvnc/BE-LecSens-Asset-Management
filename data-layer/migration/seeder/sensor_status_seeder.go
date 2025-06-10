package seeder

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type SensorStatusSeeder struct {
	db *sql.DB
}

func NewSensorStatusSeeder(db *sql.DB) *SensorStatusSeeder {
	return &SensorStatusSeeder{
		db: db,
	}
}

// Predefined Status IDs for consistent seeding
var (
	// Sensor Status IDs (matching asset sensor IDs from asset_sensor_seeder.go)
	StatusSensor01ID = uuid.MustParse("33333333-3333-3333-3333-000000000001") // pH Sensor 1
	StatusSensor02ID = uuid.MustParse("33333333-3333-3333-3333-000000000002") // pH Sensor 2
	StatusSensor03ID = uuid.MustParse("33333333-3333-3333-3333-000000000003") // Turbidity Sensor 1
	StatusSensor04ID = uuid.MustParse("33333333-3333-3333-3333-000000000004") // Turbidity Sensor 2
	StatusSensor05ID = uuid.MustParse("33333333-3333-3333-3333-000000000005") // PM2.5 Sensor 1
	StatusSensor06ID = uuid.MustParse("33333333-3333-3333-3333-000000000006") // PM2.5 Sensor 2
	StatusSensor07ID = uuid.MustParse("33333333-3333-3333-3333-000000000007") // Temperature Sensor 1
	StatusSensor08ID = uuid.MustParse("33333333-3333-3333-3333-000000000008") // Temperature Sensor 2
	StatusSensor09ID = uuid.MustParse("33333333-3333-3333-3333-000000000009") // Vibration Sensor 1
	StatusSensor10ID = uuid.MustParse("33333333-3333-3333-3333-000000000010") // Vibration Sensor 2
	StatusSensor11ID = uuid.MustParse("33333333-3333-3333-3333-000000000011") // Flow Sensor 1
	StatusSensor12ID = uuid.MustParse("33333333-3333-3333-3333-000000000012") // Flow Sensor 2
	StatusSensor13ID = uuid.MustParse("33333333-3333-3333-3333-000000000013") // Pressure Sensor 1
	StatusSensor14ID = uuid.MustParse("33333333-3333-3333-3333-000000000014") // Pressure Sensor 2
	StatusSensor15ID = uuid.MustParse("33333333-3333-3333-3333-000000000015") // Level Sensor 1
	StatusSensor16ID = uuid.MustParse("33333333-3333-3333-3333-000000000016") // Level Sensor 2
	StatusSensor17ID = uuid.MustParse("33333333-3333-3333-3333-000000000017") // Power Sensor 1
	StatusSensor18ID = uuid.MustParse("33333333-3333-3333-3333-000000000018") // Power Sensor 2
	StatusSensor19ID = uuid.MustParse("33333333-3333-3333-3333-000000000019") // Gas Detection Sensor 1
	StatusSensor20ID = uuid.MustParse("33333333-3333-3333-3333-000000000020") // Gas Detection Sensor 2
	StatusSensor21ID = uuid.MustParse("33333333-3333-3333-3333-000000000021") // CO2 Sensor 1
	StatusSensor22ID = uuid.MustParse("33333333-3333-3333-3333-000000000022") // CO2 Sensor 2
	StatusSensor23ID = uuid.MustParse("33333333-3333-3333-3333-000000000023") // DO Sensor 1
	StatusSensor24ID = uuid.MustParse("33333333-3333-3333-3333-000000000024") // DO Sensor 2
)

// getSensorStatusData returns all sensor status records to be seeded
func (s *SensorStatusSeeder) getSensorStatusData() []*entity.SensorStatus {
	now := time.Now()

	return []*entity.SensorStatus{
		// pH Sensors
		{
			ID:                   StatusSensor01ID,
			AssetSensorID:        AssetSensorPHProbeID, // From asset_sensor_seeder.go
			BatteryLevel:         floatPtr(85.5),
			BatteryVoltage:       floatPtr(3.7),
			BatteryStatus:        stringPtr("good"),
			BatteryType:          stringPtr("lithium"),
			BatteryEstimatedLife: intPtr(45),
			SignalType:           stringPtr("wifi"),
			SignalRSSI:           intPtr(-45),
			SignalSNR:            floatPtr(25.0),
			SignalQuality:        intPtr(85),
			SignalStatus:         stringPtr("good"),
			ConnectionType:       stringPtr("wifi"),
			ConnectionStatus:     "online",
			LastConnectedAt:      &now,
			CurrentIP:            stringPtr("192.168.1.101"),
			CurrentNetwork:       stringPtr("WaterTreatment-WiFi"),
			Temperature:          floatPtr(23.5),
			Humidity:             floatPtr(65.0),
			IsOnline:             true,
			LastHeartbeat:        &now,
			FirmwareVersion:      stringPtr("1.2.3"),
			ErrorCount:           intPtr(0),
			RecordedAt:           now,
			CreatedAt:            now,
		},
		{
			ID:                   StatusSensor02ID,
			AssetSensorID:        AssetSensorTurbidityID,
			BatteryLevel:         floatPtr(92.3),
			BatteryVoltage:       floatPtr(3.8),
			BatteryStatus:        stringPtr("good"),
			BatteryType:          stringPtr("lithium"),
			BatteryEstimatedLife: intPtr(52),
			SignalType:           stringPtr("wifi"),
			SignalRSSI:           intPtr(-38),
			SignalSNR:            floatPtr(28.5),
			SignalQuality:        intPtr(92),
			SignalStatus:         stringPtr("excellent"),
			ConnectionType:       stringPtr("wifi"),
			ConnectionStatus:     "online",
			LastConnectedAt:      &now,
			CurrentIP:            stringPtr("192.168.1.102"),
			CurrentNetwork:       stringPtr("WaterTreatment-WiFi"),
			Temperature:          floatPtr(24.1),
			Humidity:             floatPtr(63.2),
			IsOnline:             true,
			LastHeartbeat:        &now,
			FirmwareVersion:      stringPtr("1.2.3"),
			ErrorCount:           intPtr(0),
			RecordedAt:           now,
			CreatedAt:            now,
		},

		// Turbidity Sensors
		{
			ID:                   StatusSensor03ID,
			AssetSensorID:        AssetSensorDOProbeID,
			BatteryLevel:         floatPtr(78.9),
			BatteryVoltage:       floatPtr(3.6),
			BatteryStatus:        stringPtr("good"),
			BatteryType:          stringPtr("lithium"),
			BatteryEstimatedLife: intPtr(38),
			SignalType:           stringPtr("lora"),
			SignalRSSI:           intPtr(-95),
			SignalSNR:            floatPtr(8.5),
			SignalQuality:        intPtr(65),
			SignalStatus:         stringPtr("fair"),
			ConnectionType:       stringPtr("lora"),
			ConnectionStatus:     "online",
			LastConnectedAt:      &now,
			Temperature:          floatPtr(22.8),
			Humidity:             floatPtr(68.5),
			IsOnline:             true,
			LastHeartbeat:        &now,
			FirmwareVersion:      stringPtr("2.1.0"),
			ErrorCount:           intPtr(1),
			RecordedAt:           now,
			CreatedAt:            now,
		},
		{
			ID:                   StatusSensor04ID,
			AssetSensorID:        AssetSensorWaterConductID,
			BatteryLevel:         floatPtr(15.2),
			BatteryVoltage:       floatPtr(3.2),
			BatteryStatus:        stringPtr("low"),
			BatteryType:          stringPtr("lithium"),
			BatteryEstimatedLife: intPtr(5),
			SignalType:           stringPtr("lora"),
			SignalRSSI:           intPtr(-98),
			SignalSNR:            floatPtr(5.2),
			SignalQuality:        intPtr(45),
			SignalStatus:         stringPtr("poor"),
			ConnectionType:       stringPtr("lora"),
			ConnectionStatus:     "online",
			LastConnectedAt:      &now,
			Temperature:          floatPtr(25.3),
			Humidity:             floatPtr(71.2),
			IsOnline:             true,
			LastHeartbeat:        &now,
			FirmwareVersion:      stringPtr("2.1.0"),
			ErrorCount:           intPtr(3),
			RecordedAt:           now,
			CreatedAt:            now,
		},

		// PM2.5 Sensors
		{
			ID:                   StatusSensor05ID,
			AssetSensorID:        AssetSensorPM25ID,
			BatteryLevel:         floatPtr(95.7),
			BatteryVoltage:       floatPtr(3.9),
			BatteryStatus:        stringPtr("good"),
			BatteryType:          stringPtr("rechargeable"),
			BatteryEstimatedLife: intPtr(60),
			SignalType:           stringPtr("cellular"),
			SignalRSSI:           intPtr(-75),
			SignalSNR:            floatPtr(15.0),
			SignalQuality:        intPtr(75),
			SignalStatus:         stringPtr("good"),
			ConnectionType:       stringPtr("cellular"),
			ConnectionStatus:     "online",
			LastConnectedAt:      &now,
			CurrentIP:            stringPtr("10.0.0.15"),
			CurrentNetwork:       stringPtr("LTE-4G"),
			Temperature:          floatPtr(26.8),
			Humidity:             floatPtr(58.9),
			IsOnline:             true,
			LastHeartbeat:        &now,
			FirmwareVersion:      stringPtr("3.0.1"),
			ErrorCount:           intPtr(0),
			RecordedAt:           now,
			CreatedAt:            now,
		},
		{
			ID:                   StatusSensor06ID,
			AssetSensorID:        AssetSensorCO2ID,
			BatteryLevel:         floatPtr(8.3),
			BatteryVoltage:       floatPtr(3.0),
			BatteryStatus:        stringPtr("critical"),
			BatteryType:          stringPtr("rechargeable"),
			BatteryEstimatedLife: intPtr(1),
			SignalType:           stringPtr("cellular"),
			SignalRSSI:           intPtr(-105),
			SignalQuality:        intPtr(25),
			SignalStatus:         stringPtr("poor"),
			ConnectionType:       stringPtr("cellular"),
			ConnectionStatus:     "offline",
			LastDisconnectedAt:   &now,
			Temperature:          floatPtr(28.1),
			Humidity:             floatPtr(61.5),
			IsOnline:             false,
			FirmwareVersion:      stringPtr("3.0.1"),
			ErrorCount:           intPtr(12),
			LastErrorAt:          &now,
			RecordedAt:           now,
			CreatedAt:            now,
		},

		// Temperature Sensors
		{
			ID:                   StatusSensor07ID,
			AssetSensorID:        AssetSensorAirTempID,
			BatteryLevel:         floatPtr(88.4),
			BatteryVoltage:       floatPtr(3.7),
			BatteryStatus:        stringPtr("good"),
			BatteryType:          stringPtr("lithium"),
			BatteryEstimatedLife: intPtr(47),
			SignalType:           stringPtr("zigbee"),
			SignalRSSI:           intPtr(-65),
			SignalSNR:            floatPtr(18.0),
			SignalQuality:        intPtr(80),
			SignalStatus:         stringPtr("good"),
			ConnectionType:       stringPtr("zigbee"),
			ConnectionStatus:     "online",
			LastConnectedAt:      &now,
			Temperature:          floatPtr(24.5),
			Humidity:             floatPtr(60.0),
			IsOnline:             true,
			LastHeartbeat:        &now,
			FirmwareVersion:      stringPtr("1.5.2"),
			ErrorCount:           intPtr(0),
			RecordedAt:           now,
			CreatedAt:            now,
		},
		{
			ID:                   StatusSensor08ID,
			AssetSensorID:        AssetSensorAirHumidityID,
			BatteryLevel:         floatPtr(91.1),
			BatteryVoltage:       floatPtr(3.8),
			BatteryStatus:        stringPtr("good"),
			BatteryType:          stringPtr("lithium"),
			BatteryEstimatedLife: intPtr(51),
			SignalType:           stringPtr("zigbee"),
			SignalRSSI:           intPtr(-58),
			SignalSNR:            floatPtr(22.5),
			SignalQuality:        intPtr(85),
			SignalStatus:         stringPtr("good"),
			ConnectionType:       stringPtr("zigbee"),
			ConnectionStatus:     "online",
			LastConnectedAt:      &now,
			Temperature:          floatPtr(23.9),
			Humidity:             floatPtr(62.3),
			IsOnline:             true,
			LastHeartbeat:        &now,
			FirmwareVersion:      stringPtr("1.5.2"),
			ErrorCount:           intPtr(0),
			RecordedAt:           now,
			CreatedAt:            now,
		},

		// Add more sensors with different scenarios...
		// Vibration Sensors
		{
			ID:                   StatusSensor09ID,
			AssetSensorID:        AssetSensorColdStorageTempID,
			BatteryLevel:         floatPtr(72.6),
			BatteryVoltage:       floatPtr(3.5),
			BatteryStatus:        stringPtr("good"),
			BatteryType:          stringPtr("rechargeable"),
			BatteryEstimatedLife: intPtr(30),
			SignalType:           stringPtr("ethernet"),
			ConnectionType:       stringPtr("ethernet"),
			ConnectionStatus:     "online",
			LastConnectedAt:      &now,
			CurrentIP:            stringPtr("192.168.2.50"),
			CurrentNetwork:       stringPtr("Industrial-LAN"),
			Temperature:          floatPtr(35.2),
			Humidity:             floatPtr(45.8),
			IsOnline:             true,
			LastHeartbeat:        &now,
			FirmwareVersion:      stringPtr("4.1.0"),
			ErrorCount:           intPtr(0),
			RecordedAt:           now,
			CreatedAt:            now,
		},
		{
			ID:                   StatusSensor10ID,
			AssetSensorID:        AssetSensorAmbientTempID,
			BatteryLevel:         floatPtr(18.7),
			BatteryVoltage:       floatPtr(3.3),
			BatteryStatus:        stringPtr("low"),
			BatteryType:          stringPtr("rechargeable"),
			BatteryEstimatedLife: intPtr(7),
			SignalType:           stringPtr("ethernet"),
			ConnectionType:       stringPtr("ethernet"),
			ConnectionStatus:     "connecting",
			Temperature:          floatPtr(34.8),
			Humidity:             floatPtr(47.2),
			IsOnline:             false,
			FirmwareVersion:      stringPtr("4.1.0"),
			ErrorCount:           intPtr(2),
			RecordedAt:           now,
			CreatedAt:            now,
		},
	}
}

// statusExists checks if a sensor status exists
func (s *SensorStatusSeeder) statusExists(ctx context.Context, assetSensorID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sensor_status WHERE asset_sensor_id = $1)`
	var exists bool
	err := s.db.QueryRowContext(ctx, query, assetSensorID).Scan(&exists)
	return exists, err
}

// insertStatus inserts a sensor status into the database
func (s *SensorStatusSeeder) insertStatus(ctx context.Context, status *entity.SensorStatus) error {
	query := `
		INSERT INTO sensor_status (
			id, asset_sensor_id, battery_level, battery_voltage, battery_status,
			battery_last_charged, battery_estimated_life, battery_type, signal_type, signal_rssi,
			signal_snr, signal_quality, signal_frequency, signal_channel, signal_status,
			connection_type, connection_status, last_connected_at, last_disconnected_at,
			current_ip, current_network, temperature, humidity, is_online, last_heartbeat,
			firmware_version, error_count, last_error_at, recorded_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)
	`

	_, err := s.db.ExecContext(ctx, query,
		status.ID, status.AssetSensorID, status.BatteryLevel, status.BatteryVoltage, status.BatteryStatus,
		status.BatteryLastCharged, status.BatteryEstimatedLife, status.BatteryType, status.SignalType, status.SignalRSSI,
		status.SignalSNR, status.SignalQuality, status.SignalFrequency, status.SignalChannel, status.SignalStatus,
		status.ConnectionType, status.ConnectionStatus, status.LastConnectedAt, status.LastDisconnectedAt,
		status.CurrentIP, status.CurrentNetwork, status.Temperature, status.Humidity, status.IsOnline, status.LastHeartbeat,
		status.FirmwareVersion, status.ErrorCount, status.LastErrorAt, status.RecordedAt, status.CreatedAt,
	)

	return err
}

// Seed seeds all sensor status records
func (s *SensorStatusSeeder) Seed(ctx context.Context) error {
	log.Println("Starting sensor status seeding...")

	statuses := s.getSensorStatusData()
	createdCount := 0
	skippedCount := 0

	for _, status := range statuses {
		exists, err := s.statusExists(ctx, status.AssetSensorID)
		if err != nil {
			log.Printf("Error checking status existence for sensor %s: %v", status.AssetSensorID, err)
			continue
		}

		if exists {
			log.Printf("Status for sensor %s already exists, skipping", status.AssetSensorID)
			skippedCount++
			continue
		}

		err = s.insertStatus(ctx, status)
		if err != nil {
			log.Printf("Error creating status for sensor %s: %v", status.AssetSensorID, err)
			continue
		}

		log.Printf("Created sensor status for sensor: %s (Battery: %.1f%%, Signal: %s)", status.AssetSensorID, safeFloat(status.BatteryLevel), safeString(status.SignalStatus))
		createdCount++
	}

	log.Printf("Sensor status seeding completed. Created: %d, Skipped: %d", createdCount, skippedCount)
	return nil
}
