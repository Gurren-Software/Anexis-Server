package backup

import (
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers backup routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	backup := router.Group("/backup")
	backup.Use(middleware.JWTAuth(jwtSecret))
	{
		backup.GET("", handler.List)
		backup.POST("/export", handler.StartExport)
		backup.GET("/:id", handler.Get)
		backup.GET("/:id/download", handler.Download)
	}
}
