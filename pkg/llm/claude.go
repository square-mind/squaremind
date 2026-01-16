package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	claudeAPIURL  = "https://api.anthropic.com/v1/messages"
	claudeVersion = "2023-06-01"
)

// ClaudeProvider implements the Provider interface for Anthropic's Claude
type ClaudeProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(apiKey string) *ClaudeProvider {
	return &ClaudeProvider{
		apiKey:  apiKey,
		baseURL: claudeAPIURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		model: string(ModelClaude35Sonnet),
	}
}

// WithBaseURL sets a custom base URL
func (p *ClaudeProvider) WithBaseURL(url string) *ClaudeProvider {
	p.baseURL = url
	return p
}

// WithModel sets the default model
func (p *ClaudeProvider) WithModel(model string) *ClaudeProvider {
	p.model = model
	return p
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// Complete generates a completion
func (p *ClaudeProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	model := req.Model
	if model == "" {
		model = p.model
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	// Build Claude API request
	claudeReq := claudeRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
	}

	if req.System != "" {
		claudeReq.System = req.System
	}

	if req.Temperature > 0 {
		claudeReq.Temperature = &req.Temperature
	}

	if len(req.Stop) > 0 {
		claudeReq.StopSequences = req.Stop
	}

	// Make request
	body, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", claudeVersion)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(respBody, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract content
	content := ""
	for _, block := range claudeResp.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &CompletionResponse{
		Content:      content,
		FinishReason: claudeResp.StopReason,
		TokensUsed:   claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
	}, nil
}

// Chat implements chat completion for Claude
func (p *ClaudeProvider) Chat(ctx context.Context, req ChatRequest) (*CompletionResponse, error) {
	model := req.Model
	if model == "" {
		model = p.model
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	// Convert messages
	messages := make([]claudeMessage, 0, len(req.Messages))
	var system string
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			system = msg.Content
			continue
		}
		messages = append(messages, claudeMessage(msg))
	}

	claudeReq := claudeRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  messages,
		System:    system,
	}

	if req.Temperature > 0 {
		claudeReq.Temperature = &req.Temperature
	}

	if len(req.Stop) > 0 {
		claudeReq.StopSequences = req.Stop
	}

	body, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", claudeVersion)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(respBody, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	content := ""
	for _, block := range claudeResp.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &CompletionResponse{
		Content:      content,
		FinishReason: claudeResp.StopReason,
		TokensUsed:   claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
	}, nil
}

// Claude API types
type claudeRequest struct {
	Model         string          `json:"model"`
	MaxTokens     int             `json:"max_tokens"`
	Messages      []claudeMessage `json:"messages"`
	System        string          `json:"system,omitempty"`
	Temperature   *float64        `json:"temperature,omitempty"`
	StopSequences []string        `json:"stop_sequences,omitempty"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Role       string         `json:"role"`
	Content    []contentBlock `json:"content"`
	Model      string         `json:"model"`
	StopReason string         `json:"stop_reason"`
	Usage      claudeUsage    `json:"usage"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type claudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
