package vo_test

import (
	"testing"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestVoRole(t *testing.T) {
	t.Run("Should correct roles", func(t *testing.T) {
		tests := []vo.AccountRoleType{
			vo.RoleSuperAdmin,
			vo.RoleAdmin,
			vo.RoleUser,
		}

		for _, rt := range tests {
			role, err := vo.NewRole(rt)
			assert.NoError(t, err)
			assert.Equal(t, rt, role.Value())
		}
	})

	t.Run("Should error with invalid roles", func(t *testing.T) {
		role, err := vo.NewRole(vo.RoleUnknown)
		assert.Equal(t, vo.Role{}, role)
		assert.Error(t, err)
		assert.Equal(t, codex.Internal, errx.Code(err))

		role, err = vo.NewRole(vo.AccountRoleType(100))
		assert.Equal(t, vo.Role{}, role)
		assert.Error(t, err)
		assert.Equal(t, codex.Internal, errx.Code(err))
	})

	t.Run("Should stringify known roles", func(t *testing.T) {
		assert.Equal(t, "super admin", vo.RoleSuperAdmin.String())
		assert.Equal(t, "admin", vo.RoleAdmin.String())
		assert.Equal(t, "user", vo.RoleUser.String())
	})

	t.Run("Should stringify unknown role as unknown", func(t *testing.T) {
		s := vo.AccountRoleType(999).String()
		assert.Equal(t, vo.RoleString[vo.RoleUnknown], s)
	})
}
