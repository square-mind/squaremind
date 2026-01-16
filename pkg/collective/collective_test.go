package collective

import (
	"context"
	"testing"
	"time"

	"github.com/square-mind/squaremind/pkg/agent"
	"github.com/square-mind/squaremind/pkg/identity"
)

func TestNewCollective(t *testing.T) {
	cfg := CollectiveConfig{
		MinAgents:          2,
		MaxAgents:          10,
		ConsensusThreshold: 0.67,
		ReputationDecay:    0.01,
	}

	c := NewCollective("TestCollective", cfg)

	if c.Name != "TestCollective" {
		t.Errorf("Expected name 'TestCollective', got '%s'", c.Name)
	}

	if c.ID == "" {
		t.Error("Collective ID should not be empty")
	}

	if c.Size() != 0 {
		t.Errorf("Expected size 0, got %d", c.Size())
	}
}

func TestDefaultCollectiveConfig(t *testing.T) {
	cfg := DefaultCollectiveConfig()

	if cfg.MinAgents != 2 {
		t.Errorf("Expected MinAgents 2, got %d", cfg.MinAgents)
	}

	if cfg.MaxAgents != 100 {
		t.Errorf("Expected MaxAgents 100, got %d", cfg.MaxAgents)
	}

	if cfg.ConsensusThreshold != 0.67 {
		t.Errorf("Expected ConsensusThreshold 0.67, got %f", cfg.ConsensusThreshold)
	}
}

func TestCollective_JoinLeave(t *testing.T) {
	c := NewCollective("TestCollective", DefaultCollectiveConfig())

	a, _ := agent.NewAgent(agent.AgentConfig{
		Name:         "Agent1",
		Capabilities: []identity.CapabilityType{identity.CapCodeWrite},
	})

	err := c.Join(a)
	if err != nil {
		t.Fatalf("Failed to join: %v", err)
	}

	if c.Size() != 1 {
		t.Errorf("Expected size 1 after join, got %d", c.Size())
	}

	// Test GetAgent
	retrieved, ok := c.GetAgent(a.Identity.SID)
	if !ok {
		t.Error("GetAgent should find the agent")
	}
	if retrieved.Identity.Name != "Agent1" {
		t.Errorf("Expected agent name 'Agent1', got '%s'", retrieved.Identity.Name)
	}

	// Test Leave
	err = c.Leave(a.Identity.SID)
	if err != nil {
		t.Fatalf("Failed to leave: %v", err)
	}

	if c.Size() != 0 {
		t.Errorf("Expected size 0 after leave, got %d", c.Size())
	}
}

func TestCollective_JoinFull(t *testing.T) {
	cfg := CollectiveConfig{
		MinAgents:          1,
		MaxAgents:          2,
		ConsensusThreshold: 0.67,
	}
	c := NewCollective("SmallCollective", cfg)

	a1, _ := agent.NewAgent(agent.AgentConfig{Name: "Agent1", Capabilities: []identity.CapabilityType{identity.CapCodeWrite}})
	a2, _ := agent.NewAgent(agent.AgentConfig{Name: "Agent2", Capabilities: []identity.CapabilityType{identity.CapCodeWrite}})
	a3, _ := agent.NewAgent(agent.AgentConfig{Name: "Agent3", Capabilities: []identity.CapabilityType{identity.CapCodeWrite}})

	c.Join(a1)
	c.Join(a2)

	err := c.Join(a3)
	if err != ErrCollectiveFull {
		t.Errorf("Expected ErrCollectiveFull, got %v", err)
	}
}

func TestCollective_LeaveNotFound(t *testing.T) {
	c := NewCollective("TestCollective", DefaultCollectiveConfig())

	err := c.Leave("non-existent-sid")
	if err != ErrAgentNotFound {
		t.Errorf("Expected ErrAgentNotFound, got %v", err)
	}
}

func TestCollective_GetAgents(t *testing.T) {
	c := NewCollective("TestCollective", DefaultCollectiveConfig())

	a1, _ := agent.NewAgent(agent.AgentConfig{Name: "Agent1", Capabilities: []identity.CapabilityType{identity.CapCodeWrite}})
	a2, _ := agent.NewAgent(agent.AgentConfig{Name: "Agent2", Capabilities: []identity.CapabilityType{identity.CapCodeReview}})

	c.Join(a1)
	c.Join(a2)

	agents := c.GetAgents()
	if len(agents) != 2 {
		t.Errorf("Expected 2 agents, got %d", len(agents))
	}
}

func TestCollective_Stats(t *testing.T) {
	c := NewCollective("TestCollective", DefaultCollectiveConfig())

	a1, _ := agent.NewAgent(agent.AgentConfig{Name: "Agent1", Capabilities: []identity.CapabilityType{identity.CapCodeWrite}})
	c.Join(a1)

	stats := c.Stats()

	if stats.Name != "TestCollective" {
		t.Errorf("Expected name 'TestCollective', got '%s'", stats.Name)
	}

	if stats.AgentCount != 1 {
		t.Errorf("Expected agent count 1, got %d", stats.AgentCount)
	}

	if stats.PendingTasks != 0 {
		t.Errorf("Expected 0 pending tasks, got %d", stats.PendingTasks)
	}
}

func TestCollective_StartStop(t *testing.T) {
	c := NewCollective("TestCollective", DefaultCollectiveConfig())

	a, _ := agent.NewAgent(agent.AgentConfig{
		Name:         "Agent1",
		Capabilities: []identity.CapabilityType{identity.CapCodeWrite},
	})
	c.Join(a)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start collective: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	c.Stop()

	// Verify agent is stopped
	if a.GetState() != agent.StateTerminated {
		t.Errorf("Expected agent state Terminated, got %s", a.GetState())
	}
}

func TestCollective_GetComponents(t *testing.T) {
	c := NewCollective("TestCollective", DefaultCollectiveConfig())

	if c.GetMemory() == nil {
		t.Error("GetMemory should not return nil")
	}

	if c.GetReputation() == nil {
		t.Error("GetReputation should not return nil")
	}

	if c.GetGossip() == nil {
		t.Error("GetGossip should not return nil")
	}

	if c.GetMarket() == nil {
		t.Error("GetMarket should not return nil")
	}

	if c.GetConsensus() == nil {
		t.Error("GetConsensus should not return nil")
	}
}
