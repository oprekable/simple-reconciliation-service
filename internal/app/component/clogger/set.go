package clogger

import "github.com/google/wire"

var Set = wire.NewSet(
	NewLogger,
)
