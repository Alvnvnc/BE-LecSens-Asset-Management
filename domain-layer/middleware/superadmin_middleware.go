package middleware

import (
	"be-lecsens/asset_management/domain-layer/service"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireSuperAdminMiddleware ensures only SuperAdmin users can access certain endpoints
// SuperAdmin is a global role that can access everything without tenant restrictions
func RequireSuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("SuperAdmin Middleware: Checking user role for %s %s", c.Request.Method, c.Request.URL.Path)

		// Get user role from JWT context (set by JWT middleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			log.Printf("SuperAdmin Middleware: User role not found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "User role not found in context",
				"source":  "superadmin_middleware",
			})
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok {
			log.Printf("SuperAdmin Middleware: Invalid user role format: %v", userRole)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Invalid user role format in context",
				"source":  "superadmin_middleware",
			})
			return
		}

		log.Printf("SuperAdmin Middleware: User role is: %s", userRoleStr)

		// Check if user has SuperAdmin privileges
		// SUPERADMIN: Global role, can access everything without tenant restrictions
		if userRoleStr != "SUPERADMIN" {
			log.Printf("SuperAdmin Middleware: Access denied for role: %s", userRoleStr)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"details": "SuperAdmin role required for this operation",
				"source":  "superadmin_middleware",
			})
			return
		}

		log.Printf("SuperAdmin Middleware: Access granted for role: %s", userRoleStr)

		// For SUPERADMIN, we might want to bypass tenant restrictions in the future
		if userRoleStr == "SUPERADMIN" {
			// SuperAdmin can access resources across all tenants
			// For now, we'll just log this capability
			log.Printf("SuperAdmin Middleware: SuperAdmin detected - global access granted")
		}

		c.Next()
	}
}

// RequireAdminOrSuperAdminMiddleware allows both Admin and SuperAdmin roles
func RequireAdminOrSuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Admin/SuperAdmin Middleware: Checking user role for %s %s", c.Request.Method, c.Request.URL.Path)

		userRole, exists := c.Get("user_role")
		if !exists {
			log.Printf("Admin/SuperAdmin Middleware: User role not found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "User role not found in context",
				"source":  "admin_middleware",
			})
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok {
			log.Printf("Admin/SuperAdmin Middleware: Invalid user role format: %v", userRole)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Invalid user role format in context",
				"source":  "admin_middleware",
			})
			return
		}

		log.Printf("Admin/SuperAdmin Middleware: User role is: %s", userRoleStr)

		// Allow SUPERADMIN, ADMIN, and potentially other administrative roles
		allowedRoles := []string{"SUPERADMIN", "ADMIN", "MANAGER"}
		isAllowed := false
		for _, role := range allowedRoles {
			if userRoleStr == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			log.Printf("Admin/SuperAdmin Middleware: Access denied for role: %s", userRoleStr)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"details": "Admin, Manager, or SuperAdmin role required for this operation",
				"source":  "admin_middleware",
			})
			return
		}

		log.Printf("Admin/SuperAdmin Middleware: Access granted for role: %s", userRoleStr)
		c.Next()
	}
}

// RequireSuperAdminWithAPIValidationMiddleware ensures only SuperAdmin users can access certain endpoints
// This middleware makes direct calls to the external API to validate SuperAdmin status
func RequireSuperAdminWithAPIValidationMiddleware(userManagementService *service.UserManagementService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("SuperAdmin API Validation Middleware: Checking user role for %s %s", c.Request.Method, c.Request.URL.Path)

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("SuperAdmin API Validation Middleware: Authorization header is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Authorization header is required",
				"source":  "superadmin_middleware",
			})
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("SuperAdmin API Validation Middleware: Invalid Authorization header format: %s", authHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Authorization header must be in format: Bearer {token}",
				"source":  "superadmin_middleware",
			})
			return
		}

		tokenString := parts[1]
		log.Printf("SuperAdmin API Validation Middleware: Validating token: %s...", tokenString[:10])

		// Validate SuperAdmin status using external API
		superAdminResponse, err := userManagementService.ValidateSuperAdmin(c.Request.Context(), tokenString)
		if err != nil {
			log.Printf("SuperAdmin API Validation Middleware: Error validating SuperAdmin: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Failed to validate SuperAdmin status: " + err.Error(),
				"source":  "superadmin_middleware",
			})
			return
		}

		// Log detailed response information for debugging
		log.Printf("SuperAdmin API Validation Middleware: API Response - UserID: %s, Role: %s, IsSuperAdmin: %v, Valid: %v",
			superAdminResponse.UserID, superAdminResponse.UserRole, superAdminResponse.IsSuperAdmin, superAdminResponse.Valid)

		// Check if the user is a valid SuperAdmin
		if !superAdminResponse.Valid || !superAdminResponse.IsSuperAdmin {
			log.Printf("SuperAdmin API Validation Middleware: Access denied - not a SuperAdmin")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"details": "SuperAdmin role required for this operation",
				"source":  "superadmin_middleware",
				"role":    superAdminResponse.UserRole,
			})
			return
		}

		// Store information in context for later use
		c.Set("user_id", superAdminResponse.UserID)
		c.Set("user_role", superAdminResponse.UserRole)
		c.Set("is_superadmin", superAdminResponse.IsSuperAdmin)

		log.Printf("SuperAdmin API Validation Middleware: Access granted for SuperAdmin: %s", superAdminResponse.UserID)
		c.Next()
	}
}

// SuperAdminPassthroughMiddleware is a lightweight middleware that only checks the role in the JWT token
// without validating the signature. This is useful for development or when you need to bypass JWT validation
// but still want to check the role.
func SuperAdminPassthroughMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("SuperAdmin Passthrough Middleware: Checking user role for %s %s", c.Request.Method, c.Request.URL.Path)

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("SuperAdmin Passthrough Middleware: Authorization header is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Authorization header is required",
				"source":  "superadmin_passthrough_middleware",
			})
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("SuperAdmin Passthrough Middleware: Invalid Authorization header format: %s", authHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Authorization header must be in format: Bearer {token}",
				"source":  "superadmin_passthrough_middleware",
			})
			return
		}

		tokenString := parts[1]
		log.Printf("SuperAdmin Passthrough Middleware: Processing token: %s...", tokenString[:10])

		// Split the token into parts
		tokenParts := strings.Split(tokenString, ".")
		if len(tokenParts) != 3 {
			log.Printf("SuperAdmin Passthrough Middleware: Invalid token format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Invalid token format",
				"source":  "superadmin_passthrough_middleware",
			})
			return
		}

		// Decode the payload (second part)
		payload, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
		if err != nil {
			log.Printf("SuperAdmin Passthrough Middleware: Error decoding payload: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Error decoding token payload",
				"source":  "superadmin_passthrough_middleware",
			})
			return
		}

		// Parse the payload
		var claims map[string]interface{}
		if err := json.Unmarshal(payload, &claims); err != nil {
			log.Printf("SuperAdmin Passthrough Middleware: Error parsing payload: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Error parsing token payload",
				"source":  "superadmin_passthrough_middleware",
			})
			return
		}

		// Debug: Print all claims from the token
		log.Printf("SuperAdmin Passthrough Middleware: Token claims: %+v", claims)
		for key, value := range claims {
			log.Printf("SuperAdmin Passthrough Middleware: Claim[%s] = %v (type: %T)", key, value, value)
		}

		// Extract role from claims using role_name
		role, ok := claims["role_name"].(string)
		if !ok {
			log.Printf("SuperAdmin Passthrough Middleware: Role not found in token (role_name)")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Role not found in token (role_name)",
				"source":  "superadmin_passthrough_middleware",
			})
			return
		}

		log.Printf("SuperAdmin Passthrough Middleware: User role is: %s", role)

		// Check if user has SuperAdmin privileges
		if role != "SUPERADMIN" {
			log.Printf("SuperAdmin Passthrough Middleware: Access denied for role: %s", role)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"details": "SuperAdmin role required for this operation",
				"source":  "superadmin_passthrough_middleware",
			})
			return
		}

		// Store user information in context
		if userID, ok := claims["sub"].(string); ok {
			c.Set("user_id", userID)
		}
		c.Set("user_role", role)
		c.Set("is_superadmin", true)

		log.Printf("SuperAdmin Passthrough Middleware: Access granted for role: %s", role)
		c.Next()
	}
}
