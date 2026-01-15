package migration

import (
	"errors"
	"strconv"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/response"
	"github.com/gin-gonic/gin"
)

// Handler handles migration HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new migration handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Start godoc
// @Summary Start migration
// @Description Start a new migration job from external provider
// @Tags migration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body StartMigrationRequest true "Migration details"
// @Success 201 {object} response.Response{data=MigrationResponse}
// @Router /api/v1/migration [post]
func (h *Handler) Start(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req StartMigrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	job, err := h.service.Start(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, ErrActiveJobExists) {
			response.BadRequest(c, "ACTIVE_JOB_EXISTS", "An active migration job already exists")
			return
		}
		response.InternalError(c, "Failed to start migration")
		return
	}

	response.Created(c, ToMigrationResponse(job))
}

// Get godoc
// @Summary Get migration job
// @Description Get details of a migration job
// @Tags migration
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job ID"
// @Success 200 {object} response.Response{data=MigrationResponse}
// @Router /api/v1/migration/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	jobID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if jobID == 0 {
		response.BadRequest(c, "INVALID_ID", "Invalid job ID")
		return
	}

	job, err := h.service.GetJob(userID, uint(jobID))
	if err != nil {
		if errors.Is(err, ErrJobNotFound) {
			response.NotFound(c, "Migration job not found")
			return
		}
		response.InternalError(c, "Failed to get migration job")
		return
	}

	response.OK(c, ToMigrationResponse(job))
}

// List godoc
// @Summary List migration jobs
// @Description List user's migration jobs
// @Tags migration
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} response.Response
// @Router /api/v1/migration [get]
func (h *Handler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req ListMigrationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PerPage = 20
	}

	jobs, total, err := h.service.ListJobs(userID, &req)
	if err != nil {
		response.InternalError(c, "Failed to list migrations")
		return
	}

	jobResponses := make([]MigrationResponse, len(jobs))
	for i, j := range jobs {
		jobResponses[i] = *ToMigrationResponse(&j)
	}

	response.OKWithMeta(c, jobResponses, &response.Meta{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	})
}

// Cancel godoc
// @Summary Cancel migration
// @Description Cancel a running migration job
// @Tags migration
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job ID"
// @Success 200 {object} response.Response
// @Router /api/v1/migration/{id}/cancel [post]
func (h *Handler) Cancel(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	jobID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if jobID == 0 {
		response.BadRequest(c, "INVALID_ID", "Invalid job ID")
		return
	}

	err := h.service.Cancel(userID, uint(jobID))
	if err != nil {
		if errors.Is(err, ErrJobNotFound) {
			response.NotFound(c, "Migration job not found")
			return
		}
		response.InternalError(c, "Failed to cancel migration")
		return
	}

	response.OK(c, gin.H{"message": "Migration cancelled"})
}
