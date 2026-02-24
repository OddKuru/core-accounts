package psql

import (
	"context"
	"time"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/entity"
	"github.com/OddKuru/core-accounts/internal/domain/rco"
	"github.com/OddKuru/core-accounts/internal/domain/repository"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var _ repository.AccountQuery = (*AccountQuery)(nil)

type AccountQuery struct {
	pool *pgxpool.Pool
}

func NewAccountQuery(pool *pgxpool.Pool) (*AccountQuery, error) {
	return &AccountQuery{
		pool: pool,
	}, nil
}

func (a *AccountQuery) HasByEmail(ctx context.Context, email vo.Email) (bool, error) {
	var exists bool
	err := a.pool.QueryRow(ctx, AccountHasByEmailQuery, email.Value()).Scan(&exists)
	if err != nil {
		return false, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] pool.QueryRow")
	}
	return exists, nil
}

func (a *AccountQuery) HasByName(ctx context.Context, name vo.LoginName) (bool, error) {
	var exists bool
	err := a.pool.QueryRow(ctx, AccountHasByNameQuery, name.Value()).Scan(&exists)
	if err != nil {
		return false, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] pool.QueryRow")
	}
	return exists, nil
}

func (a *AccountQuery) GetById(ctx context.Context, id vo.ID) (*aggregate.Account, error) {
	var (
		name     string
		email    string
		password string
		role     string
		version  uint
		created  time.Time
		updated  time.Time
	)

	err := a.pool.QueryRow(ctx, AccountGetByIdQuery, id.Value()).
		Scan(&name, &email, &password, &version, &created, &updated, &role)
	if err != nil {
		return nil, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] pool.QueryRow")
	}

	voName, err := vo.NewLoginName(name)
	if err != nil {
		return nil, errors.Wrap(err, "[AccountQuery] vo.NewLoginName")
	}
	voEmail, err := vo.NewEmail(email)
	if err != nil {
		return nil, errors.Wrap(err, "[AccountQuery] vo.NewEmail")
	}
	voPassword := vo.NewHashedPassword([]byte(password))
	voRole, err := vo.NewRole(RoleValueFromDB[role])
	if err != nil {
		return nil, errors.Wrap(err, "[AccountQuery] vo.NewRole")
	}

	accEntity, err := entity.NewAccount(id, voName, voEmail, voPassword, voRole, version, created, updated)
	if err != nil {
		return nil, errors.Wrap(err, "[AccountQuery] entity.NewAccount")
	}
	agg, err := aggregate.NewAccount(accEntity)
	if err != nil {
		return nil, errors.Wrap(err, "[AccountQuery] aggregate.NewAccount")
	}

	return agg, nil
}

func (a *AccountQuery) GetByQuery(
	ctx context.Context, query rco.Query,
) (*rco.DataWithPageCount[*aggregate.Account], error) {
	return nil, nil
}
