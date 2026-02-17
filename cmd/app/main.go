package main

import (
	"context"
	"interview-project/internal/config"
	"interview-project/internal/server"
	"interview-project/storage/sqlite"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	storage, err := sqlite.InitDatabase(cfg.DB.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	srv := server.NewServer(cfg, storage.DB())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.Start(ctx); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-sigChan
	log.Println("Interrupt signal received")

	cancel()

	<-time.After(2 * time.Second)
}
