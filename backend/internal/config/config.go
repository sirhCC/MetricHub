package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Host        string
	Port        int
	Environment string
	RequestTimeoutSeconds int

	// Database configuration
	DatabaseURL string
	RedisURL    string
	NATSURL     string

	// JWT configuration
	JWTSecret string

	// Logging configuration
	LogLevel string

	// Metrics configuration
	MetricsEnabled bool

	// Migration configuration
	AutoMigrate bool
}

// Load configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Host:        getEnvWithDefault("API_HOST", "0.0.0.0"),
		Port:        getEnvAsIntWithDefault("API_PORT", 8080),
		Environment: getEnvWithDefault("ENVIRONMENT", "development"),
		RequestTimeoutSeconds: getEnvAsIntWithDefault("REQUEST_TIMEOUT_SECONDS", 5),

		DatabaseURL: getEnvWithDefault("DATABASE_URL", "postgres://metrichub:metrichub_dev@localhost:5432/metrichub?sslmode=disable"),
		RedisURL:    getEnvWithDefault("REDIS_URL", "redis://localhost:6379"),
		NATSURL:     getEnvWithDefault("NATS_URL", "nats://localhost:4222"),

		JWTSecret: getEnvWithDefault("JWT_SECRET", "your-secret-key-change-in-production"),
		LogLevel:  getEnvWithDefault("LOG_LEVEL", "info"),

		MetricsEnabled: getEnvAsBoolWithDefault("METRICS_ENABLED", true),
		AutoMigrate:    getEnvAsBoolWithDefault("AUTO_MIGRATE", true),
	}

	// Validate required configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks that all required configuration is present and valid
func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.JWTSecret == "" || c.JWTSecret == "your-secret-key-change-in-production" {
		if c.Environment == "production" {
			return fmt.Errorf("JWT_SECRET must be set to a secure value in production")
		}
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("API_PORT must be between 1 and 65535")
	}

	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, strings.ToLower(c.LogLevel)) {
		return fmt.Errorf("LOG_LEVEL must be one of: %s", strings.Join(validLogLevels, ", "))
	}

	return nil
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// Helper functions for environment variable parsing

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBoolWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
