package auth

import (
	"testing"
	"time"

	"github.com/Gurren-Software/Anexis-Server/packages/database/models"
	"github.com/google/uuid"
)

func TestToUserResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)
	user := &models.User{
		BaseModel:    models.BaseModel{ID: uuid.New(), CreatedAt: createdAt},
		Name:         "Ada",
		Email:        "ada@example.com",
		StorageQuota: 1024,
		StorageUsed:  512,
	}

	resp := ToUserResponse(user)

	if resp.ID != user.ID || resp.Name != user.Name || resp.Email != user.Email {
		t.Fatalf("unexpected user response: %#v", resp)
	}
	if resp.StorageQuota != 1024 || resp.StorageUsed != 512 {
		t.Fatalf("unexpected storage fields: %#v", resp)
	}
	if resp.CreatedAt != createdAt.Format(time.RFC3339) {
		t.Fatalf("unexpected created_at %q", resp.CreatedAt)
	}
}
