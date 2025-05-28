package common

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Key for storing tenant information in context
type contextKey string

const (
	// TenantContextKey is the key used to store tenant ID in the context
	TenantContextKey contextKey = "tenant_id"
)

// GetTenantID retrieves the tenant ID from the context
func GetTenantID(ctx context.Context) (uuid.UUID, bool) {
	// First try to get from context.Context (stored as uuid.UUID)
	if tenantID, ok := ctx.Value(TenantContextKey).(uuid.UUID); ok {
		return tenantID, true
	}

	// If not found, try to get from Gin context (stored as string)
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if tenantIDStr, exists := ginCtx.Get("tenant_id"); exists {
			if tenantIDString, ok := tenantIDStr.(string); ok {
				if tenantID, err := uuid.Parse(tenantIDString); err == nil {
					return tenantID, true
				}
			}
		}
	}

	return uuid.UUID{}, false
}

// WithTenant adds tenant ID to context
func WithTenant(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, TenantContextKey, tenantID)
}

// RequireTenant gets the tenant ID from context and panics if not found
// Use this in internal functions where tenant ID should always be present
func RequireTenant(ctx context.Context) uuid.UUID {
	tenantID, ok := GetTenantID(ctx)
	if !ok {
		panic("tenant ID is required but not found in context")
	}
	return tenantID
}

// IsSuperAdmin checks if the current user has SuperAdmin role
func IsSuperAdmin(ctx context.Context) bool {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		userRole, exists := ginCtx.Get("user_role")
		if !exists {
			return false
		}
		if roleStr, ok := userRole.(string); ok {
			return roleStr == "SUPERADMIN"
		}
	}
	return false
}

// GetUserRole retrieves the user role from context
func GetUserRole(ctx context.Context) (string, bool) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		userRole, exists := ginCtx.Get("user_role")
		if !exists {
			return "", false
		}
		if roleStr, ok := userRole.(string); ok {
			return roleStr, true
		}
	}
	return "", false
}
