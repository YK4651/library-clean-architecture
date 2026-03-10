package bookdm

import "context"

type IBookRepository interface {
	FindByID(ctx context.Context, id *BookID) (*Book, error)
}
