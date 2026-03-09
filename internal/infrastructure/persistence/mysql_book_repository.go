package persistence

import (
	"context"
	"database/sql"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
)

type MySQLBookRepository struct {
	db *sql.DB
}

func NewMySQLBookRepository(db *sql.DB) bookdm.IBookRepository {
	return &MySQLBookRepository{db: db}
}

func (r *MySQLBookRepository) FindByID(ctx context.Context, id *bookdm.BookID) (*bookdm.Book, error) {
	query := "SELECT id, title, author, isbn FROM books WHERE id = ?"

	var idStr, title, author, isbnStr string
	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&idStr, &title, &author, &isbnStr,
	)
	_ = idStr // idは引数から既知
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	isbn, err := bookdm.NewISBN(isbnStr)
	if err != nil {
		return nil, err
	}
	// total_copiesはスキーマにないため1を使用
	return bookdm.NewBook(*id, title, author, isbn, 1)
}
