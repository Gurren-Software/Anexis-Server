package files

import (
	"errors"
	"io"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/features/auth"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles file HTTP requests
type Handler struct {
	service  *Service
	authRepo *auth.Repository
}

// NewHandler creates a new files handler
func NewHandler(service *Service, authRepo *auth.Repository) *Handler {
	return &Handler{service: service, authRepo: authRepo}
}

// Upload godoc
// @Summary Upload file
// @Description Upload a new file
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param compress formData bool false "Compress file"
// @Param parent_id formData string false "Parent folder ID"
// @Param description formData string false "File description"
// @Success 201 {object} response.Response{data=FileResponse}
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

	var req UploadRequest
	_ = c.ShouldBind(&req)

	// Manually parse parent_id since Gin doesn't bind *uuid.UUID from form data
	if parentIDStr := c.PostForm("parent_id"); parentIDStr != "" {
		parentID, err := uuid.Parse(parentIDStr)
		if err != nil {
			response.BadRequest(c, "INVALID_PARENT_ID", "Invalid parent folder ID")
			return
		}
		req.ParentID = &parentID
	}

	// Get user storage quota
	user, err := h.authRepo.FindByID(userID)
	if err != nil || user == nil {
		response.InternalError(c, "Failed to get user")
		return
	}

	fileRecord, err := h.service.Upload(c.Request.Context(), userID, file, req.Compress, req.ParentID, req.Description, user.StorageQuota)
	if err != nil {
		if errors.Is(err, ErrStorageExceeded) {
			response.BadRequest(c, "STORAGE_EXCEEDED", "Storage quota exceeded")
			return
		}
		response.InternalError(c, "Failed to upload file")
		return
	}

	response.Created(c, h.service.ToFileResponse(fileRecord))
}

// Download godoc
// @Summary Download file
// @Description Download a file
// @Tags files
// @Produce octet-stream
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {file} binary
// @Router /api/v1/files/{id}/download [get]
func (h *Handler) Download(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
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
// @Description Get file metadata
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {object} response.Response{data=FileResponse}
// @Router /api/v1/files/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
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

	response.OK(c, h.service.ToFileResponse(file))
}

// List godoc
// @Summary List files
// @Description List files with pagination
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param parent_id query string false "Parent folder ID"
// @Param search query string false "Search term"
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
		req.Page = 1
		req.PerPage = 20
	}

	// Manually parse parent_id since Gin doesn't bind *uuid.UUID from query params
	if parentIDStr := c.Query("parent_id"); parentIDStr != "" {
		parentID, err := uuid.Parse(parentIDStr)
		if err != nil {
			response.BadRequest(c, "INVALID_PARENT_ID", "Invalid parent folder ID")
			return
		}
		req.ParentID = &parentID
	}

	files, total, err := h.service.ListFiles(userID, &req)
	if err != nil {
		response.InternalError(c, "Failed to list files")
		return
	}

	fileResponses := make([]FileResponse, len(files))
	for i, f := range files {
		fileResponses[i] = *h.service.ToFileResponse(&f)
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

	response.Created(c, h.service.ToFileResponse(folder))
}

// Rename godoc
// @Summary Rename file/folder
// @Description Rename a file or folder
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Param request body RenameRequest true "New name"
// @Success 200 {object} response.Response{data=FileResponse}
// @Router /api/v1/files/{id}/rename [put]
func (h *Handler) Rename(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Invalid file ID")
		return
	}

	var req RenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	file, err := h.service.Rename(userID, fileID, req.Name)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			response.NotFound(c, "File not found")
			return
		}
		response.InternalError(c, "Failed to rename file")
		return
	}

	response.OK(c, h.service.ToFileResponse(file))
}

// Move godoc
// @Summary Move file/folder
// @Description Move a file or folder
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Param request body MoveRequest true "New parent"
// @Success 200 {object} response.Response{data=FileResponse}
// @Router /api/v1/files/{id}/move [put]
func (h *Handler) Move(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Invalid file ID")
		return
	}

	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	file, err := h.service.Move(userID, fileID, req.ParentID)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			response.NotFound(c, "File not found")
			return
		}
		response.InternalError(c, "Failed to move file")
		return
	}

	response.OK(c, h.service.ToFileResponse(file))
}

// Delete godoc
// @Summary Delete file/folder
// @Description Delete a file or folder
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 204
// @Router /api/v1/files/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Invalid file ID")
		return
	}

	err = h.service.Delete(c.Request.Context(), userID, fileID)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			response.NotFound(c, "File not found")
			return
		}
		response.InternalError(c, "Failed to delete file")
		return
	}

	response.NoContent(c)
}
