package main

import "github.com/gin-gonic/gin"

func main() {
	store := NewStore()

	r := gin.Default()

	r.POST("/request", RequestHandler(store))
	r.GET("/stats", StatsHandler(store))

	r.Run(":8080")
}