package process

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"simple-reconciliation-service/internal/app/repository/_helper"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"strings"
	"time"

	"github.com/aaronjan/hunch"
	"github.com/blockloop/scan/v2"
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
			Query: QueryDropTableSystemTrx,
		},
		{
			Query: QueryDropTableBankTrx,
		},
		{
			Query: QueryDropTableReconciliationMap,
		},
	}

	return _helper.ExecTxQueries(ctx, d.db, tx, stmtData)
}

func (d *DB) dropTableWith(ctx context.Context, methodName string, extraExec hunch.ExecutableInSequence) (err error) {
	var tx *sql.Tx
	defer func() {
		err = _helper.CommitOrRollback(ctx, tx, err)
		log.Err(
			ctx,
			fmt.Sprintf(
				"[process.NewDB] Exec %s method in db",
				methodName,
			),
			err,
		)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			tx, e = d.db.BeginTx(ctx, nil)
			return nil, e
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			return tx, d.dropTables(c, tx)
		},
		extraExec,
	)

	return
}

func (d *DB) createTables(ctx context.Context, tx *sql.Tx, listBank []string, startDate time.Time, toDate time.Time) (err error) {
	return _helper.ExecTxQueries(
		ctx,
		d.db,
		tx,
		[]_helper.StmtData{
			{
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
				Query: QueryCreateTableSystemTrx,
			},
			{
				Query: QueryCreateTableBankTrx,
			},
			{
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
		err = _helper.CommitOrRollback(ctx, tx, err)
		log.Err(
			ctx,
			fmt.Sprintf("[process.NewDB] %s method to db (%d data)", methodName, reflect.ValueOf(data).Len()),
			err,
		)
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
					Query: query,
					Args: func() []any {
						return []any{
							string(i.([]byte)),
						}
					}(),
				},
			}

			return nil, _helper.ExecTxQueries(ctx, d.db, tx, stmtData)
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
		err = _helper.CommitOrRollback(ctx, tx, err)
		log.Err(
			ctx,
			fmt.Sprintf("[process.NewDB] Exec GenerateReconciliationMap method to db (Amount %f - %f)", minAmount, maxAmount),
			err,
		)
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
					Query: QueryInsertTableReconciliationMap,
					Args: func() []any {
						return []any{
							minAmount,
							maxAmount,
						}
					}(),
				},
			}

			return nil, _helper.ExecTxQueries(ctx, d.db, tx, stmtData)
		},
	)

	return
}

func (d *DB) GetReconciliationSummary(ctx context.Context) (returnData ReconciliationSummary, err error) {
	defer func() {
		log.Err(
			ctx,
			"[process.NewDB] Exec GetReconciliationSummary method from db",
			err,
		)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return d.db.PrepareContext(c, QueryGetReconciliationSummary)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return i.(*sql.Stmt).QueryContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			rows := i.(*sql.Rows)
			return nil, scan.RowStrict(&returnData, rows)
		},
	)

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
