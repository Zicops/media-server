package main

import (
	//"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Set up a route to serve the media files
	r.GET("/media/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join("media", filename)

		c.File(filePath)
	})

	// Start the web server on port 8080
	r.Run(":8080")
}
