package seeder

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"

	"be-lecsens/asset_management/data-layer/entity"
)

// SensorLogsSeeder handles seeding of sensor logs data
type SensorLogsSeeder struct {
	db *sql.DB
}

// NewSensorLogsSeeder creates a new sensor logs seeder instance
func NewSensorLogsSeeder(db *sql.DB) *SensorLogsSeeder {
	return &SensorLogsSeeder{db: db}
}

// Seed inserts sample sensor logs data
func (s *SensorLogsSeeder) Seed() error {
	log.Println("Starting sensor logs seeding...")

	// Predefined sensor IDs from existing sensors
	sensorIDs := []uuid.UUID{
		AssetSensorPHProbeID,   // pH Sensor
		AssetSensorTurbidityID, // Turbidity Sensor
		AssetSensorPM25ID,      // PM2.5 Sensor
		AssetSensorCO2ID,       // CO2 Sensor
		AssetSensorAirTempID,   // Air Temperature Sensor
	}

	// Sample sensor logs data
	sensorLogs := []entity.SensorLogs{
		{
			ID:                 uuid.New(),
			TenantID:           GetDefaultTenantID(),
			AssetSensorID:      sensorIDs[0],
			LogType:            "reading",
			LogLevel:           "info",
			Message:            "Normal pH reading recorded",
			Component:          stringPtr("sensor"),
			EventType:          stringPtr("reading"),
			ErrorCode:          stringPtr(""),
			ConnectionType:     stringPtr("wifi"),
			ConnectionStatus:   stringPtr("connected"),
			IPAddress:          stringPtr("192.168.1.100"),
			MACAddress:         stringPtr("00:11:22:33:44:55"),
			NetworkName:        stringPtr("sensor_network"),
			ConnectionDuration: int64Ptr(3600),
			Metadata:           json.RawMessage(`{"reading": 7.2, "unit": "pH", "battery_level": 85}`),
			SourceIP:           stringPtr("192.168.1.1"),
			UserAgent:          stringPtr("SensorGateway/1.0"),
			SessionID:          stringPtr(uuid.New().String()),
			RecordedAt:         time.Now(),
			CreatedAt:          time.Now(),
			UpdatedAt:          timePtr(time.Now()),
		},
		{
			ID:                 uuid.New(),
			TenantID:           GetDefaultTenantID(),
			AssetSensorID:      sensorIDs[1],
			LogType:            "reading",
			LogLevel:           "warning",
			Message:            "pH reading slightly high",
			Component:          stringPtr("sensor"),
			EventType:          stringPtr("reading"),
			ErrorCode:          stringPtr("PH_HIGH"),
			ConnectionType:     stringPtr("wifi"),
			ConnectionStatus:   stringPtr("connected"),
			IPAddress:          stringPtr("192.168.1.101"),
			MACAddress:         stringPtr("00:11:22:33:44:56"),
			NetworkName:        stringPtr("sensor_network"),
			ConnectionDuration: int64Ptr(3600),
			Metadata:           json.RawMessage(`{"reading": 8.5, "unit": "pH", "battery_level": 90}`),
			SourceIP:           stringPtr("192.168.1.1"),
			UserAgent:          stringPtr("SensorGateway/1.0"),
			SessionID:          stringPtr(uuid.New().String()),
			RecordedAt:         time.Now(),
			CreatedAt:          time.Now(),
			UpdatedAt:          timePtr(time.Now()),
		},
		{
			ID:                 uuid.New(),
			TenantID:           GetDefaultTenantID(),
			AssetSensorID:      sensorIDs[2],
			LogType:            "reading",
			LogLevel:           "info",
			Message:            "Normal turbidity reading",
			Component:          stringPtr("sensor"),
			EventType:          stringPtr("reading"),
			ErrorCode:          stringPtr(""),
			ConnectionType:     stringPtr("wifi"),
			ConnectionStatus:   stringPtr("connected"),
			IPAddress:          stringPtr("192.168.1.102"),
			MACAddress:         stringPtr("00:11:22:33:44:57"),
			NetworkName:        stringPtr("sensor_network"),
			ConnectionDuration: int64Ptr(3600),
			Metadata:           json.RawMessage(`{"reading": 2.5, "unit": "NTU", "battery_level": 95}`),
			SourceIP:           stringPtr("192.168.1.1"),
			UserAgent:          stringPtr("SensorGateway/1.0"),
			SessionID:          stringPtr(uuid.New().String()),
			RecordedAt:         time.Now(),
			CreatedAt:          time.Now(),
			UpdatedAt:          timePtr(time.Now()),
		},
		// Connection log
		{
			ID:                 uuid.New(),
			TenantID:           GetDefaultTenantID(),
			AssetSensorID:      sensorIDs[3],
			LogType:            "connection",
			LogLevel:           "info",
			Message:            "Sensor connected successfully",
			Component:          stringPtr("network"),
			EventType:          stringPtr("connected"),
			ErrorCode:          stringPtr(""),
			ConnectionType:     stringPtr("wifi"),
			ConnectionStatus:   stringPtr("connected"),
			IPAddress:          stringPtr("192.168.1.103"),
			MACAddress:         stringPtr("00:11:22:33:44:58"),
			NetworkName:        stringPtr("sensor_network"),
			ConnectionDuration: int64Ptr(7200),
			Metadata:           json.RawMessage(`{"signal_strength": -45, "connection_time": "2025-06-10T10:00:00Z"}`),
			SourceIP:           stringPtr("192.168.1.1"),
			UserAgent:          stringPtr("SensorGateway/1.0"),
			SessionID:          stringPtr(uuid.New().String()),
			RecordedAt:         time.Now(),
			CreatedAt:          time.Now(),
			UpdatedAt:          timePtr(time.Now()),
		},
		// Error log
		{
			ID:                 uuid.New(),
			TenantID:           GetDefaultTenantID(),
			AssetSensorID:      sensorIDs[4],
			LogType:            "error",
			LogLevel:           "error",
			Message:            "Sensor calibration failed",
			Component:          stringPtr("sensor"),
			EventType:          stringPtr("error"),
			ErrorCode:          stringPtr("CAL_FAIL"),
			ConnectionType:     stringPtr("wifi"),
			ConnectionStatus:   stringPtr("connected"),
			IPAddress:          stringPtr("192.168.1.104"),
			MACAddress:         stringPtr("00:11:22:33:44:59"),
			NetworkName:        stringPtr("sensor_network"),
			ConnectionDuration: int64Ptr(1800),
			Metadata:           json.RawMessage(`{"last_calibration": "2025-06-01", "error_count": 3}`),
			SourceIP:           stringPtr("192.168.1.1"),
			UserAgent:          stringPtr("SensorGateway/1.0"),
			SessionID:          stringPtr(uuid.New().String()),
			RecordedAt:         time.Now(),
			CreatedAt:          time.Now(),
			UpdatedAt:          timePtr(time.Now()),
		},
	}

	// Insert logs
	for _, logEntry := range sensorLogs {
		query := `
			INSERT INTO sensor_logs (
				id, tenant_id, asset_sensor_id, log_type, log_level, message,
				component, event_type, error_code, connection_type, connection_status,
				ip_address, mac_address, network_name, connection_duration,
				metadata, source_ip, user_agent, session_id,
				recorded_at, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
				$16, $17, $18, $19, $20, $21, $22
			)
		`

		_, err := s.db.Exec(query,
			logEntry.ID, logEntry.TenantID, logEntry.AssetSensorID, logEntry.LogType, logEntry.LogLevel, logEntry.Message,
			logEntry.Component, logEntry.EventType, logEntry.ErrorCode, logEntry.ConnectionType, logEntry.ConnectionStatus,
			logEntry.IPAddress, logEntry.MACAddress, logEntry.NetworkName, logEntry.ConnectionDuration,
			logEntry.Metadata, logEntry.SourceIP, logEntry.UserAgent, logEntry.SessionID,
			logEntry.RecordedAt, logEntry.CreatedAt, logEntry.UpdatedAt,
		)

		if err != nil {
			log.Printf("Error inserting sensor log: %v", err)
			return err
		}
	}

	log.Printf("Successfully seeded %d sensor logs", len(sensorLogs))
	return nil
}
