package dto

import (
	"github.com/google/uuid"
)

// TenantDTO represents a tenant response from the external API
type TenantDTO struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	LogoURL          string    `json:"logo_url,omitempty"`
	ContactEmail     string    `json:"contact_email"`
	ContactPhone     string    `json:"contact_phone,omitempty"`
	SubscriptionPlan string    `json:"subscription_plan"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        string    `json:"created_at"`
}

// IsSubscriptionValid checks if the tenant's subscription is still valid
func (t *TenantDTO) IsSubscriptionValid() bool {
	// For now, just check if tenant is active
	return t.IsActive
}

// Meta represents pagination metadata in the API response
type Meta struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// ExternalTenantsResponse represents the external API response for tenants
type ExternalTenantsResponse struct {
	Meta    Meta        `json:"meta"`
	Tenants []TenantDTO `json:"tenants"`
}

// ExternalTenantResponse represents the external API response for a single tenant
type ExternalTenantResponse struct {
	Tenant TenantDTO `json:"data"`
}

// TenantLimits represents the limits for a tenant
type TenantLimits struct {
	MaxUsers         int    `json:"maxUsers"`
	MaxAssets        int    `json:"maxAssets"`
	MaxRentals       int    `json:"maxRentals"`
	SubscriptionPlan string `json:"subscriptionPlan"`
}

// TenantLimitsResponse represents the response from tenant limits endpoint
type TenantLimitsResponse struct {
	TenantID uuid.UUID    `json:"tenantID"`
	Limits   TenantLimits `json:"limits"`
}

// TenantSubscriptionResponse represents the response from tenant subscription endpoint
type TenantSubscriptionResponse struct {
	TenantID              uuid.UUID `json:"tenantID"`
	SubscriptionPlan      string    `json:"subscriptionPlan"`
	SubscriptionStartDate string    `json:"subscriptionStartDate"`
	SubscriptionEndDate   string    `json:"subscriptionEndDate"`
	IsActive              bool      `json:"isActive"`
}
