package getbook

import (
	"context"

	"github.com/YK4651/library-clean-architecture/internal/application/query"
)

// GetBookUseCase - GetBook Use Case (CQRSクエリサイド)
//
// このユースケースは非常にシンプル - QueryServiceに委譲するだけです。
// すべての複雑さはBookQueryServiceにあります！
type GetBookUseCase struct {
	bookQueryService query.BookQueryService // クエリサービス（インフラ層から注入）
}

// NewGetBookUseCase - 新しいGetBookUseCaseを作成
func NewGetBookUseCase(bookQueryService query.BookQueryService) *GetBookUseCase {
	return &GetBookUseCase{
		bookQueryService: bookQueryService,
	}
}

// Execute - ユースケースを実行して書籍情報を取得
func (uc *GetBookUseCase) Execute(ctx context.Context, request *GetBookRequest) *GetBookResponse {
	// 単にQueryServiceに委譲！contextを渡す
	bookReadModel, err := uc.bookQueryService.GetBookByID(ctx, request.BookID)

	if err != nil {
		// データベースエラーなどの予期しないエラー
		return NewFailureResponse("An unexpected error occurred")
	}

	if bookReadModel == nil {
		// 書籍が見つからない場合
		return NewNotFoundResponse()
	}

	// 成功：書籍情報を返す
	return NewSuccessResponse(bookReadModel)
}
