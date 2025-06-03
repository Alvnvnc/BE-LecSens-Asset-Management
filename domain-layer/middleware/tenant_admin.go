package middleware

import (
	"be-lecsens/asset_management/helpers/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TenantAdminMiddleware ensures the user has tenant admin role
func TenantAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{
				Error:   "Unauthorized",
				Message: "User role not found",
			})
			c.Abort()
			return
		}

		// Check if user has tenant admin role
		if role != "tenant_admin" {
			c.JSON(http.StatusForbidden, common.ErrorResponse{
				Error:   "Forbidden",
				Message: "User does not have tenant admin privileges",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
