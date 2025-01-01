package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/joho/godotenv"
)

type Result struct {
	FileName    string `json:"file_name"`
	FileChangeB bool   `json:"file_change_b"`
}

const (
	network = "tcp"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("ERROR loading .env file")
	}

	var allowedIP string = os.Getenv("ALLOWED_IP")
	var port string = os.Getenv("SERVER_PORT")

	fmt.Println("TCP server in GO SERVER")

	fmt.Println("Listening on port:", port)
	connection, err := net.Listen(network, ":"+port)
	chkNilErr(err)
	defer connection.Close()

	for {
		fmt.Println("Waiting for a connection...")
		client, err := connection.Accept()
		chkNilErr(err)

		clientIP, _, err := net.SplitHostPort(client.RemoteAddr().String())
		chkNilErr(err)

		if clientIP != allowedIP {
			fmt.Printf("Connection from %s rejected.\n", clientIP)
			client.Close()
			continue
		}

		fmt.Println("Client connected:", client.RemoteAddr())
		go handleClient(client)
	}

}

func handleClient(client net.Conn) {

	decoder := json.NewDecoder(client)

	for {
		var result Result
		err := decoder.Decode(&result)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				return
			}
			fmt.Println("Error decoding JSON metadata")
			return
		}
		fmt.Printf("FileChanged: %s, change status: %t\n", result.FileName, result.FileChangeB)

		if result.FileChangeB {
			go ReceiveFile(result, client)
		}
	}
}

func chkNilErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
