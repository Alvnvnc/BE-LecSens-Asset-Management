package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MeasurementDataType defines the data type of a measurement field
type MeasurementDataType string

const (
	MeasurementDataTypeString  MeasurementDataType = "string"
	MeasurementDataTypeNumber  MeasurementDataType = "number"
	MeasurementDataTypeBoolean MeasurementDataType = "boolean"
	MeasurementDataTypeArray   MeasurementDataType = "array"
	MeasurementDataTypeObject  MeasurementDataType = "object"
)

// SensorMeasurementField defines a single field in a sensor measurement
type SensorMeasurementField struct {
	Name        string              `json:"name"`
	Label       string              `json:"label"`
	Description string              `json:"description,omitempty"`
	DataType    MeasurementDataType `json:"data_type"`
	Required    bool                `json:"required"`
	Unit        string              `json:"unit,omitempty"`
	Min         *float64            `json:"min,omitempty"`
	Max         *float64            `json:"max,omitempty"`
}

// SensorMeasurementType defines the structure of measurements for a sensor type
// This combines the previous SensorType and SensorMeasurementType entities
type SensorMeasurementType struct {
	ID               uuid.UUID                `json:"id"`
	TenantID         uuid.UUID                `json:"tenant_id"`
	Name             string                   `json:"name"` // E.g., "pH Sensor", "Microplastic Detector"
	Description      string                   `json:"description"`
	UnitOfMeasure    string                   `json:"unit_of_measure,omitempty"` // E.g., "pH", "ppm", "°C"
	MinAcceptedValue float64                  `json:"min_accepted_value,omitempty"`
	MaxAcceptedValue float64                  `json:"max_accepted_value,omitempty"`
	PropertiesSchema json.RawMessage          `json:"properties_schema,omitempty"` // JSON Schema for additional properties
	Fields           []SensorMeasurementField `json:"fields"`
	UIConfiguration  json.RawMessage          `json:"ui_configuration,omitempty"` // For UI display settings
	Version          int                      `json:"version"`
	IsActive         bool                     `json:"is_active"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        *time.Time               `json:"updated_at,omitempty"`
}

// NewSensorMeasurementType creates a new measurement type with default values
// This replaces both NewSensorType and the previous NewSensorMeasurementType
func NewSensorMeasurementType() *SensorMeasurementType {
	now := time.Now()
	return &SensorMeasurementType{
		ID:        uuid.New(),
		Version:   1,
		IsActive:  true,
		Fields:    []SensorMeasurementField{},
		CreatedAt: now,
	}
}

// SensorMeasurementTypeFromJSON parses JSON into a SensorMeasurementType
func SensorMeasurementTypeFromJSON(data []byte) (*SensorMeasurementType, error) {
	var sensorType SensorMeasurementType
	if err := json.Unmarshal(data, &sensorType); err != nil {
		return nil, err
	}
	return &sensorType, nil
}

// AddField adds a new measurement field to the type definition
func (m *SensorMeasurementType) AddField(field SensorMeasurementField) {
	m.Fields = append(m.Fields, field)
}

// ValidateMeasurement validates that a measurement conforms to this type's schema
func (m *SensorMeasurementType) ValidateMeasurement(data []byte) (bool, []string) {
	var measurement map[string]interface{}
	if err := json.Unmarshal(data, &measurement); err != nil {
		return false, []string{"Invalid JSON format"}
	}

	var errors []string

	// Check required fields
	for _, field := range m.Fields {
		if field.Required {
			if _, ok := measurement[field.Name]; !ok {
				errors = append(errors, "Required field missing: "+field.Name)
			}
		}
	}

	// Validate data types and ranges
	for fieldName, value := range measurement {
		// Find field definition
		var fieldDef *SensorMeasurementField
		for _, f := range m.Fields {
			if f.Name == fieldName {
				fieldDef = &f
				break
			}
		}

		// Skip validation for fields not defined in the schema
		if fieldDef == nil {
			continue
		}

		// Validate type
		switch fieldDef.DataType {
		case MeasurementDataTypeNumber:
			floatVal, ok := value.(float64)
			if !ok {
				errors = append(errors, fieldName+" must be a number")
				continue
			}

			// Check range if specified
			if fieldDef.Min != nil && floatVal < *fieldDef.Min {
				errors = append(errors, fieldName+" is below minimum value")
			}
			if fieldDef.Max != nil && floatVal > *fieldDef.Max {
				errors = append(errors, fieldName+" is above maximum value")
			}

		case MeasurementDataTypeString:
			if _, ok := value.(string); !ok {
				errors = append(errors, fieldName+" must be a string")
			}

		case MeasurementDataTypeBoolean:
			if _, ok := value.(bool); !ok {
				errors = append(errors, fieldName+" must be a boolean")
			}

		case MeasurementDataTypeArray:
			if _, ok := value.([]interface{}); !ok {
				errors = append(errors, fieldName+" must be an array")
			}

		case MeasurementDataTypeObject:
			if _, ok := value.(map[string]interface{}); !ok {
				errors = append(errors, fieldName+" must be an object")
			}
		}
	}

	return len(errors) == 0, errors
}

// ToJSON converts the measurement type definition to JSON
func (m *SensorMeasurementType) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
