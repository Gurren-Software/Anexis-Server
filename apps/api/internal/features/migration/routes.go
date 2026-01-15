package migration

import (
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers migration routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	migration := router.Group("/migration")
	migration.Use(middleware.JWTAuth(jwtSecret))
	{
		migration.GET("", handler.List)
		migration.POST("", handler.Start)
		migration.GET("/:id", handler.Get)
		migration.POST("/:id/cancel", handler.Cancel)
	}
}
