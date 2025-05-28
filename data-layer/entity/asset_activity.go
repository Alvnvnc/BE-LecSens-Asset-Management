package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ActivityType defines the type of activity performed on an asset
type ActivityType string

const (
	ActivityTypeMaintenance ActivityType = "maintenance"
	ActivityTypeCalibration ActivityType = "calibration"
	ActivityTypeInspection  ActivityType = "inspection"
)

// ActivityStatus defines the current status of an activity
type ActivityStatus string

const (
	ActivityStatusPending   ActivityStatus = "pending"
	ActivityStatusCompleted ActivityStatus = "completed"
	ActivityStatusFailed    ActivityStatus = "failed"
)

// AssetActivity represents a record of work performed on an asset
type AssetActivity struct {
	ID            uuid.UUID      `json:"id"`
	TenantID      uuid.UUID      `json:"tenant_id"`
	AssetID       uuid.UUID      `json:"asset_id"`
	ActivityType  ActivityType   `json:"activity_type"`
	Status        ActivityStatus `json:"status"`
	ScheduledDate time.Time      `json:"scheduled_date"`
	CompletedDate *time.Time     `json:"completed_date,omitempty"`
	Description   string         `json:"description,omitempty"`
	Notes         string         `json:"notes,omitempty"`
	AssignedTo    *uuid.UUID     `json:"assigned_to,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     *time.Time     `json:"updated_at,omitempty"`
}

// NewAssetActivity creates a new asset activity with default values
func NewAssetActivity() *AssetActivity {
	now := time.Now()
	return &AssetActivity{
		ID:        uuid.New(),
		Status:    ActivityStatusPending,
		CreatedAt: now,
	}
}

// Validate validates the asset activity fields
func (a *AssetActivity) Validate() error {
	if a.TenantID == uuid.Nil {
		return errors.New("tenant_id is required")
	}
	if a.AssetID == uuid.Nil {
		return errors.New("asset_id is required")
	}
	if a.ActivityType == "" {
		return errors.New("activity_type is required")
	}
	if a.CompletedDate != nil && a.CompletedDate.Before(a.ScheduledDate) {
		return errors.New("completed_date cannot be before scheduled_date")
	}
	return nil
}

// MarkCompleted marks the activity as completed
func (a *AssetActivity) MarkCompleted() {
	now := time.Now()
	a.Status = ActivityStatusCompleted
	a.CompletedDate = &now
	a.UpdatedAt = &now
}

// MarkFailed marks the activity as failed
func (a *AssetActivity) MarkFailed() {
	now := time.Now()
	a.Status = ActivityStatusFailed
	a.UpdatedAt = &now
}

// IsOverdue checks if the activity is overdue
func (a *AssetActivity) IsOverdue() bool {
	return a.Status == ActivityStatusPending && time.Now().After(a.ScheduledDate)
}

// GetDuration returns the duration of the activity if completed
func (a *AssetActivity) GetDuration() *time.Duration {
	if a.CompletedDate == nil {
		return nil
	}
	duration := a.CompletedDate.Sub(a.ScheduledDate)
	return &duration
}
