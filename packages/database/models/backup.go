package models

import (
	"time"

	"github.com/google/uuid"
)

// BackupType represents the type of backup operation
type BackupType string

const (
	BackupTypeExport BackupType = "export"
	BackupTypeImport BackupType = "import"
)

// BackupStatus represents the status of a backup job
type BackupStatus string

const (
	BackupStatusPending   BackupStatus = "pending"
	BackupStatusRunning   BackupStatus = "running"
	BackupStatusCompleted BackupStatus = "completed"
	BackupStatusFailed    BackupStatus = "failed"
)

// BackupJob represents a backup/restore job
type BackupJob struct {
	BaseModel
	UserID uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	Type   BackupType   `gorm:"not null" json:"type"`
	Status BackupStatus `gorm:"default:pending" json:"status"`

	// Archive details
	ArchiveKey  string `json:"archive_key,omitempty"`
	ArchiveSize int64  `json:"archive_size,omitempty"`

	// Progress tracking
	TotalFiles     int   `gorm:"default:0" json:"total_files"`
	ProcessedFiles int   `gorm:"default:0" json:"processed_files"`
	TotalBytes     int64 `gorm:"default:0" json:"total_bytes"`
	ProcessedBytes int64 `gorm:"default:0" json:"processed_bytes"`

	// Timestamps
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`

	// Error tracking
	LastError string `json:"last_error,omitempty"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the table name
func (BackupJob) TableName() string {
	return "backup_jobs"
}

// Progress returns the backup progress as a percentage
func (b *BackupJob) Progress() float64 {
	if b.TotalFiles == 0 {
		return 0
	}
	return float64(b.ProcessedFiles) / float64(b.TotalFiles) * 100
}

// GetID returns the job ID
func (b *BackupJob) GetID() uuid.UUID {
	return b.ID
}
