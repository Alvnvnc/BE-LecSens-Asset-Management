package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// JWTPayload represents the payload of a JWT token
type JWTPayload struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	RoleID   string `json:"role_id"`
	RoleName string `json:"role_name"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
	Nbf      int64  `json:"nbf"`
}

// DecodeJWTToken decodes a JWT token without validating the signature
// This is for testing purposes only - in production you should validate the signature
func DecodeJWTToken(tokenString string) (*JWTPayload, error) {
	// Split the token into its parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT token format")
	}

	// Decode the payload (second part)
	payload := parts[1]

	// Add padding if necessary
	if len(payload)%4 != 0 {
		payload += strings.Repeat("=", 4-len(payload)%4)
	}

	// Base64 decode
	decodedPayload, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	// Parse JSON
	var jwtPayload JWTPayload
	if err := json.Unmarshal(decodedPayload, &jwtPayload); err != nil {
		return nil, fmt.Errorf("failed to parse JWT payload: %w", err)
	}

	return &jwtPayload, nil
}

// Base64Encode encodes a byte slice to base64 URL encoding (for testing purposes)
func Base64Encode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// DecodeJWT decodes and validates JWT token
func DecodeJWT(tokenString, secretKey string) (*JWTClaims, error) {
	return ValidateToken(tokenString, secretKey)
}
