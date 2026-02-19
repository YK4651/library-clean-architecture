package extendloan

import (
	"context"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

type IUserRepository interface {
	FindByID(ctx context.Context, id *userdm.UserID) (*userdm.User, error)
}

type ILoanRepository interface {
	FindByIDAndUserID(ctx context.Context, id *loandm.LoanID, userID *userdm.UserID) (*loandm.Loan, error)
	ListActiveLoansByUser(ctx context.Context, userID *userdm.UserID) ([]*loandm.Loan, error)
	Save(ctx context.Context, loan *loandm.Loan) error
}

type ExtendLoanUseCase struct {
	userRepo IUserRepository
	loanRepo ILoanRepository
}

func NewExtendLoanUseCase(userRepo IUserRepository, loanRepo ILoanRepository) *ExtendLoanUseCase {
	return &ExtendLoanUseCase{userRepo: userRepo, loanRepo: loanRepo}
}

// Execute はローン延長を実行する。
// 1) ユーザーをロード 2) 停止中でないか最初にチェック 3) このユーザーの全ローンで延滞なしをチェック
// 4) findByIdAndUserId でローン取得（自分のものでなければ nil）5) 書籍一致確認 6) 延長して保存
func (uc *ExtendLoanUseCase) Execute(ctx context.Context, req *ExtendLoanRequest) (*ExtendLoanResponse, error) {
	userID, err := userdm.NewUserID(req.UserID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &UserNotFoundError{UserID: req.UserID}
	}

	// ルール1: ユーザー状態を最初にチェック（軽量）
	if user.Status() == userdm.UserStatusSuspended {
		return nil, &UserSuspendedError{UserID: req.UserID}
	}

	// ルール2: このユーザーの全ローンで延滞がないこと
	activeLoans, err := uc.loanRepo.ListActiveLoansByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if loandm.OverdueLoanChecker(activeLoans, nil) {
		return nil, &UserHasOverdueLoansError{UserID: req.UserID}
	}

	loanID, err := loandm.LoanIDFromString(req.LoanID)
	if err != nil {
		return nil, err
	}
	bookID, err := bookdm.BookIDFromString(req.BookID)
	if err != nil {
		return nil, err
	}

	// セキュリティ組み込み: 自分のローンでなければ nil
	loan, err := uc.loanRepo.FindByIDAndUserID(ctx, &loanID, userID)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, &LoanNotFoundOrForbiddenError{LoanID: req.LoanID}
	}

	if !loan.BookID().Equals(bookID) {
		return nil, &BookLoanMismatchError{BookID: req.BookID, LoanID: req.LoanID}
	}

	if loan.IsReturned() {
		return nil, &LoanAlreadyReturnedError{LoanID: req.LoanID}
	}

	extended := loan.Extend()
	if err := uc.loanRepo.Save(ctx, extended); err != nil {
		return nil, err
	}

	newDue := extended.DueDate()
	return &ExtendLoanResponse{
		NewDueDate:   newDue.Format("2006-01-02"),
		ExtendedDays: loandm.ExtendDays,
	}, nil
}

// --- エラー型（ハンドラで errors.As 用）---

type UserNotFoundError struct{ UserID string }
func (e *UserNotFoundError) Error() string { return "user not found" }

type UserSuspendedError struct{ UserID string }
func (e *UserSuspendedError) Error() string { return "user is suspended" }

type UserHasOverdueLoansError struct{ UserID string }
func (e *UserHasOverdueLoansError) Error() string { return "user has overdue loans" }

type LoanNotFoundOrForbiddenError struct{ LoanID string }
func (e *LoanNotFoundOrForbiddenError) Error() string { return "loan not found or access denied" }

type BookLoanMismatchError struct{ BookID, LoanID string }
func (e *BookLoanMismatchError) Error() string { return "book and loan do not match" }

type LoanAlreadyReturnedError struct{ LoanID string }
func (e *LoanAlreadyReturnedError) Error() string { return "loan has already been returned" }
