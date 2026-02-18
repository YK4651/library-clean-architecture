package borrowbook

import (
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
)

type BorrowBookResponse struct {
	LoanID     string
	UserID     string
	BookID     string
	BookTitle  string
	BorrowedAt time.Time
	DueDate    time.Time
}

func NewBorrowBookResponse(
	loan *loandm.Loan,
	bookTitle string,
) *BorrowBookResponse {
	return &BorrowBookResponse{
		LoanID:     loan.Id().Value(),
		UserID:     loan.UserID().Value(),
		BookID:     loan.BookID().Value(),
		BookTitle:  bookTitle,
		BorrowedAt: loan.BorrowedAt(),
		DueDate:    loan.DueDate(),
	}
}
