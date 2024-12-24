package main

import (
	"fmt"
)

func main() {
	fmt.Println("Main function started, watching for file changes...")

	// Create a channel to receive file change notifications
	notifyChan := make(chan bool)

	// Start file watching in a goroutine
	go articleFileNotify(notifyChan)

	// Infinite loop to keep main running and handle notifications
	for {
		select {
		case <-notifyChan:
			fmt.Println("Main detected file change! Perform desired actions here.")
			// You can call another function or process the event
		}
	}
}

