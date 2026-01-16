package files

import (
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers file routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	files := router.Group("/files")
	files.Use(middleware.JWTAuth(jwtSecret))
	{
		files.GET("", handler.List)
		files.POST("/upload", handler.Upload)
		files.POST("/folder", handler.CreateFolder)
		files.GET("/:id", handler.Get)
		files.GET("/:id/download", handler.Download)
		files.PUT("/:id/rename", handler.Rename)
		files.PUT("/:id/move", handler.Move)
		files.DELETE("/:id", handler.Delete)
	}
}
