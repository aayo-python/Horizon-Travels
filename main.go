package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/janto-pee/Horizon-Travels.git/controllers"
	"github.com/janto-pee/Horizon-Travels.git/util"
)

func main() {
	// Load configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load configuration environment:", err)
	}

	// Validate required configuration
	if err := validateConfig(config); err != nil {
		log.Fatal("Invalid configuration:", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in goroutine
	go func() {
		if err := runGinServer(config); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down server...")
	case <-ctx.Done():
		log.Println("Server context cancelled")
	}

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := gracefulShutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server exited gracefully")
	}
}

func validateConfig(config util.Config) error {
	if config.HTTPServerAddress == "" {
		return fmt.Errorf("HTTP server address is required")
	}
	if config.DatabaseURL == "" {
		return fmt.Errorf("database URL is required")
	}
	return nil
}

func runGinServer(config util.Config) error {
	server, err := controllers.NewServer(config)
	if err != nil {
		return fmt.Errorf("could not create server: %w", err)
	}

	log.Printf("Starting server on %s", config.HTTPServerAddress)

	if err := server.Start(config.HTTPServerAddress); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not start server: %w", err)
	}

	return nil
}

func gracefulShutdown(ctx context.Context) error {
	// Close database connections
	if err := util.CloseDB(ctx); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	// Add any other cleanup operations here
	log.Println("Cleanup completed")
	return nil
}
