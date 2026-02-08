package log

import (
	"context"
	"io"
	"time"

	"github.com/OddKuru/accounts/pkg/logger"
	"github.com/rs/zerolog"
)

const (
	CallerSkipFrameCount       = 3
	StructCallerSkipFrameCount = 4
)

type (
	Options struct {
		Out    io.Writer
		Level  logger.Level
		Pretty bool
	}

	Log struct {
		withFields []logger.Field
		logger     zerolog.Logger
	}
)

func (l *Log) Enabled(ctx context.Context, level logger.Level) bool {
	lg := zerolog.Ctx(ctx)
	if lg == nil {
		lg = &l.logger
	}
	switch level {
	case logger.DebugLvl:
		lg.Debug().Enabled()
	case logger.InfoLvl:
		lg.Info().Enabled()
	case logger.WarnLvl:
		lg.Warn().Enabled()
	case logger.ErrorLvl:
		lg.Error().Enabled()
	}
	return false
}

//nolint:gocyclo,funlen
func fieldToEvent(event *zerolog.Event, fields []logger.Field) *zerolog.Event {
	for _, f := range fields {
		switch f.Type {
		case logger.StringType:
			val, ok := f.Value.(string)
			if ok {
				event = event.Str(f.Key, val)
			}
		case logger.IntType:
			val, ok := f.Value.(int)
			if ok {
				event = event.Int(f.Key, val)
			}
		case logger.ErrorType:
			val, ok := f.Value.(error)
			if ok {
				event = event.Err(val)
			}
		case logger.AnyType:
			event = event.Interface(f.Key, f.Value)
		case logger.DurationType:
			val, ok := f.Value.(time.Duration)
			if ok {
				event = event.Dur(f.Key, val)
			}
		case logger.TimeType:
			val, ok := f.Value.(time.Time)
			if ok {
				event = event.Time(f.Key, val)
			}
		case logger.BoolType:
			val, ok := f.Value.(bool)
			if ok {
				event = event.Bool(f.Key, val)
			}
		case logger.Int8Type:
			val, ok := f.Value.(int8)
			if ok {
				event = event.Int8(f.Key, val)
			}
		case logger.Int16Type:
			val, ok := f.Value.(int16)
			if ok {
				event = event.Int16(f.Key, val)
			}
		case logger.Int32Type:
			val, ok := f.Value.(int32)
			if ok {
				event = event.Int32(f.Key, val)
			}
		case logger.Int64Type:
			val, ok := f.Value.(int64)
			if ok {
				event = event.Int64(f.Key, val)
			}
		case logger.Uint8Type:
			val, ok := f.Value.(uint8)
			if ok {
				event = event.Uint8(f.Key, val)
			}
		case logger.Uint16Type:
			val, ok := f.Value.(uint16)
			if ok {
				event = event.Uint16(f.Key, val)
			}
		case logger.Uint32Type:
			val, ok := f.Value.(uint32)
			if ok {
				event = event.Uint32(f.Key, val)
			}
		case logger.Uint64Type:
			val, ok := f.Value.(uint64)
			if ok {
				event = event.Uint64(f.Key, val)
			}
		case logger.Float32Type:
			val, ok := f.Value.(float32)
			if ok {
				event = event.Float32(f.Key, val)
			}
		case logger.Float64Type:
			val, ok := f.Value.(float64)
			if ok {
				event = event.Float64(f.Key, val)
			}
		case logger.RawJsonType:
			val, ok := f.Value.([]byte)
			if ok {
				event = event.RawJSON(f.Key, val)
			}
		case logger.GroupType:
			val, ok := f.Value.([]logger.Field)
			if ok {
				d := zerolog.Dict()
				e := fieldToEvent(d, val)
				event = event.Dict(f.Key, e)
			}
		}
	}
	return event
}

func (l *Log) combineFields(fields []logger.Field) []logger.Field {
	combined := make([]logger.Field, 0, len(l.withFields)+len(fields))
	if len(l.withFields) > 0 {
		combined = append(combined, l.withFields...)
	}
	if len(fields) > 0 {
		combined = append(combined, fields...)
	}
	return combined
}

func (l *Log) Log(ctx context.Context, level logger.Level, message string, fields ...logger.Field) {
	lg := zerolog.Ctx(ctx)
	if lg == nil {
		lg = &l.logger
	}

	log := lg.With().Timestamp().CallerWithSkipFrameCount(StructCallerSkipFrameCount).Logger()
	var e *zerolog.Event
	switch level {
	case logger.DebugLvl:
		e = log.Debug().Ctx(ctx)
	case logger.InfoLvl:
		e = log.Info().Ctx(ctx)
	case logger.WarnLvl:
		e = log.Warn().Ctx(ctx)
	case logger.ErrorLvl:
		e = log.Error().Ctx(ctx)
	default:
		return
	}

	combineFields := l.combineFields(fields)

	e = fieldToEvent(e, combineFields)
	e.Msg(message)
}

func (l *Log) Debug(ctx context.Context, message string, fields ...logger.Field) {
	l.Log(ctx, logger.DebugLvl, message, fields...)
}

func (l *Log) Info(ctx context.Context, message string, fields ...logger.Field) {
	l.Log(ctx, logger.InfoLvl, message, fields...)
}

func (l *Log) Warn(ctx context.Context, message string, fields ...logger.Field) {
	l.Log(ctx, logger.WarnLvl, message, fields...)
}

func (l *Log) Error(ctx context.Context, message string, fields ...logger.Field) {
	l.Log(ctx, logger.ErrorLvl, message, fields...)
}

func (l *Log) With(fields ...logger.Field) logger.Logger {
	combined := l.combineFields(fields)
	return &Log{
		withFields: combined,
	}
}

func NewLogger(opt *Options) *Log {
	out := opt.Out
	if opt.Out == nil {
		out = DefaultOutput(opt.Pretty)
	}
	newLogger := zerolog.New(out).Level(convertLvlToZerologLvl(opt.Level))
	zerolog.DefaultContextLogger = &newLogger
	return &Log{
		withFields: make([]logger.Field, 0),
		logger:     newLogger,
	}
}

func convertLvlToZerologLvl(lvl logger.Level) zerolog.Level {
	switch lvl {
	case logger.DebugLvl:
		return zerolog.DebugLevel
	case logger.InfoLvl:
		return zerolog.InfoLevel
	case logger.WarnLvl:
		return zerolog.WarnLevel
	case logger.ErrorLvl:
		return zerolog.ErrorLevel
	}
	return zerolog.InfoLevel
}
