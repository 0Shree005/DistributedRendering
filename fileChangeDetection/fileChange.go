package filechangedetection

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func FileChange(notifyChan chan string, dirPath string) {
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
					notifyChan <- fileName
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

	select {}
}
