package loandm_test

import (
	"testing"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

func TestUserCanBorrowWhenAllConditionsMet(t *testing.T) {
	service := loandm.NewLoanEligibilityService()

	// Create user with default state (0 loans, 0 fees)
	u := userdm.NewUser("John Doe", "john@example.com")

	// Create book
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(&bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatal(err)
	}

	if !service.CanBorrow(u, b) {
		t.Error("ユーザーは本を借りられるべきです")
	}
}

func TestUserCannotBorrowWhenMaxLoansReached(t *testing.T) {
	service := loandm.NewLoanEligibilityService()

	// Create user with max loans (5) using ReconstructUser
	userID := userdm.GenerateUserID()
	u := userdm.ReconstructUser(
		userID,
		"John Doe",
		"john@example.com",
		userdm.UserStatusActive,
		5, // currentLoanCount = MaxLoans
		0, // overdueFees
		time.Now(),
	)

	// Create book
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(&bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatal(err)
	}

	if service.CanBorrow(u, b) {
		t.Error("ユーザーは本を借りられないべきです（貸出上限に達している）")
	}
}
