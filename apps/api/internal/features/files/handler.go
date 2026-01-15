package files

import (
	"errors"
	"io"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/response"
	"github.com/gin-gonic/gin"
)

// Handler handles file HTTP requests
type Handler struct {
	service  *Service
	authRepo AuthRepository
}

// NewHandler creates a new files handler
func NewHandler(service *Service, authRepo AuthRepository) *Handler {
	return &Handler{service: service, authRepo: authRepo}
}

// Upload godoc
// @Summary Upload a file
// @Description Upload a file to storage
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param parent_id formData int false "Parent folder ID"
// @Param description formData string false "File description"
// @Param tags formData string false "Comma-separated tags"
// @Param compress formData bool false "Compress file"
// @Success 201 {object} response.Response{data=FileResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/files/upload [post]
func (h *Handler) Upload(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "INVALID_FILE", "No file provided")
		return
	}

	var opts UploadRequest
	if err := c.ShouldBind(&opts); err != nil {
		// Ignore bind errors for optional fields
	}

	// Get user storage info (simplified - in real app, get from user context or cache)
	storageQuota := int64(5 * 1024 * 1024 * 1024) // 5GB default
	storageUsed := int64(0)                       // Would be from user record

	fileRecord, err := h.service.Upload(c.Request.Context(), userID, storageQuota, storageUsed, file, &opts)
	if err != nil {
		if errors.Is(err, ErrStorageQuotaExceeded) {
			response.BadRequest(c, "QUOTA_EXCEEDED", "Storage quota exceeded")
			return
		}
		response.InternalError(c, "Failed to upload file")
		return
	}

	response.Created(c, ToFileResponse(fileRecord))
}

// Download godoc
// @Summary Download a file
// @Description Download a file from storage
// @Tags files
// @Produce octet-stream
// @Security BearerAuth
// @Param id path int true "File ID"
// @Success 200 {file} binary
// @Failure 404 {object} response.Response
// @Router /api/v1/files/{id}/download [get]
func (h *Handler) Download(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	fileID := parseUint(c.Param("id"))
	if fileID == 0 {
		response.BadRequest(c, "INVALID_ID", "Invalid file ID")
		return
	}

	reader, file, err := h.service.Download(c.Request.Context(), userID, fileID)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			response.NotFound(c, "File not found")
			return
		}
		response.InternalError(c, "Failed to download file")
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", "attachment; filename=\""+file.OriginalName+"\"")
	c.Header("Content-Type", file.MimeType)

	if _, err := io.Copy(c.Writer, reader); err != nil {
		response.InternalError(c, "Failed to stream file")
		return
	}
}

// Get godoc
// @Summary Get file details
// @Description Get metadata for a file
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param id path int true "File ID"
// @Success 200 {object} response.Response{data=FileResponse}
// @Failure 404 {object} response.Response
// @Router /api/v1/files/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	fileID := parseUint(c.Param("id"))
	if fileID == 0 {
		response.BadRequest(c, "INVALID_ID", "Invalid file ID")
		return
	}

	file, err := h.service.GetFile(userID, fileID)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			response.NotFound(c, "File not found")
			return
		}
		response.InternalError(c, "Failed to get file")
		return
	}

	response.OK(c, ToFileResponse(file))
}

// List godoc
// @Summary List files
// @Description List user's files with optional filtering
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param parent_id query int false "Parent folder ID"
// @Param search query string false "Search query"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} response.Response{data=FileListResponse}
// @Router /api/v1/files [get]
func (h *Handler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req ListFilesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		// Use defaults
		req.Page = 1
		req.PerPage = 20
	}

	files, total, err := h.service.ListFiles(userID, &req)
	if err != nil {
		response.InternalError(c, "Failed to list files")
		return
	}

	fileResponses := make([]FileResponse, len(files))
	for i, f := range files {
		fileResponses[i] = *ToFileResponse(&f)
	}

	response.OKWithMeta(c, fileResponses, &response.Meta{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	})
}

// CreateFolder godoc
// @Summary Create folder
// @Description Create a new folder
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateFolderRequest true "Folder details"
// @Success 201 {object} response.Response{data=FileResponse}
// @Router /api/v1/files/folder [post]
func (h *Handler) CreateFolder(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	folder, err := h.service.CreateFolder(userID, &req)
	if err != nil {
		response.InternalError(c, "Failed to create folder")
		return
	}

	response.Created(c, ToFileResponse(folder))
}

// Rename godoc
// @Summary Rename file/folder
// @Description Rename a file or folder
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "File ID"
// @Param request body RenameRequest true "New name"
// @Success 200 {object} response.Response{data=FileResponse}
// @Router /api/v1/files/{id}/rename [put]
func (h *Handler) Rename(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	fileID := parseUint(c.Param("id"))
	if fileID == 0 {
		response.BadRequest(c, "INVALID_ID", "Invalid file ID")
		return
	}

	var req RenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	file, err := h.service.Rename(userID, fileID, &req)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			response.NotFound(c, "File not found")
			return
		}
		response.InternalError(c, "Failed to rename")
		return
	}

	response.OK(c, ToFileResponse(file))
}

// Move godoc
// @Summary Move file/folder
// @Description Move a file or folder to another location
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "File ID"
// @Param request body MoveRequest true "Target location"
// @Success 200 {object} response.Response{data=FileResponse}
// @Router /api/v1/files/{id}/move [put]
func (h *Handler) Move(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	fileID := parseUint(c.Param("id"))
	if fileID == 0 {
		response.BadRequest(c, "INVALID_ID", "Invalid file ID")
		return
	}

	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	file, err := h.service.Move(userID, fileID, &req)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			response.NotFound(c, "File not found")
			return
		}
		response.InternalError(c, "Failed to move")
		return
	}

	response.OK(c, ToFileResponse(file))
}

// Delete godoc
// @Summary Delete file/folder
// @Description Delete a file or folder
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param id path int true "File ID"
// @Success 204
// @Failure 404 {object} response.Response
// @Router /api/v1/files/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	fileID := parseUint(c.Param("id"))
	if fileID == 0 {
		response.BadRequest(c, "INVALID_ID", "Invalid file ID")
		return
	}

	err := h.service.Delete(c.Request.Context(), userID, fileID)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			response.NotFound(c, "File not found")
			return
		}
		response.InternalError(c, "Failed to delete")
		return
	}

	response.NoContent(c)
}

func parseUint(s string) uint {
	var id uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		id = id*10 + uint(c-'0')
	}
	return id
}
