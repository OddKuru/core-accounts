package grpc_job

import (
	"testing"

	"github.com/OddKuru/core-accounts/pkg/ecosystem"
	"github.com/rs/zerolog"
)

func TestAppLoggerFromZerolog(t *testing.T) {
	output := &testOut{}
	logger := zerolog.New(output)

	log := ecosystem.NewAppLoggerWithZerolog(&logger)

	testAppLogger(t, log, output, "message")
}
