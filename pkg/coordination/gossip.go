package coordination

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MessageType represents types of gossip messages
type MessageType string

const (
	MsgAgentJoined   MessageType = "agent_joined"
	MsgAgentLeft     MessageType = "agent_left"
	MsgTaskAvailable MessageType = "task_available"
	MsgTaskBid       MessageType = "task_bid"
	MsgTaskAssigned  MessageType = "task_assigned"
	MsgTaskCompleted MessageType = "task_completed"
	MsgHeartbeat     MessageType = "heartbeat"
	MsgConsensus     MessageType = "consensus"
)

// Message represents a gossip message
type Message struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	From      string      `json:"from"`    // Sender SID
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
	TTL       int         `json:"ttl"` // Hops remaining
}

// GossipProtocol implements epidemic-style message propagation
type GossipProtocol struct {
	mu sync.RWMutex

	peers    map[string]bool // SID -> active
	seen     map[string]bool // Message ID -> seen
	handlers map[MessageType][]MessageHandler

	fanout   int           // Number of peers to forward to
	interval time.Duration // Gossip interval

	msgChan chan Message
}

// MessageHandler handles incoming gossip messages
type MessageHandler func(msg Message)

// NewGossipProtocol creates a new gossip protocol instance
func NewGossipProtocol() *GossipProtocol {
	return &GossipProtocol{
		peers:    make(map[string]bool),
		seen:     make(map[string]bool),
		handlers: make(map[MessageType][]MessageHandler),
		fanout:   3,
		interval: 100 * time.Millisecond,
		msgChan:  make(chan Message, 1000),
	}
}

// AddPeer adds a peer to the gossip network
func (g *GossipProtocol) AddPeer(sid string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.peers[sid] = true
}

// RemovePeer removes a peer from the gossip network
func (g *GossipProtocol) RemovePeer(sid string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.peers, sid)
}

// GetPeers returns list of active peers
func (g *GossipProtocol) GetPeers() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	peers := make([]string, 0, len(g.peers))
	for sid := range g.peers {
		peers = append(peers, sid)
	}
	return peers
}

// PeerCount returns the number of peers
func (g *GossipProtocol) PeerCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.peers)
}

// OnMessage registers a handler for a message type
func (g *GossipProtocol) OnMessage(msgType MessageType, handler MessageHandler) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.handlers[msgType] = append(g.handlers[msgType], handler)
}

// Broadcast sends a message to the network
func (g *GossipProtocol) Broadcast(msg Message) {
	msg.ID = uuid.New().String()
	msg.Timestamp = time.Now()
	if msg.TTL == 0 {
		msg.TTL = 10 // Max hops
	}

	select {
	case g.msgChan <- msg:
	default:
		// Channel full, drop message
	}
}

// Start begins the gossip protocol
func (g *GossipProtocol) Start(ctx context.Context) {
	go g.processMessages(ctx)
	go g.runGossipLoop(ctx)
}

// processMessages handles incoming messages
func (g *GossipProtocol) processMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-g.msgChan:
			g.handleMessage(msg)
		}
	}
}

// handleMessage processes a single message
func (g *GossipProtocol) handleMessage(msg Message) {
	g.mu.Lock()

	// Check if already seen
	if g.seen[msg.ID] {
		g.mu.Unlock()
		return
	}
	g.seen[msg.ID] = true

	// Get handlers
	handlers := g.handlers[msg.Type]
	g.mu.Unlock()

	// Execute handlers
	for _, h := range handlers {
		h(msg)
	}

	// Forward to random peers if TTL > 0
	if msg.TTL > 0 {
		msg.TTL--
		g.forward(msg)
	}
}

// forward sends message to random subset of peers
func (g *GossipProtocol) forward(msg Message) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Get list of peers (excluding sender)
	var candidates []string
	for sid := range g.peers {
		if sid != msg.From {
			candidates = append(candidates, sid)
		}
	}

	// Select random subset
	if len(candidates) <= g.fanout {
		// Send to all
		for _, sid := range candidates {
			g.sendTo(sid, msg)
		}
	} else {
		// Random selection
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		for i := 0; i < g.fanout; i++ {
			g.sendTo(candidates[i], msg)
		}
	}
}

// sendTo sends a message to a specific peer
func (g *GossipProtocol) sendTo(sid string, msg Message) {
	// In a real implementation, this would use network transport
	// For now, we just re-queue (simulating local delivery)
	go func() {
		select {
		case g.msgChan <- msg:
		default:
			// Channel full
		}
	}()
}

// runGossipLoop periodically cleans up and maintains state
func (g *GossipProtocol) runGossipLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			g.cleanup()
		}
	}
}

// cleanup removes old seen messages
func (g *GossipProtocol) cleanup() {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Keep seen map from growing unbounded
	// In production, use a time-based eviction
	if len(g.seen) > 10000 {
		g.seen = make(map[string]bool)
	}
}

// SetFanout sets the fanout parameter
func (g *GossipProtocol) SetFanout(fanout int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.fanout = fanout
}

// Stats returns gossip protocol statistics
type GossipStats struct {
	PeerCount    int
	SeenMessages int
	HandlerCount int
}

// Stats returns current gossip statistics
func (g *GossipProtocol) Stats() GossipStats {
	g.mu.RLock()
	defer g.mu.RUnlock()

	handlerCount := 0
	for _, handlers := range g.handlers {
		handlerCount += len(handlers)
	}

	return GossipStats{
		PeerCount:    len(g.peers),
		SeenMessages: len(g.seen),
		HandlerCount: handlerCount,
	}
}
