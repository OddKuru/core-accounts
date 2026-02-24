package ports

import (
	"context"

	"github.com/OddKuru/core-accounts/internal/domain/repository"
)

type Repositories interface {
	AccountRepository() repository.AccountCommand
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(repositories Repositories) error) error
}
