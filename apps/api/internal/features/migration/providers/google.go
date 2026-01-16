package providers

import (
	"context"
	"io"
)

// GoogleProvider implements Provider for Google Drive
type GoogleProvider struct {
	accessToken  string
	refreshToken string
}

// NewGoogleProvider creates a new Google Drive provider
func NewGoogleProvider(accessToken, refreshToken string) *GoogleProvider {
	return &GoogleProvider{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

// ListFiles lists all files from Google Drive
func (p *GoogleProvider) ListFiles(ctx context.Context) ([]*FileInfo, error) {
	// TODO: Implement Google Drive API integration
	// This is a stub - actual implementation would use Google Drive API
	return []*FileInfo{}, nil
}

// DownloadFile downloads a file from Google Drive
func (p *GoogleProvider) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
	// TODO: Implement Google Drive API download
	return nil, nil
}

// GetFileInfo gets info about a file from Google Drive
func (p *GoogleProvider) GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error) {
	// TODO: Implement Google Drive API file info
	return nil, nil
}

// Ensure GoogleProvider implements Provider
var _ Provider = (*GoogleProvider)(nil)
