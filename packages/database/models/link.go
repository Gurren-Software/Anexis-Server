package models

import (
	"time"

	"github.com/google/uuid"
)

// LinkType represents the type of access link
type LinkType string

const (
	LinkTypeDownload  LinkType = "download"
	LinkTypeStream    LinkType = "stream"
	LinkTypePermanent LinkType = "permanent"
	LinkTypeTemporal  LinkType = "temporal"
)

// LinkAccessType represents the access level of a link
type LinkAccessType string

const (
	LinkAccessPublic     LinkAccessType = "public"
	LinkAccessPrivate    LinkAccessType = "private"
	LinkAccessRestricted LinkAccessType = "restricted"
)

// Link represents an access link for a file
type Link struct {
	BaseModel
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	FileID uuid.UUID `gorm:"type:uuid;not null;index" json:"file_id"`
	Token  string    `gorm:"uniqueIndex;not null" json:"token"`
	Type   LinkType  `gorm:"not null" json:"type"`

	// Access control
	AccessType    LinkAccessType `gorm:"default:public" json:"access_type"`
	Password      *string        `json:"-"`
	MaxDownloads  *int           `json:"max_downloads,omitempty"`
	DownloadCount int            `gorm:"default:0" json:"download_count"`
	AllowedEmails string         `json:"allowed_emails,omitempty"`

	// Expiration
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`

	// Metadata
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
	File File `gorm:"foreignKey:FileID" json:"file,omitempty"`
}

// TableName overrides the table name
func (Link) TableName() string {
	return "links"
}

// IsExpired returns true if the link has expired
func (l *Link) IsExpired() bool {
	if l.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*l.ExpiresAt)
}

// HasDownloadsRemaining returns true if downloads are still available
func (l *Link) HasDownloadsRemaining() bool {
	if l.MaxDownloads == nil {
		return true
	}
	return l.DownloadCount < *l.MaxDownloads
}

// GetID returns the link ID
func (l *Link) GetID() uuid.UUID {
	return l.ID
}
