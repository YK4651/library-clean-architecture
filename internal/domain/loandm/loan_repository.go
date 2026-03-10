package loandm

import (
	"context"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

type ILoanRepository interface {
	// SSOT: Loanテーブルからアクティブな貸出をカウント
	CountActiveLoansForUser(ctx context.Context, userID *userdm.UserID) (uint32, error)
	CountActiveLoansForBook(ctx context.Context, bookID *bookdm.BookID) (uint32, error)

	Save(ctx context.Context, loan *Loan) error
}
