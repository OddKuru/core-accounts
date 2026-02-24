package v1

import (
	"context"

	v1 "github.com/OddKuru/core-accounts/gogen/accounts/v1"
	query "github.com/OddKuru/core-accounts/gogen/query/v1"
	"github.com/OddKuru/core-accounts/internal/app/usecase/account"
	"github.com/pkg/errors"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.AccountsServiceServer = (*Server)(nil)

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
	// TODO implement me
	panic("implement me")
}

func (s *Server) UpdateRoleById(ctx context.Context, request *v1.UpdateRoleByIdRequest) (*empty.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Server) UpdatePasswordById(ctx context.Context, request *v1.UpdatePasswordByIdRequest) (*empty.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Server) GetAccountById(ctx context.Context, id *v1.Id) (*v1.Account, error) {
	acc, err := s.accountUseCase.GetById(ctx, id.GetId())
	if err != nil {
		return nil, errors.Wrap(err, "[Server] accountUseCase.GetById")
	}
	return convertAccAggregateToDTO(acc), nil
}

func (s *Server) GetAccountsByQuery(ctx context.Context, data *query.QueryData) (*v1.QueryAccountsResponse, error) {
	// TODO implement me
	panic("implement me")
}
