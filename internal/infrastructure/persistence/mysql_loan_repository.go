package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

type MySQLLoanRepository struct {
	db *sql.DB
}

func NewMySQLLoanRepository(db *sql.DB) loandm.ILoanRepository {
	return &MySQLLoanRepository{db: db}
}

func (r *MySQLLoanRepository) CountActiveLoansForUser(ctx context.Context, userID *userdm.UserID) (int, error) {
	query := "SELECT COUNT(*) FROM loans WHERE user_id = ? AND returned_at IS NULL"

	var count int
	err := r.db.QueryRowContext(ctx, query, userID.Value()).Scan(&count)
	return count, err
}

func (r *MySQLLoanRepository) CountActiveLoansForBook(ctx context.Context, bookID *bookdm.BookID) (int, error) {
	query := "SELECT COUNT(*) FROM loans WHERE book_id = ? AND returned_at IS NULL"

	var count int
	err := r.db.QueryRowContext(ctx, query, bookID.Value()).Scan(&count)
	return count, err
}

func (r *MySQLLoanRepository) FindByID(ctx context.Context, id *loandm.LoanID) (*loandm.Loan, error) {
	query := `SELECT id, user_id, book_id, borrowed_at, due_date, returned_at, late_fee
	          FROM loans WHERE id = ?`

	var idStr, userIDStr, bookIDStr string
	var borrowedAt, dueDate time.Time
	var returnedAt sql.NullTime
	var lateFee int

	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&idStr, &userIDStr, &bookIDStr, &borrowedAt, &dueDate, &returnedAt, &lateFee,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	loanID, _ := loandm.LoanIDFromString(idStr)
	uid := userdm.ReconstructUserID(userIDStr)
	bid, err := bookdm.BookIDFromString(bookIDStr)
	if err != nil {
		return nil, err
	}

	var ret *time.Time
	if returnedAt.Valid {
		ret = &returnedAt.Time
	}

	return loandm.ReconstructLoan(&loanID, uid, &bid, borrowedAt, dueDate, ret, lateFee), nil
}

func (r *MySQLLoanRepository) Save(ctx context.Context, loan *loandm.Loan) error {
	if loan.ReturnedAt() != nil {
		return r.updateReturned(ctx, loan)
	}
	return r.insert(ctx, loan)
}

func (r *MySQLLoanRepository) insert(ctx context.Context, loan *loandm.Loan) error {
	query := `INSERT INTO loans (id, user_id, book_id, borrowed_at, due_date, returned_at, late_fee)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	var returnedAt *string
	if loan.ReturnedAt() != nil {
		t := loan.ReturnedAt().Format("2006-01-02 15:04:05")
		returnedAt = &t
	}

	_, err := r.db.ExecContext(ctx, query,
		loan.Id().Value(),
		loan.UserID().Value(),
		loan.BookID().Value(),
		loan.BorrowedAt().Format("2006-01-02 15:04:05"),
		loan.DueDate().Format("2006-01-02 15:04:05"),
		returnedAt,
		loan.LateFee(),
	)
	return err
}

func (r *MySQLLoanRepository) updateReturned(ctx context.Context, loan *loandm.Loan) error {
	query := `UPDATE loans SET returned_at = ?, late_fee = ? WHERE id = ?`
	returnedAt := loan.ReturnedAt().Format("2006-01-02 15:04:05")
	_, err := r.db.ExecContext(ctx, query, returnedAt, loan.LateFee(), loan.Id().Value())
	return err
}
