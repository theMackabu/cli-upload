package cli

import (
	"net/http"
	"time"
)

type Response struct {
	ID       string   `json:"id"`
	FileName string   `json:"fileName"`
	CreateAt string   `json:"createdAt"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	URL  string `json:"url"`
	Size Size   `json:"size"`
}

type Size struct {
	Raw       int64  `json:"raw"`
	Formatted string `json:"formatted"`
}

type Config struct {
	BaseURL string
	Timeout time.Duration
}

type FileUploader struct {
	config *Config
	client *http.Client
}
