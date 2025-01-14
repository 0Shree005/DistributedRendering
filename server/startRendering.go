package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func startRendering(fileName string) {
	fmt.Println("Rendering file: ", fileName)

	cmd := exec.Command("./scripts/startRender.sh", fileName)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %v\n", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %v\n", err)
		return
	}

	// Start blender and rendering the `fileName`
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	// Read blender's rendering stats
	go func() {
		// read rendering stats
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	go func() {
		// read rendering errors
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// wait for the feh window to close
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for command: %v\n", err)
		return
	}
	fmt.Println("Rendering process completed.")
}
