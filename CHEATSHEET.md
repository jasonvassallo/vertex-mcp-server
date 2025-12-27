# MCP Server & Gemini Quick Reference

> **TL;DR:** Type `mcp-server start` then use `gemini "your question"` anywhere in terminal!

## 🚀 Server Management (From Anywhere)

```bash
mcp-server start      # Start the MCP server
mcp-server stop       # Stop the server
mcp-server restart    # Restart the server
mcp-server status     # Check if running
mcp-server test       # Test all 3 tools
mcp-server logs       # View live logs (Ctrl+C to exit)
```

**Shortcuts:**
```bash
mcp-start             # Same as mcp-server start
mcp-stop              # Same as mcp-server stop
mcp-status            # Same as mcp-server status
```

## 🤖 Quick Gemini Access (Terminal)

### Ask Anything
```bash
gemini "Explain quantum computing"
gemini "Write a Python function to reverse a string"
gemini "What's the difference between Go and Rust?"
```

### Web Search (Current Info)
```bash
gemini-search "Latest features in Go 1.25"
gemini-search "Current Python best practices 2025"
```

### Code Review
```bash
# Review a file
gemini-review myfile.py
gemini-review script.js javascript

# Review code from clipboard (copy code first)
gemini-review-clipboard python
gemini-review-clipboard go
```

## 🎯 VS Code Integration

### Method 1: Tasks (Cmd+Shift+P → "Run Task")
- **MCP: Start Server** - Start the server
- **MCP: Stop Server** - Stop the server
- **MCP: Server Status** - Check status
- **MCP: Test Server** - Run tests
- **Gemini: Ask Question** - Interactive prompt
- **Gemini: Review Selected Code** - Select code, run task
- **Gemini: Web Search Query** - Search with current info

### Method 2: Claude Code (Me!)
I can automatically use your MCP server! Just ask me to use Gemini:
- "Use gemini_query to explain this code"
- "Use gemini_code_review to review this function"
- "Use gemini_query_with_search for latest React features"

## 📁 File Locations

```
Server:        ~/bin/mcp-server
Binary:        vertex-mcp-server/bin/vertex-mcp-server
Logs:          /tmp/mcp-server.log
PID File:      /tmp/mcp-server.pid
Config:        ~/.config/claude/mcp_settings.json
API Keys:      ~/.api_keys (secure, 600 permissions)
```

## 🔧 Common Tasks

### First Time Setup (Already Done!)
```bash
✓ Go installed
✓ Vertex AI credentials created
✓ Server built
✓ Commands in PATH
✓ Shell aliases configured
✓ Claude Code configured
✓ VS Code tasks created
```

### Daily Workflow

**Option A: Auto-start via Claude Code**
- Just use Claude Code - server starts automatically
- No manual commands needed!

**Option B: Manual start**
```bash
# Morning: Start server
mcp-server start

# Use all day via terminal or VS Code tasks
gemini "your questions"

# Evening: Stop server (or leave running)
mcp-server stop
```

### Troubleshooting

**Server won't start?**
```bash
mcp-server stop          # Force stop
mcp-server start         # Try again
mcp-server logs          # Check errors
```

**API not responding?**
```bash
mcp-status               # Is it running?
lsof -i :8080           # Is port 8080 open?
cat /tmp/mcp-server.log # Check logs
```

**Need to reload shell config?**
```bash
source ~/.zshrc         # Reload aliases
```

## 💡 Pro Tips

### Tip 1: Background Server
Server runs in background by default. You can close the terminal!

### Tip 2: Chain Commands
```bash
# Start and test in one go
mcp-server start && sleep 2 && mcp-server test
```

### Tip 3: Quick Status Check
```bash
# Add to your prompt (optional)
# Edit ~/.zshrc and add near the end:
# PROMPT="$PROMPT \$(mcp-status 2>/dev/null | grep -q running && echo '🤖' || echo '')"
```

### Tip 4: Alias Your Own Commands
Add to `~/.zshrc`:
```bash
alias ask='gemini'
alias review='gemini-review'
alias ai='gemini'
```

## 🎨 Integration Summary

| Method | Use Case | How to Access |
|--------|----------|---------------|
| **Terminal** | Quick questions anywhere | `gemini "question"` |
| **VS Code Tasks** | While coding in VS Code | `Cmd+Shift+P` → Tasks |
| **Claude Code** | Complex tasks, context-aware | Ask me to use Gemini tools |

## 📊 Costs

- **Gemini 2.5 Flash**: $0.10 input / $0.40 output per 1M tokens
- **Typical query**: ~$0.001 (less than a penny!)
- **Free credits**: $300 for 3 months (Google Cloud new users)

## 🔒 Security

```bash
# API keys location (secure)
~/.api_keys              # Your API keys (600 permissions)

# Vertex AI credentials
~/.ai_orchestrator/vertex-ai-key.json  # Service account (600 permissions)

# NEVER commit these files to git!
# Both are in .gitignore ✓
```

## ⚡ Quick Examples

```bash
# Code help
gemini "How do I read a file in Go?"

# Debugging
gemini "Why would this Python code raise TypeError?"

# Learning
gemini-search "New features in TypeScript 5.0"

# Code review
gemini-review app.py python

# Architecture
gemini "Should I use microservices or monolith for a small startup?"

# Current events
gemini-search "Latest AI breakthroughs December 2025"
```

## 📚 Full Documentation

- [README.md](README.md) - Complete reference
- [INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md) - All integration options
- [QUICKSTART.md](QUICKSTART.md) - Quick setup guide

---

**Remember:** Just type `mcp-server start` once, then `gemini "your question"` anywhere! 🚀
