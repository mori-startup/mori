package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"mori/pkg/utils"
)

// Define the upload path – adjust as needed.
const uploadPath = "./fileUploads"

// Define allowed file extensions
var allowedExtensions = []string{
	".txt", ".pdf", ".doc", ".docx", ".xls", ".xlsx",
	".jpg", ".jpeg", ".png", ".gif", ".csv",
}

// FileInfo holds information about an uploaded file
type FileInfo struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	UploadDate string `json:"uploadDate"`
}

// isPathSafe checks if the given path is safe and within the upload directory
func isPathSafe(path string) bool {
	// Get absolute paths
	absUploadPath, err := filepath.Abs(uploadPath)
	if err != nil {
		return false
	}
	absFilePath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Check if the file path is within the upload directory
	return strings.HasPrefix(absFilePath, absUploadPath)
}

// sanitizeFilename removes potentially dangerous characters from the filename
func sanitizeFilename(filename string) string {
	// Remove any path separators
	filename = filepath.Base(filename)

	// Remove any null bytes
	filename = strings.ReplaceAll(filename, "\x00", "")

	// Remove any potentially dangerous characters
	dangerous := []string{"..", "~", "/", "\\"}
	for _, d := range dangerous {
		filename = strings.ReplaceAll(filename, d, "")
	}

	return filename
}

// hasAllowedExtension checks if the file has an allowed extension
func hasAllowedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// UploadFiles handles file uploads.
// It expects a multipart form with one or more files under the key "files".
func (h *Handler) UploadFiles(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// Allow only POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ensure the upload directory exists with proper permissions
	absUploadPath, err := filepath.Abs(uploadPath)
	if err != nil {
		http.Error(w, "Error resolving upload path", http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(absUploadPath); os.IsNotExist(err) {
		// Create directory with restricted permissions (0755)
		if err := os.MkdirAll(absUploadPath, 0755); err != nil {
			http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
			return
		}
	}

	// Parse the multipart form (limit: 20MB).
	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Retrieve files with key "files".
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Sanitize the filename
		sanitizedFilename := sanitizeFilename(fileHeader.Filename)
		if sanitizedFilename == "" {
			http.Error(w, "Invalid filename", http.StatusBadRequest)
			return
		}

		// Check file extension
		if !hasAllowedExtension(sanitizedFilename) {
			http.Error(w, "File type not allowed", http.StatusBadRequest)
			return
		}

		// Create destination file path
		dstPath := filepath.Join(absUploadPath, sanitizedFilename)

		// Verify the path is safe
		if !isPathSafe(dstPath) {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		// Check if file already exists
		if _, err := os.Stat(dstPath); err == nil {
			http.Error(w, "File already exists", http.StatusConflict)
			return
		}

		// Create the file with restricted permissions (0644)
		dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file data to the destination.
		_, err = io.Copy(dst, file)
		if err != nil {
			// Clean up the file if copy fails
			os.Remove(dstPath)
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "File(s) uploaded successfully")
}

// ListFiles returns a JSON array of objects, each containing the file name, its size, and upload date.
func (h *Handler) ListFiles(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	// Ensure the upload directory exists and is safe
	absUploadPath, err := filepath.Abs(uploadPath)
	if err != nil {
		http.Error(w, "Error resolving upload path", http.StatusInternalServerError)
		return
	}

	// Verify the upload directory exists
	if _, err := os.Stat(absUploadPath); os.IsNotExist(err) {
		// Return empty list if directory doesn't exist
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]FileInfo{})
		return
	}

	// Read directory contents
	files, err := os.ReadDir(absUploadPath)
	if err != nil {
		http.Error(w, "Error reading upload directory", http.StatusInternalServerError)
		return
	}

	var filesInfo []FileInfo
	for _, file := range files {
		if !file.IsDir() {
			// Get file info
			info, err := file.Info()
			if err != nil {
				// Skip files with errors getting info
				continue
			}

			// Verify the file is within the upload directory
			filePath := filepath.Join(absUploadPath, file.Name())
			if !isPathSafe(filePath) {
				// Skip files that are not in the upload directory
				continue
			}

			// Sanitize the filename
			sanitizedName := sanitizeFilename(file.Name())
			if sanitizedName == "" {
				// Skip files with invalid names
				continue
			}

			filesInfo = append(filesInfo, FileInfo{
				Name:       sanitizedName,
				Size:       info.Size(),
				UploadDate: info.ModTime().Format("2006-01-02 15:04:05"),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filesInfo)
}

// DeleteFile deletes a file from the upload folder.
// It expects the URL pattern: /api/files/{filename}
func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	// Handle preflight OPTIONS request.
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Filename not specified", http.StatusBadRequest)
		return
	}

	// Sanitize the filename
	filename := sanitizeFilename(parts[3])
	if filename == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Crée le chemin absolu du fichier dans le dossier d'upload
	absUploadPath, err := filepath.Abs(uploadPath)
	if err != nil {
		http.Error(w, "Error resolving upload path", http.StatusInternalServerError)
		return
	}
	absFilePath := filepath.Join(absUploadPath, filename)

	// Vérifie que le chemin est sûr
	if !isPathSafe(absFilePath) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Si le fichier n'existe pas, retourne OK
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "File not found (already deleted)")
		return
	}

	if err := os.Remove(absFilePath); err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "File deleted successfully")
}

// ReadFileContent reads the content of a file securely.
func ReadFileContent(filename string) ([]byte, error) {
	// Sanitize the filename
	sanitizedFilename := sanitizeFilename(filename)
	if sanitizedFilename == "" {
		return nil, fmt.Errorf("invalid filename")
	}

	// Construct the absolute file path
	absUploadPath, err := filepath.Abs(uploadPath)
	if err != nil {
		return nil, fmt.Errorf("error resolving upload path: %v", err)
	}
	absFilePath := filepath.Join(absUploadPath, sanitizedFilename)

	// Verify the path is safe
	if !isPathSafe(absFilePath) {
		return nil, fmt.Errorf("invalid file path")
	}

	// Read the file content
	content, err := os.ReadFile(absFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return content, nil
}
