package csqlite

import (
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/component/clogger"

	"github.com/google/wire"
)

type DBPath struct {
	ReadDBPath  string
	WriteDBPath string
}

func ProviderDBSqlite(config *cconfig.Config, logger *clogger.Logger, bBPath DBPath) (*DBSqlite, func(), error) {
	return NewDBSqlite(
		config,
		logger,
		bBPath.ReadDBPath,
		bBPath.WriteDBPath,
	)
}

var Set = wire.NewSet(
	ProviderDBSqlite,
)
