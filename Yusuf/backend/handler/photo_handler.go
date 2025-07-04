package handler

import (
	"backend/model"
	"backend/repository"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadPhoto menerima foto base64 dan simpan ke file
func UploadPhoto(c *gin.Context) {
	var req model.UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	base64Data := req.Image
	if strings.Contains(base64Data, ",") {
		base64Data = strings.Split(base64Data, ",")[1]
	}

	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 data"})
		return
	}

	id, err := repository.SavePhotoWithDeviceInfo(imageBytes, req)
	if err != nil {
		log.Println("SavePhoto error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"id":      id.String(),
	})
}

// GetPhotoByID mengembalikan image/png dari file berdasarkan ID
func GetPhotoByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	image, err := repository.GetPhotoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(image)
}

// ListPhotos menampilkan semua foto + URL preview
func ListPhotos(c *gin.Context) {
	photos, err := repository.GetAllPhotos()
	if err != nil {
		log.Println("ListPhotos error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve photos"})
		return
	}

	c.JSON(http.StatusOK, photos)
}

// UploadVideo menerima video WebM dari form-data dan simpan ke file
func UploadVideo(c *gin.Context) {
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video upload"})
		return
	}
	defer file.Close()

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("video_%d_%s", timestamp, header.Filename)
	savePath := filepath.Join("tmp/videos", filename)

	out, err := os.Create(savePath)
	if err != nil {
		log.Println("Video save error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video"})
		return
	}
	defer out.Close()
	io.Copy(out, file)

	publicURL := os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		publicURL = "http://localhost:8090"
	}
	videoURL := fmt.Sprintf("%s/videos/%s", publicURL, filename)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Video uploaded successfully",
		"video_url": videoURL,
	})
}

func ServePhotoFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("tmp/photos", filename)

	// Cek apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Serve file image
	c.File(filePath)
}
