package vo_test

import (
	"testing"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestVoEmail(t *testing.T) {
	t.Run("Should correct constructor email", func(t *testing.T) {
		email := "test@example.com"
		voEmail, err := vo.NewEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, email, voEmail.Value())
	})

	t.Run("Should error with empty email", func(t *testing.T) {
		voEmail, err := vo.NewEmail("")
		assert.Equal(t, vo.Email{}, voEmail)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})

	t.Run("Should error with incorrect email", func(t *testing.T) {
		voEmail, err := vo.NewEmail("not-email")
		assert.Equal(t, vo.Email{}, voEmail)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

