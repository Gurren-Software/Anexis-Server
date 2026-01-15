package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// User represents a user account in the system
type User struct {
	gorm.Model
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	Name         string         `gorm:"not null" json:"name"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Age          uint8          `json:"age,omitempty"`
	Birthday     *time.Time     `json:"birthday,omitempty"`
	MemberNumber sql.NullString `gorm:"uniqueIndex" json:"member_number,omitempty"`
	ActivatedAt  sql.NullTime   `json:"activated_at,omitempty"`
	StorageQuota int64          `gorm:"default:5368709120" json:"storage_quota"` // 5GB default
	StorageUsed  int64          `gorm:"default:0" json:"storage_used"`

	// Relationships
	Files []File `gorm:"foreignKey:UserID" json:"files,omitempty"`
	Links []Link `gorm:"foreignKey:UserID" json:"links,omitempty"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}

// AvailableStorage returns the remaining storage quota
func (u *User) AvailableStorage() int64 {
	return u.StorageQuota - u.StorageUsed
}
