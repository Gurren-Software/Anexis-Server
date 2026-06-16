package auth

import (
	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers auth routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string, isStandalone bool) {
	auth := router.Group("/auth")
	{
		if !isStandalone {
			// Public routes - only in SaaS mode
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
			auth.POST("/refresh", handler.RefreshToken)
		}

		// Protected routes - always available but using different auth
		protected := auth.Group("")
		protected.Use(middleware.JWTAuth(jwtSecret))
		{
			protected.GET("/me", handler.Me)
			protected.PUT("/password", handler.ChangePassword)
		}
	}
}
