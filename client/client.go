package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

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

	var host string = os.Getenv("SERVER_IP")
	var port string = os.Getenv("SERVER_PORT")

	fmt.Println(host)
	fmt.Println(port)

	connection, err := net.Dial(network, host+":"+port)
	chkNilError(err)

	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		chkNilError(err)
		input = strings.TrimSpace(input)

		connection.Write([]byte(input))
	}
}

func chkNilError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
