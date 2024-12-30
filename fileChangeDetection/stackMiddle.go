package filechangedetection

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex

func Middle(wg *sync.WaitGroup, ctx context.Context, resultChan chan string, notifyChan chan string, dirPath string) {
	defer wg.Done()

	var stack Stack
	ticker := time.NewTicker(1 * time.Second)
	initStackIticker := time.NewTicker(20 * time.Second)

	initStack(&stack) // initialising stack for the FIRST 20 seconds as well

	wg.Add(2)
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
				shouldItGoToMain(&stack, resultChan, fileNameRes)
			case <-ctx.Done():
				fmt.Println("Exiting Middle")
				// close(resultChan)
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

func shouldItGoToMain(stack *Stack, resultChan chan string, fileName string) {
	mutex.Lock()
	defer mutex.Unlock()

	if !stack.IsEmpty() {
		peekedItem := stack.Peek()
		if peekedItem == false {
			stack.Push(true)
			resultChan <- fileName
		}
	}
}
