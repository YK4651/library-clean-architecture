package userdm

import "github.com/YK4651/library-clean-architecture/internal/domain/shared"

// UserID はユーザーエンティティのID
type UserID shared.ID

// NewUserID は新しいUserIDを生成
func NewUserID() UserID {
	return UserID(shared.NewID())
}

// GenerateUserID は新しいUserIDのポインタを生成（ReconstructUser 等で使用）
func GenerateUserID() *UserID {
	id := NewUserID()
	return &id
}

// Value は文字列としての値を返す
func (id UserID) Value() string {
	return shared.ID(id).Value()
}

// Equals は値オブジェクトの等価性を比較
func (id UserID) Equals(other UserID) bool {
	return shared.ID(id).Equals(shared.ID(other))
}
