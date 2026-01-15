package providers

import (
	"context"
	"io"
)

// AmazonProvider implements Provider for Amazon S3/Drive
type AmazonProvider struct {
	accessToken  string
	refreshToken string
}

// NewAmazonProvider creates a new Amazon provider
func NewAmazonProvider(accessToken, refreshToken string) *AmazonProvider {
	return &AmazonProvider{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

// ListFiles lists all files from Amazon
func (p *AmazonProvider) ListFiles(ctx context.Context) ([]*FileInfo, error) {
	// TODO: Implement Amazon Drive/S3 API integration
	return []*FileInfo{}, nil
}

// DownloadFile downloads a file from Amazon
func (p *AmazonProvider) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
	// TODO: Implement Amazon Drive/S3 download
	return nil, nil
}

// GetFileInfo gets info about a file from Amazon
func (p *AmazonProvider) GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error) {
	// TODO: Implement Amazon Drive/S3 file info
	return nil, nil
}

var _ Provider = (*AmazonProvider)(nil)
