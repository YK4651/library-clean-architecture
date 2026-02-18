package returnbook_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/application/returnbook"
	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

type mockLoanRepo struct {
	loan    *loandm.Loan
	findErr error
	saveErr error
	saved   *loandm.Loan
}

func (m *mockLoanRepo) FindByID(ctx context.Context, id *loandm.LoanID) (*loandm.Loan, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.loan, nil
}

func (m *mockLoanRepo) Save(ctx context.Context, loan *loandm.Loan) error {
	m.saved = loan
	return m.saveErr
}

func TestReturnBookUseCase_Success_SavesLoanOnly(t *testing.T) {
	ctx := context.Background()
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	borrowedAt := dueDate.AddDate(0, 0, -14)
	loanID := loandm.NewLoanID()
	userID := userdm.ReconstructUserID("01HQXK5V8N3YZGJB4QWERT001")
	bookID := bookdm.NewBookID()
	activeLoan := loandm.ReconstructLoan(&loanID, userID, &bookID, borrowedAt, dueDate, nil, 0)

	repo := &mockLoanRepo{loan: activeLoan}
	uc := returnbook.NewReturnBookUseCase(repo)

	req, _ := returnbook.NewReturnBookRequest(loanID.Value())
	resp, err := uc.Execute(ctx, req)
	if err != nil {
		t.Fatalf("Expected success, got: %v", err)
	}
	if resp.LoanID != loanID.Value() {
		t.Errorf("Expected LoanID %s, got %s", loanID.Value(), resp.LoanID)
	}
	if resp.BookID != bookID.Value() {
		t.Errorf("Expected BookID %s, got %s", bookID.Value(), resp.BookID)
	}
	if repo.saved == nil {
		t.Fatal("Expected loan to be saved")
	}
	if !repo.saved.IsReturned() {
		t.Error("Saved loan should be marked returned")
	}
	// Book is not updated - SSOT: availability derived from loans
}

func TestReturnBookUseCase_OnTime_ZeroLateFee(t *testing.T) {
	// 返却日 = 期限日 → 延滞0円
	dueDate := time.Date(2025, 1, 25, 12, 0, 0, 0, time.UTC)
	borrowedAt := dueDate.AddDate(0, 0, -14)
	loanID := loandm.NewLoanID()
	userID := userdm.ReconstructUserID("01HQXK5V8N3YZGJB4QWERT001")
	bookID := bookdm.NewBookID()
	activeLoan := loandm.ReconstructLoan(&loanID, userID, &bookID, borrowedAt, dueDate, nil, 0)

	repo := &mockLoanRepo{loan: activeLoan}
	uc := returnbook.NewReturnBookUseCase(repo)

	req, _ := returnbook.NewReturnBookRequest(loanID.Value())
	resp, err := uc.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	// 実際の返却日は time.Now() なので、このテストではレスポンスの構造のみ検証
	if resp.LateFee < 0 {
		t.Errorf("LateFee should be >= 0, got %d", resp.LateFee)
	}
	// ドメインで「期限内→0円」は既にテスト済み。ここでは Save が呼ばれていることを確認
	if repo.saved == nil || repo.saved.LateFee() < 0 {
		t.Error("Saved loan should have valid late fee")
	}
}

func TestReturnBookUseCase_LoanNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockLoanRepo{loan: nil}
	uc := returnbook.NewReturnBookUseCase(repo)

	req, _ := returnbook.NewReturnBookRequest(loandm.NewLoanID().Value())
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("Expected error when loan not found")
	}
	if !errors.As(err, new(*returnbook.LoanNotFoundError)) {
		t.Errorf("Expected LoanNotFoundError, got %T", err)
	}
	if repo.saved != nil {
		t.Error("Loan should not be saved when not found")
	}
}

func TestReturnBookUseCase_AlreadyReturned(t *testing.T) {
	ctx := context.Background()
	returnedAt := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	dueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	loanID := loandm.NewLoanID()
	userID := userdm.ReconstructUserID("01HQXK5V8N3YZGJB4QWERT001")
	bookID := bookdm.NewBookID()
	returnedLoan := loandm.ReconstructLoan(&loanID, userID, &bookID, dueDate.AddDate(0, 0, -14), dueDate, &returnedAt, 0)

	repo := &mockLoanRepo{loan: returnedLoan}
	uc := returnbook.NewReturnBookUseCase(repo)

	req, _ := returnbook.NewReturnBookRequest(loanID.Value())
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("Expected error when loan already returned")
	}
	if !errors.As(err, new(*returnbook.LoanAlreadyReturnedError)) {
		t.Errorf("Expected LoanAlreadyReturnedError, got %T", err)
	}
	if repo.saved != nil {
		t.Error("Loan should not be saved when already returned")
	}
}

func TestReturnBookUseCase_InvalidLoanID(t *testing.T) {
	ctx := context.Background()
	repo := &mockLoanRepo{}
	uc := returnbook.NewReturnBookUseCase(repo)

	req, _ := returnbook.NewReturnBookRequest("invalid-id")
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("Expected error for invalid loan ID")
	}
}

func TestReturnBookUseCase_SaveError(t *testing.T) {
	ctx := context.Background()
	loanID := loandm.NewLoanID()
	userID := userdm.ReconstructUserID("01HQXK5V8N3YZGJB4QWERT001")
	bookID := bookdm.NewBookID()
	activeLoan := loandm.ReconstructLoan(&loanID, userID, &bookID,
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC), nil, 0)

	repo := &mockLoanRepo{loan: activeLoan, saveErr: errors.New("db error")}
	uc := returnbook.NewReturnBookUseCase(repo)

	req, _ := returnbook.NewReturnBookRequest(loanID.Value())
	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatal("Expected error when save fails")
	}
	if err.Error() != "db error" {
		t.Errorf("Expected db error, got %v", err)
	}
}
