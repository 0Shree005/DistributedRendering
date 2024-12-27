package main

import (
	"fmt"
	// "log"
	// "path/filepath"
	"sync"
	"time"
	// "github.com/fsnotify/fsnotify"
)

var mutex sync.Mutex

func fileChange(resultChan chan string, fileChan chan string) {
	var stack Stack

	ticker := time.NewTicker(1 * time.Second)
	senShitToticker := time.NewTicker(3 * time.Second)
	initStackTicker := time.NewTicker(10 * time.Second)
	// notifyChan := make(chan bool)

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
			case <-senShitToticker.C:
				fmt.Println("Calling shouldItGoToMain")
				shouldItGoToMain(&stack, resultChan, fileChan)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-initStackTicker.C:
				mutex.Lock()
				fmt.Println("Stack was initialized")
				initStack(&stack)
				mutex.Unlock()
			}
		}
	}()

	select {}
}

func shouldItGoToMain(stack2 *Stack, resultChan chan string, fileChan chan string) {

	// fileName := startLooking(fileChan)

	// fmt.Printf("Value of fileName is: %v", fileName)

	mutex.Lock()
	defer mutex.Unlock()

	if !stack2.IsEmpty() {
		peekedItem := stack2.Peek()
		// fmt.Printf("Peeked Element is %T\n", peekedItem)
		if peekedItem == false {
			stack2.Push(true)
			fmt.Println("This will be sent to main, YES")
			// resultChan <- fileName
		} else if peekedItem == true {
			fmt.Println("This won't be sent to main, NO")
			// resultChan <- "NO"
		} else if peekedItem == nil {
			fmt.Println("Peeked item is NIL")
			// resultChan <- "NIL"
		}
	} else if stack2.IsEmpty() {
		fmt.Println("The stack is empty")
	}
}

func initStack(stack1 *Stack) {
	stack1.Reset()
	stack1.Push(false)
	// fmt.Println("Last item: ", stack1.Peek())
}

// func startLooking(fileChan chan string) string {
// 	go articleFile(fileChan)
// 	var fileNameResult string
// 	for fileNameResult = range fileChan {
// 		fmt.Println("File name is: ", fileNameResult)
// 	}
// 	return fileNameResult
// }
