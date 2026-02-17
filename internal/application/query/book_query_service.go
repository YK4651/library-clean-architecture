package query

import "context"

// BookQueryService - Query Serviceインターフェース (CQRSクエリサイド)
//
// Repository（エンティティを扱う）とは異なり、
// QueryServiceはRead Modelを直接返します。
//
// 注意：すべてのメソッドは context.Context を最初の引数として受け取ります。
// これにより、ミドルウェアが提供するDB接続を取得できます。
//
//go:generate mockgen -source=book_query_service.go -destination=mocks/mock_book_query_service.go -package=mocks
type BookQueryService interface {
	// GetBookByID - 現在の貸出ステータスを含めてIDで書籍を取得
	//
	// 見つかった場合、書籍のread modelを返し、そうでなければnilを返す
	GetBookByID(ctx context.Context, bookID string) (*BookReadModel, error)

	// ListBooks - すべての書籍を貸出ステータスと共に取得（Lesson 6用）
	//
	// 書籍read modelの配列を返す
	ListBooks(ctx context.Context) ([]*BookReadModel, error)
}
