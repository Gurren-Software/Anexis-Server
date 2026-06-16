package backup

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers backup routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, authMiddleware gin.HandlerFunc) {
	backup := router.Group("/backup")
	backup.Use(authMiddleware)
	{
		backup.GET("", handler.List)
		backup.POST("/export", handler.StartExport)
		backup.GET("/:id", handler.Get)
		backup.GET("/:id/download", handler.Download)
	}
}
