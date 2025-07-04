package router

import (
	"backend/handler"

	"github.com/gin-gonic/gin"
)

func Utility(router *gin.RouterGroup) {
	routes := router.Group("/mysuf")
	{
		// Upload routes
		routes.POST("/photos", handler.UploadPhoto)
		routes.POST("/videos", handler.UploadVideo)

		// Serve media files
		routes.GET("/photos/:filename", handler.ServePhotoFile)
	}
}
