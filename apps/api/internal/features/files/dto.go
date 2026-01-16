package files

import (
	"time"

	"github.com/google/uuid"
)

// UploadRequest represents file upload metadata
type UploadRequest struct {
	Compress    bool       `form:"compress"`
	ParentID    *uuid.UUID `form:"parent_id"`
	Description string     `form:"description"`
}

// FileResponse represents file data in responses
type FileResponse struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	OriginalName string     `json:"original_name"`
	MimeType     string     `json:"mime_type"`
	Size         int64      `json:"size"`
	IsFolder     bool       `json:"is_folder"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	Description  string     `json:"description,omitempty"`
	IsEncrypted  bool       `json:"is_encrypted"`
	IsCompressed bool       `json:"is_compressed"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
}

// ListFilesRequest represents file listing parameters
type ListFilesRequest struct {
	ParentID *uuid.UUID `form:"parent_id"`
	Search   string     `form:"search"`
	Page     int        `form:"page" binding:"min=1"`
	PerPage  int        `form:"per_page" binding:"min=1,max=100"`
}

// CreateFolderRequest represents folder creation payload
type CreateFolderRequest struct {
	Name     string     `json:"name" binding:"required,min=1,max=255"`
	ParentID *uuid.UUID `json:"parent_id"`
}

// RenameRequest represents file/folder rename payload
type RenameRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// MoveRequest represents file/folder move payload
type MoveRequest struct {
	ParentID *uuid.UUID `json:"parent_id"`
}

// FileListResponse represents paginated file list
type FileListResponse struct {
	Files      []FileResponse `json:"files"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
}

// ToFileResponse converts file model to response
func ToFileResponse(id uuid.UUID, name, originalName, mimeType string, size int64, isFolder bool, parentID *uuid.UUID, description string, isEncrypted, isCompressed bool, createdAt, updatedAt time.Time) *FileResponse {
	return &FileResponse{
		ID:           id,
		Name:         name,
		OriginalName: originalName,
		MimeType:     mimeType,
		Size:         size,
		IsFolder:     isFolder,
		ParentID:     parentID,
		Description:  description,
		IsEncrypted:  isEncrypted,
		IsCompressed: isCompressed,
		CreatedAt:    createdAt.Format(time.RFC3339),
		UpdatedAt:    updatedAt.Format(time.RFC3339),
	}
}
