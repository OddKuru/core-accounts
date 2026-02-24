package jobs

import (
	"context"
	"sync"
	"time"

	ayaka "github.com/OddKuru/core-accounts/pkg/core"
	"github.com/OddKuru/core-accounts/pkg/ecosystem"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var mu = sync.Mutex{}

func registerPrometheus(srv *grpc.Server) {
	mu.Lock()
	defer mu.Unlock()
	grpcPrometheus.EnableHandlingTimeHistogram(
		grpcPrometheus.WithHistogramBuckets(
			prometheus.DefBuckets,
		),
	)

	grpcPrometheus.Register(srv)
}

// Health service

type healthService struct{}

func (s *healthService) List(
	_ context.Context,
	_ *health.HealthListRequest,
) (*health.HealthListResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s *healthService) Check(context.Context, *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}

func (s *healthService) Watch(*health.HealthCheckRequest, health.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

type GRPCVar struct {
	Address        string
	RequestTimeout time.Duration
	Tracer         opentracing.Tracer
}

func NewGRPC[T any](v GRPCVar, regs ...ecosystem.GrpcRegister[T]) (ayaka.Job[T], error) {
	job, err := ecosystem.NewGrpcJobBuilder[T]().
		Address(v.Address).
		Interceptors(
			ErrorInterceptor,
		).
		RequestTimeout(v.RequestTimeout).
		Register(regs...).
		RegisterServer(func(srv *grpc.Server) error {
			// Register monitoring
			registerPrometheus(srv)
			// Register healthcheck service
			health.RegisterHealthServer(srv, new(healthService))
			// Register reflection service on gRPC server.
			reflection.Register(srv)
			return nil
		}).
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "ecosystem.NewGrpcJobBuilder.Build")
	}

	return job, nil
}
