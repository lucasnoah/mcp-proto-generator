package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// Get port from environment or use default
	port := 3333
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	// Initialize gRPC clients
	clients, err := initClients()
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer closeClients(clients)

	// Create MCP server
	server := &MCPServer{
		clients: clients,
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	
	// MCP endpoints
	mux.HandleFunc("/tools", server.handleListTools)
	mux.HandleFunc("/rpc", server.handleRPC)
	
	// Health check endpoint  
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})



	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 MCP server starting on port %d", port)
		log.Printf("📊 Serving %d tools from %d services", 99, 7)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("✅ Server stopped")
}
