package main

import (
	"fmt"
)

func main() {
	// stackFile()
	custom()
}

// func stackFile() {
// 	resultChan := make(chan string)
// 	fileChan := make(chan string)
// 	go fileChange(resultChan, fileChan)
//
// 	for stackResult := range resultChan {
// 		fmt.Println("StackResult:", stackResult)
// 	}
// }

func custom() {
	articleFileNotifyChan := make(chan string)
	go articleFile(articleFileNotifyChan)

	for articleResult := range articleFileNotifyChan {
		fmt.Println("Article result fileName is: ", articleResult)
	}
}
