package process

import (
	"context"
	"database/sql"
	"encoding/json"
	"simple-reconciliation-service/internal/pkg/reconcile/parser"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"strings"
	"time"

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
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				tx = i.(*sql.Tx)
			}

			return nil, nil
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableArguments == nil {
				return d.db.PrepareContext(c, QueryDropTableArguments)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableArguments = i.(*sql.Stmt)
			}

			return tx.StmtContext(c, d.stmtDropTableArguments).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableBanks == nil {
				return d.db.PrepareContext(c, QueryDropTableBanks)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableBanks = i.(*sql.Stmt)
			}

			return tx.StmtContext(c, d.stmtDropTableBanks).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableSystemTrx == nil {
				return d.db.PrepareContext(c, QueryDropTableSystemTrx)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableSystemTrx = i.(*sql.Stmt)
			}

			return tx.StmtContext(c, d.stmtDropTableSystemTrx).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableBankTrx == nil {
				return d.db.PrepareContext(c, QueryDropTableBankTrx)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableBankTrx = i.(*sql.Stmt)
			}

			return tx.StmtContext(c, d.stmtDropTableBankTrx).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableReconciliationMap == nil {
				return d.db.PrepareContext(c, QueryDropTableReconciliationMap)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableReconciliationMap = i.(*sql.Stmt)
			}

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
			"[sample.NewDB] Exec Pre method from db",
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
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtCreateTableArguments == nil {
				return d.db.PrepareContext(c, QueryCreateTableArguments)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtCreateTableArguments = i.(*sql.Stmt)
			}

			return tx.StmtContext(c, d.stmtCreateTableArguments).ExecContext( //nolint:sqlclosecheck
				c,
				startDate.Format("2006-01-02"),
				toDate.Format("2006-01-02"),
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtCreateTableBanks == nil {
				return d.db.PrepareContext(c, QueryCreateTableBanks)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtCreateTableBanks = i.(*sql.Stmt)
			}

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
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtCreateTableSystemTrx == nil {
				return d.db.PrepareContext(c, QueryCreateTableSystemTrx)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtCreateTableSystemTrx = i.(*sql.Stmt)
			}

			return tx.StmtContext(c, d.stmtCreateTableSystemTrx).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtCreateTableBankTrx == nil {
				return d.db.PrepareContext(c, QueryCreateTableBankTrx)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtCreateTableBankTrx = i.(*sql.Stmt)
			}

			return tx.StmtContext(c, d.stmtCreateTableBankTrx).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtCreateTableReconciliationMap == nil {
				return d.db.PrepareContext(c, QueryCreateTableReconciliationMap)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtCreateTableReconciliationMap = i.(*sql.Stmt)
			}

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
			"[sample.NewDB] Exec ImportSystemTrx method from db",
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
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtInsertTableSystemTrx == nil {
				return d.db.PrepareContext(c, QueryInsertTableSystemTrx)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtInsertTableSystemTrx = i.(*sql.Stmt)
			}

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
			"[sample.NewDB] Exec ImportBankTrx method from db",
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
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtInsertTableBankTrx == nil {
				return d.db.PrepareContext(c, QueryInsertTableBankTrx)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtInsertTableBankTrx = i.(*sql.Stmt)
			}

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

func (d *DB) GenerateReconciliationMap(ctx context.Context) (err error) {
	defer func() {
		log.Err(
			ctx,
			"[sample.NewDB] Exec ImportBankTrx method from db",
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
			if d.stmtInsertTableReconciliationMap == nil {
				return d.db.PrepareContext(c, QueryInsertTableReconciliationMap)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtInsertTableReconciliationMap = i.(*sql.Stmt)
			}

			return tx.StmtContext(c, d.stmtInsertTableReconciliationMap).ExecContext( //nolint:sqlclosecheck
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

func (d *DB) Post(ctx context.Context) (err error) {
	defer func() {
		log.Err(
			ctx,
			"[sample.NewDB] Exec Post method from db",
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
