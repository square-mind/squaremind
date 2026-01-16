# Squaremind Architecture

This document provides a detailed overview of the Squaremind architecture.

## Overview

Squaremind is built on a three-layer architecture designed for autonomous multi-agent coordination:

```
┌─────────────────────────────────────────────────────────────────┐
│                   COLLECTIVE MIND SUBSTRATE                      │
│                        (Layer 3)                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   Shared    │  │  Knowledge  │  │      Swarm Patterns     │  │
│  │   Memory    │  │    Graph    │  │  (stigmergy, quorum)    │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                   FAIR COORDINATION LAYER                        │
│                        (Layer 2)                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   Gossip    │  │    Task     │  │       Consensus         │  │
│  │  Protocol   │  │   Market    │  │   (PBFT, Reputation)    │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                   AGENT IDENTITY PROTOCOL                        │
│                        (Layer 1)                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   Ed25519   │  │ Capability  │  │      Reputation         │  │
│  │  Identity   │  │    Set      │  │        Score            │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Layer 1: Agent Identity Protocol (AIP)

The foundation layer that establishes cryptographic identity and trust primitives.

### Components

#### Squaremind Identity (SID)

Every agent has a unique cryptographic identity:

```go
type SquaremindIdentity struct {
    SID        string              // Unique identifier (UUID)
    PublicKey  ed25519.PublicKey   // 32-byte public key
    PrivateKey ed25519.PrivateKey  // 64-byte private key
    Name       string              // Human-readable name
    CreatedAt  time.Time           // Creation timestamp
    ParentSID  string              // Parent agent (for spawned agents)
    Generation int                 // Generation depth
}
```

**Key Features:**
- Ed25519 elliptic curve cryptography
- Deterministic key generation
- Hierarchical identity (parent-child relationships)
- Verifiable signatures

#### Capability Set

Agents declare their capabilities:

```go
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

type Capability struct {
    Type        CapabilityType
    Proficiency float64  // 0.0 - 1.0
    Metadata    map[string]interface{}
}
```

**Capability Matching:**
- Tasks declare required capabilities
- Agents are matched by capability coverage
- Proficiency levels affect matching scores

#### Reputation Score

Multi-dimensional reputation tracking:

```go
type Reputation struct {
    Overall        float64  // Weighted average
    Reliability    float64  // Task completion rate
    Quality        float64  // Output quality score
    Cooperation    float64  // Peer interaction score
    Honesty        float64  // Bid accuracy score
    TasksCompleted int
    TasksFailed    int
    LastActive     time.Time
}
```

**Reputation Dynamics:**
- Increases with successful task completion
- Decays over time (configurable rate)
- Affects task bidding priority
- Used as stake in task market

## Layer 2: Fair Coordination Layer (FCL)

The coordination layer that enables decentralized task allocation and consensus.

### Components

#### Gossip Protocol

Epidemic-style message propagation:

```go
type GossipConfig struct {
    Fanout     int           // Number of peers per round (default: 3)
    Interval   time.Duration // Gossip interval (default: 100ms)
    DefaultTTL int           // Message time-to-live (default: 5)
}

type Message struct {
    ID        string
    Type      MessageType
    From      string
    Payload   interface{}
    Timestamp time.Time
    TTL       int
    Signature []byte
}
```

**Message Types:**
- `TASK_ANNOUNCED` - New task available
- `BID_SUBMITTED` - Agent bid on task
- `TASK_ASSIGNED` - Task assigned to agent
- `TASK_COMPLETED` - Task finished
- `REPUTATION_UPDATE` - Reputation change
- `CONSENSUS_PROPOSE` - Consensus proposal
- `CONSENSUS_VOTE` - Consensus vote

**Properties:**
- Probabilistic delivery guarantees
- Duplicate detection via message ID
- TTL-based message expiration
- Cryptographically signed messages

#### Task Market

Decentralized task allocation:

```go
type Task struct {
    ID          string
    Description string
    Required    []CapabilityType
    Complexity  string  // low, medium, high
    Deadline    time.Time
    Reward      float64
    Status      TaskStatus
    AssignedTo  string
}

type Bid struct {
    AgentSID        string
    TaskID          string
    CapabilityScore float64
    ReputationStake float64
    EstimatedTime   time.Duration
}
```

**Allocation Algorithm:**
1. Task announced via gossip
2. Qualified agents submit bids
3. Bids scored by: `capability_match * 0.6 + reputation * 0.4`
4. Highest scoring agent assigned
5. Agent stakes reputation on completion

#### Consensus Engine

Byzantine fault-tolerant consensus:

```go
type ConsensusEngine struct {
    threshold float64  // Required vote percentage (default: 0.67)
    rounds    map[string]*ConsensusRound
}

type Vote struct {
    ProposalID string
    VoterSID   string
    Decision   VoteDecision  // approve, reject, abstain
    Signature  []byte
}
```

**Consensus Types:**
- `TASK_VALIDATION` - Validate task completion
- `AGENT_ADMISSION` - Approve new agent
- `AGENT_REMOVAL` - Remove misbehaving agent
- `PARAMETER_CHANGE` - Change collective parameters

**Process:**
1. Proposer creates signed proposal
2. Proposal gossiped to all agents
3. Agents vote within timeout
4. 2/3+ approval = consensus reached
5. Result gossiped with proof

## Layer 3: Collective Mind Substrate (CMS)

The emergent intelligence layer that enables collective learning and coordination.

### Components

#### Shared Memory

Distributed key-value store:

```go
type CollectiveMemory struct {
    shortTerm map[string]*MemoryEntry  // Recent, volatile
    longTerm  map[string]*MemoryEntry  // Persistent, important
    graph     *KnowledgeGraph          // Relationship graph
}

type MemoryEntry struct {
    Key       string
    Value     interface{}
    Source    string  // Agent SID
    Timestamp time.Time
    Relevance float64
    AccessCount int
}
```

**Features:**
- Short-term and long-term memory
- Relevance-based retention
- Access frequency tracking
- Automatic consolidation

#### Knowledge Graph

Relationship-based knowledge storage:

```go
type KnowledgeGraph struct {
    nodes map[string]*KnowledgeNode
    edges []*KnowledgeEdge
}

type KnowledgeNode struct {
    ID         string
    Type       string
    Properties map[string]interface{}
}

type KnowledgeEdge struct {
    From     string
    To       string
    Relation string
    Weight   float64
}
```

**Usage:**
- Store learned relationships
- Enable semantic queries
- Support reasoning chains

#### Swarm Patterns

Bio-inspired coordination patterns:

```go
type SwarmPattern interface {
    Name() string
    Apply(collective *Collective, context interface{}) error
}
```

**Implemented Patterns:**

1. **Stigmergy**: Indirect coordination through environment
   - Agents leave "pheromone trails" (metadata)
   - Other agents follow successful paths

2. **Quorum Sensing**: Density-dependent behavior
   - Behavior changes based on agent count
   - Enables scaling adaptations

3. **Chemotaxis**: Gradient-following behavior
   - Agents move toward "attractors" (high-value tasks)
   - Enables emergent load balancing

4. **Division of Labor**: Specialization
   - Agents specialize based on capability proficiency
   - Collective capability coverage optimized

## Data Flow

### Task Lifecycle

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Client  │────▶│  Market  │────▶│  Agent   │────▶│ Complete │
│  Submit  │     │  Bidding │     │ Execute  │     │  Verify  │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
     │                │                │                │
     ▼                ▼                ▼                ▼
  ┌──────┐        ┌──────┐        ┌──────┐        ┌──────┐
  │Gossip│        │Gossip│        │Memory│        │Gossip│
  │ Msg  │        │ Bids │        │Update│        │Result│
  └──────┘        └──────┘        └──────┘        └──────┘
```

### Consensus Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ Proposer │────▶│  Gossip  │────▶│  Voters  │────▶│  Tally   │
│  Create  │     │ Proposal │     │   Vote   │     │  Result  │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                                       │
                                       ▼
                                  ┌──────────┐
                                  │ Threshold│
                                  │  Check   │
                                  └──────────┘
```

## Package Structure

```
pkg/
├── identity/           # Layer 1: Identity primitives
│   ├── sid.go          # Squaremind Identity
│   ├── capabilities.go # Capability system
│   └── proofs.go       # Cryptographic proofs
│
├── agent/              # Agent runtime
│   ├── agent.go        # Agent implementation
│   ├── types.go        # Task, Result, Reputation
│   ├── runtime.go      # Agent runtime
│   └── lifecycle.go    # Lifecycle management
│
├── coordination/       # Layer 2: Coordination
│   ├── gossip.go       # Gossip protocol
│   ├── market.go       # Task market
│   ├── reputation.go   # Reputation registry
│   └── consensus.go    # Consensus engine
│
├── collective/         # Layer 3: Collective
│   ├── collective.go   # Collective management
│   ├── memory.go       # Shared memory
│   └── patterns.go     # Swarm patterns
│
└── llm/                # LLM integration
    ├── provider.go     # Provider interface
    ├── claude.go       # Claude provider
    └── openai.go       # OpenAI provider
```

## Security Model

### Threat Model

- **Sybil Attacks**: Mitigated by reputation staking
- **Eclipse Attacks**: Mitigated by random peer selection
- **Task Manipulation**: Mitigated by consensus validation
- **Identity Spoofing**: Prevented by Ed25519 signatures

### Trust Assumptions

- Ed25519 cryptography is secure
- Majority of agents are honest (>2/3 for BFT)
- Network is eventually consistent
- LLM providers are available

## Performance Characteristics

| Operation | Complexity | Typical Latency |
|-----------|------------|-----------------|
| Identity Creation | O(1) | <1ms |
| Gossip Round | O(fanout) | 100ms |
| Task Assignment | O(n agents) | <500ms |
| Consensus | O(n voters) | 1-5s |
| Memory Lookup | O(1) | <1ms |

## Future Considerations

1. **Persistent Storage**: Durable agent state
2. **Multi-Collective**: Cross-collective communication
3. **Sharding**: Horizontal scaling
4. **Formal Verification**: Protocol proofs
