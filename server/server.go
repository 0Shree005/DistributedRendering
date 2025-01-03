package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	fmt.Printf("File %v received successfully!", outFile)
	w.Write([]byte("File uploaded successfully!"))
}
