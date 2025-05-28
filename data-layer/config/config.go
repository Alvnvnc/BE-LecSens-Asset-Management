// filepath: /home/alvn/Documents/playground/kp/be-lecsens/asset_management/main.go
package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Environment string
	Server      ServerConfig
	DB          DatabaseConfig
	Tenant      TenantConfig
	User        UserConfig
	JWT         JWTConfig
	Cloudinary  CloudinaryConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// TenantConfig holds tenant service configuration
type TenantConfig struct {
	APIURL string
	APIKey string
}

// UserConfig holds user management service configuration
type UserConfig struct {
	APIURL                      string
	APIKey                      string
	ValidateTokenEndpoint       string
	UserInfoEndpoint            string
	ValidatePermissionsEndpoint string
	ValidateSuperAdminEndpoint  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string
	Issuer    string
	ExpiresIn int
	Debug     bool
}

// CloudinaryConfig holds Cloudinary configuration
type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Environment: getEnvOrFail("ENVIRONMENT"),
		Server: ServerConfig{
			Port: getEnvOrFail("PORT"),
		},
		DB: DatabaseConfig{
			Host:     getEnvOrFail("DB_HOST"),
			Port:     getEnvOrFail("DB_PORT"),
			User:     getEnvOrFail("DB_USER"),
			Password: getEnvOrFail("DB_PASSWORD"),
			Name:     getEnvOrFail("DB_NAME"),
		},
		Tenant: TenantConfig{
			APIURL: getEnvOrFail("TENANT_API_URL"),
			APIKey: getEnvOrFail("TENANT_API_KEY"),
		},
		User: UserConfig{
			APIURL:                      getEnvOrFail("USER_API_URL"),
			APIKey:                      getEnvOrFail("USER_API_KEY"),
			ValidateTokenEndpoint:       getEnvOrFail("USER_AUTH_VALIDATE_TOKEN_ENDPOINT"),
			UserInfoEndpoint:            getEnvOrFail("USER_AUTH_USER_INFO_ENDPOINT"),
			ValidatePermissionsEndpoint: getEnvOrFail("USER_AUTH_VALIDATE_PERMISSIONS_ENDPOINT"),
			ValidateSuperAdminEndpoint:  getEnvOrFail("USER_AUTH_VALIDATE_SUPERADMIN_ENDPOINT"),
		},
		JWT: JWTConfig{
			SecretKey: getEnvOrFail("JWT_SECRET_KEY"),
			Issuer:    getEnvOrFail("JWT_ISSUER"),
			ExpiresIn: getEnvAsIntOrFail("JWT_EXPIRES_IN"),
			Debug:     getEnvAsBoolOrFail("JWT_DEBUG"),
		},
		Cloudinary: CloudinaryConfig{
			CloudName: getEnvOrFail("CLOUDINARY_CLOUD_NAME"),
			APIKey:    getEnvOrFail("CLOUDINARY_API_KEY"),
			APISecret: getEnvOrFail("CLOUDINARY_API_SECRET"),
		},
	}
}

// getEnvOrFail gets an environment variable or fails if not found
func getEnvOrFail(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable not set: %s", key)
	}
	return value
}

// getEnvAsIntOrFail gets an environment variable as an integer or fails if not found/invalid
func getEnvAsIntOrFail(key string) int {
	valueStr := getEnvOrFail(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Environment variable %s must be an integer: %v", key, err)
	}
	return value
}

// getEnvAsBoolOrFail gets an environment variable as a boolean or fails if not found/invalid
func getEnvAsBoolOrFail(key string) bool {
	valueStr := getEnvOrFail(key)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Fatalf("Environment variable %s must be a boolean: %v", key, err)
	}
	return value
}
