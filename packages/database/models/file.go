package models

import (
	"time"

	"gorm.io/gorm"
)

// FileStatus represents the current state of a file
type FileStatus string

const (
	FileStatusPending    FileStatus = "pending"
	FileStatusUploading  FileStatus = "uploading"
	FileStatusProcessing FileStatus = "processing"
	FileStatusReady      FileStatus = "ready"
	FileStatusFailed     FileStatus = "failed"
	FileStatusDeleted    FileStatus = "deleted"
)

// File represents a stored file in the cloud storage
type File struct {
	gorm.Model
	UserID       uint       `gorm:"index;not null" json:"user_id"`
	Name         string     `gorm:"not null" json:"name"`
	OriginalName string     `gorm:"not null" json:"original_name"`
	MimeType     string     `gorm:"not null" json:"mime_type"`
	Size         int64      `gorm:"not null" json:"size"`
	StorageKey   string     `gorm:"uniqueIndex;not null" json:"storage_key"` // B2 object key
	StoragePath  string     `gorm:"not null" json:"storage_path"`            // Full path in bucket
	Checksum     string     `gorm:"not null" json:"checksum"`                // SHA256 hash
	Status       FileStatus `gorm:"default:'pending'" json:"status"`
	IsEncrypted  bool       `gorm:"default:false" json:"is_encrypted"`
	IsCompressed bool       `gorm:"default:false" json:"is_compressed"`
	ParentID     *uint      `gorm:"index" json:"parent_id,omitempty"` // For folder structure

	// Metadata
	Description string     `json:"description,omitempty"`
	Tags        string     `json:"tags,omitempty"` // Comma-separated tags
	UploadedAt  time.Time  `gorm:"autoCreateTime" json:"uploaded_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`

	// Relationships
	User  User   `gorm:"foreignKey:UserID" json:"-"`
	Links []Link `gorm:"foreignKey:FileID" json:"links,omitempty"`
}

// TableName overrides the table name
func (File) TableName() string {
	return "files"
}

// IsFolder returns true if this file entry represents a folder
func (f *File) IsFolder() bool {
	return f.MimeType == "application/x-directory"
}
