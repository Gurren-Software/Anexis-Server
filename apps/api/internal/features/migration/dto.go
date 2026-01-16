package migration

import "github.com/google/uuid"

// StartMigrationRequest represents migration job creation payload
type StartMigrationRequest struct {
	Provider     string `json:"provider" binding:"required,oneof=google amazon microsoft dropbox"`
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// MigrationResponse represents migration job data in responses
type MigrationResponse struct {
	ID             uuid.UUID `json:"id"`
	Provider       string    `json:"provider"`
	Status         string    `json:"status"`
	TotalFiles     int       `json:"total_files"`
	ProcessedFiles int       `json:"processed_files"`
	FailedFiles    int       `json:"failed_files"`
	TotalBytes     int64     `json:"total_bytes"`
	ProcessedBytes int64     `json:"processed_bytes"`
	Progress       float64   `json:"progress"`
	StartedAt      string    `json:"started_at,omitempty"`
	CompletedAt    string    `json:"completed_at,omitempty"`
	LastError      string    `json:"last_error,omitempty"`
	CreatedAt      string    `json:"created_at"`
}

// ListMigrationsRequest represents migration listing parameters
type ListMigrationsRequest struct {
	Status  string `form:"status"`
	Page    int    `form:"page" binding:"min=1"`
	PerPage int    `form:"per_page" binding:"min=1,max=100"`
}
