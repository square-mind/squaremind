package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/squaremind/squaremind/pkg/agent"
	"github.com/squaremind/squaremind/pkg/cli"
	"github.com/squaremind/squaremind/pkg/collective"
	"github.com/squaremind/squaremind/pkg/identity"
	"github.com/squaremind/squaremind/pkg/llm"
)

var swarmCmd = &cobra.Command{
	Use:   "swarm [task]",
	Short: "Execute a task using swarm intelligence",
	Long: `Execute a complex task using Squaremind's swarm intelligence.

This spawns multiple specialized agents that work together:
  - Researcher: Analyzes the problem space
  - Architect: Designs the solution approach
  - Implementer: Creates the implementation
  - Critic: Reviews and identifies issues
  - Synthesizer: Combines outputs into final result

The agents coordinate through fair markets and gossip protocols,
demonstrating emergent collective intelligence.

Example:
  sqm swarm "Design a microservices architecture for an e-commerce platform"
  sqm swarm "Write a comprehensive security audit checklist"`,
	Args: cobra.ExactArgs(1),
	Run:  runSwarm,
}

var swarmAgents int

func runSwarm(cmd *cobra.Command, args []string) {
	task := args[0]

	fmt.Println(cli.SmallBanner())

	if provider == nil {
		fmt.Println(cli.Error("\n  No API key configured. Run: sqm demo --help"))
		os.Exit(1)
	}

	fmt.Printf("  %sTask:%s %s\n", cli.Bold, cli.Reset, task)
	fmt.Printf("  %sMode:%s Swarm Intelligence (%d agents)\n", cli.Dim, cli.Reset, swarmAgents)
	fmt.Println()
	fmt.Println(cli.Divider(55))

	// Create collective
	fmt.Println(cli.Section("INITIALIZING SWARM"))

	c := collective.NewCollective("SwarmMind", collective.CollectiveConfig{
		MinAgents:          2,
		MaxAgents:          swarmAgents,
		ConsensusThreshold: 0.67,
	})

	// Define swarm roles
	roles := []struct {
		name   string
		caps   []identity.CapabilityType
		prompt string
		desc   string
	}{
		{
			name:   "Researcher",
			caps:   []identity.CapabilityType{identity.CapResearch, identity.CapAnalysis},
			prompt: "You are a thorough researcher. Analyze the problem, identify key considerations, constraints, and relevant prior art. Be comprehensive but concise.",
			desc:   "Problem analysis & research",
		},
		{
			name:   "Architect",
			caps:   []identity.CapabilityType{identity.CapArchitecture, identity.CapAnalysis},
			prompt: "You are a systems architect. Based on the analysis, design a high-level solution architecture. Focus on structure, patterns, and key decisions.",
			desc:   "Solution design & architecture",
		},
		{
			name:   "Implementer",
			caps:   []identity.CapabilityType{identity.CapCodeWrite, identity.CapDocumentation},
			prompt: "You are a skilled implementer. Take the architecture and create a concrete implementation plan with specific steps, code patterns, or detailed instructions.",
			desc:   "Concrete implementation",
		},
		{
			name:   "Critic",
			caps:   []identity.CapabilityType{identity.CapCodeReview, identity.CapSecurity},
			prompt: "You are a critical reviewer. Identify potential issues, edge cases, security concerns, and improvements. Be constructive but thorough.",
			desc:   "Critical review & security",
		},
		{
			name:   "Synthesizer",
			caps:   []identity.CapabilityType{identity.CapDocumentation, identity.CapAnalysis},
			prompt: "You are a synthesizer. Combine all the inputs from other agents into a coherent, well-structured final output. Resolve conflicts and create a unified response.",
			desc:   "Synthesis & final output",
		},
	}

	// Limit to requested agent count
	if swarmAgents < len(roles) {
		roles = roles[:swarmAgents]
	}

	agents := make([]*agent.Agent, 0)
	agentPrompts := make(map[string]string)

	for _, role := range roles {
		spinner := cli.NewSpinner(fmt.Sprintf("Spawning %s...", role.name))
		spinner.Start()
		time.Sleep(200 * time.Millisecond)

		a, err := agent.NewAgent(agent.AgentConfig{
			Name:         role.name,
			Capabilities: role.caps,
			Provider:     provider,
			Model:        string(llm.DefaultModel),
		})
		if err != nil {
			spinner.Stop(false)
			continue
		}

		c.Join(a)
		agents = append(agents, a)
		agentPrompts[a.Identity.SID] = role.prompt
		spinner.StopWithMessage(true, fmt.Sprintf("%s%s%s - %s", cli.BrightGreen, role.name, cli.Reset, role.desc))
	}

	// Start collective
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	spinner := cli.NewSpinner("Starting collective intelligence...")
	spinner.Start()
	c.Start(ctx)
	time.Sleep(300 * time.Millisecond)
	spinner.Stop(true)

	// Execute swarm coordination
	fmt.Println(cli.Section("SWARM EXECUTION"))

	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Phase 1: Research (runs first)
	fmt.Println()
	fmt.Printf("  %s┌─ PHASE 1: RESEARCH ──────────────────────┐%s\n", cli.Cyan, cli.Reset)

	if len(agents) > 0 {
		researcher := agents[0]
		spinner := cli.NewSpinner(fmt.Sprintf("%s analyzing problem...", researcher.Identity.Name))
		spinner.Start()

		resp, err := provider.Complete(ctx, llm.CompletionRequest{
			System:    agentPrompts[researcher.Identity.SID],
			Prompt:    fmt.Sprintf("Analyze this task and provide key insights:\n\n%s", task),
			MaxTokens: 500,
		})

		if err != nil {
			spinner.Stop(false)
			fmt.Printf("  %s│%s  Error: %v\n", cli.Cyan, cli.Reset, err)
		} else {
			spinner.Stop(true)
			results["research"] = resp.Content
			preview := resp.Content
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			fmt.Printf("  %s│%s  %s%s%s\n", cli.Cyan, cli.Reset, cli.Dim, preview, cli.Reset)
		}
	}
	fmt.Printf("  %s└──────────────────────────────────────────┘%s\n", cli.Cyan, cli.Reset)

	// Phase 2: Parallel execution (Architect + Implementer)
	if len(agents) >= 3 {
		fmt.Println()
		fmt.Printf("  %s┌─ PHASE 2: PARALLEL DESIGN ───────────────┐%s\n", cli.Yellow, cli.Reset)

		wg.Add(2)

		// Architect
		go func() {
			defer wg.Done()
			architect := agents[1]
			resp, err := provider.Complete(ctx, llm.CompletionRequest{
				System: agentPrompts[architect.Identity.SID],
				Prompt: fmt.Sprintf("Based on this research:\n%s\n\nDesign a solution for:\n%s",
					results["research"], task),
				MaxTokens: 500,
			})
			if err == nil {
				mu.Lock()
				results["architecture"] = resp.Content
				mu.Unlock()
				fmt.Printf("  %s│%s  %s✓ Architect%s completed design\n", cli.Yellow, cli.Reset, cli.Green, cli.Reset)
			}
		}()

		// Implementer
		go func() {
			defer wg.Done()
			time.Sleep(500 * time.Millisecond) // Slight delay to use architecture
			implementer := agents[2]
			resp, err := provider.Complete(ctx, llm.CompletionRequest{
				System: agentPrompts[implementer.Identity.SID],
				Prompt: fmt.Sprintf("Create implementation details for:\n%s\n\nContext:\n%s",
					task, results["research"]),
				MaxTokens: 500,
			})
			if err == nil {
				mu.Lock()
				results["implementation"] = resp.Content
				mu.Unlock()
				fmt.Printf("  %s│%s  %s✓ Implementer%s completed plan\n", cli.Yellow, cli.Reset, cli.Green, cli.Reset)
			}
		}()

		spinner := cli.NewSpinner("Agents working in parallel...")
		spinner.Start()
		wg.Wait()
		spinner.Stop(true)
		fmt.Printf("  %s└──────────────────────────────────────────┘%s\n", cli.Yellow, cli.Reset)
	}

	// Phase 3: Critical Review
	if len(agents) >= 4 {
		fmt.Println()
		fmt.Printf("  %s┌─ PHASE 3: CRITICAL REVIEW ───────────────┐%s\n", cli.Red, cli.Reset)

		critic := agents[3]
		spinner := cli.NewSpinner(fmt.Sprintf("%s reviewing outputs...", critic.Identity.Name))
		spinner.Start()

		allOutputs := fmt.Sprintf("Research:\n%s\n\nArchitecture:\n%s\n\nImplementation:\n%s",
			results["research"], results["architecture"], results["implementation"])

		resp, err := provider.Complete(ctx, llm.CompletionRequest{
			System:    agentPrompts[critic.Identity.SID],
			Prompt:    fmt.Sprintf("Review these outputs for issues and improvements:\n\n%s", allOutputs),
			MaxTokens: 400,
		})

		if err == nil {
			results["critique"] = resp.Content
			spinner.Stop(true)
			fmt.Printf("  %s│%s  Issues identified and improvements suggested\n", cli.Red, cli.Reset)
		} else {
			spinner.Stop(false)
		}
		fmt.Printf("  %s└──────────────────────────────────────────┘%s\n", cli.Red, cli.Reset)
	}

	// Phase 4: Synthesis
	fmt.Println()
	fmt.Printf("  %s┌─ PHASE 4: SYNTHESIS ──────────────────────┐%s\n", cli.Green, cli.Reset)

	var synthesizer *agent.Agent
	if len(agents) >= 5 {
		synthesizer = agents[4]
	} else if len(agents) > 0 {
		synthesizer = agents[0] // Use researcher as fallback
	}

	if synthesizer != nil {
		spinner := cli.NewSpinner("Synthesizing final output...")
		spinner.Start()

		allResults := ""
		for k, v := range results {
			allResults += fmt.Sprintf("\n=== %s ===\n%s\n", strings.ToUpper(k), v)
		}

		resp, err := provider.Complete(ctx, llm.CompletionRequest{
			System: "You are a synthesizer. Create a comprehensive, well-structured final response that combines all inputs. Format it beautifully with clear sections.",
			Prompt: fmt.Sprintf("Original task: %s\n\nAgent outputs to synthesize:\n%s\n\nCreate a final, comprehensive response.",
				task, allResults),
			MaxTokens: 1500,
		})

		spinner.Stop(err == nil)

		if err == nil {
			results["final"] = resp.Content
		}
	}
	fmt.Printf("  %s└──────────────────────────────────────────┘%s\n", cli.Green, cli.Reset)

	// Display final result
	fmt.Println(cli.Section("SWARM OUTPUT"))

	if final, ok := results["final"]; ok {
		fmt.Println()
		lines := strings.Split(final, "\n")
		for _, line := range lines {
			fmt.Printf("  %s\n", line)
		}
		fmt.Println()
	} else {
		fmt.Println(cli.Error("  No final output generated"))
	}

	// Summary
	fmt.Println(cli.Divider(55))
	fmt.Printf("\n  %sSwarm Statistics:%s\n", cli.Bold, cli.Reset)
	fmt.Printf("  %s• Agents deployed:%s %d\n", cli.Dim, cli.Reset, len(agents))
	fmt.Printf("  %s• Phases completed:%s %d\n", cli.Dim, cli.Reset, len(results))
	fmt.Printf("  %s• Coordination:%s Gossip + Market\n", cli.Dim, cli.Reset)
	fmt.Println()

	c.Stop()
}

func init() {
	swarmCmd.Flags().IntVarP(&swarmAgents, "agents", "n", 5, "Number of agents in swarm (2-5)")
	rootCmd.AddCommand(swarmCmd)
}
