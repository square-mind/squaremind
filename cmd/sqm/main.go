package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/square-mind/squaremind/pkg/agent"
	"github.com/square-mind/squaremind/pkg/collective"
	"github.com/square-mind/squaremind/pkg/config"
	"github.com/square-mind/squaremind/pkg/identity"
	"github.com/square-mind/squaremind/pkg/llm"
)

var (
	version = "0.1.0"

	// Global flags
	apiKey string

	// Global state for CLI session
	activeCollective *collective.Collective
	provider         llm.Provider
	cfg              *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "sqm",
	Short: "Squaremind - Many Agents. One Mind.",
	Long: `Squaremind is a next-generation multi-agent orchestration protocol
enabling truly autonomous AI collectives.

Many Agents. One Mind. Fair coordination. Transparent markets.

Learn more: https://squaremind.xyz`,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load config file
		var err error
		cfg, err = config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
			cfg = &config.Config{}
		}

		// Initialize provider with priority: CLI flag > env var > config file
		key := apiKey
		if key == "" {
			key = cfg.GetAnthropicKey()
		}
		if key != "" {
			provider = llm.NewClaudeProvider(key)
			return
		}

		// Fallback to OpenAI if no Anthropic key
		openaiKey := cfg.GetOpenAIKey()
		if openaiKey != "" {
			provider = llm.NewOpenAIProvider(openaiKey)
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new collective",
	Long:  `Initialize a new squaremind collective with the given name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		maxAgents, _ := cmd.Flags().GetInt("max-agents")
		threshold, _ := cmd.Flags().GetFloat64("threshold")

		cfg := collective.CollectiveConfig{
			MinAgents:          2,
			MaxAgents:          maxAgents,
			ConsensusThreshold: threshold,
			ReputationDecay:    0.01,
		}

		c := collective.NewCollective(name, cfg)
		activeCollective = c

		fmt.Printf("\n  Collective '%s' initialized\n\n", name)
		fmt.Printf("  ID: %s\n", c.ID)
		fmt.Printf("  Max Agents: %d\n", cfg.MaxAgents)
		fmt.Printf("  Consensus Threshold: %.0f%%\n\n", cfg.ConsensusThreshold*100)
	},
}

var spawnCmd = &cobra.Command{
	Use:   "spawn [name]",
	Short: "Spawn a new squaremind agent",
	Long:  `Spawn a new autonomous agent with the specified capabilities.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		caps, _ := cmd.Flags().GetStringSlice("capabilities")
		model, _ := cmd.Flags().GetString("model")

		// Convert string capabilities to types
		capTypes := make([]identity.CapabilityType, len(caps))
		for i, c := range caps {
			capTypes[i] = identity.CapabilityType(c)
		}

		cfg := agent.AgentConfig{
			Name:         name,
			Capabilities: capTypes,
			Model:        model,
			Provider:     provider,
		}

		a, err := agent.NewAgent(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error spawning agent: %v\n", err)
			os.Exit(1)
		}

		// Join collective if one is active
		if activeCollective != nil {
			if err := activeCollective.Join(a); err != nil {
				fmt.Fprintf(os.Stderr, "Error joining collective: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Printf("\n  Squaremind agent '%s' spawned\n\n", name)
		fmt.Printf("  SID: %s\n", a.Identity.SID)
		fmt.Printf("  Public Key: %s...\n", a.Identity.PublicKeyHex()[:16])
		fmt.Printf("  Capabilities: %v\n", caps)
		fmt.Printf("  Model: %s\n", model)
		fmt.Printf("  Reputation: %.1f\n\n", a.Reputation.Overall)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the collective",
	Long:  `Start all agents in the collective and begin autonomous operation.`,
	Run: func(cmd *cobra.Command, args []string) {
		if activeCollective == nil {
			fmt.Fprintln(os.Stderr, "No collective initialized. Run 'sqm init <name>' first.")
			os.Exit(1)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle shutdown signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		if err := activeCollective.Start(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting collective: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n  Starting Squaremind collective...")
		fmt.Printf("  Name: %s\n", activeCollective.Name)
		fmt.Printf("  Agents: %d\n", activeCollective.Size())
		fmt.Println("  Press Ctrl+C to stop")

		// Wait for shutdown
		<-sigChan
		fmt.Println("\n  Shutting down...")
		activeCollective.Stop()
		cancel()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show collective status",
	Long:  `Display the current status of the collective and its agents.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n  Collective Status")
		fmt.Println("  " + strings.Repeat("-", 40))

		if activeCollective == nil {
			fmt.Println("  No active collective")
			fmt.Println("  Run 'sqm init <name>' to create one")
			fmt.Println()
			return
		}

		stats := activeCollective.Stats()
		fmt.Printf("  Name: %s\n", stats.Name)
		fmt.Printf("  Agents: %d\n", stats.AgentCount)
		fmt.Printf("  Tasks Pending: %d\n", stats.PendingTasks)
		fmt.Printf("  Tasks Active: %d\n", stats.ActiveTasks)
		fmt.Printf("  Tasks Completed: %d\n", stats.CompletedTasks)
		fmt.Printf("  Avg Reputation: %.1f\n", stats.AvgReputation)
		fmt.Println()

		// List agents
		agents := activeCollective.GetAgents()
		if len(agents) > 0 {
			fmt.Println("  Agents:")
			for _, a := range agents {
				state := string(a.GetState())
				caps := a.Capabilities.List()
				capStrs := make([]string, len(caps))
				for i, c := range caps {
					capStrs[i] = string(c)
				}
				fmt.Printf("    - %s (%s) [%s] rep=%.1f\n",
					a.Identity.Name,
					a.Identity.SIDShort(),
					state,
					a.Reputation.Overall,
				)
				fmt.Printf("      capabilities: %s\n", strings.Join(capStrs, ", "))
			}
		}
		fmt.Println()
	},
}

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Task management commands",
}

var taskSubmitCmd = &cobra.Command{
	Use:   "submit [description]",
	Short: "Submit a task to the collective",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if activeCollective == nil {
			fmt.Fprintln(os.Stderr, "No collective initialized.")
			os.Exit(1)
		}

		description := args[0]
		complexity, _ := cmd.Flags().GetString("complexity")
		capsStr, _ := cmd.Flags().GetStringSlice("requires")
		reward, _ := cmd.Flags().GetFloat64("reward")
		async, _ := cmd.Flags().GetBool("async")

		// Convert capabilities
		caps := make([]identity.CapabilityType, len(capsStr))
		for i, c := range capsStr {
			caps[i] = identity.CapabilityType(c)
		}

		task := agent.NewTask(description, caps)
		task.Complexity = complexity
		task.Reward = reward
		task.Deadline = time.Now().Add(time.Hour)

		fmt.Printf("\n  Submitting task: %s\n", description)
		fmt.Printf("  Task ID: %s\n", task.ID)
		fmt.Printf("  Complexity: %s\n", complexity)
		fmt.Printf("  Required capabilities: %v\n\n", capsStr)

		if async {
			id, err := activeCollective.SubmitAsync(task)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("  Task submitted asynchronously. ID: %s\n\n", id)
		} else {
			result, err := activeCollective.Submit(task)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("  Task completed!\n")
			fmt.Printf("  Status: %s\n", result.Status)
			fmt.Printf("  Quality: %.2f\n", result.Quality)
			fmt.Printf("  Duration: %v\n", result.Duration)
			fmt.Printf("  Output:\n%s\n\n", result.Output)
		}
	},
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Agent management commands",
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agents",
	Run: func(cmd *cobra.Command, args []string) {
		if activeCollective == nil {
			fmt.Fprintln(os.Stderr, "No collective initialized.")
			os.Exit(1)
		}

		agents := activeCollective.GetAgents()
		fmt.Printf("\n  Agents (%d total):\n\n", len(agents))

		for _, a := range agents {
			caps := a.Capabilities.List()
			capStrs := make([]string, len(caps))
			for i, c := range caps {
				capStrs[i] = string(c)
			}
			fmt.Printf("  Name: %s\n", a.Identity.Name)
			fmt.Printf("  SID: %s\n", a.Identity.SID)
			fmt.Printf("  State: %s\n", a.GetState())
			fmt.Printf("  Reputation: %.1f\n", a.Reputation.Overall)
			fmt.Printf("  Capabilities: %s\n", strings.Join(capStrs, ", "))
			fmt.Printf("  Tasks Completed: %d\n", a.Reputation.TasksCompleted)
			fmt.Println()
		}
	},
}

var agentStopCmd = &cobra.Command{
	Use:   "stop [sid]",
	Short: "Stop an agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if activeCollective == nil {
			fmt.Fprintln(os.Stderr, "No collective initialized.")
			os.Exit(1)
		}

		sid := args[0]
		if err := activeCollective.Leave(sid); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n  Agent %s stopped and removed from collective.\n\n", sid)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		switch key {
		case "api-key":
			provider = llm.NewClaudeProvider(value)
			fmt.Printf("\n  API key configured for Claude.\n\n")
		case "openai-key":
			provider = llm.NewOpenAIProvider(value)
			fmt.Printf("\n  API key configured for OpenAI.\n\n")
		default:
			fmt.Fprintf(os.Stderr, "Unknown config key: %s\n", key)
			os.Exit(1)
		}
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Anthropic API key (overrides env and config)")

	// Init command flags
	initCmd.Flags().IntP("max-agents", "m", 100, "Maximum number of agents")
	initCmd.Flags().Float64P("threshold", "t", 0.67, "Consensus threshold (0.0-1.0)")

	// Spawn command flags
	spawnCmd.Flags().StringSliceP("capabilities", "c", []string{"code.write"}, "Agent capabilities")
	spawnCmd.Flags().StringP("model", "m", string(llm.DefaultModel), "LLM model to use")

	// Task submit flags
	taskSubmitCmd.Flags().StringP("complexity", "x", "medium", "Task complexity (low/medium/high)")
	taskSubmitCmd.Flags().StringSliceP("requires", "r", []string{}, "Required capabilities")
	taskSubmitCmd.Flags().Float64P("reward", "w", 10, "Reputation reward")
	taskSubmitCmd.Flags().BoolP("async", "a", false, "Submit asynchronously")

	// Add subcommands
	taskCmd.AddCommand(taskSubmitCmd)
	agentCmd.AddCommand(agentListCmd)
	agentCmd.AddCommand(agentStopCmd)
	configCmd.AddCommand(configSetCmd)

	// Add all commands to root
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(spawnCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(agentCmd)
	rootCmd.AddCommand(configCmd)
	// demoCmd is added in demo.go's init()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
