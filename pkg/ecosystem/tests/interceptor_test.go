package grpc_job

import (
	"context"
	"testing"
	"time"

	"github.com/OddKuru/core-accounts/pkg/ecosystem"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTimeoutInterceptor(t *testing.T) {
	t.Run("Should correct working", func(t *testing.T) {
		interceptor := ecosystem.TimeoutInterceptor(time.Second)
		ctx := context.Background()
		var req interface{}
		res := 42

		response, err := interceptor(ctx, req, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return res, nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 42, response)
	})

	t.Run("Should timout error", func(t *testing.T) {
		t.Parallel()
		interceptor := ecosystem.TimeoutInterceptor(time.Second)
		ctx := context.Background()
		var req interface{}
		res := 42

		response, err := interceptor(ctx, req, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
			ti := time.NewTimer(time.Second + (time.Millisecond * 200))
			<-ti.C
			return res, nil
		})

		assert.Nil(t, response)
		assert.Error(t, err)
		stat, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.DeadlineExceeded, stat.Code())
	})
}
