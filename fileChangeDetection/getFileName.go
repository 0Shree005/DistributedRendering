package filechangedetection

import (
	"fmt"
	"github.com/0Shree005/DistributedRendering/fileTransfer"
)

func GetFileChange(fileDir string) {
	resultChan := make(chan string)
	go Middle(resultChan, fileDir)
	for fileNameRes := range resultChan {
		fmt.Printf("FROM MAIN %s file was changed\n", fileNameRes)
		fileTransfer.SendFile(fileNameRes, fileDir)
	}
}
