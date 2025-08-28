package filechange

// package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0Shree005/DistributedRendering/client/fileTransfer"
)

var mutex sync.Mutex

func Middle(wg *sync.WaitGroup, ctx context.Context, dirPath string) {
	notifyChan := make(chan string)
	defer wg.Done()

	var stack Stack
	ticker := time.NewTicker(1 * time.Second)
	initStackIticker := time.NewTicker(20 * time.Second)

	initStack(&stack)

	wg.Add(1)
	go func() {
		defer wg.Done()
		FileChange(wg, ctx, notifyChan, dirPath)
	}()

	go func() {
		defer close(notifyChan)
		for {
			select {
			case fileNameRes, ok := <-notifyChan:
				if !ok {
					fmt.Println("notifyChan closed, exiting middle")
					return
				}
				shouldItGoToServer(&stack, fileNameRes, dirPath)
			case <-ctx.Done():
				fmt.Println("Exiting Middle")
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Tick")
			case <-initStackIticker.C:
				initStack(&stack)
			case <-ctx.Done():
				fmt.Println("Exiting Middle")
				return
			}
		}
	}()
}

func initStack(stack1 *Stack) {
	stack1.Reset()
	stack1.Push(false)
	fmt.Println("Stack was initialised")
}

func shouldItGoToServer(stack *Stack, fileName string, dirPath string) {
	mutex.Lock()
	defer mutex.Unlock()

	peekedItem := stack.Peek()
	if peekedItem == false {
		stack.Push(true)
		filetransfer.SendFile(fileName, dirPath)
	}
}
