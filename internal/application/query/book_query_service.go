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

	// ListBooks - ページネーション付きで書籍を貸出ステータスと共に取得（Lesson 6）
	//
	// limit: 1-100（デフォルト20）, offset: >=0（デフォルト0）
	// isAvailable は loans テーブルから派生（Single Source of Truth）
	ListBooks(ctx context.Context, limit, offset int) (*BookListReadModel, error)
}
