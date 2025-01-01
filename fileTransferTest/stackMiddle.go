// package filechangedetection

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex

func Middle(wg *sync.WaitGroup, ctx context.Context, resultChan chan Result, notifyChan chan Result, dirPath string) {
	defer wg.Done()

	var stack Stack
	ticker := time.NewTicker(1 * time.Second)
	initStackIticker := time.NewTicker(20 * time.Second)

	initStack(&stack) // initialising stack for the FIRST 20 seconds as well

	wg.Add(1)
	go func() {
		defer wg.Done()
		fileChange(notifyChan)
	}()

	go func() {
		defer close(resultChan)
		for {
			select {
			case fileNameRes, ok := <-notifyChan:
				if !ok {
					fmt.Println("resultChan closed, exiting middle")
					return
				}
				shouldItGoToMain(&stack, resultChan, fileNameRes)
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
	// fmt.Println("Last item: ", stack1.Peek())
}

func shouldItGoToMain(stack *Stack, resultChan chan Result, fileName Result) {
	mutex.Lock()
	defer mutex.Unlock()

	if !stack.IsEmpty() {
		peekedItem := stack.Peek()
		if peekedItem == false {
			stack.Push(true)
			resultChan <- fileName
			fmt.Printf("Sent to client.go: %+v\n", fileName)
		}
	}
}
