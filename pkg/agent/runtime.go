package agent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Runtime manages agent execution
type Runtime struct {
	mu sync.RWMutex

	agents     map[string]*Agent // SID -> Agent
	maxAgents  int
	taskQueue  chan *Task
	results    chan *TaskResult
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

// RuntimeConfig configures the runtime
type RuntimeConfig struct {
	MaxAgents   int
	TaskBuffer  int
}

// DefaultRuntimeConfig returns default configuration
func DefaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		MaxAgents:  100,
		TaskBuffer: 1000,
	}
}

// NewRuntime creates a new agent runtime
func NewRuntime(cfg RuntimeConfig) *Runtime {
	return &Runtime{
		agents:    make(map[string]*Agent),
		maxAgents: cfg.MaxAgents,
		taskQueue: make(chan *Task, cfg.TaskBuffer),
		results:   make(chan *TaskResult, cfg.TaskBuffer),
		stopChan:  make(chan struct{}),
	}
}

// Register registers an agent with the runtime
func (r *Runtime) Register(a *Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.agents) >= r.maxAgents {
		return fmt.Errorf("runtime at capacity (%d agents)", r.maxAgents)
	}

	r.agents[a.Identity.SID] = a
	return nil
}

// Unregister removes an agent from the runtime
func (r *Runtime) Unregister(sid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[sid]; !exists {
		return fmt.Errorf("agent %s not found", sid)
	}

	delete(r.agents, sid)
	return nil
}

// GetAgent returns an agent by SID
func (r *Runtime) GetAgent(sid string) (*Agent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	a, ok := r.agents[sid]
	return a, ok
}

// ListAgents returns all registered agents
func (r *Runtime) ListAgents() []*Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]*Agent, 0, len(r.agents))
	for _, a := range r.agents {
		agents = append(agents, a)
	}
	return agents
}

// Start starts the runtime
func (r *Runtime) Start(ctx context.Context) error {
	// Start all agents
	r.mu.RLock()
	for _, a := range r.agents {
		if err := a.Start(ctx); err != nil {
			r.mu.RUnlock()
			return err
		}
	}
	r.mu.RUnlock()

	// Start task dispatcher
	r.wg.Add(1)
	go r.dispatchTasks(ctx)

	// Start result collector
	r.wg.Add(1)
	go r.collectResults(ctx)

	return nil
}

// Stop stops the runtime
func (r *Runtime) Stop() {
	close(r.stopChan)

	// Stop all agents
	r.mu.RLock()
	for _, a := range r.agents {
		a.Stop()
	}
	r.mu.RUnlock()

	r.wg.Wait()
}

// SubmitTask submits a task to the runtime for execution
func (r *Runtime) SubmitTask(task *Task) error {
	select {
	case r.taskQueue <- task:
		return nil
	default:
		return fmt.Errorf("task queue full")
	}
}

// GetResults returns the results channel
func (r *Runtime) GetResults() <-chan *TaskResult {
	return r.results
}

// dispatchTasks dispatches tasks to capable agents
func (r *Runtime) dispatchTasks(ctx context.Context) {
	defer r.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.stopChan:
			return
		case task := <-r.taskQueue:
			r.assignTask(task)
		}
	}
}

// assignTask assigns a task to the best available agent
func (r *Runtime) assignTask(task *Task) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var bestAgent *Agent
	var bestScore float64

	for _, a := range r.agents {
		if a.GetState() != StateIdle {
			continue
		}

		score := a.Capabilities.MatchScore(task.Required)
		if score > bestScore {
			bestScore = score
			bestAgent = a
		}
	}

	if bestAgent != nil && bestScore > 0.5 {
		task.Status = TaskAssigned
		task.AssignedTo = bestAgent.Identity.SID
		bestAgent.SubmitTask(task)
	} else {
		// No suitable agent found, re-queue
		go func() {
			time.Sleep(time.Second)
			r.taskQueue <- task
		}()
	}
}

// collectResults collects results from all agents
func (r *Runtime) collectResults(ctx context.Context) {
	defer r.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.stopChan:
			return
		case <-ticker.C:
			r.mu.RLock()
			for _, a := range r.agents {
				select {
				case result := <-a.GetResults():
					select {
					case r.results <- result:
					default:
						// Results channel full
					}
				default:
					// No results from this agent
				}
			}
			r.mu.RUnlock()
		}
	}
}

// Stats returns runtime statistics
type RuntimeStats struct {
	TotalAgents   int
	IdleAgents    int
	WorkingAgents int
	PendingTasks  int
}

// Stats returns current runtime statistics
func (r *Runtime) Stats() RuntimeStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := RuntimeStats{
		TotalAgents:  len(r.agents),
		PendingTasks: len(r.taskQueue),
	}

	for _, a := range r.agents {
		switch a.GetState() {
		case StateIdle:
			stats.IdleAgents++
		case StateWorking:
			stats.WorkingAgents++
		}
	}

	return stats
}
