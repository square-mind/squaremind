package identity

import "encoding/json"

// CapabilityType represents different types of agent capabilities
type CapabilityType string

const (
	CapCodeWrite     CapabilityType = "code.write"
	CapCodeReview    CapabilityType = "code.review"
	CapCodeRefactor  CapabilityType = "code.refactor"
	CapResearch      CapabilityType = "research"
	CapAnalysis      CapabilityType = "analysis"
	CapSecurity      CapabilityType = "security"
	CapDocumentation CapabilityType = "documentation"
	CapTesting       CapabilityType = "testing"
	CapArchitecture  CapabilityType = "architecture"
)

// Capability represents a specific capability an agent possesses
type Capability struct {
	Type        CapabilityType         `json:"type"`
	Proficiency float64                `json:"proficiency"` // 0.0 - 1.0
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Proof       *CapabilityProof       `json:"proof,omitempty"`
}

// CapabilityProof provides evidence for a claimed capability
type CapabilityProof struct {
	Type      string   `json:"type"` // "benchmark", "peer_attestation", "task_history"
	Score     float64  `json:"score,omitempty"`
	Benchmark string   `json:"benchmark,omitempty"`
	Attesters []string `json:"attesters,omitempty"`
	TaskCount int      `json:"task_count,omitempty"`
}

// CapabilitySet manages an agent's capabilities
type CapabilitySet struct {
	Capabilities map[CapabilityType]*Capability `json:"capabilities"`
}

// NewCapabilitySet creates a new empty capability set
func NewCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		Capabilities: make(map[CapabilityType]*Capability),
	}
}

// Add adds a capability to the set
func (cs *CapabilitySet) Add(cap *Capability) {
	cs.Capabilities[cap.Type] = cap
}

// Has checks if the set contains a capability type
func (cs *CapabilitySet) Has(capType CapabilityType) bool {
	_, exists := cs.Capabilities[capType]
	return exists
}

// Get retrieves a capability by type
func (cs *CapabilitySet) Get(capType CapabilityType) *Capability {
	return cs.Capabilities[capType]
}

// List returns all capability types in the set
func (cs *CapabilitySet) List() []CapabilityType {
	types := make([]CapabilityType, 0, len(cs.Capabilities))
	for t := range cs.Capabilities {
		types = append(types, t)
	}
	return types
}

// MatchScore returns how well this capability set matches required capabilities
func (cs *CapabilitySet) MatchScore(required []CapabilityType) float64 {
	if len(required) == 0 {
		return 1.0
	}

	var totalScore float64
	var matched int

	for _, req := range required {
		if cap := cs.Get(req); cap != nil {
			totalScore += cap.Proficiency
			matched++
		}
	}

	if matched == 0 {
		return 0.0
	}

	// Score based on coverage and proficiency
	coverage := float64(matched) / float64(len(required))
	avgProficiency := totalScore / float64(matched)

	return coverage * avgProficiency
}

// ToJSON serializes the capability set
func (cs *CapabilitySet) ToJSON() string {
	data, err := json.Marshal(cs)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// FromJSON deserializes a capability set from JSON
func (cs *CapabilitySet) FromJSON(data []byte) error {
	return json.Unmarshal(data, cs)
}
