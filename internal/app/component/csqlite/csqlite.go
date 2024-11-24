package csqlite

import (
	"context"
	"database/sql"
	"fmt"
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/component/clogger"
	sqlDriver "simple-reconciliation-service/internal/pkg/driver/sql"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"sync"
)

type DBSqlite struct {
	DBWrite         *sql.DB
	DBRead          *sql.DB
	dBWriteConnOnce sync.Once
	dBReadConnOnce  sync.Once
}

func NewDBSqlite(config *cconfig.Config, logger *clogger.Logger) (rd *DBSqlite, err error) {
	rd = &DBSqlite{}
	ctx := logger.GetLogger().With().Str("component", "NewDBSqlite").Ctx(context.Background()).Logger().WithContext(logger.GetCtx())

	defer func() {
		if r := recover(); r != nil {
			errRecovery := fmt.Errorf("recovered from panic: %s", r)
			log.AddErr(ctx, errRecovery)
			return
		}
	}()

	if config.Data.Sqlite.IsEnabled {
		if config.Data.Sqlite.Write.IsEnabled {
			rd.dBWriteConnOnce.Do(func() {
				if rd.DBWrite, err = sqlDriver.NewSqliteDatabase(
					config.Data.Sqlite.Write.Options("sqlite_write"),
					logger.GetLogger(),
					config.Data.Sqlite.IsDoLogging,
				); err != nil {
					log.AddErr(ctx, err)
				}
			})
		}

		if config.Data.Sqlite.Read.IsEnabled {
			rd.dBReadConnOnce.Do(func() {
				if rd.DBRead, err = sqlDriver.NewSqliteDatabase(
					config.Data.Sqlite.Read.Options("sqlite_read"),
					logger.GetLogger(),
					config.Data.Sqlite.IsDoLogging,
				); err != nil {
					log.AddErr(ctx, err)
				}
			})
		}
	}

	log.Msg(ctx, "sqlite connection loaded")
	return rd, err
}
