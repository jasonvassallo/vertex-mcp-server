package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"vertex-mcp-server/internal/vertexai"
)

// Tool represents an MCP tool definition
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema ToolInputSchema `json:"inputSchema"`
}

// ToolInputSchema defines the JSON schema for tool input
type ToolInputSchema struct {
	Type       string                            `json:"type"`
	Properties map[string]map[string]interface{} `json:"properties"`
	Required   []string                          `json:"required"`
}

// CallToolResult represents the result of a tool call
type CallToolResult struct {
	Content []interface{} `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// TextContent represents text content in a tool result
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ToolHandler manages MCP tools and their execution
type ToolHandler struct {
	vertexClient *vertexai.Client
}

// NewToolHandler creates a new tool handler
func NewToolHandler(vertexClient *vertexai.Client) *ToolHandler {
	return &ToolHandler{
		vertexClient: vertexClient,
	}
}

// GetTools returns the list of available MCP tools
func (h *ToolHandler) GetTools() []Tool {
	return []Tool{
		{
			Name:        "gemini_query",
			Description: "Query Google's Gemini AI model with any prompt. Best for general questions, coding assistance, explanations, and creative tasks.",
			InputSchema: ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"prompt": {
						"type":        "string",
						"description": "The prompt or question to send to Gemini",
					},
				},
				Required: []string{"prompt"},
			},
		},
		{
			Name:        "gemini_query_with_search",
			Description: "Query Google's Gemini AI model optimized for current information. Use this for queries that need real-time or recent data.",
			InputSchema: ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"prompt": {
						"type":        "string",
						"description": "The prompt or question to send to Gemini",
					},
				},
				Required: []string{"prompt"},
			},
		},
		{
			Name:        "gemini_code_review",
			Description: "Have Gemini review code and provide improvement suggestions. Automatically analyzes code quality, potential bugs, best practices, and optimization opportunities.",
			InputSchema: ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"code": {
						"type":        "string",
						"description": "The code to review",
					},
					"language": {
						"type":        "string",
						"description": "Programming language (e.g., 'go', 'python', 'javascript')",
					},
					"focus": {
						"type":        "string",
						"description": "Optional: Specific aspects to focus on (e.g., 'security', 'performance', 'readability')",
					},
				},
				Required: []string{"code", "language"},
			},
		},
	}
}

// ExecuteTool executes a tool call and returns the result
func (h *ToolHandler) ExecuteTool(ctx context.Context, toolName string, arguments map[string]interface{}) (*CallToolResult, error) {
	switch toolName {
	case "gemini_query":
		return h.executeGeminiQuery(ctx, arguments)
	case "gemini_query_with_search":
		return h.executeGeminiQueryWithSearch(ctx, arguments)
	case "gemini_code_review":
		return h.executeCodeReview(ctx, arguments)
	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}

func (h *ToolHandler) executeGeminiQuery(ctx context.Context, args map[string]interface{}) (*CallToolResult, error) {
	prompt, ok := args["prompt"].(string)
	if !ok {
		return nil, fmt.Errorf("prompt must be a string")
	}

	response, err := h.vertexClient.GenerateContent(ctx, prompt)
	if err != nil {
		return &CallToolResult{
			Content: []interface{}{
				TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &CallToolResult{
		Content: []interface{}{
			TextContent{
				Type: "text",
				Text: response,
			},
		},
	}, nil
}

func (h *ToolHandler) executeGeminiQueryWithSearch(ctx context.Context, args map[string]interface{}) (*CallToolResult, error) {
	prompt, ok := args["prompt"].(string)
	if !ok {
		return nil, fmt.Errorf("prompt must be a string")
	}

	response, err := h.vertexClient.GenerateContentWithWebSearch(ctx, prompt)
	if err != nil {
		return &CallToolResult{
			Content: []interface{}{
				TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &CallToolResult{
		Content: []interface{}{
			TextContent{
				Type: "text",
				Text: response,
			},
		},
	}, nil
}

func (h *ToolHandler) executeCodeReview(ctx context.Context, args map[string]interface{}) (*CallToolResult, error) {
	code, ok := args["code"].(string)
	if !ok {
		return nil, fmt.Errorf("code must be a string")
	}

	language, ok := args["language"].(string)
	if !ok {
		return nil, fmt.Errorf("language must be a string")
	}

	focus := ""
	if f, ok := args["focus"].(string); ok {
		focus = f
	}

	// Build code review prompt
	prompt := fmt.Sprintf(`Review the following %s code and provide detailed feedback.

Code:
%s

`, language, code)

	if focus != "" {
		prompt += fmt.Sprintf("Focus specifically on: %s\n\n", focus)
	}

	prompt += `Please analyze:
1. Code quality and best practices
2. Potential bugs or issues
3. Performance considerations
4. Security concerns
5. Readability and maintainability
6. Specific improvements you'd recommend

Provide your response in JSON format with this structure:
{
  "summary": "Brief overall assessment",
  "quality_score": "1-10 rating",
  "issues": [
    {"severity": "high|medium|low", "type": "bug|style|performance|security", "description": "...", "suggestion": "..."}
  ],
  "strengths": ["strength 1", "strength 2"],
  "recommendations": ["recommendation 1", "recommendation 2"]
}`

	response, err := h.vertexClient.GenerateContent(ctx, prompt)
	if err != nil {
		return &CallToolResult{
			Content: []interface{}{
				TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Try to parse as JSON, but return raw text if parsing fails
	var reviewData interface{}
	if err := json.Unmarshal([]byte(response), &reviewData); err != nil {
		// Return as plain text if not valid JSON
		return &CallToolResult{
			Content: []interface{}{
				TextContent{
					Type: "text",
					Text: response,
				},
			},
		}, nil
	}

	// Return parsed JSON
	return &CallToolResult{
		Content: []interface{}{
			TextContent{
				Type: "text",
				Text: response,
			},
		},
	}, nil
}
