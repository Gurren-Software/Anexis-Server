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

	// Server mode: "saas" (default) or "standalone"
	ServerMode string

	// JWT settings
	JWTSecret          string
	JWTExpirationHours int

	// API Key for standalone mode
	APIKey string

	// Storage provider: "b2", "s3", or "local" (default: "local" for standalone, "b2" for saas)
	StorageProvider string

	// Backblaze B2 settings
	B2KeyID          string
	B2ApplicationKey string
	B2BucketName     string
	B2BucketID       string

	// S3-compatible storage settings
	S3Endpoint       string
	S3Region         string
	S3Bucket         string
	S3AccessKey      string
	S3SecretKey      string
	S3ForcePathStyle bool
	S3BasePath       string

	// Local storage settings
	StorageLocalPath string

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

		// Server mode
		ServerMode: getEnv("SERVER_MODE", "saas"),

		// JWT
		JWTSecret:          getEnv("JWT_SECRET", "change-this-secret-in-production"),
		JWTExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),

		// API Key for standalone mode
		APIKey: getEnv("ANEXIS_API_KEY", ""),

		// Storage provider
		StorageProvider: getEnv("STORAGE_PROVIDER", "local"),

		// Backblaze B2
		B2KeyID:          getEnv("B2_APPLICATION_KEY_ID", ""),
		B2ApplicationKey: getEnv("B2_APPLICATION_KEY", ""),
		B2BucketName:     getEnv("B2_BUCKET_NAME", ""),
		B2BucketID:       getEnv("B2_BUCKET_ID", ""),

		// S3-compatible storage
		S3Endpoint:       getEnv("S3_ENDPOINT", ""),
		S3Region:         getEnv("S3_REGION", "us-east-1"),
		S3Bucket:         getEnv("S3_BUCKET", ""),
		S3AccessKey:      getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:      getEnv("S3_SECRET_KEY", ""),
		S3ForcePathStyle: getEnvBool("S3_FORCE_PATH_STYLE", false),
		S3BasePath:       getEnv("S3_BASE_PATH", ""),

		// Local storage
		StorageLocalPath: getEnv("STORAGE_LOCAL_PATH", "./data/storage"),

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

// IsSaaSMode returns true if running in SaaS mode
func (c *Config) IsSaaSMode() bool {
	return c.ServerMode == "saas"
}

// IsStandaloneMode returns true if running in standalone mode
func (c *Config) IsStandaloneMode() bool {
	return c.ServerMode == "standalone"
}

// IsLocalStorage returns true if using local storage
func (c *Config) IsLocalStorage() bool {
	return c.StorageProvider == "local"
}

// IsB2Storage returns true if using Backblaze B2
func (c *Config) IsB2Storage() bool {
	return c.StorageProvider == "b2"
}

// IsS3Storage returns true if using S3-compatible storage
func (c *Config) IsS3Storage() bool {
	return c.StorageProvider == "s3"
}

// IsStorageConfigured returns true if storage is properly configured
func (c *Config) IsStorageConfigured() bool {
	switch c.StorageProvider {
	case "local":
		return c.StorageLocalPath != ""
	case "b2":
		return c.B2KeyID != "" && c.B2ApplicationKey != "" && c.B2BucketName != ""
	case "s3":
		return c.S3Endpoint != "" && c.S3Bucket != "" && c.S3AccessKey != "" && c.S3SecretKey != ""
	}
	return false
}

// NeedsAuth returns true if the server requires authentication
func (c *Config) NeedsAuth() bool {
	return c.IsSaaSMode() || c.APIKey != ""
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
