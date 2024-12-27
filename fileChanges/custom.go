package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
)

func articleFile(notifyChan chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&(fsnotify.Rename) != 0 {
					fileName := filepath.Base(event.Name)
					// fmt.Printf("The event was %s and of type %T\n", event, event)
					notifyChan <- fileName
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	dirPath := "./"
	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
