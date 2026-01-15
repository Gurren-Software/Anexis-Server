package auth

import (
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers auth routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	auth := router.Group("/auth")
	{
		// Public routes
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)

		// Protected routes
		protected := auth.Group("")
		protected.Use(middleware.JWTAuth(jwtSecret))
		{
			protected.GET("/me", handler.GetCurrentUser)
			protected.PUT("/password", handler.ChangePassword)
		}
	}
}
