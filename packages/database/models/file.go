package models

import (
	"time"

	"github.com/google/uuid"
)

// FileStatus represents the processing status of a file
type FileStatus string

const (
	FileStatusPending   FileStatus = "pending"
	FileStatusUploading FileStatus = "uploading"
	FileStatusReady     FileStatus = "ready"
	FileStatusFailed    FileStatus = "failed"
	FileStatusDeleted   FileStatus = "deleted"
)

// File represents a file stored in the system
type File struct {
	BaseModel
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Name         string     `gorm:"not null" json:"name"`
	OriginalName string     `gorm:"not null" json:"original_name"`
	MimeType     string     `gorm:"not null" json:"mime_type"`
	Size         int64      `gorm:"not null" json:"size"`
	StorageKey   string     `gorm:"uniqueIndex;not null" json:"storage_key"`
	StoragePath  string     `gorm:"not null" json:"storage_path"`
	Checksum     string     `gorm:"not null" json:"checksum"`
	Status       FileStatus `gorm:"default:pending" json:"status"`
	IsEncrypted  bool       `gorm:"default:false" json:"is_encrypted"`
	IsCompressed bool       `gorm:"default:false" json:"is_compressed"`

	// Folder support
	ParentID *uuid.UUID `gorm:"type:uuid;index" json:"parent_id,omitempty"`

	// Metadata
	Description string     `json:"description,omitempty"`
	Tags        string     `json:"tags,omitempty"`
	UploadedAt  *time.Time `json:"uploaded_at,omitempty"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`

	// Relationships
	User   User   `gorm:"foreignKey:UserID" json:"-"`
	Links  []Link `gorm:"foreignKey:FileID" json:"links,omitempty"`
	Parent *File  `gorm:"foreignKey:ParentID" json:"-"`
}

// TableName overrides the table name
func (File) TableName() string {
	return "files"
}

// IsFolder returns true if this file is a folder
func (f *File) IsFolder() bool {
	return f.MimeType == "application/x-directory" || f.MimeType == "inode/directory"
}

// GetID returns the file ID
func (f *File) GetID() uuid.UUID {
	return f.ID
}
