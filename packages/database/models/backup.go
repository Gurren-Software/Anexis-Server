package models

import (
	"time"

	"gorm.io/gorm"
)

// BackupStatus represents the state of a backup job
type BackupStatus string

const (
	BackupStatusPending   BackupStatus = "pending"
	BackupStatusRunning   BackupStatus = "running"
	BackupStatusCompleted BackupStatus = "completed"
	BackupStatusFailed    BackupStatus = "failed"
)

// BackupType represents the type of backup operation
type BackupType string

const (
	BackupTypeExport BackupType = "export" // Export data from server
	BackupTypeImport BackupType = "import" // Import/restore data to server
)

// BackupJob represents a backup or restore operation
type BackupJob struct {
	gorm.Model
	UserID uint         `gorm:"index;not null" json:"user_id"`
	Type   BackupType   `gorm:"not null" json:"type"`
	Status BackupStatus `gorm:"default:'pending'" json:"status"`

	// Archive details
	ArchiveKey  string `json:"archive_key,omitempty"`  // Storage key for backup archive
	ArchiveSize int64  `json:"archive_size,omitempty"` // Size in bytes

	// Progress tracking
	TotalFiles     int   `gorm:"default:0" json:"total_files"`
	ProcessedFiles int   `gorm:"default:0" json:"processed_files"`
	TotalBytes     int64 `gorm:"default:0" json:"total_bytes"`
	ProcessedBytes int64 `gorm:"default:0" json:"processed_bytes"`

	// Timing
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"` // When backup download expires

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

// IsComplete returns true if backup is finished
func (b *BackupJob) IsComplete() bool {
	return b.Status == BackupStatusCompleted || b.Status == BackupStatusFailed
}
