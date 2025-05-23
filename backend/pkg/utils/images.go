package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Patht to default image location
const defaultImage = "imageUpload/default.svg"

// Creates new file and reads image bytes into it
// returns path to new image
// returns default avatar or provided one
func SaveAvatar(r *http.Request) string {
	// Read data from request
	file, fileHeader, errRead := r.FormFile("avatar")
	if errRead != nil {
		return defaultImage
	}
	defer file.Close()
	// get content type -> png, gif or jpeg
	contentType := fileHeader.Header["Content-Type"][0]
	// Create empty local file with correct file extension
	localFile, err := createTempFile(contentType)
	// If type not recognized return default image path
	if err != nil {
		return defaultImage
	}
	defer localFile.Close()

	// read data in new file
	fileData, err := io.ReadAll(file)
	if err != nil {
		return defaultImage
	}
	localFile.Write(fileData)
	return strings.Replace(localFile.Name(), "\\", "/", -1)
}

/* ------------------------- for posts and comments ------------------------- */
// Creates new file and reads image bytes into it
// returns path to new image
// returns no path if not exist
func SaveImage(r *http.Request) string {
	// Read data from request
	file, fileHeader, errRead := r.FormFile("image")
	if errRead != nil {
		return ""
	}
	defer file.Close()
	// get content type -> png, gif or jpeg
	contentType := fileHeader.Header["Content-Type"][0]
	// Create empty local file with correct file extension
	localFile, err := createTempFile(contentType)
	// If type not ecognized return default image path
	if err != nil {
		return ""
	}
	defer localFile.Close()

	// read data in new file
	fileData, err := io.ReadAll(file)
	if err != nil {
		return ""
	}
	localFile.Write(fileData)
	return strings.Replace(localFile.Name(), "\\", "/", -1)
}

// creates empty local file based on file type
func createTempFile(fileType string) (*os.File, error) {
	var localFile *os.File
	var err error

	switch fileType {
	case "image/jpeg":
		localFile, err = os.CreateTemp("imageUpload", "*.jpg")
	case "image/png":
		localFile, err = os.CreateTemp("imageUpload", "*.png")
	case "image/gif":
		localFile, err = os.CreateTemp("imageUpload", "*.gif")
	default:
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	return localFile, nil
}
