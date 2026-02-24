package logger

import (
	"context"
)

const (
	_ = iota
	StringType
	IntType
	AnyType
	ErrorType
	DurationType
	TimeType
	BoolType
	Float32Type
	Float64Type
	Int8Type
	Int16Type
	Int32Type
	Int64Type
	Uint8Type
	Uint16Type
	Uint32Type
	Uint64Type
	RawJsonType
	GroupType
)

type Level int

const (
	DebugLvl Level = -10
	InfoLvl  Level = 0
	WarnLvl  Level = 10
	ErrorLvl Level = 20
)

const (
	ErrKey = "error"
)

type (
	NoopLogger struct{}

	Logger interface {
		Log(ctx context.Context, level Level, message string, fields ...Field)
		Debug(ctx context.Context, message string, fields ...Field)
		Info(ctx context.Context, message string, fields ...Field)
		Warn(ctx context.Context, message string, fields ...Field)
		Error(ctx context.Context, message string, fields ...Field)
		With(fields ...Field) Logger
		InjectCtx(ctx context.Context) context.Context
		Enabled(ctx context.Context, level Level) bool
	}

	Field struct {
		Key   string
		Type  int
		Value any
	}
)

func (n NoopLogger) Log(_ context.Context, _ Level, _ string, _ ...Field) {}

func (n NoopLogger) Debug(_ context.Context, _ string, _ ...Field) {}

func (n NoopLogger) Info(_ context.Context, _ string, _ ...Field) {}

func (n NoopLogger) Warn(_ context.Context, _ string, _ ...Field) {}

func (n NoopLogger) Error(_ context.Context, _ string, _ ...Field) {}

func (n NoopLogger) With(_ ...Field) Logger {
	return n
}

func (n NoopLogger) InjectCtx(_ context.Context) context.Context {
	return context.Background()
}

func (n NoopLogger) Enabled(_ context.Context, _ Level) bool {
	return false
}

var _ Logger = NoopLogger{}
