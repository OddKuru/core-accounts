package ports

import "github.com/OddKuru/core-accounts/internal/domain/vo"

type IDGenerator interface {
	GenerateID() (vo.ID, error)
}
