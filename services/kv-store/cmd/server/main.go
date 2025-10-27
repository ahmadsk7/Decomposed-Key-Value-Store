package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"kv-store/internal"
	"kv-store/internal/store"
	kv "kv-store/internal/proto/kv/v1"
)

func main() {
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = ":8081"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	store := store.New()
	kv.RegisterKeyValueServer(s, internal.NewServer(store))

	log.Printf("kv-store listening on %s", addr)

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
	s.GracefulStop()
	log.Println("shutdown complete")
}
