package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func ReceiveFile(fileResult Result, connection net.Conn) {
	fmt.Println("Receiving file...")

	file, err := os.Create(fileResult.FileName)
	if err != nil {
		log.Printf("Failed to create a file: %v", err)
		return
	}
	defer file.Close()

	log.Printf("Receiving file: %s\n", fileResult.FileName)

	bytes, err := io.Copy(file, connection)
	if err != nil {
		log.Printf("Failed to receive file: %v", err)
		return
	}
	log.Printf("File of size %.2f MB received succesfully", bToMb(bytes))
}

func bToMb(x int64) float64 {
	if x < 0 {
		return 0
	}
	return float64(x) / (1000 * 1000)
}
