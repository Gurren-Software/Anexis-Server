package migration

import (
	"strings"
	"testing"
	"time"

	"github.com/Gurren-Software/Anexis-Server/packages/database/models"
	"github.com/google/uuid"
)

func TestToMigrationResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)
	startedAt := createdAt.Add(time.Minute)
	completedAt := createdAt.Add(time.Hour)
	job := &models.MigrationJob{
		BaseModel:      models.BaseModel{ID: uuid.New(), CreatedAt: createdAt},
		Provider:       models.ProviderTypeGoogle,
		Status:         models.MigrationStatusCompleted,
		TotalFiles:     10,
		ProcessedFiles: 7,
		FailedFiles:    1,
		TotalBytes:     1000,
		ProcessedBytes: 700,
		StartedAt:      &startedAt,
		CompletedAt:    &completedAt,
	}

	resp := ToMigrationResponse(job)

	if resp.Provider != "google" || resp.Status != "completed" {
		t.Fatalf("unexpected migration response: %#v", resp)
	}
	if resp.Progress != 70 {
		t.Fatalf("expected 70%% progress, got %f", resp.Progress)
	}
	if resp.StartedAt != startedAt.Format(time.RFC3339) || resp.CompletedAt != completedAt.Format(time.RFC3339) {
		t.Fatalf("unexpected timestamps: %#v", resp)
	}
}

func TestGenerateStorageKey(t *testing.T) {
	userID := uuid.New()
	key := generateStorageKey(userID, "report.pdf")

	if !strings.Contains(key, userID.String()) || !strings.HasSuffix(key, "/report.pdf") {
		t.Fatalf("unexpected storage key %q", key)
	}
}
