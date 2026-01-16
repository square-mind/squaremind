package coordination

import (
	"context"
	"testing"
	"time"
)

func TestNewConsensusEngine(t *testing.T) {
	ce := NewConsensusEngine(0.67)
	if ce == nil {
		t.Fatal("NewConsensusEngine returned nil")
	}
}

func TestConsensusEngine_Propose(t *testing.T) {
	ce := NewConsensusEngine(0.67)

	ctx := context.Background()
	round, err := ce.Propose(ctx, "agent-1", ConsensusTypeTaskAssignment, map[string]interface{}{
		"task_id": "task-123",
	})

	if err != nil {
		t.Fatalf("Propose failed: %v", err)
	}

	if round == nil {
		t.Fatal("Propose returned nil round")
	}

	if round.Proposal.Proposer != "agent-1" {
		t.Errorf("Expected proposer 'agent-1', got '%s'", round.Proposal.Proposer)
	}

	if round.Result != "pending" {
		t.Errorf("Expected result 'pending', got '%s'", round.Result)
	}

	// Proposer should have auto-voted
	if len(round.Votes) != 1 {
		t.Errorf("Expected 1 vote (proposer), got %d", len(round.Votes))
	}
}

func TestConsensusEngine_SubmitVote(t *testing.T) {
	ce := NewConsensusEngine(0.67)

	ctx := context.Background()
	round, _ := ce.Propose(ctx, "agent-1", ConsensusTypeTaskAssignment, nil)

	err := ce.SubmitVote(Vote{
		AgentSID:   "agent-2",
		ProposalID: round.Proposal.ID,
		Value:      true,
		Reason:     "approved",
	})

	if err != nil {
		t.Fatalf("SubmitVote failed: %v", err)
	}

	updatedRound := ce.GetRound(round.Proposal.ID)
	if len(updatedRound.Votes) != 2 {
		t.Errorf("Expected 2 votes, got %d", len(updatedRound.Votes))
	}
}

func TestConsensusEngine_CheckConsensus_Accepted(t *testing.T) {
	ce := NewConsensusEngine(0.67)

	ctx := context.Background()
	round, _ := ce.Propose(ctx, "agent-1", ConsensusTypeTaskAssignment, nil)

	// Add more votes to reach threshold (67% of 3 = 2 votes needed)
	ce.SubmitVote(Vote{AgentSID: "agent-2", ProposalID: round.Proposal.ID, Value: true})

	reached, result := ce.CheckConsensus(round.Proposal.ID, 3)

	if !reached {
		t.Error("Expected consensus to be reached")
	}

	if result != "accepted" {
		t.Errorf("Expected result 'accepted', got '%s'", result)
	}
}

func TestConsensusEngine_CheckConsensus_Rejected(t *testing.T) {
	ce := NewConsensusEngine(0.67)

	ctx := context.Background()
	round, _ := ce.Propose(ctx, "agent-1", ConsensusTypeTaskAssignment, nil)

	// Add reject votes
	ce.SubmitVote(Vote{AgentSID: "agent-2", ProposalID: round.Proposal.ID, Value: false})
	ce.SubmitVote(Vote{AgentSID: "agent-3", ProposalID: round.Proposal.ID, Value: false})

	reached, result := ce.CheckConsensus(round.Proposal.ID, 3)

	if reached {
		t.Error("Expected consensus to not be reached")
	}

	if result != "rejected" {
		t.Errorf("Expected result 'rejected', got '%s'", result)
	}
}

func TestConsensusEngine_SetThreshold(t *testing.T) {
	ce := NewConsensusEngine(0.67)
	ce.SetThreshold(0.51)

	// Verify threshold change works by needing fewer votes
	ctx := context.Background()
	round, _ := ce.Propose(ctx, "agent-1", ConsensusTypeTaskAssignment, nil)

	// With 51% threshold and 2 voters, 1 vote (51%) should be enough
	// Proposer already voted yes
	reached, result := ce.CheckConsensus(round.Proposal.ID, 2)

	if !reached || result != "accepted" {
		t.Errorf("Expected acceptance with 51%% threshold, got reached=%v, result=%s", reached, result)
	}
}

func TestConsensusEngine_SetTimeout(t *testing.T) {
	ce := NewConsensusEngine(0.67)
	ce.SetTimeout(1 * time.Second)
	// No assertion, just ensure no panic
}

func TestConsensusEngine_Stats(t *testing.T) {
	ce := NewConsensusEngine(0.67)

	ctx := context.Background()
	ce.Propose(ctx, "agent-1", ConsensusTypeTaskAssignment, nil)

	stats := ce.Stats()

	if stats.PendingRounds != 1 {
		t.Errorf("Expected 1 pending round, got %d", stats.PendingRounds)
	}
}

func TestConsensusEngine_GetAllRounds(t *testing.T) {
	ce := NewConsensusEngine(0.67)

	ctx := context.Background()
	ce.Propose(ctx, "agent-1", ConsensusTypeTaskAssignment, nil)
	ce.Propose(ctx, "agent-2", ConsensusTypeAgentSpawn, nil)

	rounds := ce.GetAllRounds()

	if len(rounds) != 2 {
		t.Errorf("Expected 2 rounds, got %d", len(rounds))
	}
}
