package process

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"simple-reconciliation-service/internal/pkg/reconcile/parser"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"strings"
	"time"

	"github.com/blockloop/scan/v2"

	"github.com/aaronjan/hunch"
	"github.com/pkg/errors"
)

type DB struct {
	db                               *sql.DB
	stmtDropTableArguments           *sql.Stmt
	stmtDropTableBanks               *sql.Stmt
	stmtDropTableSystemTrx           *sql.Stmt
	stmtDropTableBankTrx             *sql.Stmt
	stmtDropTableReconciliationMap   *sql.Stmt
	stmtCreateTableArguments         *sql.Stmt
	stmtCreateTableBanks             *sql.Stmt
	stmtCreateTableSystemTrx         *sql.Stmt
	stmtCreateTableBankTrx           *sql.Stmt
	stmtCreateTableReconciliationMap *sql.Stmt
	stmtInsertTableSystemTrx         *sql.Stmt
	stmtInsertTableBankTrx           *sql.Stmt
	stmtInsertTableReconciliationMap *sql.Stmt
	stmtGetReconciliationSummary     *sql.Stmt
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
	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtDropTableArguments == nil {
				d.stmtDropTableArguments, e = d.db.PrepareContext(c, QueryDropTableArguments)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtDropTableArguments).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtDropTableBanks == nil {
				d.stmtDropTableBanks, e = d.db.PrepareContext(c, QueryDropTableBanks)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtDropTableBanks).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtDropTableSystemTrx == nil {
				d.stmtDropTableSystemTrx, e = d.db.PrepareContext(c, QueryDropTableSystemTrx)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtDropTableSystemTrx).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtDropTableBankTrx == nil {
				d.stmtDropTableBankTrx, e = d.db.PrepareContext(c, QueryDropTableBankTrx)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtDropTableBankTrx).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtDropTableReconciliationMap == nil {
				d.stmtDropTableReconciliationMap, e = d.db.PrepareContext(c, QueryDropTableReconciliationMap)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtDropTableReconciliationMap).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
	)

	return
}

func (d *DB) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time) (err error) {
	defer func() {
		log.Err(
			ctx,
			"[process.NewDB] Exec Pre method from db",
			err,
		)
	}()

	tx, er := d.db.BeginTx(ctx, nil)
	if er != nil {
		return er
	}

	e := d.db.Ping()
	if e != nil {
		return e
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return nil, d.dropTables(c, tx)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtCreateTableArguments == nil {
				d.stmtCreateTableArguments, e = d.db.PrepareContext(c, QueryCreateTableArguments)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtCreateTableArguments).ExecContext( //nolint:sqlclosecheck
				c,
				startDate.Format("2006-01-02"),
				toDate.Format("2006-01-02"),
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtCreateTableBanks == nil {
				d.stmtCreateTableBanks, e = d.db.PrepareContext(c, QueryCreateTableBanks)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			b := new(strings.Builder)
			err := json.NewEncoder(b).Encode(listBank)
			if err != nil {
				return nil, err
			}

			return tx.StmtContext(c, d.stmtCreateTableBanks).ExecContext( //nolint:sqlclosecheck
				c,
				b.String(),
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtCreateTableSystemTrx == nil {
				d.stmtCreateTableSystemTrx, e = d.db.PrepareContext(c, QueryCreateTableSystemTrx)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtCreateTableSystemTrx).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtCreateTableBankTrx == nil {
				d.stmtCreateTableBankTrx, e = d.db.PrepareContext(c, QueryCreateTableBankTrx)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtCreateTableBankTrx).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtCreateTableReconciliationMap == nil {
				d.stmtCreateTableReconciliationMap, e = d.db.PrepareContext(c, QueryCreateTableReconciliationMap)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtCreateTableReconciliationMap).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
	)

	if err != nil {
		err = errors.Wrap(tx.Rollback(), err.Error())
	} else {
		err = tx.Commit()
	}

	return
}

func (d *DB) ImportSystemTrx(ctx context.Context, data []*parser.SystemTrxData) (err error) {
	defer func() {
		log.Err(
			ctx,
			fmt.Sprintf("[process.NewDB] ImportSystemTrx method from db (%d data)", len(data)),
			err,
		)
	}()

	tx, er := d.db.BeginTx(ctx, nil)
	if er != nil {
		return er
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
			if d.stmtInsertTableSystemTrx == nil {
				d.stmtInsertTableSystemTrx, e = d.db.PrepareContext(c, QueryInsertTableSystemTrx)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtInsertTableSystemTrx).ExecContext( //nolint:sqlclosecheck
				c,
				jsonData,
			)
		},
	)

	if err != nil {
		err = errors.Wrap(tx.Rollback(), err.Error())
	} else {
		err = tx.Commit()
	}

	return
}

func (d *DB) ImportBankTrx(ctx context.Context, data []*parser.BankTrxData) (err error) {
	defer func() {
		log.Err(
			ctx,
			fmt.Sprintf("[process.NewDB] Exec ImportBankTrx method to db (%d data)", len(data)),
			err,
		)
	}()

	tx, er := d.db.BeginTx(ctx, nil)
	if er != nil {
		return er
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
			if d.stmtInsertTableBankTrx == nil {
				d.stmtInsertTableBankTrx, e = d.db.PrepareContext(c, QueryInsertTableBankTrx)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtInsertTableBankTrx).ExecContext( //nolint:sqlclosecheck
				c,
				jsonData,
			)
		},
	)

	if err != nil {
		err = errors.Wrap(tx.Rollback(), err.Error())
	} else {
		err = tx.Commit()
	}

	return
}

func (d *DB) GenerateReconciliationMap(ctx context.Context, minAmount float64, maxAmount float64) (err error) {
	defer func() {
		log.Err(
			ctx,
			fmt.Sprintf("[process.NewDB] Exec GenerateReconciliationMap method to db (Amount %f - %f)", minAmount, maxAmount),
			err,
		)
	}()

	tx, er := d.db.BeginTx(ctx, nil)
	if er != nil {
		return er
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtInsertTableReconciliationMap == nil {
				d.stmtInsertTableReconciliationMap, e = d.db.PrepareContext(c, QueryInsertTableReconciliationMap)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtInsertTableReconciliationMap).ExecContext( //nolint:sqlclosecheck
				c,
				minAmount,
				maxAmount,
			)
		},
	)

	if err != nil {
		err = errors.Wrap(tx.Rollback(), err.Error())
	} else {
		err = tx.Commit()
	}

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
			if d.stmtGetReconciliationSummary == nil {
				d.stmtGetReconciliationSummary, e = d.db.PrepareContext(c, QueryGetReconciliationSummary)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return d.stmtGetReconciliationSummary.QueryContext( //nolint:sqlclosecheck
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
	defer func() {
		log.Err(
			ctx,
			"[process.NewDB] Exec Post method from db",
			err,
		)
	}()

	tx, er := d.db.BeginTx(ctx, nil)
	if er != nil {
		return er
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return nil, d.dropTables(c, tx)
		},
	)

	if err != nil {
		err = errors.Wrap(tx.Rollback(), err.Error())
	} else {
		err = tx.Commit()
	}

	return
}

func (d *DB) Close() (err error) {
	return d.db.Close()
}
