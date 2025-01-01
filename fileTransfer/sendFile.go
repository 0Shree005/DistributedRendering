package fileTransfer

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func SendFile(fileName string, fileDir string, connection net.Conn) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Printf("Sending file: %s\n", fileName)

	_, err = io.Copy(connection, file)
	if err != nil {
		log.Printf("Error copying the file to connection: %v\n", err)
		return
	}
	fmt.Println("File sent successfully")
}
