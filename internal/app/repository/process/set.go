package process

import (
	"simple-reconciliation-service/internal/app/component"

	"github.com/google/wire"
)

func ProviderDB(comp *component.Components) (*DB, error) {
	return NewDB(comp.DBSqlite.DBWrite)
}

var Set = wire.NewSet(
	ProviderDB,
	wire.Bind(new(Repository), new(*DB)),
)
