package agent

import (
	"time"

	"github.com/google/uuid"
	"github.com/squaremind/squaremind/pkg/identity"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskAssigned  TaskStatus = "assigned"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
)

// Task represents a unit of work
type Task struct {
	ID           string                    `json:"id"`
	Description  string                    `json:"description"`
	Requirements string                    `json:"requirements"`
	Complexity   string                    `json:"complexity"` // "low", "medium", "high"
	Required     []identity.CapabilityType `json:"required_capabilities"`
	Deadline     time.Time                 `json:"deadline"`
	Reward       float64                   `json:"reward"` // Reputation points
	Status       TaskStatus                `json:"status"`
	AssignedTo   string                    `json:"assigned_to,omitempty"` // Agent SID
	CreatedAt    time.Time                 `json:"created_at"`
}

// NewTask creates a new task
func NewTask(description string, required []identity.CapabilityType) *Task {
	return &Task{
		ID:          uuid.New().String(),
		Description: description,
		Required:    required,
		Status:      TaskPending,
		Complexity:  "medium",
		CreatedAt:   time.Now(),
	}
}

// WithComplexity sets the task complexity
func (t *Task) WithComplexity(complexity string) *Task {
	t.Complexity = complexity
	return t
}

// WithDeadline sets the task deadline
func (t *Task) WithDeadline(deadline time.Time) *Task {
	t.Deadline = deadline
	return t
}

// WithReward sets the task reward
func (t *Task) WithReward(reward float64) *Task {
	t.Reward = reward
	return t
}

// WithRequirements sets the task requirements
func (t *Task) WithRequirements(requirements string) *Task {
	t.Requirements = requirements
	return t
}

// TaskResult represents the result of a completed task
type TaskResult struct {
	TaskID    string        `json:"task_id"`
	AgentSID  string        `json:"agent_sid"`
	Status    TaskStatus    `json:"status"`
	Output    string        `json:"output"`
	Error     string        `json:"error,omitempty"`
	Quality   float64       `json:"quality"` // 0.0 - 1.0
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// Reputation tracks an agent's reputation
type Reputation struct {
	Overall     float64 `json:"overall"`     // 0-100
	Reliability float64 `json:"reliability"` // Completes tasks on time
	Quality     float64 `json:"quality"`     // Quality of outputs
	Cooperation float64 `json:"cooperation"` // Works well with others
	Honesty     float64 `json:"honesty"`     // Accurate self-assessment

	TasksCompleted int `json:"tasks_completed"`
	TasksFailed    int `json:"tasks_failed"`

	LastActive time.Time `json:"last_active"`
	DecayRate  float64   `json:"decay_rate"` // Daily decay percentage
}

// NewReputation creates a new reputation starting at baseline
func NewReputation() *Reputation {
	return &Reputation{
		Overall:     50.0,
		Reliability: 50.0,
		Quality:     50.0,
		Cooperation: 50.0,
		Honesty:     50.0,
		DecayRate:   0.01, // 1% daily decay
		LastActive:  time.Now(),
	}
}

// RecordSuccess updates reputation after successful task
func (r *Reputation) RecordSuccess(quality float64) {
	r.TasksCompleted++
	r.Quality = r.Quality*0.9 + quality*100*0.1 // Exponential moving average
	r.Reliability = r.Reliability*0.95 + 100*0.05
	r.recalculateOverall()
	r.LastActive = time.Now()
}

// RecordFailure updates reputation after failed task
func (r *Reputation) RecordFailure() {
	r.TasksFailed++
	r.Reliability = r.Reliability * 0.9 // 10% penalty
	r.recalculateOverall()
	r.LastActive = time.Now()
}

// RecordCooperation updates cooperation score
func (r *Reputation) RecordCooperation(score float64) {
	r.Cooperation = r.Cooperation*0.9 + score*100*0.1
	r.recalculateOverall()
}

// recalculateOverall updates the overall score
func (r *Reputation) recalculateOverall() {
	r.Overall = (r.Reliability + r.Quality + r.Cooperation + r.Honesty) / 4
}

// ApplyDecay applies time-based reputation decay
func (r *Reputation) ApplyDecay() {
	daysSinceActive := time.Since(r.LastActive).Hours() / 24
	if daysSinceActive > 1 {
		decayFactor := 1.0 - (r.DecayRate * daysSinceActive)
		if decayFactor < 0.5 {
			decayFactor = 0.5 // Floor at 50% of original
		}
		r.Overall *= decayFactor
	}
}

// AgentMemory represents an agent's memory store
type AgentMemory struct {
	ShortTerm map[string]interface{} `json:"short_term"`
	LongTerm  map[string]interface{} `json:"long_term"`
	Episodic  []Episode              `json:"episodic"`
}

// Episode represents a memorable event
type Episode struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Content   string                 `json:"content"`
	Context   map[string]interface{} `json:"context"`
	Timestamp time.Time              `json:"timestamp"`
	Salience  float64                `json:"salience"` // 0.0-1.0 importance
}

// NewAgentMemory creates a new agent memory
func NewAgentMemory() *AgentMemory {
	return &AgentMemory{
		ShortTerm: make(map[string]interface{}),
		LongTerm:  make(map[string]interface{}),
		Episodic:  make([]Episode, 0),
	}
}

// Store stores a value in short-term memory
func (m *AgentMemory) Store(key string, value interface{}) {
	m.ShortTerm[key] = value
}

// Recall retrieves a value from memory (checks short-term first, then long-term)
func (m *AgentMemory) Recall(key string) (interface{}, bool) {
	if v, ok := m.ShortTerm[key]; ok {
		return v, true
	}
	if v, ok := m.LongTerm[key]; ok {
		return v, true
	}
	return nil, false
}

// Consolidate moves important short-term memories to long-term
func (m *AgentMemory) Consolidate() {
	// Simple implementation: move everything
	for k, v := range m.ShortTerm {
		m.LongTerm[k] = v
	}
	m.ShortTerm = make(map[string]interface{})
}

// AddEpisode adds an episodic memory
func (m *AgentMemory) AddEpisode(ep Episode) {
	ep.ID = uuid.New().String()
	ep.Timestamp = time.Now()
	m.Episodic = append(m.Episodic, ep)

	// Keep only last 100 episodes
	if len(m.Episodic) > 100 {
		m.Episodic = m.Episodic[len(m.Episodic)-100:]
	}
}
