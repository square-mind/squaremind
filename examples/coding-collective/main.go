// Coding Collective Example
//
// This example demonstrates a specialized collective for software development
// with agents that have complementary capabilities.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/squaremind/squaremind/pkg/agent"
	"github.com/squaremind/squaremind/pkg/collective"
	"github.com/squaremind/squaremind/pkg/identity"
	"github.com/squaremind/squaremind/pkg/llm"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down collective...")
		cancel()
	}()

	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println(" SQUAREMIND CODING COLLECTIVE")
	fmt.Println("=" + string(make([]byte, 50)))

	// Create the collective
	c := collective.NewCollective("CodingCollective", collective.CollectiveConfig{
		MinAgents:          5,
		MaxAgents:          20,
		ConsensusThreshold: 0.67,
	})

	// Create LLM provider (uses simulated responses if no API key set)
	var provider llm.Provider
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		provider = llm.NewClaudeProvider(apiKey)
		fmt.Println("Using Claude API\n")
	} else if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		provider = llm.NewOpenAIProvider(apiKey)
		fmt.Println("Using OpenAI API\n")
	} else {
		fmt.Println("No API key set - using simulated responses")
		fmt.Println("Set ANTHROPIC_API_KEY or OPENAI_API_KEY for real LLM responses\n")
	}

	// Spawn specialized agents
	agents := []struct {
		name string
		caps []identity.CapabilityType
		role string
	}{
		{
			name: "SeniorDev",
			caps: []identity.CapabilityType{identity.CapCodeWrite, identity.CapCodeReview, identity.CapArchitecture},
			role: "Senior Developer - writes complex code, reviews, and designs architecture",
		},
		{
			name: "JuniorDev",
			caps: []identity.CapabilityType{identity.CapCodeWrite, identity.CapTesting},
			role: "Junior Developer - writes code and tests",
		},
		{
			name: "SecurityEngineer",
			caps: []identity.CapabilityType{identity.CapSecurity, identity.CapCodeReview},
			role: "Security Engineer - audits code for vulnerabilities",
		},
		{
			name: "TechWriter",
			caps: []identity.CapabilityType{identity.CapDocumentation, identity.CapAnalysis},
			role: "Technical Writer - creates documentation",
		},
		{
			name: "QAEngineer",
			caps: []identity.CapabilityType{identity.CapTesting, identity.CapCodeReview},
			role: "QA Engineer - writes and runs tests",
		},
		{
			name: "Researcher",
			caps: []identity.CapabilityType{identity.CapResearch, identity.CapAnalysis},
			role: "Researcher - investigates solutions and analyzes options",
		},
	}

	fmt.Println("\nSpawning agents...\n")
	for _, agentDef := range agents {
		a, err := agent.NewAgent(agent.AgentConfig{
			Name:         agentDef.name,
			Capabilities: agentDef.caps,
			Provider:     provider,
			Model:        string(llm.DefaultModel),
		})
		if err != nil {
			panic(err)
		}
		c.Join(a)
		fmt.Printf("  [+] %s\n", agentDef.name)
		fmt.Printf("      Role: %s\n", agentDef.role)
		fmt.Printf("      SID: %s\n\n", a.Identity.SIDShort())
	}

	// Start the collective
	fmt.Println("Starting collective...")
	if err := c.Start(ctx); err != nil {
		panic(err)
	}

	// Show collective status
	printStatus(c)

	// Define a set of tasks
	tasks := []struct {
		description string
		complexity  string
		required    []identity.CapabilityType
	}{
		{
			description: "Design the authentication system architecture",
			complexity:  "high",
			required:    []identity.CapabilityType{identity.CapArchitecture},
		},
		{
			description: "Implement JWT token generation and validation",
			complexity:  "medium",
			required:    []identity.CapabilityType{identity.CapCodeWrite},
		},
		{
			description: "Review the authentication code for security vulnerabilities",
			complexity:  "medium",
			required:    []identity.CapabilityType{identity.CapSecurity, identity.CapCodeReview},
		},
		{
			description: "Write unit tests for the authentication module",
			complexity:  "medium",
			required:    []identity.CapabilityType{identity.CapTesting},
		},
		{
			description: "Create API documentation for the auth endpoints",
			complexity:  "low",
			required:    []identity.CapabilityType{identity.CapDocumentation},
		},
	}

	// Submit tasks
	fmt.Println("\n" + string(make([]byte, 50)))
	fmt.Println(" EXECUTING TASKS")
	fmt.Println(string(make([]byte, 50)) + "\n")

	for i, taskDef := range tasks {
		fmt.Printf("[Task %d/%d] %s\n", i+1, len(tasks), taskDef.description)
		fmt.Printf("  Complexity: %s\n", taskDef.complexity)
		fmt.Printf("  Required: %v\n", taskDef.required)

		task := agent.NewTask(taskDef.description, taskDef.required)
		task.Complexity = taskDef.complexity
		task.Deadline = time.Now().Add(time.Hour)

		result, err := c.Submit(task)
		if err != nil {
			fmt.Printf("  Status: FAILED - %v\n\n", err)
		} else {
			fmt.Printf("  Status: %s\n", result.Status)
			fmt.Printf("  Assigned to: %s\n", result.AgentSID[:8])
			fmt.Printf("  Quality: %.2f\n", result.Quality)
			fmt.Printf("  Duration: %v\n\n", result.Duration)
		}
	}

	// Final status
	printStatus(c)

	// Print agent reputations
	fmt.Println("\nAgent Reputations:")
	for _, a := range c.GetAgents() {
		fmt.Printf("  %s: %.1f (completed: %d, failed: %d)\n",
			a.Identity.Name,
			a.Reputation.Overall,
			a.Reputation.TasksCompleted,
			a.Reputation.TasksFailed,
		)
	}

	// Cleanup
	c.Stop()
	fmt.Println("\nCollective stopped. Goodbye!")
}

func printStatus(c *collective.Collective) {
	stats := c.Stats()
	fmt.Println("\n" + string(make([]byte, 50)))
	fmt.Println(" COLLECTIVE STATUS")
	fmt.Println(string(make([]byte, 50)))
	fmt.Printf("  Name: %s\n", stats.Name)
	fmt.Printf("  Agents: %d\n", stats.AgentCount)
	fmt.Printf("  Pending Tasks: %d\n", stats.PendingTasks)
	fmt.Printf("  Active Tasks: %d\n", stats.ActiveTasks)
	fmt.Printf("  Completed Tasks: %d\n", stats.CompletedTasks)
	fmt.Printf("  Average Reputation: %.1f\n", stats.AvgReputation)
	fmt.Println(string(make([]byte, 50)) + "\n")
}
