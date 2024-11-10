package server

import (
	"simple-reconciliation-service/internal/app/server/cli"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	cli.Set,
	NewServer,
)
