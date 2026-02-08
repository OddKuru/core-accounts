package service

import "github.com/OddKuru/core-accounts/internal/domain/vo"

type (
	PasswordHasher interface {
		Hash(password vo.Password) (vo.HashedPassword, error)
		Compare(hash vo.HashedPassword, password vo.Password) error
	}
)
