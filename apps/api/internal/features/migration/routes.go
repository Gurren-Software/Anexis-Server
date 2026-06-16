package migration

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers migration routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, authMiddleware gin.HandlerFunc) {
	migration := router.Group("/migration")
	migration.Use(authMiddleware)
	{
		migration.GET("", handler.List)
		migration.POST("", handler.Start)
		migration.GET("/:id", handler.Get)
		migration.POST("/:id/cancel", handler.Cancel)
	}
}
