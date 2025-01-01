package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

const (
	network = "tcp"
)

type Result struct {
	FileName    string `json:"file_name"`
	FileChangeB bool   `json:"file_change_b"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("ERROR loading .env file")
	}

	var host string = os.Getenv("SERVER_IP")
	var port string = os.Getenv("SERVER_PORT")

	connection, err := net.Dial(network, host+":"+port)
	chkNilError(err)
	defer connection.Close()

	resultChan := make(chan Result)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go Middle(wg, ctx, resultChan, make(chan Result), "./")

	for result := range resultChan {
		// go SendFile(result, connection)
		go SendFile(result, connection)
	}
	wg.Wait()
}

func chkNilError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
