package main

import "time"

const (
	maxURLsPerRequest = 20
	maxURLLength      = 2048
	defaultLimit      = 20
	maxLimit          = 100
)

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	SKU       string    `json:"sku"`
	ImageURLs []string  `json:"image_urls"`
	VideoURLs []string  `json:"video_urls"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductSummary struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	SKU          string    `json:"sku"`
	ImageCount   int       `json:"image_count"`
	VideoCount   int       `json:"video_count"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateProductRequest struct {
	Name      string   `json:"name"`
	SKU       string   `json:"sku"`
	ImageURLs []string `json:"image_urls"`
	VideoURLs []string `json:"video_urls"`
}

type AddMediaRequest struct {
	ImageURLs []string `json:"image_urls"`
	VideoURLs []string `json:"video_urls"`
}

type ListResponse struct {
	Data   []ProductSummary `json:"data"`
	Total  int              `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}
