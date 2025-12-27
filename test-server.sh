#!/bin/bash

# Test script for Vertex AI MCP Server
# This script tests all available tools

set -e

PORT="${PORT:-8080}"
SERVER_URL="http://localhost:${PORT}"
BOLD='\033[1m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BOLD}Testing Vertex AI MCP Server${NC}"
echo ""

# Check if server is running
if ! curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL" > /dev/null 2>&1; then
    echo -e "${RED}❌ Server is not running at $SERVER_URL${NC}"
    echo "Please start the server first:"
    echo "  cd vertex-mcp-server"
    echo "  source .env"
    echo "  ./bin/vertex-mcp-server"
    exit 1
fi

echo -e "${GREEN}✓ Server is running${NC}"
echo ""

# Test 1: Initialize
echo -e "${BLUE}Test 1: Initialize${NC}"
echo "----------------------------------------"
curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2025-06-18"
    }
  }' | python3 -m json.tool
echo ""
echo ""

# Test 2: List Tools
echo -e "${BLUE}Test 2: List Available Tools${NC}"
echo "----------------------------------------"
curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  }' | python3 -m json.tool
echo ""
echo ""

# Test 3: Simple Query
echo -e "${BLUE}Test 3: Simple Gemini Query${NC}"
echo "----------------------------------------"
echo "Prompt: 'Explain what Go is in 2 sentences'"
echo ""
curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "gemini_query",
      "arguments": {
        "prompt": "Explain what Go is in 2 sentences"
      }
    }
  }' | python3 -m json.tool
echo ""
echo ""

# Test 4: Web Search Query
echo -e "${BLUE}Test 4: Query with Web Search${NC}"
echo "----------------------------------------"
echo "Prompt: 'What are the latest features in Go 1.25?'"
echo ""
curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/call",
    "params": {
      "name": "gemini_query_with_search",
      "arguments": {
        "prompt": "What are the latest features in Go 1.25?"
      }
    }
  }' | python3 -m json.tool
echo ""
echo ""

# Test 5: Code Review
echo -e "${BLUE}Test 5: Code Review${NC}"
echo "----------------------------------------"
echo "Reviewing a simple Go function..."
echo ""
curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 5,
    "method": "tools/call",
    "params": {
      "name": "gemini_code_review",
      "arguments": {
        "code": "func Add(a, b int) int {\n    return a + b\n}",
        "language": "go"
      }
    }
  }' | python3 -m json.tool
echo ""
echo ""

echo -e "${GREEN}✓ All tests completed!${NC}"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "1. Configure Claude Code to use this MCP server"
echo "2. Add to VS Code settings:"
echo '   {
     "claude.mcpServers": {
       "vertex-gemini": {
         "type": "http",
         "url": "http://localhost:8080",
         "enabled": true
       }
     }
   }'
echo ""
