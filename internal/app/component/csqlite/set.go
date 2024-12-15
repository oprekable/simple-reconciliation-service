package csqlite

import (
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/component/clogger"

	"github.com/google/wire"
)

type ReadDBPath string
type WriteDBPath string

func ProviderDBSqlite(config *cconfig.Config, logger *clogger.Logger, readDBPath ReadDBPath, writeDBPath WriteDBPath) (*DBSqlite, error) {
	return NewDBSqlite(
		config,
		logger,
		string(readDBPath),
		string(writeDBPath),
	)
}

var Set = wire.NewSet(
	ProviderDBSqlite,
)
