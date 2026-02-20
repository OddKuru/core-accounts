package entity

import (
	"time"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Account struct {
	id        vo.ID
	name      vo.LoginName
	email     vo.Email
	role      vo.Role
	password  vo.HashedPassword
	version   uint
	createdAt time.Time
	updatedAt time.Time
}

type AccountViewer interface {
	ID() vo.ID
	Name() vo.LoginName
	Email() vo.Email
	Role() vo.Role
	Password() vo.HashedPassword
	Version() uint
	UpdatedAt() time.Time
	CreatedAt() time.Time
}

func NewAccount(
	id vo.ID,
	name vo.LoginName,
	email vo.Email,
	password vo.HashedPassword,
	role vo.Role,
	version uint,
	createdAt time.Time,
	updatedAt time.Time,
) (*Account, error) {
	result := &Account{
		id:        id,
		name:      name,
		email:     email,
		password:  password,
		role:      role,
		version:   version,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	if err := result.Validate(); err != nil {
		return nil, errx.WrapWithCode(err, codex.InvalidArgument, "Account.Validate")
	}

	return result, nil
}

func (a *Account) ID() vo.ID {
	return a.id
}

func (a *Account) Name() vo.LoginName {
	return a.name
}

func (a *Account) Email() vo.Email {
	return a.email
}

func (a *Account) Role() vo.Role {
	return a.role
}

func (a *Account) Password() vo.HashedPassword {
	return a.password
}

func (a *Account) Version() uint {
	return a.version
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Account) ChangeName(name vo.LoginName) error {
	if err := name.Validate(); err != nil {
		return errx.WrapWithCode(err, codex.InvalidArgument, "[Account] name.Validate")
	}
	a.name = name
	return nil
}

func (a *Account) ChangeEmail(email vo.Email) error {
	if err := email.Validate(); err != nil {
		return errx.WrapWithCode(err, codex.InvalidArgument, "[Account] email.Validate")
	}
	a.email = email
	return nil
}

func (a *Account) ChangePassword(password vo.HashedPassword) error {
	if err := password.Validate(); err != nil {
		return errx.WrapWithCode(err, codex.InvalidArgument, "[Account] password.Validate")
	}
	a.password = password
	return nil
}

func (a *Account) ChangeRole(role vo.Role) error {
	if err := role.Validate(); err != nil {
		return errx.WrapWithCode(err, codex.InvalidArgument, "[Account] role.Validate")
	}
	a.role = role
	return nil
}

func (a *Account) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.id, validation.Required),
		validation.Field(&a.name, validation.Required),
		validation.Field(&a.email, validation.Required),
		validation.Field(&a.role, validation.Required),
		validation.Field(&a.password, validation.Required),
	)
}

func (a *Account) View() AccountViewer {
	return a
}

func (a *Account) LogData() any {
	return map[string]any{
		"id":       a.id.Value(),
		"name":     a.name.Value(),
		"email":    "<hidden>",
		"role":     "<hidden>",
		"password": "<hidden>",
		"version":  a.Version(),
		"created":  a.CreatedAt(),
		"updated":  a.UpdatedAt(),
	}
}
