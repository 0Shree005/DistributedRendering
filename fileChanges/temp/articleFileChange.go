package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const debounceDuration = 500 * time.Millisecond

func articleFileNotify(notifyChan chan string, dirPath string) {
	// Setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Use goroutine to monitor events
	var mu sync.Mutex
	var lastEventTime time.Time
	var lastFileName string

	// var returnedToMain bool = false
	// var returnedToMainCount int = 0

	go func() {
		for {
			// fmt.Println("for loop startssssssssssssssssssssssssssssssssssssssssssss")
			select {
			case event := <-watcher.Events:
				// Detect relevant events
				if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Rename) != 0 {
					fileName := filepath.Base(event.Name)
					mu.Lock()
					elapsedTime := time.Now()
					if elapsedTime.Sub(lastEventTime) > debounceDuration || fileName != lastFileName {
						fmt.Println("event detected:", event)
						lastEventTime = elapsedTime
						lastFileName = fileName
						mu.Unlock()

						fmt.Println("File operation detected:", event)
						// fmt.Println("returnedToMainCount:", returnedToMainCount)
						time.AfterFunc(debounceDuration, func() {
							notifyChan <- fileName // Send the file name to main
							fmt.Println("articleFileChange just access main func")
						})
					} else {
						mu.Unlock()
					}
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
			// returnedToMainCount = 0
			// fmt.Println("for loop endsssssssssssssssssssssssssssssssssssssssssssss")
		}
	}()

	// Monitor current directory
	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Keep this function running indefinitely
	select {}
}

func shouldReturnToMain(returnedToMainCount int) bool {
	if returnedToMainCount > 1 {
		return false
	}
	return true
}

// package main
//
// import (
//
//	"fmt"
//	"log"
//
//	"github.com/fsnotify/fsnotify"
//
// )
//
//	func articleFileNotify() {
//		// setup watcher
//		watcher, err := fsnotify.NewWatcher()
//		if err != nil {
//			log.Fatal(err)
//		}
//		defer watcher.Close()
//
//		done := make(chan bool)
//		// use goroutine to start the watcher
//		go func() {
//			for {
//				select {
//				// provide the list of events to monitor
//				case event := <-watcher.Events:
//					if event.Op&fsnotify.Create == fsnotify.Create {
//						fmt.Println("File created:", event.Name)
//					}
//					if event.Op&fsnotify.Write == fsnotify.Write {
//						fmt.Println("File modified:", event.Name)
//					}
//					if event.Op&fsnotify.Remove == fsnotify.Remove {
//						fmt.Println("File removed:", event.Name)
//					}
//					if event.Op&fsnotify.Rename == fsnotify.Rename {
//						fmt.Println("File renamed:", event.Name)
//					}
//					if event.Op&fsnotify.Chmod == fsnotify.Chmod {
//						fmt.Println("File permissions modified:", event.Name)
//					}
//				case err := <-watcher.Errors:
//					log.Println("Error:", err)
//				}
//			}
//		}()
//
//		// provide the directory to monitor
//		dirPath := "./"
//		err = watcher.Add(dirPath)
//		if err != nil {
//			log.Fatal(err)
//		}
//		<-done
//	}
