package auth

import (
	"errors"
	"net/http"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/response"
	"github.com/gin-gonic/gin"
)

// Handler handles authentication HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} response.Response{data=UserResponse}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request data", err.Error())
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			response.BadRequest(c, "EMAIL_EXISTS", "Email is already registered")
			return
		}
		response.InternalError(c, "Failed to create account")
		return
	}

	response.Created(c, ToUserResponse(user))
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=TokenResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request data", err.Error())
		return
	}

	tokens, err := h.service.Login(&req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Unauthorized(c, "Invalid email or password")
			return
		}
		response.InternalError(c, "Authentication failed")
		return
	}

	response.OK(c, tokens)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response{data=TokenResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request data", err.Error())
		return
	}

	tokens, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	response.OK(c, tokens)
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get authenticated user's profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=UserResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/me [get]
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalError(c, "Failed to get user")
		return
	}

	response.OK(c, ToUserResponse(user))
}

// ChangePassword godoc
// @Summary Change password
// @Description Change authenticated user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request data", err.Error())
		return
	}

	err := h.service.ChangePassword(userID, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			response.BadRequest(c, "INVALID_PASSWORD", "Current password is incorrect")
			return
		}
		response.InternalError(c, "Failed to change password")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Password changed successfully"})
}
