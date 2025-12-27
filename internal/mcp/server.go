package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"vertex-mcp-server/internal/vertexai"
)

// Server is the MCP server implementation
type Server struct {
	toolHandler  *ToolHandler
	vertexClient *vertexai.Client
}

// NewServer creates a new MCP server
func NewServer(vertexClient *vertexai.Client) *Server {
	toolHandler := NewToolHandler(vertexClient)

	return &Server{
		toolHandler:  toolHandler,
		vertexClient: vertexClient,
	}
}

// ServeHTTP handles HTTP requests for the MCP server
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers for web clients
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON-RPC request
	var req struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      interface{}     `json:"id"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Received request: method=%s", req.Method)

	// Handle different MCP methods
	var result interface{}
	var err error

	ctx := r.Context()

	switch req.Method {
	case "initialize":
		result = s.handleInitialize()

	case "notifications/initialized":
		// Client notification, no response needed
		w.WriteHeader(http.StatusNoContent)
		return

	case "tools/list":
		result = s.handleListTools()

	case "tools/call":
		result, err = s.handleCallTool(ctx, req.Params)

	default:
		s.sendError(w, req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
		return
	}

	// Build JSON-RPC response
	if err != nil {
		s.sendError(w, req.ID, -32603, err.Error())
		return
	}

	s.sendSuccess(w, req.ID, result)
}

func (s *Server) handleInitialize() map[string]interface{} {
	return map[string]interface{}{
		"protocolVersion": "2025-06-18",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "vertex-mcp-server",
			"version": "1.0.0",
		},
	}
}

func (s *Server) handleListTools() map[string]interface{} {
	tools := s.toolHandler.GetTools()
	return map[string]interface{}{
		"tools": tools,
	}
}

func (s *Server) handleCallTool(ctx context.Context, paramsRaw json.RawMessage) (interface{}, error) {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(paramsRaw, &params); err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	return s.toolHandler.ExecuteTool(ctx, params.Name, params.Arguments)
}

func (s *Server) sendSuccess(w http.ResponseWriter, id interface{}, result interface{}) {
	resp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (s *Server) sendError(w http.ResponseWriter, id interface{}, code int, message string) {
	resp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}

// Close closes the server and cleans up resources
func (s *Server) Close() error {
	if s.vertexClient != nil {
		return s.vertexClient.Close()
	}
	return nil
}
