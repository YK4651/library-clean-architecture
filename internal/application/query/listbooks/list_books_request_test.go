package listbooks_test

import (
	"testing"

	"github.com/YK4651/library-clean-architecture/internal/application/query/listbooks"
)

func TestNewListBooksRequest_Defaults(t *testing.T) {
	req, err := listbooks.NewListBooksRequest("", "")
	if err != nil {
		t.Fatalf("Expected no error: %v", err)
	}
	if req.Limit != 20 {
		t.Errorf("Expected default limit 20, got %d", req.Limit)
	}
	if req.Offset != 0 {
		t.Errorf("Expected default offset 0, got %d", req.Offset)
	}
}

func TestNewListBooksRequest_ValidLimitAndOffset(t *testing.T) {
	req, err := listbooks.NewListBooksRequest("10", "5")
	if err != nil {
		t.Fatalf("Expected no error: %v", err)
	}
	if req.Limit != 10 || req.Offset != 5 {
		t.Errorf("Expected limit=10 offset=5, got limit=%d offset=%d", req.Limit, req.Offset)
	}
}

func TestNewListBooksRequest_InvalidLimitZero(t *testing.T) {
	_, err := listbooks.NewListBooksRequest("0", "")
	if err == nil {
		t.Error("Expected error for limit 0")
	}
	if err != nil && err.Error() != "limit must be between 1 and 100" {
		t.Errorf("Expected limit validation error, got: %v", err)
	}
}

func TestNewListBooksRequest_InvalidLimitOver100(t *testing.T) {
	_, err := listbooks.NewListBooksRequest("101", "")
	if err == nil {
		t.Error("Expected error for limit 101")
	}
}

func TestNewListBooksRequest_NegativeOffset(t *testing.T) {
	_, err := listbooks.NewListBooksRequest("", "-1")
	if err == nil {
		t.Error("Expected error for negative offset")
	}
	if err != nil && err.Error() != "offset must be greater than or equal to 0" {
		t.Errorf("Expected offset validation error, got: %v", err)
	}
}

func TestNewListBooksRequest_InvalidLimitNotInteger(t *testing.T) {
	_, err := listbooks.NewListBooksRequest("abc", "")
	if err == nil {
		t.Error("Expected error for non-integer limit")
	}
}

func TestNewListBooksRequest_InvalidOffsetNotInteger(t *testing.T) {
	_, err := listbooks.NewListBooksRequest("", "x")
	if err == nil {
		t.Error("Expected error for non-integer offset")
	}
}
