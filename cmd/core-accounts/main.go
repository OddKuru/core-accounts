package main

import (
	"context"

	"github.com/OddKuru/core-accounts/cmd/container"
	v1 "github.com/OddKuru/core-accounts/gogen/accounts/v1"
	grpcServer "github.com/OddKuru/core-accounts/internal/transport/grpc/v1"
	ayaka "github.com/OddKuru/core-accounts/pkg/core"
	"github.com/OddKuru/core-accounts/pkg/jobs"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func main() {
	cont, err := container.AppContainer()
	if err != nil {
		panic(err)
	}

	app := ayaka.NewApp[*container.Dependency](&ayaka.Options[*container.Dependency]{
		Name:        "accounts",
		Description: "accounts core service",
		Version:     "0.0.0",
		Logger:      cont.AppLogger(),
		Container:   cont,
	}).WithConfig(&ayaka.Config{
		StartTimeout:    cont.Config().App.StartTimeout,
		GracefulTimeout: cont.Config().App.GracefulTimeout,
	})

	accGRPCJob, err := jobs.NewGRPC[*container.Dependency](
		jobs.GRPCVar{
			Address:        cont.Config().App.GRPCAddress,
			RequestTimeout: cont.Config().App.GRPCTimeout,
			Tracer:         opentracing.NoopTracer{},
		},
		grpcJob,
	)

	if err != nil {
		panic(err)
	}

	app.WithJob(
		ayaka.JobEntry[*container.Dependency]{
			Key: "accounts_grpc",
			Job: accGRPCJob,
		},
	)

	err = app.Start()
	if err != nil {
		panic(err)
	}
}

func grpcJob(_ context.Context, di *container.Dependency, srv *grpc.Server) error {
	server, err := grpcServer.NewGRPCServer(di.AccountsUseCase())
	if err != nil {
		return err
	}
	v1.RegisterAccountsServiceServer(srv, server)

	return nil
}
