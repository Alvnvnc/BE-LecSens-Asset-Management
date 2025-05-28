package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// IoTSensorReading represents complex measurement data recorded by IoT devices with ML processing
type IoTSensorReading struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	AssetSensorID uuid.UUID `json:"asset_sensor_id"`
	SensorTypeID  uuid.UUID `json:"sensor_type_id"`
	MacAddress    string    `json:"mac_address"`
	Location      string    `json:"location"`
	// Dynamic measurement fields - this is the primary data structure replacing all specific fields
	MeasurementData json.RawMessage `json:"measurement_data"` // Flexible structure based on sensor type
	// Visualization-ready data extracted from MeasurementData
	StandardFields json.RawMessage `json:"standard_fields,omitempty"` // Common fields for visualizations
	ReadingTime    time.Time       `json:"reading_time"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      *time.Time      `json:"updated_at,omitempty"`

	// Deprecated fields - will be removed in future versions
	// These are only kept to ensure backward compatibility during transition
	DataX   json.RawMessage `json:"data_x,omitempty"`
	DataY   json.RawMessage `json:"data_y,omitempty"`
	PeakX   json.RawMessage `json:"peak_x,omitempty"`
	PeakY   json.RawMessage `json:"peak_y,omitempty"`
	PPM     float64         `json:"ppm,omitempty"`
	Label   string          `json:"label,omitempty"`
	RawData json.RawMessage `json:"raw_data,omitempty"`
}

// NewIoTSensorReading creates a new IoT sensor reading with default values
func NewIoTSensorReading() *IoTSensorReading {
	now := time.Now()
	return &IoTSensorReading{
		ID:          uuid.New(),
		ReadingTime: now,
		CreatedAt:   now,
	}
}

// IoTSensorReadingFromJSON creates an IoT sensor reading from a JSON payload
// This is useful for directly unmarshalling API request data
func IoTSensorReadingFromJSON(data []byte) (*IoTSensorReading, error) {
	var reading struct {
		MacAddress      string          `json:"mac_address"`
		Location        string          `json:"location"`
		MeasurementData json.RawMessage `json:"measurement_data"`
		SensorTypeID    string          `json:"sensor_type_id"`
		// Legacy fields (kept only for backward compatibility)
		DataX   json.RawMessage `json:"data_x,omitempty"`
		DataY   json.RawMessage `json:"data_y,omitempty"`
		PeakX   json.RawMessage `json:"peak_x,omitempty"`
		PeakY   json.RawMessage `json:"peak_y,omitempty"`
		PPM     float64         `json:"ppm,omitempty"`
		Label   string          `json:"label,omitempty"`
		RawData json.RawMessage `json:"raw_data,omitempty"`
	}

	if err := json.Unmarshal(data, &reading); err != nil {
		return nil, err
	}

	// Parse the UUID
	sensorTypeUUID, err := uuid.Parse(reading.SensorTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid sensor_type_id: %v", err)
	}

	now := time.Now()
	result := &IoTSensorReading{
		ID:           uuid.New(),
		SensorTypeID: sensorTypeUUID,
		MacAddress:   reading.MacAddress,
		Location:     reading.Location,
		ReadingTime:  now,
		CreatedAt:    now,
	}

	// Handle measurement data properly
	if reading.MeasurementData != nil {
		result.MeasurementData = reading.MeasurementData

		// Extract standard fields for visualization
		var measurements map[string]interface{}
		if err := json.Unmarshal(reading.MeasurementData, &measurements); err == nil {
			// Create a standardized set of fields for visualization
			standardFields := make(map[string]interface{})

			// Extract common fields if they exist
			if xValues, ok := measurements["x_values"]; ok {
				standardFields["x_values"] = xValues
			}
			if yValues, ok := measurements["y_values"]; ok {
				standardFields["y_values"] = yValues
			}
			if peaks, ok := measurements["peaks"]; ok {
				standardFields["peaks"] = peaks
			}
			if value, ok := measurements["value"]; ok {
				standardFields["value"] = value
			}

			// Set standard fields
			if len(standardFields) > 0 {
				standardFieldsJSON, err := json.Marshal(standardFields)
				if err == nil {
					result.StandardFields = standardFieldsJSON
				}
			}
		}
	} else {
		// Handle legacy data format by converting to MeasurementData format
		// This ensures backward compatibility during transition
		measurementMap := make(map[string]interface{})

		// Convert legacy fields to the new measurement data structure
		if reading.DataX != nil {
			var dataX []interface{}
			if json.Unmarshal(reading.DataX, &dataX) == nil {
				measurementMap["x_values"] = dataX
				result.DataX = reading.DataX // Keep legacy field populated
			}
		}

		if reading.DataY != nil {
			var dataY []interface{}
			if json.Unmarshal(reading.DataY, &dataY) == nil {
				measurementMap["y_values"] = dataY
				result.DataY = reading.DataY // Keep legacy field populated
			}
		}

		if reading.PeakX != nil || reading.PeakY != nil {
			peakData := make(map[string]interface{})

			var peakX []interface{}
			if reading.PeakX != nil && json.Unmarshal(reading.PeakX, &peakX) == nil {
				peakData["x"] = peakX
				result.PeakX = reading.PeakX
			}

			var peakY []interface{}
			if reading.PeakY != nil && json.Unmarshal(reading.PeakY, &peakY) == nil {
				peakData["y"] = peakY
				result.PeakY = reading.PeakY
			}

			if len(peakData) > 0 {
				measurementMap["peaks"] = peakData
			}
		}

		if reading.PPM > 0 {
			measurementMap["ppm"] = reading.PPM
			result.PPM = reading.PPM
		}

		if reading.Label != "" {
			measurementMap["label"] = reading.Label
			result.Label = reading.Label
		}

		if reading.RawData != nil {
			var rawData interface{}
			if json.Unmarshal(reading.RawData, &rawData) == nil {
				measurementMap["raw_data"] = rawData
				result.RawData = reading.RawData
			}
		}

		// Convert the consolidated measurement map to JSON
		if len(measurementMap) > 0 {
			measurementJSON, err := json.Marshal(measurementMap)
			if err == nil {
				result.MeasurementData = measurementJSON

				// Also set StandardFields for visualization
				visualizationFields := make(map[string]interface{})
				if xValues, ok := measurementMap["x_values"]; ok {
					visualizationFields["x_values"] = xValues
				}
				if yValues, ok := measurementMap["y_values"]; ok {
					visualizationFields["y_values"] = yValues
				}
				if value, ok := measurementMap["ppm"]; ok {
					visualizationFields["value"] = value
				}

				if len(visualizationFields) > 0 {
					standardFieldsJSON, err := json.Marshal(visualizationFields)
					if err == nil {
						result.StandardFields = standardFieldsJSON
					}
				}
			}
		}
	}

	return result, nil
}

// GetMeasurementData parses the measurement data into a map
func (r *IoTSensorReading) GetMeasurementData() (map[string]interface{}, error) {
	if r.MeasurementData == nil {
		return map[string]interface{}{}, nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(r.MeasurementData, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// SetMeasurementData sets measurement data from a map
func (r *IoTSensorReading) SetMeasurementData(data map[string]interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.MeasurementData = dataJSON
	return nil
}

// ValidateAgainstMeasurementType validates the reading against a measurement type schema
func (r *IoTSensorReading) ValidateAgainstMeasurementType(measurementType *SensorMeasurementType) (bool, []string) {
	if r.MeasurementData == nil {
		return false, []string{"Measurement data is missing"}
	}

	return measurementType.ValidateMeasurement(r.MeasurementData)
}
