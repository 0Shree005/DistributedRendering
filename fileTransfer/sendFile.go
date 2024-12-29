package fileTransfer

import (
	"fmt"
)

func SendFile(fileName string, fileDir string) {
	fmt.Printf("Sending file %s to server\n", fileName)
	fmt.Printf("%s is in %s", fileName, fileDir)
}
