package bookdm

import "github.com/YK4651/library-clean-architecture/internal/domain/shared"

// BookID は書籍エンティティのID
// なぜ型エイリアス？shared.IDの機能を継承しつつ、型安全性を確保
type BookID shared.ID

// NewBookID は新しいBookIDを生成
func NewBookID() BookID {
	return BookID(shared.NewID())
}

// Value は文字列としての値を返す
func (id BookID) Value() string {
	return shared.ID(id).Value()
}

// Equals は値オブジェクトの等価性を比較
func (id BookID) Equals(other BookID) bool {
	return shared.ID(id).Equals(shared.ID(other))
}
