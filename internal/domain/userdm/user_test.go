package userdm_test

import (
	"testing"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

func TestCanCreateValidUser(t *testing.T) {
	// Create user with default state
	u := userdm.NewUser("John Doe", "john@example.com")

	// ID should be auto-generated (ULID 26 chars)
	if len(u.Id().Value()) != 26 {
		t.Errorf("Expected UserID length 26, got %d", len(u.Id().Value()))
	}
	if u.Name() != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", u.Name())
	}
	if u.Email() != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", u.Email())
	}
	if u.CurrentLoanCount() != 0 {
		t.Errorf("Expected initial loan count 0, got %d", u.CurrentLoanCount())
	}
}

func TestCannotBorrowMoreWhenMaxLoansReached(t *testing.T) {
	// Create user with max loans using ReconstructUser
	userID := userdm.GenerateUserID()
	u := userdm.ReconstructUser(
		userID,
		"John Doe",
		"john@example.com",
		userdm.UserStatusActive,
		5, // currentLoanCount = MaxLoans
		0,
		time.Now(),
	)

	if u.CanBorrowMore() {
		t.Error("ユーザーは貸出上限に達しているため、これ以上借りられないべきです")
	}
}

func TestCanBorrowMoreWhenUnderLimit(t *testing.T) {
	// Create user with 2 loans (under limit of 5)
	userID := userdm.GenerateUserID()
	u := userdm.ReconstructUser(
		userID,
		"John Doe",
		"john@example.com",
		userdm.UserStatusActive,
		2, // currentLoanCount < MaxLoans
		0,
		time.Now(),
	)

	if !u.CanBorrowMore() {
		t.Error("ユーザーは貸出上限未満のため、さらに借りられるべきです")
	}
}
