package common

import (
	"errors"
	"fmt"
	"net/http"
)

// Custom error types for better error handling
var (
	// Tenant errors
	ErrTenantNotFound        = errors.New("tenant not found")
	ErrTenantInactive        = errors.New("tenant is not active")
	ErrTenantSubscriptionExp = errors.New("tenant subscription has expired")

	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user is not active")
	ErrUserNoTenantAccess = errors.New("user does not have access to the tenant")

	// API errors
	ErrAPIConnectionFailed = errors.New("failed to connect to external API")
	ErrAPIResponseInvalid  = errors.New("invalid response from external API")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrForbidden           = errors.New("forbidden access")
	ErrRateLimitExceeded   = errors.New("rate limit exceeded")

	// Validation and business logic errors
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
)

// APIError represents an error from an external API
type APIError struct {
	StatusCode int
	Message    string
	Err        error
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("API Error (HTTP %d): %s - %v", e.StatusCode, e.Message, e.Err)
}

// Unwrap returns the underlying error
func (e *APIError) Unwrap() error {
	return e.Err
}

// NewAPIError creates a new API error
func NewAPIError(statusCode int, message string, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
	Err     error
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("validation error: %s - %v", e.Message, e.Err)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// Unwrap returns the underlying error
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// Is implements the error interface for error comparison
func (e *ValidationError) Is(target error) bool {
	return target == ErrValidation
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Resource string
	ID       string
	Err      error
}

// Error implements the error interface
func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// Unwrap returns the underlying error
func (e *NotFoundError) Unwrap() error {
	return e.Err
}

// Is implements the error interface for error comparison
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

// NewValidationError creates a new validation error
func NewValidationError(message string, err error) *ValidationError {
	return &ValidationError{
		Message: message,
		Err:     err,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// BadRequestError represents a bad request error
type BadRequestError struct {
	Message string
	Err     error
}

// Error implements the error interface
func (e *BadRequestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("bad request: %s - %v", e.Message, e.Err)
	}
	return fmt.Sprintf("bad request: %s", e.Message)
}

// Unwrap returns the underlying error
func (e *BadRequestError) Unwrap() error {
	return e.Err
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{
		Message: message,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation)
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// HTTPStatusFromError determines the HTTP status code from an error
func HTTPStatusFromError(err error) int {
	// Check for API errors first
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode
	}

	// Check for specific error types
	switch {
	case errors.Is(err, ErrTenantNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrTenantInactive), errors.Is(err, ErrTenantSubscriptionExp):
		return http.StatusForbidden
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrRateLimitExceeded):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrAPIConnectionFailed), errors.Is(err, ErrAPIResponseInvalid):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
