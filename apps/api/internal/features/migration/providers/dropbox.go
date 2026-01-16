package providers

import (
	"context"
	"io"
)

// DropboxProvider implements Provider for Dropbox
type DropboxProvider struct {
	accessToken  string
	refreshToken string
}

// NewDropboxProvider creates a new Dropbox provider
func NewDropboxProvider(accessToken, refreshToken string) *DropboxProvider {
	return &DropboxProvider{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

// ListFiles lists all files from Dropbox
func (p *DropboxProvider) ListFiles(ctx context.Context) ([]*FileInfo, error) {
	// TODO: Implement Dropbox API integration
	return []*FileInfo{}, nil
}

// DownloadFile downloads a file from Dropbox
func (p *DropboxProvider) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
	// TODO: Implement Dropbox API download
	return nil, nil
}

// GetFileInfo gets info about a file from Dropbox
func (p *DropboxProvider) GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error) {
	// TODO: Implement Dropbox API file info
	return nil, nil
}

var _ Provider = (*DropboxProvider)(nil)
