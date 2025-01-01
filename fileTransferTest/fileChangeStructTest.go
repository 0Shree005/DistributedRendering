package main

import (
	// "fmt"
	"log"
	// "path/filepath"
	"regexp"

	"github.com/fsnotify/fsnotify"
)

func fileChange(notifyChan chan Result) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	validFileName := regexp.MustCompile(`[\w\-\.]+\.blend$`)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fileNamePattern := event.Name
				if validFileName.MatchString(fileNamePattern) {
					// fmt.Println("from fileChange(): ", fileNamePattern)
					result := Result{
						// FileName:    filepath.Base(event.Name),
						FileName:    event.Name,
						FileChangeB: true,
					}
					notifyChan <- result
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	dirPath := "/home/biggiecheese/code/DistributedRendering/fileTransferTest"
	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal("Error from fileChange for dirPath", err)
	}
	select {}
}
