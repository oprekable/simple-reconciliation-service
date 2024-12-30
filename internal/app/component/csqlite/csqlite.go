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

func NewDBSqlite(config *cconfig.Config, logger *clogger.Logger, readDBPath string, writeDBPath string) (rd *DBSqlite, cleanFunc func(), err error) {
	rd = &DBSqlite{}
	ctx := logger.GetLogger().With().Str("component", "NewDBSqlite").Ctx(context.Background()).Logger().WithContext(logger.GetCtx())

	defer func() {
		if r := recover(); r != nil {
			errRecovery := fmt.Errorf("recovered from panic: %s", r)
			log.AddErr(ctx, errRecovery)
			return
		}
	}()

	if !config.Data.Sqlite.IsEnabled {
		return
	}

	if config.Data.Sqlite.Write.IsEnabled {
		rd.dBWriteConnOnce.Do(func() {
			dbParameters := config.Data.Sqlite.Write
			if writeDBPath != "" {
				dbParameters.DBPath = writeDBPath
			}

			rd.DBWrite, err = sqlDriver.NewSqliteDatabase(
				dbParameters.Options("sqlite_write"),
				logger.GetLogger(),
				config.Data.Sqlite.IsDoLogging,
			)

			log.AddErr(ctx, err)
		})
	}

	if config.Data.Sqlite.Read.IsEnabled {
		rd.dBReadConnOnce.Do(func() {
			dbParameters := config.Data.Sqlite.Read
			if readDBPath != "" {
				dbParameters.DBPath = readDBPath
			}

			rd.DBRead, err = sqlDriver.NewSqliteDatabase(
				dbParameters.Options("sqlite_read"),
				logger.GetLogger(),
				config.Data.Sqlite.IsDoLogging,
			)

			log.AddErr(ctx, err)
		})
	}

	cleanFunc = func() {
		_ = rd.DBRead.Close()
		_ = rd.DBWrite.Close()
	}

	log.Msg(ctx, "sqlite connection loaded")
	return rd, cleanFunc, err
}
