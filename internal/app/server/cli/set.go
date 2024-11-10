package cli

import (
	"simple-reconciliation-service/internal/app/handler/hcli"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	hcli.Set,
	NewCli,
)
