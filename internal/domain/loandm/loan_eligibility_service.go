package loandm

import (
	"fmt"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

// LoanEligibilityService - Domain Service for loan eligibility checks
type LoanEligibilityService struct{}

func NewLoanEligibilityService() *LoanEligibilityService {
	return &LoanEligibilityService{}
}

func (s *LoanEligibilityService) CanBorrow(u *userdm.User, b *bookdm.Book) bool {
	// Rule 1: User must not have reached max loan limit
	if !u.CanBorrowMore() {
		return false
	}

	// Rule 2: User must not have overdue books
	if u.HasOverdueBooks() {
		return false
	}

	// Rule 3: Book must have available copies
	if !b.Available() {
		return false
	}

	return true
}

func (s *LoanEligibilityService) IneligibilityReason(u *userdm.User, b *bookdm.Book) *string {
	if !u.CanBorrowMore() {
		reason := fmt.Sprintf("ユーザーは最大貸出制限に達しています（%d冊）", u.GetMaxLoans())
		return &reason
	}

	if u.HasOverdueBooks() {
		reason := "ユーザーは延滞中の本があります"
		return &reason
	}

	if !b.Available() {
		reason := fmt.Sprintf("書籍「%s」の利用可能なコピーがありません", b.GetTitle())
		return &reason
	}

	return nil
}
