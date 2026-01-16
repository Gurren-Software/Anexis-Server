package links

import (
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers link routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	links := router.Group("/links")
	{
		// Public access routes (no auth required)
		links.GET("/:token/access", handler.Access)
		links.GET("/:token/stream", handler.Stream)

		// Protected management routes
		protected := links.Group("")
		protected.Use(middleware.JWTAuth(jwtSecret))
		{
			protected.GET("", handler.List)
			protected.POST("", handler.Create)
			protected.PUT("/:id", handler.Update)
			protected.DELETE("/:id", handler.Delete)
		}
	}
}
