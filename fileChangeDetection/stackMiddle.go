package filechangedetection

import (
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex

func Middle(resultChan chan string, dirPath string) {
	var stack Stack
	ticker := time.NewTicker(1 * time.Second)
	initStackIticker := time.NewTicker(20 * time.Second)
	notifyChan := make(chan string)

	initStack(&stack) // initialising stack for the FIRST 20 seconds as well

	go func() {
		for {
			go FileChange(notifyChan, dirPath)
			for fileNameRes := range notifyChan {
				shouldItGoToMain(&stack, resultChan, fileNameRes)
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
