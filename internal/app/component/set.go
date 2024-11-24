package component

import (
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/component/cerror"
	"simple-reconciliation-service/internal/app/component/clogger"
	"simple-reconciliation-service/internal/app/component/csqlite"

	"github.com/google/wire"
)

func NewComponents(config *cconfig.Config, logger *clogger.Logger, er *cerror.Error, dbsqlite *csqlite.DBSqlite) *Components {
	return &Components{
		Config:   config,
		Logger:   logger,
		Error:    er,
		DBSqlite: dbsqlite,
	}
}

var Set = wire.NewSet(
	wire.Value(
		cconfig.ConfigPaths([]string{
			"./*.toml",
			"./params/*.toml",
		}),
	),
	cconfig.Set,
	clogger.Set,
	cerror.Set,
	csqlite.Set,
	NewComponents,
)
