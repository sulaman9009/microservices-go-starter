package transport

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ride-sharing/shared/env"
	"syscall"
	"time"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func NewHTTPServer() *server {
	return &server{
		mux: http.NewServeMux(),
	}
}

type server struct {
	mux *http.ServeMux
}

func (s *server) Start() error {
	srv := &http.Server{
		Addr:    httpAddr,
		Handler: s.mux,
	}

	log.Printf("Starting HTTP server on %s\n", httpAddr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server:", err)
		}
	}()

	return s.waitForShutdown(srv)
}

func (s *server) waitForShutdown(server *http.Server) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	// block here
	<-sig
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	return nil
}
