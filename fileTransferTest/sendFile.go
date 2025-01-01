package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func SendFile(fileResult Result, connection net.Conn) {
	const maxRetries = 5
	const retryDelay = time.Millisecond * 100

	var file *os.File
	var err error

	for i := 0; i < maxRetries; i++ {
		file, err = os.Open(fileResult.FileName)
		if err == nil {
			break
		}
		if os.IsNotExist(err) {
			time.Sleep(retryDelay)
		} else {
			log.Fatalln("Unexpected error opening file:", err)
		}
	}
	if file == nil {
		log.Fatalln("Failed to open file after retries:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(connection)
	err = encoder.Encode(fileResult)
	if err != nil {
		fmt.Println("Error encoding:", err)
	}
	fmt.Println("Sent data to the server:", fileResult)

	fmt.Println("Started copying the file!")
	bytes, err := io.Copy(connection, file)
	if err != nil {
		log.Fatalln("Error in io.Copy(): ", err)
	}

	fmt.Printf("File of size %v sent successfully\n", bToMb(bytes))
}

func bToMb(b int64) float64 {
	if b < 0 {
		return 0
	}
	return float64(b) / (1000 * 1000)
}
