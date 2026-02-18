package returnbook

import (
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
)

type ReturnBookResponse struct {
	LoanID     string    `json:"loanId"`
	BookID     string    `json:"bookId"`
	ReturnedAt time.Time `json:"returnedAt"`
	LateFee    int       `json:"lateFee"`
	DaysLate   int       `json:"daysLate"`
	IsOverdue  bool      `json:"isOverdue"`
}

func NewReturnBookResponse(loan *loandm.Loan) *ReturnBookResponse {
	var returnedAt time.Time
	if loan.ReturnedAt() != nil {
		returnedAt = *loan.ReturnedAt()
	}
	return &ReturnBookResponse{
		LoanID:     loan.Id().Value(),
		BookID:     loan.BookID().Value(),
		ReturnedAt: returnedAt,
		LateFee:    loan.LateFee(),
		DaysLate:   loan.DaysLate(returnedAt),
		IsOverdue:  loan.LateFee() > 0,
	}
}
