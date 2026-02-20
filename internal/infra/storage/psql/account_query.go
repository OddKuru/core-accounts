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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var _ repository.AccountQuery = (*AccountQuery)(nil)

type SQLQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type TransactionExtractor interface {
	ExtractTransaction(ctx context.Context) (*pgx.Tx, error)
}

type AccountQuery struct {
	pool      *pgxpool.Pool
	txManager TransactionExtractor
}

func NewAccountQuery(pool *pgxpool.Pool) (*AccountQuery, error) {
	return &AccountQuery{
		pool: pool,
	}, nil
}

var RoleValueFromDB = map[string]vo.AccountRoleType{
	"super_admin": vo.RoleSuperAdmin,
	"admin":       vo.RoleAdmin,
	"user":        vo.RoleUser,
}

func (a *AccountQuery) getSqlQuery(ctx context.Context) SQLQuerier {
	tx, ok := ctx.Value(transactionKey{}).(pgx.Tx)
	if !ok {
		return a.pool
	}
	return tx
}

func (a *AccountQuery) HasByEmail(ctx context.Context, email vo.Email) (bool, error) {
	db := a.getSqlQuery(ctx)

	var exists bool
	err := db.QueryRow(ctx, AccountHasByEmailQuery, email.Value()).Scan(&exists)
	if err != nil {
		return false, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] db.QueryRow")
	}
	return exists, nil
}

func (a *AccountQuery) HasByName(ctx context.Context, name vo.LoginName) (bool, error) {
	db := a.getSqlQuery(ctx)

	var exists bool
	err := db.QueryRow(ctx, AccountHasByNameQuery, name.Value()).Scan(&exists)
	if err != nil {
		return false, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] db.QueryRow")
	}
	return exists, nil
}

func (a *AccountQuery) GetById(ctx context.Context, id vo.ID) (*aggregate.Account, error) {
	db := a.getSqlQuery(ctx)
	var (
		name     string
		email    string
		password string
		role     string
		version  uint
		created  time.Time
		updated  time.Time
	)

	err := db.QueryRow(ctx, AccountGetByIdQuery, id.Value()).Scan(&name, &email, &password, &version, &created, &updated, &role)
	
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

func (a *AccountQuery) GetByQuery(ctx context.Context, query rco.Query) (*rco.DataWithPageCount[*aggregate.Account], error) {
	return nil, nil
}
