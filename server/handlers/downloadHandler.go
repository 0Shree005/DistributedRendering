package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
		return
	}

	// Construct the full path to the rendered file.
	// It's important to sanitize the input to prevent directory traversal attacks.
	safeFileName := filepath.Base(fileName)
	filePath := filepath.Join("./renders", safeFileName)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve the file to the client.
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", safeFileName))
	w.Header().Set("Content-Type", "image/png") // Assuming PNG format based on your script
	http.ServeFile(w, r, filePath)

	// Optional: Delete the file after it has been served successfully.
	// This helps manage server disk space.
	go func() {
		fmt.Printf("Attempting to delete file: %s\n", filePath)
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", filePath, err)
		} else {
			fmt.Printf("Successfully deleted file: %s\n", filePath)
		}
		// You might also want to remove the entry from renderJobs sync.Map
		// once the file has been successfully downloaded.
		renderJobs.Delete(fileName)
	}()
}
