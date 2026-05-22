package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequestHandler(store *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body RequestBody

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and payload are required"})
			return
		}

		if strings.TrimSpace(body.UserID) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id cannot be empty"})
			return
		}

		if !store.Accept(body.UserID) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "request accepted"})
	}
}

func StatsHandler(store *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, store.Stats())
	}
}