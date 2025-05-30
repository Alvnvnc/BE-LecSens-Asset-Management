package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SensorType struct {
	ID           uuid.UUID
	Name         string
	Description  string
	Manufacturer string
	Model        string
	Version      string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type SensorTypeRepository struct {
	db *sql.DB
}

func NewSensorTypeRepository(db *sql.DB) *SensorTypeRepository {
	return &SensorTypeRepository{db: db}
}

// Create creates a new sensor type
func (r *SensorTypeRepository) Create(st *SensorType) error {
	query := `
		INSERT INTO sensor_types (
			id, name, description, manufacturer, model,
			version, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query,
		st.ID, st.Name, st.Description, st.Manufacturer, st.Model,
		st.Version, st.IsActive, st.CreatedAt, st.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating sensor type: %v", err)
	}
	return nil
}

// GetByID retrieves a sensor type by ID
func (r *SensorTypeRepository) GetByID(id uuid.UUID) (*SensorType, error) {
	query := `
		SELECT id, name, description, manufacturer, model,
			version, is_active, created_at, updated_at
		FROM sensor_types
		WHERE id = $1
	`
	st := &SensorType{}
	err := r.db.QueryRow(query, id).Scan(
		&st.ID, &st.Name, &st.Description, &st.Manufacturer, &st.Model,
		&st.Version, &st.IsActive, &st.CreatedAt, &st.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting sensor type: %v", err)
	}
	return st, nil
}

// GetAll retrieves all sensor types
func (r *SensorTypeRepository) GetAll() ([]*SensorType, error) {
	query := `
		SELECT id, name, description, manufacturer, model,
			version, is_active, created_at, updated_at
		FROM sensor_types
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting sensor types: %v", err)
	}
	defer rows.Close()

	var sensorTypes []*SensorType
	for rows.Next() {
		st := &SensorType{}
		err := rows.Scan(
			&st.ID, &st.Name, &st.Description, &st.Manufacturer, &st.Model,
			&st.Version, &st.IsActive, &st.CreatedAt, &st.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning sensor type: %v", err)
		}
		sensorTypes = append(sensorTypes, st)
	}
	return sensorTypes, nil
}

// Update updates a sensor type
func (r *SensorTypeRepository) Update(st *SensorType) error {
	query := `
		UPDATE sensor_types
		SET name = $1, description = $2, manufacturer = $3, model = $4,
			version = $5, is_active = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.Exec(query,
		st.Name, st.Description, st.Manufacturer, st.Model,
		st.Version, st.IsActive, time.Now(), st.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating sensor type: %v", err)
	}
	return nil
}

// Delete deletes a sensor type
func (r *SensorTypeRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM sensor_types WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting sensor type: %v", err)
	}
	return nil
}

// GetActive retrieves all active sensor types
func (r *SensorTypeRepository) GetActive() ([]*SensorType, error) {
	query := `
		SELECT id, name, description, manufacturer, model,
			version, is_active, created_at, updated_at
		FROM sensor_types
		WHERE is_active = true
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting active sensor types: %v", err)
	}
	defer rows.Close()

	var sensorTypes []*SensorType
	for rows.Next() {
		st := &SensorType{}
		err := rows.Scan(
			&st.ID, &st.Name, &st.Description, &st.Manufacturer, &st.Model,
			&st.Version, &st.IsActive, &st.CreatedAt, &st.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning sensor type: %v", err)
		}
		sensorTypes = append(sensorTypes, st)
	}
	return sensorTypes, nil
}
