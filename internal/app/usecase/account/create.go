package account

import (
	"context"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/app/ports"
	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/entity"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/OddKuru/core-accounts/pkg/logger"
	"github.com/pkg/errors"
)

var (
	ErrAccountNameAlreadyExists  = errx.New(codex.AlreadyExists, "name already exists")
	ErrAccountEmailAlreadyExists = errx.New(codex.AlreadyExists, "email already exists")
)

func (u *UseCase) Create(ctx context.Context, name, email, password string) (*aggregate.Account, error) {
	voName, err := vo.NewLoginName(name)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] vo.NewLoginName")
	}
	candidate, err := u.accountQuery.HasByName(ctx, voName)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] accountQuery.HasByName")
	}
	if candidate {
		return nil, ErrAccountNameAlreadyExists
	}
	voEmail, err := vo.NewEmail(email)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] vo.NewEmail")
	}
	hasEmail, err := u.accountQuery.HasByEmail(ctx, voEmail)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] accountQuery.HasByEmail")
	}
	if hasEmail {
		return nil, ErrAccountEmailAlreadyExists
	}

	voPassword, err := vo.NewPassword(password)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] vo.NewPassword")
	}

	hashedPassword, err := u.passwordManager.Hash(voPassword)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] passwordManager.Hash")
	}

	id, err := u.idGenerator.GenerateID()
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] idGenerator.GenerateID")
	}

	voRole, err := vo.NewRole(vo.RoleUser)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] vo.NewRole")
	}

	t := u.now.Now()

	acc, err := entity.NewAccount(id, voName, voEmail, hashedPassword, voRole, CreateAccountVersion, t, t)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] entity.NewAccount")
	}

	accAggregate, err := aggregate.NewAccount(acc)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] aggregate.NewAccount")
	}

	err = u.uow.Do(ctx, func(repos ports.Repositories) error {
		err := repos.AccountRepository().Save(ctx, accAggregate)
		if err != nil {
			u.log.Error(
				ctx,
				"create account error",
				logger.Any("aggregate", acc.LogData()),
			)
			return errors.Wrap(err, "[UseCase] accountCommand.Create")
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] uow.Do")
	}

	return accAggregate, nil
}
