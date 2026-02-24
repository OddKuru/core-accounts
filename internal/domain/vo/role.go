package vo

import (
	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	validation "github.com/go-ozzo/ozzo-validation"
)

type AccountRoleType int

const (
	RoleUnknown AccountRoleType = iota
	RoleSuperAdmin
	RoleAdmin
	RoleUser
)

var RoleString = map[AccountRoleType]string{
	RoleUnknown:    "unknown",
	RoleSuperAdmin: "super admin",
	RoleAdmin:      "admin",
	RoleUser:       "user",
}

type Role struct {
	role AccountRoleType
}

func NewRole(role AccountRoleType) (Role, error) {
	result := Role{role}

	if err := result.Validate(); err != nil {
		return Role{}, errx.WrapWithCode(err, codex.Internal, "Role.Validate")
	}

	return result, nil
}

func (r Role) Value() AccountRoleType {
	return r.role
}

func (r Role) Validate() error {
	return validation.Validate(
		r.role,
		validation.Required,
		validation.In(RoleSuperAdmin, RoleAdmin, RoleUser).Error("role invalid value"),
	)
}

func (r AccountRoleType) String() string {
	value, ok := RoleString[r]
	if !ok {
		return RoleString[RoleUnknown]
	}
	return value
}
