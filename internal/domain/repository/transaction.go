package repository

import "context"

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
