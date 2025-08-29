package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/0Shree005/DistributedRendering/server/scripts"
	"github.com/0Shree005/DistributedRendering/server/types"
)

var renderJobs sync.Map // to track status of the file rendering (filename -> *types.JobStatus)

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

		// Mark as in-progress with 0% progress
		renderJobs.Store(header.Filename, &types.JobStatus{Status: "in-progress", Progress: 0})

		// Run Blender async
		go func(fname string) {
			err := scripts.StartRendering(fname, &renderJobs)
			if err != nil {
				renderJobs.Store(fname, &types.JobStatus{Status: "failed", Progress: 0})
			} else {
				renderJobs.Store(fname, &types.JobStatus{Status: "done", Progress: 100})
			}
		}(header.Filename)

	} else {
		fmt.Printf("Invalid file name: %s\n", header.Filename)
		w.Write([]byte("Invalid file format.\n"))
	}
}
