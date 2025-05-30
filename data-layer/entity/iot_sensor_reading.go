package entity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// IoTSensorReading represents the main IoT sensor reading entity
// This is an alias for IoTSensorReadingFlexible for backward compatibility
type IoTSensorReading = IoTSensorReadingFlexible

// IoTSensorReadingFlexible represents a single measurement record with flexible data types
type IoTSensorReadingFlexible struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	TenantID      *uuid.UUID `json:"tenant_id" db:"tenant_id"`
	AssetSensorID uuid.UUID  `json:"asset_sensor_id" db:"asset_sensor_id"`
	SensorTypeID  uuid.UUID  `json:"sensor_type_id" db:"sensor_type_id"`
	MacAddress    *string    `json:"mac_address" db:"mac_address"`
	LocationID    *uuid.UUID `json:"location_id" db:"location_id"`
	LocationName  *string    `json:"location_name" db:"location_name"`

	// Measurement identification
	MeasurementType  string  `json:"measurement_type" db:"measurement_type"`   // 'raw_value', 'temperature', etc.
	MeasurementLabel *string `json:"measurement_label" db:"measurement_label"` // Human readable label
	MeasurementUnit  *string `json:"measurement_unit" db:"measurement_unit"`   // °C, μg/m³, %, etc.

	// Flexible value storage (only one should be non-null per record)
	NumericValue *float64 `json:"numeric_value" db:"numeric_value"`
	TextValue    *string  `json:"text_value" db:"text_value"`
	BooleanValue *bool    `json:"boolean_value" db:"boolean_value"`

	// Additional metadata
	DataSource        *string `json:"data_source" db:"data_source"`                 // 'json', 'text', 'csv'
	OriginalFieldName *string `json:"original_field_name" db:"original_field_name"` // Original field name

	ReadingTime time.Time  `json:"reading_time" db:"reading_time"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
}

// MeasurementData represents a single measurement value from flexible JSON
type MeasurementData struct {
	Type     string      `json:"type"`      // Field name (raw_value, temperature, etc.)
	Label    *string     `json:"label"`     // Human readable label
	Unit     *string     `json:"unit"`      // Unit of measurement
	Value    interface{} `json:"value"`     // The actual value (can be number, string, or boolean)
	DataType string      `json:"data_type"` // "numeric", "text", "boolean"
}

// FlexibleReadingRequest represents the request format for flexible IoT sensor readings
type FlexibleReadingRequest struct {
	AssetSensorID uuid.UUID              `json:"asset_sensor_id"`
	SensorTypeID  uuid.UUID              `json:"sensor_type_id"`
	MacAddress    *string                `json:"mac_address"`
	LocationID    *uuid.UUID             `json:"location_id"`
	LocationName  *string                `json:"location_name"`
	ReadingTime   *time.Time             `json:"reading_time"`
	Measurements  map[string]interface{} `json:"-"`           // Will be populated by custom unmarshaling
	DataSource    string                 `json:"data_source"` // 'json', 'text', 'csv'
}

// UnmarshalJSON custom unmarshaling for FlexibleReadingRequest
func (r *FlexibleReadingRequest) UnmarshalJSON(data []byte) error {
	// First unmarshal into a temporary map to extract all fields
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Extract known fields
	if val, ok := temp["asset_sensor_id"]; ok {
		if str, ok := val.(string); ok {
			if parsed, err := uuid.Parse(str); err == nil {
				r.AssetSensorID = parsed
			}
		}
	}

	if val, ok := temp["sensor_type_id"]; ok {
		if str, ok := val.(string); ok {
			if parsed, err := uuid.Parse(str); err == nil {
				r.SensorTypeID = parsed
			}
		}
	}

	if val, ok := temp["mac_address"]; ok {
		if str, ok := val.(string); ok {
			r.MacAddress = &str
		}
	}

	if val, ok := temp["location_id"]; ok {
		if str, ok := val.(string); ok {
			if parsed, err := uuid.Parse(str); err == nil {
				r.LocationID = &parsed
			}
		}
	}

	if val, ok := temp["location_name"]; ok {
		if str, ok := val.(string); ok {
			r.LocationName = &str
		}
	}

	if val, ok := temp["reading_time"]; ok {
		if str, ok := val.(string); ok {
			if parsed, err := time.Parse(time.RFC3339, str); err == nil {
				r.ReadingTime = &parsed
			}
		}
	}

	if val, ok := temp["data_source"]; ok {
		if str, ok := val.(string); ok {
			r.DataSource = str
		}
	}

	// Extract measurement fields (everything else that's not a known field)
	knownFields := map[string]bool{
		"asset_sensor_id": true,
		"sensor_type_id":  true,
		"mac_address":     true,
		"location_id":     true,
		"location_name":   true,
		"reading_time":    true,
		"data_source":     true,
	}

	r.Measurements = make(map[string]interface{})
	for key, value := range temp {
		if !knownFields[key] {
			r.Measurements[key] = value
		}
	}

	return nil
}

// GetValue returns the appropriate value based on the data type
func (r *IoTSensorReadingFlexible) GetValue() interface{} {
	if r.NumericValue != nil {
		return *r.NumericValue
	}
	if r.TextValue != nil {
		return *r.TextValue
	}
	if r.BooleanValue != nil {
		return *r.BooleanValue
	}
	return nil
}

// SetValue sets the appropriate value field based on the input type
func (r *IoTSensorReadingFlexible) SetValue(value interface{}) error {
	// Reset all values first
	r.NumericValue = nil
	r.TextValue = nil
	r.BooleanValue = nil

	switch v := value.(type) {
	case float64:
		r.NumericValue = &v
	case int:
		f := float64(v)
		r.NumericValue = &f
	case int64:
		f := float64(v)
		r.NumericValue = &f
	case string:
		// Try to parse as number first
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			r.NumericValue = &f
		} else if b, err := strconv.ParseBool(v); err == nil {
			r.BooleanValue = &b
		} else {
			r.TextValue = &v
		}
	case bool:
		r.BooleanValue = &v
	default:
		// Convert to string as fallback
		str := fmt.Sprintf("%v", v)
		r.TextValue = &str
	}

	return nil
}

// ParseMeasurementData parses measurement data from interface{} into structured format
func ParseMeasurementData(fieldName string, value interface{}) (*MeasurementData, error) {
	measurement := &MeasurementData{
		Type: fieldName,
	}

	// Handle different input formats
	switch v := value.(type) {
	case map[string]interface{}:
		// Handle structured measurement: {"unit": "°C", "label": "Temperature", "value": 25.3}
		if unit, ok := v["unit"].(string); ok {
			measurement.Unit = &unit
		}
		if label, ok := v["label"].(string); ok {
			measurement.Label = &label
		}
		if val, ok := v["value"]; ok {
			measurement.Value = val
		} else {
			measurement.Value = v
		}
	default:
		// Handle direct value
		measurement.Value = v
	}

	// Determine data type
	switch measurement.Value.(type) {
	case float64, int, int64:
		measurement.DataType = "numeric"
	case bool:
		measurement.DataType = "boolean"
	default:
		measurement.DataType = "text"
	}

	return measurement, nil
}

// ConvertToFlexibleReadings converts a FlexibleReadingRequest to multiple IoTSensorReadingFlexible records
func ConvertToFlexibleReadings(req *FlexibleReadingRequest) ([]*IoTSensorReadingFlexible, error) {
	var readings []*IoTSensorReadingFlexible

	readingTime := time.Now()
	if req.ReadingTime != nil {
		readingTime = *req.ReadingTime
	}

	dataSource := req.DataSource
	if dataSource == "" {
		dataSource = "json"
	}

	for fieldName, value := range req.Measurements {
		measurement, err := ParseMeasurementData(fieldName, value)
		if err != nil {
			continue // Skip invalid measurements
		}

		reading := &IoTSensorReadingFlexible{
			ID:                uuid.New(),
			AssetSensorID:     req.AssetSensorID,
			SensorTypeID:      req.SensorTypeID,
			MacAddress:        req.MacAddress,
			LocationID:        req.LocationID,
			LocationName:      req.LocationName,
			MeasurementType:   measurement.Type,
			MeasurementLabel:  measurement.Label,
			MeasurementUnit:   measurement.Unit,
			DataSource:        &dataSource,
			OriginalFieldName: &fieldName,
			ReadingTime:       readingTime,
			CreatedAt:         time.Now(),
		}

		// Set the appropriate value
		if err := reading.SetValue(measurement.Value); err != nil {
			continue // Skip if value setting fails
		}

		readings = append(readings, reading)
	}

	return readings, nil
}

// ParseTextToFlexibleReadings parses text input into flexible readings
func ParseTextToFlexibleReadings(textInput string, assetSensorID, sensorTypeID uuid.UUID) ([]*IoTSensorReadingFlexible, error) {
	// Auto-detect format and parse
	measurements := make(map[string]interface{})

	lines := strings.Split(strings.TrimSpace(textInput), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try different parsing patterns
		if strings.Contains(line, ":") {
			// Key-value format: "temperature: 25.3°C"
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				valueStr := strings.TrimSpace(parts[1])

				// Try to extract numeric value and unit
				value := parseValueFromText(valueStr)
				measurements[key] = value
			}
		} else if strings.Contains(line, "=") {
			// Assignment format: "temp=25.3"
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				valueStr := strings.TrimSpace(parts[1])
				value := parseValueFromText(valueStr)
				measurements[key] = value
			}
		} else if strings.Contains(line, ",") {
			// CSV format: "temperature,25.3,°C"
			parts := strings.Split(line, ",")
			if len(parts) >= 2 {
				key := strings.TrimSpace(parts[0])
				valueStr := strings.TrimSpace(parts[1])
				unit := ""
				if len(parts) > 2 {
					unit = strings.TrimSpace(parts[2])
				}

				value := parseValueFromText(valueStr)
				if unit != "" {
					measurements[key] = map[string]interface{}{
						"value": value,
						"unit":  unit,
					}
				} else {
					measurements[key] = value
				}
			}
		}
	}

	// Convert to readings
	req := &FlexibleReadingRequest{
		AssetSensorID: assetSensorID,
		SensorTypeID:  sensorTypeID,
		Measurements:  measurements,
		DataSource:    "text",
	}

	return ConvertToFlexibleReadings(req)
}

// parseValueFromText attempts to parse a value from text, trying numeric and boolean before string
func parseValueFromText(valueStr string) interface{} {
	// Remove common units to extract numeric value
	cleanValue := valueStr

	// Remove common units
	units := []string{"°C", "°F", "μg/m³", "mg/m³", "ppm", "ppb", "%", "Hz", "kHz", "MHz", "V", "mV", "A", "mA", "W", "kW", "Pa", "hPa", "bar", "m", "cm", "mm", "km", "kg", "g", "mg", "L", "mL"}
	for _, unit := range units {
		cleanValue = strings.TrimSuffix(cleanValue, unit)
	}
	cleanValue = strings.TrimSpace(cleanValue)

	// Try to parse as number
	if f, err := strconv.ParseFloat(cleanValue, 64); err == nil {
		return f
	}

	// Try to parse as boolean
	if b, err := strconv.ParseBool(strings.ToLower(cleanValue)); err == nil {
		return b
	}

	// Return as string
	return valueStr
}

// ValidateAgainstMeasurementType validates a reading against a measurement type schema
func (r *IoTSensorReadingFlexible) ValidateAgainstMeasurementType(mt *SensorMeasurementType) (bool, []string) {
	var errors []string

	// Check if measurement type matches
	if r.MeasurementType != mt.Name {
		errors = append(errors, fmt.Sprintf("measurement type '%s' does not match expected type '%s'", r.MeasurementType, mt.Name))
	}

	// Validate required fields
	for _, field := range mt.Fields {
		if field.Required {
			// Check if the field exists in the reading
			switch field.DataType {
			case "numeric":
				if r.NumericValue == nil {
					errors = append(errors, fmt.Sprintf("required numeric field '%s' is missing", field.Name))
				} else if field.Min != nil && *r.NumericValue < *field.Min {
					errors = append(errors, fmt.Sprintf("numeric field '%s' value %.2f is below minimum %.2f", field.Name, *r.NumericValue, *field.Min))
				} else if field.Max != nil && *r.NumericValue > *field.Max {
					errors = append(errors, fmt.Sprintf("numeric field '%s' value %.2f is above maximum %.2f", field.Name, *r.NumericValue, *field.Max))
				}
				// Validate unit for numeric fields
				if r.MeasurementUnit != nil && field.Unit != "" && *r.MeasurementUnit != field.Unit {
					errors = append(errors, fmt.Sprintf("unit '%s' does not match expected unit '%s' for field '%s'", *r.MeasurementUnit, field.Unit, field.Name))
				}
			case "text":
				if r.TextValue == nil {
					errors = append(errors, fmt.Sprintf("required text field '%s' is missing", field.Name))
				}
			case "boolean":
				if r.BooleanValue == nil {
					errors = append(errors, fmt.Sprintf("required boolean field '%s' is missing", field.Name))
				}
			}
		}
	}

	return len(errors) == 0, errors
}
