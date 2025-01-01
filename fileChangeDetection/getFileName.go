package filechangedetection

import (
	"context"
	"fmt"
	"net"
	"sync"
)

func GetFileChange(sendListener chan bool, wg *sync.WaitGroup, ctx context.Context, fileDir string, connection net.Conn) {
	defer wg.Done()

	resultChan := make(chan string)
	notifyChan := make(chan string)

	wg.Add(2)
	go func() {
		defer wg.Done()
		Middle(wg, ctx, resultChan, notifyChan, fileDir)
	}()

	for {
		select {
		case fileNameRes, ok := <-resultChan:
			if !ok {
				fmt.Println("ResultChan closed, exiting GetFileChange")
				return
			}
			fmt.Printf("FROM MAIN %s file was changed\n", fileNameRes)
			mainChan <- fileNameRes
		case <-ctx.Done():
			fmt.Println("Exiting GetFileChange()")
			close(resultChan)
			return
		}
	}
}
