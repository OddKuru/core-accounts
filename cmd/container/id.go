package container

import (
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/google/uuid"
)

type UUIDV4 struct{}

func (U UUIDV4) GenerateID() (vo.ID, error) {
	newID := uuid.New().String()
	voID, err := vo.NewID(newID)
	if err != nil {
		return vo.ID{}, err
	}
	return voID, nil
}
