package main

import (
	// gin
	"github.com/gin-gonic/gin"
)

const (
	PORT = ":8080"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	router.Run(PORT)
}
