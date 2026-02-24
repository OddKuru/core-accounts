package vo_test

import (
	"testing"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestVoPassword(t *testing.T) {
	t.Run("Should correct constructor password", func(t *testing.T) {
		value := "12345678"
		voPass, err := vo.NewPassword(value)
		assert.NoError(t, err)
		assert.Equal(t, value, voPass.Value())
	})

	t.Run("Should error with empty password", func(t *testing.T) {
		voPass, err := vo.NewPassword("")
		assert.Equal(t, vo.Password{}, voPass)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})

	t.Run("Should error with too short password", func(t *testing.T) {
		voPass, err := vo.NewPassword("short")
		assert.Equal(t, vo.Password{}, voPass)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})

	t.Run("Should error with too long password", func(t *testing.T) {
		long := make([]byte, vo.MaxPasswordLen+1)
		for i := range long {
			long[i] = 'a'
		}
		voPass, err := vo.NewPassword(string(long))
		assert.Equal(t, vo.Password{}, voPass)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})

	t.Run("Should validate hashed password", func(t *testing.T) {
		h := vo.NewHashedPassword([]byte("hash"))
		assert.NoError(t, h.Validate())
		assert.Equal(t, []byte("hash"), h.Value())
	})

	t.Run("Should error on empty hashed password", func(t *testing.T) {
		h := vo.NewHashedPassword(nil)
		assert.Error(t, h.Validate())
	})
}

