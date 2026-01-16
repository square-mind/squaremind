package collective

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/square-mind/squaremind/pkg/agent"
	"github.com/square-mind/squaremind/pkg/coordination"
)

var (
	ErrCollectiveFull = errors.New("collective at maximum capacity")
	ErrAgentNotFound  = errors.New("agent not found in collective")
)

// Collective represents a group of squaremind agents
type Collective struct {
	mu sync.RWMutex

	// Identity
	Name string
	ID   string

	// Agents
	agents map[string]*agent.Agent // SID -> Agent

	// Coordination
	gossip     *coordination.GossipProtocol
	market     *coordination.TaskMarket
	consensus  *coordination.ConsensusEngine
	reputation *coordination.ReputationRegistry

	// Shared Memory
	memory *CollectiveMemory

	// Configuration
	config CollectiveConfig

	// Task tracking
	pendingTasks   []*agent.Task
	activeTasks    map[string]*agent.Task
	completedTasks []*agent.TaskResult
}

// CollectiveConfig holds collective configuration
type CollectiveConfig struct {
	MinAgents          int     `json:"min_agents"`
	MaxAgents          int     `json:"max_agents"`
	ConsensusThreshold float64 `json:"consensus_threshold"` // e.g., 0.67 for 2/3
	ReputationDecay    float64 `json:"reputation_decay"`    // Daily decay rate
}

// DefaultCollectiveConfig returns sensible defaults
func DefaultCollectiveConfig() CollectiveConfig {
	return CollectiveConfig{
		MinAgents:          2,
		MaxAgents:          100,
		ConsensusThreshold: 0.67,
		ReputationDecay:    0.01,
	}
}

// NewCollective creates a new collective
func NewCollective(name string, cfg CollectiveConfig) *Collective {
	return &Collective{
		Name:           name,
		ID:             uuid.New().String(),
		agents:         make(map[string]*agent.Agent),
		gossip:         coordination.NewGossipProtocol(),
		market:         coordination.NewTaskMarket(),
		consensus:      coordination.NewConsensusEngine(cfg.ConsensusThreshold),
		reputation:     coordination.NewReputationRegistry(),
		memory:         NewCollectiveMemory(),
		config:         cfg,
		activeTasks:    make(map[string]*agent.Task),
		pendingTasks:   make([]*agent.Task, 0),
		completedTasks: make([]*agent.TaskResult, 0),
	}
}

// Join adds an agent to the collective
func (c *Collective) Join(a *agent.Agent) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.agents) >= c.config.MaxAgents {
		return ErrCollectiveFull
	}

	c.agents[a.Identity.SID] = a
	c.gossip.AddPeer(a.Identity.SID)
	c.reputation.Register(a.Identity.SID, a.Reputation)

	// Broadcast join to other agents
	c.gossip.Broadcast(coordination.Message{
		Type:    coordination.MsgAgentJoined,
		From:    a.Identity.SID,
		Payload: a.Identity,
	})

	return nil
}

// Leave removes an agent from the collective
func (c *Collective) Leave(sid string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.agents[sid]; !exists {
		return ErrAgentNotFound
	}

	delete(c.agents, sid)
	c.gossip.RemovePeer(sid)
	c.reputation.Unregister(sid)

	// Broadcast leave
	c.gossip.Broadcast(coordination.Message{
		Type: coordination.MsgAgentLeft,
		From: sid,
	})

	return nil
}

// GetAgent returns an agent by SID
func (c *Collective) GetAgent(sid string) (*agent.Agent, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	a, ok := c.agents[sid]
	return a, ok
}

// Submit submits a task to the collective
func (c *Collective) Submit(task *agent.Task) (*agent.TaskResult, error) {
	c.mu.Lock()
	c.pendingTasks = append(c.pendingTasks, task)
	c.mu.Unlock()

	// Broadcast task to market
	c.gossip.Broadcast(coordination.Message{
		Type:    coordination.MsgTaskAvailable,
		Payload: task,
	})

	// Let market handle bidding and assignment
	assignment, err := c.market.AssignTask(task, c.agents, c.reputation)
	if err != nil {
		return nil, err
	}

	// Move to active
	c.mu.Lock()
	c.activeTasks[task.ID] = task
	task.Status = agent.TaskAssigned
	task.AssignedTo = assignment.AgentSID
	c.mu.Unlock()

	// Submit to assigned agent
	assignedAgent := c.agents[assignment.AgentSID]
	assignedAgent.SubmitTask(task)

	// Wait for result
	result := <-assignedAgent.GetResults()

	// Update reputation
	if result.Status == agent.TaskCompleted {
		c.reputation.RecordTaskSuccess(assignment.AgentSID, result.Quality)
	} else {
		c.reputation.RecordTaskFailure(assignment.AgentSID)
	}

	// Record completion
	c.mu.Lock()
	delete(c.activeTasks, task.ID)
	c.completedTasks = append(c.completedTasks, result)
	c.mu.Unlock()

	return result, nil
}

// SubmitAsync submits a task without waiting for result
func (c *Collective) SubmitAsync(task *agent.Task) (string, error) {
	go func() {
		_, _ = c.Submit(task)
	}()
	return task.ID, nil
}

// GetAgents returns all agents in the collective
func (c *Collective) GetAgents() []*agent.Agent {
	c.mu.RLock()
	defer c.mu.RUnlock()

	agents := make([]*agent.Agent, 0, len(c.agents))
	for _, a := range c.agents {
		agents = append(agents, a)
	}
	return agents
}

// Size returns the number of agents
func (c *Collective) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.agents)
}

// Start begins collective operation
func (c *Collective) Start(ctx context.Context) error {
	// Start all coordination systems
	go c.gossip.Start(ctx)
	go c.market.Start(ctx)
	go c.runMaintenanceLoop(ctx)

	// Start all agents
	for _, a := range c.agents {
		if err := a.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Stop stops the collective
func (c *Collective) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, a := range c.agents {
		a.Stop()
	}

	c.market.Close()
}

// runMaintenanceLoop handles periodic collective maintenance
func (c *Collective) runMaintenanceLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.maintenance()
		}
	}
}

// maintenance performs periodic collective maintenance
func (c *Collective) maintenance() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Apply reputation decay
	c.reputation.ApplyDecayAll()

	// Reassign stalled tasks
	for id, task := range c.activeTasks {
		if !task.Deadline.IsZero() && time.Since(task.CreatedAt) > task.Deadline.Sub(task.CreatedAt)*2 {
			// Task is taking too long, consider reassignment
			delete(c.activeTasks, id)
			task.Status = agent.TaskPending
			c.pendingTasks = append(c.pendingTasks, task)
		}
	}
}

// GetMemory returns the collective memory
func (c *Collective) GetMemory() *CollectiveMemory {
	return c.memory
}

// GetReputation returns the reputation registry
func (c *Collective) GetReputation() *coordination.ReputationRegistry {
	return c.reputation
}

// GetGossip returns the gossip protocol
func (c *Collective) GetGossip() *coordination.GossipProtocol {
	return c.gossip
}

// GetMarket returns the task market
func (c *Collective) GetMarket() *coordination.TaskMarket {
	return c.market
}

// GetConsensus returns the consensus engine
func (c *Collective) GetConsensus() *coordination.ConsensusEngine {
	return c.consensus
}

// Stats returns collective statistics
type CollectiveStats struct {
	Name            string
	AgentCount      int
	ActiveTasks     int
	CompletedTasks  int
	PendingTasks    int
	AvgReputation   float64
}

// Stats returns current collective statistics
func (c *Collective) Stats() CollectiveStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return CollectiveStats{
		Name:           c.Name,
		AgentCount:     len(c.agents),
		ActiveTasks:    len(c.activeTasks),
		CompletedTasks: len(c.completedTasks),
		PendingTasks:   len(c.pendingTasks),
		AvgReputation:  c.reputation.AverageReputation(),
	}
}
