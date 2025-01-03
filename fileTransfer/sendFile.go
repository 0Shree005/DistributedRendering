package filetransfer

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	maxRetries = 5
	retryDelay = time.Millisecond * 100
)

func SendFile(fileName string) {
	var file *os.File
	var err error

	for i := 0; i < maxRetries; i++ {
		file, err = os.Open(fileName)
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

	var server string = os.Getenv("SERVER_IP")

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

	io.Copy(os.Stdout, resp.Body)
}
