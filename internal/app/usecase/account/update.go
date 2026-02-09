package account

import (
	"context"

	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/pkg/errors"
)

func (u *UseCase) UpdateNameById(ctx context.Context, id, name string) error {
	acc, err := u.GetById(ctx, id)
	if err != nil {
		return errors.Wrap(err, "[UseCase] GetById")
	}

	voName, err := vo.NewLoginName(name)
	if err != nil {
		return errors.Wrap(err, "[UseCase] vo.NewLoginName")
	}
	t := u.now.Now()

	err = acc.ChangeName(voName, t)
	if err != nil {
		return errors.Wrap(err, "[UseCase] acc.ChangeName")
	}

	err = u.accountCommand.UpdateByEvents(ctx, acc)
	if err != nil {
		return errors.Wrap(err, "[UseCase] accountCommand.UpdateByEvents")
	}
	return nil
}

func (u *UseCase) UpdateEmailById(ctx context.Context, id, email string) error {
	acc, err := u.GetById(ctx, id)
	if err != nil {
		return errors.Wrap(err, "[UseCase] GetById")
	}

	voEmail, err := vo.NewEmail(email)
	if err != nil {
		return errors.Wrap(err, "[UseCase] vo.NewEmail")
	}
	t := u.now.Now()

	err = acc.ChangeEmail(voEmail, t)
	if err != nil {
		return errors.Wrap(err, "[UseCase] acc.ChangeEmail")
	}

	err = u.accountCommand.UpdateByEvents(ctx, acc)
	if err != nil {
		return errors.Wrap(err, "[UseCase] accountCommand.UpdateByEvents")
	}
	return nil
}

func (u *UseCase) UpdatePasswordById(ctx context.Context, id, oldPassword, password string) error {
	acc, err := u.GetById(ctx, id)
	if err != nil {
		return errors.Wrap(err, "[UseCase] GetById")
	}

	voOldPassword, err := vo.NewPassword(oldPassword)
	if err != nil {
		return errors.Wrap(err, "[UseCase] vo.NewPassword")
	}
	err = u.passwordManager.Compare(acc.Account().Password(), voOldPassword)
	if err != nil {
		return errors.Wrap(err, "[UseCase] passwordManager.Compare")
	}

	voPassword, err := vo.NewPassword(password)
	if err != nil {
		return errors.Wrap(err, "[UseCase] vo.NewPassword")
	}
	hashedPassword, err := u.passwordManager.Hash(voPassword)
	if err != nil {
		return errors.Wrap(err, "[UseCase] passwordManager.Hash")
	}
	t := u.now.Now()

	err = acc.ChangePassword(hashedPassword, t)
	if err != nil {
		return errors.Wrap(err, "[UseCase] acc.ChangePassword")
	}
	
	err = u.accountCommand.UpdateByEvents(ctx, acc)
	if err != nil {
		return errors.Wrap(err, "[UseCase] accountCommand.UpdateByEvents")
	}
	return nil
}

func (u *UseCase) UpdateRoleById(ctx context.Context, id string, role vo.AccountRoleType) error {
	acc, err := u.GetById(ctx, id)
	if err != nil {
		return errors.Wrap(err, "[UseCase] GetById")
	}

	voRole, err := vo.NewRole(role)
	if err != nil {
		return errors.Wrap(err, "[UseCase] vo.NewEmail")
	}
	t := u.now.Now()

	err = acc.ChangeRole(voRole, t)
	if err != nil {
		return errors.Wrap(err, "[UseCase] acc.ChangeRole")
	}

	err = u.accountCommand.UpdateByEvents(ctx, acc)
	if err != nil {
		return errors.Wrap(err, "[UseCase] accountCommand.UpdateByEvents")
	}
	return nil
}
