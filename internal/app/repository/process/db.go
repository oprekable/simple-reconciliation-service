package process

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"simple-reconciliation-service/internal/app/err/core"
	"simple-reconciliation-service/internal/app/repository/_helper"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/aaronjan/hunch"
	"github.com/blockloop/scan/v2"
	"github.com/goccy/go-json"
)

type DB struct {
	db      *sql.DB
	stmtMap map[string]*sql.Stmt
}

var _ Repository = (*DB)(nil)

func NewDB(
	db *sql.DB,
) (*DB, error) {
	return &DB{
		db:      db,
		stmtMap: make(map[string]*sql.Stmt),
	}, nil
}

func (d *DB) dropTableWith(ctx context.Context, methodName string, extraExec hunch.ExecutableInSequence) (err error) {
	var tx *sql.Tx
	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			tx, e = d.db.BeginTx(ctx, nil)
			return nil, e
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			stmtData := []_helper.StmtData{
				{
					Name:  "QueryDropTableArguments",
					Query: QueryDropTableArguments,
				},
				{
					Name:  "QueryDropTableBanks",
					Query: QueryDropTableBanks,
				},
				{
					Name:  "QueryDropTableSystemTrx",
					Query: QueryDropTableSystemTrx,
				},
				{
					Name:  "QueryDropTableBankTrx",
					Query: QueryDropTableBankTrx,
				},
				{
					Name:  "QueryDropTableReconciliationMap",
					Query: QueryDropTableReconciliationMap,
				},
			}

			return tx, _helper.ExecTxQueries(ctx, d.db, tx, d.stmtMap, stmtData)
		},
		extraExec,
	)

	defer func() {
		log.Err(ctx, fmt.Sprintf("[process.NewDB] Exec %s method in db", methodName), _helper.CommitOrRollback(tx, err))
	}()

	return
}

func (d *DB) createTables(ctx context.Context, tx *sql.Tx, listBank []string, startDate time.Time, toDate time.Time) (err error) {
	return _helper.ExecTxQueries(
		ctx,
		d.db,
		tx,
		d.stmtMap,
		[]_helper.StmtData{
			{
				Name:  "QueryCreateTableArguments",
				Query: QueryCreateTableArguments,
				Args: func() []any {
					dateStringFormat := "2006-01-02"
					return []any{
						startDate.Format(dateStringFormat),
						toDate.Format(dateStringFormat),
					}
				}(),
			},
			{
				Name:  "QueryCreateTableBanks",
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
				Name:  "QueryCreateTableSystemTrx",
				Query: QueryCreateTableSystemTrx,
			},
			{
				Name:  "QueryCreateTableBankTrx",
				Query: QueryCreateTableBankTrx,
			},
			{
				Name:  "QueryCreateTableReconciliationMap",
				Query: QueryCreateTableReconciliationMap,
			},
		},
	)
}

func (d *DB) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time) (err error) {
	extraExec := func(c context.Context, i interface{}) (interface{}, error) {
		return nil, d.createTables(c, i.(*sql.Tx), listBank, startDate, toDate)
	}

	return d.dropTableWith(
		ctx,
		"Pre",
		extraExec,
	)
}

func (d *DB) importInterface(ctx context.Context, methodName string, query string, data interface{}) (err error) {
	var tx *sql.Tx
	defer func() {
		log.Err(ctx, fmt.Sprintf("[process.NewDB] %s method to db (%d data)", methodName, reflect.ValueOf(data).Len()), _helper.CommitOrRollback(tx, err))
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			tx, e = d.db.BeginTx(ctx, nil)
			return nil, e
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			return json.Marshal(data)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			stmtData := []_helper.StmtData{
				{
					Name:  methodName,
					Query: query,
					Args: func() []any {
						return []any{
							string(i.([]byte)),
						}
					}(),
				},
			}

			return nil, _helper.ExecTxQueries(ctx, d.db, tx, d.stmtMap, stmtData)
		},
	)

	return
}

func (d *DB) ImportSystemTrx(ctx context.Context, data []*systems.SystemTrxData) (err error) {
	return d.importInterface(ctx, "ImportSystemTrx", QueryInsertTableSystemTrx, data)
}

func (d *DB) ImportBankTrx(ctx context.Context, data []*banks.BankTrxData) (err error) {
	return d.importInterface(ctx, "ImportBankTrx", QueryInsertTableBankTrx, data)
}

func (d *DB) GenerateReconciliationMap(ctx context.Context, minAmount float64, maxAmount float64) (err error) {
	var tx *sql.Tx
	defer func() {
		log.Err(ctx, fmt.Sprintf("[process.NewDB] Exec GenerateReconciliationMap method to db (Amount %f - %f)", minAmount, maxAmount), _helper.CommitOrRollback(tx, err))
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			tx, e = d.db.BeginTx(ctx, nil)
			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			stmtData := []_helper.StmtData{
				{
					Name:  "QueryInsertTableReconciliationMap",
					Query: QueryInsertTableReconciliationMap,
					Args: func() []any {
						return []any{
							minAmount,
							maxAmount,
						}
					}(),
				},
			}

			return nil, _helper.ExecTxQueries(ctx, d.db, tx, d.stmtMap, stmtData)
		},
	)

	return
}

func (d *DB) GetReconciliationSummary(ctx context.Context) (returnData ReconciliationSummary, err error) {
	stmtName := "QueryGetReconciliationSummary"
	stmtQuery := QueryGetReconciliationSummary

	defer func() {
		log.Err(ctx, "[process.NewDB] Exec GetReconciliationSummary method from db", err)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (_ interface{}, e error) {
			_, ok := d.stmtMap[stmtName]
			if !ok {
				d.stmtMap[stmtName], e = d.db.PrepareContext(c, stmtQuery) //nolint:sqlclosecheck
			}

			return
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return d.stmtMap[stmtName].QueryContext(
				c,
			)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			rows := i.(*sql.Rows)
			return nil, scan.RowStrict(&returnData, rows)
		},
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = core.CErrDBConn.Error()
	}

	return
}

func (d *DB) Post(ctx context.Context) (err error) {
	extraExec := func(c context.Context, i interface{}) (interface{}, error) {
		return nil, nil
	}

	return d.dropTableWith(
		ctx,
		"Post",
		extraExec,
	)
}

func (d *DB) Close() (err error) {
	return d.db.Close()
}

func (d *DB) GetMatchedTrx(ctx context.Context) (returnData []MatchedTrx, err error) {
	stmtName := "QueryGetMatchedTrx"
	stmtQuery := QueryGetMatchedTrx

	defer func() {
		log.Err(ctx, "[process.NewDB] Exec GetMatchedTrx method from db", err)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (_ interface{}, e error) {
			_, ok := d.stmtMap[stmtName]
			if !ok {
				d.stmtMap[stmtName], e = d.db.PrepareContext(c, stmtQuery) //nolint:sqlclosecheck
			}

			return
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return d.stmtMap[stmtName].QueryContext(
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

func (d *DB) GetNotMatchedSystemTrx(ctx context.Context) (returnData []NotMatchedSystemTrx, err error) {
	stmtName := "QueryGetNotMatchedSystemTrx"
	stmtQuery := QueryGetNotMatchedSystemTrx

	defer func() {
		log.Err(ctx, "[process.NewDB] Exec GetNotMatchedSystemTrx method from db", err)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (_ interface{}, e error) {
			_, ok := d.stmtMap[stmtName]
			if !ok {
				d.stmtMap[stmtName], e = d.db.PrepareContext(c, stmtQuery) //nolint:sqlclosecheck
			}

			return
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return d.stmtMap[stmtName].QueryContext(
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

func (d *DB) GetNotMatchedBankTrx(ctx context.Context) (returnData []NotMatchedBankTrx, err error) {
	stmtName := "QueryGetNotMatchedBankTrx"
	stmtQuery := QueryGetNotMatchedBankTrx

	defer func() {
		log.Err(ctx, "[process.NewDB] Exec GetNotMatchedBankTrx method from db", err)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (_ interface{}, e error) {
			_, ok := d.stmtMap[stmtName]
			if !ok {
				d.stmtMap[stmtName], e = d.db.PrepareContext(c, stmtQuery) //nolint:sqlclosecheck
			}

			return
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return d.stmtMap[stmtName].QueryContext(
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
