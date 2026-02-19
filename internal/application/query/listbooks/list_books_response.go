package listbooks

import "github.com/YK4651/library-clean-architecture/internal/application/query"

// ListBooksResponse - ページネーション付き書籍一覧のレスポンス
type ListBooksResponse struct {
	Success      bool
	Data         *query.BookListReadModel
	ErrorMessage *string
}

// NewSuccessResponse - 成功レスポンスを生成
func NewSuccessResponse(data *query.BookListReadModel) *ListBooksResponse {
	return &ListBooksResponse{
		Success:      true,
		Data:         data,
		ErrorMessage: nil,
	}
}

// NewFailureResponse - 失敗レスポンスを生成
func NewFailureResponse(message string) *ListBooksResponse {
	return &ListBooksResponse{
		Success:      false,
		Data:         nil,
		ErrorMessage: &message,
	}
}
