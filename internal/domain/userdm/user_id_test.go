package userdm_test

import (
	"testing"

	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

func TestUserID_GeneratesValidID(t *testing.T) {
	id := userdm.GenerateUserID()

	if id.Value() == "" {
		t.Error("Expected non-empty ID value")
	}

	// UserID should be 8 digits
	if len(id.Value()) != 8 {
		t.Errorf("Expected UserID length 8, got: %d", len(id.Value()))
	}
}

func TestUserID_GeneratesUniqueIDs(t *testing.T) {
	id1 := userdm.GenerateUserID()
	id2 := userdm.GenerateUserID()

	if id1.Equals(id2) {
		t.Error("Expected different IDs to not be equal")
	}
}

func TestUserID_Equality(t *testing.T) {
	id1 := userdm.GenerateUserID()
	id2 := id1 // Same reference

	if !id1.Equals(id2) {
		t.Error("Expected same ID to be equal to itself")
	}
}

func TestUserID_ValueReturnsString(t *testing.T) {
	id := userdm.GenerateUserID()
	value := id.Value()

	if value == "" {
		t.Error("Value() should return non-empty string")
	}
}
