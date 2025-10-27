package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"api-gateway/internal"
)

func main() {
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	kvAddr := os.Getenv("KV_GRPC_ADDR")
	if kvAddr == "" {
		kvAddr = "localhost:8081"
	}

	client, err := internal.NewClient(kvAddr)
	if err != nil {
		log.Fatalf("failed to connect to kv store: %v", err)
	}

	handlers := internal.NewHandlers(client)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /kv/{key}", handlers.Put)
	mux.HandleFunc("GET /kv/{key}", handlers.Get)
	mux.HandleFunc("DELETE /kv/{key}", handlers.Delete)

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	log.Printf("api-gateway listening on %s", httpAddr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
	srv.Shutdown(nil)
	log.Println("shutdown complete")
}
