// Package provide the main REST API routes

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting Server")

	router := gin.Default()

	router.Static("/test", "./public")

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	router.POST("/form", func(c *gin.Context) {
		data := c.PostForm("data")

		c.String(http.StatusOK, fmt.Sprintf("Processed: %s.", data))

	})

	router.Run()
}
