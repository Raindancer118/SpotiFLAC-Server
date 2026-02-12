package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"spotiflac/backend/config"
	"spotiflac/server"
)

// main is the entry point for the HTTP server mode
// Following rule #4: Use existing configuration and documentation
func main() {
	// Load configuration from config.yml
	// Following rule #7: Don't hardcode values, use config
	configPath := "config.yml"
	if envPath := os.Getenv("SPOTIFLAC_CONFIG"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully")
	log.Printf("Server will listen on %s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Download path: %s", cfg.Download.Path)

	// Create and configure server
	srv := server.NewServer(cfg)

	// Handle graceful shutdown
	// Following rule #15: Fail Securely - proper cleanup on shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down server...")
		if err := srv.Stop(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		os.Exit(0)
	}()

	// Start server
	// Following rule #10: Run with least privilege (non-root user)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
