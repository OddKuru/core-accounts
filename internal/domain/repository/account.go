package repository

import (
	"context"

	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/rco"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
)

type AccountCommand interface {
	Create(ctx context.Context, account *aggregate.Account) error

	UpdateByEvents(ctx context.Context, account *aggregate.Account) error
}

type AccountQuery interface {
	HasByEmail(ctx context.Context, email vo.Email) (bool, error)
	HasByName(ctx context.Context, name vo.LoginName) (bool, error)

	GetById(ctx context.Context, id vo.ID) (*aggregate.Account, error)
	GetByQuery(ctx context.Context, query rco.Query) (rco.DataWithPageCount[*aggregate.Account], error)
}
