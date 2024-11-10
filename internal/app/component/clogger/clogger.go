package clogger

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	log zerolog.Logger
	ctx context.Context
}

func NewLogger(ctx context.Context) *Logger {
	re := regexp.MustCompile(`\r?\n`)
	var writers []io.Writer
	stdOut := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
	stdOut.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}

	stdOut.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("| %s:", i)
	}

	stdOut.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	stdOut.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("*** %s ****", re.ReplaceAllString(i.(string), " "))
	}

	writers = append(writers, stdOut)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	mw := io.MultiWriter(writers...)

	loggerCtx := zerolog.New(mw).
		With().
		Timestamp().
		Stack()

	logger := loggerCtx.Logger()
	ctx = context.WithValue(ctx, log.StartTime, time.Now())
	logger = logger.Hook(log.UptimeHook{})
	ctx = logger.WithContext(ctx)

	return &Logger{
		log: logger,
		ctx: ctx,
	}
}

func (l *Logger) GetLogger() zerolog.Logger {
	return l.log
}

func (l *Logger) GetCtx() context.Context {
	return l.ctx
}
