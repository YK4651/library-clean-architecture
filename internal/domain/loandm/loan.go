package loandm

import (
	"errors"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

// Loanエンティティ - 本の貸出を表す
type Loan struct {
	id         *LoanID
	userID     *userdm.UserID
	bookID     *bookdm.BookID
	borrowedAt time.Time
	dueDate    time.Time
	returnedAt *time.Time
	lateFee    int // 延滞料金（円）。返却時に計算して保存。
}

const LoanPeriodDays = 14

func NewLoan(
	id *LoanID,
	userID *userdm.UserID,
	bookID *bookdm.BookID,
	borrowedAt time.Time,
	returnedAt *time.Time,
) *Loan {
	if borrowedAt.IsZero() {
		borrowedAt = time.Now()
	}

	dueDate := calculateDueDate(borrowedAt)

	return &Loan{
		id:         id,
		userID:     userID,
		bookID:     bookID,
		borrowedAt: borrowedAt,
		dueDate:    dueDate,
		returnedAt: returnedAt,
		lateFee:    0,
	}
}

func calculateDueDate(borrowedAt time.Time) time.Time {
	return borrowedAt.AddDate(0, 0, LoanPeriodDays)
}

func (l *Loan) IsOverdue(currentDate *time.Time) bool {
	var now time.Time
	if currentDate != nil {
		now = *currentDate
	} else {
		now = time.Now()
	}
	return now.After(l.dueDate) && l.returnedAt == nil
}

func (l *Loan) IsReturned() bool {
	return l.returnedAt != nil
}

// CalculateLateFee は返却日時に対する延滞料金を計算する（1日あたり10円、猶予なし・上限なし）
func (l *Loan) CalculateLateFee(returnDate time.Time) int {
	if !returnDate.After(l.dueDate) {
		return 0
	}
	daysLate := int(returnDate.Sub(l.dueDate).Hours() / 24)
	if daysLate <= 0 {
		return 0
	}
	return daysLate * 10
}

// DaysLate は返却日時が期限を過ぎている場合の延滞日数を返す。期限内なら0。
func (l *Loan) DaysLate(returnDate time.Time) int {
	if !returnDate.After(l.dueDate) {
		return 0
	}
	return int(returnDate.Sub(l.dueDate).Hours() / 24)
}

// ReturnBook は返却済みとしてマークし、延滞料金を設定した新しいLoanを返す（不変）
func (l *Loan) ReturnBook(returnDate time.Time) (*Loan, error) {
	if l.IsReturned() {
		return nil, errors.New("loan has already been returned")
	}
	lateFee := l.CalculateLateFee(returnDate)
	return &Loan{
		id:         l.id,
		userID:     l.userID,
		bookID:     l.bookID,
		borrowedAt: l.borrowedAt,
		dueDate:    l.dueDate,
		returnedAt: &returnDate,
		lateFee:    lateFee,
	}, nil
}

func (l *Loan) MarkAsReturned(returnedAt *time.Time) (*Loan, error) {
	if l.IsReturned() {
		return nil, errors.New("loan has already been returned")
	}

	var returnTime *time.Time
	if returnedAt != nil {
		returnTime = returnedAt
	} else {
		now := time.Now()
		returnTime = &now
	}

	return &Loan{
		id:         l.id,
		userID:     l.userID,
		bookID:     l.bookID,
		borrowedAt: l.borrowedAt,
		dueDate:    l.dueDate,
		returnedAt: returnTime,
		lateFee:    l.lateFee,
	}, nil
}

func (l *Loan) DaysUntilDue(currentDate *time.Time) int {
	var now time.Time
	if currentDate != nil {
		now = *currentDate
	} else {
		now = time.Now()
	}

	duration := l.dueDate.Sub(now)
	days := int(duration.Hours() / 24)

	return days
}

// ゲッター（Goの慣習: "Get"プレフィックスなし）
func (l *Loan) Id() *LoanID {
	return l.id
}

func (l *Loan) UserID() *userdm.UserID {
	return l.userID
}

func (l *Loan) BookID() *bookdm.BookID {
	return l.bookID
}

func (l *Loan) BorrowedAt() time.Time {
	return l.borrowedAt
}

func (l *Loan) DueDate() time.Time {
	return l.dueDate
}

func (l *Loan) ReturnedAt() *time.Time {
	return l.returnedAt
}

func (l *Loan) LateFee() int {
	return l.lateFee
}

func (l *Loan) LoanPeriodDays() int {
	return LoanPeriodDays
}

const ExtendDays = 14

// Extend は返却期限に14日を追加する。延長回数制限はなし。バリデーションはユースケースで行う。
func (l *Loan) Extend() *Loan {
	return &Loan{
		id:         l.id,
		userID:     l.userID,
		bookID:     l.bookID,
		borrowedAt: l.borrowedAt,
		dueDate:    l.dueDate.AddDate(0, 0, ExtendDays),
		returnedAt: l.returnedAt,
		lateFee:    l.lateFee,
	}
}

// ReconstructLoan は永続化層から復元するためのコンストラクタ（due_date, late_fee を指定可能）
func ReconstructLoan(
	id *LoanID,
	userID *userdm.UserID,
	bookID *bookdm.BookID,
	borrowedAt time.Time,
	dueDate time.Time,
	returnedAt *time.Time,
	lateFee int,
) *Loan {
	return &Loan{
		id:         id,
		userID:     userID,
		bookID:     bookID,
		borrowedAt: borrowedAt,
		dueDate:    dueDate,
		returnedAt: returnedAt,
		lateFee:    lateFee,
	}
}
