package vo_test

import (
	"testing"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestVoID(t *testing.T) {
	t.Run("Should correct constructor id", func(t *testing.T) {
		id := uuid.New().String()
		voID, err := vo.NewID(id)
		assert.NoError(t, err)
		assert.Equal(t, id, voID.Value())
		assert.Equal(t, id, voID.String())
	})

	t.Run("Should error with incorrect id", func(t *testing.T) {
		id := "incorrect-uuidv4-value"
		voID, err := vo.NewID(id)
		assert.Equal(t, vo.ID{}, voID)
		assert.Error(t, err)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

