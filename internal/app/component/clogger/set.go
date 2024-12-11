package clogger

import (
	"context"
	"io"
	"os"

	"github.com/google/wire"
)

type IsShowLog bool

func ProviderLogger(ctx context.Context, isShowLog IsShowLog) *Logger {
	var logOutWriter io.Writer
	if isShowLog {
		logOutWriter = os.Stdout
	} else {
		logOutWriter = io.Discard
	}

	return NewLogger(
		ctx,
		logOutWriter,
	)
}

var Set = wire.NewSet(
	ProviderLogger,
)
