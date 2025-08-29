package handlers

import (
	"encoding/json"
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "Missing ?file= parameter", http.StatusBadRequest)
		return
	}

	if val, ok := renderJobs.Load(file); ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(val)
	} else {
		http.Error(w, "File not found", http.StatusNotFound)
	}
}
