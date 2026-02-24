package password

import (
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct{}

func (s *Service) Hash(password vo.Password) (vo.HashedPassword, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.Value()), bcrypt.DefaultCost)
	if err != nil {
		return vo.HashedPassword{}, errors.Wrap(err, "[Service] bcrypt.GenerateFromPassword")
	}
	result := vo.NewHashedPassword(hashedPassword)
	return result, nil
}

func (s *Service) Compare(hash vo.HashedPassword, password vo.Password) error {
	err := bcrypt.CompareHashAndPassword(hash.Value(), []byte(password.Value()))
	if err != nil {
		return errors.Wrap(err, "[Service] bcrypt.CompareHashAndPassword")
	}
	return nil
}

func New() (*Service, error) {
	return &Service{}, nil
}
