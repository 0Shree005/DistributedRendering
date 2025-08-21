package main

import (
	"fmt"
	"net/http"

	"github.com/0Shree005/DistributedRendering/server/handlers"
)

func main() {

	fmt.Println("Waiting for files...")
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/status", handlers.StatusHandler)
	http.ListenAndServe(":8080", nil)

}
