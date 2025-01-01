package main

import (
	"flag"
	"fmt"
	// "io"
	// "log"
	"net"
	"os"

	// "github.com/0Shree005/DistributedRendering/fileTransfer"
	"github.com/joho/godotenv"
)

const (
	network = "tcp"
)

func main() {
	// fileName := flag.String("fileName", "", "This file is the current working file from the client")
	flag.Parse()

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
	}

}

func chkNilErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
