package fileTransfer

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func SendFile(fileName string, fileDir string, connection net.Conn) {
	fmt.Println("send to server")
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Failed to open file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(connection, file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File sent successfully")
}
