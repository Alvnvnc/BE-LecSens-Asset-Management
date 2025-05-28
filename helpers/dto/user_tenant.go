package dto

import (
	"time"
)

// UserTenantCurrentResponse represents the response from GET /api/v1/user-tenant/current
type UserTenantCurrentResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID           string    `json:"user_id"`
		TenantID         string    `json:"tenant_id"`
		TenantName       string    `json:"tenant_name"`
		UserRoleInTenant string    `json:"user_role_in_tenant"`
		IsActive         bool      `json:"is_active"`
		JoinedAt         time.Time `json:"joined_at"`
	} `json:"data"`
}

// UserTenantsListResponse represents the response from GET /api/v1/user-tenant/tenants
type UserTenantsListResponse struct {
	Message string `json:"message"`
	Data    struct {
		Tenants []struct {
			TenantID   string    `json:"tenant_id"`
			TenantName string    `json:"tenant_name"`
			UserRole   string    `json:"user_role"`
			IsCurrent  bool      `json:"is_current"`
			IsActive   bool      `json:"is_active"`
			JoinedAt   time.Time `json:"joined_at"`
		} `json:"tenants"`
		TotalCount int `json:"total_count"`
	} `json:"data"`
}

// SwitchTenantRequest represents the request body for POST /api/v1/user-tenant/switch
type SwitchTenantRequest struct {
	TenantID string `json:"tenant_id"`
}

// SwitchTenantResponse represents the response from POST /api/v1/user-tenant/switch
type SwitchTenantResponse struct {
	Message string `json:"message"`
	Data    struct {
		PreviousTenantID string    `json:"previous_tenant_id"`
		NewTenantID      string    `json:"new_tenant_id"`
		NewTenantName    string    `json:"new_tenant_name"`
		UserRole         string    `json:"user_role"`
		SwitchedAt       time.Time `json:"switched_at"`
	} `json:"data"`
}

// UserTenantAccessValidationResponse represents the response from POST /api/v1/user-tenant/validate-access
type UserTenantAccessValidationResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID      string   `json:"user_id"`
		TenantID    string   `json:"tenant_id"`
		HasAccess   bool     `json:"has_access"`
		UserRole    string   `json:"user_role"`
		IsActive    bool     `json:"is_active"`
		Permissions []string `json:"permissions"`
	} `json:"data"`
}

// TenantUsersResponse represents the response from GET /api/v1/user-tenant/users
type TenantUsersResponse struct {
	Message string `json:"message"`
	Data    struct {
		Users []struct {
			UserID       string    `json:"user_id"`
			Username     string    `json:"username"`
			Email        string    `json:"email"`
			RoleInTenant string    `json:"role_in_tenant"`
			IsActive     bool      `json:"is_active"`
			JoinedAt     time.Time `json:"joined_at"`
			LastActive   time.Time `json:"last_active"`
		} `json:"users"`
		Pagination struct {
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
			TotalUsers  int `json:"total_users"`
			PerPage     int `json:"per_page"`
		} `json:"pagination"`
		TenantInfo struct {
			TenantID   string `json:"tenant_id"`
			TenantName string `json:"tenant_name"`
		} `json:"tenant_info"`
	} `json:"data"`
}
