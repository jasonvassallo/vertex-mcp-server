# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go-based MCP (Model Context Protocol) server exposing Google Cloud Vertex AI Gemini models to AI assistants. Implements JSON-RPC 2.0 over HTTP transport.

## Build & Run Commands

```bash
# Build
go build -o bin/gemini-mcp ./cmd/server

# Run (env loaded automatically via godotenv)
./bin/gemini-mcp

# Development run
go run ./cmd/server

# Test the running server
./test-server.sh

# Install to PATH (optional)
go build -o ~/.local/bin/gemini-mcp ./cmd/server
```

## Required Environment Variables

- `GOOGLE_CLOUD_PROJECT` - GCP project ID (required)
- `GOOGLE_APPLICATION_CREDENTIALS` - Path to service account key JSON (required)
- `GOOGLE_CLOUD_LOCATION` - GCP region (default: `global`)
- `GEMINI_MODEL` - Model name (default: `gemini-3-pro-preview`)
- `PORT` - Server port (default: `8080`)

## Architecture

```
cmd/server/main.go      → Entry point, config loading, HTTP server setup
internal/
├── mcp/
│   ├── server.go       → HTTP handler, JSON-RPC routing (initialize, tools/list, tools/call)
│   └── tools.go        → Tool definitions and execution logic
└── vertexai/
    └── client.go       → Vertex AI genai client wrapper
```

**Request flow**: HTTP POST → `Server.ServeHTTP` → JSON-RPC dispatch → `ToolHandler.ExecuteTool` → `vertexai.Client.GenerateContent` → Response

## MCP Tools Exposed

1. **gemini_query** - General prompts via `GenerateContent`
2. **gemini_query_with_search** - Prompts with "current information" hint via `GenerateContentWithWebSearch`
3. **gemini_code_review** - Structured code review returning JSON with scores and issues

## Key Dependencies

- `google.golang.org/genai` - Google Gen AI unified SDK
- `github.com/joho/godotenv` - Auto-loads `.env` file on startup

## Testing the Server

Manual testing via curl (server must be running):

```bash
# Initialize
curl -X POST http://localhost:8080 -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}'

# List tools
curl -X POST http://localhost:8080 -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'

# Call tool
curl -X POST http://localhost:8080 -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"gemini_query","arguments":{"prompt":"Hello"}}}'
```

## Adding New Tools

1. Add tool definition in `GetTools()` in [internal/mcp/tools.go](internal/mcp/tools.go)
2. Add case in `ExecuteTool()` switch statement
3. Implement handler function (e.g., `executeNewTool`)
4. If new Vertex AI functionality needed, add method to [internal/vertexai/client.go](internal/vertexai/client.go)
