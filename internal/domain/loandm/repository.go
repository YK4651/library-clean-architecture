package loandm

import (
	"context"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

type ILoanRepository interface {
	// SSOT: Loanテーブルからアクティブな貸出をカウント
	CountActiveLoansForUser(ctx context.Context, userID *userdm.UserID) (int, error)
	CountActiveLoansForBook(ctx context.Context, bookID *bookdm.BookID) (int, error)

	FindByID(ctx context.Context, id *LoanID) (*Loan, error)
	// FindByIDAndUserID はローンを取得し、userIDが一致しない場合はnilを返す（セキュリティ組み込み）
	FindByIDAndUserID(ctx context.Context, id *LoanID, userID *userdm.UserID) (*Loan, error)
	// ListActiveLoansByUser はこのユーザーの未返却ローン一覧（延滞チェック用）
	ListActiveLoansByUser(ctx context.Context, userID *userdm.UserID) ([]*Loan, error)
	Save(ctx context.Context, loan *Loan) error
}
