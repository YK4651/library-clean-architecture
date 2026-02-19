package extendloan_test

import (
	"testing"

	"github.com/YK4651/library-clean-architecture/internal/application/extendloan"
)

func TestNewExtendLoanRequest_Valid(t *testing.T) {
	req, err := extendloan.NewExtendLoanRequest("01HXXX", "01HYYY", "12345678")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.BookID != "01HXXX" || req.LoanID != "01HYYY" || req.UserID != "12345678" {
		t.Errorf("unexpected request: %+v", req)
	}
}

func TestNewExtendLoanRequest_EmptyBookID(t *testing.T) {
	_, err := extendloan.NewExtendLoanRequest("", "01HYYY", "12345678")
	if err == nil {
		t.Fatal("expected error for empty bookID")
	}
}

func TestNewExtendLoanRequest_EmptyLoanID(t *testing.T) {
	_, err := extendloan.NewExtendLoanRequest("01HXXX", "", "12345678")
	if err == nil {
		t.Fatal("expected error for empty loanID")
	}
}

func TestNewExtendLoanRequest_EmptyUserID(t *testing.T) {
	_, err := extendloan.NewExtendLoanRequest("01HXXX", "01HYYY", "")
	if err == nil {
		t.Fatal("expected error for empty userID")
	}
}
