package filechange

// package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/fsnotify/fsnotify"
)

func FileChange(wg *sync.WaitGroup, ctx context.Context, fileNameChan chan string, dirPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	validFileName := regexp.MustCompile(`^[\w\-\.]+\.blend$`) // WITHOUT the ./ as the prefix

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Exiting fileChange")
				return
			case event := <-watcher.Events:
				if event.Op&(fsnotify.Rename) != 0 {
					fileName := filepath.Base(event.Name)
					if validFileName.MatchString(fileName) {
						fileNameChan <- fileName
					}
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}
