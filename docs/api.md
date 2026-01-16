# Squaremind API Reference

Many Agents. One Mind.

## Go API

### Package: identity

#### SquaremindIdentity

```go
type SquaremindIdentity struct {
    SID        string
    PublicKey  ed25519.PublicKey
    PrivateKey ed25519.PrivateKey
    Name       string
    CreatedAt  time.Time
    ParentSID  string
    Generation int
}

// Create new identity
func NewSquaremindIdentity(name string, parentSID string) (*SquaremindIdentity, error)

// Sign data
func (s *SquaremindIdentity) Sign(data []byte) []byte

// Verify signature
func (s *SquaremindIdentity) Verify(data, signature []byte) bool

// Get hex-encoded public key
func (s *SquaremindIdentity) PublicKeyHex() string
```

#### Capabilities

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

type CapabilitySet struct {
    Capabilities map[CapabilityType]*Capability
}

func NewCapabilitySet() *CapabilitySet
func (cs *CapabilitySet) Add(cap *Capability)
func (cs *CapabilitySet) Has(capType CapabilityType) bool
func (cs *CapabilitySet) Get(capType CapabilityType) *Capability
func (cs *CapabilitySet) MatchScore(required []CapabilityType) float64
```

### Package: agent

#### Agent

```go
type Agent struct {
    Identity     *identity.SquaremindIdentity
    Capabilities *identity.CapabilitySet
    Provider     llm.Provider
    Model        string
    State        AgentState
    Reputation   *Reputation
    Memory       *AgentMemory
}

type AgentConfig struct {
    Name         string
    Capabilities []identity.CapabilityType
    Provider     llm.Provider
    Model        string
    ParentSID    string
}

func NewAgent(cfg AgentConfig) (*Agent, error)
func (a *Agent) Start(ctx context.Context) error
func (a *Agent) Stop()
func (a *Agent) SubmitTask(task *Task)
func (a *Agent) GetResults() <-chan *TaskResult
func (a *Agent) GetState() AgentState
```

#### Task

```go
type Task struct {
    ID           string
    Description  string
    Requirements string
    Complexity   string
    Required     []identity.CapabilityType
    Deadline     time.Time
    Reward       float64
    Status       TaskStatus
    AssignedTo   string
    CreatedAt    time.Time
}

func NewTask(description string, required []identity.CapabilityType) *Task
func (t *Task) WithComplexity(complexity string) *Task
func (t *Task) WithDeadline(deadline time.Time) *Task
func (t *Task) WithReward(reward float64) *Task
```

### Package: collective

#### Collective

```go
type Collective struct {
    Name string
    ID   string
}

type CollectiveConfig struct {
    MinAgents          int
    MaxAgents          int
    ConsensusThreshold float64
    ReputationDecay    float64
}

func NewCollective(name string, cfg CollectiveConfig) *Collective
func (c *Collective) Join(a *agent.Agent) error
func (c *Collective) Leave(sid string) error
func (c *Collective) Submit(task *agent.Task) (*agent.TaskResult, error)
func (c *Collective) SubmitAsync(task *agent.Task) (string, error)
func (c *Collective) GetAgents() []*agent.Agent
func (c *Collective) Size() int
func (c *Collective) Start(ctx context.Context) error
func (c *Collective) Stop()
func (c *Collective) Stats() CollectiveStats
```

### Package: coordination

#### GossipProtocol

```go
type GossipProtocol struct {}

func NewGossipProtocol() *GossipProtocol
func (g *GossipProtocol) AddPeer(sid string)
func (g *GossipProtocol) RemovePeer(sid string)
func (g *GossipProtocol) Broadcast(msg Message)
func (g *GossipProtocol) OnMessage(msgType MessageType, handler MessageHandler)
func (g *GossipProtocol) Start(ctx context.Context)
```

#### TaskMarket

```go
type TaskMarket struct {}

func NewTaskMarket() *TaskMarket
func (m *TaskMarket) ListTask(task *agent.Task) error
func (m *TaskMarket) SubmitBid(bid *Bid) error
func (m *TaskMarket) AssignTask(task, agents, reputation) (*TaskAssignment, error)
```

#### ConsensusEngine

```go
type ConsensusEngine struct {}

func NewConsensusEngine(threshold float64) *ConsensusEngine
func (c *ConsensusEngine) Propose(ctx, proposerSID, cType, data) (*ConsensusRound, error)
func (c *ConsensusEngine) SubmitVote(vote Vote) error
func (c *ConsensusEngine) CheckConsensus(proposalID string, totalVoters int) (bool, string)
func (c *ConsensusEngine) WaitForConsensus(ctx, proposalID, totalVoters) (bool, error)
```

### Package: llm

#### Provider Interface

```go
type Provider interface {
    Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
    Name() string
}

type CompletionRequest struct {
    Model       string
    Prompt      string
    MaxTokens   int
    Temperature float64
    Stop        []string
    System      string
}

type CompletionResponse struct {
    Content      string
    FinishReason string
    TokensUsed   int
}
```

#### Claude Provider

```go
func NewClaudeProvider(apiKey string) *ClaudeProvider
func (p *ClaudeProvider) WithModel(model string) *ClaudeProvider
func (p *ClaudeProvider) Complete(ctx, req) (*CompletionResponse, error)
```

#### OpenAI Provider

```go
func NewOpenAIProvider(apiKey string) *OpenAIProvider
func (p *OpenAIProvider) WithModel(model string) *OpenAIProvider
func (p *OpenAIProvider) Complete(ctx, req) (*CompletionResponse, error)
```

## CLI Reference

```bash
# Initialize a collective
sqm init <name> [--max-agents N] [--threshold F]

# Spawn an agent
sqm spawn <name> [-c capabilities] [-m model]

# Start the collective
sqm run

# Show status
sqm status

# Submit a task
sqm task submit <description> [-x complexity] [-r requires] [--async]

# List agents
sqm agent list

# Stop an agent
sqm agent stop <sid>

# Configure API keys
sqm config set api-key <key>
sqm config set openai-key <key>
```

## Learn More

- [GitHub](https://github.com/squaremind/squaremind)
- [Website](https://squaremind.xyz)
