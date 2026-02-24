package jobs

import (
	"os"
	"strconv"
	"time"

	ayaka "github.com/OddKuru/core-accounts/pkg/core"
	"github.com/OddKuru/core-accounts/pkg/ecosystem"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type GrpcJobEnvKeys struct {
	Address,
	RequestTimeout string
}

var DefaultUnaryJobEnvKeys = &GrpcJobEnvKeys{
	Address:        "GRPC_ADDRESS",
	RequestTimeout: "GRPC_REQUEST_TIMEOUT",
}

func NewGRPCWithEnv[T any](
	tracer opentracing.Tracer,
	keys *GrpcJobEnvKeys,
	regs ...ecosystem.GrpcRegister[T],
) (ayaka.Job[T], error) {
	if keys == nil {
		keys = DefaultUnaryJobEnvKeys
	}

	address := os.Getenv(keys.Address)
	if address == "" {
		return nil, errors.Errorf("environment variable %s is not set", keys.Address)
	}
	requestTimeoutEnv := os.Getenv(keys.RequestTimeout)
	if requestTimeoutEnv == "" {
		return nil, errors.Errorf("environment variable %s is not set", keys.RequestTimeout)
	}
	requestTimeout, err := strconv.Atoi(requestTimeoutEnv)
	if err != nil {
		return nil, errors.Errorf("environment variable %s is not a number", keys.RequestTimeout)
	}

	job, err := NewGRPC[T](GRPCVar{
		Address:        address,
		RequestTimeout: time.Duration(requestTimeout) * time.Second,
		Tracer:         tracer,
	}, regs...)
	if err != nil {
		return nil, err
	}

	return job, nil
}
