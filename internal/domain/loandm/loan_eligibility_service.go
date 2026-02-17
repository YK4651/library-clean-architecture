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

func (s *LoanEligibilityService) CanBorrow(u *userdm.User, currentLoanCount int, b *bookdm.Book) bool {
	// Rule 1 & 2: User must be eligible (not suspended, under max loans, no overdue fees)
	if !u.CanBorrow(currentLoanCount) {
		return false
	}

	// Rule 3: Book must have available copies
	if !b.Available() {
		return false
	}

	return true
}

func (s *LoanEligibilityService) IneligibilityReason(u *userdm.User, currentLoanCount int, b *bookdm.Book) *string {
	if !u.CanBorrow(currentLoanCount) {
		if currentLoanCount >= userdm.MaxLoans {
			reason := fmt.Sprintf("ユーザーは最大貸出制限に達しています（%d冊）", userdm.MaxLoans)
			return &reason
		}
		if u.OverdueFees() > 0 {
			reason := "ユーザーは延滞中の本があります"
			return &reason
		}
		reason := "ユーザーは貸出できません（停止中または条件を満たしていません）"
		return &reason
	}

	if !b.Available() {
		reason := fmt.Sprintf("書籍「%s」の利用可能なコピーがありません", b.GetTitle())
		return &reason
	}

	return nil
}
