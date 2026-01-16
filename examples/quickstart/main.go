// Squaremind Quick Start Example
//
// This example demonstrates how to create a basic collective
// with squaremind agents and submit tasks.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/square-mind/squaremind/pkg/agent"
	"github.com/square-mind/squaremind/pkg/collective"
	"github.com/square-mind/squaremind/pkg/identity"
	"github.com/square-mind/squaremind/pkg/llm"
)

func main() {
	// Create a context for clean shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Create a collective
	fmt.Println("Creating collective: DevSwarm")
	c := collective.NewCollective("DevSwarm", collective.CollectiveConfig{
		MinAgents:          2,
		MaxAgents:          10,
		ConsensusThreshold: 0.67,
		ReputationDecay:    0.01,
	})

	// Create LLM provider (uses simulated responses if no API key set)
	var provider llm.Provider
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		provider = llm.NewClaudeProvider(apiKey)
		fmt.Println("Using Claude API")
	} else if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		provider = llm.NewOpenAIProvider(apiKey)
		fmt.Println("Using OpenAI API")
	} else {
		fmt.Println("No API key set - using simulated responses")
		fmt.Println("Set ANTHROPIC_API_KEY or OPENAI_API_KEY for real LLM responses")
	}

	// Spawn agents
	fmt.Println("Spawning agents...")

	coder, err := agent.NewAgent(agent.AgentConfig{
		Name:         "Coder",
		Capabilities: []identity.CapabilityType{identity.CapCodeWrite, identity.CapCodeReview},
		Provider:     provider,
		Model:        string(llm.DefaultModel),
	})
	if err != nil {
		panic(err)
	}

	reviewer, err := agent.NewAgent(agent.AgentConfig{
		Name:         "Reviewer",
		Capabilities: []identity.CapabilityType{identity.CapCodeReview, identity.CapSecurity},
		Provider:     provider,
		Model:        string(llm.DefaultModel),
	})
	if err != nil {
		panic(err)
	}

	architect, err := agent.NewAgent(agent.AgentConfig{
		Name:         "Architect",
		Capabilities: []identity.CapabilityType{identity.CapArchitecture, identity.CapDocumentation},
		Provider:     provider,
		Model:        string(llm.DefaultModel),
	})
	if err != nil {
		panic(err)
	}

	// Join the collective
	fmt.Println("Agents joining collective...")
	c.Join(coder)
	c.Join(reviewer)
	c.Join(architect)

	// Start the collective
	fmt.Println("Starting collective...")
	if err := c.Start(ctx); err != nil {
		panic(err)
	}

	// Display collective status
	stats := c.Stats()
	fmt.Printf("\nCollective Status:\n")
	fmt.Printf("  Name: %s\n", stats.Name)
	fmt.Printf("  Agents: %d\n", stats.AgentCount)
	fmt.Printf("  Avg Reputation: %.1f\n\n", stats.AvgReputation)

	// List agents
	fmt.Println("Agents in collective:")
	for _, a := range c.GetAgents() {
		fmt.Printf("  - %s (SID: %s) - State: %s, Rep: %.1f\n",
			a.Identity.Name,
			a.Identity.SIDShort(),
			a.GetState(),
			a.Reputation.Overall,
		)
	}

	// Create and submit a task
	fmt.Println("\nSubmitting task...")
	task := agent.NewTask(
		"Implement a function to validate email addresses using regex",
		[]identity.CapabilityType{identity.CapCodeWrite},
	)
	task.Complexity = "medium"
	task.Deadline = time.Now().Add(time.Hour)
	task.Reward = 10

	result, err := c.Submit(task)
	if err != nil {
		fmt.Printf("Task failed: %v\n", err)
	} else {
		fmt.Printf("\nTask completed!\n")
		fmt.Printf("  Status: %s\n", result.Status)
		fmt.Printf("  Quality: %.2f\n", result.Quality)
		fmt.Printf("  Duration: %v\n", result.Duration)
		fmt.Printf("  Output: %s\n", result.Output)
	}

	// Final status
	stats = c.Stats()
	fmt.Printf("\nFinal Status:\n")
	fmt.Printf("  Tasks Completed: %d\n", stats.CompletedTasks)

	// Stop the collective
	fmt.Println("\nStopping collective...")
	c.Stop()
	fmt.Println("Done!")
}
