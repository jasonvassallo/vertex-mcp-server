# Vertex AI MCP Server

A Model Context Protocol (MCP) server that exposes Google Cloud Vertex AI Gemini models to AI assistants like Claude Code, GitHub Copilot, and other MCP-compatible tools.

## Features

- 🤖 **Three powerful tools**:
  - `gemini_query`: General-purpose AI queries
  - `gemini_query_with_search`: Real-time web search-grounded responses
  - `gemini_code_review`: AI-powered code review with detailed feedback

- 🚀 **Production-ready**:
  - HTTP transport for easy integration
  - Graceful shutdown handling
  - Comprehensive error handling
  - JSON-RPC 2.0 compliant

- 🔒 **Secure**:
  - Uses Google Cloud service account authentication
  - No hardcoded credentials
  - Environment-based configuration

## Prerequisites

1. **Go 1.24+** (Already installed ✓)
2. **Google Cloud Project** with:
   - Vertex AI API enabled
   - Service account with `roles/aiplatform.user` permission
   - Service account key JSON file

## Quick Start

### 1. Setup Google Cloud Credentials

If you haven't already set up Vertex AI, follow these steps:

```bash
# Set your project
export GCP_PROJECT_ID=ai-orchestrator-482302

# Enable Vertex AI API
gcloud services enable aiplatform.googleapis.com

# Create service account (if not exists)
gcloud iam service-accounts create gemini-api-access \
    --display-name="Gemini API Access"

# Grant permissions
gcloud projects add-iam-policy-binding $GCP_PROJECT_ID \
    --member="serviceAccount:gemini-api-access@${GCP_PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/aiplatform.user"

# Create key
mkdir -p ~/.ai_orchestrator
gcloud iam service-accounts keys create ~/.ai_orchestrator/vertex-ai-key.json \
    --iam-account=gemini-api-access@${GCP_PROJECT_ID}.iam.gserviceaccount.com

chmod 600 ~/.ai_orchestrator/vertex-ai-key.json
```

### 2. Configure Environment

```bash
# Copy example env file
cp .env.example .env

# Edit .env with your values
# Update GOOGLE_APPLICATION_CREDENTIALS with your actual username
```

### 3. Build and Run

```bash
# Build the server
go build -o bin/gemini-mcp ./cmd/server

# Run the server (env loaded automatically from .env)
./bin/gemini-mcp

# Or install to PATH for global access
go build -o ~/.local/bin/gemini-mcp ./cmd/server
gemini-mcp
```

You should see:

```
Starting Vertex AI MCP Server...
Configuration:
  Project ID: ai-orchestrator-482302
  Location: global
  Model: gemini-3-pro-preview
  Port: 8080
✓ Vertex AI client initialized
✓ MCP server created
🚀 Server listening on http://localhost:8080
```

## Testing the Server

### Test with curl

```bash
# Test initialization
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2025-06-18"
    }
  }'

# List available tools
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  }'

# Call gemini_query tool
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "gemini_query",
      "arguments": {
        "prompt": "Explain quantum computing in simple terms"
      }
    }
  }'

# Call code review tool
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/call",
    "params": {
      "name": "gemini_code_review",
      "arguments": {
        "code": "def factorial(n):\n    if n == 0: return 1\n    return n * factorial(n-1)",
        "language": "python"
      }
    }
  }'
```

## Integration with AI Tools

### Claude Code (VS Code Extension)

Add to your VS Code settings (`.vscode/settings.json` or global settings):

```json
{
  "claude.mcpServers": {
    "vertex-gemini": {
      "type": "http",
      "url": "http://localhost:8080",
      "enabled": true
    }
  }
}
```

Or using Claude Code CLI configuration (`~/.claude/config.json`):

```json
{
  "mcpServers": {
    "vertex-gemini": {
      "command": "/path/to/vertex-mcp-server/bin/vertex-mcp-server"
    }
  }
}
```

### GitHub Copilot (via MCP Bridge)

GitHub Copilot doesn't natively support MCP yet, but you can use it through tools that bridge the protocols. See the "GitHub Copilot Integration" section below.

### Cline (VS Code Extension)

Add to Cline's MCP settings (`.cline/settings.json`):

```json
{
  "mcpServers": {
    "vertex-gemini": {
      "httpUrl": "http://localhost:8080"
    }
  }
}
```

## Available Tools

### 1. gemini_query

Query Gemini for general tasks, coding help, explanations, etc.

**Parameters:**
- `prompt` (required): The question or prompt

**Example:**
```json
{
  "name": "gemini_query",
  "arguments": {
    "prompt": "Write a Python function to reverse a string"
  }
}
```

### 2. gemini_query_with_search

Query Gemini with real-time web search grounding. Best for current events, recent data, or information that might have changed.

**Parameters:**
- `prompt` (required): The question or prompt

**Example:**
```json
{
  "name": "gemini_query_with_search",
  "arguments": {
    "prompt": "What are the latest features in Go 1.25?"
  }
}
```

### 3. gemini_code_review

Get AI-powered code review with detailed feedback.

**Parameters:**
- `code` (required): The code to review
- `language` (required): Programming language (e.g., 'go', 'python', 'javascript')
- `focus` (optional): Specific aspect to focus on (e.g., 'security', 'performance')

**Example:**
```json
{
  "name": "gemini_code_review",
  "arguments": {
    "code": "func Add(a, b int) int { return a + b }",
    "language": "go",
    "focus": "best practices"
  }
}
```

## Configuration

All configuration is done via environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `GOOGLE_CLOUD_PROJECT` | Your GCP project ID | - | Yes |
| `GOOGLE_CLOUD_LOCATION` | GCP region | `global` | No |
| `GEMINI_MODEL` | Gemini model to use | `gemini-3-pro-preview` | No |
| `GOOGLE_APPLICATION_CREDENTIALS` | Path to service account key | - | Yes |
| `PORT` | Server port | `8080` | No |

### Model Options

- `gemini-3-pro-preview` (Default): Latest reasoning model (Global only)
- `gemini-3-flash-preview`: Fast reasoning model (Global only)
- `gemini-2.0-flash`: Production-ready fast model (Stable)
- `gemini-1.5-pro`: Advanced reasoning
- `gemini-1.5-flash`: Cost-effective

## Development

### Project Structure

```
vertex-mcp-server/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── mcp/
│   │   ├── server.go         # MCP server implementation
│   │   └── tools.go          # Tool definitions and handlers
│   └── vertexai/
│       └── client.go         # Vertex AI client wrapper
├── config/                   # Configuration files
├── .env.example             # Example environment file
├── README.md                # This file
└── go.mod                   # Go module definition
```

### Building

```bash
# Build for current platform
go build -o bin/vertex-mcp-server ./cmd/server

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/vertex-mcp-server-linux ./cmd/server
GOOS=darwin GOARCH=arm64 go build -o bin/vertex-mcp-server-darwin ./cmd/server
GOOS=windows GOARCH=amd64 go build -o bin/vertex-mcp-server.exe ./cmd/server
```

### Running in Development

```bash
# Run directly with go run
source .env
go run ./cmd/server
```

## Deployment

### Docker

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o vertex-mcp-server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/vertex-mcp-server .
ENV PORT=8080
EXPOSE 8080
CMD ["./vertex-mcp-server"]
```

### Cloud Run

```bash
# Build and deploy to Cloud Run
gcloud run deploy vertex-mcp-server \
  --source . \
  --region us-central1 \
  --set-env-vars GOOGLE_CLOUD_PROJECT=ai-orchestrator-482302 \
  --set-env-vars GOOGLE_CLOUD_LOCATION=us-central1 \
  --allow-unauthenticated
```

## Troubleshooting

### "Failed to create Vertex AI client"

**Solution:** Check your credentials:
```bash
echo $GOOGLE_APPLICATION_CREDENTIALS
ls -l $GOOGLE_APPLICATION_CREDENTIALS
```

### "PERMISSION_DENIED"

**Solution:** Ensure your service account has the correct role:
```bash
gcloud projects get-iam-policy $GCP_PROJECT_ID \
  --flatten="bindings[].members" \
  --filter="bindings.members:serviceAccount:gemini-api-access*"
```

### "Model not found"

**Solution:** Check available models in your region:
```bash
gcloud ai models list --region=us-central1 | grep gemini
```

### Connection refused when testing

**Solution:** Make sure the server is running:
```bash
# Check if server is running
lsof -i :8080

# Check logs
./bin/vertex-mcp-server
```

## Cost Considerations

Using this MCP server incurs Google Cloud Vertex AI costs:

| Model | Input (per 1M tokens) | Output (per 1M tokens) |
|-------|----------------------|------------------------|
| Gemini 1.5 Flash | $0.075 | $0.30 |
| Gemini 1.5 Pro | $1.25 | $5.00 |

**Tip:** Use `gemini-1.5-flash` for most tasks to minimize costs.

## Security Best Practices

1. **Never commit credentials**:
   ```bash
   echo ".env" >> .gitignore
   echo "*.json" >> .gitignore
   ```

2. **Restrict service account permissions**:
   - Only grant `roles/aiplatform.user` (not `owner` or `editor`)

3. **Rotate keys regularly**:
   ```bash
   # List keys
   gcloud iam service-accounts keys list \
     --iam-account=gemini-api-access@${GCP_PROJECT_ID}.iam.gserviceaccount.com

   # Delete old keys
   gcloud iam service-accounts keys delete KEY_ID \
     --iam-account=gemini-api-access@${GCP_PROJECT_ID}.iam.gserviceaccount.com
   ```

4. **Use environment variables**: Never hardcode credentials in source code

5. **Monitor usage**: Set up billing alerts in Google Cloud Console

## GitHub Copilot Integration

GitHub Copilot doesn't natively support MCP protocol. However, you can:

### Option 1: Use Both Tools Separately
- Use GitHub Copilot for inline code suggestions
- Use Claude Code (with this MCP server) for complex queries and code review

### Option 2: MCP Bridge (Advanced)
Create a bridge that translates between GitHub Copilot's API and MCP protocol. This would require custom middleware.

### Option 3: VS Code Tasks
Create VS Code tasks that call the MCP server via curl:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Ask Gemini",
      "type": "shell",
      "command": "curl -X POST http://localhost:8080 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"gemini_query\",\"arguments\":{\"prompt\":\"${input:prompt}\"}}}'",
      "problemMatcher": []
    }
  ],
  "inputs": [
    {
      "id": "prompt",
      "type": "promptString",
      "description": "What would you like to ask Gemini?"
    }
  ]
}
```

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - see LICENSE file for details

## Support

For issues:
1. Check the [Troubleshooting](#troubleshooting) section
2. Review [Vertex AI Documentation](https://cloud.google.com/vertex-ai/docs)
3. Check [MCP Specification](https://modelcontextprotocol.io)

## Acknowledgments

- Built with [Google Cloud Vertex AI](https://cloud.google.com/vertex-ai)
- Uses [Model Context Protocol](https://modelcontextprotocol.io)
- Inspired by [Google Codelabs MCP Tutorial](https://codelabs.developers.google.com/cloud-gemini-cli-mcp-go)
