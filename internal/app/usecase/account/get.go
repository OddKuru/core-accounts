package account

import (
	"context"

	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/rco"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/pkg/errors"
)

func (u *UseCase) GetById(ctx context.Context, id string) (*aggregate.Account, error) {
	voID, err := vo.NewID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] vo.NewID")
	}

	acc, err := u.accountQuery.GetById(ctx, voID)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] accountQuery.GetById")
	}
	return acc, nil
}

func (u *UseCase) GetByQuery(ctx context.Context, query rco.Query) (*rco.DataWithPageCount[*aggregate.Account], error) {
	res, err := u.accountQuery.GetByQuery(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "[UseCase] accountQuery.GetByQuery")
	}
	return res, nil
}
