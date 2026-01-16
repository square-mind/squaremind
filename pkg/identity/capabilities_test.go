package identity

import (
	"testing"
)

func TestNewCapabilitySet(t *testing.T) {
	cs := NewCapabilitySet()
	if cs == nil {
		t.Fatal("NewCapabilitySet returned nil")
	}
	if cs.Capabilities == nil {
		t.Fatal("Capabilities map should be initialized")
	}
}

func TestCapabilitySet_Add(t *testing.T) {
	cs := NewCapabilitySet()
	cap := &Capability{
		Type:        CapCodeWrite,
		Proficiency: 0.8,
	}

	cs.Add(cap)

	if !cs.Has(CapCodeWrite) {
		t.Error("CapabilitySet should have CapCodeWrite after Add")
	}
}

func TestCapabilitySet_Has(t *testing.T) {
	cs := NewCapabilitySet()
	cs.Add(&Capability{Type: CapCodeWrite, Proficiency: 0.5})

	if !cs.Has(CapCodeWrite) {
		t.Error("Has should return true for added capability")
	}

	if cs.Has(CapSecurity) {
		t.Error("Has should return false for non-existent capability")
	}
}

func TestCapabilitySet_Get(t *testing.T) {
	cs := NewCapabilitySet()
	cap := &Capability{Type: CapCodeWrite, Proficiency: 0.75}
	cs.Add(cap)

	retrieved := cs.Get(CapCodeWrite)
	if retrieved == nil {
		t.Fatal("Get should return the capability")
	}

	if retrieved.Proficiency != 0.75 {
		t.Errorf("Expected proficiency 0.75, got %f", retrieved.Proficiency)
	}

	if cs.Get(CapSecurity) != nil {
		t.Error("Get should return nil for non-existent capability")
	}
}

func TestCapabilitySet_List(t *testing.T) {
	cs := NewCapabilitySet()
	cs.Add(&Capability{Type: CapCodeWrite, Proficiency: 0.5})
	cs.Add(&Capability{Type: CapCodeReview, Proficiency: 0.6})

	list := cs.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 capabilities, got %d", len(list))
	}
}

func TestCapabilitySet_MatchScore(t *testing.T) {
	cs := NewCapabilitySet()
	cs.Add(&Capability{Type: CapCodeWrite, Proficiency: 0.8})
	cs.Add(&Capability{Type: CapCodeReview, Proficiency: 0.6})

	// Test empty requirements
	score := cs.MatchScore([]CapabilityType{})
	if score != 1.0 {
		t.Errorf("Expected score 1.0 for empty requirements, got %f", score)
	}

	// Test full match
	score = cs.MatchScore([]CapabilityType{CapCodeWrite})
	if score != 0.8 {
		t.Errorf("Expected score 0.8, got %f", score)
	}

	// Test partial match
	score = cs.MatchScore([]CapabilityType{CapCodeWrite, CapCodeReview})
	expected := (0.8 + 0.6) / 2 // Both matched
	if score != expected {
		t.Errorf("Expected score %f, got %f", expected, score)
	}

	// Test no match
	score = cs.MatchScore([]CapabilityType{CapSecurity})
	if score != 0.0 {
		t.Errorf("Expected score 0.0 for no match, got %f", score)
	}

	// Test partial coverage
	score = cs.MatchScore([]CapabilityType{CapCodeWrite, CapSecurity})
	// Coverage: 1/2 = 0.5, proficiency: 0.8
	expected = 0.5 * 0.8
	if score != expected {
		t.Errorf("Expected score %f, got %f", expected, score)
	}
}

func TestCapabilitySet_JSON(t *testing.T) {
	cs := NewCapabilitySet()
	cs.Add(&Capability{Type: CapCodeWrite, Proficiency: 0.8})

	json := cs.ToJSON()
	if json == "" || json == "{}" {
		t.Error("ToJSON should return non-empty JSON")
	}

	// Test FromJSON
	cs2 := NewCapabilitySet()
	err := cs2.FromJSON([]byte(json))
	if err != nil {
		t.Errorf("FromJSON failed: %v", err)
	}
}
