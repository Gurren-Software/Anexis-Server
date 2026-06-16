package files

import (
	"testing"
	"time"

	"github.com/Gurren-Software/Anexis-Server/packages/database/models"
	"github.com/google/uuid"
)

func TestToFileResponseFunction(t *testing.T) {
	id := uuid.New()
	parentID := uuid.New()
	createdAt := time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)

	resp := ToFileResponse(
		id,
		"stored.txt",
		"original.txt",
		"text/plain",
		42,
		false,
		&parentID,
		"notes",
		true,
		false,
		createdAt,
		updatedAt,
	)

	if resp.ID != id || resp.ParentID == nil || *resp.ParentID != parentID {
		t.Fatalf("unexpected ids: %#v", resp)
	}
	if resp.Name != "stored.txt" || resp.OriginalName != "original.txt" || resp.MimeType != "text/plain" {
		t.Fatalf("unexpected file fields: %#v", resp)
	}
	if resp.CreatedAt != createdAt.Format(time.RFC3339) || resp.UpdatedAt != updatedAt.Format(time.RFC3339) {
		t.Fatalf("unexpected timestamps: %#v", resp)
	}
}

func TestServiceToFileResponseDetectsFolder(t *testing.T) {
	createdAt := time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)
	file := &models.File{
		BaseModel:    models.BaseModel{ID: uuid.New(), CreatedAt: createdAt, UpdatedAt: createdAt},
		Name:         "docs",
		OriginalName: "docs",
		MimeType:     "inode/directory",
		Size:         0,
		Description:  "folder",
	}

	resp := (&Service{}).ToFileResponse(file)

	if !resp.IsFolder {
		t.Fatalf("expected folder response")
	}
	if resp.Description != "folder" {
		t.Fatalf("unexpected description %q", resp.Description)
	}
}
