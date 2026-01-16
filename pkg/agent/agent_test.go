package agent

import (
	"context"
	"testing"
	"time"

	"github.com/squaremind/squaremind/pkg/identity"
)

func TestNewAgent(t *testing.T) {
	cfg := AgentConfig{
		Name:         "TestAgent",
		Capabilities: []identity.CapabilityType{identity.CapCodeWrite, identity.CapCodeReview},
		Model:        "test-model",
	}

	agent, err := NewAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	if agent.Identity.Name != "TestAgent" {
		t.Errorf("Expected name 'TestAgent', got '%s'", agent.Identity.Name)
	}

	if !agent.Capabilities.Has(identity.CapCodeWrite) {
		t.Error("Agent should have CapCodeWrite capability")
	}

	if !agent.Capabilities.Has(identity.CapCodeReview) {
		t.Error("Agent should have CapCodeReview capability")
	}

	if agent.State != StateInitializing {
		t.Errorf("Expected state Initializing, got %s", agent.State)
	}

	if agent.Reputation == nil {
		t.Error("Reputation should be initialized")
	}

	if agent.Memory == nil {
		t.Error("Memory should be initialized")
	}
}

func TestAgent_StartStop(t *testing.T) {
	agent, _ := NewAgent(AgentConfig{
		Name:         "TestAgent",
		Capabilities: []identity.CapabilityType{identity.CapCodeWrite},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := agent.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start agent: %v", err)
	}

	// Give the goroutine time to start
	time.Sleep(10 * time.Millisecond)

	if agent.GetState() != StateIdle {
		t.Errorf("Expected state Idle after start, got %s", agent.GetState())
	}

	agent.Stop()

	// Give the goroutine time to stop
	time.Sleep(10 * time.Millisecond)

	if agent.GetState() != StateTerminated {
		t.Errorf("Expected state Terminated after stop, got %s", agent.GetState())
	}
}

func TestAgent_PauseResume(t *testing.T) {
	agent, _ := NewAgent(AgentConfig{
		Name:         "TestAgent",
		Capabilities: []identity.CapabilityType{identity.CapCodeWrite},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agent.Start(ctx)
	time.Sleep(10 * time.Millisecond)

	agent.Pause()
	if agent.GetState() != StatePaused {
		t.Errorf("Expected state Paused, got %s", agent.GetState())
	}

	agent.Resume()
	if agent.GetState() != StateIdle {
		t.Errorf("Expected state Idle after resume, got %s", agent.GetState())
	}

	agent.Stop()
}

func TestNewTask(t *testing.T) {
	task := NewTask("Test task", []identity.CapabilityType{identity.CapCodeWrite})

	if task.Description != "Test task" {
		t.Errorf("Expected description 'Test task', got '%s'", task.Description)
	}

	if task.ID == "" {
		t.Error("Task ID should not be empty")
	}

	if task.Status != TaskPending {
		t.Errorf("Expected status Pending, got %s", task.Status)
	}

	if len(task.Required) != 1 {
		t.Errorf("Expected 1 required capability, got %d", len(task.Required))
	}
}

func TestTask_Builders(t *testing.T) {
	deadline := time.Now().Add(time.Hour)

	task := NewTask("Test", nil).
		WithComplexity("high").
		WithDeadline(deadline).
		WithReward(50.0).
		WithRequirements("Must be fast")

	if task.Complexity != "high" {
		t.Errorf("Expected complexity 'high', got '%s'", task.Complexity)
	}

	if task.Deadline != deadline {
		t.Error("Deadline not set correctly")
	}

	if task.Reward != 50.0 {
		t.Errorf("Expected reward 50.0, got %f", task.Reward)
	}

	if task.Requirements != "Must be fast" {
		t.Errorf("Expected requirements 'Must be fast', got '%s'", task.Requirements)
	}
}

func TestReputation(t *testing.T) {
	rep := NewReputation()

	if rep.Overall != 50.0 {
		t.Errorf("Expected initial overall 50.0, got %f", rep.Overall)
	}

	// Test success
	rep.RecordSuccess(0.9)
	if rep.TasksCompleted != 1 {
		t.Errorf("Expected 1 task completed, got %d", rep.TasksCompleted)
	}

	// Test failure
	rep.RecordFailure()
	if rep.TasksFailed != 1 {
		t.Errorf("Expected 1 task failed, got %d", rep.TasksFailed)
	}
}

func TestAgentMemory(t *testing.T) {
	mem := NewAgentMemory()

	// Test Store and Recall
	mem.Store("key1", "value1")

	val, ok := mem.Recall("key1")
	if !ok {
		t.Error("Should find stored value")
	}
	if val != "value1" {
		t.Errorf("Expected 'value1', got '%v'", val)
	}

	// Test non-existent key
	_, ok = mem.Recall("nonexistent")
	if ok {
		t.Error("Should not find non-existent key")
	}

	// Test Consolidate
	mem.Consolidate()
	val, ok = mem.Recall("key1")
	if !ok {
		t.Error("Should find value in long-term memory after consolidation")
	}

	// Test AddEpisode
	mem.AddEpisode(Episode{
		Type:     "test",
		Content:  "test content",
		Salience: 0.8,
	})

	if len(mem.Episodic) != 1 {
		t.Errorf("Expected 1 episode, got %d", len(mem.Episodic))
	}
}
