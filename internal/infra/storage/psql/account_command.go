package psql

import (
	"context"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/event"
	"github.com/OddKuru/core-accounts/internal/domain/repository"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.AccountCommand = (*AccountCommand)(nil)

type AccountCommand struct {
	pool *pgxpool.Pool
}

func NewAccountCommand(pool *pgxpool.Pool) (*AccountCommand, error) {
	return &AccountCommand{
		pool: pool,
	}, nil
}

var RoleValueToDB = map[vo.AccountRoleType]string{
	vo.RoleUnknown:    "user",
	vo.RoleSuperAdmin: "super_admin",
	vo.RoleAdmin:      "admin",
	vo.RoleUser:       "user",
}

func (a *AccountCommand) getSqlQuery(ctx context.Context) SQLQuerier {
	tx, ok := ctx.Value(transactionKey{}).(pgx.Tx)
	if !ok {
		return a.pool
	}
	return tx
}

func (a *AccountCommand) Create(ctx context.Context, account *aggregate.Account) error {
	db := a.getSqlQuery(ctx)
	_, err := db.Exec(
		ctx, AccountCreateQuery,
		account.Account().ID().Value(),
		account.Account().Name().Value(),
		account.Account().Email().Value(),
		account.Account().Password().Value(),
		account.Account().Version(),
		RoleValueToDB[account.Account().Role().Value()],
		account.Account().UpdatedAt(),
		account.Account().CreatedAt(),
	)
	if err != nil {
		return errx.WrapWithCode(err, codex.Internal, "[AccountCommand] db.Exec")
	}
	return nil
}

func (a *AccountCommand) UpdateByEvents(ctx context.Context, account *aggregate.Account) error {
	events := account.Events()
	for _, ev := range events {
		switch ev.Type() {
		case event.AccountChangeNameType:

		}
	}
	return nil
}
