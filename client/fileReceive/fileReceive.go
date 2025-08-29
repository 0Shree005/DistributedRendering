package fileReceive

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func ReceiveFile(fileName string, dirPath string, server string) {
	fmt.Println("Starting file download...")

	// --- NEW LOGIC HERE ---
	// The rendered filename is based on your shell script's output
	// e.g., "untitled_render0001.png" from "untitled.blend"
	renderFileName := fileName
	if filepath.Ext(fileName) == ".blend" {
		base := fileName[:len(fileName)-len(filepath.Ext(fileName))]
		renderFileName = fmt.Sprintf("%s_render0001.png", base)
	}
	// --- END OF NEW LOGIC ---

	// Construct the download URL using the CORRECT rendered filename.
	// This is the CRITICAL change!
	downloadURL := fmt.Sprintf("%s/download?file=%s", server, renderFileName)

	// Polling loop to wait for the file to be available for download
	// This is a simple retry mechanism.
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(downloadURL)
		if err != nil {
			log.Printf("Download failed: %v. Retrying...", err)
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			// File found and ready to download
			fullFilePath := filepath.Join(dirPath, renderFileName)
			out, err := os.Create(fullFilePath)
			if err != nil {
				log.Printf("Error creating local file: %v", err)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Printf("Error writing to file: %v", err)
				return
			}

			fmt.Printf("✅ File downloaded successfully to %s\n", fullFilePath)
			return
		} else if resp.StatusCode == http.StatusNotFound {
			log.Printf("Server not ready for download. Status: %s. Retrying...", resp.Status)
			time.Sleep(2 * time.Second)
			continue
		} else {
			log.Printf("Unexpected status code from server: %d", resp.StatusCode)
			return
		}
	}
	log.Println("Failed to download file after multiple retries.")
}
