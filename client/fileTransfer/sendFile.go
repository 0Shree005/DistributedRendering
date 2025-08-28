package filetransfer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	maxRetries    = 5
	retryDelay    = time.Millisecond * 100
	statusPollInt = 5 * time.Second
)

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
		statusBody, _ := io.ReadAll(statusResp.Body)
		statusResp.Body.Close()

		status := string(statusBody)
		fmt.Println("Server status:", status)

		if contains(status, "done") {
			fmt.Println("✅ Rendering completed!")
			break
		}
		if contains(status, "failed") {
			fmt.Println("❌ Rendering failed!")
			break
		}
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
