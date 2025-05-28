package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AssetSensor represents a sensor attached to an asset
type AssetSensor struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	AssetID           uuid.UUID       `json:"asset_id"`
	Name              string          `json:"name"`
	SensorTypeID      uuid.UUID       `json:"sensor_type_id"`
	Status            string          `json:"status"`
	Configuration     json.RawMessage `json:"configuration,omitempty"` // Dynamic configuration based on sensor type
	LastReadingValue  *float64        `json:"last_reading_value,omitempty"`
	LastReadingTime   *time.Time      `json:"last_reading_time,omitempty"`
	LastReadingValues json.RawMessage `json:"last_reading_values,omitempty"` // Multiple readings for complex sensors
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         *time.Time      `json:"updated_at,omitempty"`
}

// NewAssetSensor creates a new asset sensor with default values
func NewAssetSensor() *AssetSensor {
	now := time.Now()
	return &AssetSensor{
		ID:        uuid.New(),
		Status:    "active",
		CreatedAt: now,
	}
}

// SetConfiguration sets the dynamic configuration for the sensor based on SensorMeasurementType schema
func (s *AssetSensor) SetConfiguration(config map[string]interface{}) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	s.Configuration = configJSON
	return nil
}

// GetConfiguration parses the sensor configuration into a map
func (s *AssetSensor) GetConfiguration() (map[string]interface{}, error) {
	if s.Configuration == nil {
		return map[string]interface{}{}, nil
	}

	var config map[string]interface{}
	if err := json.Unmarshal(s.Configuration, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// UpdateLastReadingValues updates the sensor's last reading values with multiple measurements
func (s *AssetSensor) UpdateLastReadingValues(readings map[string]interface{}) error {
	readingsJSON, err := json.Marshal(readings)
	if err != nil {
		return err
	}

	s.LastReadingValues = readingsJSON
	now := time.Now()
	s.LastReadingTime = &now

	// If a primary reading value is provided, update the simple LastReadingValue field too
	if primaryValue, ok := readings["primary"].(float64); ok {
		s.LastReadingValue = &primaryValue
	}

	return nil
}

// GetLastReadingValues parses the last reading values into a map
func (s *AssetSensor) GetLastReadingValues() (map[string]interface{}, error) {
	if s.LastReadingValues == nil {
		return map[string]interface{}{}, nil
	}

	var values map[string]interface{}
	if err := json.Unmarshal(s.LastReadingValues, &values); err != nil {
		return nil, err
	}
	return values, nil
}
