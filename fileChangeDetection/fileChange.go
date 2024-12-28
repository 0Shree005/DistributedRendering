package main

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Result struct {
	fileName    string
	fileChangeB bool
}

func fileChange(notifyChan chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fileName := filepath.Base(event.Name)
				notifyChan <- fileName
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
