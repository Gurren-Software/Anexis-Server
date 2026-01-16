package links

import (
	"time"

	"github.com/google/uuid"
)

// CreateLinkRequest represents link creation payload
type CreateLinkRequest struct {
	FileID       uuid.UUID `json:"file_id" binding:"required"`
	Type         string    `json:"type" binding:"required,oneof=download stream permanent temporal"`
	AccessType   string    `json:"access_type" binding:"oneof=public private restricted"`
	Password     *string   `json:"password,omitempty"`
	MaxDownloads *int      `json:"max_downloads,omitempty"`
	ExpiresIn    *int      `json:"expires_in,omitempty"` // seconds until expiration
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
}

// LinkResponse represents link data in responses
type LinkResponse struct {
	ID             uuid.UUID  `json:"id"`
	Token          string     `json:"token"`
	Type           string     `json:"type"`
	AccessType     string     `json:"access_type"`
	URL            string     `json:"url"`
	MaxDownloads   *int       `json:"max_downloads,omitempty"`
	DownloadCount  int        `json:"download_count"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
	Name           string     `json:"name,omitempty"`
	Description    string     `json:"description,omitempty"`
	FileID         uuid.UUID  `json:"file_id"`
	FileName       string     `json:"file_name"`
	CreatedAt      string     `json:"created_at"`
}

// UpdateLinkRequest represents link update payload
type UpdateLinkRequest struct {
	MaxDownloads *int    `json:"max_downloads,omitempty"`
	ExpiresIn    *int    `json:"expires_in,omitempty"` // seconds to extend
	Name         string  `json:"name,omitempty"`
	Description  string  `json:"description,omitempty"`
	Password     *string `json:"password,omitempty"`
}

// ValidateLinkRequest represents link validation payload
type ValidateLinkRequest struct {
	Password string `json:"password,omitempty"`
}

// ListLinksRequest represents link listing parameters
type ListLinksRequest struct {
	FileID  *uuid.UUID `form:"file_id"`
	Type    string     `form:"type"`
	Page    int        `form:"page" binding:"min=1"`
	PerPage int        `form:"per_page" binding:"min=1,max=100"`
}
