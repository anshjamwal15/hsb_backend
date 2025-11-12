package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/config"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/database"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize MongoDB connection
	db, err := database.NewMongoDB(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	// Initialize HTTP server
	srv := server.NewServer(cfg, db)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
