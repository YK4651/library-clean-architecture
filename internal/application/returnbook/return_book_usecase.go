package returnbook

import (
	"context"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
)

type ILoanRepository interface {
	FindByID(ctx context.Context, id *loandm.LoanID) (*loandm.Loan, error)
	Save(ctx context.Context, loan *loandm.Loan) error
}

type ReturnBookUseCase struct {
	loanRepo ILoanRepository
}

func NewReturnBookUseCase(loanRepo ILoanRepository) *ReturnBookUseCase {
	return &ReturnBookUseCase{loanRepo: loanRepo}
}

// Execute は返却処理を実行する。Loanのみ更新（Bookは更新しない - SSOT）
func (uc *ReturnBookUseCase) Execute(ctx context.Context, req *ReturnBookRequest) (*ReturnBookResponse, error) {
	loanID, err := loandm.LoanIDFromString(req.LoanID)
	if err != nil {
		return nil, err
	}

	loan, err := uc.loanRepo.FindByID(ctx, &loanID)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, &LoanNotFoundError{LoanID: req.LoanID}
	}
	if loan.IsReturned() {
		return nil, &LoanAlreadyReturnedError{LoanID: req.LoanID}
	}

	returnDate := time.Now().UTC()
	returned, err := loan.ReturnBook(returnDate)
	if err != nil {
		return nil, err
	}

	if err := uc.loanRepo.Save(ctx, returned); err != nil {
		return nil, err
	}

	return NewReturnBookResponse(returned), nil
}

type LoanNotFoundError struct {
	LoanID string
}

func (e *LoanNotFoundError) Error() string {
	return "loan not found"
}

type LoanAlreadyReturnedError struct {
	LoanID string
}

func (e *LoanAlreadyReturnedError) Error() string {
	return "loan has already been returned"
}
