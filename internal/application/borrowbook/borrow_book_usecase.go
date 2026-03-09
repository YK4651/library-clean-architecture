package borrowbook

//go:generate mockgen -destination=mocks/mock_repositories.go -package=mocks . IUserRepository,IBookRepository,ILoanRepository
import (
	"context"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
	"github.com/YK4651/library-clean-architecture/internal/domain/loandm"
	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

// リポジトリインターフェース（Step 2で実装を学習済み）
type IUserRepository interface {
	FindByID(ctx context.Context, id *userdm.UserID) (*userdm.User, error)
}

type IBookRepository interface {
	FindByID(ctx context.Context, id *bookdm.BookID) (*bookdm.Book, error)
}

type ILoanRepository interface {
	CountActiveLoansForUser(ctx context.Context, userID *userdm.UserID) (int, error)
	CountActiveLoansForBook(ctx context.Context, bookID *bookdm.BookID) (int, error)
	Save(ctx context.Context, loan *loandm.Loan) error
}

type BorrowBookUseCase struct {
	userRepo IUserRepository
	bookRepo IBookRepository
	loanRepo ILoanRepository
}

func NewBorrowBookUseCase(
	userRepo IUserRepository,
	bookRepo IBookRepository,
	loanRepo ILoanRepository,
) *BorrowBookUseCase {
	return &BorrowBookUseCase{
		userRepo: userRepo,
		bookRepo: bookRepo,
		loanRepo: loanRepo,
	}
}

// Execute - 書籍を借りる - 複数エンティティを調整
// ミドルウェアがトランザクションを自動管理
func (uc *BorrowBookUseCase) Execute(ctx context.Context, req *BorrowBookRequest) (*BorrowBookResponse, error) {
	// contextにはミドルウェアが設定したDB/tx接続が含まれている

	// 1. ユーザーIDを作成
	userID, err := userdm.NewUserID(req.UserID)
	if err != nil {
		return nil, err
	}

	// 2. ユーザーエンティティを取得
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. 書籍IDを作成
	bookID, err := bookdm.BookIDFromString(req.BookID)
	if err != nil {
		return nil, err
	}

	// 4. 書籍エンティティを取得
	book, err := uc.bookRepo.FindByID(ctx, &bookID)
	if err != nil {
		return nil, err
	}

	// 5. アクティブ貸出数をカウント（SSOT）
	userCurrentLoans, err := uc.loanRepo.CountActiveLoansForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	bookActiveLoans, err := uc.loanRepo.CountActiveLoansForBook(ctx, &bookID)
	if err != nil {
		return nil, err
	}

	// 6. ビジネスルールを検証
	if !user.CanBorrow(userCurrentLoans) {
		return nil, &UserCannotBorrowError{UserID: req.UserID}
	}

	if !book.IsAvailable(bookActiveLoans) {
		return nil, &BookNotAvailableError{BookID: req.BookID}
	}

	// 7. 新しい貸出エンティティを作成
	loanID := loandm.NewLoanID()
	loan := loandm.NewLoan(
		loanID,
		*userID,
		bookID,
		time.Now(),
		nil,
	)

	// 8. Loanのみ保存（User/Bookは更新なし - SSOT!）
	if err := uc.loanRepo.Save(ctx, loan); err != nil {
		return nil, err
	}

	// 9. レスポンスを返す（ミドルウェアが自動COMMIT）
	return NewBorrowBookResponse(loan, book.Title()), nil
}

// カスタムエラー
type UserCannotBorrowError struct {
	UserID string
}

func (e *UserCannotBorrowError) Error() string {
	return "user cannot borrow: limit reached"
}

type BookNotAvailableError struct {
	BookID string
}

func (e *BookNotAvailableError) Error() string {
	return "book is not available"
}
