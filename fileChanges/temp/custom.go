package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

var wasReturnedToMain bool = false

// Watches a directory and notifies main when file changes occur
func articleFileNotify(notifyChan chan bool) {
	// setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// use goroutine to monitor events
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// Detect any event and notify the channel
				if event.Op&(fsnotify.Rename) != 0 {
					wasReturnedToMain = true
					fmt.Printf("The event was %s and of type %T\n", event, event)
					if wasReturnedToMain {
						fmt.Println("File operation detected:", event)
						notifyChan <- true // Send signal to main
						wasReturnedToMain = false
					}
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
			wasReturnedToMain = false
			fmt.Println("for loop ends")
		}
	}()

	// Monitor current directory
	dirPath := "./"
	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Keep this function running indefinitely
	for {
		time.Sleep(1 * time.Second)
	}
}
