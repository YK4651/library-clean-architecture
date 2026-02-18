package borrowbook_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/application/borrowbook"
	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

type mockUserRepo struct {
	user *userdm.User
	err  error
}

func (m *mockUserRepo) FindByID(ctx context.Context, id *userdm.UserID) (*userdm.User, error) {
	return m.user, m.err
}

type mockBookRepo struct {
	book *bookdm.Book
	err  error
}

func (m *mockBookRepo) FindByID(ctx context.Context, id *bookdm.BookID) (*bookdm.Book, error) {
	return m.book, m.err
}

type mockLoanRepo struct {
	userCount int
	bookCount int
	saveErr   error
	SavedLoan *loandm.Loan
}

func (m *mockLoanRepo) CountActiveLoansForUser(ctx context.Context, userID *userdm.UserID) (int, error) {
	return m.userCount, nil
}

func (m *mockLoanRepo) CountActiveLoansForBook(ctx context.Context, bookID *bookdm.BookID) (int, error) {
	return m.bookCount, nil
}

func (m *mockLoanRepo) Save(ctx context.Context, loan *loandm.Loan) error {
	m.SavedLoan = loan
	return m.saveErr
}

func TestBorrowBookUseCase_Success(t *testing.T) {
	ctx := context.Background()
	userID, _ := userdm.NewUserID("12345678")
	u := userdm.ReconstructUser(userID, "John", "john@example.com", userdm.UserStatusActive, 0, time.Now())

	bookID := bookdm.NewBookID()
	isbn, _ := bookdm.NewISBN("9780134494166")
	b, _ := bookdm.NewBook(&bookID, "Clean Architecture", "Robert Martin", isbn, 3)

	userRepo := &mockUserRepo{user: u}
	bookRepo := &mockBookRepo{book: b}
	loanRepo := &mockLoanRepo{userCount: 0, bookCount: 0}

	useCase := borrowbook.NewBorrowBookUseCase(userRepo, bookRepo, loanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookID.Value())
	if err != nil {
		t.Fatal(err)
	}

	resp, err := useCase.Execute(ctx, req)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if resp.LoanID == "" {
		t.Error("Expected loan ID, got empty string")
	}
}

func TestBorrowBookUseCase_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{user: nil, err: errors.New("user not found")}
	bookRepo := &mockBookRepo{}
	loanRepo := &mockLoanRepo{}

	useCase := borrowbook.NewBorrowBookUseCase(userRepo, bookRepo, loanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookdm.NewBookID().Value())
	if err != nil {
		t.Fatal(err)
	}

	_, err = useCase.Execute(ctx, req)
	if err == nil {
		t.Error("Expected error when user not found")
	}
}

func TestBorrowBookUseCase_BookNotFound(t *testing.T) {
	ctx := context.Background()
	userID, _ := userdm.NewUserID("12345678")
	u := userdm.ReconstructUser(userID, "John", "john@example.com", userdm.UserStatusActive, 0, time.Now())

	userRepo := &mockUserRepo{user: u}
	bookRepo := &mockBookRepo{book: nil, err: errors.New("book not found")}
	loanRepo := &mockLoanRepo{}

	useCase := borrowbook.NewBorrowBookUseCase(userRepo, bookRepo, loanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookdm.NewBookID().Value())
	if err != nil {
		t.Fatal(err)
	}

	_, err = useCase.Execute(ctx, req)
	if err == nil {
		t.Error("Expected error when book not found")
	}
}

func TestBorrowBookUseCase_LoanLimitExceeded(t *testing.T) {
	ctx := context.Background()
	userID, _ := userdm.NewUserID("12345678")
	u := userdm.ReconstructUser(userID, "John", "john@example.com", userdm.UserStatusActive, 0, time.Now())

	bookID := bookdm.NewBookID()
	isbn, _ := bookdm.NewISBN("9780134494166")
	b, _ := bookdm.NewBook(&bookID, "Clean Architecture", "Robert Martin", isbn, 3)

	userRepo := &mockUserRepo{user: u}
	bookRepo := &mockBookRepo{book: b}
	loanRepo := &mockLoanRepo{userCount: 5, bookCount: 0}

	useCase := borrowbook.NewBorrowBookUseCase(userRepo, bookRepo, loanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookID.Value())
	if err != nil {
		t.Fatal(err)
	}

	_, err = useCase.Execute(ctx, req)
	if err == nil {
		t.Error("Expected UserCannotBorrowError when loan limit reached")
	}
	if !errors.As(err, new(*borrowbook.UserCannotBorrowError)) {
		t.Errorf("Expected UserCannotBorrowError, got %T", err)
	}
}

func TestBorrowBookUseCase_BookNotAvailable(t *testing.T) {
	ctx := context.Background()
	userID, _ := userdm.NewUserID("12345678")
	u := userdm.ReconstructUser(userID, "John", "john@example.com", userdm.UserStatusActive, 0, time.Now())

	bookID := bookdm.NewBookID()
	isbn, _ := bookdm.NewISBN("9780134494166")
	b, _ := bookdm.NewBook(&bookID, "Clean Architecture", "Robert Martin", isbn, 1)

	userRepo := &mockUserRepo{user: u}
	bookRepo := &mockBookRepo{book: b}
	loanRepo := &mockLoanRepo{userCount: 0, bookCount: 1}

	useCase := borrowbook.NewBorrowBookUseCase(userRepo, bookRepo, loanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookID.Value())
	if err != nil {
		t.Fatal(err)
	}

	_, err = useCase.Execute(ctx, req)
	if err == nil {
		t.Error("Expected BookNotAvailableError")
	}
	if !errors.As(err, new(*borrowbook.BookNotAvailableError)) {
		t.Errorf("Expected BookNotAvailableError, got %T", err)
	}
}

func TestBorrowBookUseCase_CreatesLoanWithCorrectDueDate(t *testing.T) {
	ctx := context.Background()
	userID, _ := userdm.NewUserID("12345678")
	u := userdm.ReconstructUser(userID, "John", "john@example.com", userdm.UserStatusActive, 0, time.Now())

	bookID := bookdm.NewBookID()
	isbn, _ := bookdm.NewISBN("9780134494166")
	b, _ := bookdm.NewBook(&bookID, "Clean Architecture", "Robert Martin", isbn, 3)

	userRepo := &mockUserRepo{user: u}
	bookRepo := &mockBookRepo{book: b}
	loanRepo := &mockLoanRepo{userCount: 0, bookCount: 0}

	useCase := borrowbook.NewBorrowBookUseCase(userRepo, bookRepo, loanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookID.Value())
	if err != nil {
		t.Fatal(err)
	}

	_, err = useCase.Execute(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	if loanRepo.SavedLoan == nil {
		t.Fatal("Expected loan to be saved")
	}
	dueDate := loanRepo.SavedLoan.DueDate()
	now := time.Now()
	expectedDue := now.AddDate(0, 0, loandm.LoanPeriodDays)
	if dueDate.Format("2006-01-02") != expectedDue.Format("2006-01-02") {
		t.Errorf("Expected due date %v, got %v", expectedDue.Format("2006-01-02"), dueDate.Format("2006-01-02"))
	}
}
