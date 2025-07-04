package main

import (
	"backend/config"
	"backend/middleware"
	"backend/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Buat folder tmp/photos dan tmp/videos jika belum ada
	dirs := []string{"tmp/photos", "tmp/videos"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
}

func ensureVideoDir() {
	path := "tmp/videos"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatal("Failed to create videos directory:", err)
		}
	}
}

func main() {
	config.InitDB()
	ensureVideoDir()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")
	router.Utility(api)

	// Pindahkan ke sini
	r.Static("/videos", "./tmp/videos")

	r.Run(":8090")
}
