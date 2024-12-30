package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	// "log"
	"net"
	// "net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	// "strings"
	"sync"

	"github.com/0Shree005/DistributedRendering/fileChangeDetection"
	"github.com/joho/godotenv"
)

const (
	network = "tcp"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	fileDir := flag.String("fileDir", "", "The main directory where the file is saved")
	flag.Parse()

	fmt.Println("FileDIR is :", *fileDir)

	err := godotenv.Load()
	if err != nil {
		panic("ERROR loading .env file")
	}

	var host string = os.Getenv("SERVER_IP")
	var port string = os.Getenv("SERVER_PORT")

	fmt.Println("host: ", host)
	fmt.Println("port: ", port)

	connection, err := net.Dial(network, host+":"+port)
	chkNilError(err)

	fmt.Println("connection: ", connection)
	fmt.Println("Connected to server")

	wg.Add(1)
	go func() {
		defer wg.Done()
		filechangedetection.GetFileChange(&wg, ctx, *fileDir, connection)
	}()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Press 'q' to stop all routines and exit.")
		for {
			input, _ := reader.ReadString('\n')
			if input == "q\n" {
				fmt.Println("Canceling context and stopping routines...")
				cancel()
				return
			}
		}
	}()
	printMemStats()
	fmt.Printf("Before wait(), ACTIVE routines are: %d\n", runtime.NumGoroutine())
	wg.Wait()
	fmt.Printf("After wait(), ACTIVE routines are: %d\n", runtime.NumGoroutine())
	printMemStats()
	fmt.Println("All routines stopped. Exiting Program")
}

func chkNilError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc: %v MB\n", bToMb(m.Alloc))
	fmt.Printf("TotalAlloc: %v MB\n", bToMb(m.TotalAlloc))
	fmt.Printf("Sys: %v MB\n", bToMb(m.Sys))
	fmt.Printf("NumGC: %v MB\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1000 / 1000
}

// func init() {
// 	go func() {
// 		log.Println(http.ListenAndServe("localhost:6060", nil))
// 	}()
// }
