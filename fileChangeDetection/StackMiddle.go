package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex

func middle(resultChan chan string) {
	var stack Stack
	ticker := time.NewTicker(1 * time.Second)
	initStackIticker := time.NewTicker(20 * time.Second)
	notifyChan := make(chan string)

	fmt.Println("Stack was initialised")
	initStack(&stack) // initialising stack for the FIRST 20 seconds as well

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		for {
			go fileChange(notifyChan)
			for fileNameRes := range notifyChan {
				if fileNameRes == "test1.blend" {
					shouldItGoToMain(&stack, resultChan, fileNameRes)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Tick")
			}
		}
	}()

	go func() {
		for {
			select {
			case <-initStackIticker.C:
				fmt.Println("Stack was initialised")
				initStack(&stack)
			}
		}
	}()

	wg.Wait()
}

func initStack(stack1 *Stack) {
	stack1.Reset()
	stack1.Push(false)
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
