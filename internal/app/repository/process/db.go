package process

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"strings"
	"time"

	"github.com/blockloop/scan/v2"

	"github.com/aaronjan/hunch"
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
	var executableInSequence []hunch.ExecutableInSequence
	stmtData := []string{
		QueryDropTableArguments,
		QueryDropTableBanks,
		QueryDropTableSystemTrx,
		QueryDropTableBankTrx,
		QueryDropTableReconciliationMap,
	}

	for k := range stmtData {
		executableInSequence = append(
			executableInSequence,
			func(c context.Context, _ interface{}) (interface{}, error) {
				i, e := d.db.PrepareContext(
					c,
					stmtData[k],
				)

				if e != nil {
					return nil, e
				}

				return tx.StmtContext(c, i).ExecContext( //nolint:sqlclosecheck
					c,
				)
			},
		)
	}

	_, err = hunch.Waterfall(
		ctx,
		executableInSequence...,
	)

	return
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
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return d.db.PrepareContext(c, QueryCreateTableArguments)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			dateStringFormat := "2006-01-02"
			return tx.StmtContext(c, i.(*sql.Stmt)).ExecContext( //nolint:sqlclosecheck
				c,
				startDate.Format(dateStringFormat),
				toDate.Format(dateStringFormat),
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return d.db.PrepareContext(c, QueryCreateTableBanks)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			b := new(strings.Builder)
			err := json.NewEncoder(b).Encode(listBank)
			if err != nil {
				return nil, err
			}

			return tx.StmtContext(c, i.(*sql.Stmt)).ExecContext( //nolint:sqlclosecheck
				c,
				b.String(),
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return d.db.PrepareContext(c, QueryCreateTableSystemTrx)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, i.(*sql.Stmt)).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return d.db.PrepareContext(c, QueryCreateTableBankTrx)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, i.(*sql.Stmt)).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return d.db.PrepareContext(c, QueryCreateTableReconciliationMap)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, i.(*sql.Stmt)).ExecContext( //nolint:sqlclosecheck
				c,
			)
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
