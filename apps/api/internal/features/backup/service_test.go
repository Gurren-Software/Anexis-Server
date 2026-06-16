package backup

import (
	"testing"
	"time"

	"github.com/Gurren-Software/Anexis-Server/packages/database/models"
	"github.com/google/uuid"
)

func TestToBackupResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)
	startedAt := createdAt.Add(time.Minute)
	completedAt := createdAt.Add(time.Hour)
	expiresAt := completedAt.Add(7 * 24 * time.Hour)
	job := &models.BackupJob{
		BaseModel:      models.BaseModel{ID: uuid.New(), CreatedAt: createdAt},
		Type:           models.BackupTypeExport,
		Status:         models.BackupStatusCompleted,
		ArchiveKey:     "backups/user/archive.zip",
		ArchiveSize:    2048,
		TotalFiles:     4,
		ProcessedFiles: 2,
		StartedAt:      &startedAt,
		CompletedAt:    &completedAt,
		ExpiresAt:      &expiresAt,
	}

	resp := (&Service{}).ToBackupResponse(job)

	if resp.ID != job.ID || resp.Type != "export" || resp.Status != "completed" {
		t.Fatalf("unexpected backup response: %#v", resp)
	}
	if resp.Progress != 50 {
		t.Fatalf("expected 50%% progress, got %f", resp.Progress)
	}
	if resp.StartedAt != startedAt.Format(time.RFC3339) || resp.CompletedAt != completedAt.Format(time.RFC3339) {
		t.Fatalf("unexpected timestamps: %#v", resp)
	}
}
