package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/0Shree005/DistributedRendering/fileChangeDetection"
	"github.com/joho/godotenv"
)

const (
	network = "tcp"
)

func main() {
	fileDir := flag.String("fileDir", "", "The main directory where the file is saved")
	flag.Parse()

	go func() {
		filechangedetection.GetFileChange(*fileDir)
	}()

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

	fmt.Println("Connected to server")
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
