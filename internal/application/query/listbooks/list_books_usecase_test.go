package listbooks_test

import (
	"context"
	"testing"

	"github.com/YK4651/library-clean-architecture/internal/application/query"
	"github.com/YK4651/library-clean-architecture/internal/application/query/listbooks"
	"github.com/YK4651/library-clean-architecture/internal/application/query/mocks"
	"go.uber.org/mock/gomock"
)

func TestListBooksUseCase_ReturnsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryService := mocks.NewMockBookQueryService(ctrl)
	uc := listbooks.NewListBooksUseCase(mockQueryService)

	ctx := context.Background()
	req, err := listbooks.NewListBooksRequest("10", "0")
	if err != nil {
		t.Fatal(err)
	}

	expected := &query.BookListReadModel{
		Books: []*query.BookReadModel{
			{ID: "b-1", ISBN: "978-0-13-468599-1", Title: "Clean Architecture", Author: "Robert C. Martin", IsAvailable: true, CurrentLoan: nil},
		},
		Total:  1,
		Limit:  10,
		Offset: 0,
	}

	mockQueryService.EXPECT().
		ListBooks(ctx, 10, 0).
		Return(expected, nil)

	resp := uc.Execute(ctx, req)

	if !resp.Success {
		t.Errorf("Expected success, got failure: %v", resp.ErrorMessage)
	}
	if resp.Data.Total != 1 {
		t.Errorf("Expected total 1, got %d", resp.Data.Total)
	}
	if len(resp.Data.Books) != 1 {
		t.Errorf("Expected 1 book, got %d", len(resp.Data.Books))
	}
}

func TestListBooksUseCase_EmptyListSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryService := mocks.NewMockBookQueryService(ctrl)
	uc := listbooks.NewListBooksUseCase(mockQueryService)

	ctx := context.Background()
	req, _ := listbooks.NewListBooksRequest("", "")

	mockQueryService.EXPECT().
		ListBooks(ctx, 20, 0).
		Return(&query.BookListReadModel{Books: []*query.BookReadModel{}, Total: 0, Limit: 20, Offset: 0}, nil)

	resp := uc.Execute(ctx, req)

	if !resp.Success {
		t.Error("Expected success for empty list (200 OK)")
	}
	if resp.Data.Total != 0 {
		t.Errorf("Expected total 0, got %d", resp.Data.Total)
	}
}

func TestListBooksUseCase_QueryErrorReturnsFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryService := mocks.NewMockBookQueryService(ctrl)
	uc := listbooks.NewListBooksUseCase(mockQueryService)

	ctx := context.Background()
	req, _ := listbooks.NewListBooksRequest("10", "0")

	mockQueryService.EXPECT().
		ListBooks(ctx, 10, 0).
		Return(nil, context.DeadlineExceeded)

	resp := uc.Execute(ctx, req)

	if resp.Success {
		t.Error("Expected failure on query error")
	}
	if resp.ErrorMessage == nil || *resp.ErrorMessage != "An unexpected error occurred" {
		t.Errorf("Expected unexpected error message, got: %v", resp.ErrorMessage)
	}
}
