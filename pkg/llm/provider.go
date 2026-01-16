package llm

import "context"

// Provider defines the interface for LLM backends
type Provider interface {
	// Complete generates a completion for the given request
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

	// Name returns the provider name
	Name() string
}

// CompletionRequest represents a completion request
type CompletionRequest struct {
	Model       string            `json:"model"`
	Prompt      string            `json:"prompt"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
	Stop        []string          `json:"stop,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	System      string            `json:"system,omitempty"`
}

// CompletionResponse represents a completion response
type CompletionResponse struct {
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason"`
	TokensUsed   int    `json:"tokens_used"`
}

// ModelType represents supported model types
type ModelType string

const (
	// Claude models (primary)
	ModelClaude35Sonnet ModelType = "claude-3-5-sonnet-20241022" // Default, recommended
	ModelClaude3Opus    ModelType = "claude-3-opus-20240229"
	ModelClaude3Sonnet  ModelType = "claude-3-sonnet-20240229"
	ModelClaude3Haiku   ModelType = "claude-3-haiku-20240307"

	// OpenAI models (fallback)
	ModelGPT4      ModelType = "gpt-4"
	ModelGPT4Turbo ModelType = "gpt-4-turbo"
	ModelGPT35     ModelType = "gpt-3.5-turbo"

	// Default model
	DefaultModel ModelType = ModelClaude35Sonnet
)

// ProviderConfig holds common provider configuration
type ProviderConfig struct {
	APIKey      string
	BaseURL     string
	Timeout     int // seconds
	MaxRetries  int
	DefaultModel string
}

// DefaultProviderConfig returns default configuration
func DefaultProviderConfig() ProviderConfig {
	return ProviderConfig{
		Timeout:    60,
		MaxRetries: 3,
	}
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"` // "user", "assistant", "system"
	Content string `json:"content"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stop        []string  `json:"stop,omitempty"`
}

// ChatProvider extends Provider with chat capabilities
type ChatProvider interface {
	Provider
	Chat(ctx context.Context, req ChatRequest) (*CompletionResponse, error)
}
