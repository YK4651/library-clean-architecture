package loandm_test

import (
	"testing"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

func TestDueDateIsCalculatedCorrectly(t *testing.T) {
	borrowedAt := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Create Value Objects
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()

	l := loandm.NewLoan(
		&loanID,
		userID,
		&bookID,
		borrowedAt,
		nil, // not returned yet
	)

	expectedDueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	if !l.DueDate().Equal(expectedDueDate) {
		t.Errorf("Expected due date %v, got %v", expectedDueDate, l.DueDate())
	}
}

func TestLoanIsOverdueAfterDueDate(t *testing.T) {
	borrowedAt := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC)

	// Create Value Objects
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()

	l := loandm.NewLoan(&loanID, userID, &bookID, borrowedAt, nil)

	if !l.IsOverdue(&currentDate) {
		t.Error("Expected loan to be overdue")
	}
}

func TestCanMarkLoanAsReturned(t *testing.T) {
	// Create Value Objects
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()

	l := loandm.NewLoan(&loanID, userID, &bookID, time.Now(), nil)

	if l.IsReturned() {
		t.Error("Loan should not be returned initially")
	}

	returned, err := l.MarkAsReturned(nil) // Pass nil to use current time
	if err != nil {
		t.Fatalf("Failed to mark as returned: %v", err)
	}

	if !returned.IsReturned() {
		t.Error("Loan should be marked as returned")
	}
	if returned.ReturnedAt() == nil {
		t.Error("ReturnedAt should be set")
	}
}

func TestCalculateLateFee_OnTime(t *testing.T) {
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	borrowedAt := dueDate.AddDate(0, 0, -14)
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()
	l := loandm.ReconstructLoan(&loanID, userID, &bookID, borrowedAt, dueDate, nil, 0)

	returnDate := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC) // 同日
	if fee := l.CalculateLateFee(returnDate); fee != 0 {
		t.Errorf("Expected late fee 0, got %d", fee)
	}
}

func TestCalculateLateFee_OneDayLate(t *testing.T) {
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	borrowedAt := dueDate.AddDate(0, 0, -14)
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()
	l := loandm.ReconstructLoan(&loanID, userID, &bookID, borrowedAt, dueDate, nil, 0)

	returnDate := time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC)
	if fee := l.CalculateLateFee(returnDate); fee != 10 {
		t.Errorf("Expected late fee 10, got %d", fee)
	}
}

func TestCalculateLateFee_FiveDaysLate(t *testing.T) {
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	borrowedAt := dueDate.AddDate(0, 0, -14)
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()
	l := loandm.ReconstructLoan(&loanID, userID, &bookID, borrowedAt, dueDate, nil, 0)

	returnDate := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC)
	if fee := l.CalculateLateFee(returnDate); fee != 50 {
		t.Errorf("Expected late fee 50, got %d", fee)
	}
}

func TestCalculateLateFee_TenDaysLate(t *testing.T) {
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	borrowedAt := dueDate.AddDate(0, 0, -14)
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()
	l := loandm.ReconstructLoan(&loanID, userID, &bookID, borrowedAt, dueDate, nil, 0)

	returnDate := time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC)
	if fee := l.CalculateLateFee(returnDate); fee != 100 {
		t.Errorf("Expected late fee 100, got %d", fee)
	}
}

func TestCalculateLateFee_HundredDaysLate(t *testing.T) {
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	borrowedAt := dueDate.AddDate(0, 0, -14)
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()
	l := loandm.ReconstructLoan(&loanID, userID, &bookID, borrowedAt, dueDate, nil, 0)

	returnDate := time.Date(2025, 4, 26, 0, 0, 0, 0, time.UTC) // 101 days later
	if fee := l.CalculateLateFee(returnDate); fee != 1010 {
		t.Errorf("Expected late fee 1010, got %d", fee)
	}
	returnDate = time.Date(2025, 4, 25, 0, 0, 0, 0, time.UTC) // 100 days
	if fee := l.CalculateLateFee(returnDate); fee != 1000 {
		t.Errorf("Expected late fee 1000, got %d", fee)
	}
}

func TestReturnBook_AlreadyReturned_ReturnsError(t *testing.T) {
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	returnedAt := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()
	l := loandm.ReconstructLoan(&loanID, userID, &bookID, dueDate.AddDate(0, 0, -14), dueDate, &returnedAt, 0)

	_, err := l.ReturnBook(time.Now())
	if err == nil {
		t.Fatal("Expected error when returning already returned loan")
	}
	if err.Error() != "loan has already been returned" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestExtend_Adds14Days(t *testing.T) {
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	borrowedAt := dueDate.AddDate(0, 0, -14)
	loanID := loandm.NewLoanID()
	userID := userdm.GenerateUserID()
	bookID := bookdm.NewBookID()
	l := loandm.ReconstructLoan(&loanID, userID, &bookID, borrowedAt, dueDate, nil, 0)

	extended := l.Extend()
	expected := time.Date(2025, 1, 29, 0, 0, 0, 0, time.UTC)
	if !extended.DueDate().Equal(expected) {
		t.Errorf("Expected due date %v, got %v", expected, extended.DueDate())
	}
	if extended.Id() != l.Id() || extended.UserID() != l.UserID() || extended.BookID() != l.BookID() {
		t.Error("Extend should preserve id, userID, bookID")
	}
}
