package model

import "time"

type UploadRequest struct {
	Image      string     `json:"image"`
	DeviceInfo DeviceInfo `json:"device_info"`
}

type DeviceInfo struct {
	Model           string `json:"model"`
	UserIP          string `json:"user_ip"`
	DeviceID        string `json:"device_id"`
	NetworkProvider string `json:"network_provider"`
	OSVersion       string `json:"os_version"`
}

type Photo struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Preview   string    `json:"image_url"` // URL to static image file
}
