package auth

import (
	"errors"

	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/infrastructure/http/middleware"
	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/infrastructure/http/response"
	"github.com/gin-gonic/gin"
)

// Handler handles auth HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} response.Response{data=UserResponse}
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			response.BadRequest(c, "USER_EXISTS", "User already exists")
			return
		}
		response.InternalError(c, "Failed to register user")
		return
	}

	response.Created(c, ToUserResponse(user))
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and get tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=TokenResponse}
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	tokens, err := h.service.Login(&req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Unauthorized(c, "Invalid credentials")
			return
		}
		response.InternalError(c, "Failed to login")
		return
	}

	response.OK(c, tokens)
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response{data=TokenResponse}
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	tokens, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "Invalid refresh token")
		return
	}

	response.OK(c, tokens)
}

// Me godoc
// @Summary Get current user
// @Description Get current authenticated user profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=UserResponse}
// @Router /api/v1/auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil || user == nil {
		response.NotFound(c, "User not found")
		return
	}

	response.OK(c, ToUserResponse(user))
}

// ChangePassword godoc
// @Summary Change password
// @Description Change current user password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change details"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}

	err := h.service.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.BadRequest(c, "INVALID_PASSWORD", "Current password is incorrect")
			return
		}
		response.InternalError(c, "Failed to change password")
		return
	}

	response.OK(c, gin.H{"message": "Password changed successfully"})
}
