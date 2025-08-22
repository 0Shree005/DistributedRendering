package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/0Shree005/DistributedRendering/client/fileChange"
)

func Client(dirPath string) {
	fmt.Println("Main was called")

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)

	go filechange.Middle(&wg, ctx, dirPath)

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
