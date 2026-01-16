package collective

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// CollectiveMemory represents shared memory across the collective
type CollectiveMemory struct {
	mu sync.RWMutex

	// Knowledge Graph
	knowledgeGraph *KnowledgeGraph

	// Episodic Memory
	episodes []CollectiveEpisode

	// Working Memory - active contexts
	activeContexts map[string]*SharedContext

	// Semantic Memory - concepts and embeddings
	concepts map[string]*Concept
}

// NewCollectiveMemory creates a new collective memory
func NewCollectiveMemory() *CollectiveMemory {
	return &CollectiveMemory{
		knowledgeGraph: NewKnowledgeGraph(),
		episodes:       make([]CollectiveEpisode, 0),
		activeContexts: make(map[string]*SharedContext),
		concepts:       make(map[string]*Concept),
	}
}

// CollectiveEpisode represents a memorable collective event
type CollectiveEpisode struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`         // "task_completed", "agent_joined", "consensus", etc.
	Participants []string               `json:"participants"` // Agent SIDs involved
	Content      string                 `json:"content"`
	Context      map[string]interface{} `json:"context"`
	Timestamp    time.Time              `json:"timestamp"`
	Salience     float64                `json:"salience"` // 0.0-1.0 importance
}

// SharedContext represents an active shared context
type SharedContext struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Data         map[string]interface{} `json:"data"`
	Contributors []string               `json:"contributors"` // Agent SIDs
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	TTL          time.Duration          `json:"ttl"`
}

// Concept represents a semantic concept
type Concept struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Embedding   []float64 `json:"embedding,omitempty"` // Vector embedding
	Relations   []string  `json:"relations"`           // Related concept IDs
	CreatedBy   string    `json:"created_by"`          // Agent SID
	CreatedAt   time.Time `json:"created_at"`
}

// KnowledgeGraph represents the collective knowledge graph
type KnowledgeGraph struct {
	nodes map[string]*KnowledgeNode
	edges map[string][]*KnowledgeEdge
}

// KnowledgeNode represents a node in the knowledge graph
type KnowledgeNode struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  time.Time              `json:"created_at"`
}

// KnowledgeEdge represents an edge in the knowledge graph
type KnowledgeEdge struct {
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	Type      string                 `json:"type"`
	Weight    float64                `json:"weight"`
	Properties map[string]interface{} `json:"properties"`
}

// NewKnowledgeGraph creates a new knowledge graph
func NewKnowledgeGraph() *KnowledgeGraph {
	return &KnowledgeGraph{
		nodes: make(map[string]*KnowledgeNode),
		edges: make(map[string][]*KnowledgeEdge),
	}
}

// AddNode adds a node to the knowledge graph
func (kg *KnowledgeGraph) AddNode(node *KnowledgeNode) {
	kg.nodes[node.ID] = node
}

// AddEdge adds an edge to the knowledge graph
func (kg *KnowledgeGraph) AddEdge(edge *KnowledgeEdge) {
	kg.edges[edge.From] = append(kg.edges[edge.From], edge)
}

// GetNode returns a node by ID
func (kg *KnowledgeGraph) GetNode(id string) *KnowledgeNode {
	return kg.nodes[id]
}

// GetEdges returns edges from a node
func (kg *KnowledgeGraph) GetEdges(nodeID string) []*KnowledgeEdge {
	return kg.edges[nodeID]
}

// Contribute adds a memory contribution from an agent
func (m *CollectiveMemory) Contribute(agentSID string, content string, context map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	episode := CollectiveEpisode{
		ID:           uuid.New().String(),
		Type:         "contribution",
		Participants: []string{agentSID},
		Content:      content,
		Context:      context,
		Timestamp:    time.Now(),
		Salience:     0.5,
	}

	m.episodes = append(m.episodes, episode)

	// Keep episodes bounded
	if len(m.episodes) > 1000 {
		m.episodes = m.episodes[len(m.episodes)-1000:]
	}
}

// Query searches collective memory
func (m *CollectiveMemory) Query(query string) []CollectiveEpisode {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Simple substring search - in production, use vector similarity
	var results []CollectiveEpisode
	for _, ep := range m.episodes {
		if containsIgnoreCase(ep.Content, query) {
			results = append(results, ep)
		}
	}

	return results
}

// CreateContext creates a new shared context
func (m *CollectiveMemory) CreateContext(name string, creator string, ttl time.Duration) *SharedContext {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx := &SharedContext{
		ID:           uuid.New().String(),
		Name:         name,
		Data:         make(map[string]interface{}),
		Contributors: []string{creator},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		TTL:          ttl,
	}

	m.activeContexts[ctx.ID] = ctx
	return ctx
}

// GetContext returns a shared context
func (m *CollectiveMemory) GetContext(id string) *SharedContext {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.activeContexts[id]
}

// UpdateContext updates a shared context
func (m *CollectiveMemory) UpdateContext(id string, agentSID string, key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx, ok := m.activeContexts[id]
	if !ok {
		return
	}

	ctx.Data[key] = value
	ctx.UpdatedAt = time.Now()

	// Add contributor if not already present
	found := false
	for _, c := range ctx.Contributors {
		if c == agentSID {
			found = true
			break
		}
	}
	if !found {
		ctx.Contributors = append(ctx.Contributors, agentSID)
	}
}

// AddConcept adds a semantic concept
func (m *CollectiveMemory) AddConcept(name, description string, creator string) *Concept {
	m.mu.Lock()
	defer m.mu.Unlock()

	concept := &Concept{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Relations:   make([]string, 0),
		CreatedBy:   creator,
		CreatedAt:   time.Now(),
	}

	m.concepts[concept.ID] = concept
	return concept
}

// GetConcept returns a concept by ID
func (m *CollectiveMemory) GetConcept(id string) *Concept {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.concepts[id]
}

// LinkConcepts creates a relation between concepts
func (m *CollectiveMemory) LinkConcepts(conceptID1, conceptID2 string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c1, ok := m.concepts[conceptID1]; ok {
		c1.Relations = append(c1.Relations, conceptID2)
	}
	if c2, ok := m.concepts[conceptID2]; ok {
		c2.Relations = append(c2.Relations, conceptID1)
	}
}

// AddKnowledge adds knowledge to the graph
func (m *CollectiveMemory) AddKnowledge(nodeType, label string, properties map[string]interface{}) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	node := &KnowledgeNode{
		ID:         uuid.New().String(),
		Type:       nodeType,
		Label:      label,
		Properties: properties,
		CreatedAt:  time.Now(),
	}

	m.knowledgeGraph.AddNode(node)
	return node.ID
}

// ConnectKnowledge connects two knowledge nodes
func (m *CollectiveMemory) ConnectKnowledge(fromID, toID, edgeType string, weight float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	edge := &KnowledgeEdge{
		From:   fromID,
		To:     toID,
		Type:   edgeType,
		Weight: weight,
	}

	m.knowledgeGraph.AddEdge(edge)
}

// CleanupExpiredContexts removes expired contexts
func (m *CollectiveMemory) CleanupExpiredContexts() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, ctx := range m.activeContexts {
		if ctx.TTL > 0 && now.Sub(ctx.CreatedAt) > ctx.TTL {
			delete(m.activeContexts, id)
		}
	}
}

// Stats returns memory statistics
type MemoryStats struct {
	EpisodeCount   int
	ContextCount   int
	ConceptCount   int
	KnowledgeNodes int
}

// Stats returns current memory statistics
func (m *CollectiveMemory) Stats() MemoryStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return MemoryStats{
		EpisodeCount:   len(m.episodes),
		ContextCount:   len(m.activeContexts),
		ConceptCount:   len(m.concepts),
		KnowledgeNodes: len(m.knowledgeGraph.nodes),
	}
}

// containsIgnoreCase is a simple case-insensitive substring search
func containsIgnoreCase(s, substr string) bool {
	// Simple implementation - in production use proper case folding
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) > len(substr))
}
