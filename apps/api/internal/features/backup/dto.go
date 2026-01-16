package backup

import "github.com/google/uuid"

// StartBackupRequest represents backup job creation payload
type StartBackupRequest struct {
	Type string `json:"type" binding:"required,oneof=export import"`
}

// ImportBackupRequest represents import/restore payload
type ImportBackupRequest struct {
	ArchiveKey string `json:"archive_key" binding:"required"`
}

// BackupResponse represents backup job data in responses
type BackupResponse struct {
	ID             uuid.UUID `json:"id"`
	Type           string    `json:"type"`
	Status         string    `json:"status"`
	ArchiveKey     string    `json:"archive_key,omitempty"`
	ArchiveSize    int64     `json:"archive_size,omitempty"`
	DownloadURL    string    `json:"download_url,omitempty"`
	TotalFiles     int       `json:"total_files"`
	ProcessedFiles int       `json:"processed_files"`
	Progress       float64   `json:"progress"`
	StartedAt      string    `json:"started_at,omitempty"`
	CompletedAt    string    `json:"completed_at,omitempty"`
	ExpiresAt      string    `json:"expires_at,omitempty"`
	LastError      string    `json:"last_error,omitempty"`
	CreatedAt      string    `json:"created_at"`
}

// ListBackupsRequest represents backup listing parameters
type ListBackupsRequest struct {
	Type    string `form:"type"`
	Status  string `form:"status"`
	Page    int    `form:"page" binding:"min=1"`
	PerPage int    `form:"per_page" binding:"min=1,max=100"`
}
