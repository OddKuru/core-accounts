package vo_test

import (
	"testing"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestVoLoginName(t *testing.T) {
	t.Run("Should correct constructor login name", func(t *testing.T) {
		value := "user"
		name, err := vo.NewLoginName(value)
		assert.NoError(t, err)
		assert.Equal(t, value, name.Value())
	})

	t.Run("Should error with empty login name", func(t *testing.T) {
		name, err := vo.NewLoginName("")
		assert.Equal(t, vo.LoginName{}, name)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

