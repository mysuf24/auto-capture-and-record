package handler

import (
	"backend/dto"
	"backend/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadVideoWithDeviceInfo menerima video WebM dan metadata perangkat dari multipart/form-data
func UploadVideoWithDeviceInfo(c *gin.Context) {
	// Ambil file video
	header, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video file not provided"})
		return
	}

	file, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open video file"})
		return
	}
	defer file.Close()

	// Baca byte konten video
	videoBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read video file"})
		return
	}

	// Ambil JSON string device_info dari form field
	deviceInfoStr := c.PostForm("device_info")
	if deviceInfoStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device info not provided"})
		return
	}

	var deviceInfo dto.DeviceInfoDTO
	if err := json.Unmarshal([]byte(deviceInfoStr), &deviceInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device_info format"})
		return
	}

	// Generate nama file unik
	filename := fmt.Sprintf("video_%d_%s.webm", time.Now().Unix(), uuid.New().String())

	// Simpan file + metadata ke DB
	videoURL, err := repository.SaveVideoWithDeviceInfo(videoBytes, filename, deviceInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Video uploaded successfully",
		"video_url": videoURL,
	})
}

// ServeVideoFile memberikan file video berdasarkan nama file
func ServeVideoFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("tmp/videos", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
}
