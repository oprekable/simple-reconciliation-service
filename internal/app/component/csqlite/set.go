package csqlite

import "github.com/google/wire"

var Set = wire.NewSet(
	NewDBSqlite,
)
