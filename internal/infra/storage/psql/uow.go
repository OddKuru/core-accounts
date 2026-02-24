package psql

import (
	"context"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/app/ports"
	"github.com/OddKuru/core-accounts/internal/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var (
	_ ports.UnitOfWork   = (*UnitOfWork)(nil)
	_ ports.Repositories = (*Repositories)(nil)
)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

type Repositories struct {
	accountRepository repository.AccountCommand
}

func (r Repositories) AccountRepository() repository.AccountCommand {
	return r.accountRepository
}

func newRepositories(
	tx pgx.Tx,
) (Repositories, error) {
	accRepo, err := NewAccountCommand(tx)
	if err != nil {
		return Repositories{}, errors.Wrap(err, "NewAccountCommand")
	}
	return Repositories{
		accountRepository: accRepo,
	}, nil
}

func NewUnitOfWork(pool *pgxpool.Pool) (*UnitOfWork, error) {
	if pool == nil {
		return nil, errx.New(codex.InvalidArgument, "pool is nil")
	}
	return &UnitOfWork{pool: pool}, nil
}

func (u UnitOfWork) Do(ctx context.Context, fn func(repositories ports.Repositories) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return errx.WrapWithCode(err, codex.Internal, "[UnitOfWork] pool.Begin")
	}
	repos, err := newRepositories(tx)
	if err != nil {
		return err
	}
	err = fn(repos)
	if err != nil {
		rErr := tx.Rollback(ctx)
		if rErr != nil {
			return errors.Wrapf(err, "[UnitOfWork] tx.Rollback, Rollback error: %v", rErr)
		}
		return errors.Wrap(err, "unit of work func")
	}
	err = tx.Commit(ctx)
	if err != nil {
		return errx.WrapWithCode(err, codex.Internal, "[UnitOfWork] tx.Commit")
	}
	return nil
}
