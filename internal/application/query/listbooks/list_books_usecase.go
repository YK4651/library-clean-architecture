package listbooks

import (
	"context"

	"github.com/YK4651/library-clean-architecture/internal/application/query"
)

// ListBooksUseCase - ページネーション付き書籍一覧 Use Case (CQRSクエリサイド)
type ListBooksUseCase struct {
	bookQueryService query.BookQueryService
}

// NewListBooksUseCase - 新しいListBooksUseCaseを作成
func NewListBooksUseCase(bookQueryService query.BookQueryService) *ListBooksUseCase {
	return &ListBooksUseCase{
		bookQueryService: bookQueryService,
	}
}

// Execute - ユースケースを実行して書籍一覧を取得
func (uc *ListBooksUseCase) Execute(ctx context.Context, request *ListBooksRequest) *ListBooksResponse {
	list, err := uc.bookQueryService.ListBooks(ctx, request.Limit, request.Offset)
	if err != nil {
		return NewFailureResponse("An unexpected error occurred")
	}
	return NewSuccessResponse(list)
}
