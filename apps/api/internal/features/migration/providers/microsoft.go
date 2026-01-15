package providers

import (
	"context"
	"io"
)

// MicrosoftProvider implements Provider for OneDrive
type MicrosoftProvider struct {
	accessToken  string
	refreshToken string
}

// NewMicrosoftProvider creates a new OneDrive provider
func NewMicrosoftProvider(accessToken, refreshToken string) *MicrosoftProvider {
	return &MicrosoftProvider{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

// ListFiles lists all files from OneDrive
func (p *MicrosoftProvider) ListFiles(ctx context.Context) ([]*FileInfo, error) {
	// TODO: Implement Microsoft Graph API integration
	return []*FileInfo{}, nil
}

// DownloadFile downloads a file from OneDrive
func (p *MicrosoftProvider) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
	// TODO: Implement Microsoft Graph API download
	return nil, nil
}

// GetFileInfo gets info about a file from OneDrive
func (p *MicrosoftProvider) GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error) {
	// TODO: Implement Microsoft Graph API file info
	return nil, nil
}

var _ Provider = (*MicrosoftProvider)(nil)
