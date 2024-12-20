package process

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"simple-reconciliation-service/internal/app/repository/_helper"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems"
	"simple-reconciliation-service/internal/pkg/utils/log"
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

func (d *DB) createTables(ctx context.Context, tx *sql.Tx, listBank []string, startDate time.Time, toDate time.Time) (err error) {
	stmtData := []_helper.StmtData{
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
	}

	return _helper.ExecTxQueries(ctx, d.db, tx, stmtData)
}

func (d *DB) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time) (err error) {
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
			"[process.NewDB] Exec Pre method in db",
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
			return nil, d.createTables(c, tx, listBank, startDate, toDate)
		},
	)

	return
}

func (d *DB) ImportSystemTrx(ctx context.Context, data []*systems.SystemTrxData) (err error) {
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
			fmt.Sprintf("[process.NewDB] ImportSystemTrx method to db (%d data)", len(data)),
			err,
		)
	}()

	if err != nil {
		return
	}

	var jsonData string
	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return json.Marshal(data)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				jsonData = string(i.([]byte))
			}

			return nil, nil
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return d.db.PrepareContext(c, QueryInsertTableSystemTrx)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, i.(*sql.Stmt)).ExecContext( //nolint:sqlclosecheck
				c,
				jsonData,
			)
		},
	)

	return
}

func (d *DB) ImportBankTrx(ctx context.Context, data []*banks.BankTrxData) (err error) {
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
			fmt.Sprintf("[process.NewDB] Exec ImportBankTrx method to db (%d data)", len(data)),
			err,
		)
	}()

	if err != nil {
		return
	}

	var jsonData string
	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return json.Marshal(data)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				jsonData = string(i.([]byte))
			}

			return nil, nil
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return d.db.PrepareContext(c, QueryInsertTableBankTrx)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, i.(*sql.Stmt)).ExecContext( //nolint:sqlclosecheck
				c,
				jsonData,
			)
		},
	)

	return
}

func (d *DB) GenerateReconciliationMap(ctx context.Context, minAmount float64, maxAmount float64) (err error) {
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
			fmt.Sprintf("[process.NewDB] Exec GenerateReconciliationMap method to db (Amount %f - %f)", minAmount, maxAmount),
			err,
		)
	}()

	if err != nil {
		return
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return d.db.PrepareContext(c, QueryInsertTableReconciliationMap)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, i.(*sql.Stmt)).ExecContext( //nolint:sqlclosecheck
				c,
				minAmount,
				maxAmount,
			)
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
			"[process.NewDB] Exec Post method in db",
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
