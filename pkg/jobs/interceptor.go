package jobs

import (
	"context"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var codeMap = map[codex.Code]codes.Code{
	codex.Internal:           codes.Internal,
	codex.Unknown:            codes.Unknown,
	codex.NotFound:           codes.NotFound,
	codex.InvalidArgument:    codes.InvalidArgument,
	codex.PermissionDenied:   codes.PermissionDenied,
	codex.DeadlineExceeded:   codes.DeadlineExceeded,
	codex.Canceled:           codes.Canceled,
	codex.Aborted:            codes.Aborted,
	codex.AlreadyExists:      codes.AlreadyExists,
	codex.ResourceExhausted:  codes.ResourceExhausted,
	codex.FailedPrecondition: codes.FailedPrecondition,
	codex.DataLoss:           codes.DataLoss,
	codex.OutOfRange:         codes.OutOfRange,
	codex.Unimplemented:      codes.Unimplemented,
	codex.Unavailable:        codes.Unavailable,
	codex.Unauthenticated:    codes.Unauthenticated,
	codex.OK:                 codes.OK,
}

func ErrorInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		grpcCode := codeMap[errx.Code(err)]
		return resp, status.Error(grpcCode, err.Error())
	}
	return resp, nil
}

func LogInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Debug(ctx, "request started",
			logger.String("method", info.FullMethod),
			logger.Any("req", req),
		)
		res, err := handler(ctx, req)
		log.Debug(ctx, "request finished",
			logger.String("method", info.FullMethod),
			logger.Any("res", res),
		)
		return res, err
	}
}
