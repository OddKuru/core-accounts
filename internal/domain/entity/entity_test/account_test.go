package entity_test

import (
	"testing"
	"time"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/entity"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func newValidAccountVOs(t *testing.T) (vo.ID, vo.LoginName, vo.Email, vo.HashedPassword, vo.Role) {
	t.Helper()

	id, err := vo.NewID("550e8400-e29b-41d4-a716-446655440000")
	assert.NoError(t, err)
	name, err := vo.NewLoginName("user")
	assert.NoError(t, err)
	email, err := vo.NewEmail("user@example.com")
	assert.NoError(t, err)
	password := vo.NewHashedPassword([]byte("hash"))
	role, err := vo.NewRole(vo.RoleUser)
	assert.NoError(t, err)

	return id, name, email, password, role
}

func TestNewAccount_Success(t *testing.T) {
	id, name, email, password, role := newValidAccountVOs(t)
	now := time.Now()

	acc, err := entity.NewAccount(id, name, email, password, role, 1, now, now)
	assert.NoError(t, err)
	assert.NotNil(t, acc)

	assert.Equal(t, id.Value(), acc.ID().Value())
	assert.Equal(t, name.Value(), acc.Name().Value())
	assert.Equal(t, email.Value(), acc.Email().Value())
	assert.Equal(t, string(password.Value()), string(acc.Password().Value()))
	assert.Equal(t, role.Value(), acc.Role().Value())
	assert.Equal(t, uint(1), acc.Version())
	assert.False(t, acc.CreatedAt().IsZero())
	assert.False(t, acc.UpdatedAt().IsZero())
	assert.NotNil(t, acc.View())
	assert.NotNil(t, acc.LogData())
}

func TestNewAccount_InvalidID(t *testing.T) {
	name, email, password, role := func(t *testing.T) (vo.LoginName, vo.Email, vo.HashedPassword, vo.Role) {
		t.Helper()
		n, _ := vo.NewLoginName("user")
		e, _ := vo.NewEmail("user@example.com")
		p := vo.NewHashedPassword([]byte("hash"))
		r, _ := vo.NewRole(vo.RoleUser)
		return n, e, p, r
	}(t)

	var zeroID vo.ID
	now := time.Now()
	acc, err := entity.NewAccount(zeroID, name, email, password, role, 1, now, now)
	assert.Nil(t, acc)
	assert.Error(t, err)
	assert.Equal(t, codex.InvalidArgument, errx.Code(err))
}

func TestNewAccount_InvalidName(t *testing.T) {
	id, _, email, password, role := newValidAccountVOs(t)
	var zeroName vo.LoginName
	now := time.Now()

	acc, err := entity.NewAccount(id, zeroName, email, password, role, 1, now, now)
	assert.Nil(t, acc)
	assert.Error(t, err)
	assert.Equal(t, codex.InvalidArgument, errx.Code(err))
}

func TestNewAccount_InvalidEmail(t *testing.T) {
	id, name, _, password, role := newValidAccountVOs(t)
	var zeroEmail vo.Email
	now := time.Now()

	acc, err := entity.NewAccount(id, name, zeroEmail, password, role, 1, now, now)
	assert.Nil(t, acc)
	assert.Error(t, err)
	assert.Equal(t, codex.InvalidArgument, errx.Code(err))
}

func TestNewAccount_InvalidRole(t *testing.T) {
	id, name, email, password, _ := newValidAccountVOs(t)
	var zeroRole vo.Role
	now := time.Now()

	acc, err := entity.NewAccount(id, name, email, password, zeroRole, 1, now, now)
	assert.Nil(t, acc)
	assert.Error(t, err)
	assert.Equal(t, codex.InvalidArgument, errx.Code(err))
}

func TestNewAccount_InvalidPassword(t *testing.T) {
	id, name, email, _, role := newValidAccountVOs(t)
	var zeroPassword vo.HashedPassword
	now := time.Now()

	acc, err := entity.NewAccount(id, name, email, zeroPassword, role, 1, now, now)
	assert.Nil(t, acc)
	assert.Error(t, err)
	assert.Equal(t, codex.InvalidArgument, errx.Code(err))
}

func TestAccount_ChangeName(t *testing.T) {
	id, _, email, password, role := newValidAccountVOs(t)
	now := time.Now()

	name, _ := vo.NewLoginName("user")
	acc, err := entity.NewAccount(id, name, email, password, role, 1, now, now)
	assert.NoError(t, err)

	t.Run("Should change name", func(t *testing.T) {
		newName, _ := vo.NewLoginName("new-name")
		err := acc.ChangeName(newName)
		assert.NoError(t, err)
		assert.Equal(t, "new-name", acc.Name().Value())
	})

	t.Run("Should error with invalid name", func(t *testing.T) {
		var invalidName vo.LoginName
		err := acc.ChangeName(invalidName)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

func TestAccount_ChangeEmail(t *testing.T) {
	id, name, _, password, role := newValidAccountVOs(t)
	now := time.Now()

	email, _ := vo.NewEmail("user@example.com")
	acc, err := entity.NewAccount(id, name, email, password, role, 1, now, now)
	assert.NoError(t, err)

	t.Run("Should change email", func(t *testing.T) {
		newEmail, _ := vo.NewEmail("new@example.com")
		err := acc.ChangeEmail(newEmail)
		assert.NoError(t, err)
		assert.Equal(t, "new@example.com", acc.Email().Value())
	})

	t.Run("Should error with invalid email", func(t *testing.T) {
		var invalidEmail vo.Email
		err := acc.ChangeEmail(invalidEmail)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

func TestAccount_ChangePassword(t *testing.T) {
	id, name, email, _, role := newValidAccountVOs(t)
	now := time.Now()

	password := vo.NewHashedPassword([]byte("hash"))
	acc, err := entity.NewAccount(id, name, email, password, role, 1, now, now)
	assert.NoError(t, err)

	t.Run("Should change password", func(t *testing.T) {
		newPassword := vo.NewHashedPassword([]byte("new-hash"))
		err := acc.ChangePassword(newPassword)
		assert.NoError(t, err)
		assert.Equal(t, []byte("new-hash"), acc.Password().Value())
	})

	t.Run("Should error with invalid password", func(t *testing.T) {
		invalid := vo.NewHashedPassword(nil)
		err := acc.ChangePassword(invalid)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

func TestAccount_ChangeRole(t *testing.T) {
	id, name, email, password, _ := newValidAccountVOs(t)
	now := time.Now()

	role, _ := vo.NewRole(vo.RoleUser)
	acc, err := entity.NewAccount(id, name, email, password, role, 1, now, now)
	assert.NoError(t, err)

	t.Run("Should change role", func(t *testing.T) {
		newRole, _ := vo.NewRole(vo.RoleAdmin)
		err := acc.ChangeRole(newRole)
		assert.NoError(t, err)
		assert.Equal(t, vo.RoleAdmin, acc.Role().Value())
	})

	t.Run("Should error with invalid role", func(t *testing.T) {
		var invalidRole vo.Role
		err := acc.ChangeRole(invalidRole)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

func TestAccount_Version(t *testing.T) {
	id, name, email, password, role := newValidAccountVOs(t)
	now := time.Now()

	acc, err := entity.NewAccount(id, name, email, password, role, 1, now, now)
	assert.NoError(t, err)

	acc.SetVersion(10)
	assert.Equal(t, uint(10), acc.Version())

	acc.IncrementVersion()
	assert.Equal(t, uint(11), acc.Version())
}

