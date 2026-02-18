package userdm

import "context"

type IUserRepository interface {
	FindByID(ctx context.Context, id *UserID) (*User, error)
}
