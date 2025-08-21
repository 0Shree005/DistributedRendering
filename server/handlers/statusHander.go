package handlers

import (
	"fmt"
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "Missing ?file= parameter", http.StatusBadRequest)
		return
	}

	if val, ok := renderJobs.Load(file); ok {
		w.Write([]byte(fmt.Sprintf("Status for %s: %s\n", file, val)))
	} else {
		http.Error(w, "File not found", http.StatusNotFound)
	}
}
