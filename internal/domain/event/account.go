package event

import (
	"time"

	"github.com/OddKuru/core-accounts/internal/domain/vo"
)

const (
	AccountChangeNameType     Type = "account_change_name"
	AccountChangeEmailType    Type = "account_change_email"
	AccountChangeRoleType     Type = "account_change_role"
	AccountChangePasswordType Type = "account_change_password"
)

var (
	_ Event = (*AccountChangeName)(nil)
)

type AccountChangeName struct {
	id        vo.ID
	name      vo.LoginName
	timestamp time.Time
}

func NewAccountChangeName(id vo.ID, name vo.LoginName, t time.Time) AccountChangeName {
	return AccountChangeName{
		id:        id,
		name:      name,
		timestamp: t,
	}
}

func (a AccountChangeName) ID() vo.ID            { return a.id }
func (a AccountChangeName) Type() Type           { return AccountChangeNameType }
func (a AccountChangeName) Value() any           { return a.name }
func (a AccountChangeName) Timestamp() time.Time { return a.timestamp }
func (a AccountChangeName) Name() vo.LoginName   { return a.name }

type AccountChangeEmail struct {
	id        vo.ID
	email     vo.Email
	timestamp time.Time
}

func NewAccountChangeEmail(id vo.ID, name vo.Email, t time.Time) AccountChangeEmail {
	return AccountChangeEmail{
		id:        id,
		email:     name,
		timestamp: t,
	}
}

func (a AccountChangeEmail) ID() vo.ID            { return a.id }
func (a AccountChangeEmail) Type() Type           { return AccountChangeEmailType }
func (a AccountChangeEmail) Value() any           { return a.email }
func (a AccountChangeEmail) Timestamp() time.Time { return a.timestamp }
func (a AccountChangeEmail) Email() vo.Email      { return a.email }

type AccountChangeRole struct {
	id        vo.ID
	role      vo.Role
	timestamp time.Time
}

func NewAccountChangeRole(id vo.ID, r vo.Role, t time.Time) AccountChangeRole {
	return AccountChangeRole{
		id:        id,
		role:      r,
		timestamp: t,
	}
}

func (a AccountChangeRole) ID() vo.ID            { return a.id }
func (a AccountChangeRole) Type() Type           { return AccountChangeRoleType }
func (a AccountChangeRole) Value() any           { return a.role }
func (a AccountChangeRole) Timestamp() time.Time { return a.timestamp }
func (a AccountChangeRole) Role() vo.Role        { return a.role }

type AccountChangePassword struct {
	id        vo.ID
	value     vo.HashedPassword
	timestamp time.Time
}

func NewAccountChangePassword(id vo.ID, pass vo.HashedPassword, t time.Time) AccountChangePassword {
	return AccountChangePassword{
		id:        id,
		value:     pass,
		timestamp: t,
	}
}

func (a AccountChangePassword) ID() vo.ID                   { return a.id }
func (a AccountChangePassword) Type() Type                  { return AccountChangePasswordType }
func (a AccountChangePassword) Value() any                  { return a.value }
func (a AccountChangePassword) Timestamp() time.Time        { return a.timestamp }
func (a AccountChangePassword) Password() vo.HashedPassword { return a.value }
