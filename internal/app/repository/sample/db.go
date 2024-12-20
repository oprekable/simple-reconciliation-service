package sample

import (
	"context"
	"database/sql"
	"encoding/json"
	"simple-reconciliation-service/internal/app/err/core"
	"simple-reconciliation-service/internal/app/repository/_helper"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"strconv"
	"strings"
	"time"

	"github.com/aaronjan/hunch"
	"github.com/blockloop/scan/v2"
	"github.com/pkg/errors"
)

type DB struct {
	db *sql.DB
}

var _ Repository = (*DB)(nil)

func NewDB(
	db *sql.DB,
) (*DB, error) {
	return &DB{
		db: db,
	}, nil
}

func (d *DB) dropTables(ctx context.Context, tx *sql.Tx) (err error) {
	stmtData := []_helper.StmtData{
		{
			Query: QueryDropTableArguments,
		},
		{
			Query: QueryDropTableBanks,
		},
		{
			Query: QueryDropTableBaseData,
		},
	}

	return _helper.ExecTxQueries(ctx, d.db, tx, stmtData)
}

func (d *DB) createTables(ctx context.Context, tx *sql.Tx, listBank []string, startDate time.Time, toDate time.Time, limitTrxData int64, matchPercentage int) (err error) {
	stmtData := []_helper.StmtData{
		{
			Query: QueryCreateTableArguments,
			Args: func() []any {
				dateStringFormat := "2006-01-02"
				return []any{
					startDate.Format(dateStringFormat),
					toDate.Format(dateStringFormat),
					strconv.FormatInt(limitTrxData, 10),
					strconv.Itoa(matchPercentage),
				}
			}(),
		},
		{
			Query: QueryCreateTableBanks,
			Args: func() []any {
				b := new(strings.Builder)
				_ = json.NewEncoder(b).Encode(listBank)

				return []any{
					b.String(),
				}
			}(),
		},
		{
			Query: QueryCreateTableBaseData,
		},
		{
			Query: QueryCreateIndexTableBaseData,
		},
	}

	return _helper.ExecTxQueries(ctx, d.db, tx, stmtData)
}

func (d *DB) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time, limitTrxData int64, matchPercentage int) (err error) {
	var tx *sql.Tx
	tx, err = d.db.BeginTx(ctx, nil)

	defer func() {
		if err != nil {
			err = errors.Wrap(tx.Rollback(), err.Error())
		} else {
			err = tx.Commit()
		}

		log.Err(
			ctx,
			"[sample.NewDB] Exec Pre method in db",
			err,
		)
	}()

	if err != nil {
		return
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return nil, d.dropTables(c, tx)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			return nil, d.createTables(c, tx, listBank, startDate, toDate, limitTrxData, matchPercentage)
		},
	)

	return
}

func (d *DB) GetTrx(ctx context.Context) (returnData []TrxData, err error) {
	defer func() {
		log.Err(
			ctx,
			"[sample.NewDB] Exec GetData method in db",
			err,
		)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return d.db.PrepareContext(c, QueryGetTrxData)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return i.(*sql.Stmt).QueryContext(
				c,
			)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			rows := i.(*sql.Rows)
			return nil, scan.RowsStrict(&returnData, rows)
		},
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = core.CErrDBConn.Error()
	}

	return
}

func (d *DB) Post(ctx context.Context) (err error) {
	var tx *sql.Tx
	tx, err = d.db.BeginTx(ctx, nil)

	defer func() {
		if err != nil {
			err = errors.Wrap(tx.Rollback(), err.Error())
		} else {
			err = tx.Commit()
		}

		log.Err(
			ctx,
			"[sample.NewDB] Exec Post method in db",
			err,
		)
	}()

	if err != nil {
		return
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return nil, d.dropTables(c, tx)
		},
	)

	return
}

func (d *DB) Close() (err error) {
	return d.db.Close()
}
