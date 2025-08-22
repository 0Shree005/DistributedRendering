package handlers

import (
	"fmt"
	"github.com/0Shree005/DistributedRendering/server/scripts"
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"
)

var renderJobs sync.Map // to track status of the file rendering (filename -> status {"inProgress", "done", "failed"})

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	outFile, err := os.Create("./uploads/" + header.Filename)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	io.Copy(outFile, file)
	validFileName := regexp.MustCompile(`^[\w\-\.]+\.blend$`)

	if validFileName.MatchString(header.Filename) {
		fmt.Printf("File %s received successfully!\n", header.Filename)
		w.Write([]byte("File uploaded. Rendering started.\n"))

		// Mark as in-progress
		renderJobs.Store(header.Filename, "in-progress")

		// Run Blender async
		go func(fname string) {
			err := scripts.StartRendering(fname)
			if err != nil {
				renderJobs.Store(fname, "failed")
			} else {
				renderJobs.Store(fname, "done")
			}
		}(header.Filename)

	} else {
		fmt.Printf("Invalid file name: %s\n", header.Filename)
		w.Write([]byte("Invalid file format.\n"))
	}
}
