package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type OutAndErrOutput struct {
	output    io.Writer
	errOutput io.Writer
}

func (l OutAndErrOutput) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (l OutAndErrOutput) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level > zerolog.InfoLevel {
		return l.errOutput.Write(p)
	}
	return l.output.Write(p)
}

func DefaultOutput(isPretty bool) OutAndErrOutput {
	if isPretty {
		return OutAndErrOutput{
			output:    zerolog.ConsoleWriter{Out: os.Stdout},
			errOutput: zerolog.ConsoleWriter{Out: os.Stderr},
		}
	}
	return OutAndErrOutput{
		output:    os.Stdout,
		errOutput: os.Stderr,
	}
}
