package service

import (
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// UserManagementService handles all user, tenant, and user-tenant relationship operations
type UserManagementService struct {
	userAPIURL                  string
	userTenantAPIURL            string
	tenantAPIURL                string
	apiKey                      string
	validateTokenEndpoint       string
	userInfoEndpoint            string
	validatePermissionsEndpoint string
	validateSuperAdminEndpoint  string
	userCache                   map[uuid.UUID]*dto.UserDTO
	tenantCache                 map[uuid.UUID]*dto.TenantDTO
	cacheTTL                    time.Duration
	cacheLastUpdated            map[uuid.UUID]time.Time
	rateLimiter                 *common.RateLimiter
	httpClient                  *http.Client
}

// NewUserManagementService creates a new consolidated user management service
func NewUserManagementService(
	userAPIURL,
	userTenantAPIURL,
	tenantAPIURL,
	apiKey string,
	validateTokenEndpoint,
	userInfoEndpoint,
	validatePermissionsEndpoint,
	validateSuperAdminEndpoint string,
) *UserManagementService {
	return &UserManagementService{
		userAPIURL:                  userAPIURL,
		userTenantAPIURL:            userTenantAPIURL,
		tenantAPIURL:                tenantAPIURL,
		apiKey:                      apiKey,
		validateTokenEndpoint:       validateTokenEndpoint,
		userInfoEndpoint:            userInfoEndpoint,
		validatePermissionsEndpoint: validatePermissionsEndpoint,
		validateSuperAdminEndpoint:  validateSuperAdminEndpoint,
		userCache:                   make(map[uuid.UUID]*dto.UserDTO),
		tenantCache:                 make(map[uuid.UUID]*dto.TenantDTO),
		cacheLastUpdated:            make(map[uuid.UUID]time.Time),
		cacheTTL:                    5 * time.Minute,
		rateLimiter:                 common.NewRateLimiter(10, 20),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxConnsPerHost:     100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// ValidateToken validates a JWT token against the user management service
func (s *UserManagementService) ValidateToken(ctx context.Context, tokenString string) (*dto.TokenValidationResponse, error) {
	if !s.rateLimiter.Allow("validate_token") {
		return nil, common.ErrRateLimitExceeded
	}

	url := fmt.Sprintf("%s%s", s.userAPIURL, s.validateTokenEndpoint)
	log.Printf("UserManagementService: Making request to %s", url)

	// Create request with Authorization header instead of request body
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("X-API-Key", s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("UserManagementService: Request failed: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("UserManagementService: Received response with status code: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, common.ErrUserNotFound
		case http.StatusUnauthorized:
			return nil, common.ErrUnauthorized
		case http.StatusForbidden:
			return nil, common.ErrForbidden
		case http.StatusTooManyRequests:
			return nil, common.ErrRateLimitExceeded
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			// Server errors (5xx) - these are temporary issues with the external API
			log.Printf("UserManagementService: External API server error (HTTP %d)", resp.StatusCode)
			return nil, fmt.Errorf("external API server error (HTTP %d): service temporarily unavailable", resp.StatusCode)
		default:
			log.Printf("UserManagementService: Unexpected API error (HTTP %d)", resp.StatusCode)
			return nil, fmt.Errorf("API error: %d", resp.StatusCode)
		}
	}

	// Parse the response according to the validate-token API spec
	var response struct {
		IsValid  bool `json:"isValid"`
		UserInfo struct {
			UserID   string `json:"userID"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
			RoleID   string `json:"roleID"`
			Status   string `json:"status"`
			IsActive bool   `json:"isActive"`
		} `json:"userInfo"`
	}

	decodeErr := json.NewDecoder(resp.Body).Decode(&response)
	if decodeErr != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	userID, err := uuid.Parse(response.UserInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Parse tenant ID if needed
	var tenantID uuid.UUID
	// Note: In the current response format, there might not be a tenantID
	// We'll need to update this if tenant information is added to the response

	validationResponse := &dto.TokenValidationResponse{
		Valid:    response.IsValid && response.UserInfo.IsActive,
		UserID:   userID,
		UserRole: response.UserInfo.Role,
		TenantID: tenantID, // Might be empty UUID if not in response
		Email:    response.UserInfo.Email,
	}

	return validationResponse, nil
}

// GetCurrentTenant retrieves the current tenant context for the authenticated user
func (s *UserManagementService) GetCurrentTenant(ctx context.Context, userID uuid.UUID) (*dto.UserTenantCurrentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/user-tenant/current", s.userTenantAPIURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-User-ID", userID.String())
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var response dto.UserTenantCurrentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetUserTenants retrieves all tenants that the user has access to
func (s *UserManagementService) GetUserTenants(ctx context.Context, userID uuid.UUID) (*dto.UserTenantsListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/user-tenant/tenants", s.userTenantAPIURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-User-ID", userID.String())
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var response dto.UserTenantsListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// SwitchTenant switches the user's current active tenant context
func (s *UserManagementService) SwitchTenant(ctx context.Context, userID, tenantID uuid.UUID) (*dto.SwitchTenantResponse, error) {
	url := fmt.Sprintf("%s/api/v1/user-tenant/switch", s.userTenantAPIURL)

	reqBody := dto.SwitchTenantRequest{
		TenantID: tenantID.String(),
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID.String())
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var response dto.SwitchTenantResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// ValidateUserTenantAccess validates if the user has access to a specific tenant
func (s *UserManagementService) ValidateUserTenantAccess(ctx context.Context, userID, tenantID uuid.UUID) (*dto.UserTenantAccessValidationResponse, error) {
	url := fmt.Sprintf("%s/api/v1/user-tenant/validate-access", s.userTenantAPIURL)

	reqBody := map[string]string{
		"tenant_id": tenantID.String(),
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID.String())
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var response dto.UserTenantAccessValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetTenantUsers retrieves all users within the admin's current tenant
func (s *UserManagementService) GetTenantUsers(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.TenantUsersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/user-tenant/users?page=%d&limit=%d", s.userTenantAPIURL, page, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-User-ID", userID.String())
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var response dto.TenantUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// ValidateSuperAdmin validates if the user has SuperAdmin role
func (s *UserManagementService) ValidateSuperAdmin(ctx context.Context, token string) (*dto.SuperAdminValidationResponse, error) {
	if !s.rateLimiter.Allow("validate_superadmin") {
		return nil, common.ErrRateLimitExceeded
	}

	url := fmt.Sprintf("%s%s", s.userAPIURL, s.validateSuperAdminEndpoint)
	log.Printf("ValidateSuperAdmin: Making request to %s", url)

	// Create request with Authorization header instead of request body
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("ValidateSuperAdmin: Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("X-API-Key", s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Send request with detailed logging
	log.Printf("ValidateSuperAdmin: Sending request with token: %s...", token[:10])
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("ValidateSuperAdmin: Request failed: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("ValidateSuperAdmin: Received response status: %d", resp.StatusCode)

	// Read and log the response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ValidateSuperAdmin: Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the response body
	log.Printf("ValidateSuperAdmin: Response body: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		log.Printf("ValidateSuperAdmin: API error with status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("API error: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var result dto.SuperAdminValidationResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		log.Printf("ValidateSuperAdmin: Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("ValidateSuperAdmin: Successfully validated. isSuperAdmin=%v, userID=%s, role=%s",
		result.IsSuperAdmin, result.UserID, result.UserRole)

	return &result, nil
}
