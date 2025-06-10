package middleware

import (
	"be-lecsens/asset_management/helpers/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TenantAdminMiddleware ensures the user has tenant admin role
// Supports multiple admin role formats: "ADMIN", "tenant_admin", "MANAGER"
func TenantAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Tenant Admin Middleware: Checking user role for %s %s", c.Request.Method, c.Request.URL.Path)

		// Get user role from context (set by JWT middleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			log.Printf("Tenant Admin Middleware: User role not found in context")
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{
				Error:   "Unauthorized",
				Message: "User role not found in context",
			})
			c.Abort()
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok {
			log.Printf("Tenant Admin Middleware: Invalid user role format: %v", userRole)
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid user role format in context",
			})
			c.Abort()
			return
		}

		log.Printf("Tenant Admin Middleware: User role is: %s", userRoleStr)

		// Check if user has admin privileges (support multiple formats)
		// SUPERADMIN: Global admin across all tenants
		// ADMIN: Tenant-level admin
		// MANAGER: Tenant-level manager (administrative rights)
		// tenant_admin: Legacy format for tenant admin
		allowedRoles := []string{"SUPERADMIN", "ADMIN", "MANAGER", "tenant_admin"}
		isAllowed := false
		for _, role := range allowedRoles {
			if userRoleStr == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			log.Printf("Tenant Admin Middleware: Access denied for role: %s", userRoleStr)
			c.JSON(http.StatusForbidden, common.ErrorResponse{
				Error:   "Forbidden",
				Message: "Admin or Manager role required for this operation",
			})
			c.Abort()
			return
		}

		log.Printf("Tenant Admin Middleware: Access granted for role: %s", userRoleStr)
		c.Next()
	}
}
