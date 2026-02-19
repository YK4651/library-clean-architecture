package query

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3307)/library")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// テスト前にテーブルをクリーンアップ
	cleanupTables(t, db)

	return db
}

func cleanupTables(t *testing.T, db *sql.DB) {
	// 外部キー制約を一時的に無効化
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// テーブルをトランケート
	tables := []string{"loans", "books", "users"}
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table)
	}

	// 外部キー制約を再度有効化
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

func TestBookQueryService_GetBookByID_FindsAvailableBook(t *testing.T) {
	// Arrange（準備）
	db := setupTestDB(t)
	defer db.Close()

	// テストデータを挿入
	_, err := db.Exec(`
		INSERT INTO books (id, isbn, title, author, created_at)
		VALUES ('b-12345', '978-0-123456-78-9', 'Clean Architecture',
				'Robert C. Martin', '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	queryService := NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	// Act（実行）
	book, err := queryService.GetBookByID(ctx, "b-12345")

	// Assert（検証）
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if book == nil {
		t.Fatal("Expected book to be found, got nil")
	}
	if book.ID != "b-12345" {
		t.Errorf("Expected ID 'b-12345', got: %s", book.ID)
	}
	if book.Title != "Clean Architecture" {
		t.Errorf("Expected title 'Clean Architecture', got: %s", book.Title)
	}
	if !book.IsAvailable {
		t.Error("Expected book to be available")
	}
	if book.CurrentLoan != nil {
		t.Error("Expected no current loan for available book")
	}
}

func TestBookQueryService_GetBookByID_FindsBorrowedBook(t *testing.T) {
	// Arrange（準備）
	db := setupTestDB(t)
	defer db.Close()

	// テストデータと貸出を挿入
	_, err := db.Exec(`
		INSERT INTO books (id, isbn, title, author, created_at)
		VALUES ('b-12345', '978-0-123456-78-9', 'Clean Architecture',
				'Robert C. Martin', '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert book: %v", err)
	}

	// Insert test user (required for foreign key)
	_, err = db.Exec(`
		INSERT INTO users (id, name, email, status, created_at)
		VALUES ('u-001', 'Test User', 'test@example.com', 1, '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO loans (id, book_id, user_id, borrowed_at, due_date, returned_at)
		VALUES ('l-001', 'b-12345', 'u-001', '2024-01-15 10:00:00',
				'2024-02-15 10:00:00', NULL)
	`)
	if err != nil {
		t.Fatalf("Failed to insert loan: %v", err)
	}

	queryService := NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	// Act（実行）
	book, err := queryService.GetBookByID(ctx, "b-12345")

	// Assert（検証）
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if book == nil {
		t.Fatal("Expected book to be found, got nil")
	}
	if book.IsAvailable {
		t.Error("Expected book to be borrowed (not available)")
	}
	if book.CurrentLoan == nil {
		t.Fatal("Expected current loan to exist")
	}
	if book.CurrentLoan.LoanID != "l-001" {
		t.Errorf("Expected loan ID 'l-001', got: %s", book.CurrentLoan.LoanID)
	}
}

func TestBookQueryService_GetBookByID_ReturnsNilWhenNotFound(t *testing.T) {
	// Arrange（準備）
	db := setupTestDB(t)
	defer db.Close()

	queryService := NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	// Act（実行）
	book, err := queryService.GetBookByID(ctx, "b-99999")

	// Assert（検証）
	if err != nil {
		t.Errorf("Expected no error for not found case, got: %v", err)
	}
	if book != nil {
		t.Error("Expected nil for non-existent book")
	}
}

func TestBookQueryService_ListBooks_PaginationAndTotal(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	books := []struct{ id, isbn, title, author, created string }{
		{"b-1", "978-00000001", "Book 1", "Author", "2024-01-01 00:00:00"},
		{"b-2", "978-00000002", "Book 2", "Author", "2024-01-02 00:00:00"},
		{"b-3", "978-00000003", "Book 3", "Author", "2024-01-03 00:00:00"},
	}
	for _, b := range books {
		_, err := db.Exec(`INSERT INTO books (id, isbn, title, author, created_at) VALUES (?, ?, ?, ?, ?)`,
			b.id, b.isbn, b.title, b.author, b.created)
		if err != nil {
			t.Fatalf("Failed to insert book: %v", err)
		}
	}

	queryService := NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	list, err := queryService.ListBooks(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListBooks failed: %v", err)
	}
	if list.Total != 3 {
		t.Errorf("Expected total 3, got %d", list.Total)
	}
	if list.Limit != 10 || list.Offset != 0 {
		t.Errorf("Expected limit=10 offset=0, got limit=%d offset=%d", list.Limit, list.Offset)
	}
	if len(list.Books) != 3 {
		t.Errorf("Expected 3 books, got %d", len(list.Books))
	}
}

func TestBookQueryService_ListBooks_EmptyResultReturns200(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	queryService := NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	list, err := queryService.ListBooks(ctx, 20, 0)
	if err != nil {
		t.Fatalf("ListBooks failed: %v", err)
	}
	if list.Total != 0 {
		t.Errorf("Expected total 0, got %d", list.Total)
	}
	if list.Books == nil || len(list.Books) != 0 {
		t.Errorf("Expected empty books slice, got len=%d", len(list.Books))
	}
}

func TestBookQueryService_ListBooks_AvailableWhenNoActiveLoan(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec(`
		INSERT INTO books (id, isbn, title, author, created_at)
		VALUES ('b-avail', '978-0-13-468599-1', 'Clean Architecture', 'Robert C. Martin', '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert book: %v", err)
	}

	queryService := NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	list, err := queryService.ListBooks(ctx, 20, 0)
	if err != nil {
		t.Fatalf("ListBooks failed: %v", err)
	}
	if len(list.Books) != 1 {
		t.Fatalf("Expected 1 book, got %d", len(list.Books))
	}
	if !list.Books[0].IsAvailable {
		t.Error("Expected book with no active loan to be available")
	}
	if list.Books[0].CurrentLoan != nil {
		t.Error("Expected no current loan")
	}
}

func TestBookQueryService_ListBooks_BorrowedBookIsNotAvailable(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, _ = db.Exec(`
		INSERT INTO users (id, name, email, status, created_at)
		VALUES ('u-001', 'Test User', 'test@example.com', 1, '2024-01-01 00:00:00')
	`)
	_, err := db.Exec(`
		INSERT INTO books (id, isbn, title, author, created_at)
		VALUES ('b-borrowed', '978-0-321-12742-6', 'DDD', 'Eric Evans', '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert book: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO loans (id, book_id, user_id, borrowed_at, due_date, returned_at)
		VALUES ('l-11111', 'b-borrowed', 'u-001', '2024-01-10 00:00:00', '2024-01-24 00:00:00', NULL)
	`)
	if err != nil {
		t.Fatalf("Failed to insert loan: %v", err)
	}

	queryService := NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	list, err := queryService.ListBooks(ctx, 20, 0)
	if err != nil {
		t.Fatalf("ListBooks failed: %v", err)
	}
	if len(list.Books) != 1 {
		t.Fatalf("Expected 1 book, got %d", len(list.Books))
	}
	if list.Books[0].IsAvailable {
		t.Error("Expected borrowed book to be not available")
	}
	if list.Books[0].CurrentLoan == nil {
		t.Fatal("Expected current loan")
	}
	if list.Books[0].CurrentLoan.LoanID != "l-11111" {
		t.Errorf("Expected loan id l-11111, got %s", list.Books[0].CurrentLoan.LoanID)
	}
}

func TestBookQueryService_ListBooks_ReturnedLoansIgnored(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, _ = db.Exec(`
		INSERT INTO users (id, name, email, status, created_at)
		VALUES ('u-001', 'Test User', 'test@example.com', 1, '2024-01-01 00:00:00')
	`)
	_, err := db.Exec(`
		INSERT INTO books (id, isbn, title, author, created_at)
		VALUES ('b-returned', '978-0-321-12742-7', 'Returned Book', 'Author', '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert book: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO loans (id, book_id, user_id, borrowed_at, due_date, returned_at)
		VALUES ('l-ret', 'b-returned', 'u-001', '2024-01-01 00:00:00', '2024-01-15 00:00:00', '2024-01-10 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert loan: %v", err)
	}

	queryService := NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	list, err := queryService.ListBooks(ctx, 20, 0)
	if err != nil {
		t.Fatalf("ListBooks failed: %v", err)
	}
	if len(list.Books) != 1 {
		t.Fatalf("Expected 1 book, got %d", len(list.Books))
	}
	// 返却済みなのでアクティブな貸出はなく、利用可能
	if !list.Books[0].IsAvailable {
		t.Error("Expected book with returned loan to be available")
	}
	if list.Books[0].CurrentLoan != nil {
		t.Error("Expected no current loan (returned loans ignored)")
	}
}
