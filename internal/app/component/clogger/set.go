package clogger

import (
	"context"
	"io"
	"os"
	"simple-reconciliation-service/internal/app/component/cconfig"
	"strconv"

	"github.com/google/wire"
)

func ProviderLogger(ctx context.Context, config *cconfig.Config) *Logger {
	var logOutWriter io.Writer
	var isShowLog bool
	isShowLog, _ = strconv.ParseBool(config.App.IsShowLog)
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
