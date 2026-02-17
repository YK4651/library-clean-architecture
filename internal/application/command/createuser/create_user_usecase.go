package createuser

import (
	"errors"
	"fmt"

	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

// 入力DTO
type CreateUserInput struct {
	Name  string
	Email string
}

// 出力DTO
type CreateUserOutput struct {
	ID                 string
	Name               string
	Email              string
	Status             string
	CurrentBorrowCount int
	OverdueFees        float64
}

type CreateUserUseCase struct {
	// コンストラクター: リポジトリインターフェースに依存
	userRepository userdm.UserRepository
}

func NewCreateUserUseCase(repo userdm.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: repo}
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
	// ビジネスルール: 重複メールをチェック
	existingUser, err := uc.userRepository.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("ユーザーは既に存在しています")
	}

	// ファクトリーメソッドでユーザーを作成
	u := userdm.NewUser(input.Name, input.Email)

	// リポジトリに永続化
	if err := uc.userRepository.Save(u); err != nil {
		return nil, err
	}

	// DTOを返す（新規作成時は貸出数0）
	return &CreateUserOutput{
		ID:                 u.Id().Value(),
		Name:               u.Name(),
		Email:              u.Email(),
		Status:             string(u.Status()),
		CurrentBorrowCount: 0,
		OverdueFees:        u.OverdueFees(),
	}, nil
}
