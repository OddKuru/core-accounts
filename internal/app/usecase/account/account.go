package account

import (
	"context"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/rco"
	"github.com/OddKuru/core-accounts/internal/domain/repository"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/OddKuru/core-accounts/pkg/logger"
	"github.com/OddKuru/core-accounts/pkg/utils"
)

var _ Account = (*UseCase)(nil)

const (
	CreateAccountVersion uint = 1
)

type (
	Account interface {
		Create(ctx context.Context, name, email, password string) error

		UpdateNameById(ctx context.Context, id, name string) error
		UpdateEmailById(ctx context.Context, id, email string) error
		UpdatePasswordById(ctx context.Context, id, oldPassword, password string) error
		UpdateRoleById(ctx context.Context, id string, role vo.AccountRoleType) error

		GetById(ctx context.Context, id string) (*aggregate.Account, error)
		GetByQuery(ctx context.Context, query rco.Query) (*rco.DataWithPageCount[*aggregate.Account], error)
	}

	UseCase struct {
		log             logger.Logger
		accountQuery    repository.AccountQuery
		accountCommand  repository.AccountCommand
		passwordManager repository.PasswordHasher
		idGenerator     repository.IDGenerator
		now             utils.TimeNower
	}
)

func NewUseCase(
	log logger.Logger,
	accountQuery repository.AccountQuery,
	accountCommand repository.AccountCommand,
	passwordManager repository.PasswordHasher,
	idGenerator repository.IDGenerator,
	nowTimer utils.TimeNower,
) (*UseCase, error) {
	uc := &UseCase{
		log:             log,
		accountQuery:    accountQuery,
		accountCommand:  accountCommand,
		passwordManager: passwordManager,
		idGenerator:     idGenerator,
		now:             nowTimer,
	}

	if err := uc.validate(); err != nil {
		return nil, errx.WrapWithCode(err, codex.InvalidArgument, "uc.validate")
	}

	return uc, nil
}

func (u *UseCase) validate() error {
	if u.log == nil {
		return errx.New(codex.InvalidArgument, "log is nil")
	}
	if u.accountQuery == nil {
		return errx.New(codex.InvalidArgument, "accountQuery is nil")
	}
	if u.accountCommand == nil {
		return errx.New(codex.InvalidArgument, "accountCommand is nil")
	}
	if u.passwordManager == nil {
		return errx.New(codex.InvalidArgument, "passwordManager is nil")
	}
	if u.idGenerator == nil {
		return errx.New(codex.InvalidArgument, "idGenerator is nil")
	}
	if u.now == nil {
		return errx.New(codex.InvalidArgument, "now is nil")
	}
	return nil
}
