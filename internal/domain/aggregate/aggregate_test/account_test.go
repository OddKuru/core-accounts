package aggregate_test

import (
	"testing"
	"time"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/entity"
	"github.com/OddKuru/core-accounts/internal/domain/event"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func newValidEntityAccount(t *testing.T) *entity.Account {
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
	now := time.Now()

	acc, err := entity.NewAccount(id, name, email, password, role, 1, now, now)
	assert.NoError(t, err)
	return acc
}

func TestNewAccount_Success(t *testing.T) {
	eAcc := newValidEntityAccount(t)

	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, eAcc.ID().Value(), a.Account().ID().Value())
	assert.False(t, a.HasEvents())
	assert.Len(t, a.Events(), 0)
	assert.NotNil(t, a.LogData())
}

func TestNewAccount_InvalidEntity(t *testing.T) {
	acc := &entity.Account{}
	a, err := aggregate.NewAccount(acc)
	assert.Nil(t, a)
	assert.Error(t, err)
	assert.Equal(t, codex.InvalidArgument, errx.Code(err))
}

func TestAccount_ChangeName_NoChange_NoEvent(t *testing.T) {
	eAcc := newValidEntityAccount(t)
	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)

	tm := time.Now()
	name := eAcc.Name()
	err = a.ChangeName(name, tm)
	assert.NoError(t, err)
	assert.False(t, a.HasEvents())
	assert.Len(t, a.Events(), 0)
}

func TestAccount_ChangeName_WithEvent(t *testing.T) {
	eAcc := newValidEntityAccount(t)
	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)

	tm := time.Now()

	t.Run("Should create event on name change", func(t *testing.T) {
		newName, _ := vo.NewLoginName("new-name")
		err = a.ChangeName(newName, tm)
		assert.NoError(t, err)

		events := a.Events()
		assert.Len(t, events, 1)

		ev, ok := events[0].(event.AccountChangeName)
		assert.True(t, ok)
		assert.Equal(t, event.AccountChangeNameType, ev.Type())
		assert.Equal(t, "new-name", ev.Name().Value())
		assert.True(t, ev.Timestamp().Equal(tm))
		assert.True(t, a.HasEvents())
	})

	t.Run("Should error with empty name", func(t *testing.T) {
		var empty vo.LoginName
		err = a.ChangeName(empty, tm)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

func TestAccount_ChangeEmail_NoChange_NoEvent(t *testing.T) {
	eAcc := newValidEntityAccount(t)
	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)

	tm := time.Now()
	email := eAcc.Email()
	err = a.ChangeEmail(email, tm)
	assert.NoError(t, err)
	assert.Len(t, a.Events(), 0)
	assert.False(t, a.HasEvents())
}

func TestAccount_ChangeEmail_WithEvent(t *testing.T) {
	eAcc := newValidEntityAccount(t)
	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)

	tm := time.Now()

	t.Run("Should create event on email change", func(t *testing.T) {
		newEmail, _ := vo.NewEmail("new@example.com")
		err = a.ChangeEmail(newEmail, tm)
		assert.NoError(t, err)

		events := a.Events()
		assert.Len(t, events, 1)

		ev, ok := events[0].(event.AccountChangeEmail)
		assert.True(t, ok)
		assert.Equal(t, event.AccountChangeEmailType, ev.Type())
		assert.Equal(t, "new@example.com", ev.Email().Value())
	})

	t.Run("Should error with empty email", func(t *testing.T) {
		var empty vo.Email
		err = a.ChangeEmail(empty, tm)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

func TestAccount_ChangeRole_NoChange_NoEvent(t *testing.T) {
	eAcc := newValidEntityAccount(t)
	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)

	tm := time.Now()
	role := eAcc.Role()
	err = a.ChangeRole(role, tm)
	assert.NoError(t, err)
	assert.Len(t, a.Events(), 0)
	assert.False(t, a.HasEvents())
}

func TestAccount_ChangeRole_WithEvent(t *testing.T) {
	eAcc := newValidEntityAccount(t)
	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)

	tm := time.Now()

	t.Run("Should create event on role change", func(t *testing.T) {
		newRole, _ := vo.NewRole(vo.RoleAdmin)
		err = a.ChangeRole(newRole, tm)
		assert.NoError(t, err)

		events := a.Events()
		assert.Len(t, events, 1)

		ev, ok := events[0].(event.AccountChangeRole)
		assert.True(t, ok)
		assert.Equal(t, event.AccountChangeRoleType, ev.Type())
		assert.Equal(t, vo.RoleAdmin, ev.Role().Value())
	})

	t.Run("Should error with empty role", func(t *testing.T) {
		var empty vo.Role
		err = a.ChangeRole(empty, tm)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

func TestAccount_ChangePassword_Event(t *testing.T) {
	eAcc := newValidEntityAccount(t)
	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)

	tm := time.Now()

	t.Run("Should create event on password change", func(t *testing.T) {
		newPassword := vo.NewHashedPassword([]byte("new-hash"))
		err = a.ChangePassword(newPassword, tm)
		assert.NoError(t, err)

		events := a.Events()
		assert.Len(t, events, 1)

		ev, ok := events[0].(event.AccountChangePassword)
		assert.True(t, ok)
		assert.Equal(t, event.AccountChangePasswordType, ev.Type())
		assert.Equal(t, []byte("new-hash"), ev.Password().Value())
	})

	t.Run("Should error with empty password", func(t *testing.T) {
		empty := vo.NewHashedPassword(nil)
		err = a.ChangePassword(empty, tm)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

func TestAccount_ClearEvents(t *testing.T) {
	eAcc := newValidEntityAccount(t)
	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)

	tm := time.Now()
	newName, _ := vo.NewLoginName("new-name")
	err = a.ChangeName(newName, tm)
	assert.NoError(t, err)
	assert.Greater(t, len(a.Events()), 0)

	a.ClearEvent()
	assert.Len(t, a.Events(), 0)
	assert.False(t, a.HasEvents())
}

func TestAccount_VersionDelegation(t *testing.T) {
	eAcc := newValidEntityAccount(t)
	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)

	a.SetVersion(10)
	assert.Equal(t, uint(10), a.Account().Version())
	a.IncrementVersion()
	assert.Equal(t, uint(11), a.Account().Version())
}

func TestEventsHelpers_FromAggregate(t *testing.T) {
	eAcc := newValidEntityAccount(t)
	a, err := aggregate.NewAccount(eAcc)
	assert.NoError(t, err)

	tm := time.Now()
	newName, _ := vo.NewLoginName("new-name")
	err = a.ChangeName(newName, tm)
	assert.NoError(t, err)

	events := a.Events()
	assert.Equal(t, 1, events.Len())

	filtered := events.Filter(func(ev event.Event) bool {
		return ev.Type() == event.AccountChangeNameType
	})
	assert.Equal(t, 1, filtered.Len())

	logData := events.LogData()
	assert.NotNil(t, logData)
}

