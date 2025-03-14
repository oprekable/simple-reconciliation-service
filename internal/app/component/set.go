package component

import (
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/component/cerror"
	"simple-reconciliation-service/internal/app/component/cfs"
	"simple-reconciliation-service/internal/app/component/clogger"
	"simple-reconciliation-service/internal/app/component/csqlite"

	"github.com/spf13/afero"

	"github.com/google/wire"
)

func NewComponents(config *cconfig.Config, logger *clogger.Logger, er *cerror.Error, dbsqlite *csqlite.DBSqlite, fs *cfs.Fs) *Components {
	return &Components{
		Config:   config,
		Logger:   logger,
		Error:    er,
		DBSqlite: dbsqlite,
		Fs:       fs,
	}
}

var Set = wire.NewSet(
	wire.Value(
		cconfig.ConfigPaths([]string{
			"./*.toml",
			"./params/*.toml",
		}),
	),
	wire.InterfaceValue(new(afero.Fs), afero.NewOsFs()),
	cconfig.Set,
	clogger.Set,
	cerror.Set,
	csqlite.Set,
	cfs.Set,
	NewComponents,
)
