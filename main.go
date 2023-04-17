package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

const (
	BucketName = "your-bucket-name"
	ProjectID  = "your-project-id"
)

func main() {
	http.HandleFunc("/upload", handleUpload)
	http.ListenAndServe(":8080", nil)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	filename := header.Filename

	// Upload the file to Google Cloud Storage
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "Error initializing Google Cloud Storage client", http.StatusInternalServerError)
		return
	}

	bucket := storageClient.Bucket(BucketName)
	obj := bucket.Object(filename)
	writer := obj.NewWriter(ctx)
	writer.ObjectAttrs.ContentType = contentType
	if _, err := io.Copy(writer, file); err != nil {
		http.Error(w, "Error uploading the file", http.StatusInternalServerError)
		return
	}
	if err := writer.Close(); err != nil {
		http.Error(w, "Error closing the writer", http.StatusInternalServerError)
		return
	}

	// Generate a Signed URL for the uploaded file
	url, err := obj.SignedURL(ctx, time.Now().Add(24*time.Hour), &storage.SignedURLOptions{
		GoogleAccessID: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		PrivateKey:     []byte(os.Getenv("PRIVATE_KEY")),
		Method:         "GET",
		Expires:        time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		http.Error(w, "Error generating Signed URL", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully. Shareable URL: %s", url)
}
