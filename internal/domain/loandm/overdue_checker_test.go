package loandm_test

import (
	"testing"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

func TestOverdueLoanChecker_NoOverdue(t *testing.T) {
	dueDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	now := time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC)
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()
	loan := loandm.ReconstructLoan(&loanID, userID, &bookID, dueDate.AddDate(0, 0, -14), dueDate, nil, 0)

	loans := []*loandm.Loan{loan}
	if loandm.OverdueLoanChecker(loans, &now) {
		t.Error("Expected no overdue when current date before due")
	}
}

func TestOverdueLoanChecker_HasOverdue(t *testing.T) {
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	now := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC)
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()
	loan := loandm.ReconstructLoan(&loanID, userID, &bookID, dueDate.AddDate(0, 0, -14), dueDate, nil, 0)

	loans := []*loandm.Loan{loan}
	if !loandm.OverdueLoanChecker(loans, &now) {
		t.Error("Expected overdue when current date after due")
	}
}

func TestOverdueLoanChecker_ReturnedLoanNotOverdue(t *testing.T) {
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	returnedAt := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	now := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC)
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()
	loan := loandm.ReconstructLoan(&loanID, userID, &bookID, dueDate.AddDate(0, 0, -14), dueDate, &returnedAt, 0)

	loans := []*loandm.Loan{loan}
	if loandm.OverdueLoanChecker(loans, &now) {
		t.Error("Returned loan should not count as overdue")
	}
}

func TestOverdueLoanChecker_EmptyList(t *testing.T) {
	if loandm.OverdueLoanChecker(nil, nil) {
		t.Error("Empty list should have no overdue")
	}
	if loandm.OverdueLoanChecker([]*loandm.Loan{}, nil) {
		t.Error("Empty slice should have no overdue")
	}
}
