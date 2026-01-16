package backup

import (
	"errors"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles backup HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new backup handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// StartExport godoc
// @Summary Start backup export
// @Description Start exporting all user data to a backup archive
// @Tags backup
// @Produce json
// @Security BearerAuth
// @Success 201 {object} response.Response{data=BackupResponse}
// @Router /api/v1/backup/export [post]
func (h *Handler) StartExport(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	job, err := h.service.StartExport(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrActiveJobExists) {
			response.BadRequest(c, "ACTIVE_JOB_EXISTS", "An active backup job already exists")
			return
		}
		response.InternalError(c, "Failed to start backup")
		return
	}

	response.Created(c, h.service.ToBackupResponse(job))
}

// Get godoc
// @Summary Get backup job
// @Description Get details of a backup job
// @Tags backup
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Success 200 {object} response.Response{data=BackupResponse}
// @Router /api/v1/backup/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Invalid job ID")
		return
	}

	job, err := h.service.GetJob(userID, jobID)
	if err != nil {
		if errors.Is(err, ErrJobNotFound) {
			response.NotFound(c, "Backup job not found")
			return
		}
		response.InternalError(c, "Failed to get backup job")
		return
	}

	response.OK(c, h.service.ToBackupResponse(job))
}

// List godoc
// @Summary List backup jobs
// @Description List user's backup jobs
// @Tags backup
// @Produce json
// @Security BearerAuth
// @Param type query string false "Filter by type (export/import)"
// @Param status query string false "Filter by status"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} response.Response
// @Router /api/v1/backup [get]
func (h *Handler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req ListBackupsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PerPage = 20
	}

	jobs, total, err := h.service.ListJobs(userID, &req)
	if err != nil {
		response.InternalError(c, "Failed to list backups")
		return
	}

	jobResponses := make([]BackupResponse, len(jobs))
	for i, j := range jobs {
		jobResponses[i] = *h.service.ToBackupResponse(&j)
	}

	response.OKWithMeta(c, jobResponses, &response.Meta{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	})
}

// Download godoc
// @Summary Download backup
// @Description Get download URL for a completed backup
// @Tags backup
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Success 200 {object} response.Response
// @Router /api/v1/backup/{id}/download [get]
func (h *Handler) Download(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Invalid job ID")
		return
	}

	url, err := h.service.GetDownloadURL(c.Request.Context(), userID, jobID)
	if err != nil {
		if errors.Is(err, ErrJobNotFound) {
			response.NotFound(c, "Backup not found")
			return
		}
		response.BadRequest(c, "DOWNLOAD_ERROR", err.Error())
		return
	}

	response.OK(c, gin.H{"download_url": url})
}
