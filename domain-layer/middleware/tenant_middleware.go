package middleware

import (
	"be-lecsens/asset_management/helpers/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TenantMiddleware validates tenant context and user access
// This middleware assumes JWT has already been validated by JWTMiddleware
// It ensures the user belongs to a tenant and has proper tenant context
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Tenant Middleware: Validating tenant context for %s %s", c.Request.Method, c.Request.URL.Path)

		// Get user information from context (set by JWT middleware)
		userID, userExists := c.Get("user_id")
		userRole, roleExists := c.Get("user_role")
		tenantID, tenantExists := c.Get("tenant_id")

		if !userExists || !roleExists {
			log.Printf("Tenant Middleware: User information not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User information not found in context",
				"source":  "tenant_middleware",
			})
			c.Abort()
			return
		}

		// Parse user role as string
		userRoleStr, ok := userRole.(string)
		if !ok {
			log.Printf("Tenant Middleware: Invalid user role format: %v", userRole)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid user role format",
				"source":  "tenant_middleware",
			})
			c.Abort()
			return
		}

		// Handle tenant context
		var tenantUUID uuid.UUID
		var err error

		if tenantExists {
			// Tenant ID already set by JWT middleware
			switch t := tenantID.(type) {
			case string:
				tenantUUID, err = uuid.Parse(t)
				if err != nil {
					log.Printf("Tenant Middleware: Invalid tenant ID format: %v", err)
					c.JSON(http.StatusUnauthorized, gin.H{
						"error":   "Unauthorized",
						"message": "Invalid tenant ID format",
						"source":  "tenant_middleware",
					})
					c.Abort()
					return
				}
			case uuid.UUID:
				tenantUUID = t
			default:
				log.Printf("Tenant Middleware: Unexpected tenant ID type: %T", tenantID)
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "Unauthorized",
					"message": "Unexpected tenant ID type",
					"source":  "tenant_middleware",
				})
				c.Abort()
				return
			}
		} else {
			// For SUPERADMIN, tenant might not be required (global access)
			if userRoleStr == "SUPERADMIN" {
				log.Printf("Tenant Middleware: SUPERADMIN detected - allowing global access")
				// Use a default tenant or let it be empty for global operations
				tenantUUID = uuid.Nil // or use a default global tenant ID
			} else {
				log.Printf("Tenant Middleware: No tenant ID found for non-SUPERADMIN user")
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "Unauthorized",
					"message": "Tenant context required for this operation",
					"source":  "tenant_middleware",
				})
				c.Abort()
				return
			}
		}

		// Parse user ID
		userUUID, err := uuid.Parse(userID.(string))
		if err != nil {
			log.Printf("Tenant Middleware: Invalid user ID format: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid user ID format",
				"source":  "tenant_middleware",
			})
			c.Abort()
			return
		}

		// Ensure tenant context is properly set for repository usage
		if tenantUUID != uuid.Nil {
			ctx := common.WithTenant(c.Request.Context(), tenantUUID)
			c.Request = c.Request.WithContext(ctx)
			c.Set("tenant_id", tenantUUID.String()) // Ensure string format in context
		}

		log.Printf("Tenant Middleware: Validated - UserID: %s, Role: %s, TenantID: %s",
			userUUID, userRoleStr, tenantUUID)

		c.Next()
	}
}
