package listbooks

import (
	"errors"
	"strconv"
)

const (
	DefaultLimit  = 20
	MinLimit      = 1
	MaxLimit      = 100
	DefaultOffset = 0
)

// ListBooksRequest - ページネーション付き書籍一覧のリクエスト
type ListBooksRequest struct {
	Limit  int // 1-100、デフォルト20
	Offset int // >=0、デフォルト0
}

// NewListBooksRequest - クエリパラメータからリクエストを生成（バリデーション付き）
func NewListBooksRequest(limitParam, offsetParam string) (*ListBooksRequest, error) {
	limit := DefaultLimit
	if limitParam != "" {
		l, err := strconv.Atoi(limitParam)
		if err != nil {
			return nil, errors.New("limit must be a valid integer")
		}
		if l < MinLimit || l > MaxLimit {
			return nil, errors.New("limit must be between 1 and 100")
		}
		limit = l
	}

	offset := DefaultOffset
	if offsetParam != "" {
		o, err := strconv.Atoi(offsetParam)
		if err != nil {
			return nil, errors.New("offset must be a valid integer")
		}
		if o < 0 {
			return nil, errors.New("offset must be greater than or equal to 0")
		}
		offset = o
	}

	return &ListBooksRequest{Limit: limit, Offset: offset}, nil
}
