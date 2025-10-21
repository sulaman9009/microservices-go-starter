package main

import (
	"log"

	"ride-sharing/services/api-gateway/internal/transport"
)

func main() {
	srv := transport.NewHTTPServer()
	if err := srv.Start(); err != nil {
		log.Fatal("failed to start HTTP server:", err)
	}
}
