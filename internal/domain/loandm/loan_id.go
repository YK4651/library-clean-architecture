package loandm

import (
	"github.com/YK4651/library-clean-architecture/internal/domain/shared"
	"github.com/oklog/ulid/v2"
)

// LoanID は貸出エンティティのID
type LoanID shared.ID

// NewLoanID は新しいLoanIDを生成
func NewLoanID() LoanID {
	return LoanID(shared.NewID())
}

// LoanIDFromString はULID文字列からLoanIDを復元する（永続化層・リクエスト用）
func LoanIDFromString(s string) (LoanID, error) {
	_, err := ulid.Parse(s)
	if err != nil {
		return LoanID(""), err
	}
	return LoanID(shared.ID(s)), nil
}

// Value は文字列としての値を返す
func (id LoanID) Value() string {
	return shared.ID(id).Value()
}

// Equals は値オブジェクトの等価性を比較
func (id LoanID) Equals(other LoanID) bool {
	return shared.ID(id).Equals(shared.ID(other))
}
