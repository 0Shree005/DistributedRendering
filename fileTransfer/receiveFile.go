package fileTransfer

import (
	"io"
	"log"
	"net"
	"os"
)

func ReceiveFile(fileName string, connection net.Conn) {

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Failed to create a file: %v", err)
		return
	}
	defer file.Close()

	log.Printf("Receiving file: %s\n", fileName)

	_, err = io.Copy(file, connection)
	if err != nil {
		log.Printf("Failed to receive file: %v", err)
		return
	}
	log.Println("File received succesfully")
}
