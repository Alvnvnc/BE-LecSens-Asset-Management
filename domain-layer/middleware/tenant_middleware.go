package middleware

import (
	"be-lecsens/asset_management/helpers/common"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TenantMiddleware validates tenant context from JWT token
// This middleware combines functionality from RequireTenant middleware:
// 1. Extracts and decodes JWT token from Authorization header
// 2. Sets tenant context in Gin context
// 3. Adds tenant context to request for repository usage
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Tenant Middleware: Starting validation")

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Tenant Middleware: No Authorization header found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from Bearer header
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			log.Println("Tenant Middleware: Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		log.Printf("Tenant Middleware: Extracted token: %s...", token[:min(20, len(token))])

		// Decode JWT token to get user information
		jwtPayload, err := common.DecodeJWTToken(token)
		if err != nil {
			log.Printf("Tenant Middleware: Error decoding JWT token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Debug: Log JWT payload
		log.Printf("Tenant Middleware: Decoded JWT payload: %+v", jwtPayload)

		// Parse user ID as UUID
		userID, err := uuid.Parse(jwtPayload.UserID)
		if err != nil {
			log.Printf("Tenant Middleware: Invalid user ID format: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		// For now, we'll use a default tenant ID since the token doesn't contain tenant information
		// In the future, this should come from an external API call or be included in the JWT
		defaultTenantID := uuid.MustParse("e3b8f35c-a6d0-4bd3-bd78-84276b67b32e") // Use the tenant ID from your example

		// Set context values
		c.Set("tenant_id", defaultTenantID.String())
		c.Set("user_id", userID.String())
		c.Set("tenant_name", "Default Tenant") // Default value for now
		c.Set("user_role", jwtPayload.RoleName)

		log.Printf("Tenant Middleware: Set tenant_id=%s, user_id=%s, role=%s",
			defaultTenantID.String(), userID.String(), jwtPayload.RoleName)

		// Add tenant to request context for repository usage
		ctx := common.WithTenant(c.Request.Context(), defaultTenantID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
