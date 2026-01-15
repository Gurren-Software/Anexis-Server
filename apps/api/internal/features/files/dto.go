package files

import "mime/multipart"

// UploadRequest represents file upload metadata
type UploadRequest struct {
	ParentID    *uint  `form:"parent_id"`
	Description string `form:"description"`
	Tags        string `form:"tags"`
	Encrypt     bool   `form:"encrypt"`
	Compress    bool   `form:"compress"`
}

// FileResponse represents file data in responses
type FileResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	Size         int64  `json:"size"`
	Status       string `json:"status"`
	IsEncrypted  bool   `json:"is_encrypted"`
	IsCompressed bool   `json:"is_compressed"`
	Description  string `json:"description,omitempty"`
	Tags         string `json:"tags,omitempty"`
	UploadedAt   string `json:"uploaded_at"`
	ParentID     *uint  `json:"parent_id,omitempty"`
}

// ListFilesRequest represents file listing parameters
type ListFilesRequest struct {
	ParentID *uint  `form:"parent_id"`
	Search   string `form:"search"`
	Page     int    `form:"page" binding:"min=1"`
	PerPage  int    `form:"per_page" binding:"min=1,max=100"`
}

// CreateFolderRequest represents folder creation payload
type CreateFolderRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=255"`
	ParentID *uint  `json:"parent_id"`
}

// RenameRequest represents file/folder rename payload
type RenameRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// MoveRequest represents file/folder move payload
type MoveRequest struct {
	TargetParentID *uint `json:"target_parent_id"`
}

// MultipartFormFile wraps file upload data
type MultipartFormFile struct {
	File   *multipart.FileHeader
	Upload *UploadRequest
}

// FileListResponse represents paginated file list
type FileListResponse struct {
	Files      []FileResponse `json:"files"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
}
