package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/0Shree005/DistributedRendering/server/handlers"
	"github.com/rs/cors"
)

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

func main() {
	ip, err := getLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		return
	}

	fmt.Printf("Your server's local ip: %s\n", ip)

	fmt.Println("Waiting for files...")

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", handlers.UploadHandler)
	mux.HandleFunc("/status", handlers.StatusHandler)

	handler := cors.AllowAll().Handler(mux)
	http.ListenAndServe(":8080", handler)
}
