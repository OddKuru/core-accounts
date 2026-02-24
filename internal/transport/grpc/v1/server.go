package v1

import (
	"context"

	v1 "github.com/OddKuru/core-accounts/gogen/accounts/v1"
	query "github.com/OddKuru/core-accounts/gogen/query/v1"
	"github.com/OddKuru/core-accounts/internal/app/usecase/account"
	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/rco"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/OddKuru/core-accounts/pkg/fns"
	"github.com/pkg/errors"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.AccountsServiceServer = (*Server)(nil)

var RoleToType = map[string]vo.AccountRoleType{
	"user":  vo.RoleUser,
	"admin": vo.RoleAdmin,
}

type Server struct {
	accountUseCase account.Account
}

func NewGRPCServer(
	accountUseCase account.Account,
) (*Server, error) {
	res := &Server{
		accountUseCase: accountUseCase,
	}

	if err := res.validate(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Server) validate() error {
	if s.accountUseCase == nil {
		return errors.New("accountUseCase is empty")
	}
	return nil
}

func (s *Server) Create(ctx context.Context, request *v1.CreateAccountRequest) (*v1.Account, error) {
	acc, err := s.accountUseCase.Create(ctx, request.GetName(), request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, errors.Wrap(err, "[Server] accountUseCase.Create")
	}

	return convertAccAggregateToDTO(acc), nil
}

func (s *Server) UpdateNameById(ctx context.Context, request *v1.UpdateNameByIdRequest) (*empty.Empty, error) {
	err := s.accountUseCase.UpdateNameById(ctx, request.GetAccountId(), request.GetName())
	if err != nil {
		return nil, errors.Wrap(err, "[Server] accountUseCase.UpdateNameById")
	}
	return &empty.Empty{}, nil
}

func (s *Server) UpdateEmailById(ctx context.Context, request *v1.UpdateEmailByIdRequest) (*empty.Empty, error) {
	err := s.accountUseCase.UpdateEmailById(ctx, request.GetAccountId(), request.GetEmail())
	if err != nil {
		return nil, errors.Wrap(err, "[Server] accountUseCase.UpdateEmailById")
	}
	return &empty.Empty{}, nil
}

func (s *Server) UpdateRoleById(ctx context.Context, request *v1.UpdateRoleByIdRequest) (*empty.Empty, error) {
	err := s.accountUseCase.UpdateRoleById(ctx, request.GetAccountId(), RoleToType[request.GetRole()])
	if err != nil {
		return nil, errors.Wrap(err, "[Server] accountUseCase.UpdateRoleById")
	}
	return &empty.Empty{}, nil
}

func (s *Server) UpdatePasswordById(ctx context.Context, request *v1.UpdatePasswordByIdRequest) (*empty.Empty, error) {
	err := s.accountUseCase.UpdatePasswordById(
		ctx, request.GetAccountId(), request.GetOldPassword(), request.GetNewPassword(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "[Server] accountUseCase.UpdatePasswordById")
	}
	return &empty.Empty{}, nil
}

func (s *Server) GetAccountById(ctx context.Context, id *v1.Id) (*v1.Account, error) {
	acc, err := s.accountUseCase.GetById(ctx, id.GetId())
	if err != nil {
		return nil, errors.Wrap(err, "[Server] accountUseCase.GetById")
	}
	return convertAccAggregateToDTO(acc), nil
}

func (s *Server) GetAccountsByQuery(ctx context.Context, data *query.QueryData) (*v1.QueryAccountsResponse, error) {
	q, err := rco.NewQuery(
		uint(data.GetPage()),
		uint(data.GetLimit()),
		data.GetSortBy(),
		rco.QueryOrder(data.GetOrderBy().String()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "[Server] rco.NewQuery")
	}
	list, err := s.accountUseCase.GetByQuery(ctx, q)
	if err != nil {
		return nil, errors.Wrap(err, "[Server] accountUseCase.GetByQuery")
	}

	return &v1.QueryAccountsResponse{
		Accounts: fns.Map(list.Data(), func(acc *aggregate.Account) *v1.Account {
			return convertAccAggregateToDTO(acc)
		}),
		PageCount: uint64(list.PageCount()),
	}, nil
}
