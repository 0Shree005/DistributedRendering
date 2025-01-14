package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func main() {
	fmt.Println("Waiting for files...")
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	outFile, err := os.Create("./uploads/" + header.Filename)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	io.Copy(outFile, file)
	validFileName := regexp.MustCompile(`^[\w\-\.]+\.blend$`)

	if validFileName.MatchString(header.Filename) {
		fmt.Printf("File %s received successfully!\n", header.Filename)
		w.Write([]byte("File successfully reached the server!\n"))
		go startRendering(header.Filename)
	} else {
		fmt.Printf("Invalid file name: %s\n", header.Filename)
		w.Write([]byte("Invalid file format.\n"))
	}
}
