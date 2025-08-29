package filetransfer

import (
	"bytes"
	"encoding/json" // New import for JSON handling
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/0Shree005/DistributedRendering/client/fileReceive"
	"github.com/joho/godotenv"
)

const (
	maxRetries    = 5
	retryDelay    = time.Millisecond * 100
	statusPollInt = 1 * time.Second
)

// JobStatus is a local representation of the server's status struct.
type JobStatus struct {
	Status   string `json:"status"`
	Progress int    `json:"progress"`
}

func SendFile(fileName string, dirPath string) {
	var file *os.File
	var err error

	for range maxRetries {
		fmt.Println("Final file name WITH entire dirPath", dirPath+"/"+fileName)
		file, err = os.Open(dirPath + "/" + fileName)
		if err == nil {
			break
		}
		if os.IsNotExist(err) {
			time.Sleep(retryDelay)
		} else {
			log.Fatalln("Unexpected error in opening the file:", err)
		}
	}
	if file == nil {
		log.Fatalln("Failed to open file after retries:", err)
	}
	defer file.Close()

	err = godotenv.Load()
	if err != nil {
		log.Println("Warning: failed to load .env:", err)
	}

	var server string = os.Getenv("SERVER_IP")
	if server == "" {
		log.Fatalln("SERVER_IP not set in .env")
	}

	// --- File Upload ---
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}
	writer.Close()

	resp, err := http.Post(server+"/upload", writer.FormDataContentType(), body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	uploadResp, _ := io.ReadAll(resp.Body)
	fmt.Println(string(uploadResp))

	// --- Polling status ---
	for {
		time.Sleep(statusPollInt)

		statusResp, err := http.Get(fmt.Sprintf("%s/status?file=%s", server, fileName))
		if err != nil {
			log.Println("Error getting status:", err)
			continue
		}

		statusBody, err := io.ReadAll(statusResp.Body)
		if err != nil {
			log.Println("Error reading status response body:", err)
			statusResp.Body.Close()
			continue
		}
		statusResp.Body.Close()

		var jobStatus JobStatus
		if err := json.Unmarshal(statusBody, &jobStatus); err != nil {
			log.Println("Error parsing status response:", err)
			// Fallback to old behavior if JSON parsing fails
			if bytes.Contains(statusBody, []byte("done")) {
				jobStatus.Status = "done"
			} else if bytes.Contains(statusBody, []byte("failed")) {
				jobStatus.Status = "failed"
			}
		}

		fmt.Printf("Server status for %s: %s (progress: %d%%)\n", fileName, jobStatus.Status, jobStatus.Progress)

		if jobStatus.Status == "done" {
			fmt.Println("✅ Rendering completed! Downloading the rendered image from server...")
			fileReceive.ReceiveFile(fileName, dirPath, server)
			break
		}
		if jobStatus.Status == "failed" {
			fmt.Println("❌ Rendering failed!")
			break
		}
	}
}
