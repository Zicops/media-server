package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

func uploadFileToGCS(bucketName, objectName string, r io.Reader) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

	defer client.Close()

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	w := object.NewWriter(ctx)
	if _, err = io.Copy(w, r); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}
	return attrs.MediaLink, nil
}

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Error uploading file: %v", err)
			return
		}

		f, err := file.Open()
		if err != nil {
			c.String(http.StatusBadRequest, "Error opening file: %v", err)
			return
		}
		defer f.Close()

		mediaLink, err := uploadFileToGCS("zicops-vc", file.Filename, f)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error uploading to GCS: %v", err)
			return
		}

		c.String(http.StatusOK, "File uploaded to GCS: %s", mediaLink)
	})

	r.Run(":8080")
}

