package providers

import (
	"context"
	"testing"
)

func TestStubProviders(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		provider Provider
	}{
		{name: "amazon", provider: NewAmazonProvider("access", "refresh")},
		{name: "dropbox", provider: NewDropboxProvider("access", "refresh")},
		{name: "google", provider: NewGoogleProvider("access", "refresh")},
		{name: "microsoft", provider: NewMicrosoftProvider("access", "refresh")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := tt.provider.ListFiles(ctx)
			if err != nil {
				t.Fatalf("ListFiles returned error: %v", err)
			}
			if len(files) != 0 {
				t.Fatalf("expected no files from stub provider, got %d", len(files))
			}

			reader, err := tt.provider.DownloadFile(ctx, "file-id")
			if err != nil {
				t.Fatalf("DownloadFile returned error: %v", err)
			}
			if reader != nil {
				t.Fatalf("expected nil reader from stub provider")
			}

			info, err := tt.provider.GetFileInfo(ctx, "file-id")
			if err != nil {
				t.Fatalf("GetFileInfo returned error: %v", err)
			}
			if info != nil {
				t.Fatalf("expected nil file info from stub provider")
			}
		})
	}
}
