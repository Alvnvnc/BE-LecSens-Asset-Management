package middleware

import (
	"be-lecsens/asset_management/helpers/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequireUserWithTenantMiddleware ensures the user has valid tenant access
// This middleware allows any authenticated user who belongs to a tenant
// It validates that the user has tenant context and access
func RequireUserWithTenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("User With Tenant Middleware: Checking user access for %s %s", c.Request.Method, c.Request.URL.Path)

		// Get user information from context (set by JWT middleware)
		userID, userExists := c.Get("user_id")
		userRole, roleExists := c.Get("user_role")
		tenantID, tenantExists := c.Get("tenant_id")

		if !userExists || !roleExists {
			log.Printf("User With Tenant Middleware: User information not found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "User information not found in context",
				"source":  "user_tenant_middleware",
			})
			return
		}

		// Parse user role as string
		userRoleStr, ok := userRole.(string)
		if !ok {
			log.Printf("User With Tenant Middleware: Invalid user role format: %v", userRole)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Invalid user role format",
				"source":  "user_tenant_middleware",
			})
			return
		}

		// Parse user ID
		userIDStr, ok := userID.(string)
		if !ok {
			log.Printf("User With Tenant Middleware: Invalid user ID format: %v", userID)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Invalid user ID format",
				"source":  "user_tenant_middleware",
			})
			return
		}

		// Validate user ID format
		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			log.Printf("User With Tenant Middleware: Invalid user UUID format: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Invalid user UUID format",
				"source":  "user_tenant_middleware",
			})
			return
		}

		// Handle tenant validation based on role
		if userRoleStr == "SUPERADMIN" {
			// SUPERADMIN can access without specific tenant restrictions
			log.Printf("User With Tenant Middleware: SUPERADMIN access granted - UserID: %s", userUUID)
			c.Next()
			return
		}

		// For non-SUPERADMIN users, tenant is required
		if !tenantExists {
			log.Printf("User With Tenant Middleware: No tenant ID found for role: %s", userRoleStr)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"details": "Tenant context required for this operation",
				"source":  "user_tenant_middleware",
			})
			return
		}

		// Validate tenant ID format
		var tenantUUID uuid.UUID
		switch t := tenantID.(type) {
		case string:
			tenantUUID, err = uuid.Parse(t)
			if err != nil {
				log.Printf("User With Tenant Middleware: Invalid tenant ID format: %v", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":   "Unauthorized",
					"details": "Invalid tenant ID format",
					"source":  "user_tenant_middleware",
				})
				return
			}
		case uuid.UUID:
			tenantUUID = t
		default:
			log.Printf("User With Tenant Middleware: Unexpected tenant ID type: %T", tenantID)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Unexpected tenant ID type",
				"source":  "user_tenant_middleware",
			})
			return
		}

		// Validate that tenant is not empty
		if tenantUUID == uuid.Nil {
			log.Printf("User With Tenant Middleware: Empty tenant ID for user: %s", userUUID)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"details": "Valid tenant context required",
				"source":  "user_tenant_middleware",
			})
			return
		}

		// Ensure tenant context is properly set for repository usage
		ctx := common.WithTenant(c.Request.Context(), tenantUUID)
		c.Request = c.Request.WithContext(ctx)
		c.Set("tenant_id", tenantUUID.String()) // Ensure string format in context

		log.Printf("User With Tenant Middleware: Access granted - UserID: %s, Role: %s, TenantID: %s",
			userUUID, userRoleStr, tenantUUID)

		c.Next()
	}
}

// RequireUserRoleMiddleware ensures the user has one of the specified roles
// This is a flexible middleware that can check for multiple allowed roles
func RequireUserRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("User Role Middleware: Checking user role for %s %s (allowed: %v)",
			c.Request.Method, c.Request.URL.Path, allowedRoles)

		// Get user role from context (set by JWT middleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			log.Printf("User Role Middleware: User role not found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "User role not found in context",
				"source":  "user_role_middleware",
			})
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok {
			log.Printf("User Role Middleware: Invalid user role format: %v", userRole)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Invalid user role format in context",
				"source":  "user_role_middleware",
			})
			return
		}

		log.Printf("User Role Middleware: User role is: %s", userRoleStr)

		// Check if user has one of the allowed roles
		isAllowed := false
		for _, role := range allowedRoles {
			if userRoleStr == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			log.Printf("User Role Middleware: Access denied for role: %s (allowed: %v)", userRoleStr, allowedRoles)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":         "Forbidden",
				"details":       "Insufficient role privileges for this operation",
				"source":        "user_role_middleware",
				"allowed_roles": allowedRoles,
				"user_role":     userRoleStr,
			})
			return
		}

		log.Printf("User Role Middleware: Access granted for role: %s", userRoleStr)
		c.Next()
	}
}
