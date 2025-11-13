package main

import (
	"context"
	"fmt"
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

	// Initialize MongoDB
	db, err := database.NewMongoDB(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	log.Println("✅ Connected to MongoDB successfully")

	// Initialize server
	srv := server.NewServer(cfg, db)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	printStartupInfo(cfg.Port)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func printStartupInfo(port string) {
	fmt.Println("\n" + "═══════════════════════════════════════════════════════════════")
	fmt.Println("🚀 Women's Health Backend Server")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("📡 Server URL:          http://localhost:%s\n", port)
	fmt.Printf("📚 Swagger UI:          http://localhost:%s/swagger-ui\n", port)
	fmt.Printf("📄 Swagger YAML:        http://localhost:%s/swagger.yaml\n", port)
	fmt.Printf("❤️  Health Check:        http://localhost:%s/health\n", port)
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("✨ Server is running on port %s\n", port)
	fmt.Println("🎯 Ready to accept requests!")
	fmt.Println()
}
