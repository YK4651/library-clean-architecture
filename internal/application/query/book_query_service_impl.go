package query

import (
	"context"
	"database/sql"
)

// BookQueryServiceImpl - BookQueryServiceインターフェースの実装
//
// 注意：DB接続はフィールドとして保持しません。
// 代わりに、各メソッドでcontextから取得します。
type BookQueryServiceImpl struct {
	// DBフィールドなし - contextから取得！
}

// NewBookQueryService - 新しいBookQueryServiceImplを作成
//
// DBパラメータは不要 - contextから取得します
func NewBookQueryService() *BookQueryServiceImpl {
	return &BookQueryServiceImpl{}
}

// GetBookByID - IDで書籍を取得し、現在の貸出情報も含める
func (s *BookQueryServiceImpl) GetBookByID(ctx context.Context, bookID string) (*BookReadModel, error) {
	// ミドルウェアが提供するDB接続をcontextから取得
	db := ctx.Value("db").(*sql.DB)

	// 単一JOINクエリ - 効率的！
	query := `
		SELECT
			b.id,
			b.isbn,
			b.title,
			b.author,
			l.id as loan_id,
			l.user_id,
			l.borrowed_at,
			l.due_date
		FROM books b
		LEFT JOIN loans l ON b.id = l.book_id
			AND l.returned_at IS NULL  -- アクティブな貸出のみ
		WHERE b.id = ?
	`

	var (
		id, isbn, title, author string
		loanID, userID          sql.NullString // NULLを許可
		borrowedAt, dueDate     sql.NullString
	)

	// QueryRowContext を使用してcontextを渡す
	err := db.QueryRowContext(ctx, query, bookID).Scan(
		&id, &isbn, &title, &author,
		&loanID, &userID, &borrowedAt, &dueDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil // 書籍が見つかりません
	}
	if err != nil {
		return nil, err // データベースエラー
	}

	// クエリ結果から直接Read Modelを構築
	var currentLoan *LoanReadModel
	if loanID.Valid { // 貸出がある場合
		currentLoan = &LoanReadModel{
			LoanID:       loanID.String,
			UserID:       userID.String,
			BorrowedDate: borrowedAt.String,
			DueDate:      dueDate.String,
		}
	}

	return &BookReadModel{
		ID:          id,
		ISBN:        isbn,
		Title:       title,
		Author:      author,
		IsAvailable: currentLoan == nil, // 単一の情報源！
		CurrentLoan: currentLoan,
	}, nil
}

// ListBooks - ページネーション付きで書籍を取得（単一JOIN、N+1回避）
func (s *BookQueryServiceImpl) ListBooks(ctx context.Context, limit, offset int) (*BookListReadModel, error) {
	db := ctx.Value("db").(*sql.DB)

	// 1. 総件数を取得
	var total int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM books").Scan(&total)
	if err != nil {
		return nil, err
	}

	// 2. 単一JOINで書籍＋アクティブ貸出を取得（N+1なし）
	listQuery := `
		SELECT
			b.id,
			b.isbn,
			b.title,
			b.author,
			l.id as loan_id,
			l.user_id,
			l.borrowed_at,
			l.due_date
		FROM books b
		LEFT JOIN loans l ON b.id = l.book_id
			AND l.returned_at IS NULL
		ORDER BY b.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.QueryContext(ctx, listQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*BookReadModel
	for rows.Next() {
		var (
			id, isbn, title, author   string
			loanID, userID            sql.NullString
			borrowedAt, dueDate       sql.NullString
		)
		if err := rows.Scan(
			&id, &isbn, &title, &author,
			&loanID, &userID, &borrowedAt, &dueDate,
		); err != nil {
			return nil, err
		}

		var currentLoan *LoanReadModel
		if loanID.Valid {
			currentLoan = &LoanReadModel{
				LoanID:       loanID.String,
				UserID:       userID.String,
				BorrowedDate: borrowedAt.String,
				DueDate:      dueDate.String,
			}
		}

		books = append(books, &BookReadModel{
			ID:          id,
			ISBN:        isbn,
			Title:       title,
			Author:      author,
			IsAvailable: currentLoan == nil,
			CurrentLoan: currentLoan,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if books == nil {
		books = []*BookReadModel{}
	}

	return &BookListReadModel{
		Books:  books,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
