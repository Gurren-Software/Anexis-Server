package main

import (
	"context"
	"log"
	"os"
	"strings"

	swaggerdocs "github.com/Gurren-Software/Anexis-Server/apps/api/docs"
	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/config"
	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/features/auth"
	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/features/backup"
	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/features/files"
	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/features/links"
	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/features/migration"
	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/infrastructure/http"
	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/infrastructure/http/middleware"
	"github.com/Gurren-Software/Anexis-Server/packages/database"
	"github.com/Gurren-Software/Anexis-Server/packages/database/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Anexis Cloud Storage API
// @version 1.0
// @description Cloud file storage server with multi-provider storage support
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@gurren-software.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format: Bearer {token}

// @securityDefinitions.apikey APIKeyAuth
// @in header
// @name X-API-Key
// @description Enter your API key for standalone mode

func main() {
	// Load environment variables - try multiple paths for monorepo
	envPaths := []string{".env", "../../.env", "../../../.env"}
	envLoaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Loaded environment from %s", path)
			envLoaded = true
			break
		}
	}
	if !envLoaded {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()
	configureSwaggerAuth(cfg.IsStandaloneMode())

	// Print server mode info
	if cfg.IsStandaloneMode() {
		log.Println("Running in STANDALONE mode (self-hosted)")
		log.Printf("Storage provider: %s", cfg.StorageProvider)
		if cfg.APIKey != "" {
			log.Println("API key authentication enabled")
		}
	} else {
		log.Println("Running in SAAS mode")
		log.Printf("Storage provider: %s", cfg.StorageProvider)
	}

	// Initialize database
	db, err := database.NewWithConfig(&database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Health check database
	if err := db.HealthCheck(); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}
	log.Println("Database connection established")

	// Initialize storage provider
	ctx := context.Background()
	storageProvider, err := NewStorageProvider(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	log.Printf("Storage provider (%s) initialized", cfg.StorageProvider)

	// Create HTTP server
	server := http.NewServer(cfg.ServerHost, cfg.ServerPort, cfg.Debug)
	router := server.Router()

	// Global middleware
	router.Use(middleware.CORS([]string{"*"}))
	router.Use(middleware.RequestID())

	// Health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "anexis-api",
			"mode":    cfg.ServerMode,
			"storage": cfg.StorageProvider,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	baseURL := "http://" + cfg.ServerAddress()

	// Create auth middleware based on server mode
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret, cfg.APIKey, cfg.IsStandaloneMode())
	optionalAuthMiddleware := middleware.OptionalAuthMiddleware(cfg.JWTSecret, cfg.APIKey, cfg.IsStandaloneMode())

	// Initialize repositories
	authRepo := auth.NewRepository(db.DB)
	filesRepo := files.NewRepository(db.DB)
	linksRepo := links.NewRepository(db.DB)
	migrationRepo := migration.NewRepository(db.DB)
	backupRepo := backup.NewRepository(db.DB)

	if cfg.IsStandaloneMode() {
		if err := ensureStandaloneUser(authRepo); err != nil {
			log.Fatalf("Failed to initialize standalone user: %v", err)
		}
	}

	// Initialize services
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.JWTExpirationHours)
	filesService := files.NewService(filesRepo, storageProvider, authRepo)
	linksService := links.NewService(linksRepo, filesRepo, storageProvider, baseURL)
	migrationService := migration.NewService(migrationRepo, storageProvider)
	backupService := backup.NewService(backupRepo, filesRepo, storageProvider, baseURL)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	filesHandler := files.NewHandler(filesService, authRepo)
	linksHandler := links.NewHandler(linksService)
	migrationHandler := migration.NewHandler(migrationService)
	backupHandler := backup.NewHandler(backupService)

	// Register routes
	auth.RegisterRoutes(v1, authHandler, cfg.JWTSecret, cfg.IsStandaloneMode())
	files.RegisterRoutes(v1, filesHandler, authMiddleware)
	links.RegisterRoutes(v1, linksHandler, authMiddleware, optionalAuthMiddleware)
	migration.RegisterRoutes(v1, migrationHandler, authMiddleware)
	backup.RegisterRoutes(v1, backupHandler, authMiddleware)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Starting server on %s", cfg.ServerAddress())
	log.Printf("Swagger docs available at http://%s/swagger/index.html", cfg.ServerAddress())

	// Start server with graceful shutdown
	server.StartWithGracefulShutdown()
}

func init() {
	// Create required directories
	dirs := []string{"./data/temp", "./data/uploads"}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				log.Printf("failed to create directory %s: %v", dir, err)
			}
		}
	}
}

func configureSwaggerAuth(standalone bool) {
	if !standalone {
		return
	}

	replacer := strings.NewReplacer(
		`"BearerAuth": []`, `"APIKeyAuth": []`,
		`"BearerAuth": {`, `"APIKeyAuth": {`,
		`"description": "Enter your bearer token in the format: Bearer {token}"`, `"description": "Enter your API key for standalone mode"`,
		`"name": "Authorization"`, `"name": "X-API-Key"`,
	)
	swaggerdocs.SwaggerInfo.SwaggerTemplate = replacer.Replace(swaggerdocs.SwaggerInfo.SwaggerTemplate)
}

func ensureStandaloneUser(repo *auth.Repository) error {
	userID := middleware.StandaloneUserID()
	user, err := repo.FindByID(userID)
	if err != nil || user != nil {
		return err
	}

	return repo.Create(&models.User{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		Name:         "Standalone User",
		Email:        "standalone@anexis.local",
		PasswordHash: "standalone-api-key-auth",
		StorageQuota: 5 * 1024 * 1024 * 1024,
	})
}
