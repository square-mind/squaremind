package identity

import (
	"crypto/ed25519"
	"encoding/json"
	"time"
)

// SignedAction represents an action signed by an agent
type SignedAction struct {
	Action    Action       `json:"action"`
	AgentSID  string       `json:"agent_sid"`
	Signature []byte       `json:"signature"`
	Timestamp time.Time    `json:"timestamp"`
	Proof     *ActionProof `json:"proof,omitempty"`
}

// Action represents any agent action that can be signed
type Action struct {
	Type    string                 `json:"type"`
	Target  string                 `json:"target,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// ActionProof provides evidence for an action
type ActionProof struct {
	Type     string        `json:"type"` // "capability", "delegation", "consensus"
	Evidence ProofEvidence `json:"evidence"`
}

// ProofEvidence contains the actual proof data
type ProofEvidence struct {
	CapabilityType  CapabilityType `json:"capability_type,omitempty"`
	DelegatorSID    string         `json:"delegator_sid,omitempty"`
	ConsensusRound  int            `json:"consensus_round,omitempty"`
	Votes           int            `json:"votes,omitempty"`
	RequiredVotes   int            `json:"required_votes,omitempty"`
}

// NewSignedAction creates a new signed action
func NewSignedAction(identity *SquaremindIdentity, action Action) (*SignedAction, error) {
	// Serialize action for signing
	actionData, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}

	// Sign the action
	signature := identity.Sign(actionData)

	return &SignedAction{
		Action:    action,
		AgentSID:  identity.SID,
		Signature: signature,
		Timestamp: time.Now(),
	}, nil
}

// Verify verifies the signature of a signed action
func (sa *SignedAction) Verify(publicKey ed25519.PublicKey) bool {
	actionData, err := json.Marshal(sa.Action)
	if err != nil {
		return false
	}

	return ed25519.Verify(publicKey, actionData, sa.Signature)
}

// DelegationProof represents a proof of delegated authority
type DelegationProof struct {
	DelegatorSID   string    `json:"delegator_sid"`
	DelegateSID    string    `json:"delegate_sid"`
	Capability     CapabilityType `json:"capability"`
	ExpiresAt      time.Time `json:"expires_at"`
	Signature      []byte    `json:"signature"`
}

// NewDelegationProof creates a new delegation proof
func NewDelegationProof(delegator *SquaremindIdentity, delegateSID string, capability CapabilityType, duration time.Duration) (*DelegationProof, error) {
	proof := &DelegationProof{
		DelegatorSID: delegator.SID,
		DelegateSID:  delegateSID,
		Capability:   capability,
		ExpiresAt:    time.Now().Add(duration),
	}

	// Sign the delegation
	data, err := json.Marshal(struct {
		DelegateSID string
		Capability  CapabilityType
		ExpiresAt   time.Time
	}{delegateSID, capability, proof.ExpiresAt})
	if err != nil {
		return nil, err
	}

	proof.Signature = delegator.Sign(data)
	return proof, nil
}

// IsValid checks if the delegation proof is still valid
func (dp *DelegationProof) IsValid() bool {
	return time.Now().Before(dp.ExpiresAt)
}

// ConsensusProof represents proof of collective consensus
type ConsensusProof struct {
	Round         int                   `json:"round"`
	Proposal      string                `json:"proposal"`
	Votes         map[string]Vote       `json:"votes"` // SID -> Vote
	Threshold     float64               `json:"threshold"`
	Result        string                `json:"result"` // "accepted", "rejected"
	Timestamp     time.Time             `json:"timestamp"`
}

// Vote represents an agent's vote in consensus
type Vote struct {
	AgentSID  string    `json:"agent_sid"`
	Value     bool      `json:"value"` // true = accept, false = reject
	Signature []byte    `json:"signature"`
	Timestamp time.Time `json:"timestamp"`
}

// IsAccepted returns whether the consensus was accepted
func (cp *ConsensusProof) IsAccepted() bool {
	return cp.Result == "accepted"
}

// VoteCount returns the number of accept/reject votes
func (cp *ConsensusProof) VoteCount() (accepts int, rejects int) {
	for _, v := range cp.Votes {
		if v.Value {
			accepts++
		} else {
			rejects++
		}
	}
	return
}
