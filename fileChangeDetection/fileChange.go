package filechangedetection

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

func FileChange(wg *sync.WaitGroup, ctx context.Context, notifyChan chan string, dirPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Exiting fileChange")
				return
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

	<-ctx.Done()
	// select {}
}
