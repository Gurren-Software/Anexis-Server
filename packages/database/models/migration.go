package models

import (
	"time"

	"github.com/google/uuid"
)

// MigrationStatus represents the status of a migration job
type MigrationStatus string

const (
	MigrationStatusPending   MigrationStatus = "pending"
	MigrationStatusRunning   MigrationStatus = "running"
	MigrationStatusCompleted MigrationStatus = "completed"
	MigrationStatusFailed    MigrationStatus = "failed"
	MigrationStatusCancelled MigrationStatus = "cancelled"
)

// ProviderType represents the cloud provider for migration
type ProviderType string

const (
	ProviderTypeGoogle    ProviderType = "google"
	ProviderTypeAmazon    ProviderType = "amazon"
	ProviderTypeMicrosoft ProviderType = "microsoft"
	ProviderTypeDropbox   ProviderType = "dropbox"
)

// MigrationJob represents a migration job from an external provider
type MigrationJob struct {
	BaseModel
	UserID   uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	Provider ProviderType    `gorm:"not null" json:"provider"`
	Status   MigrationStatus `gorm:"default:pending" json:"status"`

	// Progress tracking
	TotalFiles     int   `gorm:"default:0" json:"total_files"`
	ProcessedFiles int   `gorm:"default:0" json:"processed_files"`
	FailedFiles    int   `gorm:"default:0" json:"failed_files"`
	TotalBytes     int64 `gorm:"default:0" json:"total_bytes"`
	ProcessedBytes int64 `gorm:"default:0" json:"processed_bytes"`

	// Timestamps
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Error tracking
	LastError string `json:"last_error,omitempty"`

	// OAuth tokens (encrypted in production)
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the table name
func (MigrationJob) TableName() string {
	return "migration_jobs"
}

// Progress returns the migration progress as a percentage
func (m *MigrationJob) Progress() float64 {
	if m.TotalFiles == 0 {
		return 0
	}
	return float64(m.ProcessedFiles) / float64(m.TotalFiles) * 100
}

// GetID returns the job ID
func (m *MigrationJob) GetID() uuid.UUID {
	return m.ID
}
