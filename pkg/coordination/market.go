package coordination

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/squaremind/squaremind/pkg/agent"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrNoBids       = errors.New("no bids received")
	ErrMarketClosed = errors.New("market closed")
)

// Bid represents an agent's bid on a task
type Bid struct {
	AgentSID        string        `json:"agent_sid"`
	TaskID          string        `json:"task_id"`
	CapabilityScore float64       `json:"capability_score"`
	ReputationStake float64       `json:"reputation_stake"`
	EstimatedTime   time.Duration `json:"estimated_time"`
	Timestamp       time.Time     `json:"timestamp"`
}

// TaskAssignment represents the result of task matching
type TaskAssignment struct {
	TaskID   string `json:"task_id"`
	AgentSID string `json:"agent_sid"`
	Bid      *Bid   `json:"bid"`
}

// TaskMarket implements decentralized task allocation
type TaskMarket struct {
	mu sync.RWMutex

	listings map[string]*agent.Task // TaskID -> Task
	bids     map[string][]*Bid      // TaskID -> Bids

	bidTimeout time.Duration
	closed     bool
}

// NewTaskMarket creates a new task market
func NewTaskMarket() *TaskMarket {
	return &TaskMarket{
		listings:   make(map[string]*agent.Task),
		bids:       make(map[string][]*Bid),
		bidTimeout: 5 * time.Second,
	}
}

// Start begins market operation
func (m *TaskMarket) Start(ctx context.Context) {
	// Market runs passively, processing bids as they come
	go m.runMaintenanceLoop(ctx)
}

// runMaintenanceLoop cleans up expired listings
func (m *TaskMarket) runMaintenanceLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.cleanup()
		}
	}
}

// cleanup removes old listings
func (m *TaskMarket) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, task := range m.listings {
		// Remove listings older than deadline or 1 hour
		if !task.Deadline.IsZero() && now.After(task.Deadline) {
			delete(m.listings, id)
			delete(m.bids, id)
		} else if now.Sub(task.CreatedAt) > time.Hour {
			delete(m.listings, id)
			delete(m.bids, id)
		}
	}
}

// ListTask adds a task to the market
func (m *TaskMarket) ListTask(task *agent.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrMarketClosed
	}

	m.listings[task.ID] = task
	m.bids[task.ID] = make([]*Bid, 0)
	return nil
}

// UnlistTask removes a task from the market
func (m *TaskMarket) UnlistTask(taskID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.listings, taskID)
	delete(m.bids, taskID)
}

// GetListing returns a task listing
func (m *TaskMarket) GetListing(taskID string) (*agent.Task, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, ok := m.listings[taskID]
	return task, ok
}

// ListAllTasks returns all listed tasks
func (m *TaskMarket) ListAllTasks() []*agent.Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := make([]*agent.Task, 0, len(m.listings))
	for _, task := range m.listings {
		tasks = append(tasks, task)
	}
	return tasks
}

// SubmitBid submits a bid on a task
func (m *TaskMarket) SubmitBid(bid *Bid) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrMarketClosed
	}

	if _, exists := m.listings[bid.TaskID]; !exists {
		return ErrTaskNotFound
	}

	bid.Timestamp = time.Now()
	m.bids[bid.TaskID] = append(m.bids[bid.TaskID], bid)
	return nil
}

// GetBids returns all bids for a task
func (m *TaskMarket) GetBids(taskID string) []*Bid {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.bids[taskID]
}

// AssignTask matches a task to the best bidder
func (m *TaskMarket) AssignTask(
	task *agent.Task,
	agents map[string]*agent.Agent,
	reputation *ReputationRegistry,
) (*TaskAssignment, error) {
	// List the task
	if err := m.ListTask(task); err != nil {
		return nil, err
	}

	// Generate bids from capable agents
	for sid, a := range agents {
		if a.GetState() != agent.StateIdle {
			continue
		}

		score := a.Capabilities.MatchScore(task.Required)
		if score > 0.5 { // Minimum threshold
			bid := &Bid{
				AgentSID:        sid,
				TaskID:          task.ID,
				CapabilityScore: score,
				ReputationStake: a.Reputation.Overall * 0.1, // Stake 10% of reputation
				EstimatedTime:   estimateTime(task, score),
			}
			m.SubmitBid(bid)
		}
	}

	// Wait for bid collection period
	time.Sleep(m.bidTimeout)

	// Select best bid
	return m.selectBestBid(task.ID, reputation)
}

// selectBestBid chooses the winning bid
func (m *TaskMarket) selectBestBid(taskID string, reputation *ReputationRegistry) (*TaskAssignment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	bids := m.bids[taskID]
	if len(bids) == 0 {
		return nil, ErrNoBids
	}

	// Score each bid
	type scoredBid struct {
		bid   *Bid
		score float64
	}
	scored := make([]scoredBid, len(bids))

	for i, bid := range bids {
		rep := reputation.Get(bid.AgentSID)
		repScore := 50.0 // Default
		if rep != nil {
			repScore = rep.Overall
		}

		// Combined score: capability (40%), reputation (40%), stake (20%)
		score := bid.CapabilityScore*0.4 +
			(repScore/100)*0.4 +
			(bid.ReputationStake/100)*0.2

		scored[i] = scoredBid{bid, score}
	}

	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Winner is highest scored bid
	winner := scored[0].bid

	return &TaskAssignment{
		TaskID:   taskID,
		AgentSID: winner.AgentSID,
		Bid:      winner,
	}, nil
}

// estimateTime estimates task completion time based on complexity and capability
func estimateTime(task *agent.Task, capabilityScore float64) time.Duration {
	baseTime := time.Minute

	switch task.Complexity {
	case "low":
		baseTime = time.Minute
	case "medium":
		baseTime = 5 * time.Minute
	case "high":
		baseTime = 30 * time.Minute
	}

	// Higher capability = faster completion
	return time.Duration(float64(baseTime) / capabilityScore)
}

// Close closes the market
func (m *TaskMarket) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
}

// SetBidTimeout sets the bid collection timeout
func (m *TaskMarket) SetBidTimeout(timeout time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bidTimeout = timeout
}

// Stats returns market statistics
type MarketStats struct {
	ActiveListings int
	TotalBids      int
	AvgBidsPerTask float64
}

// Stats returns current market statistics
func (m *TaskMarket) Stats() MarketStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalBids := 0
	for _, bids := range m.bids {
		totalBids += len(bids)
	}

	avgBids := 0.0
	if len(m.listings) > 0 {
		avgBids = float64(totalBids) / float64(len(m.listings))
	}

	return MarketStats{
		ActiveListings: len(m.listings),
		TotalBids:      totalBids,
		AvgBidsPerTask: avgBids,
	}
}
