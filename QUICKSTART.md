# Quick Start Guide

Get your Vertex AI MCP server running in 3 steps.

## Step 1: Start the Server

```bash
cd "/Users/jasonvassallo/Library/Mobile Documents/com~apple~CloudDocs/Documents/Code/vertex-mcp-server"
source .env
./bin/vertex-mcp-server
```

**Expected output:**
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

Available tools:
  - gemini_query: Query Gemini AI
  - gemini_query_with_search: Query with web search
  - gemini_code_review: AI-powered code review

Press Ctrl+C to stop the server
```

**Keep this terminal open!** The server needs to run continuously.

## Step 2: Test It (New Terminal)

```bash
cd "/Users/jasonvassallo/Library/Mobile Documents/com~apple~CloudDocs/Documents/Code/vertex-mcp-server"
./test-server.sh
```

This tests all three tools and verifies everything works.

## Step 3: Use with Claude Code

### Option A: Auto-start via Config (Recommended)

Create `~/.config/claude/mcp_settings.json`:

```json
{
  "mcpServers": {
    "vertex-gemini": {
      "command": "/Users/jasonvassallo/Library/Mobile Documents/com~apple~CloudDocs/Documents/Code/vertex-mcp-server/bin/vertex-mcp-server",
      "env": {
        "GOOGLE_CLOUD_PROJECT": "ai-orchestrator-482302",
        "GOOGLE_CLOUD_LOCATION": "global",
        "GEMINI_MODEL": "gemini-3-pro-preview",
        "GOOGLE_APPLICATION_CREDENTIALS": "/Users/jasonvassallo/.ai_orchestrator/vertex-ai-key.json",
        "PORT": "8080"
      }
    }
  }
}
```

Claude Code will automatically start and manage the server.

### Option B: Manual Connection

If server is already running on localhost:8080, Claude Code may auto-detect it.

## Common Commands

### Start Server in Background
```bash
cd vertex-mcp-server
source .env
./bin/vertex-mcp-server > server.log 2>&1 &
echo $! > server.pid  # Save process ID
```

### Stop Background Server
```bash
kill $(cat server.pid)
rm server.pid
```

### Check if Server is Running
```bash
lsof -i :8080
# or
curl -s -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | jq .
```

### Quick Gemini Query (CLI)
```bash
curl -s -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "gemini_query",
      "arguments": {
        "prompt": "Explain Go in 2 sentences"
      }
    }
  }' | jq -r '.result.content[0].text'
```

## Troubleshooting

### "Failed to create Vertex AI client"
**Fix:** Check your service account key exists:
```bash
ls -l ~/.ai_orchestrator/vertex-ai-key.json
```

### "Address already in use"
**Fix:** Another process is using port 8080:
```bash
lsof -i :8080  # Find the process
kill <PID>     # Stop it
```

### "PERMISSION_DENIED"
**Fix:** Check service account permissions:
```bash
gcloud projects get-iam-policy ai-orchestrator-482302 \
  --flatten="bindings[].members" \
  --filter="bindings.members:serviceAccount:gemini-api-access*"
```

Should show `roles/aiplatform.user`.

## Next Steps

1. ✅ Server is running
2. ✅ Tests pass
3. ✅ Configure Claude Code (see INTEGRATION_GUIDE.md)
4. ✅ Try advanced features (code review, web search)
5. ✅ Deploy to production (see README.md)

## Need Help?

- **Full Documentation:** [README.md](README.md)
- **Integration Guide:** [INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md)
- **Vertex AI Setup:** [../ai-orchestrator/VERTEX_AI_SETUP.md](../ai-orchestrator/VERTEX_AI_SETUP.md)
