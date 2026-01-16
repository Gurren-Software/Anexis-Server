package providers

import (
	"context"
	"io"
)

// FileInfo represents file information from a cloud provider
type FileInfo struct {
	ID       string
	Name     string
	Path     string
	Size     int64
	MimeType string
	IsFolder bool
}

// Provider defines the interface for cloud storage providers
type Provider interface {
	// ListFiles returns all files from the provider
	ListFiles(ctx context.Context) ([]*FileInfo, error)

	// DownloadFile downloads a file from the provider
	DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error)

	// GetFileInfo returns info about a specific file
	GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error)
}
