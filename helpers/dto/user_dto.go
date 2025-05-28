package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserDTO represents a user response from the external user management API
type UserDTO struct {
	ID           uuid.UUID      `json:"id"`
	Email        string         `json:"email"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Role         string         `json:"role"`
	IsActive     bool           `json:"is_active"`
	TenantAccess []TenantAccess `json:"tenant_access"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// TenantAccess represents the user's access to a tenant
type TenantAccess struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	AccessLevel string    `json:"access_level"`
	IsDefault   bool      `json:"is_default"`
}

// TokenValidationRequest represents the request to validate a JWT token
type TokenValidationRequest struct {
	Token string `json:"token"`
}

// TokenValidationResponse represents the response from token validation
type TokenValidationResponse struct {
	Valid    bool      `json:"valid"`
	UserID   uuid.UUID `json:"userID"`
	UserRole string    `json:"userRole"`
	TenantID uuid.UUID `json:"tenantID"` // Add TenantID field
	Email    string    `json:"email"`
}

// PermissionValidationRequest represents the request to validate user permissions
type PermissionValidationRequest struct {
	Token        string    `json:"token"`
	TenantID     uuid.UUID `json:"tenantId,omitempty"`
	RequiredRole string    `json:"requiredRole,omitempty"`
	Permissions  []string  `json:"permissions,omitempty"`
}

// PermissionValidationResponse represents the response from permission validation
type PermissionValidationResponse struct {
	Valid             bool      `json:"valid"`
	HasRolePermission bool      `json:"hasRolePermission"`
	HasTenantAccess   bool      `json:"hasTenantAccess"`
	UserID            uuid.UUID `json:"userID"`
	UserRole          string    `json:"userRole"`
}

// UserTenantAccessRequest represents the request to validate user-tenant access
type UserTenantAccessRequest struct {
	UserID uuid.UUID `json:"userId"`
}

// UserTenantAccessResponse represents the response from user-tenant access validation
type UserTenantAccessResponse struct {
	UserID    uuid.UUID `json:"userID"`
	TenantID  uuid.UUID `json:"tenantID"`
	HasAccess bool      `json:"hasAccess"`
}

// ExternalUserResponse represents the external API response for a single user
type ExternalUserResponse struct {
	User UserDTO `json:"data"`
}

// SuperAdminValidationResponse represents the response from the validate-superadmin endpoint
type SuperAdminValidationResponse struct {
	Valid        bool   `json:"valid"`
	UserID       string `json:"userID"`
	UserRole     string `json:"userRole"`
	IsSuperAdmin bool   `json:"isSuperAdmin"`
}

// HasAccessToTenant checks if the user has access to the specified tenant
func (u *UserDTO) HasAccessToTenant(tenantID uuid.UUID) bool {
	for _, access := range u.TenantAccess {
		if access.TenantID == tenantID {
			return true
		}
	}
	return false
}

// GetDefaultTenant returns the user's default tenant
func (u *UserDTO) GetDefaultTenant() (uuid.UUID, bool) {
	for _, access := range u.TenantAccess {
		if access.IsDefault {
			return access.TenantID, true
		}
	}

	// If no default is set but user has at least one tenant, return the first one
	if len(u.TenantAccess) > 0 {
		return u.TenantAccess[0].TenantID, true
	}

	return uuid.Nil, false
}
