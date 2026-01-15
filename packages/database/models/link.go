package models

import (
	"time"

	"gorm.io/gorm"
)

// LinkType represents the type of access link
type LinkType string

const (
	LinkTypeDownload  LinkType = "download"  // Direct download
	LinkTypeStream    LinkType = "stream"    // Streaming access
	LinkTypePermanent LinkType = "permanent" // Never expires
	LinkTypeTemporal  LinkType = "temporal"  // Time-limited access
)

// LinkAccessType defines who can access the link
type LinkAccessType string

const (
	LinkAccessPublic     LinkAccessType = "public"     // Anyone with link
	LinkAccessPrivate    LinkAccessType = "private"    // Owner only
	LinkAccessRestricted LinkAccessType = "restricted" // Specific users/passwords
)

// Link represents an access link to a file
type Link struct {
	gorm.Model
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	FileID     uint           `gorm:"index;not null" json:"file_id"`
	Token      string         `gorm:"uniqueIndex;not null" json:"token"`
	Type       LinkType       `gorm:"not null" json:"type"`
	AccessType LinkAccessType `gorm:"default:'public'" json:"access_type"`

	// Access control
	Password       *string    `json:"-"`                       // Optional password protection
	MaxDownloads   *int       `json:"max_downloads,omitempty"` // Max number of downloads
	DownloadCount  int        `gorm:"default:0" json:"download_count"`
	AllowedEmails  string     `json:"allowed_emails,omitempty"` // Comma-separated allowed emails
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`     // For temporal links
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`

	// Metadata
	Name        string `json:"name,omitempty"`        // Custom link name
	Description string `json:"description,omitempty"` // Link description

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
	File File `gorm:"foreignKey:FileID" json:"file,omitempty"`
}

// TableName overrides the table name
func (Link) TableName() string {
	return "links"
}

// IsExpired checks if the link has expired
func (l *Link) IsExpired() bool {
	if l.Type == LinkTypePermanent || l.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*l.ExpiresAt)
}

// HasDownloadsRemaining checks if downloads are still allowed
func (l *Link) HasDownloadsRemaining() bool {
	if l.MaxDownloads == nil {
		return true
	}
	return l.DownloadCount < *l.MaxDownloads
}

// CanAccess checks if the link is still valid for access
func (l *Link) CanAccess() bool {
	return !l.IsExpired() && l.HasDownloadsRemaining()
}
