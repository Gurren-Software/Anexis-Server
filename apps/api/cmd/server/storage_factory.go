package main

import (
	"context"
	"fmt"
	"io"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/config"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage/backblaze"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage/local"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage/s3"
)

// LocalStorageProvider provides additional methods for local storage
type LocalStorageProvider interface {
	storage.Provider
	ServeFile(key string) (io.ReadCloser, int64, error)
	GetBasePath() string
}

// NewStorageProvider creates a new storage provider based on configuration
func NewStorageProvider(ctx context.Context, cfg *config.Config) (storage.Provider, error) {
	switch cfg.StorageProvider {
	case "local":
		client, err := local.NewClient(&local.Config{
			BasePath: cfg.StorageLocalPath,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize local storage: %w", err)
		}
		return client, nil

	case "b2":
		if cfg.B2KeyID == "" || cfg.B2ApplicationKey == "" || cfg.B2BucketName == "" {
			return nil, fmt.Errorf("B2 storage requires B2_APPLICATION_KEY_ID, B2_APPLICATION_KEY, and B2_BUCKET_NAME to be set")
		}
		client, err := backblaze.NewClient(ctx, &backblaze.Config{
			KeyID:          cfg.B2KeyID,
			ApplicationKey: cfg.B2ApplicationKey,
			BucketName:     cfg.B2BucketName,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize B2 storage: %w", err)
		}
		return client, nil

	case "s3":
		if cfg.S3Endpoint == "" || cfg.S3Bucket == "" || cfg.S3AccessKey == "" || cfg.S3SecretKey == "" {
			return nil, fmt.Errorf("S3 storage requires S3_ENDPOINT, S3_BUCKET, S3_ACCESS_KEY, and S3_SECRET_KEY to be set")
		}
		client, err := s3.NewClient(ctx, &s3.Config{
			Endpoint:       cfg.S3Endpoint,
			Region:         cfg.S3Region,
			Bucket:         cfg.S3Bucket,
			AccessKey:      cfg.S3AccessKey,
			SecretKey:      cfg.S3SecretKey,
			ForcePathStyle: cfg.S3ForcePathStyle,
			BasePath:       cfg.S3BasePath,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize S3 storage: %w", err)
		}
		return client, nil

	default:
		return nil, fmt.Errorf("unknown storage provider: %s", cfg.StorageProvider)
	}
}

// IsLocalStorage checks if the provider is a local storage provider
func IsLocalStorage(p storage.Provider) bool {
	_, ok := p.(LocalStorageProvider)
	return ok
}
