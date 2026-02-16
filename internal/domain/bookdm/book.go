package bookdm

import "errors"

// Bookエンティティ - 図書館の本を表す
type Book struct {
	id              *BookID
	title           string
	author          string
	isbn            *ISBN
	totalCopies     int
	availableCopies int // 仮実装: 本来は Loan テーブルから算出 (SSOT)
}

func NewBook(
	id *BookID,
	title string,
	author string,
	isbn *ISBN,
	totalCopies int,
) (*Book, error) {
	// コンストラクタでの検証（フェイルファスト）
	if len(title) == 0 {
		return nil, errors.New("book title cannot be empty")
	}
	if len(author) == 0 {
		return nil, errors.New("book author cannot be empty")
	}
	if totalCopies < 1 {
		return nil, errors.New("total copies must be at least 1")
	}

	return &Book{
		id:              id,
		title:           title,
		author:          author,
		isbn:            isbn,
		totalCopies:     totalCopies,
		availableCopies: totalCopies,
	}, nil
}

// ゲッター（Goの慣習: "Get"プレフィックスなし）
func (b *Book) Id() *BookID {
	return b.id
}

func (b *Book) Title() string {
	return b.title
}

// GetTitle returns title (仮実装: LoanEligibilityService.IneligibilityReason 用)
func (b *Book) GetTitle() string {
	return b.Title()
}

func (b *Book) Author() string {
	return b.author
}

func (b *Book) ISBN() *ISBN {
	return b.isbn
}

func (b *Book) TotalCopies() int {
	return b.totalCopies
}

// ビジネスロジックメソッド
// Book entity focuses on RULES, not STATE
// Available copies derived from Loan table (Single Source of Truth)

// IsAvailable validates if book has available copies based on provided active loan count
// The count is provided by the use case (derived from Loan table)
func (b *Book) IsAvailable(currentActiveLoans int) bool {
	return currentActiveLoans < b.totalCopies
}

// Available returns whether book has any copies (仮実装: 現在の貸出数は 0 とみなす。本来は UseCase が IsAvailable(count) に渡す)
func (b *Book) Available() bool {
	return b.availableCopies > 0
}

// AvailableCopies returns available copy count (仮実装: テスト用。本来は Repository から取得)
func (b *Book) AvailableCopies() int {
	return b.availableCopies
}

// BorrowCopy decrements available copies (仮実装: テスト用。本来は Loan を追加するだけ)
func (b *Book) BorrowCopy() error {
	if b.availableCopies <= 0 {
		return errors.New("no copies available to borrow")
	}
	b.availableCopies--
	return nil
}

// State changes are tracked in Loan table (no BorrowCopy/ReturnCopy methods needed)

// 状態変更メソッド
func (b *Book) UpdateTitle(newTitle string) error {
	if len(newTitle) == 0 {
		return errors.New("book title cannot be empty")
	}
	b.title = newTitle
	return nil
}

func (b *Book) UpdateAuthor(newAuthor string) error {
	if len(newAuthor) == 0 {
		return errors.New("book author cannot be empty")
	}
	b.author = newAuthor
	return nil
}
