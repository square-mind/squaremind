package identity

import (
	"testing"
)

func TestNewSquaremindIdentity(t *testing.T) {
	id, err := NewSquaremindIdentity("TestAgent", "")
	if err != nil {
		t.Fatalf("Failed to create identity: %v", err)
	}

	if id.Name != "TestAgent" {
		t.Errorf("Expected name 'TestAgent', got '%s'", id.Name)
	}

	if id.SID == "" {
		t.Error("SID should not be empty")
	}

	if len(id.PublicKey) != 32 {
		t.Errorf("Expected 32-byte public key, got %d bytes", len(id.PublicKey))
	}

	if len(id.PrivateKey) != 64 {
		t.Errorf("Expected 64-byte private key, got %d bytes", len(id.PrivateKey))
	}

	if id.Generation != 0 {
		t.Errorf("Expected generation 0 for root agent, got %d", id.Generation)
	}
}

func TestNewSquaremindIdentity_WithParent(t *testing.T) {
	parent, _ := NewSquaremindIdentity("Parent", "")
	child, err := NewSquaremindIdentity("Child", parent.SID)
	if err != nil {
		t.Fatalf("Failed to create child identity: %v", err)
	}

	if child.ParentSID != parent.SID {
		t.Errorf("Expected parent SID '%s', got '%s'", parent.SID, child.ParentSID)
	}

	if child.Generation != 1 {
		t.Errorf("Expected generation 1 for child agent, got %d", child.Generation)
	}
}

func TestSignAndVerify(t *testing.T) {
	id, _ := NewSquaremindIdentity("TestAgent", "")

	data := []byte("test message")
	signature := id.Sign(data)

	if len(signature) != 64 {
		t.Errorf("Expected 64-byte signature, got %d bytes", len(signature))
	}

	if !id.Verify(data, signature) {
		t.Error("Signature verification failed")
	}

	// Test with modified data
	modifiedData := []byte("modified message")
	if id.Verify(modifiedData, signature) {
		t.Error("Signature should not verify with modified data")
	}
}

func TestPublicKeyHex(t *testing.T) {
	id, _ := NewSquaremindIdentity("TestAgent", "")

	hex := id.PublicKeyHex()
	if len(hex) != 64 { // 32 bytes = 64 hex chars
		t.Errorf("Expected 64 hex chars, got %d", len(hex))
	}
}

func TestSIDShort(t *testing.T) {
	id, _ := NewSquaremindIdentity("TestAgent", "")

	short := id.SIDShort()
	if len(short) != 8 {
		t.Errorf("Expected 8 char short SID, got %d", len(short))
	}

	if short != id.SID[:8] {
		t.Errorf("Short SID doesn't match first 8 chars of SID")
	}
}
