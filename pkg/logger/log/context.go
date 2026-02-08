package log

import (
	"context"
)

func (l *Log) InjectCtx(ctx context.Context) context.Context {
	return l.logger.WithContext(ctx)
}
