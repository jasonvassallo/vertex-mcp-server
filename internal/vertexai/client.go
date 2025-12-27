package vertexai

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

// Client wraps Vertex AI Gemini API access
type Client struct {
	genaiClient *genai.Client
	projectID   string
	location    string
	modelName   string
}

// Config holds Vertex AI configuration
type Config struct {
	ProjectID string
	Location  string
	ModelName string
}

// NewClient creates a new Vertex AI client
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	// Validate configuration
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("GOOGLE_CLOUD_PROJECT is required")
	}
	if cfg.Location == "" {
		cfg.Location = "global"
	}
	if cfg.ModelName == "" {
		cfg.ModelName = "gemini-3-pro-preview"
	}

	// Check for credentials
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Println("WARNING: GOOGLE_APPLICATION_CREDENTIALS not set. Using application default credentials.")
	}

	// Create Vertex AI client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  cfg.ProjectID,
		Location: cfg.Location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Vertex AI client: %w", err)
	}

	return &Client{
		genaiClient: client,
		projectID:   cfg.ProjectID,
		location:    cfg.Location,
		modelName:   cfg.ModelName,
	}, nil
}

// GenerateContent sends a prompt to Gemini and returns the response
func (c *Client) GenerateContent(ctx context.Context, prompt string) (string, error) {
	// Generate content
	resp, err := c.genaiClient.Models.GenerateContent(ctx, c.modelName,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			Temperature:     genai.Ptr[float64](0.7),
			TopP:            genai.Ptr[float64](0.95),
			MaxOutputTokens: genai.Ptr[int64](8192),
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: false, // Hide raw thinking process from final output
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract text from response
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts in response")
	}

	return resp.Candidates[0].Content.Parts[0].Text, nil
}

// GenerateContentWithWebSearch sends a prompt optimized for current/recent information
// Note: Web search grounding is available through Google AI Studio API but not yet
// in the Vertex AI Go SDK. This uses the same generation logic with a prompt hint.
func (c *Client) GenerateContentWithWebSearch(ctx context.Context, prompt string) (string, error) {
	// Add a system instruction to use current knowledge
	enhancedPrompt := "Please provide current, up-to-date information if available. " + prompt

	// Generate content
	resp, err := c.genaiClient.Models.GenerateContent(ctx, c.modelName,
		genai.Text(enhancedPrompt),
		&genai.GenerateContentConfig{
			Temperature:     genai.Ptr[float64](0.7),
			TopP:            genai.Ptr[float64](0.95),
			MaxOutputTokens: genai.Ptr[int64](8192),
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract text from response
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts in response")
	}

	return resp.Candidates[0].Content.Parts[0].Text, nil
}

// Close closes the Vertex AI client
func (c *Client) Close() error {
	// The new genai.Client does not have a Close method that needs to be called
	return nil
}
