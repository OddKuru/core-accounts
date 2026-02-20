package v1

import (
	v1 "github.com/OddKuru/core-accounts/gogen/accounts/v1"
	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertAccAggregateToDTO(acc *aggregate.Account) *v1.Account {
	return &v1.Account{
		Id:        acc.Account().ID().Value(),
		Name:      acc.Account().Name().Value(),
		Email:     acc.Account().Email().Value(),
		Role:      vo.RoleString[acc.Account().Role().Value()],
		UpdatedAt: timestamppb.New(acc.Account().UpdatedAt()),
		CreatedAt: timestamppb.New(acc.Account().CreatedAt()),
	}
}
