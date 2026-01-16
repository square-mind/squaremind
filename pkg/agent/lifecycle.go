package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/square-mind/squaremind/pkg/identity"
	"github.com/square-mind/squaremind/pkg/llm"
)

// LifecycleManager manages agent lifecycle events
type LifecycleManager struct {
	mu sync.RWMutex

	runtime  *Runtime
	provider llm.Provider
	model    string

	// Event handlers
	onSpawn     []func(*Agent)
	onTerminate []func(*Agent)
}

// NewLifecycleManager creates a new lifecycle manager
func NewLifecycleManager(runtime *Runtime, provider llm.Provider, model string) *LifecycleManager {
	return &LifecycleManager{
		runtime:     runtime,
		provider:    provider,
		model:       model,
		onSpawn:     make([]func(*Agent), 0),
		onTerminate: make([]func(*Agent), 0),
	}
}

// Spawn creates and starts a new agent
func (lm *LifecycleManager) Spawn(ctx context.Context, name string, capabilities []identity.CapabilityType) (*Agent, error) {
	agent, err := NewAgent(AgentConfig{
		Name:         name,
		Capabilities: capabilities,
		Provider:     lm.provider,
		Model:        lm.model,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	// Register with runtime
	if err := lm.runtime.Register(agent); err != nil {
		return nil, fmt.Errorf("failed to register agent: %w", err)
	}

	// Start agent
	if err := agent.Start(ctx); err != nil {
		_ = lm.runtime.Unregister(agent.Identity.SID)
		return nil, fmt.Errorf("failed to start agent: %w", err)
	}

	// Notify handlers
	lm.mu.RLock()
	for _, handler := range lm.onSpawn {
		handler(agent)
	}
	lm.mu.RUnlock()

	return agent, nil
}

// SpawnChild creates a child agent from a parent
func (lm *LifecycleManager) SpawnChild(ctx context.Context, parent *Agent, name string, capabilities []identity.CapabilityType) (*Agent, error) {
	agent, err := NewAgent(AgentConfig{
		Name:         name,
		Capabilities: capabilities,
		Provider:     lm.provider,
		Model:        lm.model,
		ParentSID:    parent.Identity.SID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create child agent: %w", err)
	}

	// Inherit some reputation from parent
	agent.Reputation.Overall = parent.Reputation.Overall * 0.5
	agent.Reputation.Honesty = parent.Reputation.Honesty * 0.75

	// Register with runtime
	if err := lm.runtime.Register(agent); err != nil {
		return nil, fmt.Errorf("failed to register agent: %w", err)
	}

	// Start agent
	if err := agent.Start(ctx); err != nil {
		_ = lm.runtime.Unregister(agent.Identity.SID)
		return nil, fmt.Errorf("failed to start agent: %w", err)
	}

	// Notify handlers
	lm.mu.RLock()
	for _, handler := range lm.onSpawn {
		handler(agent)
	}
	lm.mu.RUnlock()

	return agent, nil
}

// Terminate stops and removes an agent
func (lm *LifecycleManager) Terminate(sid string) error {
	agent, ok := lm.runtime.GetAgent(sid)
	if !ok {
		return fmt.Errorf("agent %s not found", sid)
	}

	// Notify handlers before termination
	lm.mu.RLock()
	for _, handler := range lm.onTerminate {
		handler(agent)
	}
	lm.mu.RUnlock()

	// Stop agent
	agent.Stop()

	// Unregister from runtime
	return lm.runtime.Unregister(sid)
}

// TerminateAll terminates all agents
func (lm *LifecycleManager) TerminateAll() {
	for _, agent := range lm.runtime.ListAgents() {
		_ = lm.Terminate(agent.Identity.SID)
	}
}

// OnSpawn registers a callback for agent spawn events
func (lm *LifecycleManager) OnSpawn(handler func(*Agent)) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.onSpawn = append(lm.onSpawn, handler)
}

// OnTerminate registers a callback for agent termination events
func (lm *LifecycleManager) OnTerminate(handler func(*Agent)) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.onTerminate = append(lm.onTerminate, handler)
}

// HealthCheck checks agent health
type HealthStatus struct {
	SID        string
	Name       string
	State      AgentState
	Healthy    bool
	LastActive time.Time
	Uptime     time.Duration
	Message    string
}

// HealthCheck returns health status for an agent
func (lm *LifecycleManager) HealthCheck(sid string) (*HealthStatus, error) {
	agent, ok := lm.runtime.GetAgent(sid)
	if !ok {
		return nil, fmt.Errorf("agent %s not found", sid)
	}

	status := &HealthStatus{
		SID:        agent.Identity.SID,
		Name:       agent.Identity.Name,
		State:      agent.GetState(),
		LastActive: agent.LastActive,
		Uptime:     time.Since(agent.StartedAt),
	}

	// Determine health
	switch agent.GetState() {
	case StateTerminated:
		status.Healthy = false
		status.Message = "Agent terminated"
	case StatePaused:
		status.Healthy = true
		status.Message = "Agent paused"
	case StateWorking:
		status.Healthy = true
		status.Message = "Agent working"
	case StateIdle:
		if time.Since(agent.LastActive) > 5*time.Minute {
			status.Healthy = true
			status.Message = "Agent idle (inactive)"
		} else {
			status.Healthy = true
			status.Message = "Agent idle"
		}
	default:
		status.Healthy = false
		status.Message = "Unknown state"
	}

	return status, nil
}

// HealthCheckAll returns health status for all agents
func (lm *LifecycleManager) HealthCheckAll() []*HealthStatus {
	agents := lm.runtime.ListAgents()
	statuses := make([]*HealthStatus, 0, len(agents))

	for _, agent := range agents {
		status, err := lm.HealthCheck(agent.Identity.SID)
		if err == nil {
			statuses = append(statuses, status)
		}
	}

	return statuses
}
