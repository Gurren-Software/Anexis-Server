package models

import (
	"time"

	"gorm.io/gorm"
)

// MigrationStatus represents the state of a migration job
type MigrationStatus string

const (
	MigrationStatusPending   MigrationStatus = "pending"
	MigrationStatusRunning   MigrationStatus = "running"
	MigrationStatusCompleted MigrationStatus = "completed"
	MigrationStatusFailed    MigrationStatus = "failed"
	MigrationStatusCancelled MigrationStatus = "cancelled"
)

// ProviderType represents the cloud provider being migrated from
type ProviderType string

const (
	ProviderTypeGoogle    ProviderType = "google"
	ProviderTypeAmazon    ProviderType = "amazon"
	ProviderTypeMicrosoft ProviderType = "microsoft"
	ProviderTypeDropbox   ProviderType = "dropbox"
)

// MigrationJob represents a background job to migrate files from an external provider
type MigrationJob struct {
	gorm.Model
	UserID   uint            `gorm:"index;not null" json:"user_id"`
	Provider ProviderType    `gorm:"not null" json:"provider"`
	Status   MigrationStatus `gorm:"default:'pending'" json:"status"`

	// Progress tracking
	TotalFiles     int   `gorm:"default:0" json:"total_files"`
	ProcessedFiles int   `gorm:"default:0" json:"processed_files"`
	FailedFiles    int   `gorm:"default:0" json:"failed_files"`
	TotalBytes     int64 `gorm:"default:0" json:"total_bytes"`
	ProcessedBytes int64 `gorm:"default:0" json:"processed_bytes"`

	// Timing
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Error tracking
	LastError string `json:"last_error,omitempty"`

	// OAuth tokens (encrypted)
	AccessToken  string `gorm:"type:text" json:"-"`
	RefreshToken string `gorm:"type:text" json:"-"`

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

// IsComplete returns true if migration is finished (success or failure)
func (m *MigrationJob) IsComplete() bool {
	return m.Status == MigrationStatusCompleted ||
		m.Status == MigrationStatusFailed ||
		m.Status == MigrationStatusCancelled
}
