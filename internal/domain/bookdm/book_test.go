package bookdm_test

import (
	"testing"

	"github.com/YK4651/library-clean-architecture/internal/domain/bookdm"
)

func TestCanCreateBook(t *testing.T) {
	// Create Value Objects
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}

	b, err := bookdm.NewBook(&bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	// ID should be auto-generated ULID (26 characters)
	if len(b.Id().Value()) != 26 {
		t.Errorf("Expected ULID length 26, got %d", len(b.Id().Value()))
	}
	if b.Title() != "Clean Architecture" {
		t.Errorf("Expected title 'Clean Architecture', got '%s'", b.Title())
	}
	if b.AvailableCopies() != 3 {
		t.Errorf("Expected 3 copies, got %d", b.AvailableCopies())
	}
}

func TestCanBorrowWhenCopiesAvailable(t *testing.T) {
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(&bookID, "Clean Architecture", "Robert Martin", isbn, 2)
	if err != nil {
		t.Fatal(err)
	}

	if !b.Available() {
		t.Error("Expected book to be available for borrowing")
	}
}

func TestBorrowDecreasesAvailableCopies(t *testing.T) {
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(&bookID, "Clean Architecture", "Robert Martin", isbn, 2)
	if err != nil {
		t.Fatal(err)
	}

	err = b.BorrowCopy()
	if err != nil {
		t.Fatalf("Borrow failed: %v", err)
	}

	if b.AvailableCopies() != 1 {
		t.Errorf("Expected 1 copy remaining, got %d", b.AvailableCopies())
	}
}

func TestCannotBorrowWhenNoCopies(t *testing.T) {
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(&bookID, "Clean Architecture", "Robert Martin", isbn, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Borrow the only copy
	err = b.BorrowCopy()
	if err != nil {
		t.Fatal(err)
	}

	// Try to borrow again when no copies available
	err = b.BorrowCopy()
	if err == nil {
		t.Error("Expected error when borrowing with 0 available copies")
	}
}
