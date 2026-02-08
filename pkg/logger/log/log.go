package log

import (
	"context"

	"github.com/OddKuru/accounts/pkg/logger"
	"github.com/rs/zerolog"
)

func loggerFromContext(ctx context.Context) *zerolog.Logger {
	logFromCtx := zerolog.Ctx(ctx)
	if logFromCtx == nil {
		return nil
	}
	log := logFromCtx.With().Timestamp().CallerWithSkipFrameCount(CallerSkipFrameCount).Logger()
	return &log
}

func Debug(ctx context.Context, msg string, fields ...logger.Field) {
	log := loggerFromContext(ctx)
	if log == nil {
		return
	}
	e := log.Debug().Ctx(ctx)
	e = fieldToEvent(e, fields)
	e.Msg(msg)
}

func Info(ctx context.Context, msg string, fields ...logger.Field) {
	log := loggerFromContext(ctx)
	if log == nil {
		return
	}
	e := log.Info().Ctx(ctx)
	e = fieldToEvent(e, fields)
	e.Msg(msg)
}

func Warn(ctx context.Context, msg string, fields ...logger.Field) {
	log := loggerFromContext(ctx)
	if log == nil {
		return
	}
	e := log.Warn().Ctx(ctx)
	e = fieldToEvent(e, fields)
	e.Msg(msg)
}

func Error(ctx context.Context, msg string, fields ...logger.Field) {
	log := loggerFromContext(ctx)
	if log == nil {
		return
	}
	e := log.Error().Ctx(ctx)
	e = fieldToEvent(e, fields)
	e.Msg(msg)
}
