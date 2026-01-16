package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds application configuration
type Config struct {
	// Server settings
	ServerHost  string
	ServerPort  string
	Environment string
	Debug       bool

	// JWT settings
	JWTSecret          string
	JWTExpirationHours int

	// Backblaze B2 settings
	B2KeyID          string
	B2ApplicationKey string
	B2BucketName     string
	B2BucketID       string

	// Database settings (passed to database package)
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Rate limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// File upload limits
	MaxUploadSize int64 // in bytes
}

// Load returns the application configuration from environment variables
func Load() *Config {
	return &Config{
		// Server
		ServerHost:  getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		Debug:       getEnvBool("DEBUG", true),

		// JWT
		JWTSecret:          getEnv("JWT_SECRET", "change-this-secret-in-production"),
		JWTExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),

		// Backblaze B2
		B2KeyID:          getEnv("B2_APPLICATION_KEY_ID", ""),
		B2ApplicationKey: getEnv("B2_APPLICATION_KEY", ""),
		B2BucketName:     getEnv("B2_BUCKET_NAME", ""),
		B2BucketID:       getEnv("B2_BUCKET_ID", ""),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "anexis"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Rate limiting
		RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   time.Duration(getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,

		// File upload
		MaxUploadSize: int64(getEnvInt("MAX_UPLOAD_SIZE_MB", 100)) * 1024 * 1024,
	}
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// ServerAddress returns the full server address
func (c *Config) ServerAddress() string {
	return c.ServerHost + ":" + c.ServerPort
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}
