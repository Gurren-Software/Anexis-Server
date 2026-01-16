package links

import (
	"errors"
	"io"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles link HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new links handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create access link
// @Description Create a new access link for a file
// @Tags links
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateLinkRequest true "Link details"
// @Success 201 {object} response.Response{data=LinkResponse}
// @Router /api/v1/links [post]
func (h *Handler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	link, err := h.service.Create(userID, &req)
	if err != nil {
		response.InternalError(c, "Failed to create link")
		return
	}

	response.Created(c, h.service.ToLinkResponse(link))
}

// List godoc
// @Summary List links
// @Description List user's access links
// @Tags links
// @Produce json
// @Security BearerAuth
// @Param file_id query string false "Filter by file ID"
// @Param type query string false "Filter by link type"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} response.Response
// @Router /api/v1/links [get]
func (h *Handler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req ListLinksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PerPage = 20
	}

	// Manually parse file_id since Gin doesn't bind *uuid.UUID from query params
	if fileIDStr := c.Query("file_id"); fileIDStr != "" {
		fileID, err := uuid.Parse(fileIDStr)
		if err != nil {
			response.BadRequest(c, "INVALID_FILE_ID", "Invalid file ID")
			return
		}
		req.FileID = &fileID
	}

	links, total, err := h.service.List(userID, &req)
	if err != nil {
		response.InternalError(c, "Failed to list links")
		return
	}

	linkResponses := make([]LinkResponse, len(links))
	for i, l := range links {
		linkResponses[i] = *h.service.ToLinkResponse(&l)
	}

	response.OKWithMeta(c, linkResponses, &response.Meta{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	})
}

// Update godoc
// @Summary Update link
// @Description Update an access link
// @Tags links
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Link ID"
// @Param request body UpdateLinkRequest true "Updated link details"
// @Success 200 {object} response.Response{data=LinkResponse}
// @Router /api/v1/links/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	linkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Invalid link ID")
		return
	}

	var req UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	link, err := h.service.Update(userID, linkID, &req)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			response.NotFound(c, "Link not found")
			return
		}
		response.InternalError(c, "Failed to update link")
		return
	}

	response.OK(c, h.service.ToLinkResponse(link))
}

// Delete godoc
// @Summary Delete link
// @Description Delete an access link
// @Tags links
// @Produce json
// @Security BearerAuth
// @Param id path string true "Link ID"
// @Success 204
// @Router /api/v1/links/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	linkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Invalid link ID")
		return
	}

	err = h.service.Delete(userID, linkID)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			response.NotFound(c, "Link not found")
			return
		}
		response.InternalError(c, "Failed to delete link")
		return
	}

	response.NoContent(c)
}

// Access godoc
// @Summary Access file via link
// @Description Download a file using an access link
// @Tags links
// @Produce octet-stream
// @Param token path string true "Link token"
// @Param password query string false "Password (if protected)"
// @Success 200 {file} binary
// @Router /api/v1/links/{token}/access [get]
func (h *Handler) Access(c *gin.Context) {
	token := c.Param("token")
	password := c.Query("password")

	reader, link, err := h.service.AccessFile(c.Request.Context(), token, password)
	if err != nil {
		switch {
		case errors.Is(err, ErrLinkNotFound):
			response.NotFound(c, "Link not found")
		case errors.Is(err, ErrLinkExpired):
			response.BadRequest(c, "LINK_EXPIRED", "This link has expired")
		case errors.Is(err, ErrDownloadLimit):
			response.BadRequest(c, "DOWNLOAD_LIMIT", "Download limit reached")
		case errors.Is(err, ErrAccessDenied):
			response.Forbidden(c, "Password required")
		case errors.Is(err, ErrInvalidPassword):
			response.Forbidden(c, "Invalid password")
		default:
			response.InternalError(c, "Failed to access file")
		}
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", "attachment; filename=\""+link.File.OriginalName+"\"")
	c.Header("Content-Type", link.File.MimeType)

	if _, err := io.Copy(c.Writer, reader); err != nil {
		response.InternalError(c, "Failed to stream file")
		return
	}
}

// Stream godoc
// @Summary Get streaming URL
// @Description Get a streaming URL for a file
// @Tags links
// @Produce json
// @Param token path string true "Link token"
// @Param password query string false "Password (if protected)"
// @Success 200 {object} response.Response
// @Router /api/v1/links/{token}/stream [get]
func (h *Handler) Stream(c *gin.Context) {
	token := c.Param("token")
	password := c.Query("password")

	url, _, err := h.service.GetStreamURL(c.Request.Context(), token, password)
	if err != nil {
		switch {
		case errors.Is(err, ErrLinkNotFound):
			response.NotFound(c, "Link not found")
		case errors.Is(err, ErrLinkExpired):
			response.BadRequest(c, "LINK_EXPIRED", "This link has expired")
		case errors.Is(err, ErrAccessDenied):
			response.Forbidden(c, "Password required")
		default:
			response.InternalError(c, "Failed to get stream URL")
		}
		return
	}

	response.OK(c, gin.H{"stream_url": url})
}
