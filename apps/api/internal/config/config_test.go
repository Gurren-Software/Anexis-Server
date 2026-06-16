package config

import (
	"testing"
	"time"
)

func TestLoadReadsEnvironment(t *testing.T) {
	t.Setenv("SERVER_HOST", "127.0.0.1")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("DEBUG", "false")
	t.Setenv("SERVER_MODE", "standalone")
	t.Setenv("ANEXIS_API_KEY", "test-key")
	t.Setenv("STORAGE_PROVIDER", "s3")
	t.Setenv("S3_ENDPOINT", "https://s3.example.com")
	t.Setenv("S3_BUCKET", "bucket")
	t.Setenv("S3_ACCESS_KEY", "access")
	t.Setenv("S3_SECRET_KEY", "secret")
	t.Setenv("S3_FORCE_PATH_STYLE", "true")
	t.Setenv("RATE_LIMIT_WINDOW_SECONDS", "30")
	t.Setenv("MAX_UPLOAD_SIZE_MB", "2")

	cfg := Load()

	if cfg.ServerAddress() != "127.0.0.1:9090" {
		t.Fatalf("expected server address 127.0.0.1:9090, got %q", cfg.ServerAddress())
	}
	if !cfg.IsProduction() || cfg.IsDevelopment() {
		t.Fatalf("expected production environment")
	}
	if cfg.Debug {
		t.Fatalf("expected debug to be false")
	}
	if !cfg.IsStandaloneMode() || cfg.IsSaaSMode() {
		t.Fatalf("expected standalone mode")
	}
	if !cfg.IsS3Storage() || !cfg.IsStorageConfigured() {
		t.Fatalf("expected configured S3 storage")
	}
	if !cfg.S3ForcePathStyle {
		t.Fatalf("expected S3 force path style")
	}
	if cfg.RateLimitWindow != 30*time.Second {
		t.Fatalf("expected 30s rate limit window, got %s", cfg.RateLimitWindow)
	}
	if cfg.MaxUploadSize != 2*1024*1024 {
		t.Fatalf("expected 2 MiB max upload, got %d", cfg.MaxUploadSize)
	}
	if !cfg.NeedsAuth() {
		t.Fatalf("expected auth to be required")
	}
}

func TestStorageConfiguredByProvider(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want bool
	}{
		{
			name: "local requires path",
			cfg:  Config{StorageProvider: "local", StorageLocalPath: "/tmp/storage"},
			want: true,
		},
		{
			name: "b2 requires credentials",
			cfg: Config{
				StorageProvider:  "b2",
				B2KeyID:          "key-id",
				B2ApplicationKey: "app-key",
				B2BucketName:     "bucket",
			},
			want: true,
		},
		{
			name: "s3 missing secret",
			cfg: Config{
				StorageProvider: "s3",
				S3Endpoint:      "https://s3.example.com",
				S3Bucket:        "bucket",
				S3AccessKey:     "access",
			},
			want: false,
		},
		{
			name: "unknown provider",
			cfg:  Config{StorageProvider: "other"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.IsStorageConfigured(); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
