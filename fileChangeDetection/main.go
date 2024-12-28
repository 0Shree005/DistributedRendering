package main

import (
	"fmt"
)

func main() {
	resultChan := make(chan string)
	go middle(resultChan)
	for fileNameRes := range resultChan {
		fmt.Printf("FROM MAIN %s file was changed\n", fileNameRes)
	}
}
