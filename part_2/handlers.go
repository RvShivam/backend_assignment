package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func validateURLs(urls []string) string {
	if len(urls) > maxURLsPerRequest {
		return fmt.Sprintf("max %d URLs per array", maxURLsPerRequest)
	}
	for _, u := range urls {
		if len(u) > maxURLLength {
			return "URL exceeds 2048 characters"
		}
		parsed, err := url.ParseRequestURI(u)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return fmt.Sprintf("invalid URL: %s", u)
		}
	}
	return ""
}

func CreateProductHandler(store *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		if strings.TrimSpace(req.Name) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		if strings.TrimSpace(req.SKU) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "sku is required"})
			return
		}
		if msg := validateURLs(req.ImageURLs); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if msg := validateURLs(req.VideoURLs); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		p, err := store.Create(req)
		if err == errDuplicateSKU {
			c.JSON(http.StatusConflict, gin.H{"error": "sku already exists"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.JSON(http.StatusCreated, p)
	}
}

func ListProductsHandler(store *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := defaultLimit
		offset := 0

		if l := c.Query("limit"); l != "" {
			v, err := strconv.Atoi(l)
			if err != nil || v < 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a positive integer"})
				return
			}
			if v > maxLimit {
				v = maxLimit
			}
			limit = v
		}

		if o := c.Query("offset"); o != "" {
			v, err := strconv.Atoi(o)
			if err != nil || v < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "offset must be a non-negative integer"})
				return
			}
			offset = v
		}

		data, total := store.List(limit, offset)
		c.JSON(http.StatusOK, ListResponse{
			Data:   data,
			Total:  total,
			Limit:  limit,
			Offset: offset,
		})
	}
}

func GetProductHandler(store *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := store.Get(c.Param("id"))
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusOK, p)
	}
}

func AddMediaHandler(store *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AddMediaRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		if len(req.ImageURLs) == 0 && len(req.VideoURLs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "at least one of image_urls or video_urls is required"})
			return
		}
		if msg := validateURLs(req.ImageURLs); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if msg := validateURLs(req.VideoURLs); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		p, ok := store.AddMedia(c.Param("id"), req.ImageURLs, req.VideoURLs)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusOK, p)
	}
}
