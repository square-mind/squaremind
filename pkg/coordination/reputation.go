package coordination

import (
	"sync"
	"time"

	"github.com/square-mind/squaremind/pkg/agent"
)

// ReputationRegistry manages reputation scores for all agents
type ReputationRegistry struct {
	mu sync.RWMutex

	scores  map[string]*agent.Reputation // SID -> Reputation
	history map[string][]ReputationEvent // SID -> Events
}

// ReputationEvent represents a reputation change event
type ReputationEvent struct {
	AgentSID   string    `json:"agent_sid"`
	Type       string    `json:"type"` // "task_success", "task_failure", "peer_rating", "decay"
	Delta      float64   `json:"delta"`
	Reason     string    `json:"reason"`
	Timestamp  time.Time `json:"timestamp"`
}

// NewReputationRegistry creates a new reputation registry
func NewReputationRegistry() *ReputationRegistry {
	return &ReputationRegistry{
		scores:  make(map[string]*agent.Reputation),
		history: make(map[string][]ReputationEvent),
	}
}

// Register registers an agent with initial reputation
func (r *ReputationRegistry) Register(sid string, rep *agent.Reputation) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.scores[sid] = rep
	r.history[sid] = make([]ReputationEvent, 0)
}

// Unregister removes an agent from the registry
func (r *ReputationRegistry) Unregister(sid string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.scores, sid)
	delete(r.history, sid)
}

// Get returns an agent's reputation
func (r *ReputationRegistry) Get(sid string) *agent.Reputation {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.scores[sid]
}

// GetAll returns all reputations
func (r *ReputationRegistry) GetAll() map[string]*agent.Reputation {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]*agent.Reputation)
	for sid, rep := range r.scores {
		result[sid] = rep
	}
	return result
}

// Update updates an agent's reputation
func (r *ReputationRegistry) Update(sid string, rep *agent.Reputation) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.scores[sid] = rep
}

// RecordTaskSuccess records a successful task completion
func (r *ReputationRegistry) RecordTaskSuccess(sid string, quality float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rep, ok := r.scores[sid]
	if !ok {
		return
	}

	oldOverall := rep.Overall
	rep.RecordSuccess(quality)

	// Record event
	event := ReputationEvent{
		AgentSID:  sid,
		Type:      "task_success",
		Delta:     rep.Overall - oldOverall,
		Reason:    "Task completed successfully",
		Timestamp: time.Now(),
	}
	r.history[sid] = append(r.history[sid], event)

	// Keep history bounded
	if len(r.history[sid]) > 100 {
		r.history[sid] = r.history[sid][len(r.history[sid])-100:]
	}
}

// RecordTaskFailure records a failed task
func (r *ReputationRegistry) RecordTaskFailure(sid string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rep, ok := r.scores[sid]
	if !ok {
		return
	}

	oldOverall := rep.Overall
	rep.RecordFailure()

	// Record event
	event := ReputationEvent{
		AgentSID:  sid,
		Type:      "task_failure",
		Delta:     rep.Overall - oldOverall,
		Reason:    "Task failed",
		Timestamp: time.Now(),
	}
	r.history[sid] = append(r.history[sid], event)
}

// RecordPeerRating records a peer rating
func (r *ReputationRegistry) RecordPeerRating(sid string, raterSID string, rating float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rep, ok := r.scores[sid]
	if !ok {
		return
	}

	// Verify rater exists and has sufficient reputation to rate
	raterRep, ok := r.scores[raterSID]
	if !ok || raterRep.Overall < 30 {
		return // Rater needs minimum reputation
	}

	oldOverall := rep.Overall

	// Weight rating by rater's reputation
	weight := raterRep.Overall / 100
	rep.Cooperation = rep.Cooperation*0.9 + rating*100*0.1*weight
	// Recalculate overall score
	rep.Overall = (rep.Reliability + rep.Quality + rep.Cooperation + rep.Honesty) / 4

	// Record event
	event := ReputationEvent{
		AgentSID:  sid,
		Type:      "peer_rating",
		Delta:     rep.Overall - oldOverall,
		Reason:    "Rated by peer " + raterSID,
		Timestamp: time.Now(),
	}
	r.history[sid] = append(r.history[sid], event)
}

// ApplyDecayAll applies decay to all agents
func (r *ReputationRegistry) ApplyDecayAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for sid, rep := range r.scores {
		oldOverall := rep.Overall
		rep.ApplyDecay()

		if rep.Overall != oldOverall {
			event := ReputationEvent{
				AgentSID:  sid,
				Type:      "decay",
				Delta:     rep.Overall - oldOverall,
				Reason:    "Time-based decay",
				Timestamp: time.Now(),
			}
			r.history[sid] = append(r.history[sid], event)
		}
	}
}

// GetHistory returns reputation history for an agent
func (r *ReputationRegistry) GetHistory(sid string) []ReputationEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.history[sid]
}

// GetTopAgents returns the top N agents by reputation
func (r *ReputationRegistry) GetTopAgents(n int) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type agentRep struct {
		sid   string
		score float64
	}

	agents := make([]agentRep, 0, len(r.scores))
	for sid, rep := range r.scores {
		agents = append(agents, agentRep{sid, rep.Overall})
	}

	// Sort by score descending
	for i := 0; i < len(agents)-1; i++ {
		for j := i + 1; j < len(agents); j++ {
			if agents[j].score > agents[i].score {
				agents[i], agents[j] = agents[j], agents[i]
			}
		}
	}

	// Take top N
	result := make([]string, 0, n)
	for i := 0; i < n && i < len(agents); i++ {
		result = append(result, agents[i].sid)
	}

	return result
}

// AverageReputation returns the average reputation across all agents
func (r *ReputationRegistry) AverageReputation() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.scores) == 0 {
		return 0
	}

	total := 0.0
	for _, rep := range r.scores {
		total += rep.Overall
	}

	return total / float64(len(r.scores))
}

