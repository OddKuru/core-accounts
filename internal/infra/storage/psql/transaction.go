package psql

import (
	"context"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var _ repository.Transactor = (*Transaction)(nil)

type transactionKey struct{}

type Transaction struct {
	pool *pgxpool.Pool
}

func NewTransaction(pool *pgxpool.Pool) (*Transaction, error) {
	return &Transaction{
		pool: pool,
	}, nil
}

func (t *Transaction) WithTransaction(
	ctx context.Context,
	fn func(txCtx context.Context,
) error) (err error) {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return errx.New(codex.Internal, "transaction begin failed")
	}
	defer func() {
		if err != nil {
			txErr := tx.Rollback(ctx)
			if txErr != nil {
				err = errx.WrapWithCode(err, codex.Internal, "[Transaction] tx.Rollback")
			}
		}
	}()
	txCtx := context.WithValue(ctx, transactionKey{}, tx)
	err = fn(txCtx)
	if err != nil {
		return errors.Wrap(err, "[Transaction] transaction failed")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errx.New(codex.Internal, "transaction commit failed")
	}
	return nil
}
