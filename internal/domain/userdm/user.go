package userdm

import (
	"errors"
	"time"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
)

// Userエンティティ - 図書館の利用者を表す
type User struct {
	id               *UserID
	name             string
	email            string
	status           UserStatus
	currentLoanCount int // 仮実装: 本来は Loan テーブルから取得 (SSOT)
	overdueFees      float64
	createdAt        time.Time
}

const MaxLoans = 5

// NewUser creates a new active user
func NewUser(name, email string) *User {
	return &User{
		id:               GenerateUserID(),
		name:             name,
		email:            email,
		status:           UserStatusActive,
		currentLoanCount: 0,
		overdueFees:      0,
		createdAt:        time.Now(),
	}
}

// ReconstructUser rebuilds user from persistence (仮実装: currentLoanCount を含む)
func ReconstructUser(
	id *UserID,
	name, email string,
	status UserStatus,
	currentLoanCount int,
	overdueFees float64,
	createdAt time.Time,
) *User {
	return &User{id, name, email, status, currentLoanCount, overdueFees, createdAt}
}

// User entity focuses on RULES, not STATE
// Loan count is derived from Loan table (Single Source of Truth)

// CanBorrow validates if user can borrow based on provided count
// The count is provided by the use case (derived from Loan table)
func (u *User) CanBorrow(currentLoanCount int) bool {
	if u.status == UserStatusSuspended {
		return false
	}
	if currentLoanCount >= MaxLoans {
		return false
	}
	if u.overdueFees > 0 {
		return false
	}
	return true
}

// State changes are tracked in Loan table (no BorrowBook/ReturnBook methods needed)

func (u *User) AddOverdueFee(amount float64) (*User, error) {
	if amount <= 0 {
		return nil, errors.New("overdue fee must be greater than 0")
	}
	return &User{
		id:               u.id,
		name:             u.name,
		email:            u.email,
		status:           u.status,
		currentLoanCount: u.currentLoanCount,
		overdueFees:      u.overdueFees + amount,
		createdAt:        u.createdAt,
	}, nil
}

func (u *User) PayOverdueFee(amount float64) (*User, error) {
	if amount < 0 {
		return nil, errors.New("payment amount cannot be negative")
	}
	if amount > u.overdueFees {
		return nil, errors.New("payment exceeds current overdue fees")
	}
	return &User{
		id:               u.id,
		name:             u.name,
		email:            u.email,
		status:           u.status,
		currentLoanCount: u.currentLoanCount,
		overdueFees:      u.overdueFees - amount,
		createdAt:        u.createdAt,
	}, nil
}

// ゲッター（Goの慣習: "Get"プレフィックスなし）
func (u *User) Id() *UserID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Status() UserStatus {
	return u.status
}

func (u *User) OverdueFees() float64 {
	return u.overdueFees
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// Note: CurrentLoanCount is NOT stored in User entity (SSOT pattern)
// Use LoanRepository.CountActiveLoansForUser() to get the current count

// HasOverdueBooks checks if user has overdue books (overdue fees > 0)
func (u *User) HasOverdueBooks() bool {
	return u.overdueFees > 0
}

// 仮実装: LoanEligibilityService 用。本来は UseCase が貸出数を渡し CanBorrow(count) を呼ぶ
func (u *User) CanBorrowMore() bool {
	return u.CanBorrow(u.currentLoanCount)
}

// GetMaxLoans returns max loans limit (仮実装: IneligibilityReason 用)
func (u *User) GetMaxLoans() int {
	return MaxLoans
}

// CurrentLoanCount returns current loan count (仮実装: テスト用。本来は Repository から取得)
func (u *User) CurrentLoanCount() int {
	return u.currentLoanCount
}
