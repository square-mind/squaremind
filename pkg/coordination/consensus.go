package coordination

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrConsensusTimeout = errors.New("consensus timeout")
	ErrInsufficientVotes = errors.New("insufficient votes for consensus")
)

// ConsensusType represents the type of consensus being reached
type ConsensusType string

const (
	ConsensusTypeTaskAssignment ConsensusType = "task_assignment"
	ConsensusTypeAgentSpawn     ConsensusType = "agent_spawn"
	ConsensusTypeAgentTerminate ConsensusType = "agent_terminate"
	ConsensusTypeParameterChange ConsensusType = "parameter_change"
)

// Proposal represents a proposal for consensus
type Proposal struct {
	ID        string                 `json:"id"`
	Type      ConsensusType          `json:"type"`
	Proposer  string                 `json:"proposer"` // SID of proposing agent
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
}

// Vote represents a vote on a proposal
type Vote struct {
	AgentSID   string    `json:"agent_sid"`
	ProposalID string    `json:"proposal_id"`
	Value      bool      `json:"value"` // true = accept, false = reject
	Reason     string    `json:"reason,omitempty"`
	Signature  []byte    `json:"signature,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// ConsensusRound represents a single round of consensus
type ConsensusRound struct {
	Proposal  *Proposal          `json:"proposal"`
	Votes     map[string]*Vote   `json:"votes"` // SID -> Vote
	Threshold float64            `json:"threshold"`
	Timeout   time.Duration      `json:"timeout"`
	StartedAt time.Time          `json:"started_at"`
	Result    string             `json:"result"` // "pending", "accepted", "rejected", "timeout"
}

// ConsensusEngine implements PBFT-style consensus
type ConsensusEngine struct {
	mu sync.RWMutex

	rounds    map[string]*ConsensusRound // ProposalID -> Round
	threshold float64                     // Consensus threshold (e.g., 0.67 for 2/3)
	timeout   time.Duration

	// Callbacks
	onAccept func(*Proposal)
	onReject func(*Proposal)
}

// NewConsensusEngine creates a new consensus engine
func NewConsensusEngine(threshold float64) *ConsensusEngine {
	return &ConsensusEngine{
		rounds:    make(map[string]*ConsensusRound),
		threshold: threshold,
		timeout:   30 * time.Second,
	}
}

// SetThreshold sets the consensus threshold
func (c *ConsensusEngine) SetThreshold(threshold float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.threshold = threshold
}

// SetTimeout sets the consensus timeout
func (c *ConsensusEngine) SetTimeout(timeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timeout = timeout
}

// OnAccept sets the callback for accepted proposals
func (c *ConsensusEngine) OnAccept(callback func(*Proposal)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onAccept = callback
}

// OnReject sets the callback for rejected proposals
func (c *ConsensusEngine) OnReject(callback func(*Proposal)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onReject = callback
}

// Propose starts a new consensus round
func (c *ConsensusEngine) Propose(ctx context.Context, proposerSID string, cType ConsensusType, data map[string]interface{}) (*ConsensusRound, error) {
	proposal := &Proposal{
		ID:        uuid.New().String(),
		Type:      cType,
		Proposer:  proposerSID,
		Data:      data,
		CreatedAt: time.Now(),
	}

	c.mu.Lock()
	round := &ConsensusRound{
		Proposal:  proposal,
		Votes:     make(map[string]*Vote),
		Threshold: c.threshold,
		Timeout:   c.timeout,
		StartedAt: time.Now(),
		Result:    "pending",
	}
	c.rounds[proposal.ID] = round
	c.mu.Unlock()

	// Proposer automatically votes yes
	c.SubmitVote(Vote{
		AgentSID:   proposerSID,
		ProposalID: proposal.ID,
		Value:      true,
		Reason:     "proposer",
		Timestamp:  time.Now(),
	})

	return round, nil
}

// SubmitVote submits a vote for a proposal
func (c *ConsensusEngine) SubmitVote(vote Vote) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	round, ok := c.rounds[vote.ProposalID]
	if !ok {
		return errors.New("proposal not found")
	}

	if round.Result != "pending" {
		return errors.New("consensus already reached")
	}

	vote.Timestamp = time.Now()
	round.Votes[vote.AgentSID] = &vote

	return nil
}

// CheckConsensus checks if consensus has been reached for a proposal
func (c *ConsensusEngine) CheckConsensus(proposalID string, totalVoters int) (bool, string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	round, ok := c.rounds[proposalID]
	if !ok {
		return false, "not_found"
	}

	if round.Result != "pending" {
		return round.Result == "accepted", round.Result
	}

	// Check timeout
	if time.Since(round.StartedAt) > round.Timeout {
		round.Result = "timeout"
		if c.onReject != nil {
			go c.onReject(round.Proposal)
		}
		return false, "timeout"
	}

	// Count votes
	accepts := 0
	rejects := 0
	for _, vote := range round.Votes {
		if vote.Value {
			accepts++
		} else {
			rejects++
		}
	}

	// Check if threshold reached
	requiredVotes := int(float64(totalVoters) * round.Threshold)
	if requiredVotes < 1 {
		requiredVotes = 1
	}

	if accepts >= requiredVotes {
		round.Result = "accepted"
		if c.onAccept != nil {
			go c.onAccept(round.Proposal)
		}
		return true, "accepted"
	}

	// Check if rejection is certain (not enough remaining votes could change outcome)
	remainingVotes := totalVoters - len(round.Votes)
	if accepts+remainingVotes < requiredVotes {
		round.Result = "rejected"
		if c.onReject != nil {
			go c.onReject(round.Proposal)
		}
		return false, "rejected"
	}

	return false, "pending"
}

// WaitForConsensus waits for consensus to be reached
func (c *ConsensusEngine) WaitForConsensus(ctx context.Context, proposalID string, totalVoters int) (bool, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-ticker.C:
			reached, result := c.CheckConsensus(proposalID, totalVoters)
			if result != "pending" {
				if result == "timeout" {
					return false, ErrConsensusTimeout
				}
				return reached, nil
			}
		}
	}
}

// GetRound returns a consensus round
func (c *ConsensusEngine) GetRound(proposalID string) *ConsensusRound {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.rounds[proposalID]
}

// GetAllRounds returns all consensus rounds
func (c *ConsensusEngine) GetAllRounds() []*ConsensusRound {
	c.mu.RLock()
	defer c.mu.RUnlock()

	rounds := make([]*ConsensusRound, 0, len(c.rounds))
	for _, round := range c.rounds {
		rounds = append(rounds, round)
	}
	return rounds
}

// CleanupOldRounds removes old completed rounds
func (c *ConsensusEngine) CleanupOldRounds(maxAge time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, round := range c.rounds {
		if round.Result != "pending" && time.Since(round.StartedAt) > maxAge {
			delete(c.rounds, id)
		}
	}
}

// Stats returns consensus engine statistics
type ConsensusStats struct {
	PendingRounds  int
	AcceptedRounds int
	RejectedRounds int
	TimeoutRounds  int
}

// Stats returns current consensus statistics
func (c *ConsensusEngine) Stats() ConsensusStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := ConsensusStats{}
	for _, round := range c.rounds {
		switch round.Result {
		case "pending":
			stats.PendingRounds++
		case "accepted":
			stats.AcceptedRounds++
		case "rejected":
			stats.RejectedRounds++
		case "timeout":
			stats.TimeoutRounds++
		}
	}
	return stats
}
