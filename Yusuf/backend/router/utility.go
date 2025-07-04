package router

import (
	"backend/handler"

	"github.com/gin-gonic/gin"
)

// Utility mendaftarkan semua route untuk upload video dan serve file
func Utility(router *gin.RouterGroup) {
	mysuf := router.Group("/mysuf")
	{
		// Upload video dengan metadata device
		mysuf.POST("/videos", handler.UploadVideoWithDeviceInfo)

		// Serve file video langsung via URL
		mysuf.GET("/videos/:filename", handler.ServeVideoFile)
	}
}
