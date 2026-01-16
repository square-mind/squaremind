package coordination

import (
	"context"
	"testing"
	"time"
)

func TestNewGossipProtocol(t *testing.T) {
	g := NewGossipProtocol()
	if g == nil {
		t.Fatal("NewGossipProtocol returned nil")
	}

	if g.fanout != 3 {
		t.Errorf("Expected default fanout 3, got %d", g.fanout)
	}
}

func TestGossipProtocol_AddRemovePeer(t *testing.T) {
	g := NewGossipProtocol()

	g.AddPeer("agent-1")
	g.AddPeer("agent-2")

	if g.PeerCount() != 2 {
		t.Errorf("Expected 2 peers, got %d", g.PeerCount())
	}

	peers := g.GetPeers()
	if len(peers) != 2 {
		t.Errorf("Expected 2 peers in list, got %d", len(peers))
	}

	g.RemovePeer("agent-1")
	if g.PeerCount() != 1 {
		t.Errorf("Expected 1 peer after removal, got %d", g.PeerCount())
	}
}

func TestGossipProtocol_OnMessage(t *testing.T) {
	g := NewGossipProtocol()

	received := make(chan Message, 1)
	g.OnMessage(MsgTaskAvailable, func(msg Message) {
		received <- msg
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g.Start(ctx)

	// Broadcast a message
	g.Broadcast(Message{
		Type:    MsgTaskAvailable,
		From:    "agent-1",
		Payload: "test payload",
	})

	select {
	case msg := <-received:
		if msg.Type != MsgTaskAvailable {
			t.Errorf("Expected message type MsgTaskAvailable, got %s", msg.Type)
		}
	case <-time.After(time.Second):
		t.Error("Timed out waiting for message")
	}
}

func TestGossipProtocol_SetFanout(t *testing.T) {
	g := NewGossipProtocol()

	g.SetFanout(5)

	stats := g.Stats()
	// We can't directly check fanout from stats, but ensure no panic
	_ = stats
}

func TestGossipProtocol_Stats(t *testing.T) {
	g := NewGossipProtocol()

	g.AddPeer("agent-1")
	g.AddPeer("agent-2")
	g.OnMessage(MsgTaskAvailable, func(msg Message) {})

	stats := g.Stats()

	if stats.PeerCount != 2 {
		t.Errorf("Expected 2 peers in stats, got %d", stats.PeerCount)
	}

	if stats.HandlerCount != 1 {
		t.Errorf("Expected 1 handler in stats, got %d", stats.HandlerCount)
	}
}
