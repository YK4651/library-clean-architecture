package extendloan_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/application/extendloan"
	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

const testUserID = "12345678"

type mockUserRepo struct {
	user    *userdm.User
	findErr error
}

func (m *mockUserRepo) FindByID(ctx context.Context, id *userdm.UserID) (*userdm.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.user, nil
}

type mockLoanRepo struct {
	findByIDAndUserID *loandm.Loan
	listActive        []*loandm.Loan
	saveErr           error
	saved             *loandm.Loan
}

func (m *mockLoanRepo) FindByIDAndUserID(ctx context.Context, id *loandm.LoanID, userID *userdm.UserID) (*loandm.Loan, error) {
	return m.findByIDAndUserID, nil
}

func (m *mockLoanRepo) ListActiveLoansByUser(ctx context.Context, userID *userdm.UserID) ([]*loandm.Loan, error) {
	return m.listActive, nil
}

func (m *mockLoanRepo) Save(ctx context.Context, loan *loandm.Loan) error {
	m.saved = loan
	return m.saveErr
}

func TestExtendLoanUseCase_Success(t *testing.T) {
	ctx := context.Background()
	dueDate := time.Date(2030, 1, 15, 0, 0, 0, 0, time.UTC)
	borrowedAt := dueDate.AddDate(0, 0, -14)
	loanID := loandm.NewLoanID()
	userID := userdm.ReconstructUserID(testUserID)
	bookID := bookdm.NewBookID()
	activeLoan := loandm.ReconstructLoan(&loanID, userID, &bookID, borrowedAt, dueDate, nil, 0)

	user := userdm.ReconstructUser(userID, "name", "a@b.com", userdm.UserStatusActive, 0, time.Now())

	userRepo := &mockUserRepo{user: user}
	loanRepo := &mockLoanRepo{
		findByIDAndUserID: activeLoan,
		listActive:        []*loandm.Loan{activeLoan},
	}
	uc := extendloan.NewExtendLoanUseCase(userRepo, loanRepo)

	req, _ := extendloan.NewExtendLoanRequest(bookID.Value(), loanID.Value(), testUserID)
	resp, err := uc.Execute(ctx, req)
	if err != nil {
		t.Fatalf("expected success: %v", err)
	}
	if resp.NewDueDate != "2030-01-29" {
		t.Errorf("expected newDueDate 2030-01-29, got %s", resp.NewDueDate)
	}
	if resp.ExtendedDays != 14 {
		t.Errorf("expected extendedDays 14, got %d", resp.ExtendedDays)
	}
	if loanRepo.saved == nil || loanRepo.saved.DueDate().Format("2006-01-02") != "2030-01-29" {
		t.Error("expected extended loan to be saved with new due date")
	}
}

func TestExtendLoanUseCase_UserSuspended_Rejected(t *testing.T) {
	ctx := context.Background()
	userID := userdm.ReconstructUserID(testUserID)
	user := userdm.ReconstructUser(userID, "name", "a@b.com", userdm.UserStatusSuspended, 0, time.Now())
	loanID := loandm.NewLoanID()
	bookID := bookdm.NewBookID()

	userRepo := &mockUserRepo{user: user}
	loanRepo := &mockLoanRepo{listActive: nil}
	uc := extendloan.NewExtendLoanUseCase(userRepo, loanRepo)

	req, _ := extendloan.NewExtendLoanRequest(bookID.Value(), loanID.Value(), testUserID)
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("expected error when user suspended")
	}
	if !errors.As(err, new(*extendloan.UserSuspendedError)) {
		t.Errorf("expected UserSuspendedError, got %T", err)
	}
	if loanRepo.saved != nil {
		t.Error("loan should not be saved when user suspended")
	}
}

func TestExtendLoanUseCase_UserHasOverdue_Rejected(t *testing.T) {
	ctx := context.Background()
	dueDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	userID := userdm.ReconstructUserID(testUserID)
	user := userdm.ReconstructUser(userID, "name", "a@b.com", userdm.UserStatusActive, 0, time.Now())
	loanID := loandm.NewLoanID()
	bookID := bookdm.NewBookID()
	overdueLoan := loandm.ReconstructLoan(&loanID, userID, &bookID, dueDate.AddDate(0, 0, -14), dueDate, nil, 0)

	userRepo := &mockUserRepo{user: user}
	loanRepo := &mockLoanRepo{
		findByIDAndUserID: overdueLoan,
		listActive:        []*loandm.Loan{overdueLoan},
	}
	uc := extendloan.NewExtendLoanUseCase(userRepo, loanRepo)

	req, _ := extendloan.NewExtendLoanRequest(bookID.Value(), loanID.Value(), testUserID)
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("expected error when user has overdue loan")
	}
	if !errors.As(err, new(*extendloan.UserHasOverdueLoansError)) {
		t.Errorf("expected UserHasOverdueLoansError, got %T", err)
	}
}

func TestExtendLoanUseCase_LoanNotFoundOrNotOwned_Rejected(t *testing.T) {
	ctx := context.Background()
	userID := userdm.ReconstructUserID(testUserID)
	user := userdm.ReconstructUser(userID, "name", "a@b.com", userdm.UserStatusActive, 0, time.Now())
	loanID := loandm.NewLoanID()
	bookID := bookdm.NewBookID()

	userRepo := &mockUserRepo{user: user}
	loanRepo := &mockLoanRepo{
		findByIDAndUserID: nil,
		listActive:        []*loandm.Loan{},
	}
	uc := extendloan.NewExtendLoanUseCase(userRepo, loanRepo)

	req, _ := extendloan.NewExtendLoanRequest(bookID.Value(), loanID.Value(), testUserID)
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("expected error when loan not found or not user's")
	}
	if !errors.As(err, new(*extendloan.LoanNotFoundOrForbiddenError)) {
		t.Errorf("expected LoanNotFoundOrForbiddenError, got %T", err)
	}
}

func TestExtendLoanUseCase_BookLoanMismatch_Rejected(t *testing.T) {
	ctx := context.Background()
	dueDate := time.Date(2030, 1, 15, 0, 0, 0, 0, time.UTC)
	borrowedAt := dueDate.AddDate(0, 0, -14)
	loanID := loandm.NewLoanID()
	userID := userdm.ReconstructUserID(testUserID)
	bookID1 := bookdm.NewBookID()
	bookID2 := bookdm.NewBookID()
	activeLoan := loandm.ReconstructLoan(&loanID, userID, &bookID1, borrowedAt, dueDate, nil, 0)
	user := userdm.ReconstructUser(userID, "name", "a@b.com", userdm.UserStatusActive, 0, time.Now())

	userRepo := &mockUserRepo{user: user}
	loanRepo := &mockLoanRepo{
		findByIDAndUserID: activeLoan,
		listActive:        []*loandm.Loan{activeLoan},
	}
	uc := extendloan.NewExtendLoanUseCase(userRepo, loanRepo)

	req, _ := extendloan.NewExtendLoanRequest(bookID2.Value(), loanID.Value(), testUserID)
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("expected error when book and loan do not match")
	}
	if !errors.As(err, new(*extendloan.BookLoanMismatchError)) {
		t.Errorf("expected BookLoanMismatchError, got %T", err)
	}
}

func TestExtendLoanUseCase_LoanAlreadyReturned_Rejected(t *testing.T) {
	ctx := context.Background()
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	returnedAt := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	loanID := loandm.NewLoanID()
	userID := userdm.ReconstructUserID(testUserID)
	bookID := bookdm.NewBookID()
	returnedLoan := loandm.ReconstructLoan(&loanID, userID, &bookID, dueDate.AddDate(0, 0, -14), dueDate, &returnedAt, 0)
	user := userdm.ReconstructUser(userID, "name", "a@b.com", userdm.UserStatusActive, 0, time.Now())

	userRepo := &mockUserRepo{user: user}
	loanRepo := &mockLoanRepo{
		findByIDAndUserID: returnedLoan,
		listActive:        []*loandm.Loan{},
	}
	uc := extendloan.NewExtendLoanUseCase(userRepo, loanRepo)

	req, _ := extendloan.NewExtendLoanRequest(bookID.Value(), loanID.Value(), testUserID)
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("expected error when loan already returned")
	}
	if !errors.As(err, new(*extendloan.LoanAlreadyReturnedError)) {
		t.Errorf("expected LoanAlreadyReturnedError, got %T", err)
	}
}

func TestExtendLoanUseCase_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{user: nil}
	loanRepo := &mockLoanRepo{listActive: []*loandm.Loan{}}
	uc := extendloan.NewExtendLoanUseCase(userRepo, loanRepo)

	req, _ := extendloan.NewExtendLoanRequest(bookdm.NewBookID().Value(), loandm.NewLoanID().Value(), testUserID)
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("expected error when user not found")
	}
	if !errors.As(err, new(*extendloan.UserNotFoundError)) {
		t.Errorf("expected UserNotFoundError, got %T", err)
	}
}

func TestExtendLoanUseCase_InvalidUserID(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{}
	loanRepo := &mockLoanRepo{}
	uc := extendloan.NewExtendLoanUseCase(userRepo, loanRepo)

	req, _ := extendloan.NewExtendLoanRequest("01HXXX", "01HYYY", "invalid")
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("expected error for invalid user ID format")
	}
}

func TestExtendLoanUseCase_SaveError(t *testing.T) {
	ctx := context.Background()
	dueDate := time.Date(2030, 1, 15, 0, 0, 0, 0, time.UTC)
	loanID := loandm.NewLoanID()
	userID := userdm.ReconstructUserID(testUserID)
	bookID := bookdm.NewBookID()
	activeLoan := loandm.ReconstructLoan(&loanID, userID, &bookID, dueDate.AddDate(0, 0, -14), dueDate, nil, 0)
	user := userdm.ReconstructUser(userID, "name", "a@b.com", userdm.UserStatusActive, 0, time.Now())

	userRepo := &mockUserRepo{user: user}
	loanRepo := &mockLoanRepo{
		findByIDAndUserID: activeLoan,
		listActive:        []*loandm.Loan{activeLoan},
		saveErr:            errors.New("db error"),
	}
	uc := extendloan.NewExtendLoanUseCase(userRepo, loanRepo)

	req, _ := extendloan.NewExtendLoanRequest(bookID.Value(), loanID.Value(), testUserID)
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("expected error when save fails")
	}
	if err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}
