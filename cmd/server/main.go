package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vertex-mcp-server/internal/mcp"
	"vertex-mcp-server/internal/vertexai"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Vertex AI MCP Server...")

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it (using system environment)")
	}

	// Load configuration from environment
	config := vertexai.Config{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Location:  os.Getenv("GOOGLE_CLOUD_LOCATION"),
		ModelName: os.Getenv("GEMINI_MODEL"),
	}

	// Set defaults
	if config.Location == "" {
		config.Location = "global"
	}
	if config.ModelName == "" {
		config.ModelName = "gemini-3-pro-preview"
	}

	// Validate required configuration
	if config.ProjectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT environment variable is required")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Force global location for Gemini 3 models as they are global-only
	if (config.ModelName == "gemini-3-pro-preview" || config.ModelName == "gemini-3-flash-preview") && config.Location != "global" {
		log.Printf("NOTICE: Model %s requires 'global' location. Overriding configured location '%s' to 'global'.", config.ModelName, config.Location)
		config.Location = "global"
	}

	log.Printf("Configuration:")
	log.Printf("  Project ID: %s", config.ProjectID)
	log.Printf("  Location: %s", config.Location)
	log.Printf("  Model: %s", config.ModelName)
	log.Printf("  Port: %s", port)

	// Warn if using a potentially invalid model
	knownModels := map[string]bool{
		"gemini-1.5-flash":       true,
		"gemini-1.5-pro":         true,
		"gemini-2.0-flash":       true,
		"gemini-3-pro-preview":   true,
		"gemini-3-flash-preview": true,
	}
	if !knownModels[config.ModelName] {
		log.Printf("WARNING: Model '%s' is not in the known tested list. Ensure it exists in your region.", config.ModelName)
	}

	// Create context
	ctx := context.Background()

	// Initialize Vertex AI client
	vertexClient, err := vertexai.NewClient(ctx, config)
	if err != nil {
		log.Fatalf("Failed to create Vertex AI client: %v", err)
	}

	log.Println("✓ Vertex AI client initialized")

	// Create MCP server
	mcpServer := mcp.NewServer(vertexClient)
	defer mcpServer.Close()

	log.Println("✓ MCP server created")

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      mcpServer,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("🚀 Server listening on http://localhost:%s", port)
		log.Println()
		log.Println("Available tools:")
		log.Println("  - gemini_query: Query Gemini AI")
		log.Println("  - gemini_query_with_search: Query with web search")
		log.Println("  - gemini_code_review: AI-powered code review")
		log.Println()
		log.Println("To test, send JSON-RPC requests to this endpoint")
		log.Println("Press Ctrl+C to stop the server")
		log.Println()

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\nShutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Server stopped")
}
