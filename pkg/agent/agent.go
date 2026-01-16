package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/square-mind/squaremind/pkg/identity"
	"github.com/square-mind/squaremind/pkg/llm"
)

// AgentState represents the current state of an agent
type AgentState string

const (
	StateInitializing AgentState = "initializing"
	StateIdle         AgentState = "idle"
	StateWorking      AgentState = "working"
	StatePaused       AgentState = "paused"
	StateTerminated   AgentState = "terminated"
)

// Agent represents a squaremind AI agent
type Agent struct {
	mu sync.RWMutex

	// Identity
	Identity *identity.SquaremindIdentity

	// Capabilities
	Capabilities *identity.CapabilitySet

	// LLM Backend
	Provider llm.Provider
	Model    string

	// State
	State       AgentState
	CurrentTask *Task

	// Reputation
	Reputation *Reputation

	// Memory
	Memory *AgentMemory

	// Channels for coordination
	taskChan   chan *Task
	resultChan chan *TaskResult
	stopChan   chan struct{}

	// Lifecycle
	StartedAt  time.Time
	LastActive time.Time
}

// AgentConfig holds configuration for creating a new agent
type AgentConfig struct {
	Name         string
	Capabilities []identity.CapabilityType
	Provider     llm.Provider
	Model        string
	ParentSID    string
}

// NewAgent creates a new squaremind agent
func NewAgent(cfg AgentConfig) (*Agent, error) {
	// Create identity
	id, err := identity.NewSquaremindIdentity(cfg.Name, cfg.ParentSID)
	if err != nil {
		return nil, err
	}

	// Initialize capabilities
	capSet := identity.NewCapabilitySet()
	for _, capType := range cfg.Capabilities {
		capSet.Add(&identity.Capability{
			Type:        capType,
			Proficiency: 0.5, // Start at 50%, improve through tasks
		})
	}

	return &Agent{
		Identity:     id,
		Capabilities: capSet,
		Provider:     cfg.Provider,
		Model:        cfg.Model,
		State:        StateInitializing,
		Reputation:   NewReputation(),
		Memory:       NewAgentMemory(),
		taskChan:     make(chan *Task, 10),
		resultChan:   make(chan *TaskResult, 10),
		stopChan:     make(chan struct{}),
		StartedAt:    time.Now(),
		LastActive:   time.Now(),
	}, nil
}

// Start begins the agent's autonomous operation
func (a *Agent) Start(ctx context.Context) error {
	a.mu.Lock()
	a.State = StateIdle
	a.mu.Unlock()

	go a.runLoop(ctx)
	return nil
}

// runLoop is the main agent operation loop
func (a *Agent) runLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			a.terminate()
			return
		case <-a.stopChan:
			a.terminate()
			return
		case task := <-a.taskChan:
			a.executeTask(ctx, task)
		}
	}
}

// executeTask handles task execution
func (a *Agent) executeTask(ctx context.Context, task *Task) {
	a.mu.Lock()
	a.State = StateWorking
	a.CurrentTask = task
	a.LastActive = time.Now()
	a.mu.Unlock()

	startTime := time.Now()

	// Execute with LLM
	result, err := a.performTask(ctx, task)
	result.Duration = time.Since(startTime)
	result.Timestamp = time.Now()
	result.AgentSID = a.Identity.SID

	a.mu.Lock()
	a.State = StateIdle
	a.CurrentTask = nil
	a.mu.Unlock()

	// Update reputation based on result
	if err != nil {
		a.Reputation.RecordFailure()
	} else {
		a.Reputation.RecordSuccess(result.Quality)
	}

	// Add to episodic memory
	a.Memory.AddEpisode(Episode{
		Type:    "task_completion",
		Content: fmt.Sprintf("Completed task: %s", task.Description),
		Context: map[string]interface{}{
			"task_id": task.ID,
			"quality": result.Quality,
			"status":  result.Status,
		},
		Salience: result.Quality,
	})

	// Send result
	select {
	case a.resultChan <- result:
	default:
		// Channel full, drop result
	}
}

// performTask uses the LLM to perform the actual task
func (a *Agent) performTask(ctx context.Context, task *Task) (*TaskResult, error) {
	// If no provider, return simulated result
	if a.Provider == nil {
		return &TaskResult{
			TaskID:  task.ID,
			Status:  TaskCompleted,
			Output:  fmt.Sprintf("[Simulated] Completed: %s", task.Description),
			Quality: 0.75,
		}, nil
	}

	prompt := a.buildPrompt(task)

	response, err := a.Provider.Complete(ctx, llm.CompletionRequest{
		Model:  a.Model,
		Prompt: prompt,
	})
	if err != nil {
		return &TaskResult{
			TaskID: task.ID,
			Status: TaskFailed,
			Error:  err.Error(),
		}, err
	}

	return &TaskResult{
		TaskID:  task.ID,
		Status:  TaskCompleted,
		Output:  response.Content,
		Quality: 0.8, // Would be evaluated by quality assessment
	}, nil
}

// buildPrompt constructs the prompt for the LLM
func (a *Agent) buildPrompt(task *Task) string {
	capsJSON := a.Capabilities.ToJSON()

	return fmt.Sprintf(`You are a squaremind AI agent with the following identity:
Name: %s
SID: %s
Capabilities: %s

Your task:
%s

Requirements:
%s

Perform this task to the best of your ability. Be thorough and precise.`,
		a.Identity.Name,
		a.Identity.SID,
		capsJSON,
		task.Description,
		task.Requirements,
	)
}

// Stop signals the agent to stop
func (a *Agent) Stop() {
	close(a.stopChan)
}

// terminate cleans up agent resources
func (a *Agent) terminate() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.State = StateTerminated
}

// SubmitTask submits a task to the agent
func (a *Agent) SubmitTask(task *Task) {
	select {
	case a.taskChan <- task:
	default:
		// Channel full
	}
}

// GetResults returns the results channel
func (a *Agent) GetResults() <-chan *TaskResult {
	return a.resultChan
}

// GetState returns the current agent state
func (a *Agent) GetState() AgentState {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.State
}

// GetCurrentTask returns the current task being worked on
func (a *Agent) GetCurrentTask() *Task {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.CurrentTask
}

// Pause pauses the agent
func (a *Agent) Pause() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.State == StateIdle {
		a.State = StatePaused
	}
}

// Resume resumes a paused agent
func (a *Agent) Resume() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.State == StatePaused {
		a.State = StateIdle
	}
}
