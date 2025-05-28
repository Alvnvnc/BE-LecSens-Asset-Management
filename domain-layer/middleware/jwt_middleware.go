package middleware

import (
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/common"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTConfig holds JWT middleware configuration
type JWTConfig struct {
	SecretKey string
}

// JWTMiddleware validates JWT tokens and extracts claims
func JWTMiddleware(config JWTConfig, userManagementService *service.UserManagementService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("JWT Middleware: Authorization header is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("JWT Middleware: Invalid Authorization header format: %s", authHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must be in format: Bearer {token}",
			})
			return
		}

		tokenString := parts[1]
		log.Printf("JWT Middleware: Validating token: %s...", tokenString[:10])

		// Validate token with the user management service
		validationResponse, err := userManagementService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			log.Printf("JWT Middleware: Token validation error: %v", err)

			// Provide more specific error messages
			errorMsg := "Invalid or expired token"
			if strings.Contains(err.Error(), "user not found") {
				errorMsg = "user not found"
			} else if strings.Contains(err.Error(), "unauthorized") {
				errorMsg = "unauthorized access"
			} else if strings.Contains(err.Error(), "external API server error") {
				errorMsg = "authentication service temporarily unavailable"
			} else if strings.Contains(err.Error(), "API error") {
				errorMsg = "API authentication failed"
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   errorMsg,
				"details": err.Error(),
				"source":  "jwt",
			})
			return
		}

		if !validationResponse.Valid {
			log.Printf("JWT Middleware: Token is invalid according to validation response")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		log.Printf("JWT Middleware: Token validated successfully for user %s with role %s",
			validationResponse.UserID, validationResponse.UserRole)

		// Create JWT claims from validation response (tenant ID now included)
		claims := &common.JWTClaims{
			UserID:   validationResponse.UserID,
			Role:     validationResponse.UserRole,
			TenantID: validationResponse.TenantID, // Use tenant ID from validation response
		}

		log.Printf("JWT Middleware: User %s belongs to tenant %s", claims.UserID, claims.TenantID)

		// Store claims in context for later use
		c.Set("jwt_claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("tenant_id", claims.TenantID) // Set tenant_id directly for repository use

		// Add tenant ID to request context for common.RequireTenant() usage
		ctx := common.WithTenant(c.Request.Context(), claims.TenantID)
		c.Request = c.Request.WithContext(ctx)

		log.Printf("JWT Middleware: Successfully set tenant %s in context", claims.TenantID)
		c.Next()
	}
}
