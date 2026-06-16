package links

import (
	"testing"
	"time"

	"github.com/Gurren-Software/Anexis-Server/packages/database/models"
	"github.com/google/uuid"
)

func TestToLinkResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)
	expiresAt := createdAt.Add(time.Hour)
	maxDownloads := 5
	fileID := uuid.New()
	link := &models.Link{
		BaseModel:     models.BaseModel{ID: uuid.New(), CreatedAt: createdAt},
		FileID:        fileID,
		Token:         "token-123",
		Type:          models.LinkTypeDownload,
		AccessType:    models.LinkAccessRestricted,
		MaxDownloads:  &maxDownloads,
		DownloadCount: 2,
		ExpiresAt:     &expiresAt,
		Name:          "Public link",
		Description:   "Shared file",
		File: models.File{
			BaseModel:    models.BaseModel{ID: fileID},
			OriginalName: "report.pdf",
		},
	}

	resp := (&Service{baseURL: "https://files.example.com"}).ToLinkResponse(link)

	if resp.URL != "https://files.example.com/api/v1/links/token-123" {
		t.Fatalf("unexpected url %q", resp.URL)
	}
	if resp.FileName != "report.pdf" || resp.FileID != fileID {
		t.Fatalf("unexpected file fields: %#v", resp)
	}
	if resp.MaxDownloads == nil || *resp.MaxDownloads != maxDownloads {
		t.Fatalf("unexpected max downloads: %#v", resp.MaxDownloads)
	}
}

func TestGenerateTokenLength(t *testing.T) {
	token, err := generateToken(24)
	if err != nil {
		t.Fatalf("generateToken failed: %v", err)
	}
	if len(token) != 24 {
		t.Fatalf("expected token length 24, got %d", len(token))
	}
}
