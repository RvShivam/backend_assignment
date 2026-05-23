package main

import "github.com/gin-gonic/gin"

func main() {
	store := NewStore()

	r := gin.Default()

	r.POST("/products", CreateProductHandler(store))
	r.GET("/products", ListProductsHandler(store))
	r.GET("/products/:id", GetProductHandler(store))
	r.POST("/products/:id/media", AddMediaHandler(store))

	r.Run(":8081")
}
