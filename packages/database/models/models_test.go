package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestFileHelpers(t *testing.T) {
	id := uuid.New()
	file := File{
		BaseModel: BaseModel{ID: id},
		MimeType:  "application/x-directory",
	}

	if file.TableName() != "files" {
		t.Fatalf("unexpected table name %q", file.TableName())
	}
	if !file.IsFolder() {
		t.Fatalf("expected directory mime type to be folder")
	}
	if file.GetID() != id {
		t.Fatalf("expected file id %s, got %s", id, file.GetID())
	}
}

func TestLinkExpirationAndDownloads(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	limit := 2

	link := Link{
		ExpiresAt:     &past,
		MaxDownloads:  &limit,
		DownloadCount: 2,
	}

	if link.TableName() != "links" {
		t.Fatalf("unexpected table name %q", link.TableName())
	}
	if !link.IsExpired() {
		t.Fatalf("expected expired link")
	}
	if link.HasDownloadsRemaining() {
		t.Fatalf("expected downloads to be exhausted")
	}

	link.MaxDownloads = nil
	if !link.HasDownloadsRemaining() {
		t.Fatalf("expected unlimited downloads when max downloads is nil")
	}
}

func TestBackupProgress(t *testing.T) {
	job := BackupJob{TotalFiles: 4, ProcessedFiles: 1}

	if job.TableName() != "backup_jobs" {
		t.Fatalf("unexpected table name %q", job.TableName())
	}
	if got := job.Progress(); got != 25 {
		t.Fatalf("expected 25%% progress, got %f", got)
	}

	job.TotalFiles = 0
	if got := job.Progress(); got != 0 {
		t.Fatalf("expected 0 progress with no total, got %f", got)
	}
}

func TestMigrationProgress(t *testing.T) {
	job := MigrationJob{TotalFiles: 8, ProcessedFiles: 2}

	if job.TableName() != "migration_jobs" {
		t.Fatalf("unexpected table name %q", job.TableName())
	}
	if got := job.Progress(); got != 25 {
		t.Fatalf("expected 25%% progress, got %f", got)
	}

	job.TotalFiles = 0
	if got := job.Progress(); got != 0 {
		t.Fatalf("expected 0 progress with no total, got %f", got)
	}
}
