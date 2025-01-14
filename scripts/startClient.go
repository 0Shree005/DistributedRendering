package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/0Shree005/DistributedRendering/client"
)

func main() {

	var dirPath string

	fmt.Print("Please enter the save file's path: ")
	fmt.Scan(&dirPath)

	if dirPath[:1] == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get home directory: %v\n", err)
		}
		dirPath = filepath.Join(home, dirPath[1:])
	}

	dirPath = os.ExpandEnv(dirPath)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s\n", dirPath)
	}
	fmt.Println("The dirPath is: ", dirPath)

	openBLender := exec.Command("blender")
	openBLender.Stdout = nil
	openBLender.Stderr = nil
	err := openBLender.Start()
	if err != nil {
		log.Fatalf("Failed to start Blender: %v\n", err)
	}
	fmt.Println("Blender started in the background")

	// call the mainFile with the user specified directory as the file change detection path
	client.Client(dirPath)
}
