# Integration Guide: Vertex MCP Server

Complete guide for integrating your Vertex AI MCP server with Claude Code, GitHub Copilot, and other AI tools.

## Quick Start

### 1. Start the MCP Server

```bash
cd "/Users/jasonvassallo/Library/Mobile Documents/com~apple~CloudDocs/Documents/Code/vertex-mcp-server"
source .env
./bin/vertex-mcp-server
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

Keep this terminal window open - the server needs to run continuously.

## Integration with Claude Code (VS Code)

### Option 1: Command-Line Configuration (Recommended)

Claude Code typically looks for MCP server configuration in `~/.config/claude/mcp_settings.json` (Linux/Mac) or `%APPDATA%\claude\mcp_settings.json` (Windows).

Create or edit this file:

```bash
mkdir -p ~/.config/claude
cat > ~/.config/claude/mcp_settings.json << 'EOF'
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
EOF
```

### Option 2: HTTP Endpoint (Alternative)

If Claude Code supports HTTP-based MCP servers:

```json
{
  "mcpServers": {
    "vertex-gemini": {
      "url": "http://localhost:8080",
      "type": "http"
    }
  }
}
```

### Option 3: VS Code Settings (If Supported)

Some Claude Code installations read from VS Code settings. Add to `.vscode/settings.json`:

```json
{
  "claude.mcpServers": {
    "vertex-gemini": {
      "command": "/Users/jasonvassallo/Library/Mobile Documents/com~apple~CloudDocs/Documents/Code/vertex-mcp-server/bin/vertex-mcp-server",
      "env": {
        "GOOGLE_CLOUD_PROJECT": "ai-orchestrator-482302",
        "GOOGLE_CLOUD_LOCATION": "global",
        "GOOGLE_APPLICATION_CREDENTIALS": "/Users/jasonvassallo/.ai_orchestrator/vertex-ai-key.json"
      }
    }
  }
}
```

### Testing Claude Code Integration

1. Restart VS Code after configuring
2. Open a new Claude Code chat
3. Try asking: "Use the gemini_query tool to explain what Go is"
4. Claude should automatically detect and use your MCP server

## GitHub Copilot Integration

### Understanding the Limitation

**Important:** GitHub Copilot does **not** natively support the Model Context Protocol (MCP). It uses its own proprietary protocol and model infrastructure.

### Why You Can't Directly Use Gemini with Copilot

GitHub Copilot is:
- Tied to OpenAI's Codex models
- Deeply integrated with Microsoft's infrastructure
- Not designed to use external model endpoints
- Not compatible with MCP protocol

### Your Options for Using Gemini with GitHub Copilot

#### Option 1: Use Both Tools Separately (Recommended)

This is what most developers do:

**GitHub Copilot:**
- Inline code completions
- Quick suggestions while typing
- Code generation from comments

**Claude Code (with your MCP server):**
- Complex queries and explanations
- Code review and analysis
- Architecture discussions
- Access to Gemini models via your MCP server

**Practical Workflow:**
```
1. Use Copilot for: day-to-day coding, autocomplete, quick functions
2. Use Claude Code + MCP for: complex problems, architectural questions, code review
3. Use Gemini directly for: specialized tasks that benefit from Gemini's strengths
```

#### Option 2: VS Code Tasks (Quick Access to Gemini)

Create quick access to your MCP server via VS Code tasks.

Add to `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Ask Gemini",
      "type": "shell",
      "command": "curl -s -X POST http://localhost:8080 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"gemini_query\",\"arguments\":{\"prompt\":\"${input:geminiPrompt}\"}}}' | jq -r '.result.content[0].text'",
      "problemMatcher": [],
      "presentation": {
        "echo": true,
        "reveal": "always",
        "focus": false,
        "panel": "new"
      }
    },
    {
      "label": "Gemini Code Review",
      "type": "shell",
      "command": "curl -s -X POST http://localhost:8080 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"gemini_code_review\",\"arguments\":{\"code\":\"${selectedText}\",\"language\":\"${input:language}\"}}}' | jq -r '.result.content[0].text'",
      "problemMatcher": [],
      "presentation": {
        "echo": true,
        "reveal": "always",
        "focus": false,
        "panel": "new"
      }
    }
  ],
  "inputs": [
    {
      "id": "geminiPrompt",
      "type": "promptString",
      "description": "What would you like to ask Gemini?",
      "default": "Explain this code"
    },
    {
      "id": "language",
      "type": "promptString",
      "description": "Programming language?",
      "default": "python"
    }
  ]
}
```

**Usage:**
1. Press `Cmd+Shift+P` (Mac) or `Ctrl+Shift+P` (Windows/Linux)
2. Type "Tasks: Run Task"
3. Select "Ask Gemini" or "Gemini Code Review"
4. Enter your prompt

#### Option 3: Custom VS Code Extension (Advanced)

Build a custom VS Code extension that:
- Provides a command palette interface to your MCP server
- Displays results in a webview panel
- Integrates with VS Code's editor

This requires TypeScript/JavaScript knowledge but gives you full control.

#### Option 4: Terminal Aliases (Quickest)

Add to your `~/.zshrc` or `~/.bashrc`:

```bash
# Quick Gemini query
gemini() {
  curl -s -X POST http://localhost:8080 \
    -H "Content-Type: application/json" \
    -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"gemini_query\",\"arguments\":{\"prompt\":\"$*\"}}}" \
    | jq -r '.result.content[0].text'
}

# Code review from clipboard
gemini-review() {
  local code=$(pbpaste)  # macOS - use xclip on Linux
  local lang=${1:-python}
  curl -s -X POST http://localhost:8080 \
    -H "Content-Type: application/json" \
    -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"gemini_code_review\",\"arguments\":{\"code\":\"$(echo "$code" | jq -Rs .)\",\"language\":\"$lang\"}}}" \
    | jq -r '.result.content[0].text'
}
```

**Usage:**
```bash
# In terminal
gemini "Explain quantum computing"
gemini-review python  # Reviews code from clipboard
```

### What About "Codex Integration"?

**Codex** is the model behind GitHub Copilot. You cannot:
- Swap Codex for Gemini in GitHub Copilot
- Make Copilot use your Vertex AI models
- Use MCP with GitHub Copilot's infrastructure

**However, you CAN:**
- Use OpenAI's Codex API directly (if you have access)
- Build your own Copilot-like tool using Gemini
- Use Continue.dev (see below)

## Alternative: Continue.dev Extension

[Continue.dev](https://continue.dev) is an open-source Copilot alternative that DOES support custom models:

### Setup Continue.dev with Gemini

1. Install Continue.dev extension in VS Code
2. Configure `~/.continue/config.json`:

```json
{
  "models": [
    {
      "title": "Gemini via MCP",
      "provider": "openai",
      "model": "gpt-4",
      "apiBase": "http://localhost:8080",
      "apiKey": "dummy"
    }
  ]
}
```

Note: This requires adapting your MCP server to speak OpenAI-compatible API, which would require additional wrapper code.

## Summary: Best Setup for Your Use Case

Based on your requirements (access Gemini 3 Pro without Google's VS Code extension):

### Recommended Setup:

1. **Primary: Claude Code + Your MCP Server** ✅
   - Full MCP protocol support
   - Direct access to Gemini models
   - Rich tool integration
   - Already configured!

2. **Secondary: GitHub Copilot (as is)** ✅
   - Keep using it for inline completions
   - Don't try to force Gemini integration
   - Complementary to Claude Code

3. **Quick Access: VS Code Tasks** ✅
   - Add the tasks.json above
   - Quick Gemini queries without switching tools
   - Uses your MCP server

4. **Terminal Access: Shell Aliases** ✅
   - Add the aliases above
   - Use Gemini from command line
   - Perfect for quick questions

### What You've Achieved:

✅ Local MCP server running Gemini via Vertex AI
✅ No dependency on Google's VS Code extension
✅ Full control over model selection and configuration
✅ Can integrate with any MCP-compatible tool
✅ Cost-effective ($0.10/$0.40 per 1M tokens for Flash)
✅ Enterprise-grade with Vertex AI SLAs

### What You Can't Do (Limitations):

❌ Make GitHub Copilot use Gemini (architectural limitation)
❌ Replace Codex in Copilot (proprietary system)
❌ Use MCP with tools that don't support it

## Next Steps

1. **Test the Integration:**
   ```bash
   cd vertex-mcp-server
   source .env
   ./bin/vertex-mcp-server &

   # In another terminal
   ./test-server.sh
   ```

2. **Configure Claude Code:**
   - Choose one of the configuration methods above
   - Restart VS Code
   - Test with a simple query

3. **Optional: Add VS Code Tasks:**
   - Copy the tasks.json configuration
   - Test with `Cmd+Shift+P` → "Tasks: Run Task"

4. **Optional: Add Shell Aliases:**
   - Add to `~/.zshrc`
   - Run `source ~/.zshrc`
   - Test with `gemini "hello world"`

## Troubleshooting

### MCP Server Not Starting

```bash
# Check if Vertex AI key exists
ls -l ~/.ai_orchestrator/vertex-ai-key.json

# Check if project ID is set
echo $GOOGLE_CLOUD_PROJECT

# Test Vertex AI access
gcloud auth application-default print-access-token
```

### Claude Code Not Seeing MCP Server

1. Check Claude Code documentation for exact config location
2. Try all three configuration methods
3. Check Claude Code logs for errors
4. Ensure server is running (`lsof -i :8080`)

### Permission Denied Errors

```bash
# Check service account permissions
gcloud projects get-iam-policy ai-orchestrator-482302 \
  --flatten="bindings[].members" \
  --filter="bindings.members:serviceAccount:gemini-api-access*"

# Should show roles/aiplatform.user
```

## Additional Resources

- [Model Context Protocol Specification](https://modelcontextprotocol.io)
- [Vertex AI Documentation](https://cloud.google.com/vertex-ai/docs)
- [Claude Code Documentation](https://claude.com/code)
- [Continue.dev Documentation](https://continue.dev/docs)

## Support

For issues with:
- **MCP Server:** Check [README.md](README.md) and server logs
- **Vertex AI:** Review [VERTEX_AI_SETUP.md](../ai-orchestrator/VERTEX_AI_SETUP.md)
- **Claude Code:** Visit Claude Code documentation
- **GitHub Copilot:** Contact GitHub Support (they won't help with custom models)
