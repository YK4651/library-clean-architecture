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

func (s *LoanEligibilityService) CanBorrow(u *userdm.User, currentLoanCount int, b *bookdm.Book, currentActiveLoans int) bool {
	// Rule 1: User must be eligible to borrow
	if !u.CanBorrow(currentLoanCount) {
		return false
	}

	// Rule 2: Book must have available copies
	if !b.IsAvailable(currentActiveLoans) {
		return false
	}

	return true
}

func (s *LoanEligibilityService) IneligibilityReason(u *userdm.User, currentLoanCount int, b *bookdm.Book, currentActiveLoans int) *string {
	if !u.CanBorrow(currentLoanCount) {
		reason := fmt.Sprintf("ユーザーは貸出条件を満たしていません（現在の貸出数: %d）", currentLoanCount)
		return &reason
	}

	if !b.IsAvailable(currentActiveLoans) {
		reason := fmt.Sprintf("書籍「%s」の利用可能なコピーがありません", b.Title())
		return &reason
	}

	return nil
}
