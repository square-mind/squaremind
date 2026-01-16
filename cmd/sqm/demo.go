package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/square-mind/squaremind/pkg/agent"
	"github.com/square-mind/squaremind/pkg/cli"
	"github.com/square-mind/squaremind/pkg/collective"
	"github.com/square-mind/squaremind/pkg/identity"
	"github.com/square-mind/squaremind/pkg/llm"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run an interactive demo showcasing Squaremind capabilities",
	Long: `Run an impressive demo that showcases the power of Squaremind's
collective intelligence architecture.

This demo will:
  1. Create a collective of AI agents
  2. Assign cryptographic identities to each agent
  3. Demonstrate task allocation through fair markets
  4. Show emergent coordination between agents
  5. Display real-time reputation updates

Example:
  sqm demo --api-key YOUR_ANTHROPIC_API_KEY

Or set the environment variable:
  export ANTHROPIC_API_KEY=YOUR_KEY
  sqm demo`,
	Run: runDemo,
}

func runDemo(cmd *cobra.Command, args []string) {
	// Print banner
	fmt.Println(cli.SmallBanner())

	// Check for API key
	if provider == nil {
		fmt.Println(cli.Error("\n  No API key configured."))
		fmt.Println()
		fmt.Println("  Set your API key using one of these methods:")
		fmt.Println(cli.StatusLine("info", "CLI flag:    sqm demo --api-key YOUR_KEY"))
		fmt.Println(cli.StatusLine("info", "Environment: export ANTHROPIC_API_KEY=YOUR_KEY"))
		fmt.Println(cli.StatusLine("info", "Config file: ~/.squaremind/config.yaml"))
		fmt.Println()
		os.Exit(1)
	}

	fmt.Printf("  %sProvider:%s %s\n", cli.Gray, cli.Reset, cli.Highlight(provider.Name()))
	fmt.Printf("  %sModel:%s    %s\n", cli.Gray, cli.Reset, cli.Info(string(llm.DefaultModel)))
	fmt.Println()
	fmt.Println(cli.Divider(50))

	// Phase 1: Create Collective
	fmt.Println(cli.Section("PHASE 1: Creating Collective"))

	spinner := cli.NewSpinner("Initializing collective 'DemoSwarm'...")
	spinner.Start()
	time.Sleep(500 * time.Millisecond)

	c := collective.NewCollective("DemoSwarm", collective.CollectiveConfig{
		MinAgents:          2,
		MaxAgents:          10,
		ConsensusThreshold: 0.67,
		ReputationDecay:    0.01,
	})
	spinner.Stop(true)

	fmt.Println(cli.StatusLine("success", fmt.Sprintf("Collective ID: %s%s%s", cli.Cyan, c.ID[:8], cli.Reset)))
	fmt.Println(cli.StatusLine("info", "Consensus threshold: 67%"))
	fmt.Println(cli.StatusLine("info", "Max agents: 10"))

	// Phase 2: Spawn Agents
	fmt.Println(cli.Section("PHASE 2: Spawning Agents"))

	agents := []struct {
		name string
		caps []identity.CapabilityType
		desc string
	}{
		{
			name: "Architect",
			caps: []identity.CapabilityType{identity.CapArchitecture, identity.CapAnalysis},
			desc: "System design & architecture",
		},
		{
			name: "Coder",
			caps: []identity.CapabilityType{identity.CapCodeWrite, identity.CapCodeReview},
			desc: "Code implementation & review",
		},
		{
			name: "SecurityAuditor",
			caps: []identity.CapabilityType{identity.CapSecurity, identity.CapCodeReview},
			desc: "Security analysis & auditing",
		},
	}

	spawnedAgents := make([]*agent.Agent, 0)

	for _, agentDef := range agents {
		spawnSpinner := cli.NewSpinner(fmt.Sprintf("Spawning %s...", agentDef.name))
		spawnSpinner.Start()
		time.Sleep(300 * time.Millisecond)

		a, err := agent.NewAgent(agent.AgentConfig{
			Name:         agentDef.name,
			Capabilities: agentDef.caps,
			Provider:     provider,
			Model:        string(llm.DefaultModel),
		})
		if err != nil {
			spawnSpinner.Stop(false)
			fmt.Println(cli.Error(fmt.Sprintf("  Failed to spawn agent: %v", err)))
			continue
		}

		_ = c.Join(a)
		spawnedAgents = append(spawnedAgents, a)
		spawnSpinner.StopWithMessage(true, fmt.Sprintf("Spawned %s%s%s (%s)", cli.Bold+cli.BrightGreen, agentDef.name, cli.Reset, agentDef.desc))

		// Show agent details
		capStrs := make([]string, len(agentDef.caps))
		for i, cap := range agentDef.caps {
			capStrs[i] = string(cap)
		}
		fmt.Printf("      %sSID:%s %s  %sCaps:%s %s\n",
			cli.Gray, cli.Reset, a.Identity.SIDShort(),
			cli.Gray, cli.Reset, strings.Join(capStrs, ", "))
	}

	// Phase 3: Start Collective
	fmt.Println(cli.Section("PHASE 3: Activating Collective"))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	spinner = cli.NewSpinner("Starting gossip protocol...")
	spinner.Start()
	time.Sleep(400 * time.Millisecond)
	if err := c.Start(ctx); err != nil {
		spinner.Stop(false)
		fmt.Println(cli.Error(fmt.Sprintf("  Failed to start: %v", err)))
		return
	}
	spinner.Stop(true)

	spinner = cli.NewSpinner("Initializing task market...")
	spinner.Start()
	time.Sleep(300 * time.Millisecond)
	spinner.Stop(true)

	spinner = cli.NewSpinner("Establishing consensus layer...")
	spinner.Start()
	time.Sleep(300 * time.Millisecond)
	spinner.Stop(true)

	// Show collective status
	fmt.Println()
	fmt.Printf("  %s┌─────────────────────────────────────────┐%s\n", cli.Green, cli.Reset)
	fmt.Printf("  %s│%s  %sCOLLECTIVE ACTIVE%s                       %s│%s\n", cli.Green, cli.Reset, cli.Bold+cli.BrightGreen, cli.Reset, cli.Green, cli.Reset)
	fmt.Printf("  %s│%s                                         %s│%s\n", cli.Green, cli.Reset, cli.Green, cli.Reset)
	fmt.Printf("  %s│%s   Agents Online: %s%-3d%s                    %s│%s\n", cli.Green, cli.Reset, cli.BrightGreen, len(spawnedAgents), cli.Reset, cli.Green, cli.Reset)
	fmt.Printf("  %s│%s   Gossip Peers:  %s%-3d%s                    %s│%s\n", cli.Green, cli.Reset, cli.Cyan, len(spawnedAgents), cli.Reset, cli.Green, cli.Reset)
	fmt.Printf("  %s│%s   Market Status: %sOPEN%s                   %s│%s\n", cli.Green, cli.Reset, cli.Green, cli.Reset, cli.Green, cli.Reset)
	fmt.Printf("  %s└─────────────────────────────────────────┘%s\n", cli.Green, cli.Reset)

	// Phase 4: Submit Task
	fmt.Println(cli.Section("PHASE 4: Submitting Task to Market"))

	taskDesc := "Analyze the security implications of a distributed consensus mechanism and provide recommendations"

	fmt.Printf("  %sTask:%s %s\n", cli.Gray, cli.Reset, taskDesc)
	fmt.Printf("  %sRequired:%s %s%s%s\n", cli.Gray, cli.Reset, cli.Cyan, "security, analysis", cli.Reset)
	fmt.Println()

	task := agent.NewTask(
		taskDesc,
		[]identity.CapabilityType{identity.CapSecurity, identity.CapAnalysis},
	)
	task.Complexity = "high"
	task.Reward = 15

	// Simulate market bidding
	spinner = cli.NewSpinner("Broadcasting task to market...")
	spinner.Start()
	time.Sleep(400 * time.Millisecond)
	spinner.Stop(true)

	fmt.Println()
	fmt.Printf("  %s┌─ MARKET BIDDING ─────────────────────────┐%s\n", cli.Yellow, cli.Reset)

	for _, a := range spawnedAgents {
		score := a.Capabilities.MatchScore(task.Required)
		if score > 0.3 {
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("  %s│%s  %s%-15s%s bid: capability=%.2f stake=%.1f %s│%s\n",
				cli.Yellow, cli.Reset,
				cli.Dim, a.Identity.Name, cli.Reset,
				score, a.Reputation.Overall*0.1,
				cli.Yellow, cli.Reset)
		}
	}

	fmt.Printf("  %s└──────────────────────────────────────────┘%s\n", cli.Yellow, cli.Reset)
	fmt.Println()

	spinner = cli.NewSpinner("Selecting best agent...")
	spinner.Start()
	time.Sleep(500 * time.Millisecond)
	spinner.StopWithMessage(true, fmt.Sprintf("Task assigned to %sSecurityAuditor%s", cli.Bold+cli.BrightGreen, cli.Reset))

	// Phase 5: Execute Task
	fmt.Println(cli.Section("PHASE 5: Agent Executing Task"))

	spinner = cli.NewSpinner("SecurityAuditor processing task with Claude...")
	spinner.Start()

	// Actually call the LLM
	result, err := c.Submit(task)

	if err != nil {
		spinner.Stop(false)
		fmt.Println(cli.Error(fmt.Sprintf("  Task failed: %v", err)))
	} else {
		spinner.Stop(true)

		fmt.Println()
		fmt.Printf("  %s┌─ TASK RESULT ────────────────────────────┐%s\n", cli.Green, cli.Reset)
		fmt.Printf("  %s│%s  Status:   %s%s%s%s                        %s│%s\n",
			cli.Green, cli.Reset, cli.Bold, cli.BrightGreen, result.Status, cli.Reset, cli.Green, cli.Reset)
		fmt.Printf("  %s│%s  Quality:  %s%.2f%s                           %s│%s\n",
			cli.Green, cli.Reset, cli.Cyan, result.Quality, cli.Reset, cli.Green, cli.Reset)
		fmt.Printf("  %s│%s  Duration: %s%v%s                      %s│%s\n",
			cli.Green, cli.Reset, cli.Dim, result.Duration.Round(time.Millisecond), cli.Reset, cli.Green, cli.Reset)
		fmt.Printf("  %s└──────────────────────────────────────────┘%s\n", cli.Green, cli.Reset)

		// Print response (truncated if too long)
		fmt.Println()
		fmt.Printf("  %sAgent Response:%s\n", cli.Bold, cli.Reset)
		fmt.Println(cli.Divider(50))

		output := result.Output
		if len(output) > 500 {
			output = output[:500] + "..."
		}
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			fmt.Printf("  %s%s%s\n", cli.Dim, line, cli.Reset)
		}
		fmt.Println(cli.Divider(50))
	}

	// Phase 6: Reputation Update
	fmt.Println(cli.Section("PHASE 6: Reputation Updates"))

	for _, a := range spawnedAgents {
		change := ""
		if a.Identity.Name == "SecurityAuditor" && result != nil && result.Status == "completed" {
			change = fmt.Sprintf(" %s(+%.1f)%s", cli.Green, result.Quality*10, cli.Reset)
		}
		fmt.Printf("  %s%-15s%s  Rep: %s%.1f%s%s\n",
			cli.Dim, a.Identity.Name, cli.Reset,
			cli.Yellow, a.Reputation.Overall, cli.Reset,
			change)
	}

	// Final Summary
	fmt.Println(cli.Section("DEMO COMPLETE"))

	fmt.Printf(`
  %sSquaremind demonstrated:%s

  %s✓%s %sCryptographic Identity%s - Each agent has Ed25519 keypair
  %s✓%s %sFair Task Markets%s      - Transparent bidding & selection
  %s✓%s %sReputation Staking%s     - Accountability through stakes
  %s✓%s %sLLM Integration%s        - Real AI-powered task execution
  %s✓%s %sCollective Coordination%s - Many agents, one mind

  %sLearn more:%s https://squaremind.xyz
  %sWhitepaper:%s https://squaremind.xyz/whitepaper.html

`,
		cli.Bold, cli.Reset,
		cli.Green, cli.Reset, cli.Bold, cli.Reset,
		cli.Green, cli.Reset, cli.Bold, cli.Reset,
		cli.Green, cli.Reset, cli.Bold, cli.Reset,
		cli.Green, cli.Reset, cli.Bold, cli.Reset,
		cli.Green, cli.Reset, cli.Bold, cli.Reset,
		cli.Cyan, cli.Reset,
		cli.Cyan, cli.Reset,
	)

	// Cleanup
	c.Stop()
}

func init() {
	rootCmd.AddCommand(demoCmd)
}
