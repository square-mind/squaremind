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
	openaiAPIURL = "https://api.openai.com/v1/chat/completions"
)

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
	orgID      string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: openaiAPIURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		model: string(ModelGPT4),
	}
}

// WithBaseURL sets a custom base URL
func (p *OpenAIProvider) WithBaseURL(url string) *OpenAIProvider {
	p.baseURL = url
	return p
}

// WithModel sets the default model
func (p *OpenAIProvider) WithModel(model string) *OpenAIProvider {
	p.model = model
	return p
}

// WithOrgID sets the organization ID
func (p *OpenAIProvider) WithOrgID(orgID string) *OpenAIProvider {
	p.orgID = orgID
	return p
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Complete generates a completion
func (p *OpenAIProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	// Convert to chat format
	messages := []openaiMessage{
		{
			Role:    "user",
			Content: req.Prompt,
		},
	}

	if req.System != "" {
		messages = append([]openaiMessage{
			{Role: "system", Content: req.System},
		}, messages...)
	}

	return p.doRequest(ctx, req.Model, messages, req.MaxTokens, req.Temperature, req.Stop)
}

// Chat implements chat completion
func (p *OpenAIProvider) Chat(ctx context.Context, req ChatRequest) (*CompletionResponse, error) {
	messages := make([]openaiMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = openaiMessage(m)
	}

	return p.doRequest(ctx, req.Model, messages, req.MaxTokens, req.Temperature, req.Stop)
}

func (p *OpenAIProvider) doRequest(ctx context.Context, model string, messages []openaiMessage, maxTokens int, temperature float64, stop []string) (*CompletionResponse, error) {
	if model == "" {
		model = p.model
	}

	if maxTokens == 0 {
		maxTokens = 4096
	}

	openaiReq := openaiRequest{
		Model:    model,
		Messages: messages,
	}

	if maxTokens > 0 {
		openaiReq.MaxTokens = &maxTokens
	}

	if temperature > 0 {
		openaiReq.Temperature = &temperature
	}

	if len(stop) > 0 {
		openaiReq.Stop = stop
	}

	body, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	if p.orgID != "" {
		httpReq.Header.Set("OpenAI-Organization", p.orgID)
	}

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

	var openaiResp openaiResponse
	if err := json.Unmarshal(respBody, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := openaiResp.Choices[0]

	return &CompletionResponse{
		Content:      choice.Message.Content,
		FinishReason: choice.FinishReason,
		TokensUsed:   openaiResp.Usage.TotalTokens,
	}, nil
}

// OpenAI API types
type openaiRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	MaxTokens   *int            `json:"max_tokens,omitempty"`
	Temperature *float64        `json:"temperature,omitempty"`
	Stop        []string        `json:"stop,omitempty"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []openaiChoice `json:"choices"`
	Usage   openaiUsage    `json:"usage"`
}

type openaiChoice struct {
	Index        int           `json:"index"`
	Message      openaiMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

type openaiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
