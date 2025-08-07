package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/upload", uploadPDFHandler)
	router.POST("/chat", chatHandler)

	router.Run(":8000")
}

func uploadPDFHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Upload endpoint works!"})
}

func chatHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Chat endpoint works!"})
}
