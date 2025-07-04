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
	// Load environment variables dari .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Buat direktori tmp/videos jika belum ada
	if err := os.MkdirAll("tmp/videos", os.ModePerm); err != nil {
		log.Fatalf("Failed to create videos directory: %v", err)
	}
}

func main() {
	// Inisialisasi koneksi database
	config.InitDB()

	// Setup router
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Register API routes
	api := r.Group("/api")
	router.Utility(api)

	// Serve video statis langsung dari /videos
	r.Static("/videos", "./tmp/videos")

	// Jalankan server di port 8090
	if err := r.Run(":8090"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
