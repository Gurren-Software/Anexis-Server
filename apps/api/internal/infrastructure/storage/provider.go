package storage

import (
	"context"
	"io"
)

// Provider defines the interface for storage backends
type Provider interface {
	// Upload uploads a file to storage
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error

	// UploadLarge uploads a large file with concurrent upload support
	UploadLarge(ctx context.Context, key string, reader io.Reader, size int64, contentType string, concurrency int) error

	// Download downloads a file from storage
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// DownloadRange downloads a specific byte range of a file (for streaming)
	DownloadRange(ctx context.Context, key string, start, end int64) (io.ReadCloser, error)

	// Delete removes a file from storage
	Delete(ctx context.Context, key string) error

	// GetURL returns a temporary URL for direct download
	GetURL(ctx context.Context, key string, expiresIn int) (string, error)

	// GetStreamURL returns a URL suitable for streaming
	GetStreamURL(ctx context.Context, key string, expiresIn int) (string, error)

	// Exists checks if a file exists in storage
	Exists(ctx context.Context, key string) (bool, error)

	// GetMetadata returns file metadata
	GetMetadata(ctx context.Context, key string) (*FileMetadata, error)

	// List returns files under a given prefix
	List(ctx context.Context, prefix string, maxKeys int) ([]*FileMetadata, error)

	// Copy copies a file within storage
	Copy(ctx context.Context, srcKey, dstKey string) error
}

// FileMetadata represents storage file metadata
type FileMetadata struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified int64
	ETag         string
}

// UploadOptions provides additional options for upload operations
type UploadOptions struct {
	ContentType  string
	Metadata     map[string]string
	Concurrency  int
	PartSize     int64
	IsEncrypted  bool
	IsCompressed bool
}

// DownloadOptions provides additional options for download operations
type DownloadOptions struct {
	Start int64 // For range requests
	End   int64 // For range requests
}
