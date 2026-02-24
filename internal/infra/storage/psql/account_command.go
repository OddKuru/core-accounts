package psql

import (
	"context"
	"time"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/entity"
	"github.com/OddKuru/core-accounts/internal/domain/repository"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

var _ repository.AccountCommand = (*AccountCommand)(nil)

var ErrOptimisticLock = errx.New(codex.Aborted, "version not equal")

type AccountCommand struct {
	tx pgx.Tx
}

func NewAccountCommand(
	tx pgx.Tx,
) (*AccountCommand, error) {
	return &AccountCommand{
		tx: tx,
	}, nil
}

func (a *AccountCommand) Load(ctx context.Context, id vo.ID) (*aggregate.Account, error) {
	var (
		name     string
		email    string
		password string
		role     string
		version  uint
		created  time.Time
		updated  time.Time
	)

	err := a.tx.QueryRow(ctx, AccountGetByIdQuery, id.Value()).
		Scan(&name, &email, &password, &version, &created, &updated, &role)
	if err != nil {
		return nil, errx.WrapWithCode(err, codex.Aborted, "failed to load account")
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

func (a *AccountCommand) Save(ctx context.Context, account *aggregate.Account) error {
	has, err := a.hasByID(ctx, account.Account().ID())
	if err != nil {
		return err
	}
	if !has {
		return a.create(ctx, account)
	}
	err = a.update(ctx, account)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountCommand) hasByID(ctx context.Context, id vo.ID) (bool, error) {
	var exists bool
	err := a.tx.QueryRow(ctx, AccountHasByIDQuery, id.Value()).Scan(&exists)
	if err != nil {
		return false, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] tx.QueryRow")
	}
	return exists, nil
}

func (a *AccountCommand) create(ctx context.Context, account *aggregate.Account) error {
	_, err := a.tx.Exec(
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
		return errx.WrapWithCode(err, codex.Internal, "[AccountCommand] tx.Exec")
	}
	return nil
}

func (a *AccountCommand) update(ctx context.Context, account *aggregate.Account) error {
	tag, err := a.tx.Exec(ctx, AccountUpdateQuery,
		account.Account().ID().Value(),
		account.Account().Name().Value(),
		account.Account().Email().Value(),
		account.Account().Password().Value(),
		account.Account().UpdatedAt(),
		RoleValueToDB[account.Account().Role().Value()],
		account.Account().Version(),
	)
	if err != nil {
		return errx.WrapWithCode(err, codex.Internal, "[AccountCommand] tx.Exec")
	}
	if tag.RowsAffected() == 0 {
		return ErrOptimisticLock
	}
	return nil
}

//
//func (a *AccountCommand) updateNameByID(
//	ctx context.Context,
//	name event.AccountChangeName,
//	version uint,
//) error {
//	tag, err := a.tx.Exec(ctx, AccountLoginUpdateQuery, name.ID().Value(), name.Value(), version)
//	if err != nil {
//		return errx.WrapWithCode(err, codex.Internal, "[AccountCommand] tx.Exec")
//	}
//	if tag.RowsAffected() == 0 {
//		return ErrOptimisticLock
//	}
//	return nil
//}
//
//func (a *AccountCommand) updateEmailByID(
//	ctx context.Context,
//	email event.AccountChangeEmail,
//	version uint,
//) error {
//	tag, err := a.tx.Exec(ctx, AccountEmailUpdateQuery, email.ID().Value(), email.Value(), version)
//	if err != nil {
//		return errx.WrapWithCode(err, codex.Internal, "[AccountCommand] tx.Exec")
//	}
//	if tag.RowsAffected() == 0 {
//		return ErrOptimisticLock
//	}
//	return nil
//}
//
//func (a *AccountCommand) updateRoleByID(
//	ctx context.Context,
//	role event.AccountChangeRole,
//	version uint,
//) error {
//	tag, err := a.tx.Exec(
//		ctx, AccountRoleUpdateQuery, role.ID().Value(), RoleValueToDB[role.Role().Value()], version,
//	)
//	if err != nil {
//		return errx.WrapWithCode(err, codex.Internal, "[AccountCommand] tx.Exec")
//	}
//	if tag.RowsAffected() == 0 {
//		return ErrOptimisticLock
//	}
//	return nil
//}
//
//func (a *AccountCommand) updatePasswordByID(
//	ctx context.Context,
//	pass event.AccountChangePassword,
//	version uint,
//) error {
//	tag, err := a.tx.Exec(ctx, AccountPasswordUpdateQuery, pass.ID().Value(), pass.Value(), version)
//	if err != nil {
//		return errx.WrapWithCode(err, codex.Internal, "[AccountCommand] tx.Exec")
//	}
//	if tag.RowsAffected() == 0 {
//		return ErrOptimisticLock
//	}
//	return nil
//}
