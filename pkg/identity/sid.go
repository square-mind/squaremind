package identity

import (
	"crypto/ed25519"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// SquaremindIdentity represents a unique, cryptographic identity for an agent
type SquaremindIdentity struct {
	SID        string             `json:"sid"`
	PublicKey  ed25519.PublicKey  `json:"public_key"`
	PrivateKey ed25519.PrivateKey `json:"-"` // Never serialized

	// Metadata
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ParentSID  string    `json:"parent_sid,omitempty"`
	Generation int       `json:"generation"`
}

// NewSquaremindIdentity creates a new squaremind identity with fresh keypair
func NewSquaremindIdentity(name string, parentSID string) (*SquaremindIdentity, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	generation := 0
	if parentSID != "" {
		generation = 1 // Would look up parent's generation in real impl
	}

	return &SquaremindIdentity{
		SID:        uuid.New().String(),
		PublicKey:  pub,
		PrivateKey: priv,
		Name:       name,
		CreatedAt:  time.Now(),
		ParentSID:  parentSID,
		Generation: generation,
	}, nil
}

// Sign signs data with the agent's private key
func (s *SquaremindIdentity) Sign(data []byte) []byte {
	return ed25519.Sign(s.PrivateKey, data)
}

// Verify verifies a signature against this identity's public key
func (s *SquaremindIdentity) Verify(data, signature []byte) bool {
	return ed25519.Verify(s.PublicKey, data, signature)
}

// PublicKeyHex returns the public key as a hex string
func (s *SquaremindIdentity) PublicKeyHex() string {
	return hex.EncodeToString(s.PublicKey)
}

// SIDShort returns a shortened version of the SID for display
func (s *SquaremindIdentity) SIDShort() string {
	if len(s.SID) > 8 {
		return s.SID[:8]
	}
	return s.SID
}
