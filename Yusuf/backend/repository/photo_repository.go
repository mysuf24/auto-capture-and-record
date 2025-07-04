package repository

import (
	"backend/config"
	"backend/model"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func SavePhotoWithDeviceInfo(imageBytes []byte, req model.UploadRequest) (uuid.UUID, error) {
	id := uuid.New()
	createdAt := time.Now()
	filename := fmt.Sprintf("photo_%d_%s.png", createdAt.Unix(), id.String())
	savePath := filepath.Join("tmp/photos", filename)

	// Simpan ke file
	if err := os.WriteFile(savePath, imageBytes, 0644); err != nil {
		return uuid.Nil, err
	}

	// Simpan metadata ke database
	_, err := config.DB.Exec(`
		INSERT INTO photos (id, file_path, created_at, model, user_ip, device_id, network_provider, os_version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, id, filename, createdAt, req.DeviceInfo.Model, req.DeviceInfo.UserIP, req.DeviceInfo.DeviceID, req.DeviceInfo.NetworkProvider, req.DeviceInfo.OSVersion)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func GetPhotoByID(id uuid.UUID) ([]byte, error) {
	var filePath string
	err := config.DB.QueryRow(`SELECT file_path FROM photos WHERE id = $1`, id).Scan(&filePath)
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join("tmp/photos", filePath)
	return os.ReadFile(fullPath)
}

func GetAllPhotos() ([]model.Photo, error) {
	rows, err := config.DB.Query(`SELECT id, file_path, created_at FROM photos ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	publicURL := os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		publicURL = "http://localhost:8090"
	}

	var photos []model.Photo
	for rows.Next() {
		var id uuid.UUID
		var filePath string
		var created time.Time

		if err := rows.Scan(&id, &filePath, &created); err != nil {
			return nil, err
		}

		imageURL := fmt.Sprintf("%s/photos/%s", publicURL, filePath)
		photos = append(photos, model.Photo{
			ID:        id.String(),
			CreatedAt: created,
			Preview:   imageURL,
		})
	}

	return photos, nil
}
