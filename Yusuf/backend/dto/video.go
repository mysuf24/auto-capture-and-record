package dto

type DeviceInfoDTO struct {
	DeviceModel string  `json:"device_model"`
	IPAddress   string  `json:"ip_address"`
	UserAgent   string  `json:"user_agent"`
	Network     string  `json:"network"`
	Platform    string  `json:"platform"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
