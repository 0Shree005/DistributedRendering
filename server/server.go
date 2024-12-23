package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/joho/godotenv"
)

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

func handleClient(client io.Reader) {
	for {
		buff, beof := readBuffer(client)
		if beof {
			fmt.Println("client disconnected...")
			break
		} else {
			fmt.Printf("Message received: %s\n", string(buff))
		}
	}
}

func readBuffer(reader io.Reader) ([]byte, bool) {
	b := make([]byte, 1024)
	bn, err := reader.Read(b)

	if err != nil {
		if err == io.EOF {
			return nil, true
		} else {
			panic(err)
		}
	}

	return b[:bn], false
}

func chkNilErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
