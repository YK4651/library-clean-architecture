package loandm

import "time"

// OverdueLoanChecker はこのユーザーの全ローンに延滞があるかどうかを判定するドメインサービス
func OverdueLoanChecker(loans []*Loan, asOf *time.Time) bool {
	for _, l := range loans {
		if l.IsOverdue(asOf) {
			return true
		}
	}
	return false
}
