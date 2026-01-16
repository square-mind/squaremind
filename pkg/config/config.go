package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	AnthropicAPIKey string `yaml:"anthropic_api_key"`
	OpenAIAPIKey    string `yaml:"openai_api_key"`
	DefaultModel    string `yaml:"default_model"`
}

// DefaultConfigPath returns the default config file path
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".squaremind", "config.yaml")
}

// Load reads configuration from the config file
func Load() (*Config, error) {
	return LoadFromPath(DefaultConfigPath())
}

// LoadFromPath reads configuration from a specific path
func LoadFromPath(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Return empty config if file doesn't exist
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save writes configuration to the config file
func (c *Config) Save() error {
	return c.SaveToPath(DefaultConfigPath())
}

// SaveToPath writes configuration to a specific path
func (c *Config) SaveToPath(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// GetAnthropicKey returns the Anthropic API key with priority:
// 1. Environment variable ANTHROPIC_API_KEY
// 2. Config file
func (c *Config) GetAnthropicKey() string {
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return key
	}
	return c.AnthropicAPIKey
}

// GetOpenAIKey returns the OpenAI API key with priority:
// 1. Environment variable OPENAI_API_KEY
// 2. Config file
func (c *Config) GetOpenAIKey() string {
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return key
	}
	return c.OpenAIAPIKey
}
